package automatons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNFABuilderImpl_AddState(t *testing.T) {
	builder := NewNFABuilder()

	state0 := builder.AddState()
	state1 := builder.AddState()
	state2 := builder.AddState()

	assert.Equal(t, 0, state0)
	assert.Equal(t, 1, state1)
	assert.Equal(t, 2, state2)
}

func TestNFABuilderImpl_AddTransitionValidation(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()
	chars := NewRuneSetRune('a')

	assert.Panics(t, func() {
		builder.AddTransitionForRuneSet(-1, s1, chars)
	})

	assert.Panics(t, func() {
		builder.AddTransitionForRuneSet(2, s1, chars)
	})

	// Test invalid target state
	assert.Panics(t, func() {
		builder.AddTransitionForRuneSet(s0, -1, chars)
	})

	assert.Panics(t, func() {
		builder.AddTransitionForRuneSet(s0, 2, chars)
	})

	builder.AddTransitionForRuneSet(s0, s1, chars)
}

func TestNFABuilderImpl_AddTransitions(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()
	chars := NewRuneSetRune('a')

	builder.AddTransitionForRuneSet(s0, s1, chars)
	builder.SetStartState(s0)
	builder.AcceptState(s1)
	nfa := builder.Build()

	// Verify the transition exists
	transitionsBySource := nfa.TransitionsBySource
	targets, exists := transitionsBySource[s0]
	assert.True(t, exists, "Expected transitions for source state not found")

	// Check if the transition contains our character and target
	found := false
	for info := range targets.All() {
		if info.Range != nil && info.Range.Contains('a') {
			for _, target := range info.Values {
				if target == s1 {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}

	assert.True(t, found, "Expected transition from s0 to s1 with character 'a' not found")
}

func TestNFABuilderImpl_SetStartStateValidation(t *testing.T) {
	builder := NewNFABuilder()

	assert.Panics(t, func() {
		builder.SetStartState(0)
	})

	builder.AddState()

	assert.Panics(t, func() {
		builder.SetStartState(-1)
	})

	assert.Panics(t, func() {
		builder.SetStartState(1)
	})

	builder.SetStartState(0)
}

func TestNFABuilderImpl_AcceptStateValidation(t *testing.T) {
	builder := NewNFABuilder()

	// Test accepting state with no states
	assert.Panics(t, func() {
		builder.AcceptState(0)
	})

	builder.AddState()

	// Test invalid accepting state
	assert.Panics(t, func() {
		builder.AcceptState(-1)
	})

	assert.Panics(t, func() {
		builder.AcceptState(1)
	})

	builder.AcceptState(0)
}

func TestNFABuilderImpl_BuildValidation(t *testing.T) {
	builder := NewNFABuilder()

	assert.Panics(t, func() {
		builder.Build()
	})

	builder.AddState()

	assert.Panics(t, func() {
		builder.Build()
	})
}

func TestNFABuilderImpl_BuildValidNFA(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()

	builder.SetStartState(s0)
	builder.AcceptState(s1)

	chars := NewRuneSetRune('a')
	builder.AddTransitionForRuneSet(s0, s1, chars)

	nfa := builder.Build()

	assert.Equal(t, s0, nfa.StartState)
	assert.Equal(t, 2, nfa.StateCount)
	assert.True(t, nfa.AcceptingStates[s1])
	assert.NotNil(t, nfa.TransitionsBySource[s0])
}

func TestNFABuilderImpl_CopyFrom(t *testing.T) {
	// Build original NFA
	builder1 := NewNFABuilder()
	s0 := builder1.AddState()
	s1 := builder1.AddState()

	builder1.SetStartState(s0)
	builder1.AcceptState(s1)
	chars := NewRuneSetRune('b')
	builder1.AddTransitionForRuneSet(s0, s1, chars)
	nfa1 := builder1.Build()

	// Copy to new builder
	builder2 := NewNFABuilder()
	stateMapping := builder2.CopyFrom(nfa1)
	builder2.SetStartState(stateMapping.Start)
	for _, acc := range stateMapping.Acceptings {
		builder2.AcceptState(acc)
	}
	nfa2 := builder2.Build()

	assert.Equal(t, nfa1.StateCount, nfa2.StateCount)
	assert.Equal(t, len(nfa1.AcceptingStates), len(nfa2.AcceptingStates))
	assert.Equal(t, len(nfa1.TransitionsBySource), len(nfa2.TransitionsBySource))
}
