package automatons

import (
	"testing"
)

func TestNFABuilderImpl_AddState(t *testing.T) {
	builder := NewNFABuilder()

	state0 := builder.AddState()
	state1 := builder.AddState()
	state2 := builder.AddState()

	if state0 != 0 {
		t.Errorf("Expected first state to be 0, got %d", state0)
	}
	if state1 != 1 {
		t.Errorf("Expected second state to be 1, got %d", state1)
	}
	if state2 != 2 {
		t.Errorf("Expected third state to be 2, got %d", state2)
	}
}

func TestNFABuilderImpl_AddTransitionValidation(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()
	chars := NewRuneSet_Rune('a')

	Expect(func() {
		builder.AddTransitionForRuneSet(-1, s1, chars)
	}).ToPanic()

	Expect(func() {
		builder.AddTransitionForRuneSet(2, s1, chars)
	}).ToPanic()

	// Test invalid target state
	Expect(func() {
		builder.AddTransitionForRuneSet(s0, -1, chars)
	}).ToPanic()

	Expect(func() {
		builder.AddTransitionForRuneSet(s0, 2, chars)
	}).ToPanic()

	builder.AddTransitionForRuneSet(s0, s1, chars)
}

func TestNFABuilderImpl_AddTransitions(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()
	chars := NewRuneSet_Rune('a')

	builder.AddTransitionForRuneSet(s0, s1, chars)
	builder.SetStartState(s0)
	builder.AcceptState(s1)
	nfa := builder.Build()

	// Verify the transition exists
	transitionsBySource := nfa.TransitionsBySource
	targets, exists := transitionsBySource[s0]
	if !exists {
		t.Fatal("No transitions found for source state")
	}

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

	if !found {
		t.Error("Expected transition from s0 to s1 with character 'a' not found")
	}
}

func TestNFABuilderImpl_SetStartStateValidation(t *testing.T) {
	builder := NewNFABuilder()

	Expect(func() {
		builder.SetStartState(0)
	}).ToPanic()

	builder.AddState()

	Expect(func() {
		builder.SetStartState(-1)
	}).ToPanic()

	Expect(func() {
		builder.SetStartState(1)
	}).ToPanic()

	builder.SetStartState(0)
}

func TestNFABuilderImpl_AcceptStateValidation(t *testing.T) {
	builder := NewNFABuilder()

	// Test accepting state with no states
	Expect(func() {
		builder.AcceptState(0)
	}).ToPanic()

	builder.AddState()

	// Test invalid accepting state
	Expect(func() {
		builder.AcceptState(-1)
	}).ToPanic()

	Expect(func() {
		builder.AcceptState(1)
	}).ToPanic()

	builder.AcceptState(0)
}

func TestNFABuilderImpl_BuildValidation(t *testing.T) {
	builder := NewNFABuilder()

	Expect(func() {
		builder.Build()
	}).ToPanic()

	builder.AddState()

	Expect(func() {
		builder.Build()
	}).ToPanic()
}

func TestNFABuilderImpl_BuildValidNFA(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()

	builder.SetStartState(s0)
	builder.AcceptState(s1)

	chars := NewRuneSet_Rune('a')
	builder.AddTransitionForRuneSet(s0, s1, chars)

	nfa := builder.Build()

	Expect(nfa.StartState).ToEqual(s0)
	Expect(nfa.StateCount).ToEqual(2)
	Expect(bool(nfa.AcceptingStates[s1])).ToEqual(true)
	Expect(nfa.TransitionsBySource[s0] != nil).ToEqual(true)
}

func TestNFABuilderImpl_CopyFrom(t *testing.T) {
	// Build original NFA
	builder1 := NewNFABuilder()
	s0 := builder1.AddState()
	s1 := builder1.AddState()

	builder1.SetStartState(s0)
	builder1.AcceptState(s1)
	chars := NewRuneSet_Rune('b')
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

	Expect(nfa2.StateCount).ToEqual(nfa1.StateCount)
	Expect(len(nfa2.AcceptingStates)).ToEqual(len(nfa1.AcceptingStates))
	Expect(len(nfa2.TransitionsBySource)).ToEqual(len(nfa1.TransitionsBySource))
}
