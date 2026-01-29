package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type BasicBlockStartState interface {
	BlockStartState
}

type BasicBlockStartStateData struct {
	BlockStartStateData
}

func NewBasicBlockStartStateData(atn *ATN, production *generated.Element, rule *generated.ParserRule, stateNumber int, end *BlockEndState, decision int) *BasicBlockStartStateData {
	return &BasicBlockStartStateData{
		BlockStartStateData: *NewBlockStartStateData(atn, production, rule, stateNumber, decision, end, ATN_BASIC),
	}
}
