package automatons

// Minimize minimizes a DFA using the table-filling algorithm
// This removes unreachable states and merges equivalent states
func Minimize(dfa NFA) NFA {
	stateCount := dfa.StateCount
	acceptingStates := dfa.AcceptingStates
	areDifferent := initializeDistinguishabilityTable(stateCount)
	prefillByAccept(stateCount, acceptingStates, areDifferent)
	for fillIterate(stateCount, areDifferent, dfa) {
	}
	groups := getEquivalentGroups(stateCount, areDifferent)
	return buildNewAutomaton(groups, dfa, acceptingStates)
}

func buildNewAutomaton(groups [][]int, dfa NFA, acceptingStates map[int]bool) NFA {
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
	builder.SetStartState(oldToNew[dfa.StartState])

	// Mark accepting states
	for oldState := range acceptingStates {
		if newState, exists := oldToNew[oldState]; exists {
			builder.AcceptState(newState)
		}
	}

	//Add transitions
	transitionsBySource := dfa.TransitionsBySource
	for newState, oldState := range newToOld {
		if targets, exists := transitionsBySource[oldState]; exists {
			for info := range targets.All() {
				if info.Range != nil {
					charset := NewRuneSet_Range(info.Range.Start, info.Range.End)
					// Since this is a DFA, there should be only one target per character
					if len(info.Values) > 0 {
						targetOldState := info.Values[0]
						if targetNewState, exists := oldToNew[targetOldState]; exists {
							builder.AddTransitionForRuneSet(newState, targetNewState, charset)
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
	return *result
}

func getEquivalentGroups(stateCount int, areDifferent [][]bool) [][]int {
	groups := make([][]int, 0)
	processed := make([]bool, stateCount)
	for i := 0; i < stateCount; i++ {
		if processed[i] {
			continue
		}

		group := []int{i}
		processed[i] = true

		for j := i + 1; j < stateCount; j++ {
			if !processed[j] && !areDifferent[i][j] {
				group = append(group, j)
				processed[j] = true
			}
		}
		groups = append(groups, group)
	}
	return groups
}

func fillIterate(stateCount int, areDifferent [][]bool, dfa NFA) bool {
	changed := false
	for i := 0; i < stateCount; i++ {
		for j := i + 1; j < stateCount; j++ {
			if !areDifferent[i][j] {
				// Check if these states should be marked as distinguishable
				if shouldMarkDistinguishable(dfa, areDifferent, i, j) {
					areDifferent[i][j] = true
					areDifferent[j][i] = true
					changed = true
				}
			}
		}
	}
	return changed
}

func prefillByAccept(stateCount int, acceptingStates map[int]bool, areDifferent [][]bool) {
	for i := 0; i < stateCount; i++ {
		for j := i + 1; j < stateCount; j++ {
			if (acceptingStates[i] && !acceptingStates[j]) || (!acceptingStates[i] && acceptingStates[j]) {
				areDifferent[i][j] = true
				areDifferent[j][i] = true
			}
		}
	}
}

func initializeDistinguishabilityTable(stateCount int) [][]bool {
	result := make([][]bool, stateCount)
	for i := range result {
		result[i] = make([]bool, stateCount)
	}
	return result
}

func shouldMarkDistinguishable(dfa NFA, areDifferent [][]bool, state1, state2 int) bool {
	targets := NewRuneRangeTargets()
	dfa.TransitionsBySource[state1].MergeNonEpsilonInto(&targets.RuneRangeMappingBase)
	dfa.TransitionsBySource[state2].MergeNonEpsilonInto(&targets.RuneRangeMappingBase)

	for _, section := range targets.Ranges {
		if !section.Range.Includes {
			continue
		}
		if len(section.Values) != 2 {
			return true
		}
		target1 := section.Values[0]
		target2 := section.Values[1]

		if areDifferent[target1][target2] {
			return true
		}
	}

	return false
}
