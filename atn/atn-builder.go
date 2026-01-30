package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type ATNBuilder interface {
	InititializeStartAndStopStates(rules []generated.ParserRule)

	AddBasicState(production *generated.Element, rule *generated.ParserRule) BasicState
	AddBasicBlockStartState(production *generated.Element, rule *generated.ParserRule, end *BlockEndState, loopback PlusLoopbackState, decision int) BasicBlockStartState
	AddBlockEndState(production *generated.Element, rule *generated.ParserRule, start *BlockStartState) BlockEndState

	AddRuleStartState(production *generated.Element, rule *generated.ParserRule, stop *RuleStopState) RuleStartState
	AddRuleStopState(production *generated.Element, rule *generated.ParserRule) RuleStopState

	AddStarBlockStartState(production *generated.Element, rule *generated.ParserRule, end *BlockEndState) StarBlockStartState
	AddPlusBlockStartState(production *generated.Element, rule *generated.ParserRule, decisionend *BlockEndState, loopback PlusLoopbackState, decision int) PlusBlockStartState

	AddStarLoopbackState(production *generated.Element, rule *generated.ParserRule) StarLoopbackState
	AddPlusLoopbackState(production *generated.Element, rule *generated.ParserRule, decision int) PlusLoopbackState

	AddStarLoopEntryState(production *generated.Element, rule *generated.ParserRule, loopback *StarLoopbackState, decision int) StarLoopEntryState
	AddLoopEndState(production *generated.Element, rule *generated.ParserRule, loopback *ATNState) LoopEndState

	AddEpsilonTransition(from ATNState, to ATNState) EpsilonTransition
	AddAtomTransition(from ATNState, to ATNState, atom int) AtomTransition
	AddRuleTransition(from ATNState, to ATNState, rule *generated.ParserRule, followState *ATNState) RuleTransition

	Build() *ATN
}

type ATNBuilderData struct {
	atn ATN
}

func NewATNBuilder() ATNBuilderData {
	return ATNBuilderData{
		atn: NewATN(),
	}
}

func (b *ATNBuilderData) InititializeStartAndStopStates(rules []generated.ParserRule) {
	for _, rule := range rules {
		stop := b.AddRuleStopState(&rule, nil)
		start := b.AddRuleStartState(&rule, nil, &stop)
		b.atn.AddStartAndStopState(&rule, &start, &stop)
	}
}

func (b *ATNBuilderData) AddBasicState(rule *generated.ParserRule, production *generated.Element) BasicState {
	state := NewBasicStateData(&b.atn, production, rule, len(b.atn.States()))
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddBlockEndState(rule *generated.ParserRule, production *generated.Element, start *BlockStartState) BlockEndState {
	state := NewBlockEndStateData(&b.atn, production, rule, len(b.atn.States()), start)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddRuleStartState(rule *generated.ParserRule, production *generated.Element, stop *RuleStopState) RuleStartState {
	state := NewRuleStartStateData(&b.atn, production, rule, len(b.atn.States()), stop)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddRuleStopState(rule *generated.ParserRule, production *generated.Element) RuleStopState {
	state := NewRuleStopStateData(&b.atn, production, rule, len(b.atn.States()))
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddStarBlockStartState(rule *generated.ParserRule, production *generated.Element, end *BlockEndState) StarBlockStartState {
	state := NewStarBlockStartStateData(&b.atn, production, rule, len(b.atn.States()), end)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddPlusBlockStartState(rule *generated.ParserRule, production *generated.Element, end *BlockEndState, loopback PlusLoopbackState, decision int) PlusBlockStartState {
	state := NewPlusBlockStartStateData(&b.atn, production, rule, len(b.atn.States()), end, loopback, decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddBasicBlockStartState(rule *generated.ParserRule, production *generated.Element, end *BlockEndState, loopback PlusLoopbackState, decision int) BasicBlockStartState {
	state := NewBasicBlockStartStateData(&b.atn, production, rule, len(b.atn.States()), end, decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddStarLoopbackState(rule *generated.ParserRule, production *generated.Element) StarLoopbackState {
	state := NewStarLoopbackStateData(&b.atn, production, rule, len(b.atn.States()))
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddPlusLoopbackState(rule *generated.ParserRule, production *generated.Element, decision int) PlusLoopbackState {
	state := NewPlusLoopbackStateData(&b.atn, production, rule, len(b.atn.States()), decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddLoopEndState(rule *generated.ParserRule, production *generated.Element, loopback *ATNState) LoopEndState {
	state := NewLoopEndStateData(&b.atn, production, rule, len(b.atn.States()), loopback)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddStarLoopEntryState(rule *generated.ParserRule, production *generated.Element, loopback *StarLoopbackState, decision int) StarLoopEntryState {
	state := NewStarLoopEntryStateData(&b.atn, production, rule, len(b.atn.States()), loopback, decision)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) AddEpsilonTransition(from *ATNState, to *ATNState) EpsilonTransition {
	transition := NewEpsilonTransitionData(to)
	(*from).AddTransition(transition)
	return transition
}

func (b *ATNBuilderData) AddAtomTransition(from *ATNState, to *ATNState, atom int) AtomTransition {
	transition := NewAtomTransitionData(to, atom)
	(*from).AddTransition(transition)
	return transition
}

func (b *ATNBuilderData) AddRuleTransition(from *ATNState, to *ATNState, rule *generated.ParserRule, followState *ATNState) RuleTransition {
	transition := NewRuleTransitionData(to, rule, followState)
	(*from).AddTransition(transition)
	return transition
}

func (b *ATNBuilderData) Build() *ATN {
	return &b.atn
}
