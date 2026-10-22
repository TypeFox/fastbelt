// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"context"
	"iter"
	"slices"
	"unique"

	core "typefox.dev/fastbelt"
)

// LanguageCompletionAdapter bridges the framework's generic completion
// provider and the generated artefacts of a language (CompletionParser, ATN,
// synthetic node factories, completion dispatch table). It lives in this
// package because it is a parser-level contract that does not touch LSP or
// server-side types. The code generator emits the implementation.
type LanguageCompletionAdapter interface {
	// Parse runs the generated CompletionParser over the given prefix
	// tokens (the document's tokens up to the cursor).
	Parse(tokens []core.Token) *CompletionParseResult
	// ATN returns the RuntimeATN the simulator queries.
	ATN() *RuntimeATN
	// SyntheticOwnerFor returns a fresh, detached AST node of the given
	// node type (e.g. "Transition"). The completion provider uses it for
	// the owner of a cross-reference, and for its containers, when the main
	// parser did not create them. Returns (nil, false) for unknown types.
	SyntheticOwnerFor(nodeType string) (core.AstNode, bool)
	// DispatchCompletion runs the per-field scope+filter chain for the
	// given CompletionHint.Key on the supplied synthetic owner. Returns
	// (nil, false) if the field key isn't registered.
	DispatchCompletion(
		ctx context.Context,
		field string,
		owner core.AstNode,
	) (iter.Seq[*core.SymbolDescription], bool)
	// HasAssignment reports whether node's property field has been
	// assigned. Used by the completion engine to decide whether to wrap
	// the existing AST node in a synthetic parent via a PrecedingAction.
	// Grammars without cross-reference assignments emit a body that
	// always returns false.
	HasAssignment(node core.AstNode, property string) bool
	// ApplyAction allocates a new AST node of actionType and assigns
	// value to its property field. Whether the setter is SetX or
	// SetXItem is decided at generation time from the field's grammar
	// type (single vs array). Returns nil when no matching
	// (actionType, property) pair is registered or when value's dynamic
	// type is not assignable to the property. Grammars without
	// tree-rewrite actions emit a body that always returns nil.
	ApplyAction(actionType, property string, value core.AstNode) core.AstNode
	// AssignsCurrent reports whether the child at the given field and index
	// of the container was added by a tree-rewriting action or as the operand
	// of an infix rule. Such a child is created by the same rule call as the
	// container, in contrast to the child of an assigned rule call.
	AssignsCurrent(container core.AstNode, field unique.Handle[string], index int) bool
}

// CompletionParseResult is the output of a CompletionParser run. It contains
// every artefact the completion provider needs to drive the ATN simulator.
type CompletionParseResult struct {
	// Tokens is the token slice the parser was given (typically the document's
	// prefix tokens up to the cursor).
	Tokens []core.Token
	// Snapshots is the time-ordered sequence of snapshots recorded at every
	// rule entry and every Sync. SimulateAt picks the simulator's start from it.
	Snapshots []ATNSnapshot
	// RuleStack is the chain of rule contexts the parser was inside when it
	// stopped, outermost first. The completion provider falls back to it for
	// the owner of a cross-reference that it cannot derive from the AST of
	// the main parser.
	RuleStack []RuleContext
}

// ATNSnapshot records "at TokenIdx tokens consumed, the parser was positioned
// at ATN state ATNStateIdx, inside RuleStack". The generated CompletionParser
// records one at every rule entry and every Sync; the simulator consumes the
// remaining tokens from there.
//
// FollowStates is the stack of the ATN states that the rule calls in progress
// return to. It lets the simulator leave the rule that the snapshot was taken
// in. Fresh and Left describe how these rule calls changed since the last
// token was consumed: the topmost Fresh entries of FollowStates were added
// since then, and Left holds the FollowState indices of the assigned rule
// calls that ended since then, in that order.
//
// RuleStack is a copy of the live rule stack at the time of the snapshot.
type ATNSnapshot struct {
	TokenIdx     int
	ATNStateIdx  int
	RuleStack    []RuleContext
	FollowStates []int
	Fresh        int
	Left         []int
}

// RuleContext is a single frame of CompletionParseResult.RuleStack. Type is
// the name of the node type that the rule returns, e.g. "Statemachine", as
// used by [LanguageCompletionAdapter.SyntheticOwnerFor]. It is empty for
// rules that do not create nodes, such as composite rules.
type RuleContext struct {
	Type string
}

// CompletionParserState is the helper a generated CompletionParser uses to
// record snapshots and maintain its rule stack. It wraps a regular ParserState
// (which still does the consume/lookahead/sync work) and adds the bookkeeping
// completion needs.
//
// The generator emits, per rule:
//
//	func (p *CompletionParser) ParseFoo() {
//	    p.cp.EnterRule("Foo", Foo__Start)
//	    defer p.cp.ExitRule()
//	    // ...mirrors the main parser's control flow, but every Sync(i) is
//	    // preceded by p.cp.RecordSnapshot(i), and all AST mutation is omitted.
//	}
type CompletionParserState struct {
	state *ParserState

	snapshots []ATNSnapshot
	ruleStack []RuleContext
	// entries is parallel to ruleStack and holds the token index at which
	// each rule was entered.
	entries []int
	// left holds the FollowState indices of the assigned rule calls that
	// ended at the token index leftAt.
	left   []int
	leftAt int
}

// NewCompletionParserState wraps an existing ParserState so a generated
// CompletionParser can record its progress.
func NewCompletionParserState(state *ParserState) *CompletionParserState {
	return &CompletionParserState{state: state}
}

// State returns the underlying ParserState; the generated CompletionParser
// uses it for the same Consume/Sync/Lookahead/EnterRule/ExitRule calls the
// main parser makes.
func (cp *CompletionParserState) State() *ParserState {
	return cp.state
}

// EnterRule pushes a RuleContext onto the rule stack and records a snapshot at
// the rule's start state. The generator emits one call at the start of every
// rule function.
//
// While the parser is in error mode, EnterRule still pushes a RuleContext (so
// the rule stack reflects what the generated code is doing) but does not
// record a snapshot, as there is no parsing progress to capture. Snapshots
// resume once recovery clears the error mode.
func (cp *CompletionParserState) EnterRule(nodeType string, ruleStartStateIdx int) {
	cp.ruleStack = append(cp.ruleStack, RuleContext{Type: nodeType})
	cp.entries = append(cp.entries, cp.state.Index)
	if cp.state.ErrorMode != ErrorModeNone {
		return
	}
	cp.snapshots = append(cp.snapshots, cp.snapshot(ruleStartStateIdx))
}

// snapshot describes the current state of the parser at the given ATN state.
func (cp *CompletionParserState) snapshot(atnStateIdx int) ATNSnapshot {
	index := cp.state.Index
	followStates := cp.state.followStates
	// Count the rules entered at the current token. The entry rule has no
	// follow state and is not counted.
	fresh := 0
	for i := len(cp.entries) - 1; i >= 0 && cp.entries[i] == index && fresh < len(followStates); i-- {
		fresh++
	}
	// leaveRule never changes cp.left in place, so the snapshot can share it
	var left []int
	if cp.leftAt == index {
		left = cp.left
	}
	return ATNSnapshot{
		TokenIdx:     index,
		ATNStateIdx:  atnStateIdx,
		RuleStack:    slices.Clone(cp.ruleStack),
		FollowStates: slices.Clone(followStates),
		Fresh:        fresh,
		Left:         left,
	}
}

// ExitRule pops the top RuleContext. The generator emits one deferred call at
// the start of every rule function.
//
// ExitRule always pops, even while the parser is in error mode, so the live
// stack stays balanced when the parser recovers and enters more rules. Each
// snapshot holds its own copy of the stack, so popping loses nothing.
func (cp *CompletionParserState) ExitRule() {
	if last := len(cp.ruleStack) - 1; last >= 0 {
		cp.leaveRule(cp.entries[last])
		cp.ruleStack = cp.ruleStack[:last]
		cp.entries = cp.entries[:last]
	}
}

// leaveRule records the rule call that ends, if it is an assigned rule call
// that consumed a token. The rule was entered at the token index entry.
func (cp *CompletionParserState) leaveRule(entry int) {
	index := cp.state.Index
	if cp.leftAt != index {
		cp.left, cp.leftAt = nil, index
	}
	followStates := cp.state.followStates
	if entry == index || len(followStates) == 0 || cp.state.atn == nil {
		return
	}
	followIdx := followStates[len(followStates)-1]
	if cp.state.atn.ruleCallAt(followIdx).leavesToContainer() {
		cp.left = append(slices.Clip(cp.left), followIdx)
	}
}

// RecordSnapshot records a snapshot at the given ATN state. The generator
// emits one call right before every p.state.Sync(idx) call so the simulator
// has a known starting point at every branch boundary.
//
// While the parser is in error mode, RecordSnapshot is a no-op: parsing has
// stopped advancing, so further snapshots would only point past the actual
// progress.
func (cp *CompletionParserState) RecordSnapshot(atnStateIdx int) {
	if cp.state.ErrorMode != ErrorModeNone {
		return
	}
	cp.snapshots = append(cp.snapshots, cp.snapshot(atnStateIdx))
}

// Result builds the CompletionParseResult. The generated Parse method calls
// it once after the entry rule returns.
func (cp *CompletionParserState) Result(tokens []core.Token) *CompletionParseResult {
	snapshots := append([]ATNSnapshot(nil), cp.snapshots...)
	return &CompletionParseResult{
		Tokens:    tokens,
		Snapshots: snapshots,
		RuleStack: deepestRuleStack(snapshots),
	}
}

// deepestRuleStack returns the deepest RuleStack among the snapshots taken at
// the latest token position the parser reached, i.e. the most-nested rule the
// parser was in at that position.
func deepestRuleStack(snapshots []ATNSnapshot) []RuleContext {
	if len(snapshots) == 0 {
		return nil
	}
	maxTokenIdx := snapshots[0].TokenIdx
	for _, s := range snapshots {
		if s.TokenIdx > maxTokenIdx {
			maxTokenIdx = s.TokenIdx
		}
	}
	var best []RuleContext
	for _, s := range snapshots {
		if s.TokenIdx != maxTokenIdx {
			continue
		}
		if len(s.RuleStack) > len(best) {
			best = s.RuleStack
		}
	}
	return append([]RuleContext(nil), best...)
}

// SimulateAt simulates the ATN forward to the given cursor token index,
// trying candidate snapshots until one yields a non-empty result. Returns the
// live set and the snapshot the simulator started from.
//
// Snapshots are tried earliest first. The earliest snapshot is the broadest
// context, since the parser had not yet committed to any LL(k) decision, so
// simulating forward from it keeps every alternative consistent with the
// typed input alive. A later snapshot reflects a single branch the parser
// committed to and would miss completions from sibling branches.
//
// Later snapshots are fallbacks: when the prefix is unparseable from an
// earlier snapshot (a syntax error stops the simulation), a snapshot closer
// to the cursor may still yield a useful result.
//
// "Non-empty result" means that the simulator consumed the whole prefix and
// that NextCompletionsFromSet finds reachable atom transitions; the latter
// rules out snapshots that end at a RuleStop with an empty return stack.
func (r *CompletionParseResult) SimulateAt(atn *RuntimeATN, cursor int) (live []simPath, snap ATNSnapshot, ok bool) {
	if len(r.Snapshots) == 0 || cursor < 0 {
		return nil, ATNSnapshot{}, false
	}
	try := func(s ATNSnapshot) ([]simPath, bool) {
		// Seed the return stack with the parser's rule calls, so that the
		// simulator can leave the rule that the snapshot was taken in
		seed := simPath{stateIdx: s.ATNStateIdx, stack: s.FollowStates, fresh: s.Fresh, left: s.Left, consumedAt: -1}
		seed.hints = make([]*CompletionHint, len(s.FollowStates))
		for i, followIdx := range s.FollowStates {
			if rt := atn.ruleTransitionAt(followIdx); rt != nil {
				seed.hints[i] = rt.CompletionHint
			}
		}
		l := atn.simulateFrom(seed, r.Tokens[s.TokenIdx:cursor], DefaultSimConfig)
		if len(l) == 0 {
			return nil, false
		}
		info := atn.NextCompletionsFromSet(l)
		if len(info.Tokens) == 0 && len(info.Hints) == 0 {
			return nil, false
		}
		return l, true
	}

	// Only snapshots within maxBacktrack tokens of the cursor are tried:
	// simulating from far back is expensive, and any syntax error in between
	// makes it fail anyway. The latest snapshot is always tried.
	const maxBacktrack = 32
	var snapshots []ATNSnapshot
	for i := len(r.Snapshots) - 1; i >= 0; i-- {
		snapshots = r.Snapshots[i:]
		if r.Snapshots[i].TokenIdx+maxBacktrack < cursor {
			break
		}
	}

	for _, s := range snapshots {
		if l, ok := try(s); ok {
			return l, s, true
		}
	}
	return nil, ATNSnapshot{}, false
}
