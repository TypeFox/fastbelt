package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type BasicState interface {
	ATNState
}

type BasicStateData struct {
	ATNStateData
}

func NewBasicStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int) *BasicStateData {
	return &BasicStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ATN_BASIC),
	}
}
