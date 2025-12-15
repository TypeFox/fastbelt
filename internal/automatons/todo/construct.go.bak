package automatons

import "fmt"

// ConstructionKit provides methods for constructing NFAs from various patterns
type ConstructionKit struct{}

// NewConstructionKit creates a new ConstructionKit instance
func NewConstructionKit() *ConstructionKit {
	return &ConstructionKit{}
}

func (ck *ConstructionKit) Empty() (NFA, error) {
	builder := NewNFABuilder()
	start := builder.AddState()
	builder.SetStartState(start)
	builder.AcceptState(start)
	return builder.Build()
}

func (ck *ConstructionKit) Reject() (NFA, error) {
	builder := NewNFABuilder()
	start := builder.AddState()
	builder.SetStartState(start)
	builder.AddTransition(start, start, NewRuneSet_Full())
	return builder.Build()
}

// Consume creates an NFA that matches the given character set
func (ck *ConstructionKit) Consume(characters *RuneSet) (NFA, error) {
	builder := NewNFABuilder()
	start := builder.AddState()
	end := builder.AddState()

	if err := builder.SetStartState(start); err != nil {
		return nil, fmt.Errorf("failed to set start state: %v", err)
	}

	if err := builder.AcceptState(end); err != nil {
		return nil, fmt.Errorf("failed to set accepting state: %v", err)
	}

	if err := builder.AddTransition(start, end, characters); err != nil {
		return nil, fmt.Errorf("failed to add transition: %v", err)
	}

	return builder.Build()
}

// Alternate creates an NFA that matches any of the given automata
func (ck *ConstructionKit) Alternate(automata ...NFA) (NFA, error) {
	if len(automata) == 0 {
		return nil, fmt.Errorf("no automata provided for alternation")
	}

	builder := NewNFABuilder()
	start := builder.AddState()
	end := builder.AddState()

	if err := builder.SetStartState(start); err != nil {
		return nil, fmt.Errorf("failed to set start state: %v", err)
	}

	if err := builder.AcceptState(end); err != nil {
		return nil, fmt.Errorf("failed to set accepting state: %v", err)
	}

	for _, automaton := range automata {
		stateMapping, err := builder.CopyFrom(automaton)
		if err != nil {
			return nil, fmt.Errorf("failed to copy automaton: %v", err)
		}

		// Add epsilon transition from new start to copied start
		if err := builder.AddTransition(start, stateMapping.Start, NewRuneSet_Empty()); err != nil {
			return nil, fmt.Errorf("failed to add start transition: %v", err)
		}

		// Add epsilon transitions from all accepting states to new end
		for _, accepting := range stateMapping.Acceptings {
			if err := builder.AddTransition(accepting, end, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add accepting transition: %v", err)
			}
		}
	}

	return ck.finalize(builder)
}

func (*ConstructionKit) finalize(builder *NFABuilderImpl) (NFA, error) {
	nfa, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build alternation NFA: %v", err)
	}

	// Minimize and determinize the result
	determinized := Determinize(nfa)
	return Minimize(determinized), nil
}

// Concat creates an NFA that matches the concatenation of the given automata
func (ck *ConstructionKit) Concat(automata ...NFA) (NFA, error) {
	if len(automata) == 0 {
		return nil, fmt.Errorf("no automata provided for concatenation")
	}

	builder := NewNFABuilder()
	start := builder.AddState()

	if err := builder.SetStartState(start); err != nil {
		return nil, fmt.Errorf("failed to set start state: %v", err)
	}

	currentStart := start

	for _, automaton := range automata {
		end := builder.AddState()

		stateMapping, err := builder.CopyFrom(automaton)
		if err != nil {
			return nil, fmt.Errorf("failed to copy automaton: %v", err)
		}

		// Add epsilon transition from current start to copied start
		if err := builder.AddTransition(currentStart, stateMapping.Start, NewRuneSet_Empty()); err != nil {
			return nil, fmt.Errorf("failed to add start transition: %v", err)
		}

		// Add epsilon transitions from all accepting states to end
		for _, accepting := range stateMapping.Acceptings {
			if err := builder.AddTransition(accepting, end, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add accepting transition: %v", err)
			}
		}

		currentStart = end
	}

	// Mark the final state as accepting
	if err := builder.AcceptState(currentStart); err != nil {
		return nil, fmt.Errorf("failed to set final accepting state: %v", err)
	}

	return ck.finalize(builder)
}

// Repeat creates an NFA that matches the given automaton repeated min to max times
// If max is -1, there is no upper limit
func (ck *ConstructionKit) Repeat(automaton NFA, min, max int) (NFA, error) {
	if min < 0 || (max >= 0 && min > max) {
		return nil, fmt.Errorf("invalid range: min=%d, max=%d", min, max)
	}

	builder := NewNFABuilder()
	start := builder.AddState()
	accept := builder.AddState()
	previousStart := -1

	if err := builder.SetStartState(start); err != nil {
		return nil, fmt.Errorf("failed to set start state: %v", err)
	}

	if err := builder.AcceptState(accept); err != nil {
		return nil, fmt.Errorf("failed to set accepting state: %v", err)
	}

	var count *int
	if max != -1 {
		c := max - min
		count = &c
	}

	// Handle required repetitions (min)
	for min > 0 {
		end := builder.AddState()

		stateMapping, err := builder.CopyFrom(automaton)
		if err != nil {
			return nil, fmt.Errorf("failed to copy automaton: %v", err)
		}

		if err := builder.AddTransition(start, stateMapping.Start, NewRuneSet_Empty()); err != nil {
			return nil, fmt.Errorf("failed to add start transition: %v", err)
		}

		for _, accepting := range stateMapping.Acceptings {
			if err := builder.AddTransition(accepting, end, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add accepting transition: %v", err)
			}
		}

		min--
		previousStart = start
		start = end
	}

	// Add epsilon transition to accept state
	if err := builder.AddTransition(start, accept, NewRuneSet_Empty()); err != nil {
		return nil, fmt.Errorf("failed to add transition to accept: %v", err)
	}

	if count == nil { // No upper limit, loop!
		if previousStart == -1 {
			end := builder.AddState()

			stateMapping, err := builder.CopyFrom(automaton)
			if err != nil {
				return nil, fmt.Errorf("failed to copy automaton: %v", err)
			}

			if err := builder.AddTransition(start, stateMapping.Start, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add start transition: %v", err)
			}

			if err := builder.AddTransition(end, accept, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add transition to accept: %v", err)
			}

			for _, accepting := range stateMapping.Acceptings {
				if err := builder.AddTransition(accepting, end, NewRuneSet_Empty()); err != nil {
					return nil, fmt.Errorf("failed to add accepting transition: %v", err)
				}
			}

			previousStart = start
			start = end
		}

		// Add loop back transition
		if err := builder.AddTransition(start, previousStart, NewRuneSet_Empty()); err != nil {
			return nil, fmt.Errorf("failed to add loop transition: %v", err)
		}
	} else { // Existing upper limit, no loop!
		for *count > 0 {
			end := builder.AddState()

			stateMapping, err := builder.CopyFrom(automaton)
			if err != nil {
				return nil, fmt.Errorf("failed to copy automaton: %v", err)
			}

			if err := builder.AddTransition(start, stateMapping.Start, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add start transition: %v", err)
			}

			for _, accepting := range stateMapping.Acceptings {
				if err := builder.AddTransition(accepting, end, NewRuneSet_Empty()); err != nil {
					return nil, fmt.Errorf("failed to add accepting transition: %v", err)
				}
			}

			if err := builder.AddTransition(end, accept, NewRuneSet_Empty()); err != nil {
				return nil, fmt.Errorf("failed to add transition to accept: %v", err)
			}

			*count--
			start = end
		}
	}

	return ck.finalize(builder)
}

// Complement creates an NFA that matches the complement of the given automaton
func (ck *ConstructionKit) Complement(automaton NFA) (NFA, error) {
	builder := NewNFABuilder()

	stateMapping, err := builder.CopyFrom(automaton)
	if err != nil {
		return nil, fmt.Errorf("failed to copy automaton: %v", err)
	}

	if err := builder.SetStartState(stateMapping.Start); err != nil {
		return nil, fmt.Errorf("failed to set start state: %v", err)
	}

	// Accept all non-accepting states from the original automaton
	acceptingStates := automaton.GetAcceptingStates()
	for oldState, newState := range stateMapping.Mapping {
		if !acceptingStates[oldState] {
			if err := builder.AcceptState(newState); err != nil {
				return nil, fmt.Errorf("failed to set accepting state: %v", err)
			}
		}
	}

	return ck.finalize(builder)
}

// IntersectNFA creates an NFA that matches the intersection of two automata
// Uses De Morgan's law: A ∩ B = ¬(¬A ∪ ¬B)
func (ck *ConstructionKit) IntersectNFA(a, b NFA) (NFA, error) {
	// Get complements
	notA, err := ck.Complement(a)
	if err != nil {
		return nil, fmt.Errorf("failed to complement first automaton: %v", err)
	}

	notB, err := ck.Complement(b)
	if err != nil {
		return nil, fmt.Errorf("failed to complement second automaton: %v", err)
	}

	// Get union of complements
	notAOrB, err := ck.Alternate(notA, notB)
	if err != nil {
		return nil, fmt.Errorf("failed to alternate complements: %v", err)
	}

	// Normalize the result
	normalized := Minimize(Determinize(notAOrB))

	// Return complement of the union
	return ck.Complement(normalized)
}

// Default instance for package-level functions
var defaultKit = NewConstructionKit()

// Package-level convenience functions

// Consume creates an NFA that matches the given character set
func Consume(characters *RuneSet) (NFA, error) {
	return defaultKit.Consume(characters)
}

// Alternate creates an NFA that matches any of the given automata
func Alternate(automata ...NFA) (NFA, error) {
	return defaultKit.Alternate(automata...)
}

// Concat creates an NFA that matches the concatenation of the given automata
func Concat(automata ...NFA) (NFA, error) {
	return defaultKit.Concat(automata...)
}

// Repeat creates an NFA that matches the given automaton repeated min to max times
func Repeat(automaton NFA, min, max int) (NFA, error) {
	return defaultKit.Repeat(automaton, min, max)
}

// Complement creates an NFA that matches the complement of the given automaton
func Complement(automaton NFA) (NFA, error) {
	return defaultKit.Complement(automaton)
}

// IntersectNFA creates an NFA that matches the intersection of two automata
func IntersectNFA(a, b NFA) (NFA, error) {
	return defaultKit.IntersectNFA(a, b)
}
