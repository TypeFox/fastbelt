// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

// atnConfig is one configuration in an ALL(*) simulation: an ATN state reached
// while predicting a particular alternative, together with the call-stack
// context that led there.
//
// alt is 0-based to match Fastbelt's generated switch-case alternative
// numbering (the i-th epsilon transition out of a decision state is alt i).
type atnConfig struct {
	state   *RuntimeATNState
	alt     int
	context *predictionContext
	// reachesIntoOuterContext counts how far closure escaped the start rule via
	// rule-stop pops with an empty/exhausted context. A non-zero value means the
	// config "dips into the outer context" and SLL conflict resolution must be
	// treated conservatively.
	reachesIntoOuterContext int
}

func newATNConfig(state *RuntimeATNState, alt int, context *predictionContext) *atnConfig {
	return &atnConfig{state: state, alt: alt, context: context}
}

// newATNConfigWithContext copies c but with a different context, preserving the
// alt and reachesIntoOuterContext (mirrors ANTLR's NewATNConfig4).
func newATNConfigWithContext(c *atnConfig, context *predictionContext) *atnConfig {
	return &atnConfig{
		state:                   c.state,
		alt:                     c.alt,
		context:                 context,
		reachesIntoOuterContext: c.reachesIntoOuterContext,
	}
}

// newATNConfigWithState copies c but at a different state.
func newATNConfigWithState(c *atnConfig, state *RuntimeATNState, context *predictionContext) *atnConfig {
	return &atnConfig{
		state:                   state,
		alt:                     c.alt,
		context:                 context,
		reachesIntoOuterContext: c.reachesIntoOuterContext,
	}
}

// Equals reports configuration identity: same state, same alternative, and
// equal context. (Semantic context is intentionally absent.)
func (c *atnConfig) Equals(o *atnConfig) bool {
	if c == o {
		return true
	}
	if o == nil {
		return false
	}
	if c.state != o.state || c.alt != o.alt {
		return false
	}
	if c.context == nil || o.context == nil {
		return c.context == o.context
	}
	return c.context.Equals(o.context)
}

func (c *atnConfig) Hash() uint64 {
	h := uint64(7)
	h = hashCombine(h, uint64(c.state.StateNumber))
	h = hashCombine(h, uint64(c.alt))
	if c.context != nil {
		h = hashCombine(h, c.context.Hash())
	}
	return h
}
