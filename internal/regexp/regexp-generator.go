package regexp

import (
	"fmt"
	"sort"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/automatons"
)

func GenerateTransitionsUsingBranching(bySource *automatons.RuneRangeTargetsMapping) generator.Node {
	n := generator.NewNode()
	targets := make(map[int]automatons.RuneSet)
	keys := make([]int, 0)
	for transition := range bySource.All() {
		target := transition.Values[0]
		runeSet, exists := targets[target]
		if !exists {
			runeSet = *automatons.NewRuneSetEmpty()
			keys = append(keys, target)
		}
		runeSet.AddRange(transition.Range.Start, transition.Range.End)
		targets[target] = runeSet
	}
	sort.Ints(keys)
	for _, target := range keys {
		runeSet := targets[target]
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

func GenerateTransitionsUsingSwitchCasing(bySource *automatons.RuneRangeTargetsMapping) generator.Node {
	n := generator.NewNode()
	n.AppendLine("switch r {")
	var last *automatons.RuneRangeMappingSection[automatons.Targets] = nil
	for transition := range bySource.All() {
		if last != nil {
			if last.Values[0] != transition.Values[0] {
				n.Indent(func(n generator.Node) {
					n.AppendLine(fmt.Sprintf("state = %d", last.Values[0]))
				})
			} else {
				n.Indent(func(n generator.Node) {
					n.AppendLine(" fallthrough")
				})
			}
		}
		for char := transition.Range.Start; char <= transition.Range.End; char++ {
			n.AppendLine(fmt.Sprintf("case %s:", automatons.FormatInt(char)))
			if char != transition.Range.End {
				n.Indent(func(n generator.Node) {
					n.AppendLine(" fallthrough")
				})
			}
		}
		last = &transition
	}
	n.Indent(func(n generator.Node) {
		n.AppendLine(fmt.Sprintf("state = %d", last.Values[0]))
	})
	n.AppendLine("default:")
	n.Indent(func(n generator.Node) {
		n.AppendLine("break loop")
	})
	n.AppendLine("}")
	return n
}

func GenerateTransitionsUsingBinarySearch(bySource *automatons.RuneRangeTargetsMapping) generator.Node {
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
	n.AppendLine("nextState := -1")
	n.AppendLine("searchIndex := sort.Search(len(next), func(i int) bool {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("return lookup[i*2] > r")
	})
	n.AppendLine("}) - 1")
	n.AppendLine("if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {")
	n.Indent(func(n generator.Node) {
		n.AppendLine("return next[searchIndex]")
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
			n.AppendLine("r, runeSize := utf8.DecodeRuneInString(input[index:])")
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
					if len(bySource.Ranges) < 4 {
						n.AppendNode(GenerateTransitionsUsingBranching(bySource))
					} else if automatons.GetNumberOfRunes(bySource.Ranges) < 50 {
						n.AppendNode(GenerateTransitionsUsingSwitchCasing(bySource))
					} else {
						n.AppendNode(GenerateTransitionsUsingBinarySearch(bySource))
					}
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
	root.AppendLine("},")
	return root
}
