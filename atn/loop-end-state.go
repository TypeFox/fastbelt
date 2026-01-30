package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type LoopEndState interface {
	ATNState
	Loopback() ATNState
}

type LoopEndStateData struct {
	ATNStateData
	loopback ATNState
}

func NewLoopEndStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int, loopback ATNState) *LoopEndStateData {
	return &LoopEndStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ATN_LOOP_END),
		loopback:     loopback,
	}
}

func (l *LoopEndStateData) Loopback() ATNState {
	return l.loopback
}
