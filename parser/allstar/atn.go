// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import "fmt"

// ATNStateType is the discriminator for ATN state kinds.
type ATNStateType int

const (
	ATNInvalidType    ATNStateType = 0
	ATNBasic          ATNStateType = 1
	ATNRuleStart      ATNStateType = 2
	ATNPlusBlockStart ATNStateType = 4
	ATNStarBlockStart ATNStateType = 5
	ATNTokenStart     ATNStateType = 6
	ATNRuleStop       ATNStateType = 7
	ATNBlockEnd       ATNStateType = 8
	ATNStarLoopBack   ATNStateType = 9
	ATNStarLoopEntry  ATNStateType = 10
	ATNPlusLoopBack   ATNStateType = 11
	ATNLoopEnd        ATNStateType = 12
)

// ATNState is the single concrete ATN state type.
// Fields specific to certain state kinds are non-nil only for those kinds.
type ATNState struct {
	ATN                    *ATN
	Production             Production // nil for rule start/stop
	StateNumber            int
	Rule                   *Rule
	EpsilonOnlyTransitions bool
	Transitions            []Transition
	Type                   ATNStateType

	// Decision index; -1 if this state is not a decision state.
	Decision int

	// Populated for BlockStartState kinds.
	End      *ATNState
	Loopback *ATNState // PlusBlockStart.loopback, StarLoopEntry.loopback, LoopEnd.loopback

	// Populated for BlockEndState.
	Start *ATNState

	// Populated for RuleStartState.
	Stop *ATNState
}

// ATN is the Augmented Transition Network built from a set of parser rules.
type ATN struct {
	DecisionMap      map[string]*ATNState
	States           []*ATNState
	DecisionStates   []*ATNState
	RuleToStartState map[*Rule]*ATNState
	RuleToStopState  map[*Rule]*ATNState
}

// Transition is the interface implemented by all ATN transitions.
type Transition interface {
	Target() *ATNState
	IsEpsilon() bool
}

// AtomTransition fires on a specific token type.
// CategoryMatches holds the IDs of all token types that match via category
// inheritance; populated from the Terminal's CategoryMatches at ATN-build time.
type AtomTransition struct {
	target          *ATNState
	TokenTypeID     int
	CategoryMatches []int
}

func (t *AtomTransition) Target() *ATNState { return t.target }
func (t *AtomTransition) IsEpsilon() bool   { return false }

// EpsilonTransition fires without consuming a token.
type EpsilonTransition struct {
	target *ATNState
}

func (t *EpsilonTransition) Target() *ATNState { return t.target }
func (t *EpsilonTransition) IsEpsilon() bool   { return true }

// RuleTransition enters a sub-rule and returns to FollowState.
type RuleTransition struct {
	target      *ATNState // the rule's RuleStartState
	Rule        *Rule
	FollowState *ATNState
}

func (t *RuleTransition) Target() *ATNState { return t.target }
func (t *RuleTransition) IsEpsilon() bool   { return true }

// atnHandle is an internal pair of (entry, exit) ATN states for a sub-network.
type atnHandle struct {
	left  *ATNState
	right *ATNState
}

// BuildATNKey returns the decision-map key for a production occurrence in a rule.
// Format: "<ruleName>_<prodType>_<occurrence>"
func BuildATNKey(rule *Rule, prodType string, occurrence int) string {
	return fmt.Sprintf("%s_%s_%d", rule.Name, prodType, occurrence)
}

// CreateATN builds an ATN from the given parser rules.
func CreateATN(rules []*Rule) *ATN {
	atn := &ATN{
		DecisionMap:      map[string]*ATNState{},
		States:           []*ATNState{},
		DecisionStates:   []*ATNState{},
		RuleToStartState: map[*Rule]*ATNState{},
		RuleToStopState:  map[*Rule]*ATNState{},
	}
	createRuleStartAndStopATNStates(atn, rules)
	for _, rule := range rules {
		ruleBlock := block(atn, rule, rule.Definition)
		if ruleBlock == nil {
			continue
		}
		buildRuleHandle(atn, rule, ruleBlock)
	}
	return atn
}

func createRuleStartAndStopATNStates(atn *ATN, rules []*Rule) {
	for _, rule := range rules {
		start := newATNState(atn, rule, nil, ATNRuleStart)
		stop := newATNState(atn, rule, nil, ATNRuleStop)
		start.Stop = stop
		atn.RuleToStartState[rule] = start
		atn.RuleToStopState[rule] = stop
	}
}

func atom(atn *ATN, rule *Rule, prod Production) *atnHandle {
	switch prod.Kind() {
	case ProdTerminal:
		t := prod.(*Terminal)
		return tokenRef(atn, rule, t.TokenTypeID, t.CategoryMatches, prod)
	case ProdNonTerminal:
		return ruleRef(atn, rule, prod.(*NonTerminal))
	case ProdAlternation:
		return alternation(atn, rule, prod.(*Alternation))
	case ProdOption:
		return option(atn, rule, prod.(*Option))
	case ProdRepetition:
		return repetition(atn, rule, prod.(*Repetition))
	case ProdRepetitionMandatory:
		return repetitionMandatory(atn, rule, prod.(*RepetitionMandatory))
	case ProdAlternative:
		return block(atn, rule, prod.(*Alternative).Definition)
	default:
		return nil
	}
}

func repetition(atn *ATN, rule *Rule, rep *Repetition) *atnHandle {
	starState := newATNState(atn, rule, rep, ATNStarBlockStart)
	defineDecisionState(atn, starState)
	handle := makeAlts(atn, rule, starState, rep, block(atn, rule, rep.Definition))
	return star(atn, rule, rep, handle)
}

func repetitionMandatory(atn *ATN, rule *Rule, rep *RepetitionMandatory) *atnHandle {
	plusState := newATNState(atn, rule, rep, ATNPlusBlockStart)
	defineDecisionState(atn, plusState)
	handle := makeAlts(atn, rule, plusState, rep, block(atn, rule, rep.Definition))
	return plus(atn, rule, rep, handle)
}

func alternation(atn *ATN, rule *Rule, alt *Alternation) *atnHandle {
	start := newATNState(atn, rule, alt, ATNBasic)
	defineDecisionState(atn, start)
	alts := make([]*atnHandle, len(alt.Alternatives))
	for i, a := range alt.Alternatives {
		alts[i] = atom(atn, rule, a)
	}
	return makeAlts(atn, rule, start, alt, alts...)
}

func option(atn *ATN, rule *Rule, opt *Option) *atnHandle {
	start := newATNState(atn, rule, opt, ATNBasic)
	defineDecisionState(atn, start)
	handle := makeAlts(atn, rule, start, opt, block(atn, rule, opt.Definition))
	return optional(atn, rule, opt, handle)
}

func block(atn *ATN, rule *Rule, children []Production) *atnHandle {
	handles := make([]*atnHandle, 0, len(children))
	for _, child := range children {
		h := atom(atn, rule, child)
		if h != nil {
			handles = append(handles, h)
		}
	}
	switch len(handles) {
	case 0:
		return nil
	case 1:
		return handles[0]
	default:
		return makeBlock(atn, handles)
	}
}

func plus(atn *ATN, rule *Rule, prod Production, handle *atnHandle) *atnHandle {
	blkStart := handle.left
	blkEnd := handle.right

	loop := newATNState(atn, rule, prod, ATNPlusLoopBack)
	defineDecisionState(atn, loop)
	end := newATNState(atn, rule, prod, ATNLoopEnd)
	blkStart.Loopback = loop
	end.Loopback = loop

	typeName, _ := ProductionTypeName(prod)
	atn.DecisionMap[BuildATNKey(rule, typeName, prod.Occurrence())] = loop

	addEpsilon(blkEnd, loop)  // block can see loop back
	addEpsilon(loop, blkStart) // loop back to start
	addEpsilon(loop, end)      // exit

	return &atnHandle{left: blkStart, right: end}
}

func star(atn *ATN, rule *Rule, prod Production, handle *atnHandle) *atnHandle {
	start := handle.left
	end := handle.right

	entry := newATNState(atn, rule, prod, ATNStarLoopEntry)
	defineDecisionState(atn, entry)
	loopEnd := newATNState(atn, rule, prod, ATNLoopEnd)
	loop := newATNState(atn, rule, prod, ATNStarLoopBack)
	entry.Loopback = loop
	loopEnd.Loopback = loop

	addEpsilon(entry, start)    // loop enter edge (alt 0)
	addEpsilon(entry, loopEnd)  // bypass loop edge (alt 1)
	addEpsilon(end, loop)       // block end hits loop back
	addEpsilon(loop, entry)     // loop back to entry/exit decision

	typeName, _ := ProductionTypeName(prod)
	atn.DecisionMap[BuildATNKey(rule, typeName, prod.Occurrence())] = entry

	return &atnHandle{left: entry, right: loopEnd}
}

func optional(atn *ATN, rule *Rule, opt *Option, handle *atnHandle) *atnHandle {
	start := handle.left
	end := handle.right

	addEpsilon(start, end)

	atn.DecisionMap[BuildATNKey(rule, "Option", opt.Idx)] = start
	return handle
}

func makeAlts(atn *ATN, rule *Rule, start *ATNState, prod Production, alts ...*atnHandle) *atnHandle {
	end := newATNState(atn, rule, prod, ATNBlockEnd)
	end.Start = start
	start.End = end

	for _, alt := range alts {
		if alt != nil {
			addEpsilon(start, alt.left)
			addEpsilon(alt.right, end)
		} else {
			addEpsilon(start, end)
		}
	}

	typeName, ok := ProductionTypeName(prod)
	if ok {
		atn.DecisionMap[BuildATNKey(rule, typeName, prod.Occurrence())] = start
	}

	return &atnHandle{left: start, right: end}
}

func makeBlock(atn *ATN, alts []*atnHandle) *atnHandle {
	for i := 0; i < len(alts)-1; i++ {
		handle := alts[i]
		next := alts[i+1].left
		var t Transition
		if len(handle.left.Transitions) == 1 {
			t = handle.left.Transitions[0]
		}
		rt, isRule := t.(*RuleTransition)
		if handle.left.Type == ATNBasic &&
			handle.right.Type == ATNBasic &&
			t != nil &&
			((isRule && rt.FollowState == handle.right) ||
				(!isRule && t.Target() == handle.right)) {
			// avoid epsilon edge to next element
			if isRule {
				rt.FollowState = next
			} else {
				switch at := t.(type) {
				case *AtomTransition:
					at.target = next
				case *EpsilonTransition:
					at.target = next
				}
			}
			removeState(atn, handle.right) // we skipped over this state
		} else {
			addEpsilon(handle.right, next)
		}
	}

	first := alts[0]
	last := alts[len(alts)-1]
	return &atnHandle{left: first.left, right: last.right}
}

func tokenRef(atn *ATN, rule *Rule, tokenTypeID int, categoryMatches []int, prod Production) *atnHandle {
	left := newATNState(atn, rule, prod, ATNBasic)
	right := newATNState(atn, rule, prod, ATNBasic)
	addTransition(left, &AtomTransition{
		target:          right,
		TokenTypeID:     tokenTypeID,
		CategoryMatches: categoryMatches,
	})
	return &atnHandle{left: left, right: right}
}

func ruleRef(atn *ATN, currentRule *Rule, nt *NonTerminal) *atnHandle {
	ruleStart := atn.RuleToStartState[nt.ReferencedRule]
	left := newATNState(atn, currentRule, nt, ATNBasic)
	right := newATNState(atn, currentRule, nt, ATNBasic)
	addTransition(left, &RuleTransition{
		target:      ruleStart,
		Rule:        nt.ReferencedRule,
		FollowState: right,
	})
	return &atnHandle{left: left, right: right}
}

func buildRuleHandle(atn *ATN, rule *Rule, b *atnHandle) {
	start := atn.RuleToStartState[rule]
	addEpsilon(start, b.left)
	stop := atn.RuleToStopState[rule]
	addEpsilon(b.right, stop)
}

func addEpsilon(from, to *ATNState) {
	addTransition(from, &EpsilonTransition{target: to})
}

func newATNState(atn *ATN, rule *Rule, prod Production, typ ATNStateType) *ATNState {
	s := &ATNState{
		ATN:         atn,
		Production:  prod,
		StateNumber: len(atn.States),
		Rule:        rule,
		Type:        typ,
		Decision:    -1,
	}
	atn.States = append(atn.States, s)
	return s
}

func addTransition(state *ATNState, t Transition) {
	if len(state.Transitions) == 0 {
		state.EpsilonOnlyTransitions = t.IsEpsilon()
	}
	state.Transitions = append(state.Transitions, t)
}

func removeState(atn *ATN, state *ATNState) {
	for i, s := range atn.States {
		if s == state {
			atn.States = append(atn.States[:i], atn.States[i+1:]...)
			return
		}
	}
}

func defineDecisionState(atn *ATN, state *ATNState) int {
	atn.DecisionStates = append(atn.DecisionStates, state)
	state.Decision = len(atn.DecisionStates) - 1
	return state.Decision
}
