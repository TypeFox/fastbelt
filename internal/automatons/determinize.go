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
			transitions := nfa.TransitionsBySource[sourceState]
			if transitions != nil {
				transitions.MergeNonEpsilonInto(&nfaTargets.RuneRangeMappingBase)
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

	// Complete the DFA: add a dead/trap state (lazily) and route every state's
	// uncovered input characters to it so the transition function is total.
	// The dead state is non-accepting, which is what complement relies on.
	var dead int = -1
	coveredByState := make(map[int]*RuneSet)
	for _, trans := range transitions {
		sourceState := dfaStateMapping[trans.Source]
		prev := coveredByState[sourceState]
		if prev == nil {
			prev = NewRuneSetEmpty()
		}
		ch := NewRuneSetRange(trans.Range.Start, trans.Range.End)
		coveredByState[sourceState] = Union(prev, ch)
	}

	for _, newState := range dfaStateMapping {
		covered := coveredByState[newState]
		if covered == nil {
			covered = NewRuneSetEmpty()
		}
		gap := Negate(covered)
		if gap.Length() > 0 {
			if dead == -1 {
				dead = builder.AddState()
				builder.AddTransitionForRuneSet(dead, dead, NewRuneSetFull())
			}
			builder.AddTransitionForRuneSet(newState, dead, gap)
		}
	}

	return builder.Build()
}
