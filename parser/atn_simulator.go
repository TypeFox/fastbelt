// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"strconv"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

// SimConfig holds the configurable cost/loop limits for an ATN simulation.
// The defaults are conservative; callers needing different limits can
// instantiate their own copy.
type SimConfig struct {
	// MaxLiveSet caps the size of the live set at each step. The simulation
	// is an NFA exploration; pathological grammars can still grow it, so we
	// drop the lowest-priority entries (oldest insertion order) when the cap
	// is exceeded rather than refusing the request.
	MaxLiveSet int
	// MaxReturnStack caps the depth of the per-path return stack. The parser
	// itself has finite recursion in practice, so a generous default is fine.
	MaxReturnStack int
}

// DefaultSimConfig is the SimConfig used by Simulate/NextCompletionsFromSet
// when the convenience entry points (without explicit config) are called.
var DefaultSimConfig = SimConfig{
	MaxLiveSet:     512,
	MaxReturnStack: 64,
}

// simPath is one node in the NFA live set: a state index plus the return-stack
// of FollowState indices the simulator must continue at when leaving rules
// through RuleStop transitions. Identity is the triple (stateIdx, stack,
// hints): two paths at the same state with the same return stack but
// different hint contexts must stay distinct because the completion
// provider surfaces different items for each.
//
// hints is parallel to stack: hints[i] is the CompletionHint that became
// active when stack frame i was pushed, or nil if that push did not change
// the active hint. The currently active hint is the topmost non-nil entry,
// implemented as `activeHintAfter(p.stack, p.hints)` below.
type simPath struct {
	stateIdx int
	stack    []int
	hints    []*CompletionHint
}

func (p simPath) key() string {
	var b strings.Builder
	b.Grow(8 + 6*len(p.stack) + 8*len(p.hints))
	b.WriteString(strconv.Itoa(p.stateIdx))
	b.WriteByte('|')
	for _, s := range p.stack {
		b.WriteString(strconv.Itoa(s))
		b.WriteByte(',')
	}
	b.WriteByte('|')
	for _, h := range p.hints {
		if h != nil {
			b.WriteString(h.Field)
		}
		b.WriteByte(',')
	}
	return b.String()
}

// activeHint returns the topmost non-nil entry on the path's hint stack,
// or nil if no cross-reference rule call is currently in progress.
func (p simPath) activeHint() *CompletionHint {
	for i := len(p.hints) - 1; i >= 0; i-- {
		if p.hints[i] != nil {
			return p.hints[i]
		}
	}
	return nil
}

// Simulate advances an NFA-style live set from startStateIdx by consuming the
// given tokens. The returned slice contains the encoded keys of every distinct
// (state, return-stack) pair that survived consumption; pass it to
// NextCompletionsFromSet to obtain the valid next tokens + hints.
//
// The simulator is deliberately path-aware (return stacks for RuleTransition /
// RuleStop) so that completions inside deeply-nested rules resolve correctly.
// It does not apply LL(k) lookahead - every alternative consistent with the
// consumed tokens stays live. That is what we want for completion: the user
// hasn't picked a branch yet.
func (atn *RuntimeATN) Simulate(startStateIdx int, tokens []core.Token) []simPath {
	return atn.SimulateWithConfig(startStateIdx, tokens, DefaultSimConfig)
}

// SimulateWithConfig is Simulate with explicit cost limits.
func (atn *RuntimeATN) SimulateWithConfig(startStateIdx int, tokens []core.Token, cfg SimConfig) []simPath {
	if startStateIdx < 0 || startStateIdx >= len(atn.States) {
		return nil
	}
	// Seed the live set with the closure of the start state.
	live := atn.epsilonClosure([]simPath{{stateIdx: startStateIdx, stack: nil}}, cfg)
	for _, tok := range tokens {
		live = atn.advance(live, tok.Type, cfg)
		if len(live) == 0 {
			return nil
		}
	}
	return live
}

// TokenCompletion is one atom-transition emission: a TokenType that the
// simulator reports as valid at the cursor, along with the source ATN state
// index of the transition that produced it. Two cursors inside the same
// rule but at different grammar positions surface different ATNStateIdx
// values for the same TokenType - the completion contributor branches on
// this to identify the exact grammar element being completed.
type TokenCompletion struct {
	TokenType   *core.TokenType
	ATNStateIdx int
}

// HintCompletion is one rule-call hint emission. ATNStateIdx is the source
// ATN state of the rule transition (or the atom transition that inherited
// the hint from a surrounding rule call).
type HintCompletion struct {
	Hint        *CompletionHint
	ATNStateIdx int
}

// CompletionInfo bundles everything the completion provider needs from a
// single live-set walk: every atom-transition emission (TokenType + source
// ATN state), every cross-reference hint emission (CompletionHint + source
// ATN state), and the subset of TokenType IDs that appeared ONLY on hinted
// transitions in the closure. The completion provider uses HintedOnlyIDs
// to suppress hinted ID tokens from the token pass - the dispatch table
// supplies the real candidates with their actual names.
//
// All slices/maps are owned by the caller after this returns; the simulator
// keeps no internal state.
type CompletionInfo struct {
	Tokens        []TokenCompletion
	Hints         []HintCompletion
	HintedOnlyIDs collections.Set[int]
}

// HasToken reports whether tokenID appears in any TokenCompletion.
// O(len(Tokens)); use it for assertions and one-off lookups, not for
// hot-path filtering.
func (info CompletionInfo) HasToken(tokenID int) bool {
	for _, tc := range info.Tokens {
		if tc.TokenType != nil && tc.TokenType.Id == tokenID {
			return true
		}
	}
	return false
}

// HasHintField reports whether any HintCompletion carries the given
// CompletionHint.Field value.
func (info CompletionInfo) HasHintField(field string) bool {
	for _, h := range info.Hints {
		if h.Hint != nil && h.Hint.Field == field {
			return true
		}
	}
	return false
}

// NextCompletionsFromSet returns a CompletionInfo covering every atom
// transition reachable from the given live set's ε-closure. It is the
// single public entry point for inspecting what the simulator yielded;
// callers that only need set-membership use info.HasToken / HasHintField.
//
// Tokens and hints are deduplicated by (sourceStateIdx, Id) and
// (sourceStateIdx, Field) respectively - so the same TokenType emitted
// from two different grammar positions surfaces twice, letting the
// completion contributor distinguish them via ATNStateIdx.
//
// hintedOnlyIDs holds TokenType.Id values that appeared on at least one
// hinted (cross-reference) atom transition AND on no unhinted atom
// transition. The completion provider uses this to skip emitting an "ID"
// keyword item when the only valid use of ID at this position is a
// cross-reference whose candidates the dispatch table will produce.
func (atn *RuntimeATN) NextCompletionsFromSet(live []simPath) CompletionInfo {
	var tokenComps []TokenCompletion
	var hints []HintCompletion
	seenToken := make(collections.Set[[2]int])
	seenHint := make(collections.Set[struct {
		state int
		field string
	}])
	hinted := make(collections.Set[int])
	unhinted := make(collections.Set[int])
	closure := atn.epsilonClosure(live, DefaultSimConfig)
	for _, p := range closure {
		if p.stateIdx < 0 || p.stateIdx >= len(atn.States) {
			continue
		}
		state := atn.States[p.stateIdx]
		if state == nil {
			continue
		}
		pathHint := p.activeHint()
		for _, t := range state.Transitions {
			at, ok := t.(*RuntimeAtomTransition)
			if !ok || at.TokenType == nil {
				continue
			}
			tk := [2]int{p.stateIdx, at.TokenType.Id}
			if seenToken.Add(tk) {
				tokenComps = append(tokenComps, TokenCompletion{
					TokenType:   at.TokenType,
					ATNStateIdx: p.stateIdx,
				})
			}
			// Atoms inside a cross-reference rule call inherit the rule
			// call's hint when they don't carry one of their own. That
			// makes every token of a multi-token cross-reference text
			// (e.g. the ID + "." + ID atoms of an FQN reference) count as
			// hinted, so HintedOnlyIDs suppresses them from the token
			// pass and the cross-ref dispatch supplies the real candidates.
			effectiveHint := at.CompletionHint
			if effectiveHint == nil {
				effectiveHint = pathHint
			}
			if effectiveHint != nil {
				hinted.Add(at.TokenType.Id)
				hk := struct {
					state int
					field string
				}{p.stateIdx, effectiveHint.Field}
				if seenHint.Add(hk) {
					hints = append(hints, HintCompletion{
						Hint:        effectiveHint,
						ATNStateIdx: p.stateIdx,
					})
				}
			} else {
				unhinted.Add(at.TokenType.Id)
			}
		}
	}
	hintedOnly := make(collections.Set[int])
	for id := range hinted {
		if !unhinted.Has(id) {
			hintedOnly.Add(id)
		}
	}
	return CompletionInfo{
		Tokens:        tokenComps,
		Hints:         hints,
		HintedOnlyIDs: hintedOnly,
	}
}

// epsilonClosure expands the live set across every reachable epsilon-style
// transition: plain epsilons, rule entries (push FollowState), and rule exits
// (pop the top of the return stack). The result includes the seed states.
func (atn *RuntimeATN) epsilonClosure(seed []simPath, cfg SimConfig) []simPath {
	out := make([]simPath, 0, len(seed))
	seen := make(collections.Set[string], len(seed)*2)
	stack := make([]simPath, 0, len(seed))
	for _, p := range seed {
		k := p.key()
		if !seen.Add(k) {
			continue
		}
		out = append(out, p)
		stack = append(stack, p)
	}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.stateIdx < 0 || cur.stateIdx >= len(atn.States) {
			continue
		}
		state := atn.States[cur.stateIdx]
		if state == nil {
			continue
		}
		// A rule-stop state pops the top of the return stack (if any) and
		// continues at the FollowState recorded by the RuleTransition that
		// brought us here. The matching hint frame is popped alongside.
		if state.Type == ATNRuleStop && len(cur.stack) > 0 {
			top := cur.stack[len(cur.stack)-1]
			nextStack := append([]int{}, cur.stack[:len(cur.stack)-1]...)
			nextHints := append([]*CompletionHint{}, cur.hints[:len(cur.hints)-1]...)
			np := simPath{stateIdx: top, stack: nextStack, hints: nextHints}
			k := np.key()
			if seen.Add(k) {
				out = append(out, np)
				stack = append(stack, np)
			}
			// Stop states have no outgoing transitions of their own.
			continue
		}
		for _, t := range state.Transitions {
			switch tt := t.(type) {
			case *RuntimeEpsilonTransition:
				idx := atn.stateIndex(tt.Target)
				if idx < 0 {
					continue
				}
				np := simPath{stateIdx: idx, stack: cur.stack, hints: cur.hints}
				k := np.key()
				if !seen.Add(k) {
					continue
				}
				out = append(out, np)
				stack = append(stack, np)
			case *RuntimeRuleTransition:
				targetIdx := atn.stateIndex(tt.Target)
				followIdx := atn.stateIndex(tt.FollowState)
				if targetIdx < 0 {
					continue
				}
				if len(cur.stack) >= cfg.MaxReturnStack {
					continue
				}
				newStack := make([]int, len(cur.stack)+1)
				copy(newStack, cur.stack)
				newStack[len(cur.stack)] = followIdx
				newHints := make([]*CompletionHint, len(cur.hints)+1)
				copy(newHints, cur.hints)
				newHints[len(cur.hints)] = tt.CompletionHint
				np := simPath{stateIdx: targetIdx, stack: newStack, hints: newHints}
				k := np.key()
				if !seen.Add(k) {
					continue
				}
				out = append(out, np)
				stack = append(stack, np)
			}
		}
		if cfg.MaxLiveSet > 0 && len(out) > cfg.MaxLiveSet {
			// Trim from the end (most recently added) - this corresponds to
			// the deepest exploration in our DFS; the caller already has the
			// nearer-front paths, which dominate the keyword set.
			out = out[:cfg.MaxLiveSet]
			break
		}
	}
	return out
}

// advance consumes one token, returning the new live set.
func (atn *RuntimeATN) advance(live []simPath, tokenType *core.TokenType, cfg SimConfig) []simPath {
	closure := atn.epsilonClosure(live, cfg)
	next := make([]simPath, 0, len(closure))
	seen := make(collections.Set[string], len(closure))
	for _, p := range closure {
		if p.stateIdx < 0 || p.stateIdx >= len(atn.States) {
			continue
		}
		state := atn.States[p.stateIdx]
		if state == nil {
			continue
		}
		for _, t := range state.Transitions {
			at, ok := t.(*RuntimeAtomTransition)
			if !ok || at.TokenType == nil || !at.TokenType.Matches(tokenType) {
				continue
			}
			targetIdx := atn.stateIndex(at.Target)
			if targetIdx < 0 {
				continue
			}
			np := simPath{stateIdx: targetIdx, stack: p.stack, hints: p.hints}
			k := np.key()
			if !seen.Add(k) {
				continue
			}
			next = append(next, np)
		}
	}
	return next
}
