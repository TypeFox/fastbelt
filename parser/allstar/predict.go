// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"fmt"
	"strings"
	"sync"

	core "typefox.dev/fastbelt"
)

// AmbiguityReport is a callback invoked when an ambiguity is detected.
type AmbiguityReport func(message string)

// TokenSource is the minimal interface required by the prediction algorithm.
// *parser.ParserState satisfies this interface.
type TokenSource interface {
	LA(offset int) *core.Token
}

// tokenTypeID returns the token type ID for a token.
// nil is treated as EOF (ID 0).
func tokenTypeID(t *core.Token) int {
	if t == nil {
		return core.EOF.Id
	}
	return t.TypeId
}

// PredicateSet records which alternatives are guarded and whether their
// gate predicate currently evaluates to true.
type PredicateSet struct {
	predicates []bool
}

// Is returns true when index is out of range (unconstrained) or when
// predicates[index] is true.
func (p *PredicateSet) Is(index int) bool {
	if p == nil || index >= len(p.predicates) {
		return true
	}
	return p.predicates[index]
}

// Set sets the predicate value for the given alternative index.
func (p *PredicateSet) Set(index int, value bool) {
	for len(p.predicates) <= index {
		p.predicates = append(p.predicates, true)
	}
	p.predicates[index] = value
}

// String serialises the predicate set as a binary string ("10110…").
func (p *PredicateSet) String() string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	for _, v := range p.predicates {
		if v {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

// EmptyPredicates is the zero-value PredicateSet used when there are no gates.
var EmptyPredicates = &PredicateSet{}

// dfaCache returns a DFA for the given predicate configuration.
// Different predicate sets may produce different prediction decisions.
type dfaCache func(predicates *PredicateSet) *DFA

// predictError is an unexported error returned when the lookahead finds no
// valid path. Callers that need field access use errors.As.
type predictError struct {
	tokenPath         []int // token type IDs of consumed lookahead
	possibleTypeIDs   []int
	actualTokenTypeID int
}

func (e *predictError) Error() string {
	return fmt.Sprintf("adaptive predict error: unexpected token type %d after path %v (possible: %v)",
		e.actualTokenTypeID, e.tokenPath, e.possibleTypeIDs)
}

// newPredictError constructs a predictError from the lookahead path, the
// previous DFA state (used to collect possible token types), and the actual
// token type ID.
func newPredictError(path []int, prev *DFAState, actualTypeID int) error {
	possible := make([]int, 0)
	seen := map[int]bool{}
	for _, c := range prev.Configs.Elements() {
		for _, t := range c.State.Transitions {
			at, ok := t.(*AtomTransition)
			if !ok {
				continue
			}
			if !seen[at.TokenTypeID] {
				seen[at.TokenTypeID] = true
				possible = append(possible, at.TokenTypeID)
			}
		}
	}
	return &predictError{
		tokenPath:         path,
		possibleTypeIDs:   possible,
		actualTokenTypeID: actualTypeID,
	}
}

// newDFACache creates a dfaCache backed by a mutex-protected map.
func newDFACache(start *ATNState, decision int) dfaCache {
	var mu sync.RWMutex
	m := map[string]*DFA{}
	return func(predicates *PredicateSet) *DFA {
		key := predicates.String()
		mu.RLock()
		existing, ok := m[key]
		mu.RUnlock()
		if ok {
			return existing
		}
		mu.Lock()
		defer mu.Unlock()
		// double-check after upgrading lock
		if existing, ok = m[key]; ok {
			return existing
		}
		dfa := &DFA{
			ATNStartState: start,
			Decision:      decision,
			States:        map[string]*DFAState{},
		}
		m[key] = dfa
		return dfa
	}
}

// initDFACaches creates one dfaCache per decision state.
func initDFACaches(atn *ATN) []dfaCache {
	caches := make([]dfaCache, len(atn.DecisionStates))
	for i, ds := range atn.DecisionStates {
		caches[i] = newDFACache(ds, i)
	}
	return caches
}

// adaptivePredict runs the ALL(*) algorithm and returns the chosen alternative
// (0-based), or -1 on error.
func adaptivePredict(src TokenSource, dfas []dfaCache, decision int, preds *PredicateSet, log AmbiguityReport) (int, error) {
	dfa := dfas[decision](preds)
	if dfa.Start == nil {
		configs := computeStartState(dfa.ATNStartState)
		dfa.Start = addDFAState(dfa, newDFAState(configs))
	}
	return performLookahead(src, dfa, dfa.Start, preds, log)
}

// performLookahead walks the DFA, extending it as needed.
func performLookahead(src TokenSource, dfa *DFA, s0 *DFAState, preds *PredicateSet, log AmbiguityReport) (int, error) {
	prev := s0
	i := 1
	path := []int{}
	t := src.LA(i)
	i++

	for {
		tID := tokenTypeID(t)
		d := getExistingTargetState(prev, tID)
		if d == nil {
			d = computeLookaheadTarget(src, dfa, prev, tID, i, preds, log)
		}

		if d == DFAError {
			return -1, newPredictError(path, prev, tID)
		}

		if d.IsAcceptState {
			return d.Prediction, nil
		}

		prev = d
		path = append(path, tID)
		t = src.LA(i)
		i++
	}
}

// computeLookaheadTarget computes (and memoizes) the DFA transition from prev
// on the given token type ID.
func computeLookaheadTarget(src TokenSource, dfa *DFA, prev *DFAState, tID int, lookahead int, preds *PredicateSet, log AmbiguityReport) *DFAState {
	reach := computeReachSet(prev.Configs, tID, preds)
	if reach.Len() == 0 {
		addDFAEdge(dfa, prev, tID, DFAError)
		return DFAError
	}

	next := newDFAState(reach)
	predictedAlt, hasUnique := getUniqueAlt(reach, preds)
	if hasUnique {
		next.IsAcceptState = true
		next.Prediction = predictedAlt
		next.Configs.UniqueAlt = predictedAlt
	} else if hasConflictTerminatingPrediction(reach) {
		prediction := minAlt(reach.Alts())
		next.IsAcceptState = true
		next.Prediction = prediction
		next.Configs.UniqueAlt = prediction
		reportLookaheadAmbiguity(src, dfa, lookahead, reach.Alts(), log)
	}

	next = addDFAEdge(dfa, prev, tID, next)
	return next
}

// getExistingTargetState returns the edge for tokenTypeID, or nil if absent.
func getExistingTargetState(state *DFAState, tID int) *DFAState {
	if state.Edges == nil {
		return nil
	}
	return state.Edges[tID]
}

// computeReachSet advances the simulation one token.
func computeReachSet(configs *ATNConfigSet, tID int, preds *PredicateSet) *ATNConfigSet {
	intermediate := NewATNConfigSet()
	var skippedStopStates []*ATNConfig

	for _, c := range configs.Elements() {
		if !preds.Is(c.Alt) {
			continue
		}
		if c.State.Type == ATNRuleStop {
			skippedStopStates = append(skippedStopStates, c)
			continue
		}
		for _, t := range c.State.Transitions {
			target := getReachableTarget(t, tID)
			if target != nil {
				intermediate.Add(&ATNConfig{
					State: target,
					Alt:   c.Alt,
					Stack: c.Stack,
				})
			}
		}
	}

	var reach *ATNConfigSet
	if len(skippedStopStates) == 0 && intermediate.Len() == 1 {
		reach = intermediate
	}

	if reach == nil {
		reach = NewATNConfigSet()
		for _, c := range intermediate.Elements() {
			closure(c, reach)
		}
	}

	if len(skippedStopStates) > 0 && !hasConfigInRuleStopState(reach) {
		for _, c := range skippedStopStates {
			reach.Add(c)
		}
	}

	return reach
}

// getReachableTarget returns the target state if t matches tID, nil otherwise.
func getReachableTarget(t Transition, tID int) *ATNState {
	at, ok := t.(*AtomTransition)
	if !ok {
		return nil
	}
	if tokenMatches(tID, at.TokenTypeID, at.CategoryMatches) {
		return at.target
	}
	return nil
}

// tokenMatches checks whether tokenTypeID matches the transition's type,
// including via category inheritance.
func tokenMatches(tokenTypeID, transitionTypeID int, categoryMatches []int) bool {
	if tokenTypeID == transitionTypeID {
		return true
	}
	for _, cat := range categoryMatches {
		if tokenTypeID == cat {
			return true
		}
	}
	return false
}

// closure follows all epsilon transitions from config, adding reachable
// configs to configs set.
func closure(config *ATNConfig, configs *ATNConfigSet) {
	p := config.State

	if p.Type == ATNRuleStop {
		if len(config.Stack) > 0 {
			stack := make([]*ATNState, len(config.Stack))
			copy(stack, config.Stack)
			followState := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			closure(&ATNConfig{
				State: followState,
				Alt:   config.Alt,
				Stack: stack,
			}, configs)
		} else {
			// Dip into outer context: add as-is.
			configs.Add(config)
		}
		return
	}

	if !p.EpsilonOnlyTransitions {
		configs.Add(config)
	}

	for _, t := range p.Transitions {
		c := getEpsilonTarget(config, t)
		if c != nil {
			closure(c, configs)
		}
	}
}

// getEpsilonTarget returns the config reached via an epsilon transition, or nil.
func getEpsilonTarget(config *ATNConfig, t Transition) *ATNConfig {
	switch et := t.(type) {
	case *EpsilonTransition:
		return &ATNConfig{
			State: et.target,
			Alt:   config.Alt,
			Stack: config.Stack,
		}
	case *RuleTransition:
		stack := make([]*ATNState, len(config.Stack)+1)
		copy(stack, config.Stack)
		stack[len(config.Stack)] = et.FollowState
		return &ATNConfig{
			State: et.target,
			Alt:   config.Alt,
			Stack: stack,
		}
	}
	return nil
}

// computeStartState builds the initial DFAState for a decision.
func computeStartState(atnState *ATNState) *ATNConfigSet {
	configs := NewATNConfigSet()
	for i, t := range atnState.Transitions {
		config := &ATNConfig{
			State: t.Target(),
			Alt:   i,
			Stack: []*ATNState{},
		}
		closure(config, configs)
	}
	return configs
}

// getUniqueAlt returns the single alt shared by all configs (filtered by preds),
// and true. Returns 0, false when there are multiple distinct alts.
func getUniqueAlt(configs *ATNConfigSet, preds *PredicateSet) (int, bool) {
	alt := -1
	for _, c := range configs.Elements() {
		if !preds.Is(c.Alt) {
			continue
		}
		if alt == -1 {
			alt = c.Alt
		} else if alt != c.Alt {
			return 0, false
		}
	}
	if alt == -1 {
		return 0, false
	}
	return alt, true
}

// hasConflictTerminatingPrediction returns true when the config set indicates
// that prediction can terminate with a conflict (ambiguity).
func hasConflictTerminatingPrediction(configs *ATNConfigSet) bool {
	if allConfigsInRuleStopStates(configs) {
		return true
	}
	altSets := getConflictingAltSets(configs.Elements())
	return hasConflictingAltSet(altSets) && !hasStateAssociatedWithOneAlt(altSets)
}

func allConfigsInRuleStopStates(configs *ATNConfigSet) bool {
	for _, c := range configs.Elements() {
		if c.State.Type != ATNRuleStop {
			return false
		}
	}
	return true
}

func hasConfigInRuleStopState(configs *ATNConfigSet) bool {
	for _, c := range configs.Elements() {
		if c.State.Type == ATNRuleStop {
			return true
		}
	}
	return false
}

// getConflictingAltSets groups configs by their state/stack key (alt excluded)
// and returns the set of alts associated with each group.
func getConflictingAltSets(configs []*ATNConfig) map[string]map[int]bool {
	result := map[string]map[int]bool{}
	for _, c := range configs {
		key := atnConfigKey(c, false)
		alts := result[key]
		if alts == nil {
			alts = map[int]bool{}
			result[key] = alts
		}
		alts[c.Alt] = true
	}
	return result
}

func hasConflictingAltSet(altSets map[string]map[int]bool) bool {
	for _, alts := range altSets {
		if len(alts) > 1 {
			return true
		}
	}
	return false
}

func hasStateAssociatedWithOneAlt(altSets map[string]map[int]bool) bool {
	for _, alts := range altSets {
		if len(alts) == 1 {
			return true
		}
	}
	return false
}

func newDFAState(configs *ATNConfigSet) *DFAState {
	return &DFAState{
		Configs:    configs,
		Edges:      map[int]*DFAState{},
		Prediction: -1,
	}
}

func addDFAEdge(dfa *DFA, from *DFAState, tID int, to *DFAState) *DFAState {
	to = addDFAState(dfa, to)
	if from.Edges == nil {
		from.Edges = map[int]*DFAState{}
	}
	from.Edges[tID] = to
	return to
}

func addDFAState(dfa *DFA, state *DFAState) *DFAState {
	if state == DFAError {
		return state
	}
	key := state.Configs.Key()
	if existing, ok := dfa.States[key]; ok {
		return existing
	}
	state.Configs.Finalize()
	dfa.States[key] = state
	return state
}

func reportLookaheadAmbiguity(src TokenSource, dfa *DFA, lookahead int, alts []int, log AmbiguityReport) {
	atnState := dfa.ATNStartState
	rule := atnState.Rule
	production := atnState.Production
	msg := buildAmbiguityError(rule, production, lookahead, alts)
	log(msg)
}

func buildAmbiguityError(rule *Rule, prod Production, lookahead int, alts []int) string {
	altStrs := make([]string, len(alts))
	for i, a := range alts {
		altStrs[i] = fmt.Sprintf("%d", a)
	}
	// Chevrotain stores first-occurrence idx as 0 (no suffix) and subsequent
	// ones as 1, 2, … Our grammar model uses 1-based Idx, so we subtract 1
	// to get the Chevrotain-compatible index used in message formatting.
	var occurrence string
	if prod != nil {
		tsIdx := prod.Occurrence() - 1
		if tsIdx > 0 {
			occurrence = fmt.Sprintf("%d", tsIdx)
		}
	}
	dslName := productionDSLName(prod)
	return fmt.Sprintf(
		"Ambiguous Alternatives Detected: <%s> in <%s%s> inside <%s> Rule,\n"+
			"lookahead depth %d may appear as a prefix path in all these alternatives.\n"+
			"See: https://chevrotain.io/docs/guide/resolving_grammar_errors.html#AMBIGUOUS_ALTERNATIVES\n"+
			"For Further details.",
		strings.Join(altStrs, ", "),
		dslName,
		occurrence,
		rule.Name,
		lookahead,
	)
}

func productionDSLName(prod Production) string {
	if prod == nil {
		return "unknown"
	}
	switch prod.Kind() {
	case ProdNonTerminal:
		return "SUBRULE"
	case ProdOption:
		return "OPTION"
	case ProdAlternation:
		return "OR"
	case ProdRepetitionMandatory:
		return "AT_LEAST_ONE"
	case ProdRepetition:
		return "MANY"
	case ProdTerminal:
		return "CONSUME"
	default:
		return "unknown"
	}
}

// minAlt returns the smallest value in the slice.
// Returns -1 for empty slices.
func minAlt(alts []int) int {
	if len(alts) == 0 {
		return -1
	}
	m := alts[0]
	for _, v := range alts[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
