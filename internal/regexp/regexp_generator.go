package regexp

import (
	"fmt"
	"slices"

	"typefox.dev/fastbelt/internal/automatons"
	"typefox.dev/fastbelt/util/codegen"
)

type GenerateTransitionsUsingBinarySearchResult struct {
	Lookup codegen.Node
	Next   codegen.Node
	Code   codegen.Node
}

func shiftStatesByDeadState(state int, deadState int) int {
	//With this trick we transform the DFA to an NFA without dead state again.
	//This makes the generated code smaller and faster, because we don't have to check for the dead state at all.
	//1. get the DFA dead state (if any)
	//If any:
	//  2. remove the lines corresponding to the dead state in arrays LOOKUP and NEXT
	//  3. correct all state usages as if the dead state has never existed
	//     (start state, next states, accepting states)
	if deadState > -1 && state > deadState {
		return state - 1
	}
	return state
}

func GenerateTransitionsUsingBinarySearch(bySource *automatons.RuneRangeTargetsMapping, source int, tokenName string, imports map[string]bool, deadState int) GenerateTransitionsUsingBinarySearchResult {
	lookup := codegen.NewNode()
	next := codegen.NewNode()
	transitions := make([]automatons.RuneRangeMappingSection[automatons.Targets], 0)
	if source != deadState {
		for transition := range bySource.All() {
			if transition.Values[0] == deadState {
				continue
			}
			transitions = append(transitions, transition)
		}

		lookup.Append("{")
		for _, transition := range transitions {
			lookup.Append(fmt.Sprintf("%s, ", automatons.FormatLowHighInts(transition.Range.Start, transition.Range.End)))
		}
		lookup.AppendLine("},")

		next.Append("{")
		for _, transition := range transitions {
			target := transition.Values[0]
			target = shiftStatesByDeadState(target, deadState)
			next.Append(fmt.Sprintf("%d, ", target))
		}
		next.AppendLine("},")
	}

	n := codegen.NewNode()
	n.AppendLine("nextState := -1")
	n.AppendLine(fmt.Sprintf("next := %s_Next[%d]", tokenName, shiftStatesByDeadState(source, deadState)))
	n.AppendLine(fmt.Sprintf("lookup := %s_Lookup[%d]", tokenName, shiftStatesByDeadState(source, deadState)))

	if len(transitions) < 16 {
		n.AppendLine("for i, lowHigh := range lookup {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("if rune(lowHigh&0xFFFFFFFF) <= r && r <= rune(lowHigh>>32) {")
			n.Indent(func(n codegen.Node) {
				n.AppendLine("nextState = next[i]")
				n.AppendLine("break")
			})
			n.AppendLine("}")
		})
		n.AppendLine("}")
	} else {
		imports["sort"] = true
		n.AppendLine("searchIndex := sort.Search(len(next), func(i int) bool {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("return rune(lookup[i] & 0xFFFFFFFF) > r")
		})
		n.AppendLine("}) - 1")
		n.AppendLine("if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("nextState = next[searchIndex]")
		})
		n.AppendLine("}")
	}
	n.AppendLine("if nextState > -1 {")
	n.Indent(func(n codegen.Node) {
		n.AppendLine("state = nextState")
	})
	n.AppendLine("} else {")
	n.Indent(func(n codegen.Node) {
		n.AppendLine("break loop")
	})
	n.AppendLine("}")
	return GenerateTransitionsUsingBinarySearchResult{
		Lookup: lookup,
		Next:   next,
		Code:   n,
	}
}

type GenerateRegExpResult struct {
	Imports map[string]bool
	Vars    codegen.Node
	Code    codegen.Node
}

func (r *RegexpImpl) GenerateRegExp(funcName string, tokenName string) GenerateRegExpResult {
	deadState := r.dfa.DeadState()

	vars := codegen.NewNode()
	lookup := codegen.NewNode()
	lookup.AppendLine(fmt.Sprintf("var %s_Lookup = [][]int64{", tokenName))
	next := codegen.NewNode()
	next.AppendLine(fmt.Sprintf("var %s_Next = [][]int{", tokenName))
	accepting := codegen.NewNode()
	acceptingName := fmt.Sprintf("%s_Accepting", tokenName)
	accepting.AppendLine(fmt.Sprintf("var %s = [%d]bool{", acceptingName, r.dfa.StateCount-1))
	accepting.Indent(func(n codegen.Node) {
		acceptingStates := r.dfa.AcceptingStates
		stateIDs := make([]int, 0, len(acceptingStates))
		for state, isAccepting := range acceptingStates {
			if isAccepting {
				stateIDs = append(stateIDs, state)
			}
		}
		slices.Sort(stateIDs)
		for _, state := range stateIDs {
			n.AppendLine(fmt.Sprintf("%d: true,", shiftStatesByDeadState(state, deadState)))
		}
	})
	accepting.AppendLine("}")

	imports := map[string]bool{"unicode/utf8": true}
	root := codegen.NewNode()
	root.AppendLine(fmt.Sprintf("func %s(s string, offset int) int {", funcName))
	root.Indent(func(n codegen.Node) {
		n.AppendLine("input := s[offset:]")
		n.AppendLine("length := len(input)")
		n.AppendLine(fmt.Sprintf("state := %d", shiftStatesByDeadState(r.dfa.StartState, deadState)))
		n.AppendLine("acceptedIndex := 0")
		n.AppendLine("index := 0")
		n.AppendLine("loop: for index < length {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("r, runeSize := utf8.DecodeRuneInString(input[index:])")
			n.AppendLine("switch state {")
			transitions := r.dfa.TransitionsBySource
			for source := 0; source < r.dfa.StateCount; source++ {
				if source == deadState {
					continue
				}
				bySource := transitions[source]
				n.AppendLine(fmt.Sprintf("case %d:", shiftStatesByDeadState(source, deadState)))
				n.Indent(func(n codegen.Node) {
					result := GenerateTransitionsUsingBinarySearch(bySource, source, tokenName, imports, deadState)
					n.AppendNode(result.Code)
					lookup.AppendNode(result.Lookup)
					next.AppendNode(result.Next)
				})
			}
			n.AppendLine("default:")
			n.Indent(func(n codegen.Node) {
				n.AppendLine("break loop")
			})
			n.AppendLine("}")
			n.AppendLine("index += runeSize")
			n.AppendLine(fmt.Sprintf("if %s[state] {", acceptingName))
			n.Indent(func(n codegen.Node) {
				n.AppendLine("acceptedIndex = index")
			})
			n.AppendLine("}")
		})
		n.AppendLine("}")
		n.AppendLine("return acceptedIndex")
	})
	root.Append("}")
	lookup.AppendLine("}")
	next.AppendLine("}")
	vars.AppendNode(lookup)
	vars.AppendNode(next)
	vars.AppendNode(accepting)
	return GenerateRegExpResult{
		Imports: imports,
		Code:    root,
		Vars:    vars,
	}
}
