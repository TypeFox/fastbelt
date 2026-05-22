// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"context"
	"iter"

	core "typefox.dev/fastbelt"
)

// LanguageCompletionAdapter bridges the framework's generic completion
// provider and the per-language generated artefacts (CompletionParser,
// SyntheticFactories, CompletionDispatch). The interface lives in this
// package because it is purely a parser-level contract - none of its
// methods touch LSP or server-side types - and so the generated language
// package can implement it without pulling in the LSP server module.
//
// Adopters do not typically write this implementation by hand; the code
// generator emits one alongside the rest of the generated files.
type LanguageCompletionAdapter interface {
	// Parse runs the generated CompletionParser over the given prefix
	// tokens (the document's tokens up to the cursor) and returns the
	// snapshots + rule stack the completion provider needs.
	Parse(tokens []core.Token) *CompletionParseResult
	// ATN returns the RuntimeATN the simulator queries.
	ATN() *RuntimeATN
	// SyntheticOwnerFor returns a fresh, detached AST node for the given
	// parser-rule key (e.g. "Transition"). Returns (nil, false) for
	// unknown keys; the completion provider aborts cross-reference
	// completion when it can't build the synthetic chain.
	SyntheticOwnerFor(ruleKey string) (core.AstNode, bool)
	// DispatchCompletion runs the per-field scope+filter chain for the
	// given CompletionHint.Field on the supplied synthetic owner. Returns
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
}

// CompletionParseResult is the output of a CompletionParser run. It contains
// every artefact the completion provider needs to drive the ATN simulator and
// build the synthetic-owner chain.
type CompletionParseResult struct {
	// Tokens is the token slice the parser was given (typically the document's
	// prefix tokens up to the cursor).
	Tokens []core.Token
	// NextTokenIndex is the index of the next unconsumed token at the moment
	// the parser stopped. For a complete prefix this equals len(Tokens); for a
	// bailed parse it's the index of the offending token.
	NextTokenIndex int
	// Snapshots is the time-ordered sequence of (tokenIdx, atnStateIdx) pairs
	// recorded at every rule entry and every Sync. The completion provider
	// picks the latest entry with TokenIdx <= cursor as the simulator's start.
	Snapshots []ATNSnapshot
	// RuleStack is the chain of rule contexts the parser was inside when it
	// stopped, outermost first. It drives synthetic-owner construction.
	RuleStack []RuleContext
}

// ATNSnapshot records "at TokenIdx tokens consumed, we were positioned at ATN
// state ATNStateIdx, and the parser was inside RuleStack". Snapshots are
// recorded by the generated CompletionParser at every interesting point (rule
// entry, Sync) - between snapshots, the simulator advances the live set by
// consuming tokens.
//
// RuleStack captures the live rule context at the moment the snapshot was
// taken. It is what makes the completion provider's synthetic-owner chain
// robust against error recovery: by the time Parse returns, the live stack
// has been popped back to empty, but each snapshot preserves the context the
// parser was in when it was recorded.
type ATNSnapshot struct {
	TokenIdx    int
	ATNStateIdx int
	RuleStack   []RuleContext
}

// RuleContext is a single frame of CompletionParseResult.RuleStack. RuleKey
// is a stable identifier the generator emits per parser rule (typically the
// rule's exported Go name, e.g. "Statemachine"); the synthetic-owner factory
// table is keyed by this string. Assignment is the property name currently
// being assigned (empty when the parser is not inside an assignment).
type RuleContext struct {
	RuleKey    string
	Assignment string
}

// CompletionParserState is the embedded helper a generated CompletionParser
// uses to record snapshots and maintain its rule stack. It wraps a regular
// ParserState (which still does the consume/lookahead/sync work) and adds the
// bookkeeping completion needs.
//
// The generator emits, per rule:
//
//	func (p *MyCompletionParser) ParseFoo() {
//	    p.cp.EnterRule("Foo", FooRuleStartStateIdx)
//	    defer p.cp.ExitRule()
//	    // ...mirrors the main parser's control flow, but every Sync(i) is
//	    //  preceded by p.cp.RecordSnapshot(i), every assignment site is
//	    //  preceded by p.cp.MarkAssignment("PropertyName"), and every
//	    //  AST-mutation call (AssignToken, SetName, etc.) is omitted.
//	}
type CompletionParserState struct {
	state *ParserState

	snapshots []ATNSnapshot
	ruleStack []RuleContext
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

// EnterRule pushes a RuleContext onto the rule stack and snapshots the
// (currentTokenIdx, ruleStartStateIdx, ruleStack) triple. The generator emits
// one call at the start of every rule function.
//
// While the parser is in error mode, EnterRule still pushes a RuleContext
// (so the rule stack reflects what the generated code is doing) but does not
// record a snapshot - snapshots should describe actual parsing progress, and
// while InError is set there is no progress to capture. With a recovering
// strategy InError is normally cleared again at the next ExitRule/Sync, and
// snapshots resume at that point.
func (cp *CompletionParserState) EnterRule(ruleKey string, ruleStartStateIdx int) {
	cp.ruleStack = append(cp.ruleStack, RuleContext{RuleKey: ruleKey})
	if cp.state.ErrorMode != ErrorModeNone {
		return
	}
	cp.snapshots = append(cp.snapshots, ATNSnapshot{
		TokenIdx:    cp.state.Index,
		ATNStateIdx: ruleStartStateIdx,
		RuleStack:   append([]RuleContext(nil), cp.ruleStack...),
	})
}

// ExitRule pops the top RuleContext. The matching push happens in EnterRule;
// the generator emits one call (via defer) at the start of every rule
// function.
//
// ExitRule always pops, even while the parser is in error mode. The cursor's
// rule context is preserved by ATNSnapshot.RuleStack - each snapshot copies
// the live stack at the moment it was recorded, so unwinding (whether
// recovery-driven or EOF-driven) doesn't destroy the context the completion
// provider needs. Always popping keeps the live stack balanced under
// recovery, which is essential when the parser continues past a transient
// error and enters more rules.
func (cp *CompletionParserState) ExitRule() {
	if len(cp.ruleStack) > 0 {
		cp.ruleStack = cp.ruleStack[:len(cp.ruleStack)-1]
	}
}

// RecordSnapshot snapshots the current token index alongside the given ATN
// state index and live rule stack. The generator emits one call right before
// every p.state.Sync(idx) call so the simulator has a known starting point at
// every branch boundary.
//
// While the parser is in error mode, RecordSnapshot is a no-op: parsing has
// stopped advancing, so further snapshots would only point past the actual
// cursor and confuse FindSnapshotAt.
func (cp *CompletionParserState) RecordSnapshot(atnStateIdx int) {
	if cp.state.ErrorMode != ErrorModeNone {
		return
	}
	cp.snapshots = append(cp.snapshots, ATNSnapshot{
		TokenIdx:    cp.state.Index,
		ATNStateIdx: atnStateIdx,
		RuleStack:   append([]RuleContext(nil), cp.ruleStack...),
	})
}

// MarkAssignment sets the assignment property on the top RuleContext. The
// generator emits this immediately before each Consume that corresponds to an
// assignment (`Property=...` in the grammar). The assignment value is cleared
// automatically at the next EnterRule/ExitRule boundary; callers that need to
// clear it earlier can call ClearAssignment.
func (cp *CompletionParserState) MarkAssignment(property string) {
	if len(cp.ruleStack) == 0 {
		return
	}
	cp.ruleStack[len(cp.ruleStack)-1].Assignment = property
}

// ClearAssignment resets the top frame's Assignment field.
func (cp *CompletionParserState) ClearAssignment() {
	if len(cp.ruleStack) == 0 {
		return
	}
	cp.ruleStack[len(cp.ruleStack)-1].Assignment = ""
}

// Result builds a CompletionParseResult snapshot. The generator's top-level
// Parse method calls this once after the entry rule returns.
//
// The returned RuleStack is the deepest rule stack recorded at the latest
// snapshot's token index. This effectively answers "what rules did the
// parser enter at the cursor's token position?" - by the time Parse returns
// the live stack has unwound to empty, but the snapshots preserve the
// context at each point of progress.
func (cp *CompletionParserState) Result(tokens []core.Token) *CompletionParseResult {
	snapshots := append([]ATNSnapshot(nil), cp.snapshots...)
	return &CompletionParseResult{
		Tokens:         tokens,
		NextTokenIndex: cp.state.Index,
		Snapshots:      snapshots,
		RuleStack:      deepestRuleStack(snapshots),
	}
}

// deepestRuleStack returns the deepest RuleStack among the snapshots whose
// TokenIdx equals the maximum TokenIdx seen. This is the cursor's rule
// context: the latest token position the parser reached, expanded to the
// most-nested rule the parser was in at that position.
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
// trying candidate snapshots in priority order until one yields a non-empty
// result. Returns the live set and the snapshot the simulator started from.
//
// Snapshots are tried EARLIEST FIRST. The earliest snapshot is the broadest
// context - the parser hadn't yet committed to any LL(k) sub-decision - so
// simulating forward from it through the prefix tokens lets the NFA keep
// every alternative still consistent with the typed input alive. A narrower
// (later) snapshot would reflect a single branch the parser committed to,
// and would miss completions from the sibling branches that are still
// reachable from the typed prefix.
//
// Later snapshots are used as fallbacks: when the prefix is unparseable
// from an earlier snapshot (the parser bailed mid-input and the broad
// simulation can't consume past the bad token), a later snapshot taken
// closer to the cursor may still yield a useful result.
//
// "Yields a non-empty result" means both that the simulator advanced
// through the prefix without losing the live set AND that
// NextCompletionsFromSet finds reachable atom transitions; the latter
// rules out snapshots that land at a RuleStop with an empty return stack
// (the parser-committed-too-deep case).
//
// Callers that only want the cursor's snapshot without simulating should
// use FindSnapshotAt instead.
func (r *CompletionParseResult) SimulateAt(atn *RuntimeATN, cursor int) (live []simPath, snap ATNSnapshot, ok bool) {
	if len(r.Snapshots) == 0 || cursor < 0 {
		// No snapshots means the parser never entered any rule - not even the entry rule
		return nil, ATNSnapshot{}, false
	}
	// try simulates the ATN from the specified snapshot
	// If it yields a non-empty completion set, that likely means that the snapshot is valid (no syntax errors in the prefix)
	// If it fails, we should try the next snapshot, which is closer to the cursor, but might be more specific (less broad)
	try := func(s ATNSnapshot) ([]simPath, bool) {
		l := atn.Simulate(s.ATNStateIdx, r.Tokens[s.TokenIdx:cursor])
		if len(l) == 0 {
			return nil, false
		}
		info := atn.NextCompletionsFromSet(l)
		if len(info.Tokens) == 0 && len(info.Hints) == 0 {
			return nil, false
		}
		return l, true
	}

	const maxBacktrack = 32
	var snapshots []ATNSnapshot
	// Find latest snapshots which are within reach of the cursor
	for i := len(r.Snapshots) - 1; i >= 0; i-- {
		snapshots = r.Snapshots[i:]
		s := r.Snapshots[i]
		if s.TokenIdx+maxBacktrack < cursor {
			// Prevent simulating from snapshots that are too far back from the cursor
			// Theoretically, we could simulate from the first snapshot (the broadest context)
			// However, for very long inputs, this could be very expensive
			// And even then, a single syntax error in the prefix will prevent any useful output
			// Note that the simulator will always try at least one snapshot, even if it's far back
			break
		}
	}

	// Simulate from each snapshot (broadest first)
	for _, s := range snapshots {
		if l, ok := try(s); ok {
			return l, s, true
		}
	}
	return nil, ATNSnapshot{}, false
}

// FindSnapshotAt returns the snapshot that gives the broadest context at the
// given cursor token index. Selection rules:
//
//   - If any snapshot has TokenIdx == cursor, return the EARLIEST such
//     snapshot. Multiple snapshots at the same token index represent the
//     parser making successive branch decisions at that position; the
//     earliest reflects the parser BEFORE those decisions, so simulating
//     from there exposes every alternative still consistent with the input
//     consumed so far.
//   - Otherwise (cursor lies strictly between two snapshots, or past all of
//     them), return the LATEST snapshot with TokenIdx < cursor. The
//     simulator advances forward from that snapshot through the intervening
//     tokens.
//
// Returns (ATNSnapshot{}, false) if no snapshot fits - which only happens
// for an empty result (no rule was ever entered).
func (r *CompletionParseResult) FindSnapshotAt(cursor int) (ATNSnapshot, bool) {
	var best ATNSnapshot
	found := false
	for _, s := range r.Snapshots {
		if s.TokenIdx > cursor {
			break
		}
		if s.TokenIdx == cursor {
			// First snapshot AT the cursor wins; later ones at the same idx
			// represent narrower contexts we don't want to commit to.
			if !found || best.TokenIdx < cursor {
				best = s
				found = true
			}
			continue
		}
		// s.TokenIdx < cursor: track the latest such snapshot.
		best = s
		found = true
	}
	return best, found
}
