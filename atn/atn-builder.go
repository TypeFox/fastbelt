package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type ATNBuilder interface {
	InititializeStartAndStopStates(rules []generated.ParserRule)

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

	MakeBlock(rule generated.ParserRule, element []generated.Element) ATNHandle
	MakeRuleHandle(rule generated.ParserRule, block ATNHandle) ATNHandle

	Build() ATN
}

type ATNBuilderData struct {
	atn ATN
}

func (b *ATNBuilderData) MakeBlock(rule generated.ParserRule, element []generated.Element) *ATNHandle {
	return nil
}

func (b *ATNBuilderData) MakeRuleHandle(rule generated.ParserRule, block *ATNHandle) *ATNHandle {
	start := b.atn.RuleStartState(rule)
	stop := b.atn.RuleStopState(rule)
	b.AddEpsilonTransition(start, block.Start)
	b.AddEpsilonTransition(block.Stop, stop)
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

func (b *ATNBuilderData) InititializeStartAndStopStates(rules []generated.ParserRule) {
	for _, rule := range rules {
		stop := b.AddRuleStopState(rule, nil)
		start := b.AddRuleStartState(rule, nil, stop)
		b.atn.AddStartAndStopState(rule, start, stop)
	}
}

func (b *ATNBuilderData) AddBasicState(rule generated.ParserRule, production generated.Element) BasicState {
	state := NewBasicStateData(b.atn, production, rule, len(b.atn.States()))
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddBlockEndState(rule generated.ParserRule, production generated.Element, start BlockStartState) BlockEndState {
	state := NewBlockEndStateData(b.atn, production, rule, len(b.atn.States()), start)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddRuleStartState(rule generated.ParserRule, production generated.Element, stop RuleStopState) RuleStartState {
	state := NewRuleStartStateData(b.atn, production, rule, len(b.atn.States()), stop)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddRuleStopState(rule generated.ParserRule, production generated.Element) RuleStopState {
	state := NewRuleStopStateData(b.atn, production, rule, len(b.atn.States()))
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddStarBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState) StarBlockStartState {
	state := NewStarBlockStartStateData(b.atn, production, rule, len(b.atn.States()), end)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddPlusBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState, loopback PlusLoopbackState, decision int) PlusBlockStartState {
	state := NewPlusBlockStartStateData(b.atn, production, rule, len(b.atn.States()), end, loopback, decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddBasicBlockStartState(rule generated.ParserRule, production generated.Element, end BlockEndState, loopback PlusLoopbackState, decision int) BasicBlockStartState {
	state := NewBasicBlockStartStateData(b.atn, production, rule, len(b.atn.States()), end, decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddStarLoopbackState(rule generated.ParserRule, production generated.Element) StarLoopbackState {
	state := NewStarLoopbackStateData(b.atn, production, rule, len(b.atn.States()))
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddPlusLoopbackState(rule generated.ParserRule, production generated.Element, decision int) PlusLoopbackState {
	state := NewPlusLoopbackStateData(b.atn, production, rule, len(b.atn.States()), decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddLoopEndState(rule generated.ParserRule, production generated.Element, loopback ATNState) LoopEndState {
	state := NewLoopEndStateData(b.atn, production, rule, len(b.atn.States()), loopback)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddStarLoopEntryState(rule generated.ParserRule, production generated.Element, loopback StarLoopbackState, decision int) StarLoopEntryState {
	state := NewStarLoopEntryStateData(b.atn, production, rule, len(b.atn.States()), loopback, decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddEpsilonTransition(from ATNState, to ATNState) EpsilonTransition {
	transition := NewEpsilonTransitionData(to)
	from.AddTransition(transition)
	return transition
}

func (b *ATNBuilderData) AddAtomTransition(from ATNState, to ATNState, atom int) AtomTransition {
	transition := NewAtomTransitionData(to, atom)
	from.AddTransition(transition)
	return transition
}

func (b *ATNBuilderData) AddRuleTransition(from ATNState, to ATNState, rule generated.ParserRule, followState ATNState) RuleTransition {
	transition := NewRuleTransitionData(to, rule, followState)
	from.AddTransition(transition)
	return transition
}

func (b *ATNBuilderData) Build() ATN {
	return b.atn
}
