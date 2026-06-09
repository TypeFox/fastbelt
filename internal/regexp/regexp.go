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
	reachabilities := r.dfa.ComputeAcceptanceReachability()
	transitions := r.dfa.TransitionsBySource[startState]
	if transitions == nil {
		return startCharsSet
	}
	for transition := range transitions.All() {
		if !transition.Range.Includes {
			continue
		}
		acceptedStateReachable := false
		for _, target := range transition.Values {
			if !reachabilities[target] {
				continue
			}
			acceptedStateReachable = true
			break
		}
		if !acceptedStateReachable {
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
		return buildConcat(kit, op.Sub)
	case syntax.OpAlternate:
		alternatives := make([]*automatons.NFA, len(op.Sub))
		for i, sub := range op.Sub {
			alternatives[i] = newNFAFromSyntax(sub)
		}
		return kit.Alternate(alternatives...)
	case syntax.OpStar:
		if isNonGreedy(op) {
			panic(unsupportedLazyMsg(op))
		}
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, 0, -1)
	case syntax.OpPlus:
		if isNonGreedy(op) {
			panic(unsupportedLazyMsg(op))
		}
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, 1, -1)
	case syntax.OpQuest:
		if isNonGreedy(op) {
			panic(unsupportedLazyMsg(op))
		}
		nfa := newNFAFromSyntax(op.Sub[0])
		return kit.Repeat(nfa, 0, 1)
	case syntax.OpRepeat:
		if isNonGreedy(op) {
			panic(unsupportedLazyMsg(op))
		}
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

// isNonGreedy reports whether op carries the syntax.NonGreedy flag.
func isNonGreedy(op *syntax.Regexp) bool {
	return op.Flags&syntax.NonGreedy != 0
}

// isLazyQuantifier reports whether op is a quantifier with the non-greedy flag set.
func isLazyQuantifier(op *syntax.Regexp) bool {
	switch op.Op {
	case syntax.OpStar, syntax.OpPlus, syntax.OpQuest, syntax.OpRepeat:
		return isNonGreedy(op)
	}
	return false
}

func unsupportedLazyMsg(op *syntax.Regexp) string {
	return fmt.Sprintf("non-greedy quantifier %v is only supported when followed by a fixed tail in a concatenation", op)
}

// quantifierBounds extracts the (min, max) bounds of a quantifier op.
// max == -1 means unbounded.
func quantifierBounds(op *syntax.Regexp) (int, int) {
	switch op.Op {
	case syntax.OpStar:
		return 0, -1
	case syntax.OpPlus:
		return 1, -1
	case syntax.OpQuest:
		return 0, 1
	case syntax.OpRepeat:
		return int(op.Min), int(op.Max)
	}
	panic(fmt.Sprintf("not a quantifier: %v", op.Op))
}

// buildConcat builds the NFA for an OpConcat's children. When a child is a
// non-greedy (lazy) quantifier, it delegates to ConstructionKit.LazyMatch to
// produce the shortest-match construction. Recurses on the tail so multiple
// lazy quantifiers compose.
func buildConcat(kit *automatons.ConstructionKit, subs []*syntax.Regexp) *automatons.NFA {
	for i, sub := range subs {
		if !isLazyQuantifier(sub) {
			// Skip non-lazy expressions in the regexp
			// If there are non, simply build the concatenation as usual
			continue
		}
		// Recurse into buildConcat to handle multiple lazy quantifiers
		prefixNFA := buildConcat(kit, subs[:i])
		tailNFA := buildConcat(kit, subs[i+1:])
		bodyNFA := newNFAFromSyntax(sub.Sub[0])
		minRep, maxRep := quantifierBounds(sub)
		return kit.Concat(prefixNFA, kit.LazyMatch(bodyNFA, minRep, maxRep, tailNFA))
	}

	// Simplifications for common cases
	if len(subs) == 0 {
		return kit.Empty()
	}
	if len(subs) == 1 {
		return newNFAFromSyntax(subs[0])
	}
	// Simply concatenate all subexpressions
	// This is the default case if there are no lazy quantifiers
	chain := make([]*automatons.NFA, len(subs))
	for i, sub := range subs {
		chain[i] = newNFAFromSyntax(sub)
	}
	return kit.Concat(chain...)
}
