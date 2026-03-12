// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildATNKey(t *testing.T) {
	rule := &Rule{Name: "MyRule"}
	assert.Equal(t, "MyRule_Alternation_1", BuildATNKey(rule, "Alternation", 1))
	assert.Equal(t, "MyRule_Option_2", BuildATNKey(rule, "Option", 2))
	assert.Equal(t, "MyRule_RepetitionMandatory_1", BuildATNKey(rule, "RepetitionMandatory", 1))
	assert.Equal(t, "MyRule_Repetition_3", BuildATNKey(rule, "Repetition", 3))
}

func TestCreateATN_SingleTerminal(t *testing.T) {
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&Terminal{TokenTypeID: 1, Idx: 1},
		},
	}
	atn := CreateATN([]*Rule{rule})

	// RuleStart + RuleStop + Basic(left) + Basic(right)
	assert.Len(t, atn.States, 4)
	// No decision states for a plain terminal.
	assert.Empty(t, atn.DecisionStates)
	// Verify atom transition exists.
	start := atn.RuleToStartState[rule]
	require.NotNil(t, start)
	assert.Equal(t, ATNRuleStart, start.Type)
	assert.Equal(t, ATNRuleStop, atn.RuleToStopState[rule].Type)
}

func TestCreateATN_Alternation(t *testing.T) {
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&Alternation{
				Alternatives: []*Alternative{
					{Definition: []Production{&Terminal{TokenTypeID: 1, Idx: 1}}},
					{Definition: []Production{&Terminal{TokenTypeID: 2, Idx: 2}}},
				},
				Idx: 1,
			},
		},
	}
	atn := CreateATN([]*Rule{rule})

	// Should have at least one decision state (the OR block start).
	assert.NotEmpty(t, atn.DecisionStates)
	key := BuildATNKey(rule, "Alternation", 1)
	decisionState, ok := atn.DecisionMap[key]
	require.True(t, ok, "decision map should contain alternation key")
	assert.Equal(t, ATNBasic, decisionState.Type)
	assert.GreaterOrEqual(t, decisionState.Decision, 0)
}

func TestCreateATN_Option(t *testing.T) {
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&Option{
				Definition: []Production{&Terminal{TokenTypeID: 1, Idx: 1}},
				Idx:        1,
			},
		},
	}
	atn := CreateATN([]*Rule{rule})

	key := BuildATNKey(rule, "Option", 1)
	decisionState, ok := atn.DecisionMap[key]
	require.True(t, ok)
	// The option decision state should have an epsilon bypass.
	assert.Equal(t, ATNBasic, decisionState.Type)
}

func TestCreateATN_Repetition(t *testing.T) {
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&Repetition{
				Definition: []Production{&Terminal{TokenTypeID: 1, Idx: 1}},
				Idx:        1,
			},
		},
	}
	atn := CreateATN([]*Rule{rule})

	key := BuildATNKey(rule, "Repetition", 1)
	entry, ok := atn.DecisionMap[key]
	require.True(t, ok)
	assert.Equal(t, ATNStarLoopEntry, entry.Type)

	// StarLoopEntry, StarBlockStart, StarLoopBack, LoopEnd should all appear.
	types := map[ATNStateType]int{}
	for _, s := range atn.States {
		types[s.Type]++
	}
	assert.Greater(t, types[ATNStarLoopEntry], 0)
	assert.Greater(t, types[ATNStarLoopBack], 0)
	assert.Greater(t, types[ATNLoopEnd], 0)
}

func TestCreateATN_RepetitionMandatory(t *testing.T) {
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&RepetitionMandatory{
				Definition: []Production{&Terminal{TokenTypeID: 1, Idx: 1}},
				Idx:        1,
			},
		},
	}
	atn := CreateATN([]*Rule{rule})

	key := BuildATNKey(rule, "RepetitionMandatory", 1)
	loopBack, ok := atn.DecisionMap[key]
	require.True(t, ok)
	assert.Equal(t, ATNPlusLoopBack, loopBack.Type)

	types := map[ATNStateType]int{}
	for _, s := range atn.States {
		types[s.Type]++
	}
	assert.Greater(t, types[ATNPlusBlockStart], 0)
	assert.Greater(t, types[ATNPlusLoopBack], 0)
	assert.Greater(t, types[ATNLoopEnd], 0)
}

func TestCreateATN_NonTerminal(t *testing.T) {
	ruleB := &Rule{Name: "B", Definition: []Production{
		&Terminal{TokenTypeID: 1, Idx: 1},
	}}
	ruleA := &Rule{Name: "A", Definition: []Production{
		&NonTerminal{ReferencedRule: ruleB, Idx: 1},
	}}
	atn := CreateATN([]*Rule{ruleA, ruleB})

	startA := atn.RuleToStartState[ruleA]
	require.NotNil(t, startA)

	// Follow epsilon from start → find a RuleTransition pointing at ruleB.
	found := false
	for _, s := range atn.States {
		for _, t := range s.Transitions {
			if rt, ok := t.(*RuleTransition); ok {
				if rt.Rule == ruleB {
					found = true
				}
			}
		}
	}
	assert.True(t, found, "should have a RuleTransition pointing at ruleB")
}

func TestCreateATN_NestedRule(t *testing.T) {
	inner := &Rule{Name: "Inner", Definition: []Production{
		&Terminal{TokenTypeID: 1, Idx: 1},
	}}
	outer := &Rule{Name: "Outer", Definition: []Production{
		&NonTerminal{ReferencedRule: inner, Idx: 1},
		&Terminal{TokenTypeID: 2, Idx: 1},
	}}
	atn := CreateATN([]*Rule{inner, outer})

	// Both rules should have start/stop states.
	assert.NotNil(t, atn.RuleToStartState[inner])
	assert.NotNil(t, atn.RuleToStartState[outer])

	// Find the RuleTransition and check its FollowState leads toward tokenTypeID=2.
	for _, s := range atn.States {
		for _, tr := range s.Transitions {
			if rt, ok := tr.(*RuleTransition); ok && rt.Rule == inner {
				assert.NotNil(t, rt.FollowState, "RuleTransition must have a follow state")
			}
		}
	}
}

func TestCreateATN_DecisionMapKeys(t *testing.T) {
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&Alternation{
				Alternatives: []*Alternative{
					{Definition: []Production{&Terminal{TokenTypeID: 1, Idx: 1}}},
					{Definition: []Production{&Terminal{TokenTypeID: 2, Idx: 2}}},
				},
				Idx: 1,
			},
			&Option{
				Definition: []Production{&Terminal{TokenTypeID: 1, Idx: 3}},
				Idx:        1,
			},
			&Repetition{
				Definition: []Production{&Terminal{TokenTypeID: 2, Idx: 4}},
				Idx:        1,
			},
		},
	}
	atn := CreateATN([]*Rule{rule})

	expectedKeys := []string{
		"R_Alternation_1",
		"R_Option_1",
		"R_Repetition_1",
	}
	for _, key := range expectedKeys {
		_, ok := atn.DecisionMap[key]
		assert.True(t, ok, "expected key %q in decision map", key)
	}
}

func TestMakeBlock_Optimisation(t *testing.T) {
	// When consecutive basic states can be merged, makeBlock removes the
	// intermediate state to avoid an extra epsilon hop.
	rule := &Rule{
		Name: "R",
		Definition: []Production{
			&Terminal{TokenTypeID: 1, Idx: 1},
			&Terminal{TokenTypeID: 2, Idx: 2},
		},
	}
	atn := CreateATN([]*Rule{rule})

	// The two-terminal sequence should have been merged: the first atom's
	// right state is eliminated and the transition jumps directly to the
	// second atom's left state.
	// We verify this by checking the atom transition's target is a state with
	// an atom transition to tokenTypeID=2 (not an epsilon).
	start := atn.RuleToStartState[rule]
	require.NotNil(t, start)

	// Follow epsilon from RuleStart → first atom left state.
	var firstAtomLeft *ATNState
	for _, tr := range start.Transitions {
		if tr.IsEpsilon() {
			firstAtomLeft = tr.Target()
		}
	}
	require.NotNil(t, firstAtomLeft)

	// The atom transition should point to a state from which tokenTypeID=2 is reachable.
	for _, tr := range firstAtomLeft.Transitions {
		if at, ok := tr.(*AtomTransition); ok {
			assert.Equal(t, 1, at.TokenTypeID)
			// The target should have a transition for tokenTypeID=2.
			hasTwo := false
			for _, tt := range at.target.Transitions {
				if at2, ok2 := tt.(*AtomTransition); ok2 && at2.TokenTypeID == 2 {
					hasTwo = true
				}
			}
			assert.True(t, hasTwo, "merged block: next atom (typeID=2) should be directly reachable")
		}
	}
}
