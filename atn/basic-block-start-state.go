package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type BasicBlockStartState interface {
	BlockStartState
}

type BasicBlockStartStateData struct {
	BlockStartStateData
}

func NewBasicBlockStartStateData(atn *ATN, production *generated.Element, end *BlockEndState, stateNumber int, rule *generated.ParserRule, ty int) *BasicBlockStartStateData {
	return &BasicBlockStartStateData{
		BlockStartStateData: *NewBlockStartStateData(atn, production, rule, stateNumber, 0, end, ty),
	}
}
