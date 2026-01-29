package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type StarLoopEntryState interface {
	DecisionState
	Loopback() *StarLoopbackState
}

type StarLoopEntryStateData struct {
	DecisionStateData
	loopback *StarLoopbackState
}

func NewStarLoopEntryStateData(atn *ATN, production *generated.Element, rule *generated.ParserRule, stateNumber int, decision int, loopback *StarLoopbackState) *StarLoopEntryStateData {
	return &StarLoopEntryStateData{
		DecisionStateData: *NewDecisionStateData(atn, production, rule, stateNumber, decision, ATN_STAR_LOOP_ENTRY),
		loopback:          loopback,
	}
}

func (s *StarLoopEntryStateData) Loopback() *StarLoopbackState {
	return s.loopback
}
