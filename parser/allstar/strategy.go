// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import "fmt"

// Strategy is the prediction interface shared by LL(k) and ALL(*) algorithms.
// Implementations are chosen at parser-construction time to swap prediction
// behaviour without changing parsing logic.
//
// PredictAlternation returns the chosen alternative index (0-based) for the
// OR decision identified by key, or -1 when no prediction is possible.
//
// PredictOptional returns true when the optional block identified by key
// should be entered (alt 0 in the underlying decision state).
type Strategy interface {
	PredictAlternation(src TokenSource, key string) int
	PredictOptional(src TokenSource, key string) bool
}

// LLStarLookaheadOptions configures the strategy.
type LLStarLookaheadOptions struct {
	// Logging is called whenever an ambiguity is detected.
	// Defaults to a function that prints to stdout.
	Logging AmbiguityReport
}

// LLStarLookahead is the main entry point. It holds the RuntimeATN and the
// DFA caches produced during initialisation.
type LLStarLookahead struct {
	atn     *RuntimeATN
	dfas    []dfaCache
	logging AmbiguityReport
}

// NewLLStarLookahead creates a new strategy, builds the ATN from rules, converts
// it to the minimal RuntimeATN, and initialises the DFA cache array.
// The build-time ATN is discarded after conversion.
func NewLLStarLookahead(rules []*Rule, opts *LLStarLookaheadOptions) *LLStarLookahead {
	rtn := BuildRuntimeATN(CreateATN(rules))
	return newFromRuntime(rtn, opts)
}

// NewLLStarLookaheadFromRuntime creates a strategy from a pre-built RuntimeATN.
// Use this with ATNs constructed directly from generated Go code to skip the
// build-time grammar processing entirely.
func NewLLStarLookaheadFromRuntime(rtn *RuntimeATN, opts *LLStarLookaheadOptions) *LLStarLookahead {
	return newFromRuntime(rtn, opts)
}

func newFromRuntime(rtn *RuntimeATN, opts *LLStarLookaheadOptions) *LLStarLookahead {
	logging := AmbiguityReport(func(msg string) { fmt.Println(msg) })
	if opts != nil && opts.Logging != nil {
		logging = opts.Logging
	}
	return &LLStarLookahead{
		atn:     rtn,
		dfas:    initDFACaches(rtn),
		logging: logging,
	}
}

// AdaptivePredict runs the ALL(*) algorithm.
// Returns the chosen alternative index (0-based), or -1 on parse error.
func (s *LLStarLookahead) AdaptivePredict(src TokenSource, decision int, predicates *PredicateSet) int {
	if predicates == nil {
		predicates = EmptyPredicates
	}
	alt, _ := adaptivePredict(src, s.dfas, decision, predicates, s.logging)
	return alt
}

// PredictAlternation implements Strategy.
// It resolves the decision state by key and delegates to AdaptivePredict.
func (s *LLStarLookahead) PredictAlternation(src TokenSource, key string) int {
	ds := s.atn.DecisionMap[key]
	if ds == nil {
		return -1
	}
	return s.AdaptivePredict(src, ds.Decision, EmptyPredicates)
}

// PredictOptional implements Strategy.
// Returns true when the adaptive prediction chooses alt 0 (enter the block).
//
// Fast path: when the grammar is LL(1)-decidable at this decision point (no
// token overlap between the enter and skip/exit alternatives), the LL(1) table
// gives the answer in O(1).  This is critical for StarLoopEntry decisions where
// the exit alternative has no first tokens (it only fires at EOF), because
// adaptivePredict would otherwise scan the entire remaining input before
// returning the exit alternative.
func (s *LLStarLookahead) PredictOptional(src TokenSource, key string) bool {
	ds := s.atn.DecisionMap[key]
	if ds == nil {
		return false
	}
	if table, ok := buildLL1Table(ds); ok {
		tID := tokenTypeID(src.LA(1))
		alt, found := table[tID]
		if !found {
			return false
		}
		return alt == 0
	}
	alt, err := adaptivePredict(src, s.dfas, ds.Decision, EmptyPredicates, s.logging)
	return err == nil && alt == 0
}

// BuildLookaheadForAlternation returns a function that wraps AdaptivePredict
// for a specific alternation occurrence in a rule. Falls back to an LL(1)
// table when the grammar is deterministic at depth 1.
func (s *LLStarLookahead) BuildLookaheadForAlternation(
	rule *Rule,
	occurrence int,
	hasPredicates bool,
) func(src TokenSource, gates []func() bool) int {
	key := BuildATNKey(rule, "Alternation", occurrence)
	decisionState := s.atn.DecisionMap[key]
	if decisionState == nil {
		return func(_ TokenSource, _ []func() bool) int { return -1 }
	}
	decision := decisionState.Decision
	dfas := s.dfas
	logging := s.logging

	// Attempt LL(1) fast path: inspect the first token reachable from each
	// alternative's start state. If every token maps to exactly one alt, we
	// can skip the full ATN simulation.
	if !hasPredicates {
		if table, ok := buildLL1Table(decisionState); ok {
			return func(src TokenSource, _ []func() bool) int {
				tID := tokenTypeID(src.LA(1))
				alt, found := table[tID]
				if !found {
					return -1
				}
				return alt
			}
		}
	}

	if hasPredicates {
		return func(src TokenSource, gates []func() bool) int {
			preds := &PredicateSet{}
			for i, gate := range gates {
				if gate != nil {
					preds.Set(i, gate())
				}
			}
			alt, _ := adaptivePredict(src, dfas, decision, preds, logging)
			return alt
		}
	}
	return func(src TokenSource, _ []func() bool) int {
		alt, _ := adaptivePredict(src, dfas, decision, EmptyPredicates, logging)
		return alt
	}
}

// BuildLookaheadForOptional returns a function that wraps AdaptivePredict
// for OPTION, MANY, and AT_LEAST_ONE productions.
func (s *LLStarLookahead) BuildLookaheadForOptional(
	rule *Rule,
	occurrence int,
	prodType string,
) func(src TokenSource) bool {
	key := BuildATNKey(rule, prodType, occurrence)
	decisionState := s.atn.DecisionMap[key]
	if decisionState == nil {
		return func(_ TokenSource) bool { return false }
	}
	decision := decisionState.Decision
	dfas := s.dfas
	logging := s.logging

	// Fast path: same LL(1) optimisation as PredictOptional.
	if table, ok := buildLL1Table(decisionState); ok {
		return func(src TokenSource) bool {
			tID := tokenTypeID(src.LA(1))
			alt, found := table[tID]
			return found && alt == 0
		}
	}

	return func(src TokenSource) bool {
		alt, err := adaptivePredict(src, dfas, decision, EmptyPredicates, logging)
		if err != nil {
			return false
		}
		return alt == 0
	}
}

// buildLL1Table attempts to build a token-type-to-alt map from the decision
// state. Returns (table, true) when the grammar is LL(1) at this point.
// The algorithm inspects the immediate atom transitions reachable via epsilon
// from each alternative.
func buildLL1Table(decision *RuntimeATNState) (map[int]int, bool) {
	table := map[int]int{}
	for i, t := range decision.Transitions {
		if !t.IsEpsilon() {
			continue
		}
		tokens := firstTokens(t.GetTarget(), 0)
		for _, tID := range tokens {
			if _, conflict := table[tID]; conflict {
				return nil, false
			}
			table[tID] = i
		}
	}
	return table, len(table) > 0
}

// firstTokens returns the set of token type IDs immediately reachable (via
// epsilon) from state, limited to a small depth to avoid infinite loops.
func firstTokens(state *RuntimeATNState, depth int) []int {
	if depth > 8 {
		return nil
	}
	var result []int
	for _, t := range state.Transitions {
		switch at := t.(type) {
		case *RuntimeAtomTransition:
			result = append(result, at.TokenTypeID)
			for cat := range at.CategoryMatches {
				result = append(result, cat)
			}
		case *RuntimeEpsilonTransition:
			result = append(result, firstTokens(at.Target, depth+1)...)
		case *RuntimeRuleTransition:
			result = append(result, firstTokens(at.Target, depth+1)...)
		}
	}
	return result
}
