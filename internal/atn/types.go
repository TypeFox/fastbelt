// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package atn

import (
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
)

// ATNState is the single concrete ATN state type.
// Fields specific to certain state kinds are non-nil only for those kinds.
type ATNState struct {
	ATN *ATN
	// Production is the grammar element this decision state represents; set by
	// post-processing and used to generate decision-state constants.
	Production grammar.Element
	// RuleCallEntry is the rule call that caused this state to be the entry point
	// of a parser-rule invocation; used to locate the corresponding follow state.
	// Never overwritten by post-processing.
	RuleCallEntry grammar.RuleCall
	// ConsumedElement is the grammar element (keyword, token rule call, or
	// cross-reference) whose token is consumed at this state; used to generate
	// element-state constants and AssignToken calls.
	ConsumedElement        grammar.Element
	StateNumber            int
	Rule                   grammar.AbstractRuleWithBody
	EpsilonOnlyTransitions bool
	Transitions            []Transition
	Type                   parser.ATNStateType

	// Decision index; -1 if this state is not a decision state.
	Decision int

	// Populated for BlockStartState kinds.
	End      *ATNState
	Loopback *ATNState // PlusBlockStart.loopback, StarLoopEntry.loopback, LoopEnd.loopback

	// Populated for BlockEndState.
	Start *ATNState

	// Populated for RuleStartState.
	Stop *ATNState
}

// ATN is the Augmented Transition Network built from a set of parser rules.
type ATN struct {
	DecisionMap      map[string]*ATNState
	States           []*ATNState
	DecisionStates   []*ATNState
	RuleToStartState map[grammar.AbstractRuleWithBody]*ATNState
	RuleToStopState  map[grammar.AbstractRuleWithBody]*ATNState

	// OrDecision maps an Alternatives element to the decision state for its
	// alternative choice (the MakeAlternatives start state). LoopDecision maps
	// any element carrying a cardinality (?/*/+) to the decision state for its
	// enter-vs-exit choice. Both are populated during construction and used by
	// the parser generator to emit AdaptivePredict calls for decisions that a
	// single lookahead token cannot disambiguate.
	OrDecision   map[grammar.Element]*ATNState
	LoopDecision map[grammar.Element]*ATNState
}

// Transition is the interface implemented by all ATN transitions.
type Transition interface {
	IsEpsilon() bool
	Target() *ATNState
	SetTarget(*ATNState)
}

// AtomTransition fires on a specific token type.
// CategoryMatches holds the IDs of all token types that match via category
// inheritance; populated from the Terminal's CategoryMatches at ATN-build time.
//
// CompletionHint mirrors parser.CompletionHint: set on transitions that come
// from a cross-reference assignment, propagated to the runtime ATN by the
// generator so the completion provider can dispatch per-field.
type AtomTransition struct {
	TargetState    *ATNState
	TokenTypeId    int
	CompletionHint *parser.CompletionHint
}

func (t *AtomTransition) Target() *ATNState          { return t.TargetState }
func (t *AtomTransition) IsEpsilon() bool            { return false }
func (t *AtomTransition) SetTarget(target *ATNState) { t.TargetState = target }

// EpsilonTransition fires without consuming a token.
type EpsilonTransition struct {
	TargetState *ATNState
}

func (t *EpsilonTransition) Target() *ATNState          { return t.TargetState }
func (t *EpsilonTransition) IsEpsilon() bool            { return true }
func (t *EpsilonTransition) SetTarget(target *ATNState) { t.TargetState = target }

// RuleTransition enters a sub-rule and returns to FollowState.
//
// CompletionHint mirrors parser.CompletionHint: set when this rule call
// comes from a cross-reference whose text-form is a parser rule (e.g.
// `Ref=[Decl:FQN]`). The runtime simulator propagates the hint to every
// atom transition reached inside the called rule, so the completion
// provider can dispatch per-field even when the cross-reference text
// spans multiple tokens.
type RuleTransition struct {
	TargetState    *ATNState // the rule's RuleStartState
	Rule           grammar.AbstractRuleWithBody
	FollowState    *ATNState
	CompletionHint *parser.CompletionHint
}

func (t *RuleTransition) Target() *ATNState          { return t.TargetState }
func (t *RuleTransition) IsEpsilon() bool            { return true }
func (t *RuleTransition) SetTarget(target *ATNState) { t.TargetState = target }

// ATNHandle is an internal pair of (entry, exit) ATN states for a sub-network.
type ATNHandle struct {
	Left  *ATNState
	Right *ATNState
}
