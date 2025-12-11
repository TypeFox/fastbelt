package automatons

import "fmt"

type ReducerState struct {
	State       int
	Index       int
	Input       string
	AcceptedIdx int
	Halted      bool
}

func (nfa NFAImpl) InitializeReducerState(input string) ReducerState {
	start := nfa.GetStartState()
	return ReducerState{
		State:       start,
		Index:       0,
		Input:       input,
		AcceptedIdx: -1,
		Halted:      false,
	}
}

func (nfa NFAImpl) Step(state ReducerState) (ReducerState, error) {
	if state.Halted {
		return state, fmt.Errorf("cannot step from halted state")
	}

	if state.Index >= len(state.Input) {
		return ReducerState{
			State:       state.State,
			Index:       state.Index,
			Input:       state.Input,
			AcceptedIdx: state.AcceptedIdx,
			Halted:      true,
		}, nil
	}

	runeSet := NewRuneSet_Rune(rune(state.Input[state.Index]))
	nextStates := nfa.GetTransitionsBySource()[state.State].GetTargets(runeSet)
	if len(nextStates) == 0 {
		return ReducerState{
			State:       state.State,
			Index:       state.Index,
			Input:       state.Input,
			AcceptedIdx: state.AcceptedIdx,
			Halted:      true,
		}, nil
	}
	if len(nextStates) > 1 {
		return state, fmt.Errorf("automaton is non-deterministic: multiple next states for state %d on input '%c'", state.State, state.Input[state.Index])
	}
	// For DFA, there should be exactly one next state
	nextState := nextStates[0]
	nextIndex := state.Index + 1
	acceptedIdx := state.AcceptedIdx
	if nfa.GetAcceptingStates()[nextState] {
		acceptedIdx = nextIndex
	}

	return ReducerState{
		State:       nextState,
		Index:       nextIndex,
		Input:       state.Input,
		AcceptedIdx: acceptedIdx,
		Halted:      false,
	}, nil
}
