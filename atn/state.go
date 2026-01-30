package atn

import "typefox.dev/fastbelt/internal/grammar/generated"

const (
	ATN_INVALID_TYPE = iota + 1
	ATN_BASIC
	ATN_RULE_START
	ATN_PLUS_BLOCK_START
	ATN_STAR_BLOCK_START
	ATN_RULE_STOP
	ATN_BLOCK_END
	ATN_STAR_LOOP_BACK
	ATN_STAR_LOOP_ENTRY
	ATN_PLUS_LOOP_BACK
	ATN_LOOP_END
)

type ATNState interface {
	Type() int
	ATN() ATN
	Production() generated.Element
	StateNumber() int
	Rule() generated.ParserRule
	EpsilonOnlyTransitions() bool
	Transitions() []Transition
	NextTokenWithinRule() []int
	AddTransition(transition Transition)
}

type ATNStateData struct {
	atn                    ATN
	production             generated.Element
	stateNumber            int
	rule                   generated.ParserRule
	epsilonOnlyTransitions bool
	transitions            []Transition
	nextTokenWithinRule    []int
	ty                     int
}

func NewATNStateData(atn ATN, production generated.Element, rule generated.ParserRule, stateNumber int, ty int) *ATNStateData {
	return &ATNStateData{
		atn:                    atn,
		production:             production,
		epsilonOnlyTransitions: false,
		rule:                   rule,
		transitions:            []Transition{},
		nextTokenWithinRule:    []int{},
		stateNumber:            stateNumber,
		ty:                     ty,
	}
}

func (b *ATNStateData) Type() int {
	return b.ty
}

func (b *ATNStateData) ATN() ATN {
	return b.atn
}

func (b *ATNStateData) Production() generated.Element {
	return b.production
}

func (b *ATNStateData) StateNumber() int {
	return b.stateNumber
}

func (b *ATNStateData) Rule() generated.ParserRule {
	return b.rule
}

func (b *ATNStateData) EpsilonOnlyTransitions() bool {
	return b.epsilonOnlyTransitions
}

func (b *ATNStateData) Transitions() []Transition {
	return b.transitions
}

func (b *ATNStateData) NextTokenWithinRule() []int {
	return b.nextTokenWithinRule
}

func (b *ATNStateData) AddTransition(transition Transition) {
	b.transitions = append(b.transitions, transition)
}
