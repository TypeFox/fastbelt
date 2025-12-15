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
	transitions     map[int]*NFATargets
	acceptingStates map[int]bool
}

// NewNFABuilder creates a new NFA builder
func NewNFABuilder() *NFABuilderImpl {
	return &NFABuilderImpl{
		stateCounter:    0,
		startState:      -1,
		transitions:     make(map[int]*NFATargets),
		acceptingStates: make(map[int]bool),
	}
}

// AddState adds a new state and returns its ID
func (builder *NFABuilderImpl) AddState() int {
	stateID := builder.stateCounter
	builder.stateCounter++
	return stateID
}

func (builder *NFABuilderImpl) AddTransitionForRuneSet(source int, target int, characters *RuneSet) error {
	if source < 0 || source >= builder.stateCounter {
		return fmt.Errorf("invalid source state: %d", source)
	}
	if target < 0 || target >= builder.stateCounter {
		return fmt.Errorf("invalid target state: %d", target)
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewNFATargets()
	}

	targets := builder.transitions[source]
	if characters == nil {
		targets.AddEpsilonTargets(target)
		return nil
	}
	for _, rng := range characters.Ranges {
		if rng.Includes {
			targets.AddRuneRangeTargets(rng.Start, rng.End, target)
		}
	}
	return nil
}

func (builder *NFABuilderImpl) AddTransitionForRuneRange(source int, target int, rng *RuneRange) error {
	if source < 0 || source >= builder.stateCounter {
		return fmt.Errorf("invalid source state: %d", source)
	}
	if target < 0 || target >= builder.stateCounter {
		return fmt.Errorf("invalid target state: %d", target)
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewNFATargets()
	}

	targets := builder.transitions[source]
	if rng == nil {
		targets.AddEpsilonTargets(target)
		return nil
	}
	targets.AddRuneRangeTargets(rng.Start, rng.End, target)
	return nil
}

func (builder *NFABuilderImpl) AddTransitionForSingleRune(source int, target int, rn rune) error {
	if source < 0 || source >= builder.stateCounter {
		return fmt.Errorf("invalid source state: %d", source)
	}
	if target < 0 || target >= builder.stateCounter {
		return fmt.Errorf("invalid target state: %d", target)
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewNFATargets()
	}

	targets := builder.transitions[source]
	targets.AddRuneRangeTargets(rn, rn, target)
	return nil
}

// SetStartState sets the start state
func (builder *NFABuilderImpl) SetStartState(state int) error {
	if state < 0 || state >= builder.stateCounter {
		return fmt.Errorf("invalid state: %d", state)
	}
	builder.startState = state
	return nil
}

// AcceptState marks a state as accepting
func (builder *NFABuilderImpl) AcceptState(state int) error {
	if state < 0 || state >= builder.stateCounter {
		return fmt.Errorf("invalid state: %d", state)
	}
	builder.acceptingStates[state] = true
	return nil
}

// Build constructs and returns the NFA
func (builder *NFABuilderImpl) Build() (*NFA, error) {
	if builder.stateCounter == 0 {
		return nil, fmt.Errorf("no states defined")
	}
	if builder.startState == -1 {
		return nil, fmt.Errorf("no start state defined")
	}

	nfa := &NFA{
		StartState:          builder.startState,
		StateCount:          builder.stateCounter,
		TransitionsBySource: CopyMap(builder.transitions),
		AcceptingStates:     CopyMap(builder.acceptingStates),
	}

	return nfa, nil
}

// CopyFrom copies states and transitions from another NFA into this builder
func (builder *NFABuilderImpl) CopyFrom(nfa NFA) (*StateMapping, error) {
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
			for _, target := range section.Targets {
				targetState := stateMapping[target]
				err := builder.AddTransitionForRuneRange(sourceState, targetState, section.Range)
				if err != nil {
					return nil, fmt.Errorf("failed to copy transition: %v", err)
				}
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
	}, nil
}
