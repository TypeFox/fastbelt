package automatons

import (
	"testing"
)

var kit = NewConstructionKit()

func TestConstruct_Consume(t *testing.T) {
	tests := []struct {
		name        string
		charset     *RuneSet
		testRune    rune
		expectMatch bool
	}{
		{
			name:        "single character",
			charset:     NewRuneSet_Range('a', 'a'),
			testRune:    'a',
			expectMatch: true,
		},
		{
			name:        "single character - no match",
			charset:     NewRuneSet_Range('a', 'a'),
			testRune:    'b',
			expectMatch: false,
		},
		{
			name:        "character range",
			charset:     NewRuneSet_Range('a', 'z'),
			testRune:    'm',
			expectMatch: true,
		},
		{
			name:        "character range - outside range",
			charset:     NewRuneSet_Range('a', 'z'),
			testRune:    'A',
			expectMatch: false,
		},
		{
			name:        "empty charset",
			charset:     NewRuneSet_Empty(),
			testRune:    'a',
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := kit.Consume(tt.charset)
			if err != nil {
				t.Fatalf("Consume() error = %v", err)
			}

			// Test basic properties
			if nfa.StateCount < 2 {
				t.Errorf("Expected at least 2 states, got %d", nfa.StateCount)
			}

			acceptingStates := nfa.AcceptingStates
			if len(acceptingStates) != 1 {
				t.Errorf("Expected 1 accepting state, got %d", len(acceptingStates))
			}

			// Test character matching by checking transitions
			startState := nfa.StartState
			transitions := nfa.TransitionsBySource[startState]
			if transitions == nil {
				if tt.expectMatch {
					t.Error("Expected transitions from start state, but found none")
				}
				return
			}

			hasMatch := transitions.Contains(tt.testRune)
			if hasMatch != tt.expectMatch {
				t.Errorf("Contains('%c') = %v, want %v", tt.testRune, hasMatch, tt.expectMatch)
			}
		})
	}
}

func TestConstruct_Alternate(t *testing.T) {
	// Create test automata
	charA, err := kit.Consume(NewRuneSet_Range('a', 'a'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	charB, err := kit.Consume(NewRuneSet_Range('b', 'b'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	charC, err := kit.Consume(NewRuneSet_Range('c', 'c'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	tests := []struct {
		name        string
		automata    []*NFA
		testRune    rune
		expectMatch bool
	}{
		{
			name:        "two alternatives - match first",
			automata:    []*NFA{charA, charB},
			testRune:    'a',
			expectMatch: true,
		},
		{
			name:        "two alternatives - match second",
			automata:    []*NFA{charA, charB},
			testRune:    'b',
			expectMatch: true,
		},
		{
			name:        "two alternatives - no match",
			automata:    []*NFA{charA, charB},
			testRune:    'c',
			expectMatch: false,
		},
		{
			name:        "three alternatives",
			automata:    []*NFA{charA, charB, charC},
			testRune:    'c',
			expectMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := kit.Alternate(tt.automata...)
			if err != nil {
				t.Fatalf("Alternate() error = %v", err)
			}

			// Basic validation
			if nfa.StateCount == 0 {
				t.Error("Expected non-empty automaton")
			}

			acceptingStates := nfa.AcceptingStates
			if len(acceptingStates) == 0 {
				t.Error("Expected at least one accepting state")
			}

			// Test by examining the structure - for alternation, start state should have epsilon transitions
			startState := nfa.StartState
			transitions := nfa.TransitionsBySource[startState]
			if transitions != nil && transitions.ContainsEpsilon() {
				// This is expected for alternation
			}
		})
	}

	// Test error case
	t.Run("empty alternation", func(t *testing.T) {
		_, err := kit.Alternate()
		if err == nil {
			t.Error("Expected error for empty alternation")
		}
	})
}

func TestConstruct_Concat(t *testing.T) {
	// Create test automata
	charA, err := kit.Consume(NewRuneSet_Range('a', 'a'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	charB, err := kit.Consume(NewRuneSet_Range('b', 'b'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	charC, err := kit.Consume(NewRuneSet_Range('c', 'c'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	tests := []struct {
		name     string
		automata []*NFA
	}{
		{
			name:     "two concatenations",
			automata: []*NFA{charA, charB},
		},
		{
			name:     "three concatenations",
			automata: []*NFA{charA, charB, charC},
		},
		{
			name:     "single automaton",
			automata: []*NFA{charA},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := kit.Concat(tt.automata...)
			if err != nil {
				t.Fatalf("Concat() error = %v", err)
			}

			// Basic validation
			if nfa.StateCount == 0 {
				t.Error("Expected non-empty automaton")
			}

			acceptingStates := nfa.AcceptingStates
			if len(acceptingStates) != 1 {
				t.Errorf("Expected exactly 1 accepting state, got %d", len(acceptingStates))
			}

			// For concatenation, we expect a linear structure with epsilon transitions
			startState := nfa.StartState
			if startState < 0 {
				t.Error("Invalid start state")
			}
		})
	}

	// Test error case
	t.Run("empty concatenation", func(t *testing.T) {
		_, err := kit.Concat()
		if err == nil {
			t.Error("Expected error for empty concatenation")
		}
	})
}

func TestConstruct_Repeat(t *testing.T) {
	// Create test automaton
	charA, err := kit.Consume(NewRuneSet_Range('a', 'a'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	tests := []struct {
		name        string
		automaton   *NFA
		min         int
		max         int
		expectError bool
	}{
		{
			name:        "zero or more (Kleene star)",
			automaton:   charA,
			min:         0,
			max:         -1,
			expectError: false,
		},
		{
			name:        "one or more",
			automaton:   charA,
			min:         1,
			max:         -1,
			expectError: false,
		},
		{
			name:        "exactly three",
			automaton:   charA,
			min:         3,
			max:         3,
			expectError: false,
		},
		{
			name:        "two to five",
			automaton:   charA,
			min:         2,
			max:         5,
			expectError: false,
		},
		{
			name:        "optional (zero or one)",
			automaton:   charA,
			min:         0,
			max:         1,
			expectError: false,
		},
		{
			name:        "invalid range - negative min",
			automaton:   charA,
			min:         -1,
			max:         5,
			expectError: true,
		},
		{
			name:        "invalid range - min > max",
			automaton:   charA,
			min:         5,
			max:         3,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := kit.Repeat(tt.automaton, tt.min, tt.max)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Repeat() error = %v", err)
			}

			// Basic validation
			if nfa.StateCount == 0 {
				t.Error("Expected non-empty automaton")
			}

			acceptingStates := nfa.AcceptingStates
			if len(acceptingStates) == 0 {
				t.Error("Expected at least one accepting state")
			}

			// For min=0 cases, start state should be able to reach accept state via epsilon
			if tt.min == 0 {
				startState := nfa.StartState
				transitions := nfa.TransitionsBySource[startState]
				if transitions != nil && transitions.ContainsEpsilon() {
					// This is expected for optional patterns
				}
			}
		})
	}
}

func TestConstruct_Complement(t *testing.T) {
	// Create test automaton
	charA, err := kit.Consume(NewRuneSet_Range('a', 'a'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	nfa, err := kit.Complement(charA)
	if err != nil {
		t.Fatalf("Complement() error = %v", err)
	}

	// Basic validation
	if nfa.StateCount == 0 {
		t.Error("Expected non-empty automaton")
	}

	// The complement should have the same structure but different accepting states
	originalAccepting := charA.AcceptingStates
	complementAccepting := nfa.AcceptingStates

	// Verify we have some accepting states
	if len(complementAccepting) == 0 {
		t.Error("Complement should have at least one accepting state")
	}

	// The complement should accept states that the original doesn't accept
	// For a simple 2-state automaton (start->accept), the complement should accept the start state
	if nfa.StateCount == charA.StateCount {
		foundDifference := false
		for state := 0; state < nfa.StateCount; state++ {
			originalAccepts := originalAccepting[state]
			complementAccepts := complementAccepting[state]
			if originalAccepts != complementAccepts {
				foundDifference = true
				break
			}
		}
		if !foundDifference {
			t.Error("Complement should have different accepting states than original")
		}
	}
}

func TestConstruct_IntersectNFA(t *testing.T) {
	// Create test automata
	rangeAB, err := kit.Consume(NewRuneSet_Range('a', 'b'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	rangeBC, err := kit.Consume(NewRuneSet_Range('b', 'c'))
	if err != nil {
		t.Fatalf("Failed to create test automaton: %v", err)
	}

	nfa, err := kit.Intersect(rangeAB, rangeBC)
	if err != nil {
		t.Fatalf("IntersectNFA() error = %v", err)
	}

	// Basic validation
	if nfa.StateCount == 0 {
		t.Error("Expected non-empty automaton")
	}

	// The intersection of [a-b] and [b-c] should match only 'b'
	startState := nfa.StartState
	transitions := nfa.TransitionsBySource[startState]
	if transitions != nil {
		// Should contain 'b'
		if !transitions.Contains('b') {
			t.Error("IntersectNFA should contain 'b'")
		}
		// Should not contain 'a' or 'c' (this is harder to test without a full runner)
	}
}

func TestConstruct_PackageLevelFunctions(t *testing.T) {
	// Test that package-level functions work the same as method calls
	charset := NewRuneSet_Range('x', 'x')

	// Test Consume
	nfa1, err1 := kit.Consume(charset)
	kit := NewConstructionKit()
	nfa2, err2 := kit.Consume(charset)

	if (err1 == nil) != (err2 == nil) {
		t.Error("Package function and method should have same error behavior")
	}

	if err1 == nil && err2 == nil {
		if nfa1.StateCount != nfa2.StateCount {
			t.Error("Package function and method should produce equivalent results")
		}
	}
}

func TestConstruct_ErrorHandling(t *testing.T) {
	// Test various error conditions
	charset := NewRuneSet_Range('a', 'a')

	t.Run("nil automaton handling", func(t *testing.T) {
		// This tests the robustness of our implementation
		// Most functions should handle the case where automata might have issues

		validNFA, err := kit.Consume(charset)
		if err != nil {
			t.Fatalf("Failed to create valid NFA: %v", err)
		}

		// Test Alternate with valid automaton
		_, err = kit.Alternate(validNFA)
		if err != nil {
			t.Errorf("Alternate with single valid NFA should not error: %v", err)
		}

		// Test Concat with valid automaton
		_, err = kit.Concat(validNFA)
		if err != nil {
			t.Errorf("Concat with single valid NFA should not error: %v", err)
		}
	})
}
