package automatons

import (
	"testing"
)

// TestConstruct_IntegrationExample demonstrates a complex example using multiple construction functions
func TestConstruct_IntegrationExample(t *testing.T) {
	// Create an NFA that matches: (a|b)*c+
	// This means: zero or more 'a' or 'b', followed by one or more 'c'

	// First create automata for individual characters
	charA := kit.Consume(NewRuneSet_Range('a', 'a'))
	charB := kit.Consume(NewRuneSet_Range('b', 'b'))
	charC := kit.Consume(NewRuneSet_Range('c', 'c'))

	// Create (a|b)
	aOrB := kit.Alternate(charA, charB)

	// Create (a|b)*
	aOrBStar := kit.Repeat(aOrB, 0, -1)

	// Create c+
	cPlus := kit.Repeat(charC, 1, -1)

	// Create (a|b)*c+
	final := kit.Concat(aOrBStar, cPlus)

	// Basic validation
	Expect(final.StateCount).ToBeGreaterThan(0)

	// Should have exactly one accepting state
	Expect(len(final.AcceptingStates)).ToEqual(1)

	// Should have a valid start state
	Expect(final.StartState).ToBeGreaterThanOrEqual(0)
}

// TestConstruct_RegexLikePatterns tests common regex-like patterns
func TestConstruct_RegexLikePatterns(t *testing.T) {
	tests := []struct {
		name        string
		description string
		builder     func() *NFA
	}{
		{
			name:        "optional_char",
			description: "a?  (zero or one 'a')",
			builder: func() *NFA {
				charA := kit.Consume(NewRuneSet_Range('a', 'a'))
				return kit.Repeat(charA, 0, 1)
			},
		},
		{
			name:        "kleene_star",
			description: "a*  (zero or more 'a')",
			builder: func() *NFA {
				charA := kit.Consume(NewRuneSet_Range('a', 'a'))
				return kit.Repeat(charA, 0, -1)
			},
		},
		{
			name:        "kleene_plus",
			description: "a+  (one or more 'a')",
			builder: func() *NFA {
				charA := kit.Consume(NewRuneSet_Range('a', 'a'))
				return kit.Repeat(charA, 1, -1)
			},
		},
		{
			name:        "exact_count",
			description: "a{3}  (exactly three 'a')",
			builder: func() *NFA {
				charA := kit.Consume(NewRuneSet_Range('a', 'a'))
				return kit.Repeat(charA, 3, 3)
			},
		},
		{
			name:        "range_count",
			description: "a{2,5}  (two to five 'a')",
			builder: func() *NFA {
				charA := kit.Consume(NewRuneSet_Range('a', 'a'))
				return kit.Repeat(charA, 2, 5)
			},
		},
		{
			name:        "character_class",
			description: "[a-z]  (any lowercase letter)",
			builder: func() *NFA {
				return kit.Consume(NewRuneSet_Range('a', 'z'))
			},
		},
		{
			name:        "alternation",
			description: "(hello|world)  (either 'hello' or 'world')",
			builder: func() *NFA {
				// Create "hello"
				h := kit.Consume(NewRuneSet_Range('h', 'h'))
				e := kit.Consume(NewRuneSet_Range('e', 'e'))
				l1 := kit.Consume(NewRuneSet_Range('l', 'l'))
				l2 := kit.Consume(NewRuneSet_Range('l', 'l'))
				o := kit.Consume(NewRuneSet_Range('o', 'o'))
				hello := kit.Concat(h, e, l1, l2, o)

				// Create "world"
				w := kit.Consume(NewRuneSet_Range('w', 'w'))
				o2 := kit.Consume(NewRuneSet_Range('o', 'o'))
				r := kit.Consume(NewRuneSet_Range('r', 'r'))
				l3 := kit.Consume(NewRuneSet_Range('l', 'l'))
				d := kit.Consume(NewRuneSet_Range('d', 'd'))
				world := kit.Concat(w, o2, r, l3, d)

				return kit.Alternate(hello, world)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa := tt.builder()
			Expect(nfa.StateCount).ToBeGreaterThan(0)
			Expect(len(nfa.AcceptingStates)).ToBeGreaterThan(0)
			Expect(nfa.StartState).ToBeGreaterThanOrEqual(0)
		})
	}
}

// TestConstruct_SetOperations tests set-like operations on automata
func TestConstruct_SetOperations(t *testing.T) {
	// Create test sets: [a-m] and [h-z]
	setAM := kit.Consume(NewRuneSet_Range('a', 'm'))
	setHZ := kit.Consume(NewRuneSet_Range('h', 'z'))

	// Test complement
	notAM := kit.Complement(setAM)

	// Test intersection: [a-m] ∩ [h-z] should be [h-m]
	intersection := kit.Intersect(setAM, setHZ)

	// Test union using alternation: [a-m] ∪ [h-z]
	union := kit.Alternate(setAM, setHZ)

	// Basic validation for all results
	automata := map[string]*NFA{
		"setAM":        setAM,
		"setHZ":        setHZ,
		"notAM":        notAM,
		"intersection": intersection,
		"union":        union,
	}

	for _, nfa := range automata {
		Expect(nfa.StateCount).ToBeGreaterThan(0)
		Expect(len(nfa.AcceptingStates)).ToBeGreaterThan(0)
		Expect(nfa.StartState).ToBeGreaterThanOrEqual(0)
	}
}
