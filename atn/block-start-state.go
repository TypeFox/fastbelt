package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

type BlockStartState interface {
	DecisionState
	End() *BlockEndState
}

type BlockStartStateData struct {
	DecisionStateData
	end *BlockEndState
}

func NewBlockStartStateData(atn ATN, production generated.Element, rule *generated.ParserRule, stateNumber int, decision int, end *BlockEndState, ty int) *BlockStartStateData {
	return &BlockStartStateData{
		DecisionStateData: *NewDecisionStateData(atn, production, rule, stateNumber, decision, ty),
		end:               end,
	}
}

func (b *BlockStartStateData) End() *BlockEndState {
	return b.end
}
