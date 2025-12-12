package automatons

import "fmt"

// NFABuilder interface defines methods for building an NFA
type NFABuilder interface {
	AddState() int
	AddTransition(source int, target int, characters *RuneSet) error
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
	transitions     map[int]*TransitionTargetsImpl
	acceptingStates map[int]bool
}

// NewNFABuilder creates a new NFA builder
func NewNFABuilder() *NFABuilderImpl {
	return &NFABuilderImpl{
		stateCounter:    0,
		startState:      -1,
		transitions:     make(map[int]*TransitionTargetsImpl),
		acceptingStates: make(map[int]bool),
	}
}

// AddState adds a new state and returns its ID
func (builder *NFABuilderImpl) AddState() int {
	stateID := builder.stateCounter
	builder.stateCounter++
	return stateID
}

// AddTransition adds a transition from source to target with the given character set
func (builder *NFABuilderImpl) AddTransition(source int, target int, characters *RuneSet) error {
	if source < 0 || source >= builder.stateCounter {
		return fmt.Errorf("invalid source state: %d", source)
	}
	if target < 0 || target >= builder.stateCounter {
		return fmt.Errorf("invalid target state: %d", target)
	}

	if builder.transitions[source] == nil {
		builder.transitions[source] = NewTransitionTargets()
	}

	targets := builder.transitions[source]
	targets.Add(characters, target)
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
func (builder *NFABuilderImpl) Build() (NFA, error) {
	if builder.stateCounter == 0 {
		return nil, fmt.Errorf("no states defined")
	}
	if builder.startState == -1 {
		return nil, fmt.Errorf("no start state defined")
	}

	nfa := NewNFA(builder.startState, builder.stateCounter)

	// Copy accepting states
	for state := range builder.acceptingStates {
		nfa.AddAcceptingState(state)
	}

	// Copy transitions
	for source, targets := range builder.transitions {
		nfa.transitionsBySource[source] = targets
	}

	return nfa, nil
}

// CopyFrom copies states and transitions from another NFA into this builder
func (builder *NFABuilderImpl) CopyFrom(nfa NFA) (*StateMapping, error) {
	stateMapping := make(map[int]int)

	// Create new states for each state in the source NFA
	for state := 0; state < nfa.GetStateCount(); state++ {
		newState := builder.AddState()
		stateMapping[state] = newState
	}

	// Copy transitions
	for transition := range nfa.AllTransitions() {
		sourceState := stateMapping[transition.Source]
		targetState := stateMapping[transition.Target]
		err := builder.AddTransition(sourceState, targetState, transition.Chars)
		if err != nil {
			return nil, fmt.Errorf("failed to copy transition: %v", err)
		}
	}

	// Collect accepting states
	acceptings := make([]int, 0)
	for state := range nfa.GetAcceptingStates() {
		if newState, exists := stateMapping[state]; exists {
			acceptings = append(acceptings, newState)
		}
	}

	return &StateMapping{
		Mapping:    stateMapping,
		Start:      stateMapping[nfa.GetStartState()],
		Acceptings: acceptings,
	}, nil
}
