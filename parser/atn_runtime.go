// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

// ATNStateType is the discriminator for ATN state kinds.
type ATNStateType int

const (
	ATNInvalidType    ATNStateType = 0
	ATNBasic          ATNStateType = 1
	ATNRuleStart      ATNStateType = 2
	ATNPlusBlockStart ATNStateType = 4
	ATNStarBlockStart ATNStateType = 5
	ATNTokenStart     ATNStateType = 6
	ATNRuleStop       ATNStateType = 7
	ATNBlockEnd       ATNStateType = 8
	ATNStarLoopBack   ATNStateType = 9
	ATNStarLoopEntry  ATNStateType = 10
	ATNPlusLoopBack   ATNStateType = 11
	ATNLoopEnd        ATNStateType = 12
)

// RuntimeATNState holds only the fields required for prediction at runtime.
// Back-pointers to the build-time grammar objects (Rule, Production) and
// structural scaffolding (End, Loopback, Start, Stop) are absent.
//
// RuleName and Prod* are populated only on decision states and are used
// exclusively for ambiguity-report messages.
type RuntimeATNState struct {
	StateNumber            int
	Type                   ATNStateType
	Decision               int
	EpsilonOnlyTransitions bool
	Transitions            []RuntimeTransition
	// Ambiguity reporting — only set on decision states.
	RuleName string
	ProdKind string
	ProdIdx  int
}

// RuntimeTransition is the interface implemented by all runtime ATN transitions.
type RuntimeTransition interface {
	GetTarget() *RuntimeATNState
	IsEpsilon() bool
}

// RuntimeAtomTransition fires on a specific token type.
// Target, TokenTypeID, and CategoryMatches are exported so that generated
// Go code in external packages can construct them as struct literals.
type RuntimeAtomTransition struct {
	Target          *RuntimeATNState
	TokenTypeID     int
	CategoryMatches []int
}

func (t *RuntimeAtomTransition) GetTarget() *RuntimeATNState { return t.Target }
func (t *RuntimeAtomTransition) IsEpsilon() bool             { return false }

// RuntimeEpsilonTransition fires without consuming a token.
type RuntimeEpsilonTransition struct {
	Target *RuntimeATNState
}

func (t *RuntimeEpsilonTransition) GetTarget() *RuntimeATNState { return t.Target }
func (t *RuntimeEpsilonTransition) IsEpsilon() bool             { return true }

// RuntimeRuleTransition enters a sub-rule and returns to FollowState.
type RuntimeRuleTransition struct {
	Target      *RuntimeATNState // the rule's RuleStartState
	FollowState *RuntimeATNState
}

func (t *RuntimeRuleTransition) GetTarget() *RuntimeATNState { return t.Target }
func (t *RuntimeRuleTransition) IsEpsilon() bool             { return true }

// RuntimeATN is the minimal ATN structure used at prediction time.
// It can be built directly from generated Go code without running CreateATN.
type RuntimeATN struct {
	States         []*RuntimeATNState
	DecisionStates []*RuntimeATNState          // indexed by Decision
	DecisionMap    map[string]*RuntimeATNState // key → decision state
}
