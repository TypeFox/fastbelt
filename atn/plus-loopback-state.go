package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type PlusLoopbackState interface {
	DecisionState
}

type PlusLoopbackStateData struct {
	DecisionStateData
}

func NewPlusLoopbackStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int, decision int) *PlusLoopbackStateData {
	return &PlusLoopbackStateData{
		DecisionStateData: *NewDecisionStateData(atn, production, rule, stateNumber, decision, ATN_PLUS_LOOP_BACK),
	}
}
