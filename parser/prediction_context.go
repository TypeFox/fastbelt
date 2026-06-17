// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import "strconv"

// predictionContext is a graph-structured call stack used by the adaptive predictor.
//
// A context is one of three shapes:
//   - empty ("$"): nobody called us; the bottom of every stack.
//   - singleton: a single (parent, returnState) frame.
//   - array: a set of (parent, returnState) frames sharing a node, kept sorted
//     by returnState with the empty return state ($) last.
type predictionContext struct {
	cachedHash uint64

	parents      []*predictionContext
	returnStates []int
}

// predictionContextEmptyReturnState marks the "$" return state inside an array
// context (full-context mode). It can never collide with a real ATN state
// index.
const predictionContextEmptyReturnState = 0x7FFFFFFF

// basePredictionContextEMPTY is the shared empty context ("$").
var basePredictionContextEMPTY = &predictionContext{
	returnStates: []int{predictionContextEmptyReturnState},
	cachedHash:   emptyContextHash(),
}

func emptyContextHash() uint64 { return hashCombine(1, 0) }

func emptyPredictionContext() *predictionContext { return basePredictionContextEMPTY }

func singletonPredictionContext(parent *predictionContext, returnState int) *predictionContext {
	if returnState == predictionContextEmptyReturnState && parent == nil {
		return basePredictionContextEMPTY
	}
	pc := &predictionContext{
		parents:      []*predictionContext{parent},
		returnStates: []int{returnState},
	}
	if parent != nil {
		pc.cachedHash = hashCombine(parent.cachedHash, uint64(returnState))
	} else {
		pc.cachedHash = emptyContextHash()
	}
	return pc
}

func newArrayPredictionContext(parents []*predictionContext, returnStates []int) *predictionContext {
	h := uint64(1)
	for _, p := range parents {
		var ph uint64
		if p != nil {
			ph = p.cachedHash
		}
		h = hashCombine(h, ph)
	}
	for _, rs := range returnStates {
		h = hashCombine(h, uint64(rs))
	}
	return &predictionContext{
		parents:      parents,
		returnStates: returnStates,
		cachedHash:   hashCombine(h, uint64(len(parents))<<1),
	}
}

func (p *predictionContext) length() int {
	return len(p.returnStates)
}

func (p *predictionContext) getParent(i int) *predictionContext {
	if i < 0 || i >= len(p.parents) {
		return nil
	}
	return p.parents[i]
}

func (p *predictionContext) getReturnState(i int) int {
	if i < 0 || i >= len(p.returnStates) {
		return predictionContextEmptyReturnState
	}
	return p.returnStates[i]
}

func (p *predictionContext) isEmpty() bool {
	return p.getReturnState(0) == predictionContextEmptyReturnState
}

func (p *predictionContext) hasEmptyPath() bool {
	return p.getReturnState(p.length()-1) == predictionContextEmptyReturnState
}

func (p *predictionContext) Hash() uint64 { return p.cachedHash }

func (p *predictionContext) Equals(other *predictionContext) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.cachedHash != other.cachedHash {
		return false
	}
	if len(p.returnStates) != len(other.returnStates) {
		return false
	}
	for i := range p.returnStates {
		if p.returnStates[i] != other.returnStates[i] {
			return false
		}
		a, b := p.parents[i], other.parents[i]
		if a == nil || b == nil {
			if a != b {
				return false
			}
			continue
		}
		if !a.Equals(b) {
			return false
		}
	}
	return true
}

func (p *predictionContext) String() string {
	s := "["
	for i := range p.returnStates {
		if i > 0 {
			s += ", "
		}
		if p.returnStates[i] == predictionContextEmptyReturnState {
			s += "$"
			continue
		}
		s += strconv.Itoa(p.returnStates[i])
	}
	return s + "]"
}

// mergeCache memoizes prediction-context merges within a single prediction.
type mergeCache struct {
	m map[[2]*predictionContext]*predictionContext
}

func newMergeCache() *mergeCache {
	return &mergeCache{m: map[[2]*predictionContext]*predictionContext{}}
}

func (c *mergeCache) get(a, b *predictionContext) (*predictionContext, bool) {
	if c == nil {
		return nil, false
	}
	v, ok := c.m[[2]*predictionContext{a, b}]
	return v, ok
}

func (c *mergeCache) put(a, b, v *predictionContext) {
	if c == nil {
		return
	}
	c.m[[2]*predictionContext{a, b}] = v
}

func mergePredictionContexts(a, b *predictionContext, rootIsWildcard bool, cache *mergeCache) *predictionContext {
	if a == b || a.Equals(b) {
		return a
	}
	if rootIsWildcard {
		if a.isEmpty() {
			return a
		}
		if b.isEmpty() {
			return b
		}
	}
	return mergeArrays(a, b, rootIsWildcard, cache)
}

func mergeArrays(a, b *predictionContext, rootIsWildcard bool, cache *mergeCache) *predictionContext {
	if v, ok := cache.get(a, b); ok {
		return v
	}
	if v, ok := cache.get(b, a); ok {
		return v
	}

	i, j, k := 0, 0, 0
	mergedReturnStates := make([]int, len(a.returnStates)+len(b.returnStates))
	mergedParents := make([]*predictionContext, len(a.returnStates)+len(b.returnStates))

	for i < len(a.returnStates) && j < len(b.returnStates) {
		aParent, bParent := a.parents[i], b.parents[j]
		switch {
		case a.returnStates[i] == b.returnStates[j]:
			payload := a.returnStates[i]
			bothDollars := payload == predictionContextEmptyReturnState && aParent == nil && bParent == nil
			axAX := aParent != nil && bParent != nil && aParent.Equals(bParent)
			if bothDollars || axAX {
				mergedParents[k] = aParent
			} else {
				mergedParents[k] = mergePredictionContexts(aParent, bParent, rootIsWildcard, cache)
			}
			mergedReturnStates[k] = payload
			i++
			j++
		case a.returnStates[i] < b.returnStates[j]:
			mergedParents[k] = aParent
			mergedReturnStates[k] = a.returnStates[i]
			i++
		default:
			mergedParents[k] = bParent
			mergedReturnStates[k] = b.returnStates[j]
			j++
		}
		k++
	}
	if i < len(a.returnStates) {
		for p := i; p < len(a.returnStates); p++ {
			mergedParents[k] = a.parents[p]
			mergedReturnStates[k] = a.returnStates[p]
			k++
		}
	} else {
		for p := j; p < len(b.returnStates); p++ {
			mergedParents[k] = b.parents[p]
			mergedReturnStates[k] = b.returnStates[p]
			k++
		}
	}

	if k < len(mergedParents) {
		if k == 1 {
			pc := singletonPredictionContext(mergedParents[0], mergedReturnStates[0])
			cache.put(a, b, pc)
			return pc
		}
		mergedParents = mergedParents[:k]
		mergedReturnStates = mergedReturnStates[:k]
	}

	m := newArrayPredictionContext(mergedParents, mergedReturnStates)
	if m.Equals(a) {
		cache.put(a, b, a)
		return a
	}
	if m.Equals(b) {
		cache.put(a, b, b)
		return b
	}
	cache.put(a, b, m)
	return m
}

// hashCombine mixes a running hash with a 64-bit value. The exact algorithm is
// unimportant; it only needs to be stable and reasonably well-distributed,
// since equality (not hashing) decides context identity.
func hashCombine(h, v uint64) uint64 {
	h ^= v + 0x9e3779b97f4a7c15 + (h << 6) + (h >> 2)
	return h
}
