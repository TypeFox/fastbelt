// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

// dfa is the per-decision deterministic-automaton cache built lazily during
// adaptive prediction.
//
// The dfa is shared across concurrent parses (it hangs off the shared
// RuntimeATN), so all mutation of start, and edges is guarded by mu.
type dfa struct {
	decision      int
	atnStartState *RuntimeATNState

	mu     sync.RWMutex
	start  *dfaState
	states *collections.BucketMap[*dfaState, *dfaState]
}

func newDFA(decision int, atnStartState *RuntimeATNState) *dfa {
	return &dfa{
		decision:      decision,
		atnStartState: atnStartState,
		states:        collections.NewBucketMap[*dfaState, *dfaState](),
	}
}

func (d *dfa) getStart() *dfaState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.start
}

func (d *dfa) setStart(s *dfaState) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.start = s
}

// addState deduplicates d against the existing states and returns the canonical
// instance. The error sentinel is never cached. Caller must not hold d.mu.
func (d *dfa) addState(s *dfaState) *dfaState {
	if s == errorDFAState {
		return s
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	// Upsert the state into the bucket map, which returns the existing state if present.
	result, _ := d.states.GetOrInsert(s, s)
	return result
}

// getExistingEdge returns the cached target for (from, t), or nil.
func (d *dfa) getExistingEdge(from *dfaState, t *core.TokenType) *dfaState {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return from.getEdge(t)
}

// addEdge canonicalizes to and links from --t--> to, returning the canonical to.
func (d *dfa) addEdge(from *dfaState, t *core.TokenType, to *dfaState) *dfaState {
	to = d.addState(to)
	d.mu.Lock()
	defer d.mu.Unlock()
	from.setEdge(t, to)
	return to
}
