package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type PlusBlockStartState interface {
	BlockStartState
	Loopback() PlusLoopbackState
}

type PlusBlockStartStateData struct {
	BlockStartStateData
	loopback PlusLoopbackState
}

func NewPlusBlockStartStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int, end BlockEndState, loopback PlusLoopbackState, decision int) *PlusBlockStartStateData {
	return &PlusBlockStartStateData{
		BlockStartStateData: *NewBlockStartStateData(atn, production, rule, stateNumber, decision, end, ATN_PLUS_BLOCK_START),
		loopback:            loopback,
	}
}

func (p *PlusBlockStartStateData) Loopback() PlusLoopbackState {
	return p.loopback
}
