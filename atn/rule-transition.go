package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type RuleTransition interface {
	Rule() *generated.ParserRule
	FollowState() *ATNState
}

type RuleTransitionData struct {
	AbstractTransitionData
	rule        *generated.ParserRule
	followState *ATNState
}

func NewRuleTransitionData(target ATNState, rule *generated.ParserRule, followState *ATNState) *RuleTransitionData {
	return &RuleTransitionData{
		AbstractTransitionData: NewTransitionData(target),
		rule:                   rule,
		followState:            followState,
	}
}

func (r *RuleTransitionData) Rule() *generated.ParserRule {
	return r.rule
}

func (r *RuleTransitionData) FollowState() *ATNState {
	return r.followState
}
