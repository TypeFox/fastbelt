package automatons

import "fmt"

type ConstructionKit struct{}

func NewConstructionKit() *ConstructionKit {
	return &ConstructionKit{}
}

func (ck *ConstructionKit) Empty() *NFA {
	builder := NewNFABuilder()
	start := builder.AddState()
	builder.SetStartState(start)
	builder.AcceptState(start)
	return builder.Build()
}

func (ck *ConstructionKit) Reject() *NFA {
	builder := NewNFABuilder()
	start := builder.AddState()
	builder.SetStartState(start)
	builder.AddTransitionForRuneSet(start, start, NewRuneSetFull())
	return builder.Build()
}

func (ck *ConstructionKit) Consume(characters *RuneSet) *NFA {
	builder := NewNFABuilder()
	start := builder.AddState()
	end := builder.AddState()

	builder.SetStartState(start)
	builder.AcceptState(end)
	builder.AddTransitionForRuneSet(start, end, characters)

	return builder.Build()
}

func (ck *ConstructionKit) Alternate(automata ...*NFA) *NFA {
	if len(automata) == 0 {
		panic("no automata provided for alternation")
	}

	builder := NewNFABuilder()
	start := builder.AddState()
	end := builder.AddState()

	builder.SetStartState(start)
	builder.AcceptState(end)

	for _, automaton := range automata {
		stateMapping := builder.CopyFrom(automaton)

		builder.AddTransitionForRuneSet(start, stateMapping.Start, nil)
		for _, accepting := range stateMapping.Acceptings {
			builder.AddTransitionForRuneSet(accepting, end, nil)
		}
	}

	return ck.finalize(builder)
}

func (*ConstructionKit) finalize(builder *NFABuilderImpl) *NFA {
	nfa := builder.Build()
	dfa := nfa.Determinize()
	return dfa.Minimize()
}

// Concat creates an NFA that matches the concatenation of the given automata
func (ck *ConstructionKit) Concat(automata ...*NFA) *NFA {
	if len(automata) == 0 {
		panic("no automata provided for concatenation")
	}

	builder := NewNFABuilder()
	start := builder.AddState()
	builder.SetStartState(start)

	currentStart := start

	for _, automaton := range automata {
		end := builder.AddState()

		stateMapping := builder.CopyFrom(automaton)

		builder.AddTransitionForRuneSet(currentStart, stateMapping.Start, nil)
		for _, accepting := range stateMapping.Acceptings {
			builder.AddTransitionForRuneSet(accepting, end, nil)
		}

		currentStart = end
	}

	builder.AcceptState(currentStart)

	return ck.finalize(builder)
}

// Repeat creates an NFA that matches the given automaton repeated min to max times
// If max is -1, there is no upper limit
func (ck *ConstructionKit) Repeat(automaton *NFA, min, max int) *NFA {
	if min < 0 || (max >= 0 && min > max) {
		panic(fmt.Sprintf("invalid range: min=%d, max=%d", min, max))
	}

	builder := NewNFABuilder()
	start := builder.AddState()
	accept := builder.AddState()
	previousStart := -1

	builder.SetStartState(start)
	builder.AcceptState(accept)

	var count *int
	if max != -1 {
		c := max - min
		count = &c
	}

	// Handle required repetitions (min)
	for min > 0 {
		end := builder.AddState()

		stateMapping := builder.CopyFrom(automaton)

		builder.AddTransitionForRuneSet(start, stateMapping.Start, nil)
		for _, accepting := range stateMapping.Acceptings {
			builder.AddTransitionForRuneSet(accepting, end, nil)
		}

		min--
		previousStart = start
		start = end
	}

	// Add epsilon transition to accept state
	builder.AddTransitionForRuneSet(start, accept, nil)

	if count == nil { // No upper limit, loop!
		if previousStart == -1 {
			end := builder.AddState()

			stateMapping := builder.CopyFrom(automaton)

			builder.AddTransitionForRuneSet(start, stateMapping.Start, nil)
			builder.AddTransitionForRuneSet(end, accept, nil)
			for _, accepting := range stateMapping.Acceptings {
				builder.AddTransitionForRuneSet(accepting, end, nil)
			}

			previousStart = start
			start = end
		}

		// Add loop back transition
		builder.AddTransitionForRuneSet(start, previousStart, nil)
	} else { // Existing upper limit, no loop!
		for *count > 0 {
			end := builder.AddState()

			stateMapping := builder.CopyFrom(automaton)

			builder.AddTransitionForRuneSet(start, stateMapping.Start, nil)
			for _, accepting := range stateMapping.Acceptings {
				builder.AddTransitionForRuneSet(accepting, end, nil)
			}
			builder.AddTransitionForRuneSet(end, accept, nil)

			*count--
			start = end
		}
	}

	return ck.finalize(builder)
}

// Complement creates an NFA that matches the complement of the given automaton
func (ck *ConstructionKit) Complement(automaton *NFA) *NFA {
	builder := NewNFABuilder()

	stateMapping := builder.CopyFrom(automaton)
	builder.SetStartState(stateMapping.Start)

	// Accept all non-accepting states from the original automaton
	acceptingStates := automaton.AcceptingStates
	for oldState, newState := range stateMapping.Mapping {
		if !acceptingStates[oldState] {
			builder.AcceptState(newState)
		}
	}

	return ck.finalize(builder)
}

// IntersectNFA creates an NFA that matches the intersection of two automata
// Uses De Morgan's law: A ∩ B = ¬(¬A ∪ ¬B)
func (ck *ConstructionKit) Intersect(a, b *NFA) *NFA {
	notA := ck.Complement(a)
	notB := ck.Complement(b)
	notAOrB := ck.Alternate(notA, notB)
	return ck.Complement(notAOrB)
}
