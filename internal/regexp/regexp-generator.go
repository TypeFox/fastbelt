package regexp

import (
	"fmt"
	"slices"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/automatons"
)

type GenerateTransitionsUsingBinarySearchResult struct {
	Lookup generator.Node
	Next   generator.Node
	Code   generator.Node
}

func GenerateTransitionsUsingBinarySearch(bySource *automatons.RuneRangeTargetsMapping, source int, tokenName string, imports map[string]bool) GenerateTransitionsUsingBinarySearchResult {
	transitions := make([]automatons.RuneRangeMappingSection[automatons.Targets], 0)
	if bySource != nil {
		for transition := range bySource.All() {
			transitions = append(transitions, transition)
		}
	}

	lookup := generator.NewNode()
	lookup.Append("{")
	for _, transition := range transitions {
		lookup.Append(fmt.Sprintf("%s, ", automatons.FormatLowHighInts(transition.Range.Start, transition.Range.End)))
	}
	lookup.AppendLine("},")

	next := generator.NewNode()
	next.Append("{")
	for _, transition := range transitions {
		next.Append(fmt.Sprintf("%d, ", transition.Values[0]))
	}
	next.AppendLine("},")

	n := generator.NewNode()
	n.AppendLine("nextState := -1")
	n.AppendLine(fmt.Sprintf("next := %s_Next[%d]", tokenName, source))
	n.AppendLine(fmt.Sprintf("lookup := %s_Lookup[%d]", tokenName, source))

	if len(transitions) < 16 {
		imports["slices"] = true
		n.AppendLine("searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("lo := rune(lowHigh & 0xFFFFFFFF)")
			n.AppendLine("hi := rune(lowHigh >> 32)")
			n.AppendLine("return lo <= r && r <= hi")
		})
		n.AppendLine("})")
		n.AppendLine("if searchIndex > -1 {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("nextState = next[searchIndex]")
		})
		n.AppendLine("}")
	} else {
		imports["sort"] = true
		n.AppendLine("searchIndex := sort.Search(len(next), func(i int) bool {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("return rune(lookup[i] & 0xFFFFFFFF) > r")
		})
		n.AppendLine("}) - 1")
		n.AppendLine("if searchIndex > -1 && rune(lookup[searchIndex] & 0xFFFFFFFF) <= r && r <= rune(lookup[searchIndex] >> 32) {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("nextState = next[searchIndex]")
		})
		n.AppendLine("}")
	}
	n.AppendLine("if nextState > -1 {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("state = nextState")
	})
	n.AppendLine("} else {")
	n.Indent(func(n generator.Node) {
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
	Lookup  generator.Node
	Next    generator.Node
	Code    generator.Node
}

func (r *RegexpImpl) GenerateRegExp(funcName string, tokenName string) GenerateRegExpResult {
	lookup := generator.NewNode()
	lookup.AppendLine(fmt.Sprintf("var %s_Lookup = [][]int64{", tokenName))
	next := generator.NewNode()
	next.AppendLine(fmt.Sprintf("var %s_Next = [][]int{", tokenName))
	imports := map[string]bool{"unicode/utf8": true}
	root := generator.NewNode()
	root.AppendLine(fmt.Sprintf("func %s(s string, offset int) int {", funcName))
	root.Indent(func(n generator.Node) {
		n.AppendLine("input := s[offset:]")
		n.AppendLine("length := len(input)")
		n.Append("accepted := map[int]bool{")
		acceptingStates := r.dfa.AcceptingStates
		stateIDs := make([]int, 0, len(acceptingStates))
		for state, isAccepting := range acceptingStates {
			if isAccepting {
				stateIDs = append(stateIDs, state)
			}
		}
		slices.Sort(stateIDs)
		for _, state := range stateIDs {
			n.Append(fmt.Sprintf("%d: true, ", state))
		}
		n.AppendLine("}")
		n.AppendLine(fmt.Sprintf("state := %d", r.dfa.StartState))
		n.AppendLine("acceptedIndex := 0")
		n.AppendLine("index := 0")
		n.AppendLine("loop: for index < length {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("r, runeSize := utf8.DecodeRuneInString(input[index:])")
			n.AppendLine("switch state {")
			transitions := r.dfa.TransitionsBySource

			for source := 0; source < r.dfa.StateCount; source++ {
				bySource := transitions[source]
				n.AppendLine(fmt.Sprintf("case %d:", source))
				n.Indent(func(n generator.Node) {
					result := GenerateTransitionsUsingBinarySearch(bySource, source, tokenName, imports)
					n.AppendNode(result.Code)
					lookup.AppendNode(result.Lookup)
					next.AppendNode(result.Next)
				})
			}
			n.AppendLine("default:")
			n.Indent(func(n generator.Node) {
				n.AppendLine("break loop")
			})
			n.AppendLine("}")
			n.AppendLine("index += runeSize")
			n.AppendLine("if accepted[state] {")
			n.Indent(func(n generator.Node) {
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
	return GenerateRegExpResult{
		Imports: imports,
		Code:    root,
		Lookup:  lookup,
		Next:    next,
	}
}
