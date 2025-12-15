package automatons

import (
	"slices"
)

func statesToBitMask(states []int, stateCount int) BitMask {
	mask := NewBitMask_Empty(stateCount)
	for _, state := range states {
		mask.Set(state)
	}
	return mask
}

func bitMaskToStates(mask BitMask) []int {
	states := make([]int, 0)
	for index := range len(mask) {
		if mask.IsSet(index) {
			states = append(states, index)
		}
	}
	return states
}

func rangeToRuneSet(runeRange *RuneRange) *RuneSet {
	if runeRange == nil {
		return NewRuneSet_Empty()
	}
	if runeRange.Includes {
		return NewRuneSet_Range(runeRange.Start, runeRange.End)
	}
	full := NewRuneSet_Full()
	excluded := NewRuneSet_Range(runeRange.Start, runeRange.End)
	return Except(full, excluded)
}

// Determinize converts an NFA to a DFA using the subset construction algorithm
func Determinize(nfa NFA) NFA {
	builder := NewNFABuilder()

	// Map from BitMask to new DFA state ID
	newStates := make(map[string]int) // Using BitMask.String() as key

	// Get initial epsilon closure of start state
	startClosure := GetEpsilonClosure(nfa, nfa.GetStartState())
	startHash := startClosure.String()

	// Queue of state sets to process (as BitMasks)
	queue := []BitMask{startClosure}

	// Transitions to be added after all states are created
	type transitionInfo struct {
		sourceDFAState string
		targetDFAState string
		charset        *RuneSet
	}
	transitions := make([]transitionInfo, 0)

	// Process queue
	for len(queue) > 0 {
		// Pop from front
		sourceSet := queue[0]
		queue = queue[1:]

		sourceDFAState := sourceSet.String()

		// Skip if already processed
		if _, exists := newStates[sourceDFAState]; exists {
			continue
		}

		// Create new DFA state
		newState := builder.AddState()
		newStates[sourceDFAState] = newState

		// Set as start state if this is the initial set
		if sourceDFAState == startHash {
			builder.SetStartState(newState)
		}

		// Check if any of the NFA states in this set are accepting
		sourceStates := bitMaskToStates(sourceSet)
		acceptingStates := nfa.GetAcceptingStates()
		for _, state := range sourceStates {
			if acceptingStates[state] {
				builder.AcceptState(newState)
				break
			}
		}

		// Get all possible input character sets by partitioning the alphabet
		inputs := getDistinctInputs(nfa, sourceStates)

		// For each distinct input character set, find target states
		for _, inputCharSet := range inputs {
			if inputCharSet.Length() == 0 {
				continue
			}

			// Collect all target states for this input
			targetStates := make(map[int]bool)
			for _, sourceState := range sourceStates {
				transitionTargets := nfa.GetTransitionsBySource()[sourceState]
				if transitionTargets != nil {
					targets := transitionTargets.GetTargets(inputCharSet)
					for _, target := range targets {
						targetStates[target] = true
					}
				}
			}

			if len(targetStates) > 0 {
				// Convert target states to slice and get epsilon closure
				targetSlice := make([]int, 0, len(targetStates))
				for target := range targetStates {
					targetSlice = append(targetSlice, target)
				}

				epsilonTarget := GetEpsilonClosure(nfa, targetSlice...)
				targetDFAState := epsilonTarget.String()

				// Record transition to be added later
				transitions = append(transitions, transitionInfo{
					sourceDFAState: sourceDFAState,
					targetDFAState: targetDFAState,
					charset:        inputCharSet,
				})

				// Add target set to queue if not already processed
				if _, exists := newStates[targetDFAState]; !exists {
					queue = append(queue, epsilonTarget)
				}
			}
		}
	}

	// Add all transitions
	for _, trans := range transitions {
		sourceState := newStates[trans.sourceDFAState]
		targetState := newStates[trans.targetDFAState]
		builder.AddTransition(sourceState, targetState, trans.charset)
	}

	result, err := builder.Build()
	if err != nil {
		panic("Failed to build DFA: " + err.Error())
	}

	return result
}

// getDistinctInputs partitions the alphabet based on all transition character sets
// This creates a set of disjoint character sets that cover all possible inputs
func getDistinctInputs(nfa NFA, states []int) []*RuneSet {
	var inputs []*RuneSet

	// Collect all character ranges from transitions
	for _, state := range states {
		transitionTargets := nfa.GetTransitionsBySource()[state]
		if transitionTargets == nil {
			continue
		}

		// Process all character transitions (skip epsilon)
		for transitionInfo := range transitionTargets.AllTransitions() {
			if transitionInfo.CharRange == nil {
				continue // Skip epsilon transitions
			}

			lhs := rangeToRuneSet(transitionInfo.CharRange)
			if lhs.Length() == 0 {
				continue
			}

			// Partition existing inputs with the new character set
			distinct := true
			newInputs := make([]*RuneSet, 0)

			for _, rhs := range inputs {
				first := Except(lhs, rhs)    // lhs - rhs
				second := Except(rhs, lhs)   // rhs - lhs
				third := Intersect(lhs, rhs) // lhs ∩ rhs

				if third.Length() > 0 {
					distinct = false
				}

				// Add non-empty sets
				if first.Length() > 0 {
					newInputs = append(newInputs, first)
				}
				if second.Length() > 0 {
					newInputs = append(newInputs, second)
				}
				if third.Length() > 0 {
					newInputs = append(newInputs, third)
				}
			}

			inputs = newInputs

			// If the character set is distinct from all existing ones, add it
			if distinct {
				inputs = append(inputs, lhs)
			}
		}
	}

	// Remove duplicates and sort for deterministic behavior
	return removeDuplicateCharSets(inputs)
}

// removeDuplicateCharSets removes duplicate character sets from the input slice
func removeDuplicateCharSets(inputs []*RuneSet) []*RuneSet {
	seen := make(map[string]bool)
	result := make([]*RuneSet, 0)

	for _, input := range inputs {
		key := input.String() // Use string representation as key
		if !seen[key] {
			seen[key] = true
			result = append(result, input)
		}
	}

	// Sort for deterministic output
	slices.SortFunc(result, func(a, b *RuneSet) int {
		if a.String() < b.String() {
			return -1
		} else if a.String() > b.String() {
			return 1
		}
		return 0
	})

	return result
}
