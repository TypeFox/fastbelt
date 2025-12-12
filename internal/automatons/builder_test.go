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

	// Test invalid source state
	err := builder.AddTransition(-1, s1, chars)
	if err == nil {
		t.Error("Expected error for invalid source state")
	}

	err = builder.AddTransition(2, s1, chars)
	if err == nil {
		t.Error("Expected error for source state out of range")
	}

	// Test invalid target state
	err = builder.AddTransition(s0, -1, chars)
	if err == nil {
		t.Error("Expected error for invalid target state")
	}

	err = builder.AddTransition(s0, 2, chars)
	if err == nil {
		t.Error("Expected error for target state out of range")
	}

	// Test valid transition
	err = builder.AddTransition(s0, s1, chars)
	if err != nil {
		t.Errorf("Expected no error for valid transition, got: %v", err)
	}
}

func TestNFABuilderImpl_AddTransitions(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()
	chars := NewRuneSet_Rune('a')

	err := builder.AddTransition(s0, s1, chars)
	if err != nil {
		t.Fatalf("Failed to add transition: %v", err)
	}

	err = builder.SetStartState(s0)
	if err != nil {
		t.Fatalf("Failed to set start state: %v", err)
	}

	err = builder.AcceptState(s1)
	if err != nil {
		t.Fatalf("Failed to accept state: %v", err)
	}

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build NFA: %v", err)
	}

	// Verify the transition exists
	transitionsBySource := nfa.GetTransitionsBySource()
	targets, exists := transitionsBySource[s0]
	if !exists {
		t.Fatal("No transitions found for source state")
	}

	// Check if the transition contains our character and target
	found := false
	for info := range targets.AllTransitions() {
		if info.CharRange != nil && info.CharRange.Contains('a') {
			for _, target := range info.Targets {
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

	// Test setting start state with no states
	err := builder.SetStartState(0)
	if err == nil {
		t.Error("Expected error when setting start state with no states")
	}

	builder.AddState()

	// Test invalid start state
	err = builder.SetStartState(-1)
	if err == nil {
		t.Error("Expected error for negative start state")
	}

	err = builder.SetStartState(1)
	if err == nil {
		t.Error("Expected error for start state out of range")
	}

	// Test valid start state
	err = builder.SetStartState(0)
	if err != nil {
		t.Errorf("Expected no error for valid start state, got: %v", err)
	}
}

func TestNFABuilderImpl_AcceptStateValidation(t *testing.T) {
	builder := NewNFABuilder()

	// Test accepting state with no states
	err := builder.AcceptState(0)
	if err == nil {
		t.Error("Expected error when accepting state with no states")
	}

	builder.AddState()

	// Test invalid accepting state
	err = builder.AcceptState(-1)
	if err == nil {
		t.Error("Expected error for negative accepting state")
	}

	err = builder.AcceptState(1)
	if err == nil {
		t.Error("Expected error for accepting state out of range")
	}

	// Test valid accepting state
	err = builder.AcceptState(0)
	if err != nil {
		t.Errorf("Expected no error for valid accepting state, got: %v", err)
	}
}

func TestNFABuilderImpl_BuildValidation(t *testing.T) {
	builder := NewNFABuilder()

	// Test building with no states
	_, err := builder.Build()
	if err == nil || err.Error() != "no states defined" {
		t.Errorf("Expected 'no states defined' error, got: %v", err)
	}

	builder.AddState()

	// Test building with no start state
	_, err = builder.Build()
	if err == nil || err.Error() != "no start state defined" {
		t.Errorf("Expected 'no start state defined' error, got: %v", err)
	}
}

func TestNFABuilderImpl_BuildValidNFA(t *testing.T) {
	builder := NewNFABuilder()
	s0 := builder.AddState()
	s1 := builder.AddState()

	err := builder.SetStartState(s0)
	if err != nil {
		t.Fatalf("Failed to set start state: %v", err)
	}

	err = builder.AcceptState(s1)
	if err != nil {
		t.Fatalf("Failed to accept state: %v", err)
	}

	chars := NewRuneSet_Rune('a')
	err = builder.AddTransition(s0, s1, chars)
	if err != nil {
		t.Fatalf("Failed to add transition: %v", err)
	}

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build NFA: %v", err)
	}

	if nfa.GetStartState() != s0 {
		t.Errorf("Expected start state %d, got %d", s0, nfa.GetStartState())
	}

	if nfa.GetStateCount() != 2 {
		t.Errorf("Expected state count 2, got %d", nfa.GetStateCount())
	}

	acceptingStates := nfa.GetAcceptingStates()
	if !acceptingStates[s1] {
		t.Error("Expected s1 to be an accepting state")
	}

	transitionsBySource := nfa.GetTransitionsBySource()
	if _, exists := transitionsBySource[s0]; !exists {
		t.Error("Expected transitions from s0")
	}
}

func TestNFABuilderImpl_CopyFrom(t *testing.T) {
	// Build original NFA
	builder1 := NewNFABuilder()
	s0 := builder1.AddState()
	s1 := builder1.AddState()

	err := builder1.SetStartState(s0)
	if err != nil {
		t.Fatalf("Failed to set start state: %v", err)
	}

	err = builder1.AcceptState(s1)
	if err != nil {
		t.Fatalf("Failed to accept state: %v", err)
	}

	chars := NewRuneSet_Rune('b')
	err = builder1.AddTransition(s0, s1, chars)
	if err != nil {
		t.Fatalf("Failed to add transition: %v", err)
	}

	nfa1, err := builder1.Build()
	if err != nil {
		t.Fatalf("Failed to build original NFA: %v", err)
	}

	// Copy to new builder
	builder2 := NewNFABuilder()
	stateMapping, err := builder2.CopyFrom(nfa1)
	if err != nil {
		t.Fatalf("Failed to copy NFA: %v", err)
	}

	err = builder2.SetStartState(stateMapping.Start)
	if err != nil {
		t.Fatalf("Failed to set start state in copy: %v", err)
	}

	for _, acc := range stateMapping.Acceptings {
		err = builder2.AcceptState(acc)
		if err != nil {
			t.Fatalf("Failed to accept state in copy: %v", err)
		}
	}

	nfa2, err := builder2.Build()
	if err != nil {
		t.Fatalf("Failed to build copied NFA: %v", err)
	}

	if nfa2.GetStateCount() != nfa1.GetStateCount() {
		t.Errorf("Expected copied NFA to have %d states, got %d",
			nfa1.GetStateCount(), nfa2.GetStateCount())
	}

	if len(nfa2.GetAcceptingStates()) != len(nfa1.GetAcceptingStates()) {
		t.Errorf("Expected copied NFA to have %d accepting states, got %d",
			len(nfa1.GetAcceptingStates()), len(nfa2.GetAcceptingStates()))
	}

	if len(nfa2.GetTransitionsBySource()) != len(nfa1.GetTransitionsBySource()) {
		t.Errorf("Expected copied NFA to have %d transition sources, got %d",
			len(nfa1.GetTransitionsBySource()), len(nfa2.GetTransitionsBySource()))
	}
}
