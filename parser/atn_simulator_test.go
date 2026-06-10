// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"testing"

	core "typefox.dev/fastbelt"
)

// tok constructs a token of the given type for tests.
func tok(tt *core.TokenType) core.Token {
	return core.Token{Type: tt, TypeId: tt.Id, Image: tt.Name}
}

// linearATN builds an ATN with two basic states connected by a single atom
// transition on tokenType. The atom transition is the only transition on the
// source state; the target has no transitions. Useful as a small leaf rule.
func linearATN(t *testing.T, tokenType *core.TokenType) *RuntimeATN {
	t.Helper()
	s0 := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, Decision: -1}
	s1 := &RuntimeATNState{StateNumber: 1, Type: ATNBasic, Decision: -1}
	s0.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: s1, TokenType: tokenType},
	}
	atn := NewRuntimeATN([]*RuntimeATNState{s0, s1}, nil, nil)
	return atn
}

func createTokenType(id int, name string) *core.TokenType {
	return core.NewTokenType(id, name, name, 0, 0, 0, false, nil, nil)
}

// TestSimulator_AlternativeCoverage builds an ATN for `rule: 'a' 'x' | 'a' 'y'`
// and asserts that simulating consumption of 'a' from the entry state yields a
// live set whose NextCompletionsFromSet contains BOTH 'x' and 'y'. This is the
// regression test for the user's primary concern: the parser would have
// committed to one branch via LL(k) lookahead, but the simulator must keep
// every alternative consistent with the consumed input live.
func TestSimulator_AlternativeCoverage(t *testing.T) {
	tA := createTokenType(1, "a")
	tX := createTokenType(2, "x")
	tY := createTokenType(3, "y")

	// States:
	//   0 -ε-> 1 (alt A)         1 -a-> 2 -x-> 3
	//   0 -ε-> 4 (alt B)         4 -a-> 5 -y-> 6
	s0 := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, Decision: -1, EpsilonOnlyTransitions: true}
	s1 := &RuntimeATNState{StateNumber: 1, Type: ATNBasic, Decision: -1}
	s2 := &RuntimeATNState{StateNumber: 2, Type: ATNBasic, Decision: -1}
	s3 := &RuntimeATNState{StateNumber: 3, Type: ATNBasic, Decision: -1}
	s4 := &RuntimeATNState{StateNumber: 4, Type: ATNBasic, Decision: -1}
	s5 := &RuntimeATNState{StateNumber: 5, Type: ATNBasic, Decision: -1}
	s6 := &RuntimeATNState{StateNumber: 6, Type: ATNBasic, Decision: -1}
	s0.Transitions = []RuntimeTransition{
		&RuntimeEpsilonTransition{Target: s1},
		&RuntimeEpsilonTransition{Target: s4},
	}
	s1.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: s2, TokenType: tA},
	}
	s2.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: s3, TokenType: tX},
	}
	s4.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: s5, TokenType: tA},
	}
	s5.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: s6, TokenType: tY},
	}
	atn := NewRuntimeATN([]*RuntimeATNState{s0, s1, s2, s3, s4, s5, s6}, nil, nil)

	live := atn.Simulate(0, []core.Token{tok(tA)})
	if len(live) == 0 {
		t.Fatalf("simulator returned empty live set after consuming 'a'")
	}
	info := atn.NextCompletionsFromSet(live)
	if !info.HasToken(tX.Id) {
		t.Errorf("expected 'x' to be in valid next tokens; got %v", info.Tokens)
	}
	if !info.HasToken(tY.Id) {
		t.Errorf("expected 'y' to be in valid next tokens; got %v", info.Tokens)
	}
}

// TestSimulator_MidGroupAdvance builds an ATN for `rule: 'a' 'b' 'c'` and
// asserts that after consuming "a b" the simulator yields a live set whose
// NextCompletionsFromSet returns exactly {'c'}. This is the regression for the
// reviewer's concern that Sync isn't called between consecutive consumes:
// without the simulator advancing the live set token-by-token, the answer
// would wrongly be the rule's FIRST set.
func TestSimulator_MidGroupAdvance(t *testing.T) {
	tA := createTokenType(1, "a")
	tB := createTokenType(2, "b")
	tC := createTokenType(3, "c")

	s0 := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, Decision: -1}
	s1 := &RuntimeATNState{StateNumber: 1, Type: ATNBasic, Decision: -1}
	s2 := &RuntimeATNState{StateNumber: 2, Type: ATNBasic, Decision: -1}
	s3 := &RuntimeATNState{StateNumber: 3, Type: ATNBasic, Decision: -1}
	s0.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: s1, TokenType: tA}}
	s1.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: s2, TokenType: tB}}
	s2.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: s3, TokenType: tC}}
	atn := NewRuntimeATN([]*RuntimeATNState{s0, s1, s2, s3}, nil, nil)

	live := atn.Simulate(0, []core.Token{tok(tA), tok(tB)})
	if len(live) == 0 {
		t.Fatalf("simulator returned empty live set after consuming 'a b'")
	}
	info := atn.NextCompletionsFromSet(live)
	if !info.HasToken(tC.Id) {
		t.Errorf("expected 'c' to be in valid next tokens; got %v", info.Tokens)
	}
	if info.HasToken(tA.Id) {
		t.Errorf("'a' should NOT be in valid next tokens after consuming 'a b'; the simulator did not advance past the start state")
	}
}

// TestSimulator_RuleCallReturnsToFollow exercises the return-stack semantics:
// outer rule calls inner rule via RuleTransition, inner consumes one token,
// then RuleStop must pop and continue at the outer FollowState which expects
// a second token. The expected next-token after the inner call must be the
// outer's "after-inner" token, not anything from the inner rule.
func TestSimulator_RuleCallReturnsToFollow(t *testing.T) {
	tInner := createTokenType(1, "inner")
	tOuter := createTokenType(2, "outer")

	// Inner rule: ruleStart -ε-> innerCall -inner-> ruleStop
	innerStart := &RuntimeATNState{StateNumber: 10, Type: ATNRuleStart, Decision: -1, EpsilonOnlyTransitions: true}
	innerMid := &RuntimeATNState{StateNumber: 11, Type: ATNBasic, Decision: -1}
	innerEnd := &RuntimeATNState{StateNumber: 12, Type: ATNBasic, Decision: -1}
	innerStop := &RuntimeATNState{StateNumber: 13, Type: ATNRuleStop, Decision: -1}
	innerStart.Transitions = []RuntimeTransition{&RuntimeEpsilonTransition{Target: innerMid}}
	innerMid.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: innerEnd, TokenType: tInner}}
	innerEnd.Transitions = []RuntimeTransition{&RuntimeEpsilonTransition{Target: innerStop}}

	// Outer rule: start -RuleTransition(target=innerStart, follow=outerFollow)-> ...
	//   outerFollow -outer-> outerEnd
	outerStart := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, Decision: -1, EpsilonOnlyTransitions: true}
	outerFollow := &RuntimeATNState{StateNumber: 1, Type: ATNBasic, Decision: -1}
	outerEnd := &RuntimeATNState{StateNumber: 2, Type: ATNBasic, Decision: -1}
	outerStart.Transitions = []RuntimeTransition{
		&RuntimeRuleTransition{Target: innerStart, FollowState: outerFollow},
	}
	outerFollow.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: outerEnd, TokenType: tOuter},
	}

	atn := NewRuntimeATN(
		[]*RuntimeATNState{outerStart, outerFollow, outerEnd, innerStart, innerMid, innerEnd, innerStop},
		nil, nil,
	)

	// Before consuming anything, only 'inner' is valid (first token of inner rule).
	live := atn.Simulate(0, nil)
	info := atn.NextCompletionsFromSet(live)
	if !info.HasToken(tInner.Id) {
		t.Errorf("expected 'inner' before any consumption; got %v", info.Tokens)
	}
	if info.HasToken(tOuter.Id) {
		t.Errorf("did NOT expect 'outer' before consuming 'inner'; got %v", info.Tokens)
	}

	// After consuming 'inner', the simulator must pop the return stack and
	// expose 'outer' from the FollowState.
	live = atn.Simulate(0, []core.Token{tok(tInner)})
	if len(live) == 0 {
		t.Fatalf("simulator returned empty live set after consuming 'inner'")
	}
	info = atn.NextCompletionsFromSet(live)
	if !info.HasToken(tOuter.Id) {
		t.Errorf("expected 'outer' after consuming 'inner' (return-stack pop); got %v", info.Tokens)
	}
	if info.HasToken(tInner.Id) {
		t.Errorf("did NOT expect 'inner' to remain valid after consuming it; got %v", info.Tokens)
	}
}

// TestSimulator_HintsSurface verifies that CompletionHints attached to atom
// transitions are returned alongside the token bitset. Hints are deduplicated
// by Field value.
func TestSimulator_HintsSurface(t *testing.T) {
	tID := createTokenType(1, "ID")
	hint := &CompletionHint{Field: "Transition.Event"}

	s0 := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, Decision: -1}
	s1 := &RuntimeATNState{StateNumber: 1, Type: ATNBasic, Decision: -1}
	s0.Transitions = []RuntimeTransition{
		&RuntimeAtomTransition{Target: s1, TokenType: tID, CompletionHint: hint},
	}
	atn := NewRuntimeATN([]*RuntimeATNState{s0, s1}, nil, nil)

	live := atn.Simulate(0, nil)
	info := atn.NextCompletionsFromSet(live)
	if !info.HasToken(tID.Id) {
		t.Fatalf("expected ID token valid; got %v", info.Tokens)
	}
	if len(info.Hints) != 1 || info.Hints[0].Hint.Field != "Transition.Event" {
		t.Errorf("expected exactly one hint with Field='Transition.Event'; got %#v", info.Hints)
	}
}

// TestSimulator_EmptyLiveSetOnMismatch confirms that consuming a token with no
// matching transition produces an empty live set (and that subsequent advances
// short-circuit).
func TestSimulator_EmptyLiveSetOnMismatch(t *testing.T) {
	tA := createTokenType(1, "a")
	tB := createTokenType(2, "b")
	atn := linearATN(t, tA)

	// Consuming 'b' is invalid at the start.
	live := atn.Simulate(0, []core.Token{tok(tB)})
	if len(live) != 0 {
		t.Errorf("expected empty live set after mismatched token; got %d entries", len(live))
	}
}
