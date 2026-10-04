// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"testing"

	core "typefox.dev/fastbelt"
)

// TestCompletionParserState_SnapshotsAndStack verifies the bookkeeping
// methods: EnterRule/ExitRule push and pop, RecordSnapshot/EnterRule both
// append to Snapshots.
func TestCompletionParserState_SnapshotsAndStack(t *testing.T) {
	tA := core.NewTokenType(1, "a", "a", 0, nil, nil)

	// Build a small ATN so NewParserState has something to hold.
	s0 := &RuntimeATNState{StateNumber: 0, Type: ATNBasic, Decision: -1}
	atn := NewRuntimeATN([]*RuntimeATNState{s0}, nil, nil)

	tokens := []core.Token{tok(tA), tok(tA), tok(tA)}
	state := NewParserState(tokens, atn, &BailErrorRecovery{}, &DefaultErrorMessageProvider{})
	cp := NewCompletionParserState(state)

	// Outer rule enters at token 0, ATN start state 7.
	cp.EnterRule("Outer", 7)
	if got := len(cp.snapshots); got != 1 {
		t.Fatalf("expected 1 snapshot after EnterRule; got %d", got)
	}
	if got := cp.snapshots[0]; got.TokenIdx != 0 || got.ATNStateIdx != 7 {
		t.Fatalf("unexpected first snapshot: %+v", got)
	}
	if got := cp.snapshots[0].RuleStack; len(got) != 1 || got[0].RuleKey != "Outer" {
		t.Fatalf("expected RuleStack=[Outer] on first snapshot; got %+v", got)
	}

	// Simulate consuming one token then a Sync to ATN state 12.
	state.Index = 1
	cp.RecordSnapshot(12)
	if got := len(cp.snapshots); got != 2 {
		t.Fatalf("expected 2 snapshots after RecordSnapshot; got %d", got)
	}
	if got := cp.snapshots[1]; got.TokenIdx != 1 || got.ATNStateIdx != 12 {
		t.Fatalf("unexpected second snapshot: %+v", got)
	}

	// Enter nested rule.
	state.Index = 2
	cp.EnterRule("Inner", 33)
	if got := len(cp.ruleStack); got != 2 {
		t.Fatalf("expected stack depth 2; got %d", got)
	}

	res := cp.Result(tokens)
	// Result.RuleStack is the deepest stack at the latest snapshot's
	// TokenIdx - here, [Outer, Inner] from the Inner EnterRule snapshot at
	// TokenIdx=2.
	if len(res.RuleStack) != 2 || res.RuleStack[0].RuleKey != "Outer" || res.RuleStack[1].RuleKey != "Inner" {
		t.Errorf("unexpected RuleStack snapshot: %+v", res.RuleStack)
	}

	// ExitRule pops the inner frame.
	cp.ExitRule()
	if got := len(cp.ruleStack); got != 1 {
		t.Fatalf("expected stack depth 1 after ExitRule; got %d", got)
	}
}
