package automatons

func bitMaskToStates(dfaState string) []int {
	buffer := []byte(dfaState)
	states := make([]int, 0)
	for index, b := range buffer {
		if b == '1' {
			states = append(states, index)
		}
	}
	return states
}

type transitionInfo struct {
	Source string
	Target string
	Range  RuneRange
}

// Determinize converts an NFA to a DFA using the subset construction algorithm
func (nfa *NFA) Determinize() *NFA {
	builder := NewNFABuilder()
	dfaStateMapping := make(map[string]int) // Using BitMask.String() as key

	startDFAState := nfa.GetEpsilonClosure(nfa.StartState).String()

	transitions := make([]transitionInfo, 0)
	queue := []string{startDFAState}
	for len(queue) > 0 {
		sourceDFAState := queue[0]
		queue = queue[1:]

		if _, exists := dfaStateMapping[sourceDFAState]; exists {
			continue
		}

		newState := builder.AddState()
		dfaStateMapping[sourceDFAState] = newState

		sourceStates := bitMaskToStates(sourceDFAState)
		acceptingStates := nfa.AcceptingStates
		for _, state := range sourceStates {
			if acceptingStates[state] {
				builder.AcceptState(newState)
				break
			}
		}

		nfaTargets := NewRuneRangeTargets()
		for _, sourceState := range sourceStates {
			xxx := nfa.TransitionsBySource[sourceState]
			if xxx != nil {
				xxx.MergeNonEpsilonInto(&nfaTargets.RuneRangeMappingBase)
			}
		}
		for _, section := range nfaTargets.Ranges {
			if section.Range.Includes {
				targetDFAState := nfa.GetEpsilonClosure(section.Values...).String()
				transitions = append(transitions, transitionInfo{
					Source: sourceDFAState,
					Target: targetDFAState,
					Range:  *section.Range,
				})
				queue = append(queue, targetDFAState)
			}
		}
	}

	builder.startState = dfaStateMapping[startDFAState]

	for _, trans := range transitions {
		sourceState := dfaStateMapping[trans.Source]
		targetState := dfaStateMapping[trans.Target]
		builder.AddTransitionForRuneRange(sourceState, targetState, &trans.Range)
	}

	result, err := builder.Build()
	if err != nil {
		panic("Failed to build DFA: " + err.Error())
	}

	return result
}
