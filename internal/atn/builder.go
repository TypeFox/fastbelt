// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package atn

import (
	"slices"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
)

type ATNRuleBuilder interface {
	Assign(handle *ATNHandle)

	Plus(lookaheadName string, handle *ATNHandle) *ATNHandle
	Star(lookaheadName string, handle *ATNHandle) *ATNHandle
	Optional(lookaheadName string, handle *ATNHandle) *ATNHandle
	MakeAlternatives(lookaheadName string, start *ATNState, alts []*ATNHandle) *ATNHandle
	MakeConcatenation(alts []*ATNHandle) *ATNHandle
	TokenRef(tokenTypeId int) *ATNHandle
	RuleRef(otherRule *grammar.ParserRule) *ATNHandle

	NewEpsilonTransition(source *ATNState, target *ATNState)

	NewState(stateType parser.ATNStateType) *ATNState
	RemoveState(state *ATNState)

	GetTokenTypeByName(name string) int
	GetRuleByName(name string) *grammar.ParserRule
	GetCompositeByName(name string) grammar.CompositeRule
	GetLookaheadNameByElement(el grammar.Element) string
}

type ATNBuilder interface {
	DeclareRule(rule *grammar.ParserRule) ATNRuleBuilder
	Build() *ATN
}

type ATNBuilderImpl struct {
	rules            map[*grammar.ParserRule]*ATNRuleBuilderImpl
	atn              *ATN
	names            map[grammar.Element]string
	tokenTypeIds     map[string]int
	rulesByName      map[string]*grammar.ParserRule
	compositesByName map[string]grammar.CompositeRule
}

func NewATNBuilder(names map[grammar.Element]string, tokenTypeIds map[string]int, rulesByName map[string]*grammar.ParserRule, compositesByName map[string]grammar.CompositeRule) *ATNBuilderImpl {
	return &ATNBuilderImpl{
		rules: map[*grammar.ParserRule]*ATNRuleBuilderImpl{},
		atn: &ATN{
			DecisionMap:      map[string]*ATNState{},
			States:           []*ATNState{},
			DecisionStates:   []*ATNState{},
			RuleToStartState: map[grammar.ParserRule]*ATNState{},
			RuleToStopState:  map[grammar.ParserRule]*ATNState{},
		},
		names:            names,
		tokenTypeIds:     tokenTypeIds,
		rulesByName:      rulesByName,
		compositesByName: compositesByName,
	}
}

func (rb *ATNRuleBuilderImpl) GetTokenTypeByName(name string) int {
	if id, ok := rb.parent.tokenTypeIds[name]; ok {
		return id
	}
	return -1
}

func (rb *ATNRuleBuilderImpl) GetRuleByName(name string) *grammar.ParserRule {
	if rule, ok := rb.parent.rulesByName[name]; ok {
		return rule
	}
	return nil
}

func (rb *ATNRuleBuilderImpl) GetCompositeByName(name string) grammar.CompositeRule {
	if composite, ok := rb.parent.compositesByName[name]; ok {
		return composite
	}
	return nil
}

func (rb *ATNRuleBuilderImpl) GetLookaheadNameByElement(el grammar.Element) string {
	if name, ok := rb.parent.names[el]; ok {
		return name
	}
	panic("no lookahead name found for element")
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

func (rb *ATNRuleBuilderImpl) MakeAlternatives(lookaheadName string, start *ATNState, alts []*ATNHandle) *ATNHandle {
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

func (rb *ATNRuleBuilderImpl) MakeConcatenation(alts []*ATNHandle) *ATNHandle {
	for i := 0; i < len(alts)-1; i++ {
		handle := alts[i]
		next := alts[i+1].Left
		if rb.isConcatenationOptimizable(handle) {
			rb.optimizeConcatenation(handle, handle.Left.Transitions[0], next)
		} else {
			rb.NewEpsilonTransition(handle.Right, next)
		}
	}
	first := alts[0]
	last := alts[len(alts)-1]
	return &ATNHandle{Left: first.Left, Right: last.Right}
}

func (rb *ATNRuleBuilderImpl) isConcatenationOptimizable(handle *ATNHandle) bool {
	//Without this optimization: alts[i].Left:ATNBasic -- :RuleTransition --> alts[i].Right:ATNBasic -epsilon-> alts[i+1].Left
	//With this optimization:    alts[i].Left:ATNBasic -- :RuleTransition --> alts[i+1].Left
	//saves one ATN state and one epsilon transition per concatenation, which can add up in large grammars
	if len((*handle).Left.Transitions) != 1 {
		return false
	}
	transition := handle.Left.Transitions[0]
	ruleTransition, isRuleTransition := transition.(*RuleTransition)
	right := handle.Right
	return handle.Left.Type == parser.ATNBasic && right.Type == parser.ATNBasic &&
		((isRuleTransition && ruleTransition.FollowState == right) || transition.Target() == right)
}

func (rb *ATNRuleBuilderImpl) optimizeConcatenation(handle *ATNHandle, transition Transition, next *ATNState) {
	ruleTransition, isRuleTransition := transition.(*RuleTransition)
	if isRuleTransition {
		ruleTransition.FollowState = next
	} else {
		transition.SetTarget(next)
	}
	rb.RemoveState(handle.Right)
}

func (rb *ATNRuleBuilderImpl) TokenRef(tokenTypeId int) *ATNHandle {
	left := rb.NewState(parser.ATNBasic)
	right := rb.NewState(parser.ATNBasic)
	addTransition(left, &AtomTransition{
		TargetState: right,
		TokenTypeId: tokenTypeId,
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

func (rb *ATNRuleBuilderImpl) RemoveState(state *ATNState) {
	index := slices.Index(rb.parent.atn.States, state)
	if index > -1 {
		rb.parent.atn.States = append(rb.parent.atn.States[:index], rb.parent.atn.States[index+1:]...)
	}
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
