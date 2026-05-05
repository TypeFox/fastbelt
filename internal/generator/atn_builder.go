// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
)

type ATNRuleBuilder interface {
	Assign(handle *ATNHandle)

	Plus(lookaheadName string, handle *ATNHandle) *ATNHandle
	Star(lookaheadName string, handle *ATNHandle) *ATNHandle
	Optional(lookaheadName string, handle *ATNHandle) *ATNHandle
	MakeAlts(lookaheadName string, start *ATNState, alts []*ATNHandle) *ATNHandle

	MakeBlock(alts []*ATNHandle) *ATNHandle
	TokenRef(tokenType TokenType) *ATNHandle
	RuleRef(otherRule *grammar.ParserRule) *ATNHandle

	NewEpsilonTransition(source *ATNState, target *ATNState)

	NewState(stateType parser.ATNStateType) *ATNState
}

type ATNBuilder interface {
	DeclareRule(rule *grammar.ParserRule) ATNRuleBuilder
	Build() *ATN
}

type ATNBuilderImpl struct {
	rules map[*grammar.ParserRule]*ATNRuleBuilderImpl
	atn   *ATN
}

func NewATNBuilder() *ATNBuilderImpl {
	return &ATNBuilderImpl{
		rules: map[*grammar.ParserRule]*ATNRuleBuilderImpl{},
		atn: &ATN{
			DecisionMap:      map[string]*ATNState{},
			States:           []*ATNState{},
			DecisionStates:   []*ATNState{},
			RuleToStartState: map[grammar.ParserRule]*ATNState{},
			RuleToStopState:  map[grammar.ParserRule]*ATNState{},
		},
	}
}

func (b *ATNBuilderImpl) DeclareRule(rule *grammar.ParserRule) ATNRuleBuilder {
	left := newATNState(b.atn, rule, parser.ATNRuleStart)
	right := newATNState(b.atn, rule, parser.ATNRuleStop)
	ruleBuilder := NewATNRuleBuilder(b, rule, &ATNHandle{Left: left, Right: right})
	b.rules[rule] = ruleBuilder
	b.atn.RuleToStartState[*rule] = left
	b.atn.RuleToStopState[*rule] = right
	return ruleBuilder
}

func (b *ATNBuilderImpl) Build() *ATN {
	return b.atn
}

type ATNRuleBuilderImpl struct {
	parent *ATNBuilderImpl
	rule   *grammar.ParserRule
	handle *ATNHandle
}

func NewATNRuleBuilder(parent *ATNBuilderImpl, rule *grammar.ParserRule, handle *ATNHandle) *ATNRuleBuilderImpl {
	return &ATNRuleBuilderImpl{
		parent: parent,
		rule:   rule,
		handle: handle,
	}
}

func (rb *ATNRuleBuilderImpl) Assign(handle *ATNHandle) {
	rb.NewEpsilonTransition(rb.handle.Left, handle.Left)
	rb.NewEpsilonTransition(handle.Right, rb.handle.Right)
}

func (rb *ATNRuleBuilderImpl) Plus(lookaheadName string, handle *ATNHandle) *ATNHandle {
	atn := rb.parent.atn

	blkStart := handle.Left
	blkEnd := handle.Right

	loop := rb.NewState(parser.ATNPlusLoopBack)
	defineDecisionState(atn, loop)
	end := rb.NewState(parser.ATNLoopEnd)
	blkStart.Loopback = loop
	end.Loopback = loop

	atn.DecisionMap[lookaheadName] = loop

	rb.NewEpsilonTransition(blkEnd, loop)   // block can see loop back
	rb.NewEpsilonTransition(loop, blkStart) // loop back to start
	rb.NewEpsilonTransition(loop, end)      // exit

	return &ATNHandle{Left: blkStart, Right: end}
}

func (rb *ATNRuleBuilderImpl) Star(lookaheadName string, handle *ATNHandle) *ATNHandle {
	atn := rb.parent.atn

	start := handle.Left
	end := handle.Right

	entry := rb.NewState(parser.ATNStarLoopEntry)
	defineDecisionState(atn, entry)
	loopEnd := rb.NewState(parser.ATNLoopEnd)
	loop := rb.NewState(parser.ATNStarLoopBack)
	entry.Loopback = loop
	loopEnd.Loopback = loop

	rb.NewEpsilonTransition(entry, start)   // loop enter edge (alt 0)
	rb.NewEpsilonTransition(entry, loopEnd) // bypass loop edge (alt 1)
	rb.NewEpsilonTransition(end, loop)      // block end hits loop back
	rb.NewEpsilonTransition(loop, entry)    // loop back to entry/exit decision

	atn.DecisionMap[lookaheadName] = entry

	return &ATNHandle{Left: entry, Right: loopEnd}
}

func (rb *ATNRuleBuilderImpl) Optional(lookaheadName string, handle *ATNHandle) *ATNHandle {
	start := handle.Left
	end := handle.Right

	rb.NewEpsilonTransition(start, end)

	rb.parent.atn.DecisionMap[lookaheadName] = start
	return handle
}

func (rb *ATNRuleBuilderImpl) MakeAlts(lookaheadName string, start *ATNState, alts []*ATNHandle) *ATNHandle {
	atn := rb.parent.atn
	end := rb.NewState(parser.ATNBlockEnd)
	end.Start = start
	start.End = end

	for _, alt := range alts {
		if alt != nil {
			rb.NewEpsilonTransition(start, alt.Left)
			rb.NewEpsilonTransition(alt.Right, end)
		} else {
			rb.NewEpsilonTransition(start, end)
		}
	}

	atn.DecisionMap[lookaheadName] = start

	return &ATNHandle{Left: start, Right: end}
}

func (rb *ATNRuleBuilderImpl) MakeBlock(alts []*ATNHandle) *ATNHandle {
	for i := 0; i < len(alts)-1; i++ {
		handle := alts[i]
		next := alts[i+1].Left
		rb.NewEpsilonTransition(handle.Right, next)
	}
	first := alts[0]
	last := alts[len(alts)-1]
	return &ATNHandle{Left: first.Left, Right: last.Right}
}

func (rb *ATNRuleBuilderImpl) TokenRef(tokenType TokenType) *ATNHandle {
	left := rb.NewState(parser.ATNBasic)
	right := rb.NewState(parser.ATNBasic)
	addTransition(left, &AtomTransition{
		TargetState: right,
		TokenType:   tokenType,
	})
	return &ATNHandle{Left: left, Right: right}
}

func (rb *ATNRuleBuilderImpl) RuleRef(otherRule *grammar.ParserRule) *ATNHandle {
	ruleStart := rb.parent.atn.RuleToStartState[*otherRule]
	left := rb.NewState(parser.ATNBasic)
	right := rb.NewState(parser.ATNBasic)
	addTransition(left, &RuleTransition{
		TargetState: ruleStart,
		Rule:        otherRule,
		FollowState: right,
	})
	return &ATNHandle{Left: left, Right: right}
}

func (rb *ATNRuleBuilderImpl) NewEpsilonTransition(source *ATNState, target *ATNState) {
	addTransition(source, &EpsilonTransition{TargetState: target})
}

func (rb *ATNRuleBuilderImpl) NewState(stateType parser.ATNStateType) *ATNState {
	return newATNState(rb.parent.atn, rb.rule, stateType)
}

func newATNState(atn *ATN, rule *grammar.ParserRule, typ parser.ATNStateType) *ATNState {
	s := &ATNState{
		ATN:         atn,
		Production:  nil,
		StateNumber: len(atn.States),
		Rule:        rule,
		Type:        typ,
		Decision:    -1,
	}
	atn.States = append(atn.States, s)
	return s
}

func defineDecisionState(atn *ATN, state *ATNState) int {
	atn.DecisionStates = append(atn.DecisionStates, state)
	state.Decision = len(atn.DecisionStates) - 1
	return state.Decision
}

func addTransition(state *ATNState, t Transition) {
	if len(state.Transitions) == 0 {
		state.EpsilonOnlyTransitions = t.IsEpsilon()
	}
	state.Transitions = append(state.Transitions, t)
}
