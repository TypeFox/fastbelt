package automatons

import (
	"fmt"
)

// Minimize minimizes a DFA using the table-filling algorithm
// This removes unreachable states and merges equivalent states
func Minimize(dfa NFA) NFA {
	stateCount := dfa.GetStateCount()
	acceptingStates := dfa.GetAcceptingStates()

	// Create a distinguishability table
	table := make([][]bool, stateCount)
	for i := range table {
		table[i] = make([]bool, stateCount)
	}

	// Mark pairs where one state is accepting and the other is not
	for i := 0; i < stateCount; i++ {
		for j := i + 1; j < stateCount; j++ {
			if (acceptingStates[i] && !acceptingStates[j]) || (!acceptingStates[i] && acceptingStates[j]) {
				table[i][j] = true
				table[j][i] = true
			}
		}
	}

	// Iteratively mark distinguishable pairs
	changed := true
	for changed {
		changed = false
		for i := 0; i < stateCount; i++ {
			for j := i + 1; j < stateCount; j++ {
				if !table[i][j] {
					// Check if these states should be marked as distinguishable
					if shouldMarkDistinguishable(dfa, table, i, j) {
						table[i][j] = true
						table[j][i] = true
						changed = true
					}
				}
			}
		}
	}

	// Group equivalent states
	groups := make([][]int, 0)
	processed := make([]bool, stateCount)

	for i := 0; i < stateCount; i++ {
		if processed[i] {
			continue
		}

		group := []int{i}
		processed[i] = true

		for j := i + 1; j < stateCount; j++ {
			if !processed[j] && !table[i][j] {
				group = append(group, j)
				processed[j] = true
			}
		}
		groups = append(groups, group)
	}

	// Create mappings from old states to new states
	oldToNew := make(map[int]int)
	newToOld := make(map[int]int)

	builder := NewNFABuilder()

	// Create new states for each group
	for _, group := range groups {
		newState := builder.AddState()
		representative := group[0] // Use first state as representative
		newToOld[newState] = representative

		for _, oldState := range group {
			oldToNew[oldState] = newState
		}
	}

	// Set start state
	builder.SetStartState(oldToNew[dfa.GetStartState()])

	// Mark accepting states
	for oldState := range acceptingStates {
		if newState, exists := oldToNew[oldState]; exists {
			builder.AcceptState(newState)
		}
	}

	// Add transitions
	transitionsBySource := dfa.GetTransitionsBySource()
	for newState, oldState := range newToOld {
		if targets, exists := transitionsBySource[oldState]; exists {
			for info := range targets.AllTransitions() {
				if info.CharRange != nil {
					charset := rangeToRuneSet(info.CharRange)
					// Since this is a DFA, there should be only one target per character
					if len(info.Targets) > 0 {
						targetOldState := info.Targets[0]
						if targetNewState, exists := oldToNew[targetOldState]; exists {
							builder.AddTransition(newState, targetNewState, charset)
						}
					}
				}
			}
		}
	}

	result, err := builder.Build()
	if err != nil {
		panic("Failed to build minimized DFA: " + err.Error())
	}
	return result
}

// shouldMarkDistinguishable checks if two states should be marked as distinguishable
// based on their transitions
func shouldMarkDistinguishable(dfa NFA, table [][]bool, state1, state2 int) bool {
	transitionsBySource := dfa.GetTransitionsBySource()

	// Get all possible input characters from both states
	inputs1 := getInputCharSets(transitionsBySource, state1)
	inputs2 := getInputCharSets(transitionsBySource, state2)

	// Combine and deduplicate inputs
	allInputs := make(map[string]*RuneSet)
	for _, charset := range inputs1 {
		allInputs[charsetKey(charset)] = charset
	}
	for _, charset := range inputs2 {
		allInputs[charsetKey(charset)] = charset
	}

	// Check each input character
	for _, charset := range allInputs {
		targets1 := getTargets(transitionsBySource, state1, charset)
		targets2 := getTargets(transitionsBySource, state2, charset)

		// If different number of targets, states are distinguishable
		if len(targets1) != len(targets2) {
			return true
		}

		// For DFA, there should be exactly one target per input
		if len(targets1) == 1 && len(targets2) == 1 {
			target1 := targets1[0]
			target2 := targets2[0]

			// If target states are already marked as distinguishable,
			// then these states are distinguishable
			if target1 != target2 && table[target1][target2] {
				return true
			}
		} else if len(targets1) == 0 && len(targets2) == 0 {
			// Both states have no transition for this input - continue
			continue
		} else {
			// Different behavior (one has transition, other doesn't)
			return true
		}
	}

	return false
}

// Helper functions for working with character sets and transitions

func getInputCharSets(transitionsBySource map[int]TransitionTargets, state int) []*RuneSet {
	var charSets []*RuneSet
	if targets, exists := transitionsBySource[state]; exists {
		for info := range targets.AllTransitions() {
			if info.CharRange != nil {
				charSets = append(charSets, rangeToRuneSet(info.CharRange))
			}
		}
	}
	return charSets
}

func getTargets(transitionsBySource map[int]TransitionTargets, state int, charset *RuneSet) []int {
	if targets, exists := transitionsBySource[state]; exists {
		return targets.GetTargets(charset)
	}
	return nil
}

func unionCharSets(charSets ...*RuneSet) *RuneSet {
	if len(charSets) == 0 {
		return NewRuneSet_Empty()
	}
	result := charSets[0]
	for i := 1; i < len(charSets); i++ {
		result = Union(result, charSets[i])
	}
	return result
}

func charSetsEqual(a, b *RuneSet) bool {
	// Two charsets are equal if their difference is empty in both directions
	diff1 := Except(a, b)
	diff2 := Except(b, a)
	return diff1.Length() == 0 && diff2.Length() == 0
}

// charsetKey creates a string representation of a RuneSet for use as map key
func charsetKey(charset *RuneSet) string {
	// Simple string representation using ranges
	key := ""
	for _, r := range charset.Ranges {
		if r.Includes {
			key += fmt.Sprintf("[%d-%d]", r.Start, r.End)
		}
	}
	return key
}
