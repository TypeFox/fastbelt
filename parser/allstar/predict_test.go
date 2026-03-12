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

func buildTwoAltDecision() *ATNState {
	// Build a minimal decision state with two epsilon transitions to basic states.
	atn := &ATN{States: []*ATNState{}, DecisionStates: []*ATNState{}, DecisionMap: map[string]*ATNState{}}
	rule := &Rule{Name: "R"}
	decState := newATNState(atn, rule, nil, ATNBasic)
	alt0 := newATNState(atn, rule, nil, ATNBasic)
	alt1 := newATNState(atn, rule, nil, ATNBasic)
	addEpsilon(decState, alt0)
	addEpsilon(decState, alt1)
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
	atn := &ATN{States: []*ATNState{}, DecisionMap: map[string]*ATNState{}}
	rule := &Rule{Name: "R"}
	s0 := newATNState(atn, rule, nil, ATNBasic)
	s1 := newATNState(atn, rule, nil, ATNBasic)
	s2 := newATNState(atn, rule, nil, ATNBasic)
	addEpsilon(s0, s1)
	addEpsilon(s1, s2)
	// s2 is not epsilon-only, so it gets added.
	addTransition(s2, &AtomTransition{target: s0, TokenTypeID: 1})

	configs := NewATNConfigSet()
	closure(&ATNConfig{State: s0, Alt: 0, Stack: []*ATNState{}}, configs)
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
	atn := &ATN{States: []*ATNState{}, DecisionMap: map[string]*ATNState{}}
	rule := &Rule{Name: "R"}

	// Build a stub sub-rule: start → (atom) → stop
	startRule := newATNState(atn, rule, nil, ATNRuleStart)
	inner := newATNState(atn, rule, nil, ATNBasic)
	addTransition(inner, &AtomTransition{target: newATNState(atn, rule, nil, ATNBasic), TokenTypeID: 99})
	addEpsilon(startRule, inner)

	// Build call site: s0 --RuleTransition--> startRule; follow = s1
	s0 := newATNState(atn, rule, nil, ATNBasic)
	s1 := newATNState(atn, rule, nil, ATNBasic)
	addTransition(s1, &AtomTransition{target: newATNState(atn, rule, nil, ATNBasic), TokenTypeID: 7})
	addTransition(s0, &RuleTransition{target: startRule, Rule: rule, FollowState: s1})

	configs := NewATNConfigSet()
	closure(&ATNConfig{State: s0, Alt: 0, Stack: []*ATNState{}}, configs)

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
	atn := &ATN{States: []*ATNState{}, DecisionMap: map[string]*ATNState{}}
	rule := &Rule{Name: "R"}
	stopState := newATNState(atn, rule, nil, ATNRuleStop)

	configs := NewATNConfigSet()
	config := &ATNConfig{State: stopState, Alt: 1, Stack: []*ATNState{}}
	closure(config, configs)

	// Empty stack at rule stop → add to configs directly.
	assert.Equal(t, 1, configs.Len())
	assert.Equal(t, stopState, configs.Elements()[0].State)
}

// ────────────────────────────────────────────────────────────────────────────
// getReachableTarget / tokenMatches
// ────────────────────────────────────────────────────────────────────────────

func TestGetReachableTarget_Match(t *testing.T) {
	target := &ATNState{StateNumber: 99}
	at := &AtomTransition{target: target, TokenTypeID: 5}
	result := getReachableTarget(at, 5)
	assert.Equal(t, target, result)
}

func TestGetReachableTarget_NoMatch(t *testing.T) {
	at := &AtomTransition{target: &ATNState{}, TokenTypeID: 5}
	assert.Nil(t, getReachableTarget(at, 6))
}

func TestGetReachableTarget_CategoryMatch(t *testing.T) {
	target := &ATNState{StateNumber: 1}
	at := &AtomTransition{target: target, TokenTypeID: 5, CategoryMatches: []int{10, 11}}
	assert.Equal(t, target, getReachableTarget(at, 10))
	assert.Nil(t, getReachableTarget(at, 12))
}

// ────────────────────────────────────────────────────────────────────────────
// computeReachSet
// ────────────────────────────────────────────────────────────────────────────

func TestComputeReachSet_SingleAlt(t *testing.T) {
	atn := &ATN{States: []*ATNState{}, DecisionMap: map[string]*ATNState{}}
	rule := &Rule{Name: "R"}
	target := newATNState(atn, rule, nil, ATNBasic)
	// Make target non-epsilon-only so it gets added by closure.
	addTransition(target, &AtomTransition{target: newATNState(atn, rule, nil, ATNBasic), TokenTypeID: 0})

	src := newATNState(atn, rule, nil, ATNBasic)
	addTransition(src, &AtomTransition{target: target, TokenTypeID: 7})

	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: src, Alt: 0, Stack: []*ATNState{}})

	reach := computeReachSet(configs, 7, EmptyPredicates)
	assert.Equal(t, 1, reach.Len())
	assert.Equal(t, target, reach.Elements()[0].State)
}

func TestComputeReachSet_SkipsGatedAlt(t *testing.T) {
	atn := &ATN{States: []*ATNState{}, DecisionMap: map[string]*ATNState{}}
	rule := &Rule{Name: "R"}
	target := newATNState(atn, rule, nil, ATNBasic)
	addTransition(target, &AtomTransition{target: newATNState(atn, rule, nil, ATNBasic), TokenTypeID: 0})

	src := newATNState(atn, rule, nil, ATNBasic)
	addTransition(src, &AtomTransition{target: target, TokenTypeID: 7})

	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: src, Alt: 0, Stack: []*ATNState{}})

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
	configs.Add(&ATNConfig{State: &ATNState{}, Alt: 2, Stack: []*ATNState{}})
	configs.Add(&ATNConfig{State: &ATNState{StateNumber: 1}, Alt: 2, Stack: []*ATNState{}})
	alt, ok := getUniqueAlt(configs, EmptyPredicates)
	assert.True(t, ok)
	assert.Equal(t, 2, alt)
}

func TestGetUniqueAlt_Mixed(t *testing.T) {
	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: &ATNState{}, Alt: 0, Stack: []*ATNState{}})
	configs.Add(&ATNConfig{State: &ATNState{StateNumber: 1}, Alt: 1, Stack: []*ATNState{}})
	_, ok := getUniqueAlt(configs, EmptyPredicates)
	assert.False(t, ok)
}

// ────────────────────────────────────────────────────────────────────────────
// allConfigsInRuleStopStates
// ────────────────────────────────────────────────────────────────────────────

func TestAllConfigsInRuleStopStates(t *testing.T) {
	stop := &ATNState{Type: ATNRuleStop}
	basic := &ATNState{Type: ATNBasic}

	allStop := NewATNConfigSet()
	allStop.Add(&ATNConfig{State: stop, Alt: 0, Stack: []*ATNState{}})
	assert.True(t, allConfigsInRuleStopStates(allStop))

	mixed := NewATNConfigSet()
	mixed.Add(&ATNConfig{State: stop, Alt: 0, Stack: []*ATNState{}})
	mixed.Add(&ATNConfig{State: basic, Alt: 1, Stack: []*ATNState{}})
	assert.False(t, allConfigsInRuleStopStates(mixed))
}

// ────────────────────────────────────────────────────────────────────────────
// hasConflictTerminatingPrediction
// ────────────────────────────────────────────────────────────────────────────

func TestHasConflictTerminatingPrediction_AllAtStop(t *testing.T) {
	stop := &ATNState{Type: ATNRuleStop}
	configs := NewATNConfigSet()
	configs.Add(&ATNConfig{State: stop, Alt: 0, Stack: []*ATNState{}})
	configs.Add(&ATNConfig{State: stop, Alt: 1, Stack: []*ATNState{}})
	assert.True(t, hasConflictTerminatingPrediction(configs))
}

func TestHasConflictTerminatingPrediction_Conflicting(t *testing.T) {
	s := &ATNState{StateNumber: 5, Type: ATNBasic}
	configs := NewATNConfigSet()
	// Two configs with same state/stack but different alts → conflict.
	configs.Add(&ATNConfig{State: s, Alt: 0, Stack: []*ATNState{}})
	configs.Add(&ATNConfig{State: s, Alt: 1, Stack: []*ATNState{}})
	assert.True(t, hasConflictTerminatingPrediction(configs))
}

// ────────────────────────────────────────────────────────────────────────────
// getConflictingAltSets
// ────────────────────────────────────────────────────────────────────────────

func TestGetConflictingAltSets(t *testing.T) {
	s := &ATNState{StateNumber: 3, Type: ATNBasic}
	configs := []*ATNConfig{
		{State: s, Alt: 0, Stack: []*ATNState{}},
		{State: s, Alt: 1, Stack: []*ATNState{}},
	}
	altSets := getConflictingAltSets(configs)
	key := atnConfigKey(&ATNConfig{State: s, Alt: 0, Stack: []*ATNState{}}, false)
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
	configs.Add(&ATNConfig{State: &ATNState{StateNumber: 1}, Alt: 0, Stack: []*ATNState{}})
	state := newDFAState(configs)
	added1 := addDFAState(dfa, state)
	added2 := addDFAState(dfa, state)
	assert.Equal(t, added1, added2, "same config-set key should return same DFAState")
}

func TestAddDFAEdge(t *testing.T) {
	dfa := &DFA{States: map[string]*DFAState{}}

	fromConfigs := NewATNConfigSet()
	fromConfigs.Add(&ATNConfig{State: &ATNState{StateNumber: 1}, Alt: 0, Stack: []*ATNState{}})
	from := newDFAState(fromConfigs)

	toConfigs := NewATNConfigSet()
	toConfigs.Add(&ATNConfig{State: &ATNState{StateNumber: 2}, Alt: 1, Stack: []*ATNState{}})
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
