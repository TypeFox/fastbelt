package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type StarLoopbackState interface {
	ATNState
}

type StarLoopbackStateData struct {
	ATNStateData
}

func NewStarLoopbackStateData(atn ATN, production generated.Element, rule *generated.ParserRule, stateNumber int) *StarLoopbackStateData {
	return &StarLoopbackStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ATN_STAR_LOOP_BACK),
	}
}
