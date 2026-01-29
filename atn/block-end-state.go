package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type BlockEndState interface {
	ATNState
	Start() BlockStartState
}

type BlockEndStateData struct {
	ATNStateData
	start BlockStartState
}

func NewBlockEndStateData(atn ATN, production generated.Element, start BlockStartState, stateNumber int, rule *generated.ParserRule) *BlockEndStateData {
	return &BlockEndStateData{
		ATNStateData: *NewATNStateData(atn, production, rule, stateNumber, ATN_BLOCK_END),
		start:        start,
	}
}

func (b *BlockEndStateData) Start() BlockStartState {
	return b.start
}
