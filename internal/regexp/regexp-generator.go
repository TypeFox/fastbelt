package regexp

import (
	"fmt"
	"sort"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/automatons"
)

func GenerateTransitions_UsingBranching(bySource *automatons.RuneRangeTargetsMapping) generator.Node {
	n := generator.NewNode()
	targets := make(map[int]automatons.RuneSet)
	for transition := range bySource.All() {
		target := transition.Values[0]
		runeSet, exists := targets[target]
		if !exists {
			runeSet = *automatons.NewRuneSet_Empty()
		}
		runeSet.AddRange(transition.Range.Start, transition.Range.End)
		targets[target] = runeSet
	}
	for target, runeSet := range targets {
		n.Append("if ")
		first := true
		comment := ""
		for _, charRange := range runeSet.Ranges {
			if !charRange.Includes {
				continue
			}
			if !first {
				n.Append(" || ")
			} else {
				first = false
			}
			if charRange.Start == charRange.End {
				comment += fmt.Sprintf("%s, ", automatons.FormatRune(charRange.Start))
				n.Append(fmt.Sprintf("r == %s", automatons.FormatInt(charRange.Start)))
			} else {
				comment += fmt.Sprintf("%s..%s, ", automatons.FormatRune(charRange.Start), automatons.FormatRune(charRange.End))
				n.Append(fmt.Sprintf("r >= %s && r <= %s", automatons.FormatInt(charRange.Start), automatons.FormatInt(charRange.End)))
			}
		}
		n.AppendLine(fmt.Sprintf(" { // %s", comment))
		n.Indent(func(n generator.Node) {
			n.AppendLine(fmt.Sprintf("state = %d", target))
		})
		n.Append("} else ")
	}
	n.AppendLine("{")
	n.Indent(func(n generator.Node) {
		n.AppendLine("break loop")
	})
	n.AppendLine("}")
	return n
}

func GenerateTransitions_UsingBinarySearch(bySource *automatons.RuneRangeTargetsMapping) generator.Node {
	n := generator.NewNode()
	n.Append("lookup := []rune{")
	for transition := range bySource.All() {
		n.Append(fmt.Sprintf("%s, %s, ", automatons.FormatInt(transition.Range.Start), automatons.FormatInt(transition.Range.End)))
	}
	n.AppendLine("}")
	n.Append("next := []int{")
	for transition := range bySource.All() {
		n.Append(fmt.Sprintf("%d, ", transition.Values[0]))
	}
	n.AppendLine("}")
	n.AppendLine("nextState := regexp.BinarySearch_NextState(r, lookup, next)")
	n.AppendLine("if nextState > -1 {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("state = nextState")
	})
	n.AppendLine("} else {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("break loop")
	})
	n.AppendLine("}")
	return n
}

func (r *RegexpImpl) GenerateRegExp() generator.Node {
	root := generator.NewNode()
	root.AppendLine("func (s string, offset int) int {")
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
		n.AppendLine("acceptedIndex := -1")
		n.AppendLine("if accepted[state] {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("acceptedIndex = 0")
		})
		n.AppendLine("}")
		n.AppendLine("index := 0")
		n.AppendLine("loop: for index < length {")
		n.Indent(func(n generator.Node) {
			n.AppendLine("r := rune(input[index])")
			n.AppendLine("switch state {")
			transitions := r.dfa.TransitionsBySource

			sources := make([]int, 0, len(transitions))
			for k := range transitions {
				sources = append(sources, k)
			}
			sort.Ints(sources)

			for _, source := range sources {
				bySource := transitions[source]
				n.AppendLine(fmt.Sprintf("case %d:", source))
				n.Indent(func(n generator.Node) {
					if len(bySource.Ranges) >= 4 {
						n.AppendNode(GenerateTransitions_UsingBinarySearch(bySource))
					} else {
						n.AppendNode(GenerateTransitions_UsingBranching(bySource))
					}
				})
			}
			n.AppendLine("default:")
			n.Indent(func(n generator.Node) {
				n.AppendLine("break loop")
			})
			n.AppendLine("}")
			n.AppendLine("index++")
			n.AppendLine("if accepted[state] {")
			n.Indent(func(n generator.Node) {
				n.AppendLine("acceptedIndex = index")
			})
			n.AppendLine("}")
		})
		n.AppendLine("}")
		n.AppendLine("return acceptedIndex")
	})
	root.AppendLine("},")
	return root
}

func BinarySearch_NextState(r rune, lookup []rune, next []int) int {
	searchIndex := sort.Search(len(next), func(i int) bool {
		return lookup[i*2] > r
	}) - 1
	if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
		return next[searchIndex]
	}
	return -1
}
