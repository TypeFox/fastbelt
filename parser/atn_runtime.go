// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"sync"

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

// RuntimeAtomTransition fires on a specific token type.
type RuntimeAtomTransition struct {
	Target    *RuntimeATNState
	TokenType *fastbelt.TokenType
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
	DecisionStates []*RuntimeATNState // indexed by Decision
	DecisionMap    []*RuntimeATNState // key → decision state

	stateIdxOnce  sync.Once
	stateIdxCache map[*RuntimeATNState]int // lazily built pointer → array index

	nextTokensOnce  sync.Once
	nextTokensCache [][]bool // stateIdx → bitset indexed by TokenType.Id
	tokenSetSize    int      // common length of every entry in nextTokensCache
}

// stateIndex returns the array index of s in States, building a cache on first call.
func (atn *RuntimeATN) stateIndex(s *RuntimeATNState) int {
	atn.stateIdxOnce.Do(func() {
		atn.stateIdxCache = make(map[*RuntimeATNState]int, len(atn.States))
		for i, st := range atn.States {
			atn.stateIdxCache[st] = i
		}
	})
	if i, ok := atn.stateIdxCache[s]; ok {
		return i
	}
	return -1
}

// TokenSetSize returns the length of any slice returned by NextTokensAt.
// Callers can use it to allocate compatible token bitsets.
func (atn *RuntimeATN) TokenSetSize() int {
	atn.nextTokensOnce.Do(atn.buildNextTokensCache)
	return atn.tokenSetSize
}

// NextTokensAt returns the set of token type IDs reachable from the state at
// stateIdx via epsilon closure (including rule entries for FIRST-set tokens).
// The returned slice is indexed by TokenType.Id; it is shared and must not be
// mutated by callers. Returns nil if stateIdx is out of bounds.
func (atn *RuntimeATN) NextTokensAt(stateIdx int) []bool {
	atn.nextTokensOnce.Do(atn.buildNextTokensCache)
	if stateIdx < 0 || stateIdx >= len(atn.nextTokensCache) {
		return nil
	}
	return atn.nextTokensCache[stateIdx]
}

func (atn *RuntimeATN) buildNextTokensCache() {
	maxId := 0
	for _, st := range atn.States {
		if st == nil {
			continue
		}
		for _, t := range st.Transitions {
			if at, ok := t.(*RuntimeAtomTransition); ok && at.TokenType != nil {
				if at.TokenType.Id > maxId {
					maxId = at.TokenType.Id
				}
			}
		}
	}
	atn.tokenSetSize = maxId + 1
	atn.nextTokensCache = make([][]bool, len(atn.States))
	for i := range atn.States {
		atn.nextTokensCache[i] = atn.computeNextTokensAt(i)
	}
}

func (atn *RuntimeATN) computeNextTokensAt(stateIdx int) []bool {
	result := make([]bool, atn.tokenSetSize)
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
					result[at.TokenType.Id] = true
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
	return result
}
