package automatons

// Helper function for max
func max(a, b rune) rune {
	if a > b {
		return a
	}
	return b
}

// Helper function for min
func min(a, b rune) rune {
	if a < b {
		return a
	}
	return b
}

// GetEpsilonClosure computes the epsilon closure of the given states in the NFA.
//
// The epsilon closure of a set of states is the set of all states that can be reached
// from the input states by following zero or more epsilon (empty) transitions.
// This is a fundamental operation in NFA algorithms, particularly useful for:
//   - Converting NFAs to DFAs (subset construction)
//   - Simulating NFA execution
//   - Analyzing reachability in finite automata
//
// The algorithm uses a depth-first search approach to explore all epsilon transitions
// from the input states. It maintains a set of visited states to avoid infinite loops
// in case of epsilon cycles.
//
// Parameters:
//   - nfa: The NFA to analyze
//   - states: Zero or more state IDs to start the closure computation from
//
// Returns:
//
//	A map where keys are state IDs that belong to the epsilon closure,
//	and values are always true (the map serves as a set).
//
// Example:
//
//	Given an NFA with states 0,1,2,3 and epsilon transitions:
//	0 -ε-> 1, 1 -ε-> 2, 0 -ε-> 3
//	GetEpsilonClosure(nfa, 0) returns map[0:true, 1:true, 2:true, 3:true]
//
// Time Complexity: O(V + E) where V is the number of states and E is the number of epsilon transitions.
// Space Complexity: O(V) for the closure set and internal queue.
func GetEpsilonClosure(nfa NFA, states ...int) BitMask {
	closure := NewBitMask_Empty(nfa.GetStateCount())
	queue := make([]int, 0, len(states))

	// Initialize with input states
	queue = append(queue, states...)

	// Process queue until empty
	for len(queue) > 0 {
		// Pop from end (stack-like behavior, but order doesn't matter for correctness)
		source := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		// Skip if already processed
		if closure.IsSet(source) {
			continue
		}

		// Mark as visited
		closure.Set(source)

		// Get transitions from this state
		targets := nfa.GetTransitionsBySource()[source]
		if targets == nil {
			continue
		}

		// Add epsilon targets to queue
		for _, target := range targets.GetEpsilonTargets() {
			if !closure.IsSet(target) {
				queue = append(queue, target)
			}
		}
	}

	return closure
}
