package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type ATNBuilder interface {
	InitializeStartAndStopStates(rules []generated.ParserRule)

	AddBasicState(rule generated.ParserRule, production generated.Element) BasicState
	AddBasicBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState, loopback PlusLoopbackState, decision int) BasicBlockStartState
	AddBlockEndState(rule generated.ParserRule, production generated.Element, start BlockStartState) BlockEndState

	AddRuleStartState(rule generated.ParserRule, production generated.Element, stop RuleStopState) RuleStartState
	AddRuleStopState(rule generated.ParserRule, production generated.Element) RuleStopState
	AddStarBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState) StarBlockStartState
	AddPlusBlockStartState(rule generated.ParserRule, production generated.Element, decisionend BlockEndState, loopback PlusLoopbackState, decision int) PlusBlockStartState

	AddStarLoopbackState(rule generated.ParserRule, production generated.Element) StarLoopbackState
	AddPlusLoopbackState(rule generated.ParserRule, production generated.Element, decision int) PlusLoopbackState

	AddStarLoopEntryState(rule generated.ParserRule, production generated.Element, loopback StarLoopbackState, decision int) StarLoopEntryState
	AddLoopEndState(rule generated.ParserRule, production generated.Element, loopback ATNState) LoopEndState

	AddEpsilonTransition(from ATNState, to ATNState) EpsilonTransition
	AddAtomTransition(from ATNState, to ATNState, atom int) AtomTransition
	AddRuleTransition(from ATNState, to ATNState, rule generated.ParserRule, followState ATNState) RuleTransition

	Block(rule generated.ParserRule, elements []generated.Element) *ATNHandle
	Atom(rule generated.ParserRule, element generated.Element) *ATNHandle

	MakeBlock(alternatives []ATNHandle) *ATNHandle
	MakeRuleHandle(rule generated.ParserRule, block *ATNHandle) *ATNHandle

	Build() ATN
}

type ATNBuilderData struct {
	rules map[string]generated.ParserRule
	atn   ATN
}

func (a *ATNBuilderData) Plus(rule generated.ParserRule, plus generated.Element, handle *ATNHandle, sep *ATNHandle) *ATNHandle {
	return nil
}

func (a *ATNBuilderData) Star(rule generated.ParserRule, star generated.Element, handle *ATNHandle, sep *ATNHandle) *ATNHandle {
	return nil
}

func (a *ATNBuilderData) Optional(rule generated.ParserRule, optional generated.Element, handle *ATNHandle) *ATNHandle {
	return nil
}

func (a *ATNBuilderData) Atom(rule generated.ParserRule, element generated.Element) *ATNHandle {
	/**
	 * Element
	 * |- Action
	 * |- Group
	 * |- Assignment
	 * |- Assignable
	 *    |- Alternatives
	 *    |- RuleCall
	 *    |- CrossReference
	 *    |- Keyword
	 */
	switch casted := element.(type) {
	case generated.Action:
		return nil
	case generated.Group:
		return nil
	case generated.Assignment:
		return nil
	case generated.Alternatives:
		return nil
	case generated.RuleCall:
		nonTerminal := a.rules[casted.Rule()]
		start := a.AddBasicBlockStartState(rule, casted, nil, nil, 0)
		stop := a.AddBasicBlockStartState(rule, casted, nil, nil, 0)
		a.AddRuleTransition(start, stop, nonTerminal, stop)
		return &ATNHandle{
			Start: start,
			Stop:  stop,
		}
	case generated.CrossRef:
		start := a.AddBasicState(rule, element)
		stop := a.AddBasicState(rule, element)
		a.AddAtomTransition(start, stop, casted.TypeToken().TypeId)
		return &ATNHandle{
			Start: start,
			Stop:  stop,
		}
	case generated.Keyword:
		start := a.AddBasicState(rule, element)
		stop := a.AddBasicState(rule, element)
		a.AddAtomTransition(start, stop, casted.ValueToken().Type.Id)
		return &ATNHandle{
			Start: start,
			Stop:  stop,
		}
	default:
		return nil
	}
}

func (a *ATNBuilderData) Block(rule generated.ParserRule, elements []generated.Element) *ATNHandle {
	handles := make([]ATNHandle, 0)
	for _, elem := range elements {
		handle := a.Atom(rule, elem)
		if handle != nil {
			handles = append(handles, *handle)
		}
	}
	if len(handles) == 1 {
		return &handles[0]
	} else if len(handles) == 0 {
		return nil
	} else {
		return a.MakeBlock(handles)
	}
}

func (a *ATNBuilderData) MakeBlock(alternatives []ATNHandle) *ATNHandle {
	alternativesLength := len(alternatives)
	for i := 0; i < alternativesLength-1; i++ {
		handle := alternatives[i]
		var transition Transition
		if len(handle.Start.Transitions()) == 1 {
			transition = handle.Start.Transitions()[0]
		}
		ruleTransition, isRuleTransition := transition.(RuleTransition)
		next := alternatives[i+1].Start
		_, isStartBasic := handle.Start.(BasicState)
		_, isStopBasic := handle.Stop.(BasicState)
		if isStartBasic &&
			isStopBasic &&
			transition != nil &&
			((isRuleTransition && ruleTransition.FollowState() == handle.Stop) ||
				transition.Target() == handle.Stop) {
			if isRuleTransition {
				ruleTransition.SetFollowState(next)
			} else {
				transition.SetTarget(next)
			}
		} else {
			a.AddEpsilonTransition(handle.Stop, next)
		}
	}

	first := alternatives[0]
	last := alternatives[alternativesLength-1]
	return &ATNHandle{
		Start: first.Start,
		Stop:  last.Stop,
	}
}

func (a *ATNBuilderData) MakeRuleHandle(rule generated.ParserRule, block *ATNHandle) *ATNHandle {
	start := a.atn.RuleStartState(rule)
	stop := a.atn.RuleStopState(rule)
	a.AddEpsilonTransition(start, block.Start)
	a.AddEpsilonTransition(block.Stop, stop)
	handle := ATNHandle{
		Start: start,
		Stop:  stop,
	}
	return &handle
}

func NewATNBuilder() ATNBuilderData {
	return ATNBuilderData{
		atn: NewATN(),
	}
}

func (a *ATNBuilderData) InititializeStartAndStopStates(rules []generated.ParserRule) {
	for _, rule := range rules {
		stop := a.AddRuleStopState(rule, nil)
		start := a.AddRuleStartState(rule, nil, stop)
		a.atn.AddStartAndStopState(rule, start, stop)
	}
}

func (a *ATNBuilderData) AddBasicState(rule generated.ParserRule, production generated.Element) BasicState {
	state := NewBasicStateData(a.atn, production, rule, len(a.atn.States()))
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddBlockEndState(rule generated.ParserRule, production generated.Element, start BlockStartState) BlockEndState {
	state := NewBlockEndStateData(a.atn, production, rule, len(a.atn.States()), start)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddRuleStartState(rule generated.ParserRule, production generated.Element, stop RuleStopState) RuleStartState {
	state := NewRuleStartStateData(a.atn, production, rule, len(a.atn.States()), stop)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddRuleStopState(rule generated.ParserRule, production generated.Element) RuleStopState {
	state := NewRuleStopStateData(a.atn, production, rule, len(a.atn.States()))
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddStarBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState) StarBlockStartState {
	state := NewStarBlockStartStateData(a.atn, production, rule, len(a.atn.States()), end)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddPlusBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState, loopback PlusLoopbackState, decision int) PlusBlockStartState {
	state := NewPlusBlockStartStateData(a.atn, production, rule, len(a.atn.States()), end, loopback, decision)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddBasicBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState, loopback PlusLoopbackState, decision int) BasicBlockStartState {
	state := NewBasicBlockStartStateData(a.atn, production, rule, len(a.atn.States()), end, decision)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddStarLoopbackState(rule generated.ParserRule, production generated.Element) StarLoopbackState {
	state := NewStarLoopbackStateData(a.atn, production, rule, len(a.atn.States()))
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddPlusLoopbackState(rule generated.ParserRule, production generated.Element, decision int) PlusLoopbackState {
	state := NewPlusLoopbackStateData(a.atn, production, rule, len(a.atn.States()), decision)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddLoopEndState(rule generated.ParserRule, production generated.Element, loopback ATNState) LoopEndState {
	state := NewLoopEndStateData(a.atn, production, rule, len(a.atn.States()), loopback)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddStarLoopEntryState(rule generated.ParserRule, production generated.Element, loopback StarLoopbackState, decision int) StarLoopEntryState {
	state := NewStarLoopEntryStateData(a.atn, production, rule, len(a.atn.States()), loopback, decision)
	a.atn.AddState(state)
	return state
}

func (a *ATNBuilderData) AddEpsilonTransition(from ATNState, to ATNState) EpsilonTransition {
	transition := NewEpsilonTransitionData(to)
	from.AddTransition(transition)
	return transition
}

func (a *ATNBuilderData) AddAtomTransition(from ATNState, to ATNState, atom int) AtomTransition {
	transition := NewAtomTransitionData(to, atom)
	from.AddTransition(transition)
	return transition
}

func (a *ATNBuilderData) AddRuleTransition(from ATNState, to ATNState, rule generated.ParserRule, followState ATNState) RuleTransition {
	transition := NewRuleTransitionData(to, rule, followState)
	from.AddTransition(transition)
	return transition
}

func (a *ATNBuilderData) Build() ATN {
	return a.atn
}
