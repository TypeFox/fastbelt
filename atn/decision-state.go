package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type DecisionState interface {
	ATNState
	Decision() int
}

type DecisionStateData struct {
	ATNStateData
	decision int
}

func NewDecisionStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int, decision int, ty int) *DecisionStateData {
	return &DecisionStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ty),
		decision:     decision,
	}
}

func (st *DecisionStateData) Decision() int {
	return st.decision
}
