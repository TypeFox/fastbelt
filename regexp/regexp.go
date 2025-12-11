package regexp

import (
	"fmt"
	"regexp/syntax"

	"abc.de/regex/automatons"
)

type Regexp interface {
	FindStringIndex(s string) (loc []int)
}

type regexpImpl struct {
	pattern string
	dfa     automatons.NFA
}

func (re regexpImpl) String() string {
	return re.dfa.(*automatons.NFAImpl).String()
}

func CompileRegexp(pattern string) (Regexp, error) {
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
	return &regexpImpl{
		pattern: pattern,
		dfa:     nfa,
	}, nil
}

func MustCompilRegexp(pattern string) Regexp {
	regexp, error := CompileRegexp(pattern)
	if error != nil {
		panic(error)
	}
	return regexp
}

func (r *regexpImpl) FindStringIndex(s string) (loc []int) {
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

func (r *regexpImpl) GenerateFindStringIndex(name string) Node {
	root := NewNode()
	root.AppendLine(fmt.Sprintf("func %s_FindStringIndex(input string) (loc []int) {", name))
	root.Indent(func(n Node) {
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
		n.Indent(func(n Node) {
			n.AppendLine("acceptedIndex = 0")
		})
		n.AppendLine("}")
		n.AppendLine("index := 0")
		n.AppendLine("halted := false")
		n.AppendLine("for !halted {")
		n.Indent(func(n Node) {
			n.AppendLine("if index >= length {")
			n.Indent(func(n Node) {
				n.AppendLine("halted = true")
				n.AppendLine("continue")
			})
			n.AppendLine("} else {")
			n.Indent(func(n Node) {
				n.AppendLine("r := rune(input[index])")
				n.AppendLine("switch state {")
				transitions := r.dfa.GetTransitionsBySource()
				for source, bySource := range transitions {
					n.AppendLine(fmt.Sprintf("case %d:", source))
					n.Indent(func(n Node) {
						for transition := range bySource.AllTransitions() {
							n.Append("if ")
							if transition.CharRange.Start == transition.CharRange.End {
								n.Append(fmt.Sprintf("r == '%c'", transition.CharRange.Start))
							} else {
								n.Append(fmt.Sprintf("r >= '%c' && r <= '%c'", transition.CharRange.Start, transition.CharRange.End))
							}
							n.AppendLine(" {")
							n.Indent(func(n Node) {
								n.AppendLine(fmt.Sprintf("state = %d", transition.Targets[0]))
							})
							n.Append("} else ")
						}
						n.AppendLine("{")
						n.Indent(func(n Node) {
							n.AppendLine("halted = true")
						})
						n.AppendLine("}")
					})
				}
				n.AppendLine("}")
			})
			n.AppendLine("}")
			n.AppendLine("if !halted && accepted[state] {")
			n.Indent(func(n Node) {
				n.AppendLine("acceptedIndex = index")
			})
			n.AppendLine("}")
			n.AppendLine("index++")
		})
		n.AppendLine("}")
		n.AppendLine("return []int{0, acceptedIndex}")
	})
	root.AppendLine("}")
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
