package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type RuleStopState interface {
	ATNState
}

type RuleStopStateData struct {
	ATNStateData
}

func NewRuleStopStateData(atn ATN, production generated.Element, rule *generated.ParserRule, stateNumber int) *RuleStopStateData {
	return &RuleStopStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ATN_RULE_STOP),
	}
}
