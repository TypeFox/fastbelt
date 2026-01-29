package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type StarBlockStartState interface {
	BlockStartState
}

type StarBlockStartStateData struct {
	BlockStartStateData
}

func NewStarBlockStartStateData(atn *ATN, production *generated.Element, end *BlockEndState, stateNumber int, rule *generated.ParserRule) *StarBlockStartStateData {
	return &StarBlockStartStateData{
		BlockStartStateData: *NewBlockStartStateData(atn, production, rule, stateNumber, 0, end, ATN_STAR_BLOCK_START),
	}
}
