// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	ctx "context"
	"slices"
	"strconv"
	"strings"

	core "typefox.dev/fastbelt"
	internalATN "typefox.dev/fastbelt/internal/atn"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/codegen"
)

// parserATNData holds all ATN-derived maps shared across the three parser
// generators. Build it once with BuildParserATNData and pass it to each
// generator, so the ATN is only constructed once per generation run.
type parserATNData struct {
	followStateName   map[core.AstNode]string                   // grammar.RuleCall -> follow-state constant name
	decisionStateName map[core.AstNode]string                   // grammar.Element -> decision-state constant name
	elementStateNames map[grammar.Element]string                // token-consuming element -> ATN state constant name
	ruleStartName     map[core.AstNode]string                   // rule -> RuleStart state constant name (completion only)
	orAdaptive        map[core.AstNode]int                      // grammar.Alternatives -> ALL(*) decision index
	loopAdaptive      map[core.AstNode]int                      // grammar.Element -> ALL(*) decision index
	orDecision        map[grammar.Element]*internalATN.ATNState // grammar.Alternatives -> alternative-choice state
	loopDecision      map[grammar.Element]*internalATN.ATNState // grammar.Element -> cardinality-guard state
	tokenVarNames     []string                                  // ATN token id (0-based) -> generated token var name
}

// BuildParserATNData builds the ATN and all derived name/decision maps used by
// the three parser generators. Call it once and pass the result to
// GenerateParser, GenerateCompletionParser, and GenerateParserLookahead.
// Returns nil when the ATN cannot be built (invalid grammar).
func BuildParserATNData(grammr grammar.Grammar, tokenTypes GenerateTokenTypesResult) *parserATNData {
	builtATN, _ := internalATN.CreateATN(grammr, tokenTypes.TokenTypeIds)
	if builtATN == nil {
		return nil
	}
	elementNames := internalATN.BuildElementNames(grammr)
	stateNames := internalATN.BuildStateNameMap(builtATN, elementNames)
	stateIdx := internalATN.BuildStateIndexMap(builtATN)
	orAdaptive, loopAdaptive := buildAdaptiveDecisionMaps(builtATN)
	return &parserATNData{
		followStateName:   buildFollowStateNameMap(builtATN, stateIdx, stateNames),
		decisionStateName: buildDecisionStateNameMap(builtATN, stateIdx, stateNames),
		elementStateNames: buildElementStateNameMap(builtATN, stateNames),
		ruleStartName:     buildRuleStartNameMap(builtATN, stateIdx, stateNames),
		orAdaptive:        orAdaptive,
		loopAdaptive:      loopAdaptive,
		orDecision:        builtATN.OrDecision,
		loopDecision:      builtATN.LoopDecision,
		tokenVarNames:     tokenTypes.TokenTypeVarNames,
	}
}

// followState returns the ATN follow-state constant name for a rule call site,
// or "0" if the ATN was not built or the call site is not found.
func (d *parserATNData) followState(ruleCall core.AstNode) string {
	if d == nil {
		return "0"
	}
	if name, ok := d.followStateName[ruleCall]; ok {
		return name
	}
	return "0"
}

// elementStateName returns the ATN state constant name for a token-consuming
// element, used in AssignToken calls. Returns "" when the ATN is unavailable.
func (d *parserATNData) elementStateName(el grammar.Element) string {
	if d == nil {
		return ""
	}
	return d.elementStateNames[el]
}

// ruleStart returns the ATN RuleStart state constant name for a rule, or "0"
// when the ATN is unavailable.
func (d *parserATNData) ruleStart(rule core.AstNode) string {
	if d == nil {
		return "0"
	}
	if name, ok := d.ruleStartName[rule]; ok {
		return name
	}
	return "0"
}

// ParserGeneratorContext carries the grammar, per-decision lookahead values,
// and generator-mode flags shared by all emit functions.
type ParserGeneratorContext struct {
	grammar         grammar.Grammar
	lookaheads      map[core.AstNode]LookaheadValue // cardinality-guard decisions
	orLookaheads    map[core.AstNode]LookaheadValue // alternative-selection decisions
	inCompositeRule bool
	atnData         *parserATNData // nil when ATN not built (degenerate grammar)
	// completion distinguishes parser_gen.go (false) from completion_parser_gen.go
	// (true). In completion mode, AST-construction emits are suppressed and
	// CompletionParserState bookkeeping (EnterRule, RecordSnapshot,
	// MarkAssignment) is inserted instead.
	completion bool
	// counter for unique Go loop labels (loop0, loop1, ...)
	// Only used when the generator can combine an Alternatives node with its loop
	loopLabelSeq int
}

type LookaheadValue struct {
	name   string
	lookup LL1Decision
}

// canCombineAltsWithLoop reports whether an Alternatives node can have its loop
// guard merged with its OR switch. When combined, the OR lookahead drives both
// the entry/continuation decision and alternative selection with a single call,
// eliminating the separate guard variable. Requires both decisions to be static
// LL(1) — adaptive decisions cannot share a switch without a separate bool result.
func (ctx *ParserGeneratorContext) canCombineAltsWithLoop(alts grammar.Alternatives) bool {
	if len(alts.Alts()) <= 1 {
		return false
	}
	if ctx.atnData == nil {
		return false
	}
	_, loopIsAdaptive := ctx.atnData.loopAdaptive[alts]
	_, orIsAdaptive := ctx.atnData.orAdaptive[alts]
	return !loopIsAdaptive && !orIsAdaptive
}

// canCombineAlternatives is canCombineAltsWithLoop restricted to Alternatives
// nodes whose own cardinality is non-One (so the loop guard is on the Alternatives
// itself, not on a wrapping Assignment or Group).
func (ctx *ParserGeneratorContext) canCombineAlternatives(alts grammar.Alternatives) bool {
	return alts.Cardinality() != CardinalityOne && ctx.canCombineAltsWithLoop(alts)
}

func (ctx *ParserGeneratorContext) nextLoopLabel() string {
	label := "loop" + strconv.Itoa(ctx.loopLabelSeq)
	ctx.loopLabelSeq++
	return label
}

// GenerateParserLookahead emits parser_lookahead_gen.go — the <Grammar>ParserLookahead
// service. The service exposes one method per lookahead/prediction decision the
// generated parser performs, so the prediction strategy can be swapped out (e.g.
// for semantic or context-sensitive prediction) without touching the parser.
//
// The default implementation delegates each decision to the same primitive the
// inline parser uses: static LL(1) (Lookahead) for decisions a single token
// disambiguates, ALL(*) (AdaptivePredict) for the rest. Because the generated
// file shares a package with parser_gen.go, it reuses the lookahead tables and
// ATN-decision indices defined there.
func GenerateParserLookahead(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult, atnData *parserATNData) string {
	context := &ParserGeneratorContext{
		grammar:      grammr,
		lookaheads:   make(map[core.AstNode]LookaheadValue),
		orLookaheads: make(map[core.AstNode]LookaheadValue),
		atnData:      atnData,
	}
	populateContext(context)
	methods := generateLookaheadMethods(context)

	interfaceName := grammr.Name() + "ParserLookahead"
	defaultName := "Default" + interfaceName

	signature := func(m lookaheadMethod) string {
		if m.isOr {
			return m.name + "(state *parser.ParserState) (int, *parser.PredictionFailure)"
		}
		return m.name + "(state *parser.ParserState) bool"
	}

	// core is needed only when at least one static LL(1) table is emitted.
	needsCoreImport := false
	for _, m := range methods {
		if m.table != nil {
			needsCoreImport = true
			break
		}
	}

	// Build the maps once; generateLL1Lookahead needs them per table.
	varNameToId := buildVarNameToId(tokenTypes)
	groupMembers := buildGroupVarNameToMembers(grammr)

	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		if needsCoreImport {
			n.AppendLine("core \"typefox.dev/fastbelt\"")
		}
		n.AppendLine("\"typefox.dev/fastbelt/parser\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	// Generate named constants for every ALL(*) decision index
	hasAdaptive := false
	for _, m := range methods {
		if m.adaptiveIdx != nil {
			hasAdaptive = true
			break
		}
	}
	if hasAdaptive {
		node.AppendLine("const (")
		node.Indent(func(n codegen.Node) {
			for _, m := range methods {
				if m.adaptiveIdx != nil {
					n.AppendLine(decisionConstName(m.name), " = ", strconv.Itoa(*m.adaptiveIdx))
				}
			}
		})
		node.AppendLine(")")
		node.AppendLine()
	}

	// Static LL(1) lookup tables, one per decision that uses the static path.
	for _, m := range methods {
		if m.table != nil {
			generateLL1Lookahead(node, m.name, *m.table, varNameToId, groupMembers)
			node.AppendLine()
		}
	}

	node.AppendLine("// ", interfaceName, " abstracts every lookahead/prediction decision performed by")
	node.AppendLine("// the generated parser. Each method corresponds to a single decision point;")
	node.AppendLine("// implementations may override individual decisions while delegating the rest")
	node.AppendLine("// to ", defaultName, ".")
	node.AppendLine("type ", interfaceName, " interface {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("parser.ParserLookahead").AppendLine()
		for _, m := range methods {
			n.AppendLine(signature(m))
		}
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("// ", defaultName, " resolves every decision with the parser state's built-in")
	node.AppendLine("// static LL(1) and adaptive ALL(*) prediction.")
	node.AppendLine("type ", defaultName, " struct {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("parser.DefaultParserLookahead")
	})
	node.AppendLine("}")
	node.AppendLine()
	node.AppendLine("func New", defaultName, "() ", interfaceName, " {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return &", defaultName, "{}")
	})
	node.AppendLine("}")
	node.AppendLine()

	for _, m := range methods {
		node.AppendLine("func (l *", defaultName, ") ", signature(m), " {")
		node.Indent(func(n codegen.Node) {
			for _, line := range m.body {
				n.AppendLine(line)
			}
		})
		node.AppendLine("}")
		node.AppendLine()
	}

	return FormatIfPossible(node.String())
}

// lookaheadMethod represents one decision point in the generated service.
// OR decisions (alternative selection) return both the chosen alternative
// and a failure descriptor; guard decisions (cardinality enter/continue)
// return a bool. table is the static LL(1) table to emit as a package var,
// or nil when the decision uses adaptive prediction or a direct keyword check.
type lookaheadMethod struct {
	name        string
	body        []string
	isOr        bool
	table       *LL1Decision
	adaptiveIdx *int
}

func generateLookaheadMethods(context *ParserGeneratorContext) []lookaheadMethod {
	methods := make([]lookaheadMethod, 0, len(context.orLookaheads)+len(context.lookaheads))

	for alts, lv := range context.orLookaheads {
		m := lookaheadMethod{name: lv.name, isOr: true}
		if idx, ok := context.atnData.orAdaptive[alts]; ok {
			idx := idx
			m.adaptiveIdx = &idx
			m.body = []string{"return state.AdaptivePredict(" + decisionConstName(m.name) + ", l.PredictionMode())"}
		} else {
			table := lv.lookup
			m.table = &table
			m.body = []string{"return state.Lookahead(" + lv.name + ")"}
		}
		methods = append(methods, m)
	}

	// Token-group var names need the membership-aware Matches check; single
	// token types (keywords, terminals) can use a direct identity comparison.
	groupNames := make(map[string]bool)
	for _, tg := range context.grammar.TokenGroups() {
		groupNames[GeneratedTokenName(tg)] = true
	}

	for el, lv := range context.lookaheads {
		// The adaptive decision map is keyed by the assignable value for
		// assignments (what the ATN actually sees), not the Assignment wrapper.
		adaptiveKey := el
		if assignment, ok := el.(grammar.Assignment); ok {
			adaptiveKey = assignment.Value()
		}
		m := lookaheadMethod{name: lv.name, isOr: false}
		if idx, ok := context.atnData.loopAdaptive[adaptiveKey]; ok {
			// Complex, non-LL(1) decision that requires adaptive prediction
			idx := idx
			m.adaptiveIdx = &idx
			m.body = []string{
				"prediction, _ := state.AdaptivePredict(" + decisionConstName(m.name) + ", l.PredictionMode())",
				"return prediction == 0",
			}
		} else if token, ok := singleTokenGuard(lv.lookup); ok {
			// The guard is a true/false decision whose FIRST set is a single token
			// type, so a one-token check decides it - skip the lookup table entirely.
			// This is an optimization that applies to cardinality guards ONLY.
			if groupNames[token] {
				// Call Matches on token groups, will internally check all member types
				// using a bitset lookup.
				m.body = []string{"return " + token + ".Matches(state.LA(1).Type)"}
			} else {
				// Keyword or simple token type - direct comparison is the fastest way
				m.body = []string{"return state.LA(1).Type == " + token}
			}
		} else {
			// Requires a full lookup table, since multiple token types are involved in the decision
			table := lv.lookup
			m.table = &table
			m.body = []string{
				"prediction, _ := state.Lookahead(" + lv.name + ")",
				"return prediction == 0",
			}
		}
		methods = append(methods, m)
	}

	slices.SortFunc(methods, func(a, b lookaheadMethod) int {
		return strings.Compare(a.name, b.name)
	})
	return methods
}

// singleTokenGuard reports whether a cardinality guard's LL(1) table reduces to
// exactly one token type and, if so, returns its generated var name. A guard
// table always holds a single option (the enter FIRST set); when that option
// names just one token type, a one-token check fully decides the guard and the
// lookup table can be skipped.
func singleTokenGuard(lookup LL1Decision) (string, bool) {
	if len(lookup) == 1 && len(lookup[0]) == 1 {
		return lookup[0][0], true
	}
	return "", false
}

func decisionConstName(methodName string) string {
	return "Decision" + methodName
}

func GenerateParser(grammr grammar.Grammar, entryRule grammar.ParserRule, packageName string, tokenTypes GenerateTokenTypesResult, atnData *parserATNData) string {
	context := &ParserGeneratorContext{
		grammar:      grammr,
		lookaheads:   make(map[core.AstNode]LookaheadValue),
		orLookaheads: make(map[core.AstNode]LookaheadValue),
		atnData:      atnData,
	}

	populateContext(context)
	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/parser\"")
		n.AppendLine("\"typefox.dev/fastbelt/util/service\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	node.AppendLine("type Parser struct {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("state *parser.ParserState")
		n.AppendLine("sc *service.Container")
		n.AppendLine("referencesConstructor ", grammr.Name(), "ReferencesConstructor")
		n.AppendLine("lookahead ", grammr.Name(), "ParserLookahead")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (p *Parser) Parse(document *core.Document) *parser.ParseResult {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("recovery := service.MustGet[parser.ErrorRecoveryStrategy](p.sc)")
		n.AppendLine("messages := service.MustGet[parser.ErrorMessageProvider](p.sc)")
		n.AppendLine("referencesConstructor := service.MustGet[", grammr.Name(), "ReferencesConstructor](p.sc)")
		n.AppendLine("lookahead := service.MustGet[", grammr.Name(), "ParserLookahead](p.sc)")
		n.AppendLine("cp := &Parser{sc: p.sc, referencesConstructor: referencesConstructor, lookahead: lookahead, state: parser.NewParserState(document.Tokens, ATN(), recovery, messages)}")
		n.AppendLine("result := cp.Parse", entryRule.Name(), "()")
		n.AppendLine("cp.state.ExpectEndOfInput()")
		n.AppendLine("core.AssignContainers(document, result)")
		n.AppendLine("return &parser.ParseResult{Node: result, Errors: cp.state.Errors()}")
	})
	node.AppendLine("}")
	node.AppendLine()
	node.AppendLine("func NewParser(sc *service.Container) *Parser {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return &Parser{")
		n.AppendLine("	sc: sc,")
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	for _, rule := range grammr.Rules() {
		generateParseFunction(node, context, rule)
	}
	for _, composite := range grammr.Composites() {
		generateCompositeParseFunction(node, context, composite)
	}

	return FormatIfPossible(node.String())
}

// GenerateCompletionParser emits completion_parser_gen.go — a peer of
// parser_gen.go that mirrors the main parser's control flow but skips every
// AST-mutation call and instead records the CompletionParserState bookkeeping
// (rule stack, ATN snapshots, assignment markers) the completion provider needs.
//
// The generated file reuses the lookahead tables and ATN builder defined by
// GenerateParser/EmitGoSource, so it must be emitted into the same package.
func GenerateCompletionParser(grammr grammar.Grammar, entryRule grammar.ParserRule, packageName string, tokenTypes GenerateTokenTypesResult, atnData *parserATNData) string {
	context := &ParserGeneratorContext{
		grammar:      grammr,
		lookaheads:   make(map[core.AstNode]LookaheadValue),
		orLookaheads: make(map[core.AstNode]LookaheadValue),
		atnData:      atnData,
		completion:   true,
	}

	populateContext(context)
	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("\"sync\"")
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/parser\"")
		n.AppendLine("\"typefox.dev/fastbelt/util/service\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	node.AppendLine("type CompletionParser struct {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("state *parser.ParserState")
		n.AppendLine("cp *parser.CompletionParserState")
		n.AppendLine("sc *service.Container")
		n.AppendLine("atn func() *parser.RuntimeATN")
		n.AppendLine("lookahead ", grammr.Name(), "ParserLookahead")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func NewCompletionParser(sc *service.Container) *CompletionParser {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return &CompletionParser{")
		n.AppendLine("	sc: sc,")
		n.AppendLine("	atn: sync.OnceValue(BuildATN),")
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("// Parse runs the completion parser over the given prefix tokens (typically")
	node.AppendLine("// the document's tokens up to the cursor) and returns the recorded")
	node.AppendLine("// snapshots and rule stack. The completion provider feeds that result into")
	node.AppendLine("// the ATN simulator.")
	node.AppendLine("func (p *CompletionParser) Parse(tokens []core.Token) *parser.CompletionParseResult {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("messages := service.MustGet[parser.ErrorMessageProvider](p.sc)")
		n.AppendLine("recovery := service.MustGet[parser.ErrorRecoveryStrategy](p.sc)")
		n.AppendLine("lookahead := service.MustGet[", grammr.Name(), "ParserLookahead](p.sc)")
		n.AppendLine("cp := &CompletionParser{sc: p.sc, atn: p.atn, lookahead: lookahead}")
		n.AppendLine("cp.state = parser.NewParserState(tokens, p.atn(), recovery, messages)")
		n.AppendLine("cp.cp = parser.NewCompletionParserState(cp.state)")
		n.AppendLine("cp.Parse", entryRule.Name(), "()")
		n.AppendLine("return cp.cp.Result(tokens)")
	})
	node.AppendLine("}")
	node.AppendLine()

	for _, rule := range grammr.Rules() {
		generateParseFunction(node, context, rule)
	}
	for _, composite := range grammr.Composites() {
		generateCompositeParseFunction(node, context, composite)
	}

	return FormatIfPossible(node.String())
}

// buildFollowStateNameMap maps each grammar.RuleCall to the constant name of
// the ATN follow state for that call site.
func buildFollowStateNameMap(builtATN *internalATN.ATN, stateIdx map[*internalATN.ATNState]int, stateNames []string) map[core.AstNode]string {
	result := make(map[core.AstNode]string)
	for _, state := range builtATN.States {
		if state.RuleCallEntry == nil {
			continue
		}
		for _, t := range state.Transitions {
			if rt, ok := t.(*internalATN.RuleTransition); ok {
				result[state.RuleCallEntry] = stateNames[stateIdx[rt.FollowState]]
				break
			}
		}
	}
	return result
}

// buildDecisionStateNameMap maps each grammar.Element to the constant name of
// its ATN decision state.
func buildDecisionStateNameMap(builtATN *internalATN.ATN, stateIdx map[*internalATN.ATNState]int, stateNames []string) map[core.AstNode]string {
	result := make(map[core.AstNode]string, len(builtATN.DecisionMap))
	for _, state := range builtATN.DecisionMap {
		if state.Production != nil {
			result[state.Production] = stateNames[stateIdx[state]]
		}
	}
	return result
}

// buildAdaptiveDecisionMaps decides, per decision in the built ATN, whether to
// emit ALL(*) prediction (AdaptivePredict) instead of the static LL(1) Lookahead.
// The static fast path is kept wherever a single token already disambiguates;
// everything else routes through the ATN predictor.
//
//   - Alternative choices ("or" decisions) use adaptive prediction unless LL(1)
//     uniquely separates the alternatives. Decisions whose transition count does
//     not match the alternative count are excluded (e.g. `(a|b)?` has a skip
//     edge that pollutes the state) — those keep the static path.
//   - Cardinality guards (enter-vs-exit) use adaptive prediction when the enter
//     and exit FIRST sets overlap, i.e. a single token cannot decide. The exit
//     set is the loop's full FOLLOW, crossing rule-stop boundaries so the guard
//     is correct even when the loop sits at the end of its rule.
func buildAdaptiveDecisionMaps(a *internalATN.ATN) (orMap, loopMap map[core.AstNode]int) {
	orMap = map[core.AstNode]int{}
	for el, ds := range a.OrDecision {
		alts, ok := el.(grammar.Alternatives)
		if !ok || len(alts.Alts()) < 2 {
			continue
		}
		if len(ds.Transitions) != len(alts.Alts()) {
			continue // skip edge present (e.g. `(a|b)?`) — not a clean 1:1 mapping
		}
		if ll1UniqueDecision(ds) {
			continue
		}
		orMap[el] = ds.Decision
	}
	// Build callers so loopNeedsAdaptive can compute a loop's true (cross-rule)
	// FOLLOW set. A loop at the end of its rule depends on the enclosing context,
	// which a static intra-rule lookahead cannot see.
	callers := map[grammar.AbstractRuleWithBody][]*internalATN.ATNState{}
	for _, state := range a.States {
		for _, t := range state.Transitions {
			if rt, ok := t.(*internalATN.RuleTransition); ok {
				callers[rt.Rule] = append(callers[rt.Rule], rt.FollowState)
			}
		}
	}
	loopMap = map[core.AstNode]int{}
	for el, ds := range a.LoopDecision {
		if loopNeedsAdaptive(ds, callers) {
			loopMap[el] = ds.Decision
		}
	}
	return orMap, loopMap
}

// ll1UniqueDecision reports whether a single token uniquely identifies each
// outgoing transition of a decision state (the LL(1) fast-path condition for an
// "or" decision). It holds when the per-transition FIRST sets are pairwise
// disjoint, so one lookahead token always selects exactly one alternative.
func ll1UniqueDecision(ds *internalATN.ATNState) bool {
	seen := map[int]bool{}
	for _, t := range ds.Transitions {
		for id := range firstSetFromATNState(t.Target()) {
			if seen[id] {
				return false
			}
			seen[id] = true
		}
	}
	return true
}

// loopNeedsAdaptive reports whether a cardinality guard requires adaptive
// prediction. It does when the enter set (tokens that continue the loop) and
// the exit FOLLOW set (tokens that can appear after the loop) overlap — meaning
// a single lookahead token cannot decide. Only clean two-transition states
// (enter at [0], exit at [1]) are considered.
func loopNeedsAdaptive(ds *internalATN.ATNState, callers map[grammar.AbstractRuleWithBody][]*internalATN.ATNState) bool {
	if len(ds.Transitions) != 2 {
		return false
	}
	enter := firstSetFromATNState(ds.Transitions[0].Target())
	exit := followAwareFirstSet(ds.Transitions[1].Target(), callers)
	for id := range enter {
		if exit[id] {
			return true
		}
	}
	return false
}

// firstSetFromATNState returns the set of token-type IDs reachable first from
// state, following epsilon and rule (descend into callee) edges but not crossing
// rule-stop boundaries.
func firstSetFromATNState(state *internalATN.ATNState) map[int]bool {
	result := map[int]bool{}
	visited := map[*internalATN.ATNState]bool{}
	queue := []*internalATN.ATNState{state}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil || visited[cur] {
			continue
		}
		visited[cur] = true
		for _, t := range cur.Transitions {
			switch tt := t.(type) {
			case *internalATN.AtomTransition:
				result[tt.TokenTypeId] = true
			case *internalATN.EpsilonTransition:
				queue = append(queue, tt.Target())
			case *internalATN.RuleTransition:
				queue = append(queue, tt.Target())
			}
		}
	}
	return result
}

// followAwareFirstSet extends firstSetFromATNState with cross-rule FOLLOW: when
// the traversal reaches a rule-stop state it continues into the follow states of
// every caller of that rule (callers), so the result includes tokens that can
// appear after the current rule returns. The shared visited set guards against
// infinite recursion through recursive or transitive follows.
func followAwareFirstSet(state *internalATN.ATNState, callers map[grammar.AbstractRuleWithBody][]*internalATN.ATNState) map[int]bool {
	result := map[int]bool{}
	visited := map[*internalATN.ATNState]bool{}
	queue := []*internalATN.ATNState{state}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == nil || visited[cur] {
			continue
		}
		visited[cur] = true
		if cur.Type == parser.ATNRuleStop {
			// Fell off the end of cur.Rule: continue into every caller's follow.
			queue = append(queue, callers[cur.Rule]...)
			continue
		}
		for _, t := range cur.Transitions {
			switch tt := t.(type) {
			case *internalATN.AtomTransition:
				result[tt.TokenTypeId] = true
			case *internalATN.EpsilonTransition:
				queue = append(queue, tt.Target())
			case *internalATN.RuleTransition:
				queue = append(queue, tt.Target())
			}
		}
	}
	return result
}

// orSwitchHead emits the switch header for an alternatives decision with a
// default arm that can use the failure descriptor.
func orSwitchHead(context *ParserGeneratorContext, alts core.AstNode) string {
	return "switch prediction, failure := p.lookahead." + context.orLookaheads[alts].name + "(p.state); prediction {"
}

// orSwitchHeadAlt emits the switch header when no default arm needs the failure
// descriptor (discarding it avoids an unused-variable compiler error).
func orSwitchHeadAlt(context *ParserGeneratorContext, alts core.AstNode) string {
	return "switch prediction, _ := p.lookahead." + context.orLookaheads[alts].name + "(p.state); prediction {"
}

// guardCall returns the boolean entry/continuation guard expression for a
// cardinality decision.
func guardCall(context *ParserGeneratorContext, element core.AstNode) string {
	return "p.lookahead." + context.lookaheads[element].name + "(p.state)"
}

// buildRuleStartNameMap maps each parser/composite rule to the constant name of
// its ATN RuleStart state. Used in completion mode for EnterRule calls.
func buildRuleStartNameMap(builtATN *internalATN.ATN, stateIdx map[*internalATN.ATNState]int, stateNames []string) map[core.AstNode]string {
	result := make(map[core.AstNode]string, len(builtATN.RuleToStartState))
	for rule, state := range builtATN.RuleToStartState {
		result[rule] = stateNames[stateIdx[state]]
	}
	return result
}

// buildElementStateNameMap maps token-consuming grammar elements to their ATN
// state constant name. Used in AssignToken calls to identify which grammar
// position produced the token.
func buildElementStateNameMap(builtATN *internalATN.ATN, stateNames []string) map[grammar.Element]string {
	result := make(map[grammar.Element]string)
	for i, s := range builtATN.States {
		if s.ConsumedElement != nil {
			result[s.ConsumedElement] = stateNames[i]
		}
	}
	return result
}

// buildSyncCall returns the Go statement to emit for a Sync call at the given
// element's decision state, or "" if no decision state is known.
//
// In completion mode the statement also records an ATN snapshot at the same
// state index, giving the completion provider a starting point for the simulator
// at every branch boundary.
func buildSyncCall(context *ParserGeneratorContext, el core.AstNode) string {
	if context.atnData == nil {
		return ""
	}
	name, ok := context.atnData.decisionStateName[el]
	if !ok {
		return ""
	}
	if context.completion {
		return "p.cp.RecordSnapshot(" + name + "); p.state.Sync(" + name + ")"
	}
	return "p.state.Sync(" + name + ")"
}

// lookaheadDecl is one prediction decision discovered while walking the grammar,
// before its final (collision-free) name is assigned.
type lookaheadDecl struct {
	node   core.AstNode
	isOr   bool
	base   string
	lookup LL1Decision
}

func populateContext(context *ParserGeneratorContext) {
	var decls []lookaheadDecl
	for _, rule := range context.grammar.Rules() {
		decls = collectLookaheadDecls(context, decls, rule.Name(), rule.Body())
	}
	for _, composite := range context.grammar.Composites() {
		decls = collectLookaheadDecls(context, decls, composite.Name(), composite.Body())
	}

	// Names are derived structurally from the rule/property path plus the kind of
	// decision (e.g. ObjAlternatives, QualifiedNameLoop). A numeric suffix is
	// appended only when two decisions share a base name, mirroring the ATN
	// state-name scheme.
	counts := make(map[string]int, len(decls))
	for _, d := range decls {
		counts[d.base]++
	}
	seen := make(map[string]int, len(decls))
	for _, d := range decls {
		name := d.base
		if counts[d.base] > 1 {
			name = d.base + "_" + strconv.Itoa(seen[d.base])
			seen[d.base]++
		}
		if d.isOr {
			context.orLookaheads[d.node] = LookaheadValue{name: name, lookup: d.lookup}
		} else {
			context.lookaheads[d.node] = LookaheadValue{name: name, lookup: d.lookup}
		}
	}
}

// cardinalityWord names a cardinality guard by what it does.
func cardinalityWord(cardinality string) string {
	if cardinality == CardinalityOptional {
		return "Optional"
	}
	return "Loop"
}

func collectLookaheadDecls(context *ParserGeneratorContext, decls []lookaheadDecl, prefix string, node core.AstNode) []lookaheadDecl {
	switch n := node.(type) {
	case grammar.Alternatives:
		if len(n.Alts()) > 1 {
			decls = append(decls, lookaheadDecl{node: n, isOr: true, base: prefix + "Alternatives", lookup: ll1DecisionOr(context.atnData, n)})
		}
		// When the OR lookahead can drive both entry/continuation and alternative
		// selection, skip the separate guard variable entirely — the OR switch's
		// default case serves as the exit.
		if n.Cardinality() != CardinalityOne && !context.canCombineAlternatives(n) {
			decls = append(decls, lookaheadDecl{node: n, base: prefix + cardinalityWord(n.Cardinality()), lookup: ll1DecisionOpt(context.atnData, n)})
		}
		for _, alt := range n.Alts() {
			decls = collectLookaheadDecls(context, decls, prefix, alt)
		}
	case grammar.Group:
		if n.Cardinality() != CardinalityOne {
			decls = append(decls, lookaheadDecl{node: n, base: prefix + cardinalityWord(n.Cardinality()), lookup: ll1DecisionOpt(context.atnData, n)})
		}
		for _, element := range n.Elements() {
			decls = collectLookaheadDecls(context, decls, prefix, element)
		}
	case grammar.Keyword:
		if n.Cardinality() != CardinalityOne {
			decls = append(decls, lookaheadDecl{node: n, base: prefix + grammar.KeywordName(n) + cardinalityWord(n.Cardinality()), lookup: ll1DecisionOpt(context.atnData, n)})
		}
	case grammar.Assignment:
		if n.Cardinality() != CardinalityOne {
			// Skip the guard variable when the value is a combinable Alternatives:
			// the OR lookahead alone will drive both entry and selection.
			alts, valueIsAlts := n.Value().(grammar.Alternatives)
			if !valueIsAlts || !context.canCombineAltsWithLoop(alts) {
				decls = append(decls, lookaheadDecl{node: n, base: prefix + n.Property().Text() + cardinalityWord(n.Cardinality()), lookup: ll1DecisionOpt(context.atnData, n)})
			}
		}
		decls = collectLookaheadDecls(context, decls, prefix+n.Property().Text(), n.Value())
	case grammar.CrossRef:
		decls = collectLookaheadDecls(context, decls, prefix, n.Rule())
	case grammar.RuleCall:
		if n.Cardinality() != CardinalityOne {
			decls = append(decls, lookaheadDecl{node: n, base: prefix + cardinalityWord(n.Cardinality()), lookup: ll1DecisionOpt(context.atnData, n)})
		}
	}
	return decls
}

func generateParseFunction(node codegen.Node, context *ParserGeneratorContext, rule grammar.ParserRule) {
	returnType := grammar.FindReturnType(rule, ctx.Background())
	if returnType == nil {
		panic("Unable to find return type for rule: " + rule.Name())
	}
	context.inCompositeRule = false
	receiverType := context.receiverType()
	if context.completion {
		node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "() {")
		node.Indent(func(n codegen.Node) {
			n.AppendLine("p.cp.EnterRule(", strconv.Quote(rule.Name()), ", ", context.atnData.ruleStart(rule), ")")
			n.AppendLine("defer p.cp.ExitRule()")
			generateAbstractElementParser(n, context, rule.Body())
		})
		node.AppendLine("}")
		node.AppendLine()
	} else {
		node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "() ", returnType.Name(), " {")
		node.Indent(func(n codegen.Node) {
			n.AppendLine("current := New", returnType.Name(), "()")
			n.AppendLine("current.SetSegmentStartToken(p.state.LA(1))")
			// Generate new lexical scope for actions that immediately trigger on rule start
			n.AppendLine("{")
			generateAbstractElementParser(n, context, rule.Body())
			n.AppendLine("}")
			n.AppendLine("current.SetSegmentEndToken(p.state.LA(0))")
			n.AppendLine("return current")
		})
		node.AppendLine("}")
		node.AppendLine()
	}
}

func generateCompositeParseFunction(node codegen.Node, context *ParserGeneratorContext, rule grammar.CompositeRule) {
	context.inCompositeRule = true
	receiverType := context.receiverType()
	if context.completion {
		node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "() {")
		node.Indent(func(n codegen.Node) {
			n.AppendLine("p.cp.EnterRule(", strconv.Quote(rule.Name()), ", ", context.atnData.ruleStart(rule), ")")
			n.AppendLine("defer p.cp.ExitRule()")
			generateAbstractElementParser(n, context, rule.Body())
		})
		node.AppendLine("}")
		node.AppendLine()
	} else {
		node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "(current core.CompositeNode) {")
		node.Indent(func(n codegen.Node) {
			generateAbstractElementParser(n, context, rule.Body())
		})
		node.AppendLine("}")
		node.AppendLine()
	}
}

// receiverType is the Go struct name that owns the generated Parse methods.
func (c *ParserGeneratorContext) receiverType() string {
	if c.completion {
		return "CompletionParser"
	}
	return "Parser"
}

func generateAbstractElementParser(node codegen.Node, context *ParserGeneratorContext, element grammar.Element) {
	switch e := element.(type) {
	case grammar.Alternatives:
		generateAlternativesParser(node, context, e)
	case grammar.Group:
		generateGroupParser(node, context, e)
	case grammar.Action:
		if context.completion {
			// Actions only mutate AST structure (re-anchoring current); they
			// consume no tokens and have no ATN impact, so completion parser
			// emits nothing for them.
			return
		}
		node.AppendLine("{")
		node.Indent(func(n codegen.Node) {
			n.AppendLine("result := New", e.Type().Text(), "()")
			// Inherit segment from previous node
			n.AppendLine("result.SetSegment(current.Segment())")
			if e.Property() != nil {
				if e.Operator() == "+=" {
					n.AppendLine("result.Set", e.Property().Text(), "Item(current)")
				} else {
					n.AppendLine("result.Set", e.Property().Text(), "(current)")
				}
				// Ensure that the previous node has a valid segment ending
				n.AppendLine("current.SetSegmentEndToken(p.state.LA(0))")
				n.AppendLine("current = result")
			} else {
				// If there is no property to assign, just merge tokens
				n.AppendLine("core.AssignTokens(result, current.Tokens())")
				n.AppendLine("current = result")
			}
		})
		node.AppendLine("}")
		node.AppendLine("current := current.(", e.Type().Text(), ")")
	case grammar.Keyword:
		node.AppendLine("{")
		node.Indent(func(indent codegen.Node) {
			generateKeywordParser(indent, context, e)
		})
		node.AppendLine("}")
	case grammar.RuleCall:
		node.AppendLine("{")
		node.Indent(func(indent codegen.Node) {
			resultName := generateRuleCallParser(indent, context, e)
			if context.completion {
				return
			}
			if resultName == "result" && !context.inCompositeRule {
				// Unassigned rule call result — merge into current node
				indent.AppendLine("core.MergeTokens(result, current.Tokens())")
				indent.AppendLine("current = result")
			}
		})
		node.AppendLine("}")
	case grammar.Assignment:
		node.AppendLine("{")
		node.Indent(func(indent codegen.Node) {
			// Decision states in the ATN are keyed by the assignable value
			// (RuleCall, Keyword, etc.), not by the Assignment wrapper itself.
			syncCall := buildSyncCall(context, e.Value())
			// When the value is a combinable Alternatives, use the combined path
			// that folds the loop guard into the OR switch's default case.
			if alts, ok := e.Value().(grammar.Alternatives); ok && element.Cardinality() != CardinalityOne && context.canCombineAltsWithLoop(alts) {
				generateCombinedAssignmentParser(indent, context, e, alts, syncCall, element.Cardinality())
				return
			}
			generateCardinality(indent, func(n codegen.Node) {
				if context.completion {
					n.AppendLine("p.cp.MarkAssignment(", strconv.Quote(e.Property().Text()), ")")
				}
				generateAssignable(n, context, e.Value(), func(n2 codegen.Node, resultName string) {
					if context.completion {
						return
					}
					n2.AppendLine("if ", resultName, " != nil {")
					if _, ok := e.Value().(grammar.CrossRef); ok {
						declaringInterfaceName := getDeclaringInterface(e)
						// For cross-references, we need to create a Reference object
						resultName = "p.referencesConstructor." + declaringInterfaceName + e.Property().Text() + "(current, " + resultName + ")"
					}
					n2.Indent(func(in codegen.Node) {
						switch e.Operator() {
						case "+=":
							in.AppendLine("current.Set", e.Property().Text(), "Item(", resultName, ")")
						default:
							in.AppendLine("current.Set", e.Property().Text(), "(", resultName, ")")
						}
					})
					n2.AppendLine("}")
				})
				if context.completion {
					n.AppendLine("p.cp.ClearAssignment()")
				}
			}, func(n codegen.Node) {
				n.Append(guardCall(context, element))
			}, syncCall, element.Cardinality())
		})
		node.AppendLine("}")
	}
}

func generateAssignable(node codegen.Node, context *ParserGeneratorContext, assignable grammar.Assignable, cb func(node codegen.Node, resultName string)) {
	switch a := assignable.(type) {
	case grammar.CrossRef:
		cb(node, generateRuleCallParser(node, context, a.Rule()))
	case grammar.Keyword:
		cb(node, generateKeywordParser(node, context, a))
	case grammar.RuleCall:
		cb(node, generateRuleCallParser(node, context, a))
	case grammar.Alternatives:
		generateAssignableAlternatives(node, context, a, cb)
	default:
		panic("Unresolved assignment assignable")
	}
}

func generateAssignableAlternatives(node codegen.Node, context *ParserGeneratorContext, alts grammar.Alternatives, cb func(node codegen.Node, resultName string)) {
	node.AppendLine(orSwitchHeadAlt(context, alts))
	for i, alt := range alts.Alts() {
		node.AppendLine("case ", strconv.Itoa(i), ":")
		node.Indent(func(in codegen.Node) {
			if assignable, ok := alt.(grammar.Assignable); ok {
				generateAssignable(in, context, assignable, cb)
			}
		})
	}
	node.AppendLine("}")
}

func generateGroupParser(node codegen.Node, context *ParserGeneratorContext, group grammar.Group) {
	syncCall := buildSyncCall(context, group)
	generateCardinality(node, func(n codegen.Node) {
		n.Indent(func(in codegen.Node) {
			for _, element := range group.Elements() {
				generateAbstractElementParser(in, context, element)
			}
		})
	}, func(n codegen.Node) {
		n.Append(guardCall(context, group))
	}, syncCall, group.Cardinality())
}

func generateKeywordParser(node codegen.Node, context *ParserGeneratorContext, keyword grammar.Keyword) string {
	generateCardinality(node, func(n codegen.Node) {
		if context.completion {
			n.AppendLine("p.state.Consume(", GeneratedTokenName(keyword), ")")
		} else {
			n.AppendLine("token := p.state.Consume(", GeneratedTokenName(keyword), ")")
			n.AppendLine("core.AssignToken(current, token, ", context.atnData.elementStateName(keyword), ")")
		}
	}, func(n codegen.Node) { n.Append(guardCall(context, keyword)) }, "", keyword.Cardinality())
	return "token"
}

func generateRuleCallParser(node codegen.Node, context *ParserGeneratorContext, ruleCall grammar.RuleCall) string {
	target := ruleCall.Rule().Ref(ctx.Background())
	var result string
	switch target.(type) {
	case grammar.AbstractTokenRule:
		result = "token"
	default:
		result = "result"
	}
	first := true
	generateCardinality(node, func(n codegen.Node) {
		eq := "="
		if first {
			eq = ":="
			first = false
		}
		switch t := target.(type) {
		case grammar.AbstractTokenRule:
			if context.completion {
				n.AppendLine("p.state.Consume(", GeneratedTokenName(t), ")")
			} else {
				n.AppendLine("token ", eq, " p.state.Consume(", GeneratedTokenName(t), ")")
				n.AppendLine("core.AssignToken(current, token, ", context.atnData.elementStateName(ruleCall), ")")
			}
		case grammar.ParserRule:
			followName := context.atnData.followState(ruleCall)
			n.AppendLine("p.state.EnterRule(", followName, ")")
			if context.completion {
				n.AppendLine("p.Parse", t.Name(), "()")
			} else {
				n.AppendLine("result ", eq, " p.Parse", t.Name(), "()")
			}
			n.AppendLine("p.state.ExitRule()")
		case grammar.CompositeRule:
			followName := context.atnData.followState(ruleCall)
			if context.completion {
				// Composite rules in completion mode have the same no-arg shape
				// as a regular rule; the generated Parse function has no `current` param.
				n.AppendLine("p.state.EnterRule(", followName, ")")
				n.AppendLine("p.Parse", t.Name(), "()")
				n.AppendLine("p.state.ExitRule()")
				return
			}
			if context.inCompositeRule {
				n.AppendLine("p.state.EnterRule(", followName, ")")
				n.AppendLine("p.Parse", t.Name(), "(current)")
				n.AppendLine("p.state.ExitRule()")
			} else {
				n.AppendLine("result ", eq, " core.NewCompositeNode()")
				n.AppendLine("result.SetSegmentStartToken(p.state.LA(1))")
				n.AppendLine("p.state.EnterRule(", followName, ")")
				n.AppendLine("p.Parse", t.Name(), "(result)")
				n.AppendLine("p.state.ExitRule()")
				n.AppendLine("result.SetSegmentEndToken(p.state.LA(0))")
			}
		}
	}, func(n codegen.Node) { n.Append(guardCall(context, ruleCall)) }, "", ruleCall.Cardinality())
	return result
}

func generateAlternativesParser(node codegen.Node, context *ParserGeneratorContext, alts grammar.Alternatives) {
	syncCall := buildSyncCall(context, alts)
	if context.canCombineAlternatives(alts) {
		generateCombinedAlternativesParser(node, context, alts, syncCall)
		return
	}
	generateCardinality(node, func(n codegen.Node) {
		n.AppendLine(orSwitchHead(context, alts))
		for i, alt := range alts.Alts() {
			n.AppendLine("case ", strconv.Itoa(i), ":")
			n.Indent(func(in codegen.Node) {
				generateAbstractElementParser(in, context, alt)
			})
		}
		n.AppendLine("default:")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("p.state.AppendError(p.state.Messages().NoViableAlternative(failure), failure.Token)")
		})
		n.AppendLine("}")
	}, func(n codegen.Node) {
		n.Append(guardCall(context, alts))
	}, syncCall, alts.Cardinality())
}

// generateCombinedLoop emits the shared loop/optional pattern for decisions
// where the OR lookahead drives both the entry/continuation guard and the
// alternative selection. The loop structure is identical for both bare
// Alternatives and Assignment-wrapped Alternatives; only the case bodies differ,
// so emitCases is the caller's responsibility.
func generateCombinedLoop(node codegen.Node, context *ParserGeneratorContext, alts grammar.Alternatives, syncCall, cardinality string, emitCases func(codegen.Node)) {
	switch cardinality {
	case CardinalityOptional:
		if syncCall != "" {
			node.AppendLine(syncCall)
		}
		node.AppendLine(orSwitchHeadAlt(context, alts))
		emitCases(node)
		node.AppendLine("}")
	case CardinalityZeroOrMore:
		if syncCall != "" {
			node.AppendLine(syncCall)
		}
		label := context.nextLoopLabel()
		node.AppendLine(label + ":")
		node.AppendLine("for {")
		node.Indent(func(n codegen.Node) {
			n.AppendLine(orSwitchHeadAlt(context, alts))
			emitCases(n)
			n.AppendLine("default:")
			n.Indent(func(in codegen.Node) {
				in.AppendLine("break " + label)
			})
			n.AppendLine("}")
			if syncCall != "" {
				n.AppendLine(syncCall)
			}
		})
		node.AppendLine("}")
	case CardinalityOneOrMore:
		label := context.nextLoopLabel()
		node.AppendLine(label + ":")
		node.AppendLine("for first := true; ; first = false {")
		node.Indent(func(n codegen.Node) {
			n.AppendLine(orSwitchHead(context, alts))
			emitCases(n)
			n.AppendLine("default:")
			n.Indent(func(in codegen.Node) {
				in.AppendLine("if !first {")
				in.Indent(func(iin codegen.Node) {
					iin.AppendLine("break " + label)
				})
				in.AppendLine("}")
				in.AppendLine("p.state.AppendError(p.state.Messages().NoViableAlternative(failure), failure.Token)")
			})
			n.AppendLine("}")
			if syncCall != "" {
				n.AppendLine(syncCall)
			}
		})
		node.AppendLine("}")
	}
}

// generateCombinedAlternativesParser emits the combined loop/optional pattern
// for a bare Alternatives node — the OR lookahead drives both the
// entry/continuation guard and the alternative selection.
func generateCombinedAlternativesParser(node codegen.Node, context *ParserGeneratorContext, alts grammar.Alternatives, syncCall string) {
	emitCases := func(n codegen.Node) {
		for i, alt := range alts.Alts() {
			n.AppendLine("case ", strconv.Itoa(i), ":")
			n.Indent(func(in codegen.Node) {
				generateAbstractElementParser(in, context, alt)
			})
		}
	}
	generateCombinedLoop(node, context, alts, syncCall, alts.Cardinality(), emitCases)
}

// generateCombinedAssignmentParser emits the combined loop/optional pattern for
// an Assignment whose value is an Alternatives node.
func generateCombinedAssignmentParser(node codegen.Node, context *ParserGeneratorContext, e grammar.Assignment, alts grammar.Alternatives, syncCall, cardinality string) {
	resultCb := func(n codegen.Node, resultName string) {
		if context.completion {
			return
		}
		n.AppendLine("if ", resultName, " != nil {")
		if _, ok := e.Value().(grammar.CrossRef); ok {
			interfaceName := getDeclaringInterface(e)
			resultName = "p.referencesConstructor." + interfaceName + e.Property().Text() + "(current, " + resultName + ")"
		}
		n.Indent(func(in codegen.Node) {
			switch e.Operator() {
			case "+=":
				in.AppendLine("current.Set", e.Property().Text(), "Item(", resultName, ")")
			default:
				in.AppendLine("current.Set", e.Property().Text(), "(", resultName, ")")
			}
		})
		n.AppendLine("}")
	}
	emitCases := func(n codegen.Node) {
		for i, alt := range alts.Alts() {
			n.AppendLine("case ", strconv.Itoa(i), ":")
			n.Indent(func(in codegen.Node) {
				if context.completion {
					in.AppendLine("p.cp.MarkAssignment(", strconv.Quote(e.Property().Text()), ")")
				}
				if assignable, ok := alt.(grammar.Assignable); ok {
					generateAssignable(in, context, assignable, resultCb)
				}
				if context.completion {
					in.AppendLine("p.cp.ClearAssignment()")
				}
			})
		}
	}
	generateCombinedLoop(node, context, alts, syncCall, cardinality, emitCases)
}

// generateCardinality emits the loop/guard wrapper for a single grammar element.
// syncCall (if non-empty) is a Sync statement emitted before the guard and
// re-emitted at the bottom of each loop iteration to discard unexpected tokens
// before the next guard check. It replaces the old codegen.Callback preGuard
// pattern to avoid scattered nil checks.
func generateCardinality(node codegen.Node, element, lookahead codegen.Callback, syncCall, cardinality string) {
	switch cardinality {
	case CardinalityOne:
		element(node)
	case CardinalityOptional:
		if syncCall != "" {
			node.AppendLine(syncCall)
		}
		node.Append("if ")
		lookahead(node)
		node.AppendLine(" {")
		node.Indent(element)
		node.AppendLine("}")
	case CardinalityZeroOrMore:
		if syncCall != "" {
			node.AppendLine(syncCall)
		}
		node.Append("for ")
		lookahead(node)
		node.AppendLine(" {")
		node.Indent(func(n codegen.Node) {
			element(n)
			if syncCall != "" {
				n.AppendLine(syncCall)
			}
		})
		node.AppendLine("}")
	case CardinalityOneOrMore:
		node.Append("for ok := true; ok; ok = ")
		lookahead(node)
		node.AppendLine(" {")
		node.Indent(func(n codegen.Node) {
			element(n)
			if syncCall != "" {
				n.AppendLine(syncCall)
			}
		})
		node.AppendLine("}")
	}
}

type LL1Decision [][]string

// ll1DecisionOr builds the LL(1) lookahead table for an alternatives
// decision from the ATN. Each outgoing transition of the alternative-choice
// state corresponds to one alternative (in order); its FIRST set becomes the
// lookahead option. Decisions that a single token cannot disambiguate are
// routed to adaptive prediction (ALL(*)) instead, so a one-token FIRST set is
// always sufficient here.
func ll1DecisionOr(atnData *parserATNData, element grammar.Alternatives) LL1Decision {
	if atnData == nil {
		return nil
	}
	ds := atnData.orDecision[element]
	if ds == nil {
		return nil
	}
	lookahead := make(LL1Decision, 0, len(ds.Transitions))
	for _, t := range ds.Transitions {
		lookahead = append(lookahead, atnData.firstSetOption(t.Target()))
	}
	return lookahead
}

// firstSetOption converts the FIRST set reachable from state into a lookahead
// option: one single-token path per discriminating token type, ordered by ATN
// token id for deterministic output. Token-group ids are emitted as-is and
// expanded to leaf ids later by generateLL1Lookahead.
func (d *parserATNData) firstSetOption(state *internalATN.ATNState) []string {
	first := firstSetFromATNState(state)
	ids := []int{}
	for id := range first {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	option := []string{}
	for _, id := range ids {
		option = append(option, d.tokenVarNames[id])
	}
	return option
}

func buildVarNameToId(tokenTypes GenerateTokenTypesResult) map[string]int {
	varNameToId := make(map[string]int, len(tokenTypes.TokenTypeVarNames))
	for i, varName := range tokenTypes.TokenTypeVarNames {
		varNameToId[varName] = i + 1 // token IDs start at 1
	}
	return varNameToId
}

func buildGroupVarNameToMembers(grammr grammar.Grammar) map[string][]string {
	keywords := GetAllKeywords(grammr)
	groups := make(map[string][]string)
	for _, tg := range grammr.TokenGroups() {
		varName := GeneratedTokenName(tg)
		groups[varName] = getAllTokenGroupMembers(tg, keywords)
	}
	return groups
}

func getLeafIdsForVarName(varName string, varNameToId map[string]int, groupMembers map[string][]string) []int {
	if members, ok := groupMembers[varName]; ok {
		seen := map[int]bool{}
		var result []int
		for _, member := range members {
			for _, id := range getLeafIdsForVarName(member, varNameToId, groupMembers) {
				if !seen[id] {
					seen[id] = true
					result = append(result, id)
				}
			}
		}
		return result
	}
	if id, ok := varNameToId[varName]; ok {
		return []int{id}
	}
	return nil
}

func generateLL1Lookahead(node codegen.Node, name string, lookahead LL1Decision, varNameToId map[string]int, groupMembers map[string][]string) {
	maxId := len(varNameToId) // IDs are 1..N
	lookup := make([]int, maxId+1)
	for i := range lookup {
		lookup[i] = -1
	}

	seen := map[string]bool{}
	var typeNames []string
	for altIdx, option := range lookahead {
		for _, varName := range option {
			if !seen[varName] {
				seen[varName] = true
				typeNames = append(typeNames, varName)
			}
			for _, leafId := range getLeafIdsForVarName(varName, varNameToId, groupMembers) {
				if leafId > 0 && leafId <= maxId {
					lookup[leafId] = altIdx
				}
			}
		}
	}

	lookupParts := make([]string, 0, len(lookup))
	for i, v := range lookup {
		if v >= 0 {
			lookupParts = append(lookupParts, strconv.Itoa(i)+": "+strconv.Itoa(v+1))
		}
	}

	node.AppendLine("var ", name, " = parser.LL1Lookahead{")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("Types:  []*core.TokenType{", strings.Join(typeNames, ", "), "},")
		n.AppendLine("Lookup: []int{", strings.Join(lookupParts, ", "), "},")
	})
	node.AppendLine("}")
}

// ll1DecisionOpt builds the single-option LL(1) table for a cardinality
// guard (?, *, +) from the ATN. The guard enters/continues the element when the
// next token is in the FIRST set of the loop body, which is the FIRST set of the
// guard state's enter transition (transition 0; transition 1 is the exit edge).
func ll1DecisionOpt(atnData *parserATNData, element grammar.Element) LL1Decision {
	if atnData == nil {
		return nil
	}
	// The ATN keys a cardinality guard by the element it actually wraps. For an
	// assignment that is the assignable value (the keyword/rule/alternatives),
	// not the Assignment wrapper — mirror the adaptive-map keying.
	key := element
	if assignment, ok := element.(grammar.Assignment); ok {
		if value, ok := assignment.Value().(grammar.Element); ok {
			key = value
		}
	}
	ds := atnData.loopDecision[key]
	if ds == nil || len(ds.Transitions) == 0 {
		return nil
	}
	return LL1Decision{
		atnData.firstSetOption(ds.Transitions[0].Target()),
	}
}

// getDeclaringInterface returns the name of the interface that declares
// the property being assigned by the given assignment. Searches backwards
// through the sequential context for the nearest preceding Action (which
// re-types 'current') and falls back to FindReturnType on the containing
// ParserRule.
func getDeclaringInterface(assignment grammar.Assignment) string {
	action, parserRule := findPrecedingAction(assignment)
	if action != nil {
		return action.Type().Text()
	}
	if parserRule != nil {
		if iface := grammar.FindReturnType(parserRule, ctx.Background()); iface != nil {
			return iface.Name()
		}
		panic("No return type for rule: " + parserRule.Name())
	}
	panic("Unable to find containing parser rule for cross-reference assignment")
}

// findPrecedingAction searches backwards from element through its containing
// sequential contexts (Groups) for the most recent preceding grammar.Action.
// Stops at ParserRule boundaries (returning the rule). Passes through other
// container types such as Alternatives by recursing upward without scanning
// their siblings. Panics for unexpected container types (e.g. CompositeRule,
// which cannot contain assignments by design).
func findPrecedingAction(element grammar.Element) (grammar.Action, grammar.ParserRule) {
	container := element.Container()
	switch c := container.(type) {
	case grammar.Group:
		elements := c.Elements()
		idx := slices.Index(elements, element)
		for j := idx - 1; j >= 0; j-- {
			if act := lastActionIn(elements[j]); act != nil {
				return act, nil
			}
		}
		return findPrecedingAction(c)
	case grammar.ParserRule:
		return nil, c
	case grammar.Element:
		return findPrecedingAction(c)
	default:
		panic("Unable to find parser rule for element")
	}
}

// lastActionIn returns the last grammar.Action in the sequential tail of
// element, scanning backwards. Used to check whether a Group preceding the
// assignment contains a trailing action that re-typed 'current'. Returns nil
// for Alternatives (conditional branches are not scanned).
func lastActionIn(element grammar.Element) grammar.Action {
	switch e := element.(type) {
	case grammar.Action:
		return e
	case grammar.Group:
		for i := len(e.Elements()) - 1; i >= 0; i-- {
			if act := lastActionIn(e.Elements()[i]); act != nil {
				return act
			}
		}
	}
	return nil
}
