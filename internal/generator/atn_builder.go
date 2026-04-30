// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
)

type ATNRuleBuilder interface {
	RuleHandle() *ATNHandle

	Plus(lookaheadName string, handle *ATNHandle) *ATNHandle
	Star(lookaheadName string, handle *ATNHandle) *ATNHandle
	Optional(lookaheadName string, handle *ATNHandle) *ATNHandle
	MakeAlts(lookaheadName string, start *ATNState, alts []*ATNHandle) *ATNHandle

	MakeBlock(alts []*ATNHandle) *ATNHandle
	TokenRef(tokenType TokenInfo) *ATNHandle
	RuleRef(otherRule *grammar.ParserRule) *ATNHandle

	NewEpsilon(source *ATNState, target *ATNState)

	NewState(stateType ATNStateType) *ATNState
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
		atn:   &ATN{},
	}
}

func (b *ATNBuilderImpl) DeclareRule(rule *grammar.ParserRule) ATNRuleBuilder {
	left := newATNState(b.atn, rule, ATNRuleStart)
	right := newATNState(b.atn, rule, ATNRuleStop)
	ruleBuilder := NewATNRuleBuilder(b, rule, &ATNHandle{Left: left, Right: right})
	b.rules[rule] = ruleBuilder
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

func (rb *ATNRuleBuilderImpl) RuleHandle() *ATNHandle {
	return rb.handle
}

func (rb *ATNRuleBuilderImpl) Plus(lookaheadName string, handle *ATNHandle) *ATNHandle {
	atn := rb.parent.atn

	blkStart := handle.Left
	blkEnd := handle.Right

	loop := rb.NewState(ATNPlusLoopBack)
	defineDecisionState(atn, loop)
	end := rb.NewState(ATNLoopEnd)
	blkStart.Loopback = loop
	end.Loopback = loop

	atn.DecisionMap[lookaheadName] = loop

	addEpsilon(blkEnd, loop)   // block can see loop back
	addEpsilon(loop, blkStart) // loop back to start
	addEpsilon(loop, end)      // exit

	return &ATNHandle{Left: blkStart, Right: end}
}

func (rb *ATNRuleBuilderImpl) Star(lookaheadName string, handle *ATNHandle) *ATNHandle {
	atn := rb.parent.atn

	start := handle.Left
	end := handle.Right

	entry := rb.NewState(ATNStarLoopEntry)
	defineDecisionState(atn, entry)
	loopEnd := rb.NewState(ATNLoopEnd)
	loop := rb.NewState(ATNStarLoopBack)
	entry.Loopback = loop
	loopEnd.Loopback = loop

	addEpsilon(entry, start)   // loop enter edge (alt 0)
	addEpsilon(entry, loopEnd) // bypass loop edge (alt 1)
	addEpsilon(end, loop)      // block end hits loop back
	addEpsilon(loop, entry)    // loop back to entry/exit decision

	atn.DecisionMap[lookaheadName] = entry

	return &ATNHandle{Left: entry, Right: loopEnd}
}

func (rb *ATNRuleBuilderImpl) Optional(lookaheadName string, handle *ATNHandle) *ATNHandle {
	start := handle.Left
	end := handle.Right

	addEpsilon(start, end)

	rb.parent.atn.DecisionMap[lookaheadName] = start
	return handle
}

func (rb *ATNRuleBuilderImpl) MakeAlts(lookaheadName string, start *ATNState, alts []*ATNHandle) *ATNHandle {
	atn := rb.parent.atn
	end := rb.NewState(ATNBlockEnd)
	end.Start = start
	start.End = end

	for _, alt := range alts {
		if alt != nil {
			addEpsilon(start, alt.Left)
			addEpsilon(alt.Right, end)
		} else {
			addEpsilon(start, end)
		}
	}

	atn.DecisionMap[lookaheadName] = start

	return &ATNHandle{Left: start, Right: end}
}

func (rb *ATNRuleBuilderImpl) MakeBlock(alts []*ATNHandle) *ATNHandle {
	for i := 0; i < len(alts)-1; i++ {
		handle := alts[i]
		next := alts[i+1].Left
		var t Transition
		if len(handle.Left.Transitions) == 1 {
			t = handle.Left.Transitions[0]
		}
		rt, isRule := t.(*RuleTransition)
		if handle.Left.Type == ATNBasic &&
			handle.Right.Type == ATNBasic &&
			t != nil &&
			((isRule && rt.FollowState == handle.Right) ||
				(!isRule && t.Target() == handle.Right)) {
			// avoid epsilon edge to next element
			if isRule {
				rt.FollowState = next
			} else {
				switch at := t.(type) {
				case *AtomTransition:
					at.TargetState = next
				case *EpsilonTransition:
					at.TargetState = next
				}
			}
			removeState(rb.parent.atn, handle.Right) // we skipped over this state
		} else {
			addEpsilon(handle.Right, next)
		}
	}

	first := alts[0]
	last := alts[len(alts)-1]
	return &ATNHandle{Left: first.Left, Right: last.Right}
}

func (rb *ATNRuleBuilderImpl) TokenRef(tokenType TokenInfo) *ATNHandle {
	left := rb.NewState(ATNBasic)
	right := rb.NewState(ATNBasic)
	addTransition(left, &AtomTransition{
		TargetState:     right,
		TokenTypeID:     tokenType.ID,
		CategoryMatches: tokenType.CategoryMatches,
	})
	return &ATNHandle{Left: left, Right: right}
}

func (rb *ATNRuleBuilderImpl) RuleRef(otherRule *grammar.ParserRule) *ATNHandle {
	ruleStart := rb.parent.atn.RuleToStartState[otherRule]
	left := rb.NewState(ATNBasic)
	right := rb.NewState(ATNBasic)
	addTransition(left, &RuleTransition{
		TargetState: ruleStart,
		Rule:        otherRule,
		FollowState: right,
	})
	return &ATNHandle{Left: left, Right: right}
}

func (rb *ATNRuleBuilderImpl) NewEpsilon(source *ATNState, target *ATNState) {
	addEpsilon(source, target)
}

func (rb *ATNRuleBuilderImpl) NewState(stateType ATNStateType) *ATNState {
	return newATNState(rb.parent.atn, rb.rule, stateType)
}

func newATNState(atn *ATN, rule *grammar.ParserRule, typ ATNStateType) *ATNState {
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

func addEpsilon(from, to *ATNState) {
	addTransition(from, &EpsilonTransition{TargetState: to})
}

func addTransition(state *ATNState, t Transition) {
	if len(state.Transitions) == 0 {
		state.EpsilonOnlyTransitions = t.IsEpsilon()
	}
	state.Transitions = append(state.Transitions, t)
}

func removeState(atn *ATN, state *ATNState) {
	for i, s := range atn.States {
		if s == state {
			atn.States = append(atn.States[:i], atn.States[i+1:]...)
			return
		}
	}
}
