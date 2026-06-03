package automatons

import (
	"fmt"
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
	for source, targets := range nfa.TransitionsBySource {
		for info := range targets.All() {
			var charset *RuneSet
			if info.Range != nil {
				charset = NewRuneSetRange(info.Range.Start, info.Range.End)
			} else {
				charset = NewRuneSetEmpty()
			}
			for _, target := range info.Values {
				result += fmt.Sprintf("  %d -> %d [label=\"%v\"];\n", source, target, charset)
			}
		}
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
