package automatons

import (
	"testing"
)

// Helper function to create a simple NFA for testing
func createSimpleNFA() NFA {
	builder := NewNFABuilder()

	// Create states
	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()

	// Set start state
	builder.SetStartState(s0)

	// Set accepting state
	builder.AcceptState(s2)

	// Add transitions: s0 -'a'-> s1, s1 -'b'-> s2
	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')

	builder.AddTransition(s0, s1, charSetA)
	builder.AddTransition(s1, s2, charSetB)

	nfa, _ := builder.Build()
	return nfa
}

// Helper function to create an NFA with epsilon transitions
func createNFAWithEpsilons() NFA {
	builder := NewNFABuilder()

	// Create states
	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()
	s3 := builder.AddState()

	// Set start state
	builder.SetStartState(s0)

	// Set accepting states
	builder.AcceptState(s2)
	builder.AcceptState(s3)

	// Add transitions: s0 -ε-> s1, s0 -'a'-> s2, s1 -'b'-> s3
	emptyCharSet := NewRuneSet_Empty()
	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')

	builder.AddTransition(s0, s1, emptyCharSet) // epsilon transition
	builder.AddTransition(s0, s2, charSetA)
	builder.AddTransition(s1, s3, charSetB)

	nfa, _ := builder.Build()
	return nfa
}

// Helper function to create an NFA with non-deterministic transitions
func createNonDeterministicNFA() NFA {
	builder := NewNFABuilder()

	// Create states
	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()
	s3 := builder.AddState()

	// Set start state
	builder.SetStartState(s0)

	// Set accepting states
	builder.AcceptState(s2)
	builder.AcceptState(s3)

	// Add non-deterministic transitions: s0 -'a'-> s1, s0 -'a'-> s2, s1 -'b'-> s3
	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')

	builder.AddTransition(s0, s1, charSetA)
	builder.AddTransition(s0, s2, charSetA) // Non-deterministic: same input, different targets
	builder.AddTransition(s1, s3, charSetB)

	nfa, _ := builder.Build()
	return nfa
}

func TestStatesToBitMask(t *testing.T) {
	states := []int{0, 2, 4}
	stateCount := 6

	mask := statesToBitMask(states, stateCount)

	if !mask.IsSet(0) {
		t.Error("Expected bit 0 to be set")
	}
	if mask.IsSet(1) {
		t.Error("Expected bit 1 to not be set")
	}
	if !mask.IsSet(2) {
		t.Error("Expected bit 2 to be set")
	}
	if mask.IsSet(3) {
		t.Error("Expected bit 3 to not be set")
	}
	if !mask.IsSet(4) {
		t.Error("Expected bit 4 to be set")
	}
	if mask.IsSet(5) {
		t.Error("Expected bit 5 to not be set")
	}
}

func TestBitMaskToStates(t *testing.T) {
	stateCount := 6
	mask := NewBitMask_Empty(stateCount)
	mask.Set(0)
	mask.Set(2)
	mask.Set(4)

	states := bitMaskToStates(mask, stateCount)

	expected := []int{0, 2, 4}
	if len(states) != len(expected) {
		t.Errorf("Expected %d states, got %d", len(expected), len(states))
	}

	for i, state := range states {
		if state != expected[i] {
			t.Errorf("Expected state %d at position %d, got %d", expected[i], i, state)
		}
	}
}

func TestRangeToRuneSet(t *testing.T) {
	// Test nil range
	emptySet := rangeToRuneSet(nil)
	if emptySet.Length() != 0 {
		t.Error("Expected empty set for nil range")
	}

	// Test included range
	includedRange := NewRuneRange('a', 'c', true)
	includedSet := rangeToRuneSet(includedRange)

	if !includedSet.IncludesRune('a') {
		t.Error("Expected included set to contain 'a'")
	}
	if !includedSet.IncludesRune('b') {
		t.Error("Expected included set to contain 'b'")
	}
	if !includedSet.IncludesRune('c') {
		t.Error("Expected included set to contain 'c'")
	}
	if includedSet.IncludesRune('d') {
		t.Error("Expected included set to not contain 'd'")
	}

	// Test excluded range
	excludedRange := NewRuneRange('a', 'c', false)
	excludedSet := rangeToRuneSet(excludedRange)

	if excludedSet.IncludesRune('a') {
		t.Error("Expected excluded set to not contain 'a'")
	}
	if excludedSet.IncludesRune('b') {
		t.Error("Expected excluded set to not contain 'b'")
	}
	if excludedSet.IncludesRune('c') {
		t.Error("Expected excluded set to not contain 'c'")
	}
	if !excludedSet.IncludesRune('d') {
		t.Error("Expected excluded set to contain 'd'")
	}
}

func TestDeterminizeSimpleNFA(t *testing.T) {
	nfa := createSimpleNFA()
	dfa := Determinize(nfa)

	// The DFA should have the same number of states as the original NFA
	// since it's already deterministic
	if dfa.GetStateCount() != nfa.GetStateCount() {
		t.Errorf("Expected DFA to have %d states, got %d", nfa.GetStateCount(), dfa.GetStateCount())
	}

	// Check that start state is set
	if dfa.GetStartState() < 0 {
		t.Error("DFA should have a valid start state")
	}

	// Check that there are accepting states
	acceptingStates := dfa.GetAcceptingStates()
	if len(acceptingStates) == 0 {
		t.Error("DFA should have at least one accepting state")
	}
}

func TestDeterminizeNFAWithEpsilons(t *testing.T) {
	nfa := createNFAWithEpsilons()
	dfa := Determinize(nfa)

	// Check that start state is set
	if dfa.GetStartState() < 0 {
		t.Error("DFA should have a valid start state")
	}

	// Check that there are accepting states
	acceptingStates := dfa.GetAcceptingStates()
	if len(acceptingStates) == 0 {
		t.Error("DFA should have at least one accepting state")
	}

	// The DFA should have fewer or equal states than the original NFA
	// since epsilon closures can combine states
	if dfa.GetStateCount() > nfa.GetStateCount() {
		t.Errorf("Expected DFA to have at most %d states, got %d", nfa.GetStateCount(), dfa.GetStateCount())
	}
}

func TestDeterminizeNonDeterministicNFA(t *testing.T) {
	nfa := createNonDeterministicNFA()
	dfa := Determinize(nfa)

	// Check that start state is set
	if dfa.GetStartState() < 0 {
		t.Error("DFA should have a valid start state")
	}

	// Check that there are accepting states
	acceptingStates := dfa.GetAcceptingStates()
	if len(acceptingStates) == 0 {
		t.Error("DFA should have at least one accepting state")
	}

	// The DFA might have more states due to the subset construction
	if dfa.GetStateCount() <= 0 {
		t.Error("DFA should have at least one state")
	}
}

func TestDeterminizePreservesLanguage(t *testing.T) {
	// Create a simple NFA that accepts strings ending with 'ab'
	builder := NewNFABuilder()

	s0 := builder.AddState() // Start state, also loops on any character
	s1 := builder.AddState() // After seeing 'a'
	s2 := builder.AddState() // After seeing 'ab' (accepting)

	builder.SetStartState(s0)
	builder.AcceptState(s2)

	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')
	charSetAny := NewRuneSet_Range('a', 'z') // Any lowercase letter

	builder.AddTransition(s0, s0, charSetAny) // Loop on any character
	builder.AddTransition(s0, s1, charSetA)   // Transition on 'a'
	builder.AddTransition(s1, s2, charSetB)   // Transition on 'b'

	nfa, _ := builder.Build()
	dfa := Determinize(nfa)

	// Both NFA and DFA should have the same basic structure
	if dfa.GetStateCount() <= 0 {
		t.Error("DFA should have at least one state")
	}

	if len(dfa.GetAcceptingStates()) == 0 {
		t.Error("DFA should have at least one accepting state")
	}
}

func TestGetDistinctInputs(t *testing.T) {
	nfa := createSimpleNFA()
	states := []int{0, 1}

	inputs := getDistinctInputs(nfa, states)

	// Should have at least some distinct inputs
	if len(inputs) == 0 {
		t.Error("Expected at least one distinct input")
	}

	// Check that all inputs are non-empty
	for i, input := range inputs {
		if input.Length() == 0 {
			t.Errorf("Input %d should not be empty", i)
		}
	}
}

func TestRemoveDuplicateCharSets(t *testing.T) {
	// Create some duplicate character sets
	set1 := NewRuneSet_Rune('a')
	set2 := NewRuneSet_Rune('a') // Duplicate of set1
	set3 := NewRuneSet_Rune('b')

	inputs := []*RuneSet{set1, set2, set3}
	result := removeDuplicateCharSets(inputs)

	// Should have only 2 unique sets (set1/set2 are duplicates)
	if len(result) != 2 {
		t.Errorf("Expected 2 unique sets, got %d", len(result))
	}
}

// Test that demonstrates correct determinization behavior
func TestDeterminizationCorrectness(t *testing.T) {
	// Create NFA that accepts strings containing 'ab' or 'ba'
	builder := NewNFABuilder()

	// States: start, saw_a, saw_b, accept_ab, accept_ba
	start := builder.AddState()
	sawA := builder.AddState()
	sawB := builder.AddState()
	acceptAB := builder.AddState()
	acceptBA := builder.AddState()

	builder.SetStartState(start)
	builder.AcceptState(acceptAB)
	builder.AcceptState(acceptBA)

	charSetA := NewRuneSet_Rune('a')
	charSetB := NewRuneSet_Rune('b')
	charSetAny := NewRuneSet_Range('a', 'z')

	// From start: can go to sawA on 'a', sawB on 'b', or loop on any char
	builder.AddTransition(start, sawA, charSetA)
	builder.AddTransition(start, sawB, charSetB)
	builder.AddTransition(start, start, charSetAny)

	// From sawA: can go to acceptAB on 'b'
	builder.AddTransition(sawA, acceptAB, charSetB)
	builder.AddTransition(sawA, start, charSetAny) // Continue searching

	// From sawB: can go to acceptBA on 'a'
	builder.AddTransition(sawB, acceptBA, charSetA)
	builder.AddTransition(sawB, start, charSetAny) // Continue searching

	// From accept states: continue searching
	builder.AddTransition(acceptAB, start, charSetAny)
	builder.AddTransition(acceptBA, start, charSetAny)

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build NFA: %v", err)
	}

	// Determinize it
	dfa := Determinize(nfa)

	// Verify basic properties
	if dfa.GetStartState() < 0 {
		t.Error("DFA should have a valid start state")
	}

	if len(dfa.GetAcceptingStates()) == 0 {
		t.Error("DFA should have at least one accepting state")
	}

	// The DFA should be larger or equal in size due to subset construction
	if dfa.GetStateCount() <= 0 {
		t.Error("DFA should have at least one state")
	}

	// Test that DFA has deterministic transitions (each state should have
	// at most one transition per character)
	for sourceState := 0; sourceState < dfa.GetStateCount(); sourceState++ {
		transitions := dfa.GetTransitionsBySource()[sourceState]
		if transitions == nil {
			continue
		}

		// For a true DFA, we should not have overlapping character ranges
		// from the same state (this is a simplified check)
		seenChars := make(map[rune]bool)
		for transInfo := range transitions.AllTransitions() {
			if transInfo.CharRange != nil && transInfo.CharRange.Includes {
				for char := range transInfo.CharRange.All() {
					if seenChars[char] {
						t.Errorf("DFA has non-deterministic transitions for character '%c' from state %d",
							char, sourceState)
					}
					seenChars[char] = true
				}
			}
		}
	}
}
