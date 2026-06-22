// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"slices"

	"typefox.dev/fastbelt/util/collections"
)

// invalidAlt is the "no unique alternative" sentinel. Alternatives are 0-based,
// so -1 can never be a real alt.
const invalidAlt = -1

// configKey identifies configs that should be merged on Add: same state and
// alternative, independent of context (ANTLR keys on (s, i, semctx); Fastbelt
// has no semantic context, so the key is (state, alt)).
type configKey struct {
	state *RuntimeATNState
	alt   int
}

// atnConfigSet is a set of ATNConfig that merges configs sharing (state, alt)
// by combining their prediction contexts, and tracks the bookkeeping the
// predictor needs (uniqueAlt, conflictingAlts, dipsIntoOuterContext, fullCtx).
type atnConfigSet struct {
	configs []*atnConfig
	lookup  map[configKey]*atnConfig

	fullCtx              bool
	uniqueAlt            int
	conflictingAlts      *collections.BitSet
	dipsIntoOuterContext bool
}

func newATNConfigSet(fullCtx bool) *atnConfigSet {
	return &atnConfigSet{
		lookup:    map[configKey]*atnConfig{},
		fullCtx:   fullCtx,
		uniqueAlt: invalidAlt,
	}
}

// Add merges config into the set, combining contexts with any existing config
// sharing (state, alt).
func (s *atnConfigSet) Add(config *atnConfig, cache *mergeCache) {
	if config.reachesIntoOuterContext > 0 {
		s.dipsIntoOuterContext = true
	}
	key := configKey{state: config.state, alt: config.alt}
	existing, present := s.lookup[key]
	if !present {
		s.lookup[key] = config
		s.configs = append(s.configs, config)
		return
	}
	rootIsWildcard := !s.fullCtx
	merged := mergePredictionContexts(existing.context, config.context, rootIsWildcard, cache)
	if config.reachesIntoOuterContext > existing.reachesIntoOuterContext {
		existing.reachesIntoOuterContext = config.reachesIntoOuterContext
	}
	existing.context = merged
}

func (s *atnConfigSet) isEmpty() bool { return len(s.configs) == 0 }

// Compare requires same order and element-wise Equals (matches ANTLR).
func (s *atnConfigSet) compare(o *atnConfigSet) bool {
	if len(s.configs) != len(o.configs) {
		return false
	}
	for i := range s.configs {
		if !s.configs[i].Equals(o.configs[i]) {
			return false
		}
	}
	return true
}

func (s *atnConfigSet) Equals(o *atnConfigSet) bool {
	if s == o {
		return true
	}
	if o == nil {
		return false
	}
	return s.fullCtx == o.fullCtx &&
		s.uniqueAlt == o.uniqueAlt &&
		s.dipsIntoOuterContext == o.dipsIntoOuterContext &&
		s.conflictingAlts.Equals(o.conflictingAlts) &&
		s.compare(o)
}

func (s *atnConfigSet) Hash() uint64 {
	h := uint64(1)
	for _, c := range s.configs {
		h = 31*h + c.Hash()
	}
	return h
}

// isRuleStop reports whether the config sits in a rule-stop state.
func isRuleStop(c *atnConfig) bool { return c.state.Type == ATNRuleStop }

// allConfigsInRuleStopStates reports whether every config has reached a rule
// stop state (no config can match further input).
func allConfigsInRuleStopStates(configs *atnConfigSet) bool {
	for _, c := range configs.configs {
		if !isRuleStop(c) {
			return false
		}
	}
	return true
}

// hasConfigInRuleStopState reports whether any config has reached a rule stop.
func hasConfigInRuleStopState(configs *atnConfigSet) bool {
	return slices.ContainsFunc(configs.configs, isRuleStop)
}

// hasConflictTerminatingPrediction reports whether SLL prediction must stop
// at this config set because of an unresolvable conflict.
func hasConflictTerminatingPrediction(configs *atnConfigSet) bool {
	if allConfigsInRuleStopStates(configs) {
		return true
	}
	altsets := getConflictingAltSubsets(configs)
	return hasConflictingAltSet(altsets) && !hasStateAssociatedWithOneAlt(configs)
}

type bucketEntry struct {
	state   *RuntimeATNState
	context *predictionContext
}

func (e *bucketEntry) Hash() uint64 {
	hash := uint64(e.state.StateNumber)
	if e.context != nil {
		hash = hashCombine(hash, e.context.Hash())
	}
	return hash
}

func (e *bucketEntry) Equals(other *bucketEntry) bool {
	if e.state.StateNumber != other.state.StateNumber {
		return false
	}
	if e.context == nil || other.context == nil {
		return e.context == other.context
	}
	return e.context.Equals(other.context)
}

// getConflictingAltSubsets groups configs by (state, context), collecting the
// set of alternatives for each group. Two configs conflict when they share a
// state and context but predict different alternatives.
func getConflictingAltSubsets(configs *atnConfigSet) []*collections.BitSet {
	buckets := collections.NewBucketMap[*bucketEntry, *collections.BitSet]()
	var order []*collections.BitSet
	for _, c := range configs.configs {
		entry := &bucketEntry{state: c.state, context: c.context}
		found, exists := buckets.Get(entry)
		if !exists {
			found = collections.NewBitset()
			buckets.Set(entry, found)
			order = append(order, found)
		}
		found.Insert(c.alt)
	}
	return order
}

// getSingleViableAlt returns the alternative all subsets resolve to (by minimum
// of each subset), or invalidAlt if they differ. Also used as
// resolvesToJustOneViableAlt for full LL.
func getSingleViableAlt(altsets []*collections.BitSet) int {
	result := invalidAlt
	for _, alts := range altsets {
		minAlt := alts.Min()
		if result == invalidAlt {
			result = minAlt
		} else if result != minAlt {
			return invalidAlt
		}
	}
	return result
}

func hasConflictingAltSet(altsets []*collections.BitSet) bool {
	for _, alts := range altsets {
		if alts.Cardinality() > 1 {
			return true
		}
	}
	return false
}

// hasStateAssociatedWithOneAlt reports whether some ATN state appears with
// exactly one alternative across the config set.
func hasStateAssociatedWithOneAlt(configs *atnConfigSet) bool {
	stateToAlts := map[int]*collections.BitSet{}
	for _, c := range configs.configs {
		alts := stateToAlts[c.state.StateNumber]
		if alts == nil {
			alts = collections.NewBitset()
			stateToAlts[c.state.StateNumber] = alts
		}
		alts.Insert(c.alt)
	}
	for _, alts := range stateToAlts {
		if alts.Cardinality() == 1 {
			return true
		}
	}
	return false
}
