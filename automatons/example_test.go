package automatons

import (
	"fmt"
	"sort"
)

// Example demonstrating the epsilon closure algorithm
func ExampleGetEpsilonClosure() {
	// Create a simple NFA with epsilon transitions
	builder := NewNFABuilder()

	// Create states: 0 -> 1 -> 2
	//                 \-> 3 -> 4
	s0 := builder.AddState()
	s1 := builder.AddState()
	s2 := builder.AddState()
	s3 := builder.AddState()
	s4 := builder.AddState()

	emptyCharset := NewRuneSet_Empty()

	// Add epsilon transitions
	builder.AddTransition(s0, s1, emptyCharset) // 0 -ε-> 1
	builder.AddTransition(s1, s2, emptyCharset) // 1 -ε-> 2
	builder.AddTransition(s0, s3, emptyCharset) // 0 -ε-> 3
	builder.AddTransition(s3, s4, emptyCharset) // 3 -ε-> 4

	builder.SetStartState(s0)
	builder.AcceptState(s4)

	nfa, _ := builder.Build()

	// Get epsilon closure of state 0
	closure := GetEpsilonClosure(nfa, 0)

	// Convert to sorted slice for consistent output
	states := make([]int, 0, len(closure))
	for state := range closure {
		states = append(states, state)
	}
	sort.Ints(states)

	fmt.Printf("Epsilon closure of state 0: %v\n", states)
	// Output: Epsilon closure of state 0: [0 1 2 3 4]
}
