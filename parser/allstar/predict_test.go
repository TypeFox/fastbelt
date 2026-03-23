// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
)

// ────────────────────────────────────────────────────────────────────────────
// PredicateSet
// ────────────────────────────────────────────────────────────────────────────

func TestPredicateSet_IsOutOfBounds(t *testing.T) {
	p := &PredicateSet{}
	assert.True(t, p.Is(100), "out-of-bounds index should be unconstrained (true)")
}

func TestPredicateSet_IsInBounds(t *testing.T) {
	p := &PredicateSet{}
	p.Set(0, false)
	assert.False(t, p.Is(0))
	p.Set(0, true)
	assert.True(t, p.Is(0))
}

func TestPredicateSet_String(t *testing.T) {
	p := &PredicateSet{}
	p.Set(0, true)
	p.Set(1, false)
	p.Set(2, true)
	assert.Equal(t, "101", p.String())
}

func TestPredicateSet_Empty(t *testing.T) {
	assert.Equal(t, "", EmptyPredicates.String())
	assert.True(t, EmptyPredicates.Is(99))
}

// ────────────────────────────────────────────────────────────────────────────
// computeStartState
// ────────────────────────────────────────────────────────────────────────────

// buildTwoAltDecision constructs a minimal RuntimeATNState with two epsilon
// transitions to two basic states, representing a two-alternative decision.
func buildTwoAltDecision() *RuntimeATNState {
	decState := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, EpsilonOnlyTransitions: true}
	alt0 := &RuntimeATNState{StateNumber: 1, Type: ATNBasic}
	alt1 := &RuntimeATNState{StateNumber: 2, Type: ATNBasic}
	decState.Transitions = []RuntimeTransition{
		&RuntimeEpsilonTransition{Target: alt0},
		&RuntimeEpsilonTransition{Target: alt1},
	}
	return decState
}

func TestComputeStartState_Simple(t *testing.T) {
	dec := buildTwoAltDecision()
	configs := computeStartState(dec)
	// Each epsilon transition adds one config.
	assert.Equal(t, 2, configs.Len())
	alts := configs.Alts()
	assert.Contains(t, alts, 0)
	assert.Contains(t, alts, 1)
}

// ────────────────────────────────────────────────────────────────────────────
// closure
// ────────────────────────────────────────────────────────────────────────────

func TestClosure_Epsilon(t *testing.T) {
	s0 := &RuntimeATNState{StateNumber: 0, EpsilonOnlyTransitions: true}
	s1 := &RuntimeATNState{StateNumber: 1, EpsilonOnlyTransitions: true}
	s2 := &RuntimeATNState{StateNumber: 2}
	s0.Transitions = []RuntimeTransition{&RuntimeEpsilonTransition{Target: s1}}
	s1.Transitions = []RuntimeTransition{&RuntimeEpsilonTransition{Target: s2}}
	// s2 is not epsilon-only, so it gets added.
	s2.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: s0, TokenTypeID: 1}}

	configs := NewATNConfigSet()
	closure(&ATNConfig{State: s0, Alt: 0, Stack: []*RuntimeATNState{}}, configs)
	// Should reach s2 via two epsilons.
	found := false
	for _, c := range configs.Elements() {
		if c.State == s2 {
			found = true
		}
	}
	assert.True(t, found, "closure should follow epsilon transitions and add s2")
}

func TestClosure_RuleTransition(t *testing.T) {
	// Build a stub sub-rule: startRule → (inner) → innerEnd
	startRule := &RuntimeATNState{StateNumber: 0, Type: ATNRuleStart, EpsilonOnlyTransitions: true}
	inner := &RuntimeATNState{StateNumber: 1}
	innerEnd := &RuntimeATNState{StateNumber: 2}
	inner.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: innerEnd, TokenTypeID: 99}}
	startRule.Transitions = []RuntimeTransition{&RuntimeEpsilonTransition{Target: inner}}

	// Build call site: s0 --RuleTransition--> startRule; follow = s1
	s0 := &RuntimeATNState{StateNumber: 3, EpsilonOnlyTransitions: true}
	s1 := &RuntimeATNState{StateNumber: 4}
	s1End := &RuntimeATNState{StateNumber: 5}
	s1.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: s1End, TokenTypeID: 7}}
	s0.Transitions = []RuntimeTransition{&RuntimeRuleTransition{Target: startRule, FollowState: s1}}

	configs := NewATNConfigSet()
	closure(&ATNConfig{State: s0, Alt: 0, Stack: []*RuntimeATNState{}}, configs)

	// Should have followed into the sub-rule start, pushing s1 onto the stack.
	foundInner := false
	for _, c := range configs.Elements() {
		if c.State == inner {
			require.Len(t, c.Stack, 1)
			assert.Equal(t, s1, c.Stack[0])
			foundInner = true
		}
	}
	assert.True(t, foundInner)
}

func TestClosure_RuleStop_EmptyStack(t *testing.T) {
	stopState := &RuntimeATNState{StateNumber: 0, Type: ATNRuleStop}

	configs := NewATNConfigSet()
	config := &ATNConfig{State: stopState, Alt: 1, Stack: []*RuntimeATNState{}}
	closure(config, configs)

	// Empty stack at rule stop → add to configs directly.
	assert.Equal(t, 1, configs.Len())
	assert.Equal(t, stopState, configs.Elements()[0].State)
}

// ────────────────────────────────────────────────────────────────────────────
// getReachableTarget / tokenMatches
// ────────────────────────────────────────────────────────────────────────────

func TestGetReachableTarget_Match(t *testing.T) {
	target := &RuntimeATNState{StateNumber: 99}
	at := &RuntimeAtomTransition{Target: target, TokenTypeID: 5}
	result := getReachableTarget(at, 5)
	assert.Equal(t, target, result)
}

func TestGetReachableTarget_NoMatch(t *testing.T) {
	at := &RuntimeAtomTransition{Target: &RuntimeATNState{}, TokenTypeID: 5}
	assert.Nil(t, getReachableTarget(at, 6))
}

func TestGetReachableTarget_CategoryMatch(t *testing.T) {
	target := &RuntimeATNState{StateNumber: 1}
	at := &RuntimeAtomTransition{Target: target, TokenTypeID: 5, CategoryMatches: []int{10, 11}}
	assert.Equal(t, target, getReachableTarget(at, 10))
	assert.Nil(t, getReachableTarget(at, 12))
}

// ────────────────────────────────────────────────────────────────────────────
// computeReachSet
// ────────────────────────────────────────────────────────────────────────────

func TestComputeReachSet_SingleAlt(t *testing.T) {
	targetEnd := &RuntimeATNState{StateNumber: 1}
	target := &RuntimeATNState{StateNumber: 0}
	// Make target non-epsilon-only so it gets added by closure.
	target.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: targetEnd, TokenTypeID: 0}}

	src := &RuntimeATNState{StateNumber: 2}
	src.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: target, TokenTypeID: 7}}

	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: src, Alt: 0, Stack: []*RuntimeATNState{}})

	reach := computeReachSet(configs, 7, EmptyPredicates)
	assert.Equal(t, 1, reach.Len())
	assert.Equal(t, target, reach.Elements()[0].State)
}

func TestComputeReachSet_SkipsGatedAlt(t *testing.T) {
	targetEnd := &RuntimeATNState{StateNumber: 1}
	target := &RuntimeATNState{StateNumber: 0}
	target.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: targetEnd, TokenTypeID: 0}}

	src := &RuntimeATNState{StateNumber: 2}
	src.Transitions = []RuntimeTransition{&RuntimeAtomTransition{Target: target, TokenTypeID: 7}}

	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: src, Alt: 0, Stack: []*RuntimeATNState{}})

	// Gate alt 0 to false.
	preds := &PredicateSet{}
	preds.Set(0, false)

	reach := computeReachSet(configs, 7, preds)
	assert.Equal(t, 0, reach.Len(), "gated alt should be skipped")
}

// ────────────────────────────────────────────────────────────────────────────
// getUniqueAlt
// ────────────────────────────────────────────────────────────────────────────

func TestGetUniqueAlt_Unique(t *testing.T) {
	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: &RuntimeATNState{}, Alt: 2, Stack: []*RuntimeATNState{}})
	configs.Add(&ATNConfig{State: &RuntimeATNState{StateNumber: 1}, Alt: 2, Stack: []*RuntimeATNState{}})
	alt, ok := getUniqueAlt(configs, EmptyPredicates)
	assert.True(t, ok)
	assert.Equal(t, 2, alt)
}

func TestGetUniqueAlt_Mixed(t *testing.T) {
	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: &RuntimeATNState{}, Alt: 0, Stack: []*RuntimeATNState{}})
	configs.Add(&ATNConfig{State: &RuntimeATNState{StateNumber: 1}, Alt: 1, Stack: []*RuntimeATNState{}})
	_, ok := getUniqueAlt(configs, EmptyPredicates)
	assert.False(t, ok)
}

// ────────────────────────────────────────────────────────────────────────────
// allConfigsInRuleStopStates
// ────────────────────────────────────────────────────────────────────────────

func TestAllConfigsInRuleStopStates(t *testing.T) {
	stop := &RuntimeATNState{Type: ATNRuleStop}
	basic := &RuntimeATNState{Type: ATNBasic}

	allStop := NewATNConfigSet()
	allStop.Add(&ATNConfig{State: stop, Alt: 0, Stack: []*RuntimeATNState{}})
	assert.True(t, allConfigsInRuleStopStates(allStop))

	mixed := NewATNConfigSet()
	mixed.Add(&ATNConfig{State: stop, Alt: 0, Stack: []*RuntimeATNState{}})
	mixed.Add(&ATNConfig{State: basic, Alt: 1, Stack: []*RuntimeATNState{}})
	assert.False(t, allConfigsInRuleStopStates(mixed))
}

// ────────────────────────────────────────────────────────────────────────────
// hasConflictTerminatingPrediction
// ────────────────────────────────────────────────────────────────────────────

func TestHasConflictTerminatingPrediction_AllAtStop(t *testing.T) {
	stop := &RuntimeATNState{Type: ATNRuleStop}
	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: stop, Alt: 0, Stack: []*RuntimeATNState{}})
	configs.Add(&ATNConfig{State: stop, Alt: 1, Stack: []*RuntimeATNState{}})
	assert.True(t, hasConflictTerminatingPrediction(configs))
}

func TestHasConflictTerminatingPrediction_Conflicting(t *testing.T) {
	s := &RuntimeATNState{StateNumber: 5, Type: ATNBasic}
	configs := NewATNConfigSet()
	// Two configs with same state/stack but different alts → conflict.
	configs.Add(&ATNConfig{State: s, Alt: 0, Stack: []*RuntimeATNState{}})
	configs.Add(&ATNConfig{State: s, Alt: 1, Stack: []*RuntimeATNState{}})
	assert.True(t, hasConflictTerminatingPrediction(configs))
}

// ────────────────────────────────────────────────────────────────────────────
// getConflictingAltSets
// ────────────────────────────────────────────────────────────────────────────

func TestGetConflictingAltSets(t *testing.T) {
	s := &RuntimeATNState{StateNumber: 3, Type: ATNBasic}
	configs := []*ATNConfig{
		{State: s, Alt: 0, Stack: []*RuntimeATNState{}},
		{State: s, Alt: 1, Stack: []*RuntimeATNState{}},
	}
	altSets := getConflictingAltSets(configs)
	key := atnConfigKey(&ATNConfig{State: s, Alt: 0, Stack: []*RuntimeATNState{}}, false)
	alts, ok := altSets[key]
	assert.True(t, ok)
	assert.True(t, alts[0])
	assert.True(t, alts[1])
}

// ────────────────────────────────────────────────────────────────────────────
// addDFAState / addDFAEdge
// ────────────────────────────────────────────────────────────────────────────

func TestAddDFAState_Dedup(t *testing.T) {
	dfa := &DFA{States: map[string]*DFAState{}}
	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: &RuntimeATNState{StateNumber: 1}, Alt: 0, Stack: []*RuntimeATNState{}})
	state := newDFAState(configs)
	added1 := addDFAState(dfa, state)
	added2 := addDFAState(dfa, state)
	assert.Equal(t, added1, added2, "same config-set key should return same DFAState")
}

func TestAddDFAEdge(t *testing.T) {
	dfa := &DFA{States: map[string]*DFAState{}}

	fromConfigs := NewATNConfigSet()
	fromConfigs.Add(&ATNConfig{State: &RuntimeATNState{StateNumber: 1}, Alt: 0, Stack: []*RuntimeATNState{}})
	from := newDFAState(fromConfigs)

	toConfigs := NewATNConfigSet()
	toConfigs.Add(&ATNConfig{State: &RuntimeATNState{StateNumber: 2}, Alt: 1, Stack: []*RuntimeATNState{}})
	to := newDFAState(toConfigs)

	addDFAEdge(dfa, from, 5, to)
	assert.Equal(t, to, from.Edges[5])
}

// ────────────────────────────────────────────────────────────────────────────
// tokenTypeID helper
// ────────────────────────────────────────────────────────────────────────────

func TestTokenTypeID_Nil(t *testing.T) {
	assert.Equal(t, core.EOF.Id, tokenTypeID(nil))
}

func TestTokenTypeID_Token(t *testing.T) {
	tok := &core.Token{TypeId: 42}
	assert.Equal(t, 42, tokenTypeID(tok))
}
