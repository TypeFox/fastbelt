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

func GenerateTransitionsUsingBinarySearch(bySource *automatons.RuneRangeTargetsMapping, source int, tokenName string, imports map[string]bool, isAcceptanceReachable map[int]bool) GenerateTransitionsUsingBinarySearchResult {
	transitions := make([]automatons.RuneRangeMappingSection[automatons.Targets], 0)
	if bySource != nil && isAcceptanceReachable[source] {
		for transition := range bySource.All() {
			if !isAcceptanceReachable[transition.Values[0]] {
				continue
			}
			transitions = append(transitions, transition)
		}
	}

	lookup := codegen.NewNode()
	lookup.Append("{")
	for _, transition := range transitions {
		lookup.Append(fmt.Sprintf("%s, ", automatons.FormatLowHighInts(transition.Range.Start, transition.Range.End)))
	}
	lookup.AppendLine("},")

	next := codegen.NewNode()
	next.Append("{")
	for _, transition := range transitions {
		next.Append(fmt.Sprintf("%d, ", transition.Values[0]))
	}
	next.AppendLine("},")

	n := codegen.NewNode()
	n.AppendLine("nextState := -1")
	n.AppendLine(fmt.Sprintf("next := %s_Next[%d]", tokenName, source))
	n.AppendLine(fmt.Sprintf("lookup := %s_Lookup[%d]", tokenName, source))

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
	vars := codegen.NewNode()
	lookup := codegen.NewNode()
	lookup.AppendLine(fmt.Sprintf("var %s_Lookup = [][]int64{", tokenName))
	next := codegen.NewNode()
	next.AppendLine(fmt.Sprintf("var %s_Next = [][]int{", tokenName))
	accepting := codegen.NewNode()
	acceptingName := fmt.Sprintf("%s_Accepting", tokenName)
	accepting.AppendLine(fmt.Sprintf("var %s = [%d]bool{", acceptingName, r.dfa.StateCount))
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
			n.AppendLine(fmt.Sprintf("%d: true,", state))
		}
	})
	accepting.AppendLine("}")

	imports := map[string]bool{"unicode/utf8": true}
	root := codegen.NewNode()
	root.AppendLine(fmt.Sprintf("func %s(s string, offset int) int {", funcName))
	root.Indent(func(n codegen.Node) {
		n.AppendLine("input := s[offset:]")
		n.AppendLine("length := len(input)")
		n.AppendLine(fmt.Sprintf("state := %d", r.dfa.StartState))
		n.AppendLine("acceptedIndex := 0")
		n.AppendLine("index := 0")
		n.AppendLine("loop: for index < length {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("r, runeSize := utf8.DecodeRuneInString(input[index:])")
			n.AppendLine("switch state {")
			transitions := r.dfa.TransitionsBySource
			isAcceptanceReachable := r.dfa.ComputeAcceptanceReachability()

			for source := 0; source < r.dfa.StateCount; source++ {
				bySource := transitions[source]
				n.AppendLine(fmt.Sprintf("case %d:", source))
				n.Indent(func(n codegen.Node) {
					result := GenerateTransitionsUsingBinarySearch(bySource, source, tokenName, imports, isAcceptanceReachable)
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
