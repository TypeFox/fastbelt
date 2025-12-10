package automatons

import (
	"testing"
)

// TestConstruct_IntegrationExample demonstrates a complex example using multiple construction functions
func TestConstruct_IntegrationExample(t *testing.T) {
	// Create an NFA that matches: (a|b)*c+
	// This means: zero or more 'a' or 'b', followed by one or more 'c'

	// First create automata for individual characters
	charA, err := Consume(NewRuneSet_Range('a', 'a'))
	if err != nil {
		t.Fatalf("Failed to create automaton for 'a': %v", err)
	}

	charB, err := Consume(NewRuneSet_Range('b', 'b'))
	if err != nil {
		t.Fatalf("Failed to create automaton for 'b': %v", err)
	}

	charC, err := Consume(NewRuneSet_Range('c', 'c'))
	if err != nil {
		t.Fatalf("Failed to create automaton for 'c': %v", err)
	}

	// Create (a|b)
	aOrB, err := Alternate(charA, charB)
	if err != nil {
		t.Fatalf("Failed to create alternation a|b: %v", err)
	}

	// Create (a|b)*
	aOrBStar, err := Repeat(aOrB, 0, -1)
	if err != nil {
		t.Fatalf("Failed to create (a|b)*: %v", err)
	}

	// Create c+
	cPlus, err := Repeat(charC, 1, -1)
	if err != nil {
		t.Fatalf("Failed to create c+: %v", err)
	}

	// Create (a|b)*c+
	final, err := Concat(aOrBStar, cPlus)
	if err != nil {
		t.Fatalf("Failed to create final automaton: %v", err)
	}

	// Basic validation
	if final.GetStateCount() == 0 {
		t.Error("Expected non-empty final automaton")
	}

	// Should have exactly one accepting state
	acceptingStates := final.GetAcceptingStates()
	if len(acceptingStates) != 1 {
		t.Errorf("Expected exactly 1 accepting state, got %d", len(acceptingStates))
	}

	// Should have a valid start state
	if final.GetStartState() < 0 {
		t.Error("Invalid start state")
	}

	t.Logf("Successfully created complex automaton with %d states", final.GetStateCount())
}

// TestConstruct_RegexLikePatterns tests common regex-like patterns
func TestConstruct_RegexLikePatterns(t *testing.T) {
	tests := []struct {
		name        string
		description string
		builder     func() (NFA, error)
	}{
		{
			name:        "optional_char",
			description: "a?  (zero or one 'a')",
			builder: func() (NFA, error) {
				charA, err := Consume(NewRuneSet_Range('a', 'a'))
				if err != nil {
					return nil, err
				}
				return Repeat(charA, 0, 1)
			},
		},
		{
			name:        "kleene_star",
			description: "a*  (zero or more 'a')",
			builder: func() (NFA, error) {
				charA, err := Consume(NewRuneSet_Range('a', 'a'))
				if err != nil {
					return nil, err
				}
				return Repeat(charA, 0, -1)
			},
		},
		{
			name:        "kleene_plus",
			description: "a+  (one or more 'a')",
			builder: func() (NFA, error) {
				charA, err := Consume(NewRuneSet_Range('a', 'a'))
				if err != nil {
					return nil, err
				}
				return Repeat(charA, 1, -1)
			},
		},
		{
			name:        "exact_count",
			description: "a{3}  (exactly three 'a')",
			builder: func() (NFA, error) {
				charA, err := Consume(NewRuneSet_Range('a', 'a'))
				if err != nil {
					return nil, err
				}
				return Repeat(charA, 3, 3)
			},
		},
		{
			name:        "range_count",
			description: "a{2,5}  (two to five 'a')",
			builder: func() (NFA, error) {
				charA, err := Consume(NewRuneSet_Range('a', 'a'))
				if err != nil {
					return nil, err
				}
				return Repeat(charA, 2, 5)
			},
		},
		{
			name:        "character_class",
			description: "[a-z]  (any lowercase letter)",
			builder: func() (NFA, error) {
				return Consume(NewRuneSet_Range('a', 'z'))
			},
		},
		{
			name:        "alternation",
			description: "(hello|world)  (either 'hello' or 'world')",
			builder: func() (NFA, error) {
				// Create "hello"
				h, _ := Consume(NewRuneSet_Range('h', 'h'))
				e, _ := Consume(NewRuneSet_Range('e', 'e'))
				l1, _ := Consume(NewRuneSet_Range('l', 'l'))
				l2, _ := Consume(NewRuneSet_Range('l', 'l'))
				o, _ := Consume(NewRuneSet_Range('o', 'o'))
				hello, err := Concat(h, e, l1, l2, o)
				if err != nil {
					return nil, err
				}

				// Create "world"
				w, _ := Consume(NewRuneSet_Range('w', 'w'))
				o2, _ := Consume(NewRuneSet_Range('o', 'o'))
				r, _ := Consume(NewRuneSet_Range('r', 'r'))
				l3, _ := Consume(NewRuneSet_Range('l', 'l'))
				d, _ := Consume(NewRuneSet_Range('d', 'd'))
				world, err := Concat(w, o2, r, l3, d)
				if err != nil {
					return nil, err
				}

				return Alternate(hello, world)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa, err := tt.builder()
			if err != nil {
				t.Fatalf("Failed to build %s: %v", tt.description, err)
			}

			// Basic validation
			if nfa.GetStateCount() == 0 {
				t.Errorf("Pattern %s produced empty automaton", tt.description)
			}

			acceptingStates := nfa.GetAcceptingStates()
			if len(acceptingStates) == 0 {
				t.Errorf("Pattern %s has no accepting states", tt.description)
			}

			if nfa.GetStartState() < 0 {
				t.Errorf("Pattern %s has invalid start state", tt.description)
			}

			t.Logf("Pattern %s: %d states, %d accepting states",
				tt.description, nfa.GetStateCount(), len(acceptingStates))
		})
	}
}

// TestConstruct_SetOperations tests set-like operations on automata
func TestConstruct_SetOperations(t *testing.T) {
	// Create test sets: [a-m] and [h-z]
	setAM, err := Consume(NewRuneSet_Range('a', 'm'))
	if err != nil {
		t.Fatalf("Failed to create set [a-m]: %v", err)
	}

	setHZ, err := Consume(NewRuneSet_Range('h', 'z'))
	if err != nil {
		t.Fatalf("Failed to create set [h-z]: %v", err)
	}

	// Test complement
	notAM, err := Complement(setAM)
	if err != nil {
		t.Fatalf("Failed to create complement of [a-m]: %v", err)
	}

	// Test intersection: [a-m] ∩ [h-z] should be [h-m]
	intersection, err := IntersectNFA(setAM, setHZ)
	if err != nil {
		t.Fatalf("Failed to create intersection: %v", err)
	}

	// Test union using alternation: [a-m] ∪ [h-z]
	union, err := Alternate(setAM, setHZ)
	if err != nil {
		t.Fatalf("Failed to create union: %v", err)
	}

	// Basic validation for all results
	automata := map[string]NFA{
		"setAM":        setAM,
		"setHZ":        setHZ,
		"notAM":        notAM,
		"intersection": intersection,
		"union":        union,
	}

	for name, nfa := range automata {
		if nfa.GetStateCount() == 0 {
			t.Errorf("%s produced empty automaton", name)
		}

		acceptingStates := nfa.GetAcceptingStates()
		if len(acceptingStates) == 0 {
			t.Errorf("%s has no accepting states", name)
		}

		if nfa.GetStartState() < 0 {
			t.Errorf("%s has invalid start state", name)
		}

		t.Logf("%s: %d states, %d accepting states",
			name, nfa.GetStateCount(), len(acceptingStates))
	}
}
