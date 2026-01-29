package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type ATNBuilder interface {
	AddRuleStartState(production *generated.Element, rule *generated.ParserRule, stop *RuleStopState) RuleStartState
	Build() ATN
}

type ATNBuilderData struct {
	atn ATN
}

func NewATNBuilder() ATNBuilderData {
	return ATNBuilderData{
		atn: NewATN(),
	}
}

func (b *ATNBuilderData) AddRuleStartState(production *generated.Element, rule *generated.ParserRule, stop *RuleStopState) RuleStartState {
	state := NewRuleStartStateData(&b.atn, production, rule, len(b.atn.States()), stop)
	b.atn.AddState(state)
	return state
}

func (b *ATNBuilderData) Build() ATN {
	return b.atn
}
