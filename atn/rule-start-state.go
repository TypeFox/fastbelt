package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type RuleStartState interface {
	ATNState
	Stop() RuleStopState
}

type RuleStartStateData struct {
	ATNStateData

	stop RuleStopState
}

func NewRuleStartStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int, stop RuleStopState) RuleStartState {
	return &RuleStartStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ATN_RULE_START),
		stop:         stop,
	}
}

func (r *RuleStartStateData) Stop() RuleStopState {
	return r.stop
}
