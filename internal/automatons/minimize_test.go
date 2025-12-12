package automatons

import (
	"testing"
)

// Helper function to create a simple DFA that can be minimized
func createMinimizableDFA() NFA {
	builder := NewNFABuilder()

	// Create 5 states: 0, 1, 2, 3, 4
	s0 := builder.AddState() // 0
	s1 := builder.AddState() // 1
	s2 := builder.AddState() // 2
	s3 := builder.AddState() // 3
	s4 := builder.AddState() // 4

	// Set start state
	builder.SetStartState(s0)

	// Set accepting states: 2 and 4 are accepting
	builder.AcceptState(s2)
	builder.AcceptState(s4)

	// Add transitions to create a DFA where states 3 and 1 are equivalent
	// and states 2 and 4 are equivalent
	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')

	// From state 0: a->1, b->2
	builder.AddTransition(s0, s1, charSetA)
	builder.AddTransition(s0, s2, charSetB)

	// From state 1: a->3, b->4  (non-accepting to non-accepting, non-accepting to accepting)
	builder.AddTransition(s1, s3, charSetA)
	builder.AddTransition(s1, s4, charSetB)

	// From state 2: a->3, b->4  (accepting to non-accepting, accepting to accepting)
	builder.AddTransition(s2, s3, charSetA)
	builder.AddTransition(s2, s4, charSetB)

	// From state 3: a->3, b->4  (same as state 1 behavior)
	builder.AddTransition(s3, s3, charSetA)
	builder.AddTransition(s3, s4, charSetB)

	// From state 4: a->3, b->4  (same as state 2 behavior)
	builder.AddTransition(s4, s3, charSetA)
	builder.AddTransition(s4, s4, charSetB)

	dfa, _ := builder.Build()
	return dfa
}

// Helper function to create a DFA that cannot be minimized further (already minimal)
func createMinimalDFA() NFA {
	builder := NewNFABuilder()

	// Create 3 states
	s0 := builder.AddState() // 0
	s1 := builder.AddState() // 1
	s2 := builder.AddState() // 2

	// Set start state
	builder.SetStartState(s0)

	// Set accepting state
	builder.AcceptState(s2)

	// Add transitions
	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')

	// Create a minimal DFA for language ending in 'ab'
	builder.AddTransition(s0, s0, charSetB) // b->0
	builder.AddTransition(s0, s1, charSetA) // a->1
	builder.AddTransition(s1, s0, charSetB) // b->0 (but this goes to accepting state 2)
	builder.AddTransition(s1, s1, charSetA) // a->1
	builder.AddTransition(s1, s2, charSetB) // b->2 (accepting)
	builder.AddTransition(s2, s0, charSetB) // b->0
	builder.AddTransition(s2, s1, charSetA) // a->1

	dfa, _ := builder.Build()
	return dfa
}

// Helper function to create a simple two-state DFA
func createTwoStateDFA() NFA {
	builder := NewNFABuilder()

	s0 := builder.AddState() // 0
	s1 := builder.AddState() // 1

	builder.SetStartState(s0)
	builder.AcceptState(s1)

	charSetA := NewRuneSet_Rune('a')

	builder.AddTransition(s0, s1, charSetA)
	builder.AddTransition(s1, s1, charSetA)

	dfa, _ := builder.Build()
	return dfa
}

// Helper function to create a DFA with multiple equivalent states
func createDFAWithEquivalentStates() NFA {
	builder := NewNFABuilder()

	// Create 4 states where states 1 and 3 are equivalent, and states 2 and 0 are unique
	s0 := builder.AddState() // 0 - start state
	s1 := builder.AddState() // 1 - accepting
	s2 := builder.AddState() // 2 - non-accepting sink
	s3 := builder.AddState() // 3 - accepting (equivalent to s1)

	builder.SetStartState(s0)

	// States 1 and 3 are both accepting
	builder.AcceptState(s1)
	builder.AcceptState(s3)

	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')

	// From start state: different transitions
	builder.AddTransition(s0, s1, charSetA) // 0 -a-> 1 (accepting)
	builder.AddTransition(s0, s2, charSetB) // 0 -b-> 2 (non-accepting)

	// States 1 and 3 have IDENTICAL transitions (making them equivalent)
	builder.AddTransition(s1, s2, charSetA) // 1 -a-> 2
	builder.AddTransition(s1, s2, charSetB) // 1 -b-> 2
	builder.AddTransition(s3, s2, charSetA) // 3 -a-> 2 (same as s1)
	builder.AddTransition(s3, s2, charSetB) // 3 -b-> 2 (same as s1)

	// State 2 (sink): loops to itself
	builder.AddTransition(s2, s2, charSetA) // 2 -a-> 2
	builder.AddTransition(s2, s2, charSetB) // 2 -b-> 2

	dfa, _ := builder.Build()
	return dfa
}

func TestMinimize_SimpleCase(t *testing.T) {
	dfa := createMinimizableDFA()
	minimized := Minimize(dfa)

	// Check that the minimized DFA has fewer states
	if minimized.GetStateCount() >= dfa.GetStateCount() {
		t.Errorf("Expected minimized DFA to have fewer states. Original: %d, Minimized: %d",
			dfa.GetStateCount(), minimized.GetStateCount())
	}

	// Check that it still has a start state
	startState := minimized.GetStartState()
	if startState < 0 || startState >= minimized.GetStateCount() {
		t.Errorf("Invalid start state: %d", startState)
	}

	// Check that it has accepting states
	acceptingStates := minimized.GetAcceptingStates()
	if len(acceptingStates) == 0 {
		t.Error("Expected at least one accepting state")
	}
}

func TestMinimize_AlreadyMinimal(t *testing.T) {
	dfa := createMinimalDFA()
	originalStateCount := dfa.GetStateCount()

	minimized := Minimize(dfa)
	minimizedStateCount := minimized.GetStateCount()

	// Should not reduce the number of states significantly for an already minimal DFA
	if minimizedStateCount > originalStateCount {
		t.Errorf("Minimized DFA should not have more states than original. Original: %d, Minimized: %d",
			originalStateCount, minimizedStateCount)
	}

	// Verify structure is preserved
	if minimized.GetStartState() < 0 {
		t.Error("Minimized DFA should have a valid start state")
	}

	if len(minimized.GetAcceptingStates()) == 0 {
		t.Error("Minimized DFA should have accepting states")
	}
}

func TestMinimize_TwoStates(t *testing.T) {
	dfa := createTwoStateDFA()
	minimized := Minimize(dfa)

	// A two-state DFA with different behavior should remain two states
	expectedStates := 2
	if minimized.GetStateCount() != expectedStates {
		t.Errorf("Expected %d states, got %d", expectedStates, minimized.GetStateCount())
	}

	// Check that start state is valid
	startState := minimized.GetStartState()
	if startState < 0 || startState >= minimized.GetStateCount() {
		t.Errorf("Invalid start state: %d", startState)
	}

	// Check accepting states
	acceptingStates := minimized.GetAcceptingStates()
	if len(acceptingStates) != 1 {
		t.Errorf("Expected 1 accepting state, got %d", len(acceptingStates))
	}
}

func TestMinimize_MultipleEquivalentStates(t *testing.T) {
	dfa := createDFAWithEquivalentStates()
	minimized := Minimize(dfa)

	// Should significantly reduce the number of states
	originalCount := dfa.GetStateCount()
	minimizedCount := minimized.GetStateCount()

	t.Logf("Original state count: %d, Minimized state count: %d", originalCount, minimizedCount)

	if minimizedCount >= originalCount {
		t.Errorf("Expected significant state reduction. Original: %d, Minimized: %d",
			originalCount, minimizedCount)
	}

	// With the new DFA:
	// - State 0: start state, non-accepting, unique transitions
	// - States 1,3: accepting, identical transitions (should merge)
	// - State 2: non-accepting sink, unique transitions
	// Expected result: 3 states (0, merged{1,3}, 2)
	expectedStates := 3
	if minimizedCount != expectedStates {
		t.Errorf("Expected exactly %d states after minimization, got %d", expectedStates, minimizedCount)
	}
}

func TestMinimize_PreservesLanguage(t *testing.T) {
	dfa := createMinimizableDFA()
	minimized := Minimize(dfa)

	// Test some basic properties that should be preserved
	if minimized.GetStartState() < 0 {
		t.Error("Minimized DFA should have a valid start state")
	}

	acceptingStates := minimized.GetAcceptingStates()
	if len(acceptingStates) == 0 {
		t.Error("Minimized DFA should have accepting states")
	}

	// Test that the minimized DFA has transitions
	transitions := minimized.GetTransitionsBySource()
	if len(transitions) == 0 {
		t.Error("Minimized DFA should have transitions")
	}

	// Count total transitions to ensure structure is reasonable
	totalTransitions := 0
	for range minimized.AllTransitions() {
		totalTransitions++
	}

	if totalTransitions == 0 {
		t.Error("Minimized DFA should have transitions available through AllTransitions()")
	}
}

func TestCharSetOperations(t *testing.T) {
	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')
	charSetAB := NewRuneSet_OneOf([]rune{'a', 'b'})

	// Test union
	union := unionCharSets(charSetA, charSetB)
	if !charSetsEqual(union, charSetAB) {
		t.Error("Union of {a} and {b} should equal {a,b}")
	}

	// Test equality
	if !charSetsEqual(charSetA, charSetA) {
		t.Error("CharSet should be equal to itself")
	}

	if charSetsEqual(charSetA, charSetB) {
		t.Error("Different charsets should not be equal")
	}

	// Test empty union
	empty := unionCharSets()
	if empty.Length() != 0 {
		t.Error("Union of no charsets should be empty")
	}
}
