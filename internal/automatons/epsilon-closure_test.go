package automatons

import (
	"reflect"
	"testing"
)

// Test helper to create a simple NFA for testing epsilon closures
func createTestNFA() *NFA {
	builder := NewNFABuilder()

	// Create states
	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()
	s3 := builder.AddState()
	s4 := builder.AddState()

	// s0 -ε-> s1
	builder.AddTransitionForRuneSet(s0, s1, nil)
	// s1 -ε-> s2
	builder.AddTransitionForRuneSet(s1, s2, nil)
	// s0 -ε-> s3
	builder.AddTransitionForRuneSet(s0, s3, nil)
	// s3 -ε-> s4
	builder.AddTransitionForRuneSet(s3, s4, nil)
	// s2 -ε-> s4
	builder.AddTransitionForRuneSet(s2, s4, nil)
	// Add a regular (non-epsilon) transition for completeness
	builder.AddTransitionForRuneSet(s1, s3, NewRuneSet_Rune('a'))

	builder.SetStartState(s0)
	builder.AcceptState(s4)

	return builder.Build()
}

func TestGetEpsilonClosure_SingleState(t *testing.T) {
	nfa := createTestNFA()

	// Test closure of state 0 (should include 0, 1, 2, 3, 4)
	closure := nfa.GetEpsilonClosure(0)
	expected := NewBitMask_Bits(5, []bool{true, true, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of state 0: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_MultipleStates(t *testing.T) {
	nfa := createTestNFA()

	// Test closure of states 1 and 3
	closure := nfa.GetEpsilonClosure(1, 3)
	expected := NewBitMask_Bits(5, []bool{false, true, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of states 1,3: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_LeafState(t *testing.T) {
	nfa := createTestNFA()

	// Test closure of state 4 (leaf state with no epsilon transitions)
	closure := nfa.GetEpsilonClosure(4)
	expected := NewBitMask_Bits(5, []bool{false, false, false, false, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of state 4: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_EmptyInput(t *testing.T) {
	nfa := createTestNFA()

	// Test closure with no input states
	closure := nfa.GetEpsilonClosure()
	expected := NewBitMask_Empty(5)

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of no states: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_DuplicateStates(t *testing.T) {
	nfa := createTestNFA()

	// Test closure with duplicate input states
	closure := nfa.GetEpsilonClosure(2, 2, 2)
	expected := NewBitMask_Bits(5, []bool{false, false, true, false, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of duplicate states: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_LinearChain(t *testing.T) {
	// Create a simple linear chain: 0 -ε-> 1 -ε-> 2 -ε-> 3
	builder := NewNFABuilder()

	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()
	s3 := builder.AddState()

	builder.AddTransitionForRuneSet(s0, s1, nil)
	builder.AddTransitionForRuneSet(s1, s2, nil)
	builder.AddTransitionForRuneSet(s2, s3, nil)

	builder.SetStartState(s0)
	builder.AcceptState(s3)

	nfa := builder.Build()

	// Test closure of state 0
	closure := nfa.GetEpsilonClosure(s0)
	expected := NewBitMask_Bits(4, []bool{true, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of linear chain: expected %v, got %v", expected, closure)
	}

	// Test closure of state 1
	closure = nfa.GetEpsilonClosure(s1)
	expected = NewBitMask_Bits(4, []bool{false, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of state 1 in linear chain: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_CyclicGraph(t *testing.T) {
	// Create a cyclic epsilon graph: 0 -ε-> 1 -ε-> 2 -ε-> 0
	builder := NewNFABuilder()

	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()

	builder.AddTransitionForRuneSet(s0, s1, nil)
	builder.AddTransitionForRuneSet(s1, s2, nil)
	builder.AddTransitionForRuneSet(s2, s0, nil)

	builder.SetStartState(s0)
	builder.AcceptState(s2)

	nfa := builder.Build()

	// Test closure should include all states in the cycle
	closure := nfa.GetEpsilonClosure(s0)
	expected := NewBitMask_Bits(3, []bool{true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of cyclic graph: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_IsolatedStates(t *testing.T) {
	// Create states with no epsilon transitions
	builder := NewNFABuilder()

	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()

	// Add only regular transitions, no epsilon
	builder.AddTransitionForRuneSet(s0, s1, NewRuneSet_Rune('a'))
	builder.AddTransitionForRuneSet(s1, s2, NewRuneSet_Rune('b'))

	builder.SetStartState(s0)
	builder.AcceptState(s2)

	nfa := builder.Build()

	// Each state should only contain itself in epsilon closure
	for _, state := range []int{s0, s1, s2} {
		closure := nfa.GetEpsilonClosure(state)
		expected := NewBitMask_Empty(3)
		expected.Set(state)
		if !reflect.DeepEqual(closure, expected) {
			t.Errorf("Epsilon closure of isolated state %d: expected %v, got %v", state, expected, closure)
		}
	}
}

// Test helper functions
func TestMaxMin(t *testing.T) {
	testCases := []struct {
		a, b     rune
		expected rune
	}{
		{'a', 'b', 'b'},
		{'z', 'a', 'z'},
		{'中', '文', '文'},
		{0, 1000, 1000},
	}

	for _, tc := range testCases {
		result := maxRune(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("max(%c, %c): expected %c, got %c", tc.a, tc.b, tc.expected, result)
		}
	}

	minTestCases := []struct {
		a, b     rune
		expected rune
	}{
		{'a', 'b', 'a'},
		{'z', 'a', 'a'},
		{'中', '文', '中'},
		{0, 1000, 0},
	}

	for _, tc := range minTestCases {
		result := minRune(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("min(%c, %c): expected %c, got %c", tc.a, tc.b, tc.expected, result)
		}
	}
}

func TestGetEpsilonClosure_ConsistentResults(t *testing.T) {
	nfa := createTestNFA()

	// Run the same closure multiple times and ensure consistent results
	closure1 := nfa.GetEpsilonClosure(0)
	closure2 := nfa.GetEpsilonClosure(0)
	closure3 := nfa.GetEpsilonClosure(0)

	if !reflect.DeepEqual(closure1, closure2) || !reflect.DeepEqual(closure2, closure3) {
		t.Error("Epsilon closure results are not consistent across multiple runs")
	}

	// Test with different input orderings
	closure4 := nfa.GetEpsilonClosure(0)
	closure5 := nfa.GetEpsilonClosure(0, 0, 0) // duplicates

	if !reflect.DeepEqual(closure1, closure4) || !reflect.DeepEqual(closure1, closure5) {
		t.Error("Epsilon closure results vary with different input orderings")
	}
}
