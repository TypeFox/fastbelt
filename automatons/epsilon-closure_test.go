package automatons

import (
	"reflect"
	"sort"
	"testing"
)

// Test helper to create a simple NFA for testing epsilon closures
func createTestNFA() NFA {
	builder := NewNFABuilder()

	// Create states
	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()
	s3 := builder.AddState()
	s4 := builder.AddState()

	// Set up epsilon transitions (using empty RuneSet)
	emptyCharset := NewRuneSet_Empty()

	// s0 -ε-> s1
	builder.AddTransition(s0, s1, emptyCharset)
	// s1 -ε-> s2
	builder.AddTransition(s1, s2, emptyCharset)
	// s0 -ε-> s3
	builder.AddTransition(s0, s3, emptyCharset)
	// s3 -ε-> s4
	builder.AddTransition(s3, s4, emptyCharset)
	// s2 -ε-> s4
	builder.AddTransition(s2, s4, emptyCharset)

	// Add a regular (non-epsilon) transition for completeness
	builder.AddTransition(s1, s3, NewRuneSet_Rune('a'))

	builder.SetStartState(s0)
	builder.AcceptState(s4)

	nfa, err := builder.Build()
	if err != nil {
		panic("Failed to build test NFA: " + err.Error())
	}

	return nfa
}

func TestGetEpsilonClosure_SingleState(t *testing.T) {
	nfa := createTestNFA()

	// Test closure of state 0 (should include 0, 1, 2, 3, 4)
	closure := GetEpsilonClosure(nfa, 0)
	expected := NewBitMask_Bits(5, []bool{true, true, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of state 0: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_MultipleStates(t *testing.T) {
	nfa := createTestNFA()

	// Test closure of states 1 and 3
	closure := GetEpsilonClosure(nfa, 1, 3)
	expected := NewBitMask_Bits(5, []bool{false, true, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of states 1,3: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_LeafState(t *testing.T) {
	nfa := createTestNFA()

	// Test closure of state 4 (leaf state with no epsilon transitions)
	closure := GetEpsilonClosure(nfa, 4)
	expected := NewBitMask_Bits(5, []bool{false, false, false, false, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of state 4: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_EmptyInput(t *testing.T) {
	nfa := createTestNFA()

	// Test closure with no input states
	closure := GetEpsilonClosure(nfa)
	expected := NewBitMask_Empty(5)

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of no states: expected %v, got %v", expected, closure)
	}
}

func TestGetEpsilonClosure_DuplicateStates(t *testing.T) {
	nfa := createTestNFA()

	// Test closure with duplicate input states
	closure := GetEpsilonClosure(nfa, 2, 2, 2)
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

	emptyCharset := NewRuneSet_Empty()
	builder.AddTransition(s0, s1, emptyCharset)
	builder.AddTransition(s1, s2, emptyCharset)
	builder.AddTransition(s2, s3, emptyCharset)

	builder.SetStartState(s0)
	builder.AcceptState(s3)

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build linear chain NFA: %v", err)
	}

	// Test closure of state 0
	closure := GetEpsilonClosure(nfa, s0)
	expected := NewBitMask_Bits(4, []bool{true, true, true, true})

	if !reflect.DeepEqual(closure, expected) {
		t.Errorf("Epsilon closure of linear chain: expected %v, got %v", expected, closure)
	}

	// Test closure of state 1
	closure = GetEpsilonClosure(nfa, s1)
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

	emptyCharset := NewRuneSet_Empty()
	builder.AddTransition(s0, s1, emptyCharset)
	builder.AddTransition(s1, s2, emptyCharset)
	builder.AddTransition(s2, s0, emptyCharset)

	builder.SetStartState(s0)
	builder.AcceptState(s2)

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build cyclic NFA: %v", err)
	}

	// Test closure should include all states in the cycle
	closure := GetEpsilonClosure(nfa, s0)
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
	builder.AddTransition(s0, s1, NewRuneSet_Rune('a'))
	builder.AddTransition(s1, s2, NewRuneSet_Rune('b'))

	builder.SetStartState(s0)
	builder.AcceptState(s2)

	nfa, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build isolated states NFA: %v", err)
	}

	// Each state should only contain itself in epsilon closure
	for _, state := range []int{s0, s1, s2} {
		closure := GetEpsilonClosure(nfa, state)
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
		result := max(tc.a, tc.b)
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
		result := min(tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("min(%c, %c): expected %c, got %c", tc.a, tc.b, tc.expected, result)
		}
	}
}

// Benchmark the epsilon closure algorithm
func BenchmarkGetEpsilonClosure_LinearChain(b *testing.B) {
	// Create a long linear chain for benchmarking
	builder := NewNFABuilder()

	const chainLength = 100
	states := make([]int, chainLength)
	emptyCharset := NewRuneSet_Empty()

	for i := 0; i < chainLength; i++ {
		states[i] = builder.AddState()
		if i > 0 {
			builder.AddTransition(states[i-1], states[i], emptyCharset)
		}
	}

	builder.SetStartState(states[0])
	builder.AcceptState(states[chainLength-1])

	nfa, err := builder.Build()
	if err != nil {
		b.Fatalf("Failed to build benchmark NFA: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetEpsilonClosure(nfa, states[0])
	}
}

// Helper function to convert map[int]bool to sorted slice for easier testing
func closureToSortedSlice(closure map[int]bool) []int {
	result := make([]int, 0, len(closure))
	for state := range closure {
		result = append(result, state)
	}
	sort.Ints(result)
	return result
}

func TestGetEpsilonClosure_ConsistentResults(t *testing.T) {
	nfa := createTestNFA()

	// Run the same closure multiple times and ensure consistent results
	closure1 := GetEpsilonClosure(nfa, 0)
	closure2 := GetEpsilonClosure(nfa, 0)
	closure3 := GetEpsilonClosure(nfa, 0)

	if !reflect.DeepEqual(closure1, closure2) || !reflect.DeepEqual(closure2, closure3) {
		t.Error("Epsilon closure results are not consistent across multiple runs")
	}

	// Test with different input orderings
	closure4 := GetEpsilonClosure(nfa, 0)
	closure5 := GetEpsilonClosure(nfa, 0, 0, 0) // duplicates

	if !reflect.DeepEqual(closure1, closure4) || !reflect.DeepEqual(closure1, closure5) {
		t.Error("Epsilon closure results vary with different input orderings")
	}
}
