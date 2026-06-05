package automatons

import (
	"fmt"
	"sort"
)

type NFA struct {
	StartState          int
	StateCount          int
	AcceptingStates     map[int]bool
	TransitionsBySource map[int]*RuneRangeTargetsMapping
}

func (nfa NFA) String() string {
	result := "NFA:\n"
	result += "Start State: " + fmt.Sprintf("%d", nfa.StartState) + "\n"
	result += "Accepting States: "
	for state := range nfa.AcceptingStates {
		result += fmt.Sprintf("%d ", state)
	}
	result += "\nTransitions:\n"
	for source, targets := range nfa.TransitionsBySource {
		for info := range targets.All() {
			var charset *RuneSet
			if info.Range != nil {
				charset = NewRuneSetRange(info.Range.Start, info.Range.End)
			} else {
				charset = NewRuneSetEmpty()
			}
			for _, target := range info.Values {
				result += fmt.Sprintf("  %d --%v--> %d\n", source, charset, target)
			}
		}
	}
	return result
}

type dotFileTransitionInfo struct {
	Source int
	Target int
	Runes  *RuneSet
}

func (nfa NFA) DotFile() string {
	result := "digraph NFA {\n"
	result += "  rankdir=LR;\n"
	result += "  node [shape=circle];\n"

	// Start state
	result += "  start [shape=point];\n"
	result += fmt.Sprintf("  start -> %d;\n", nfa.StartState)

	// Accepting states
	for state := range nfa.AcceptingStates {
		result += fmt.Sprintf("  %d [shape=doublecircle];\n", state)
	}

	// Transitions
	transitions := make([]dotFileTransitionInfo, 0)
	for source, targets := range nfa.TransitionsBySource {
		for info := range targets.All() {
			var charset *RuneSet
			if info.Range != nil {
				charset = NewRuneSetRange(info.Range.Start, info.Range.End)
			} else {
				charset = NewRuneSetEmpty()
			}
			for _, target := range info.Values {
				transitions = append(transitions, dotFileTransitionInfo{
					Source: source,
					Target: target,
					Runes:  charset,
				})
			}
		}
	}
	//sort transitions by source and then target for deterministic output
	sort.Slice(transitions, func(i, j int) bool {
		if transitions[i].Source != transitions[j].Source {
			return transitions[i].Source < transitions[j].Source
		} else if transitions[i].Target != transitions[j].Target {
			return transitions[i].Target < transitions[j].Target
		} else {
			lhs := transitions[i].Runes.FirstRune()
			rhs := transitions[j].Runes.FirstRune()
			return lhs < rhs
		}
	})
	for _, transition := range transitions {
		result += fmt.Sprintf("  %d -> %d [label=\"%v\"];\n", transition.Source, transition.Target, transition.Runes)
	}

	result += "}\n"
	return result
}

func (nfa NFA) ComputeAcceptanceReachability() map[int]bool {
	revertedTransitions := make(map[int][]int)
	for source, targets := range nfa.TransitionsBySource {
		for info := range targets.All() {
			for _, target := range info.Values {
				revertedTransitions[target] = append(revertedTransitions[target], source)
			}
		}
	}

	canReach := make(map[int]bool)
	queue := []int{}

	for node := range nfa.AcceptingStates {
		canReach[node] = true
		queue = append(queue, node)
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, prev := range revertedTransitions[curr] {
			if !canReach[prev] {
				canReach[prev] = true
				queue = append(queue, prev)
			}
		}
	}

	return canReach
}

func (nfa NFA) DeadState() int {
	isAcceptanceReachable := nfa.ComputeAcceptanceReachability()
	deadState := -1
	for state := 0; state < nfa.StateCount; state++ {
		if !isAcceptanceReachable[state] {
			if deadState == -1 {
				deadState = state
			} else {
				panic(fmt.Sprintf("Multiple dead states found: %d and %d", deadState, state))
			}
		}
	}
	return deadState
}
