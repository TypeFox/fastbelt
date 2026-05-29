// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"typefox.dev/fastbelt"
)

// ATNStateType is the discriminator for ATN state kinds.
type ATNStateType int

const (
	ATNInvalidType ATNStateType = iota
	ATNBasic
	ATNRuleStart
	ATNPlusBlockStart
	ATNStarBlockStart
	ATNTokenStart
	ATNRuleStop
	ATNBlockEnd
	ATNStarLoopBack
	ATNStarLoopEntry
	ATNPlusLoopBack
	ATNLoopEnd
)

// RuntimeATNState holds only the fields required for prediction at runtime.
// Back-pointers to the build-time grammar objects (Rule, Production) and
// structural scaffolding (End, Loopback, Start, Stop) are absent.
type RuntimeATNState struct {
	StateNumber            int
	Type                   ATNStateType
	Decision               int
	EpsilonOnlyTransitions bool
	Transitions            []RuntimeTransition
}

// RuntimeTransition is the interface implemented by all runtime ATN transitions.
type RuntimeTransition interface {
	GetTarget() *RuntimeATNState
	IsEpsilon() bool
}

// CompletionHint marks an atom transition that originated from a cross-reference
// assignment in the grammar. The Field key (e.g. "Transition.Event") is used by
// the completion provider to dispatch to the right per-field scope/filter on the
// generated language CompletionProvider. The same key seeds the synthetic-owner
// chain when the cursor's container does not yet exist.
//
// PrecedingAction is non-nil when this atom is the first token-consuming
// element after a grammar Action (tree-rewrite) in its enclosing group. The
// completion provider uses it to wrap the existing AST node in a freshly
// allocated parent when the existing node's assignment slot is already filled
// - mirroring Langium's `NextFeature.type`/`property` synthesis.
type CompletionHint struct {
	Field           string
	PrecedingAction *ActionInfo
}

// ActionInfo describes a grammar Action that fires immediately before the
// atom carrying its CompletionHint. The completion provider applies it by
// allocating a new node of TargetType and assigning the existing AST node
// to the named Property. The choice of single vs append assignment is
// baked into the generated adapter from the field's grammar type, so the
// operator does not need to travel through the ATN.
type ActionInfo struct {
	TargetType string // e.g. "MemberCall"
	Property   string // e.g. "Previous"
}

// RuntimeAtomTransition fires on a specific token type.
//
// CompletionHint is non-nil when this transition's token was emitted by a
// cross-reference assignment; it lets the completion engine know that "match
// Token_ID here" means "complete a reference to the field named by Hint.Field".
type RuntimeAtomTransition struct {
	Target         *RuntimeATNState
	TokenType      *fastbelt.TokenType
	CompletionHint *CompletionHint
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
//
// CompletionHint is non-nil when this rule call was emitted by a
// cross-reference assignment whose text-form is a rule (e.g.
// `Ref=[Decl:FQN]` where FQN is itself a parser rule). Every atom
// transition reached inside the called rule then represents one token of
// the cross-reference's text, so the simulator propagates this hint
// onto its live-set paths until the matching RuleStop pops it off again.
type RuntimeRuleTransition struct {
	Target         *RuntimeATNState // the rule's RuleStartState
	FollowState    *RuntimeATNState
	CompletionHint *CompletionHint
}

func (t *RuntimeRuleTransition) GetTarget() *RuntimeATNState { return t.Target }
func (t *RuntimeRuleTransition) IsEpsilon() bool             { return true }

// RuntimeATN is the minimal ATN structure used at prediction time.
// Call [NewRuntimeATN] to construct from the full build-time ATN to ensure proper initialization.
type RuntimeATN struct {
	States         []*RuntimeATNState
	DecisionStates []*RuntimeATNState // indexed by Decision
	DecisionMap    []*RuntimeATNState // key -> decision state

	stateIdxCache map[*RuntimeATNState]int // pointer -> array index

	nextTokensCache []*fastbelt.BitSet // stateIdx -> bitset indexed by TokenType.Id
}

func NewRuntimeATN(states []*RuntimeATNState, decisionStates []*RuntimeATNState, decisionMap []*RuntimeATNState) *RuntimeATN {
	atn := &RuntimeATN{
		States:         states,
		DecisionStates: decisionStates,
		DecisionMap:    decisionMap,
	}
	atn.Init()
	return atn
}

// Initializes the different caches of the ATN.
func (atn *RuntimeATN) Init() {
	atn.buildIdxCache()
	atn.buildNextTokensCache()
}

func (atn *RuntimeATN) buildIdxCache() {
	atn.stateIdxCache = make(map[*RuntimeATNState]int, len(atn.States))
	for i, st := range atn.States {
		atn.stateIdxCache[st] = i
	}
}

// stateIndex returns the array index of s in States, building a cache on first call.
func (atn *RuntimeATN) stateIndex(s *RuntimeATNState) int {
	if i, ok := atn.stateIdxCache[s]; ok {
		return i
	}
	return -1
}

// NextTokensAt returns the set of token type IDs reachable from the state at
// stateIdx via epsilon closure (including rule entries for FIRST-set tokens).
func (atn *RuntimeATN) NextTokensAt(stateIdx int) *fastbelt.BitSet {
	if stateIdx < 0 || stateIdx >= len(atn.nextTokensCache) {
		return fastbelt.NewBitset()
	}
	return atn.nextTokensCache[stateIdx]
}

func (atn *RuntimeATN) buildNextTokensCache() {
	atn.nextTokensCache = make([]*fastbelt.BitSet, len(atn.States))
	for i := range atn.States {
		atn.nextTokensCache[i] = atn.computeNextTokensAt(i)
	}
}

func (atn *RuntimeATN) computeNextTokensAt(stateIdx int) *fastbelt.BitSet {
	var sets []*fastbelt.BitSet
	visited := make([]bool, len(atn.States))
	queue := []int{stateIdx}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if visited[cur] {
			continue
		}
		visited[cur] = true
		state := atn.States[cur]
		if state == nil {
			continue
		}
		for _, t := range state.Transitions {
			switch at := t.(type) {
			case *RuntimeAtomTransition:
				if at.TokenType != nil {
					sets = append(sets, at.TokenType.Bitset())
				}
			case *RuntimeEpsilonTransition:
				if tidx := atn.stateIndex(at.Target); tidx >= 0 {
					queue = append(queue, tidx)
				}
			case *RuntimeRuleTransition:
				// Descend into rule to collect FIRST tokens; do not enqueue FollowState.
				if tidx := atn.stateIndex(at.Target); tidx >= 0 {
					queue = append(queue, tidx)
				}
			}
		}
	}
	return fastbelt.MergeBitSets(sets)
}
