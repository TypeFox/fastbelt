// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"fmt"
	"strings"
)

// DFA is the memoization automaton for a single decision.
type DFA struct {
	Start         *DFAState
	States        map[string]*DFAState
	Decision      int
	ATNStartState *ATNState
}

// DFAState represents a set of ATN configurations.
type DFAState struct {
	Configs       *ATNConfigSet
	Edges         map[int]*DFAState
	IsAcceptState bool
	Prediction    int // valid when IsAcceptState == true
}

// DFAError is a sentinel value returned for dead-end states.
var DFAError = &DFAState{}

// ATNConfig is one thread of the ATN simulation.
type ATNConfig struct {
	State *ATNState
	Alt   int
	Stack []*ATNState
}

// ATNConfigSet holds a deduplicated set of ATN configurations.
type ATNConfigSet struct {
	configMap map[string]int
	configs   []*ATNConfig
	// UniqueAlt is set when all configs share the same alternative; -1 otherwise.
	UniqueAlt int
}

// NewATNConfigSet creates a new, empty ATNConfigSet.
func NewATNConfigSet() *ATNConfigSet {
	return &ATNConfigSet{
		configMap: map[string]int{},
		UniqueAlt: -1,
	}
}

// Add inserts c into the set if no equivalent config is already present.
func (s *ATNConfigSet) Add(c *ATNConfig) {
	key := atnConfigKey(c, true)
	if _, exists := s.configMap[key]; !exists {
		s.configMap[key] = len(s.configs)
		s.configs = append(s.configs, c)
	}
}

// Finalize releases the deduplication map to free memory.
func (s *ATNConfigSet) Finalize() {
	s.configMap = map[string]int{}
}

// Len returns the number of configs in the set.
func (s *ATNConfigSet) Len() int {
	return len(s.configs)
}

// Elements returns the configs in insertion order.
func (s *ATNConfigSet) Elements() []*ATNConfig {
	return s.configs
}

// Alts returns a slice of alternative indices, one per config.
func (s *ATNConfigSet) Alts() []int {
	alts := make([]int, len(s.configs))
	for i, c := range s.configs {
		alts[i] = c.Alt
	}
	return alts
}

// Key returns a string that uniquely identifies the set of configs.
func (s *ATNConfigSet) Key() string {
	var b strings.Builder
	for k := range s.configMap {
		b.WriteString(k)
		b.WriteByte(':')
	}
	return b.String()
}

// atnConfigKey produces the deduplication key for a config.
// When includeAlt is false the alt index is omitted (used for conflict detection).
func atnConfigKey(c *ATNConfig, includeAlt bool) string {
	var altPart string
	if includeAlt {
		altPart = fmt.Sprintf("a%d", c.Alt)
	}
	stackParts := make([]string, len(c.Stack))
	for i, s := range c.Stack {
		stackParts[i] = fmt.Sprintf("%d", s.StateNumber)
	}
	return fmt.Sprintf("%ss%d:%s", altPart, c.State.StateNumber, strings.Join(stackParts, "_"))
}
