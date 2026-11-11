// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
)

// TestNewRuntimeATNFromTables builds the ATN of
//
//	Start: Items+=Item*;
//	Item:  name=ID ref=[Item:ID];
//
// from tables and checks the resulting structure.
func TestNewRuntimeATNFromTables(t *testing.T) {
	id := createTokenType(1, "ID")
	const (
		startStart = iota
		startStop
		itemStart
		itemStop
		loopEntry
		callItem
		loopEnd
		itemName
		itemRef
		itemEnd
	)
	tables := ATNTables{
		States: []ATNStateEntry{
			startStart: {Type: ATNRuleStart, EpsilonOnly: true, Decision: -1, Transitions: 1},
			startStop:  {Type: ATNRuleStop, Decision: -1},
			itemStart:  {Type: ATNRuleStart, EpsilonOnly: true, Decision: -1, Transitions: 1},
			itemStop:   {Type: ATNRuleStop, Decision: -1},
			loopEntry:  {Type: ATNLoopEntry, EpsilonOnly: true, Decision: 1, Transitions: 2},
			callItem:   {Type: ATNBasic, EpsilonOnly: true, Decision: -1, Transitions: 1},
			loopEnd:    {Type: ATNLoopEnd, EpsilonOnly: true, Decision: -1, Transitions: 1},
			itemName:   {Type: ATNBasic, Decision: -1, Transitions: 1},
			itemRef:    {Type: ATNBasic, Decision: -1, Transitions: 1},
			itemEnd:    {Type: ATNBasic, EpsilonOnly: true, Decision: -1, Transitions: 1},
		},
		Transitions: []ATNTransitionEntry{
			// startStart
			{Kind: ATNEpsilonTransition, Target: loopEntry},
			// itemStart
			{Kind: ATNEpsilonTransition, Target: itemName},
			// loopEntry
			{Kind: ATNEpsilonTransition, Target: callItem},
			{Kind: ATNEpsilonTransition, Target: loopEnd},
			// callItem
			{Kind: ATNRuleTransition, Target: itemStart, Follow: loopEntry, Call: 1},
			// loopEnd
			{Kind: ATNEpsilonTransition, Target: startStop},
			// itemName
			{Kind: ATNAtomTransition, Target: itemRef, Token: 1},
			// itemRef
			{Kind: ATNAtomTransition, Target: itemEnd, Token: 1, Hint: 1},
			// itemEnd
			{Kind: ATNEpsilonTransition, Target: itemStop},
		},
		DecisionStates:  []int32{-1, loopEntry},
		TokenTypes:      []*core.TokenType{1: id},
		CompletionHints: []CompletionHint{{Owner: "Item", Property: "Ref"}},
		RuleCalls:       []RuleCallInfo{{Field: "Items", List: true, Type: "Item", Repeated: true}},
	}

	atn := NewRuntimeATNFromTables(tables)

	require.Len(t, atn.States, len(tables.States))
	for i, s := range atn.States {
		assert.Equal(t, i, s.StateNumber)
		assert.Equal(t, tables.States[i].Type, s.Type)
		assert.Equal(t, int(tables.States[i].Decision), s.Decision)
		assert.Equal(t, tables.States[i].EpsilonOnly, s.EpsilonOnlyTransitions)
		assert.Len(t, s.Transitions, int(tables.States[i].Transitions))
	}
	assert.Equal(t, []*RuntimeATNState{nil, atn.States[loopEntry]}, atn.DecisionStates)
	assert.Equal(t, atn.DecisionStates, atn.DecisionMap)

	entry := atn.States[loopEntry].Transitions
	assert.Same(t, atn.States[callItem], entry[0].GetTarget())
	assert.Same(t, atn.States[loopEnd], entry[1].GetTarget())

	call, ok := atn.States[callItem].Transitions[0].(*RuntimeRuleTransition)
	require.True(t, ok)
	assert.Same(t, atn.States[itemStart], call.Target)
	assert.Same(t, atn.States[loopEntry], call.FollowState)
	assert.Nil(t, call.CompletionHint)
	assert.Equal(t, &RuleCallInfo{Field: "Items", List: true, Type: "Item", Repeated: true}, call.Call)
	assert.Same(t, call.Call, atn.ruleCallAt(loopEntry))

	name, ok := atn.States[itemName].Transitions[0].(*RuntimeAtomTransition)
	require.True(t, ok)
	assert.Same(t, id, name.TokenType)
	assert.Nil(t, name.CompletionHint)
	ref, ok := atn.States[itemRef].Transitions[0].(*RuntimeAtomTransition)
	require.True(t, ok)
	assert.Equal(t, &CompletionHint{Owner: "Item", Property: "Ref"}, ref.CompletionHint)

	assert.True(t, atn.NextTokensAt(startStart).At(id.Id))
	assert.Nil(t, atn.States[startStop].Transitions)

	// The transitions of the states share one backing array; appending to
	// those of one state must not change those of the next one.
	atn.States[callItem].AppendTransitions(&RuntimeEpsilonTransition{Target: atn.States[startStop]})
	assert.Len(t, atn.States[callItem].Transitions, 2)
	assert.Same(t, atn.States[startStop], atn.States[loopEnd].Transitions[0].GetTarget())
}

func TestNewRuntimeATNFromTablesPanicsOnInconsistentTables(t *testing.T) {
	assert.Panics(t, func() {
		NewRuntimeATNFromTables(ATNTables{
			States:      []ATNStateEntry{{Type: ATNBasic, Decision: -1, Transitions: 2}},
			Transitions: []ATNTransitionEntry{{Kind: ATNEpsilonTransition, Target: 0}},
		})
	})
}
