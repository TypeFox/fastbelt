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
				charset = NewRuneSet_Range(info.Range.Start, info.Range.End)
			} else {
				charset = NewRuneSet_Empty()
			}
			for _, target := range info.Values {
				result += fmt.Sprintf("  %d --%v--> %d\n", source, charset, target)
			}
		}
	}
	return result
}
