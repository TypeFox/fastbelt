package regexp

import (
	"fmt"
	"regexp/syntax"
	"sort"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/automatons"
)

type Regexp interface {
	FindStringIndex(s string) (loc []int)
}

type RegexpImpl struct {
	pattern string
	dfa     automatons.NFA
}

func (re RegexpImpl) String() string {
	return re.dfa.(*automatons.NFAImpl).String()
}

func Compile(pattern string) (Regexp, error) {
	op, error := syntax.Parse(pattern, syntax.Perl)
	if error != nil {
		return nil, error
	}
	nfa, err := newNFAFromSyntax(op)
	nfa = automatons.Determinize(nfa)
	nfa = automatons.Minimize(nfa)
	if err != nil {
		return nil, err
	}
	return &RegexpImpl{
		pattern: pattern,
		dfa:     nfa,
	}, nil
}

func MustCompile(pattern string) Regexp {
	regexp, error := Compile(pattern)
	if error != nil {
		panic(error)
	}
	return regexp
}

func (r *RegexpImpl) FindStringIndex(s string) (loc []int) {
	dfa := r.dfa.(*automatons.NFAImpl)
	state := dfa.InitializeReducerState(s)

	for !state.Halted {
		nextState, err := dfa.Step(state)
		if err != nil {
			return nil
		}
		state = nextState
	}

	if state.AcceptedIdx != -1 {
		return []int{0, state.AcceptedIdx}
	}
	return nil
}

func (r *RegexpImpl) GetStartChars() *automatons.RuneSet {
	startCharsSet := automatons.NewRuneSet_Empty()
	transitions := r.dfa.GetTransitionsBySource()
	startState := r.dfa.GetStartState()
	for transition := range transitions[startState].AllTransitions() {
		if !transition.CharRange.Includes {
			continue
		}
		startCharsSet.AddRange(transition.CharRange.Start, transition.CharRange.End)
	}
	return startCharsSet
}

func (r *RegexpImpl) GenerateLambda() generator.Node {
	root := generator.NewNode()
	root.AppendLine("func (s string, offset int) int {")
	root.Indent(func(n generator.Node) {
		n.AppendLine("input := s[offset:]")
		n.AppendLine("length := len(input)")
		n.Append("accepted := map[int]bool{")
		acceptingStates := r.dfa.GetAcceptingStates()
		for state, isAccepting := range acceptingStates {
			if isAccepting {
				n.Append(fmt.Sprintf("%d: true, ", state))
			}
		}
		n.AppendLine("}")
		n.AppendLine(fmt.Sprintf("state := %d", r.dfa.GetStartState()))
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
			transitions := r.dfa.GetTransitionsBySource()

			sources := make([]int, 0, len(transitions))
			for k := range transitions {
				sources = append(sources, k)
			}
			sort.Ints(sources)

			for _, source := range sources {
				bySource := transitions[source]
				n.AppendLine(fmt.Sprintf("case %d:", source))
				n.Indent(func(n generator.Node) {
					targets := make(map[int]automatons.RuneSet)
					for transition := range bySource.AllTransitions() {
						target := transition.Targets[0]
						runeSet, exists := targets[target]
						if !exists {
							runeSet = *automatons.NewRuneSet_Empty()
						}
						runeSet.AddRange(transition.CharRange.Start, transition.CharRange.End)
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
								n.Append(fmt.Sprintf("r == %d", charRange.Start))
							} else {
								comment += fmt.Sprintf("%s..%s, ", automatons.FormatRune(charRange.Start), automatons.FormatRune(charRange.End))
								n.Append(fmt.Sprintf("r >= %d && r <= %d", charRange.Start, charRange.End))
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
				})
			}
			n.AppendLine("}")
			n.AppendLine("if accepted[state] {")
			n.Indent(func(n generator.Node) {
				n.AppendLine("acceptedIndex = index")
			})
			n.AppendLine("}")
			n.AppendLine("index++")
		})
		n.AppendLine("}")
		n.AppendLine("return acceptedIndex")
	})
	root.AppendLine("},")
	return root
}

func newNFAFromSyntax(op *syntax.Regexp) (automatons.NFA, error) {
	kit := automatons.NewConstructionKit()
	switch op.Op {
	case syntax.OpLiteral:
		chain := make([]automatons.NFA, len(op.Rune))
		for i, r := range op.Rune {
			nfa, error := kit.Consume(automatons.NewRuneSet_Rune(r))
			if error != nil {
				return nil, error
			}
			chain[i] = nfa
		}
		return kit.Concat(chain...)
	case syntax.OpCharClass:
		runeSet := automatons.NewRuneSet_Empty()
		for i := 0; i < len(op.Rune); i += 2 {
			start := op.Rune[i]
			end := op.Rune[i+1]
			runeSet.AddRange(start, end)
		}
		return kit.Consume(runeSet)
	case syntax.OpAnyChar:
		runeSet := automatons.NewRuneSet_Full()
		return kit.Consume(runeSet)
	case syntax.OpConcat:
		chain := make([]automatons.NFA, len(op.Sub))
		for i, sub := range op.Sub {
			nfa, error := newNFAFromSyntax(sub)
			if error != nil {
				return nil, error
			}
			chain[i] = nfa
		}
		return kit.Concat(chain...)
	case syntax.OpAlternate:
		alternatives := make([]automatons.NFA, len(op.Sub))
		for i, sub := range op.Sub {
			nfa, error := newNFAFromSyntax(sub)
			if error != nil {
				return nil, error
			}
			alternatives[i] = nfa
		}
		return kit.Alternate(alternatives...)
	case syntax.OpStar:
		nfa, error := newNFAFromSyntax(op.Sub[0])
		if error != nil {
			return nil, error
		}
		return kit.Repeat(nfa, 0, -1)
	case syntax.OpPlus:
		nfa, error := newNFAFromSyntax(op.Sub[0])
		if error != nil {
			return nil, error
		}
		return kit.Repeat(nfa, 1, -1)
	case syntax.OpQuest:
		nfa, error := newNFAFromSyntax(op.Sub[0])
		if error != nil {
			return nil, error
		}
		return kit.Repeat(nfa, 0, 1)
	case syntax.OpRepeat:
		nfa, error := newNFAFromSyntax(op.Sub[0])
		if error != nil {
			return nil, error
		}
		return kit.Repeat(nfa, int(op.Min), int(op.Max))
	case syntax.OpCapture:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpAnyCharNotNL:
		runeSet := automatons.NewRuneSet_Full()
		runeSet.RemoveRune('\n')
		return kit.Consume(runeSet)
	case syntax.OpBeginLine:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpEndLine:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpBeginText:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpEndText:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpWordBoundary:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpNoWordBoundary:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpEmptyMatch:
		return kit.Empty()
	case syntax.OpNoMatch:
		return kit.Reject()
	default:
		return nil, fmt.Errorf("unsupported syntax operation: %v", op.Op)
	}
}
