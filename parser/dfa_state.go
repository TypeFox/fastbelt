// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	core "typefox.dev/fastbelt"
)

// dfaState is a node in a decision's DFA. Each state wraps the ATNConfigSet
// reachable on the input consumed so far.
//
// edges is indexed by tokenType.Id + 1, so the EOF symbol (t == -1) maps to slot 0
type dfaState struct {
	stateNumber         int
	configs             *atnConfigSet
	edges               []*dfaState
	isAcceptState       bool
	prediction          int
	requiresFullContext bool
}

func newDFAState(stateNumber int, configs *atnConfigSet) *dfaState {
	return &dfaState{stateNumber: stateNumber, configs: configs, prediction: invalidAlt}
}

// errorDFAState is the shared sentinel meaning "no viable target".
var errorDFAState = newDFAState(0x7FFFFFFF, newATNConfigSet(false))

func (d *dfaState) getEdge(t *core.TokenType) *dfaState {
	i := t.Id + 1
	if i < 0 || i >= len(d.edges) {
		return nil
	}
	return d.edges[i]
}

func (d *dfaState) setEdge(t *core.TokenType, target *dfaState) {
	i := t.Id + 1
	for i >= len(d.edges) {
		d.edges = append(d.edges, nil)
	}
	d.edges[i] = target
}

// Equals for DFA-state deduplication compares the underlying config sets.
func (d *dfaState) Equals(o *dfaState) bool {
	if d == o {
		return true
	}
	if o == nil {
		return false
	}
	return d.configs.Equals(o.configs)
}

func (d *dfaState) Hash() uint64 { return d.configs.Hash() }
