package generator

import "typefox.dev/fastbelt/internal/grammar"

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

// ATNState is the single concrete ATN state type.
// Fields specific to certain state kinds are non-nil only for those kinds.
type ATNState struct {
	ATN                    *ATN
	Production             grammar.Element // nil for rule start/stop
	StateNumber            int
	Rule                   *grammar.ParserRule
	EpsilonOnlyTransitions bool
	Transitions            []Transition
	Type                   ATNStateType

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
	RuleToStartState map[*grammar.ParserRule]*ATNState
	RuleToStopState  map[*grammar.ParserRule]*ATNState
}

// Transition is the interface implemented by all ATN transitions.
type Transition interface {
	Target() *ATNState
	IsEpsilon() bool
}

// AtomTransition fires on a specific token type.
// CategoryMatches holds the IDs of all token types that match via category
// inheritance; populated from the Terminal's CategoryMatches at ATN-build time.
type AtomTransition struct {
	TargetState     *ATNState
	TokenTypeID     int
	CategoryMatches []int
}

func (t *AtomTransition) Target() *ATNState { return t.TargetState }
func (t *AtomTransition) IsEpsilon() bool   { return false }

// EpsilonTransition fires without consuming a token.
type EpsilonTransition struct {
	TargetState *ATNState
}

func (t *EpsilonTransition) Target() *ATNState { return t.TargetState }
func (t *EpsilonTransition) IsEpsilon() bool   { return true }

// RuleTransition enters a sub-rule and returns to FollowState.
type RuleTransition struct {
	TargetState *ATNState // the rule's RuleStartState
	Rule        *grammar.ParserRule
	FollowState *ATNState
}

func (t *RuleTransition) Target() *ATNState { return t.TargetState }
func (t *RuleTransition) IsEpsilon() bool   { return true }

// ATNHandle is an internal pair of (entry, exit) ATN states for a sub-network.
type ATNHandle struct {
	Left  *ATNState
	Right *ATNState
}

// TokenInfo carries the type ID and category-match IDs for a token type.
type TokenInfo struct {
	ID              int
	CategoryMatches []int
}
