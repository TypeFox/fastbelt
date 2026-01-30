package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type ATN interface {
	States() []*ATNState
	RuleStartState(rule *generated.ParserRule) *RuleStartState
	RuleStopState(rule *generated.ParserRule) *RuleStopState
	AddState(state ATNState)
	AddStartAndStopState(rule *generated.ParserRule, start *RuleStartState, stop *RuleStopState)
}

type ATNData struct {
	states           []*ATNState
	ruleToStartState map[*generated.ParserRule]*RuleStartState
	ruleToStopState  map[*generated.ParserRule]*RuleStopState
}

func (a *ATNData) States() []*ATNState {
	return a.states
}

func (a *ATNData) AddState(state ATNState) {
	a.states = append(a.states, &state)
}

func (a *ATNData) AddStartAndStopState(rule *generated.ParserRule, start *RuleStartState, stop *RuleStopState) {
	a.ruleToStartState[rule] = start
	a.ruleToStopState[rule] = stop
}

func (a *ATNData) RuleStartState(rule *generated.ParserRule) *RuleStartState {
	return a.ruleToStartState[rule]
}

func (a *ATNData) RuleStopState(rule *generated.ParserRule) *RuleStopState {
	return a.ruleToStopState[rule]
}

func NewATN() ATN {
	return &ATNData{
		states:           []*ATNState{},
		ruleToStartState: make(map[*generated.ParserRule]*RuleStartState),
		ruleToStopState:  make(map[*generated.ParserRule]*RuleStopState),
	}
}

func NewATNFromGrammar(grammar generated.Grammar) *ATN {
	atnBuilder := NewATNBuilder()
	atnBuilder.InititializeStartAndStopStates(grammar.Rules())
	return atnBuilder.Build()
}
