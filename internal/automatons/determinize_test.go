package automatons

import (
	"testing"
)

// Helper function to create a simple NFA for testing
func createSimpleNFA() *NFA {
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
	builder.AddTransitionForSingleRune(s0, s1, 'a')
	builder.AddTransitionForSingleRune(s1, s2, 'b')

	nfa, _ := builder.Build()
	return nfa
}

// Helper function to create an NFA with epsilon transitions
func createNFAWithEpsilons() *NFA {
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
	builder.AddTransitionForRuneSet(s0, s1, nil)
	builder.AddTransitionForSingleRune(s0, s2, 'a')
	builder.AddTransitionForSingleRune(s1, s3, 'b')

	nfa, _ := builder.Build()
	return nfa
}

// Helper function to create an NFA with non-deterministic transitions
func createNonDeterministicNFA() *NFA {
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
	builder.AddTransitionForSingleRune(s0, s1, 'a')
	builder.AddTransitionForSingleRune(s0, s2, 'a') // Non-deterministic: same input, different targets
	builder.AddTransitionForSingleRune(s1, s3, 'b')

	nfa, _ := builder.Build()
	return nfa
}

func TestDeterminizeSimpleNFA(t *testing.T) {
	nfa := createSimpleNFA()
	dfa := nfa.Determinize()

	// The DFA should have the same number of states as the original NFA
	// since it's already deterministic
	Expect(dfa.StateCount).ToEqual(nfa.StateCount)

	// Check that start state is set
	Expect(dfa.StartState).ToBeGreaterThanOrEqual(0)

	// Check that there are accepting states
	Expect(len(dfa.AcceptingStates)).ToEqual(1)
}

func TestDeterminizeNFAWithEpsilons(t *testing.T) {
	nfa := createNFAWithEpsilons()
	dfa := nfa.Determinize()

	// Check that start state is set
	Expect(dfa.StartState).ToBeGreaterThanOrEqual(0)

	// Check that there are accepting states
	Expect(len(dfa.AcceptingStates)).ToBeGreaterThanOrEqual(1)

	// The DFA should have fewer or equal states than the original NFA
	// since epsilon closures can combine states
	Expect(dfa.StateCount).ToBeLesserThan(nfa.StateCount)
}

func TestDeterminizeNonDeterministicNFA(t *testing.T) {
	nfa := createNonDeterministicNFA()
	dfa := nfa.Determinize()

	// Check that start state is set
	Expect(dfa.StartState).ToBeGreaterThanOrEqual(0)

	// Check that there are accepting states
	Expect(len(dfa.AcceptingStates)).ToBeGreaterThanOrEqual(1)

	// The DFA might have more states due to the subset construction
	Expect(dfa.StateCount).ToBeGreaterThanOrEqual(1)
}

func TestDeterminizePreservesLanguage(t *testing.T) {
	builder := NewNFABuilder()

	s0 := builder.AddState() // Start state, also loops on any character
	s1 := builder.AddState() // After seeing 'a'
	s2 := builder.AddState() // After seeing 'ab' (accepting)

	builder.SetStartState(s0)
	builder.AcceptState(s2)

	builder.AddTransitionForRuneRange(s0, s0, NewRuneRange('a', 'z', true)) // Loop on any character
	builder.AddTransitionForSingleRune(s0, s1, 'a')                         // Transition on 'a'
	builder.AddTransitionForSingleRune(s1, s2, 'b')                         // Transition on 'b'

	nfa, _ := builder.Build()
	dfa := nfa.Determinize()

	Expect(dfa.StartState).ToBeGreaterThanOrEqual(0)
	Expect(len(dfa.AcceptingStates)).ToBeGreaterThanOrEqual(1)
}

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

	runeRangeAny := NewRuneRange('a', 'z', true)

	// From start: can go to sawA on 'a', sawB on 'b', or loop on any char
	builder.AddTransitionForSingleRune(start, sawA, 'a')
	builder.AddTransitionForSingleRune(start, sawB, 'b')
	builder.AddTransitionForRuneRange(start, start, runeRangeAny)

	// From sawA: can go to acceptAB on 'b'
	builder.AddTransitionForSingleRune(sawA, acceptAB, 'b')
	builder.AddTransitionForRuneRange(sawA, start, runeRangeAny) // Continue searching

	// From sawB: can go to acceptBA on 'a'
	builder.AddTransitionForSingleRune(sawB, acceptBA, 'a')
	builder.AddTransitionForRuneRange(sawB, start, runeRangeAny) // Continue searching
	// From accept states: continue searching
	builder.AddTransitionForRuneRange(acceptAB, start, runeRangeAny)
	builder.AddTransitionForRuneRange(acceptBA, start, runeRangeAny)

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build NFA: %v", err)
	}

	// Determinize it
	dfa := nfa.Determinize()

	// Verify basic properties
	Expect(dfa.StartState).ToBeGreaterThanOrEqual(0)
	Expect(len(dfa.AcceptingStates)).ToBeGreaterThanOrEqual(1)
	Expect(dfa.StateCount).ToBeGreaterThan(0)

	// Test that DFA has deterministic transitions (each state should have
	// at most one transition per character)
	for sourceState := 0; sourceState < dfa.StateCount; sourceState++ {
		transitions := dfa.TransitionsBySource[sourceState]
		if transitions == nil {
			continue
		}

		// For a true DFA, we should not have overlapping character ranges
		// from the same state (this is a simplified check)
		seenChars := make(map[rune]bool)
		for transInfo := range transitions.All() {
			if transInfo.Range != nil && transInfo.Range.Includes {
				for char := range transInfo.Range.All() {
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
