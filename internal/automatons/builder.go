package automatons

import "fmt"

type NFABuilder interface {
	AddState() int
	AddTransitionForRuneSet(source int, target int, characters RuneSet) error
	AddTransitionForRuneRange(source int, target int, rng RuneRange) error
	SetStartState(state int) error
	AcceptState(state int) error
	Build() (NFA, error)
	CopyFrom(nfa NFA) (*StateMapping, error)
}

// StateMapping represents the mapping of states when copying from another NFA
type StateMapping struct {
	Mapping    map[int]int
	Start      int
	Acceptings []int
}

// NFABuilderImpl implements the NFABuilder interface
type NFABuilderImpl struct {
	stateCounter    int
	startState      int
	transitions     map[int]*RuneRangeTargetsMapping
	acceptingStates map[int]bool
}

// NewNFABuilder creates a new NFA builder
func NewNFABuilder() *NFABuilderImpl {
	return &NFABuilderImpl{
		stateCounter:    0,
		startState:      -1,
		transitions:     make(map[int]*RuneRangeTargetsMapping),
		acceptingStates: make(map[int]bool),
	}
}

// AddState adds a new state and returns its ID
func (builder *NFABuilderImpl) AddState() int {
	stateID := builder.stateCounter
	builder.stateCounter++
	return stateID
}

func (builder *NFABuilderImpl) AddTransitionForRuneSet(source int, target int, characters *RuneSet) {
	if source < 0 || source >= builder.stateCounter {
		panic(fmt.Sprintf("invalid source state: %d", source))
	}
	if target < 0 || target >= builder.stateCounter {
		panic(fmt.Sprintf("invalid target state: %d", target))
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewRuneRangeTargets()
	}

	targets := builder.transitions[source]
	if characters == nil {
		targets.AddEpsilonValues(Targets{target})
	} else {
		for _, rng := range characters.Ranges {
			if rng.Includes {
				targets.AddRuneRangeValues(rng.Start, rng.End, Targets{target})
			}
		}
	}
}

func (builder *NFABuilderImpl) AddTransitionForRuneRange(source int, target int, rng *RuneRange) {
	if source < 0 || source >= builder.stateCounter {
		panic(fmt.Sprintf("invalid source state: %d", source))
	}
	if target < 0 || target >= builder.stateCounter {
		panic(fmt.Sprintf("invalid target state: %d", target))
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewRuneRangeTargets()
	}

	targets := builder.transitions[source]
	if rng == nil {
		targets.AddEpsilonValues(Targets{target})
	} else {
		targets.AddRuneRangeValues(rng.Start, rng.End, Targets{target})
	}
}

func (builder *NFABuilderImpl) AddTransitionForSingleRune(source int, target int, rn rune) {
	if source < 0 || source >= builder.stateCounter {
		panic(fmt.Sprintf("invalid source state: %d", source))
	}
	if target < 0 || target >= builder.stateCounter {
		panic(fmt.Sprintf("invalid target state: %d", target))
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewRuneRangeTargets()
	}

	targets := builder.transitions[source]
	targets.AddRuneRangeValues(rn, rn, Targets{target})
}

func (builder *NFABuilderImpl) SetStartState(state int) {
	if state < 0 || state >= builder.stateCounter {
		panic(fmt.Sprintf("invalid state: %d", state))
	}
	builder.startState = state
}

func (builder *NFABuilderImpl) AcceptState(state int) {
	if state < 0 || state >= builder.stateCounter {
		panic(fmt.Sprintf("invalid state: %d", state))
	}
	builder.acceptingStates[state] = true
}

func (builder *NFABuilderImpl) Build() *NFA {
	if builder.stateCounter == 0 {
		panic("no states defined")
	}
	if builder.startState == -1 {
		panic("no start state defined")
	}

	return &NFA{
		StartState:          builder.startState,
		StateCount:          builder.stateCounter,
		TransitionsBySource: CopyMap(builder.transitions),
		AcceptingStates:     CopyMap(builder.acceptingStates),
	}
}

func (builder *NFABuilderImpl) CopyFrom(nfa *NFA) *StateMapping {
	stateMapping := make(map[int]int)

	// Create new states for each state in the source NFA
	for state := 0; state < nfa.StateCount; state++ {
		newState := builder.AddState()
		stateMapping[state] = newState
	}

	// Copy transitions
	for source, targets := range nfa.TransitionsBySource {
		sourceState := stateMapping[source]
		for section := range targets.All() {
			for _, target := range section.Values {
				targetState := stateMapping[target]
				builder.AddTransitionForRuneRange(sourceState, targetState, section.Range)
			}
		}
	}

	// Collect accepting states
	acceptings := make([]int, 0)
	for state := range nfa.AcceptingStates {
		if newState, exists := stateMapping[state]; exists {
			acceptings = append(acceptings, newState)
		}
	}

	return &StateMapping{
		Mapping:    stateMapping,
		Start:      stateMapping[nfa.StartState],
		Acceptings: acceptings,
	}
}
