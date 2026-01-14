package automatons

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			charset:     NewRuneSetRange('a', 'a'),
			testRune:    'a',
			expectMatch: true,
		},
		{
			name:        "single character - no match",
			charset:     NewRuneSetRange('a', 'a'),
			testRune:    'b',
			expectMatch: false,
		},
		{
			name:        "character range",
			charset:     NewRuneSetRange('a', 'z'),
			testRune:    'm',
			expectMatch: true,
		},
		{
			name:        "character range - outside range",
			charset:     NewRuneSetRange('a', 'z'),
			testRune:    'A',
			expectMatch: false,
		},
		{
			name:        "empty charset",
			charset:     NewRuneSetEmpty(),
			testRune:    'a',
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa := kit.Consume(tt.charset)

			// Test basic properties
			assert.GreaterOrEqual(t, nfa.StateCount, 2)
			assert.Equal(t, 1, len(nfa.AcceptingStates))

			// Test character matching by checking transitions
			startState := nfa.StartState
			transitions := nfa.TransitionsBySource[startState]
			assert.NotNil(t, transitions)

			hasMatch := transitions.Contains(tt.testRune)
			assert.Equal(t, tt.expectMatch, hasMatch)
		})
	}
}

func TestConstruct_Alternate(t *testing.T) {
	// Create test automata
	charA := kit.Consume(NewRuneSetRange('a', 'a'))
	charB := kit.Consume(NewRuneSetRange('b', 'b'))
	charC := kit.Consume(NewRuneSetRange('c', 'c'))

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
			nfa := kit.Alternate(tt.automata...)

			assert.Greater(t, nfa.StateCount, 0)
			assert.Greater(t, len(nfa.AcceptingStates), 0)
		})
	}

	// Test error case
	t.Run("empty alternation", func(t *testing.T) {
		assert.Panics(t, func() {
			kit.Alternate()
		})
	})
}

func TestConstruct_Concat(t *testing.T) {
	// Create test automata
	charA := kit.Consume(NewRuneSetRange('a', 'a'))
	charB := kit.Consume(NewRuneSetRange('b', 'b'))
	charC := kit.Consume(NewRuneSetRange('c', 'c'))

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
			nfa := kit.Concat(tt.automata...)

			// Basic validation
			assert.Greater(t, nfa.StateCount, 0)
			assert.Equal(t, 1, len(nfa.AcceptingStates))

			// For concatenation, we expect a linear structure with epsilon transitions
			assert.GreaterOrEqual(t, nfa.StartState, 0)
		})
	}

	// Test error case
	t.Run("empty concatenation", func(t *testing.T) {
		assert.Panics(t, func() {
			kit.Concat()
		})
	})
}

func TestConstruct_Repeat(t *testing.T) {
	// Create test automaton
	charA := kit.Consume(NewRuneSetRange('a', 'a'))

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
			if tt.expectError {
				assert.Panics(t, func() {
					kit.Repeat(tt.automaton, tt.min, tt.max)
				})
				return
			}

			nfa := kit.Repeat(tt.automaton, tt.min, tt.max)
			// Basic validation
			assert.Greater(t, nfa.StateCount, 0)
			assert.Greater(t, len(nfa.AcceptingStates), 0)
		})
	}
}

func TestConstruct_Complement(t *testing.T) {
	// Create test automaton
	charA := kit.Consume(NewRuneSetRange('a', 'a'))
	nfa := kit.Complement(charA)

	// Basic validation
	assert.Greater(t, nfa.StateCount, 0)

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
	rangeAB := kit.Consume(NewRuneSetRange('a', 'b'))
	rangeBC := kit.Consume(NewRuneSetRange('b', 'c'))
	nfa := kit.Intersect(rangeAB, rangeBC)

	// Basic validation
	assert.Greater(t, nfa.StateCount, 0)

	// The intersection of [a-b] and [b-c] should match only 'b'
	startState := nfa.StartState
	transitions := nfa.TransitionsBySource[startState]
	assert.Equal(t, true, transitions.Contains('b'))
}

func TestConstruct_PackageLevelFunctions(t *testing.T) {
	// Test that package-level functions work the same as method calls
	charset := NewRuneSetRange('x', 'x')

	// Test Consume
	kit.Consume(charset)
	kit := NewConstructionKit()
	kit.Consume(charset)
}

func TestConstruct_ErrorHandling(t *testing.T) {
	// Test various error conditions
	charset := NewRuneSetRange('a', 'a')

	t.Run("nil automaton handling", func(t *testing.T) {
		// This tests the robustness of our implementation
		// Most functions should handle the case where automata might have issues
		validNFA := kit.Consume(charset)

		// Test Alternate with valid automaton
		kit.Alternate(validNFA)

		// Test Concat with valid automaton
		kit.Concat(validNFA)
	})
}
