package regexp

import (
	"fmt"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/automatons"
)

type GenerateTransitionsUsingBinarySearchResult struct {
	Lookup generator.Node
	Next   generator.Node
	Code   generator.Node
}

func GenerateTransitionsUsingBinarySearch(bySource *automatons.RuneRangeTargetsMapping, source int, tokenName string) GenerateTransitionsUsingBinarySearchResult {
	lookup := generator.NewNode()
	lookup.Append("{")
	if bySource != nil {
		for transition := range bySource.All() {
			lookup.Append(fmt.Sprintf("%s, %s, ", automatons.FormatInt(transition.Range.Start), automatons.FormatInt(transition.Range.End)))
		}
	}
	lookup.AppendLine("},")

	next := generator.NewNode()
	next.Append("{")
	if bySource != nil {
		for transition := range bySource.All() {
			next.Append(fmt.Sprintf("%d, ", transition.Values[0]))
		}
	}
	next.AppendLine("},")

	n := generator.NewNode()
	n.AppendLine("nextState := -1")
	n.AppendLine(fmt.Sprintf("next := %s_Next[%d]", tokenName, source))
	n.AppendLine(fmt.Sprintf("lookup := %s_Lookup[%d]", tokenName, source))
	n.AppendLine("searchIndex := sort.Search(len(next), func(i int) bool {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("return lookup[i*2] > r")
	})
	n.AppendLine("}) - 1")
	n.AppendLine("if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("nextState = next[searchIndex]")
	})
	n.AppendLine("}")
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
	lookup.AppendLine(fmt.Sprintf("var %s_Lookup = [][]rune{", tokenName))
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
		for state, isAccepting := range acceptingStates {
			if isAccepting {
				n.Append(fmt.Sprintf("%d: true, ", state))
			}
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
					result := GenerateTransitionsUsingBinarySearch(bySource, source, tokenName)
					n.AppendNode(result.Code)
					lookup.AppendNode(result.Lookup)
					next.AppendNode(result.Next)
					imports["sort"] = true
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
