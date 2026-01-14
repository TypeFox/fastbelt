package regexp

import (
	"fmt"
	"regexp/syntax"

	"typefox.dev/fastbelt/internal/automatons"
)

type Regexp interface {
	FindStringIndex(s string) (loc []int)
}

type RegexpImpl struct {
	pattern string
	dfa     *automatons.NFA
}

func (re RegexpImpl) String() string {
	return re.dfa.String()
}

func Compile(pattern string) (Regexp, error) {
	op, error := syntax.Parse(pattern, syntax.Perl)
	if error != nil {
		return nil, error
	}
	nfa := newNFAFromSyntax(op)
	nfa = nfa.Determinize()
	nfa = nfa.Minimize()
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
	state := r.dfa.InitializeReducerState(s)

	for !state.Halted {
		nextState, err := r.dfa.Step(state)
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
	startCharsSet := automatons.NewRuneSetEmpty()
	startState := r.dfa.StartState
	transitions := r.dfa.TransitionsBySource[startState]
	if transitions == nil {
		return startCharsSet
	}
	for transition := range transitions.All() {
		if !transition.Range.Includes {
			continue
		}
		//only save lowest byte
		modStart := transition.Range.Start % 0x100
		modEnd := transition.Range.End % 0x100
		if transition.Range.End-transition.Range.Start+1 < 0x100 {
			if modStart <= modEnd {
				startCharsSet.AddRange(modStart, modEnd)
			} else {
				startCharsSet.AddRange(modStart, 0xFF)
				startCharsSet.AddRange(0x00, modEnd)
			}
		} else {
			startCharsSet.AddRange(0x00, 0xFF)
		}
	}
	return startCharsSet
}

func newNFAFromSyntax(op *syntax.Regexp) *automatons.NFA {
	kit := automatons.NewConstructionKit()
	switch op.Op {
	case syntax.OpLiteral:
		chain := make([]*automatons.NFA, len(op.Rune))
		for i, r := range op.Rune {
			chain[i] = kit.Consume(automatons.NewRuneSetRune(r))
		}
		return kit.Concat(chain...)
	case syntax.OpCharClass:
		runeSet := automatons.NewRuneSetEmpty()
		for i := 0; i < len(op.Rune); i += 2 {
			start := op.Rune[i]
			end := op.Rune[i+1]
			runeSet.AddRange(start, end)
		}
		return kit.Consume(runeSet)
	case syntax.OpAnyChar:
		runeSet := automatons.NewRuneSetFull()
		return kit.Consume(runeSet)
	case syntax.OpConcat:
		chain := make([]*automatons.NFA, len(op.Sub))
		for i, sub := range op.Sub {
			chain[i] = newNFAFromSyntax(sub)
		}
		return kit.Concat(chain...)
	case syntax.OpAlternate:
		alternatives := make([]*automatons.NFA, len(op.Sub))
		for i, sub := range op.Sub {
			alternatives[i] = newNFAFromSyntax(sub)
		}
		return kit.Alternate(alternatives...)
	case syntax.OpStar:
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, 0, -1)
	case syntax.OpPlus:
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, 1, -1)
	case syntax.OpQuest:
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, 0, 1)
	case syntax.OpRepeat:
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, int(op.Min), int(op.Max))
	case syntax.OpCapture:
		return newNFAFromSyntax(op.Sub[0])
	case syntax.OpAnyCharNotNL:
		runeSet := automatons.NewRuneSetFull()
		runeSet.RemoveRune('\n')
		return kit.Consume(runeSet)
	case syntax.OpBeginLine:
		fallthrough
	case syntax.OpEndLine:
		fallthrough
	case syntax.OpBeginText:
		fallthrough
	case syntax.OpEndText:
		fallthrough
	case syntax.OpWordBoundary:
		fallthrough
	case syntax.OpNoWordBoundary:
		panic(fmt.Sprintf("\\B, \\b, ^, $ are not supported yet: %v", op.Op))
	case syntax.OpEmptyMatch:
		return kit.Empty()
	case syntax.OpNoMatch:
		return kit.Reject()
	default:
		panic(fmt.Sprintf("unsupported syntax operation: %v", op.Op))
	}
}
