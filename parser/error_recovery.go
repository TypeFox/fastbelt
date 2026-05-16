// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import core "typefox.dev/fastbelt"

// ErrorRecoveryStrategy is the pluggable interface for parser error recovery.
type ErrorRecoveryStrategy interface {
	// RecoverInline is called by Consume when the next token does not match
	// the expected type. Implementations either return a recovered token
	// (real, via single-token deletion; or synthetic, via single-token
	// insertion) so the caller can continue, or return nil to bail.
	// Returning nil causes Consume to record an error and set inError,
	// which hard-halts the parser — only the bail strategy should do this.
	// DefaultErrorRecovery only returns nil when the input is already at EOF.
	RecoverInline(state *ParserState, expectedTokenType *core.TokenType) *core.Token
	// Recover resynchronises the input stream after a hard-halt error by
	// consuming tokens until one in the FOLLOW set of the enclosing rules
	// appears, then clears both inError and errorRecoveryMode so parsing
	// can resume. Currently invoked from Sync when it observes inError.
	Recover(state *ParserState)
	// Sync is called before optional/loop guards to discard unexpected tokens
	// that are neither in FIRST(decision state) nor in FOLLOW(enclosing rules).
	// It reports errors via reportError (subject to errorRecoveryMode dedup)
	// and advances Index itself, so the caller need not check anything.
	Sync(state *ParserState, decisionStateIdx int)
}

// DefaultErrorRecovery implements single-token deletion/insertion in
// RecoverInline plus a consume-until-FOLLOW strategy in Sync. RecoverInline
// never returns nil except at true EOF, so this strategy never hard-halts
// the parser on its own.
type DefaultErrorRecovery struct{}

// BailErrorRecovery stops on the first mismatch without attempting any
// recovery: its RecoverInline always returns nil, which sets inError and
// unwinds the parser via the LA-nil mechanism. Useful for two-stage parsing
// or scenarios where the parse result is discarded on any error.
type BailErrorRecovery struct{}

func (DefaultErrorRecovery) RecoverInline(state *ParserState, expectedTokenType *core.TokenType) *core.Token {
	la1 := state.laRaw(1)
	if la1 == nil {
		state.appendError("Unexpected end of input.", nil)
		return nil
	}
	// Single-token deletion: if the next-next token matches, skip the current one.
	if state.LAId(2) == expectedTokenType.Id {
		state.reportError("Extraneous input '"+la1.Image+"'.", la1)
		state.Index++
		token := state.LA(1)
		state.Index++
		// A real token was matched after deletion — exit error-recovery mode.
		state.reportMatch()
		return token
	}
	// Single-token insertion: synthesise a zero-width token of the expected type.
	// Note: Index is NOT advanced and errorRecoveryMode stays on, so subsequent
	// reportError calls at the same position remain suppressed.
	state.reportError("Missing '"+expectedTokenType.Name+"', got '"+la1.Image+"'.", la1)
	synthetic := makeSyntheticToken(expectedTokenType, la1)
	return &synthetic
}

func (DefaultErrorRecovery) Recover(state *ParserState) {
	if !state.inError {
		return
	}
	followSet := state.computeFollowSet()
	if len(followSet) > 0 {
		for {
			la := state.laRaw(1)
			if la == nil || tokenInSet(followSet, la.TypeId) {
				break
			}
			state.Index++
		}
	}
	state.inError = false
	state.errorRecoveryMode = false
}

func (DefaultErrorRecovery) Sync(state *ParserState, decisionStateIdx int) {
	if state.inError {
		state.recovery.Recover(state)
		if state.inError {
			return
		}
	}
	tok := state.laRaw(1)
	if tok == nil {
		return
	}
	valid := state.atn.NextTokensAt(decisionStateIdx)
	if tokenInSet(valid, tok.TypeId) {
		return
	}
	follow := state.computeFollowSet()
	if tokenInSet(follow, tok.TypeId) {
		return
	}
	// Single-token deletion: if the *next* token is in the valid set, treat
	// the current token as extraneous and skip it.
	la2 := state.laRaw(2)
	if la2 != nil && tokenInSet(valid, la2.TypeId) {
		state.reportError("Extraneous input '"+tok.Image+"'.", tok)
		state.Index++
		return
	}
	// Otherwise consume garbage tokens until we land on something that
	// belongs to either the decision's valid set or any enclosing follow set.
	// Skipping only one token can leave the parser stuck when several
	// unexpected tokens cluster together (e.g. after a broken rule body).
	state.reportError("Extraneous input '"+tok.Image+"'.", tok)
	for {
		state.Index++
		la := state.laRaw(1)
		if la == nil || tokenInSet(valid, la.TypeId) || tokenInSet(follow, la.TypeId) {
			return
		}
	}
}

// tokenInSet reports whether id is set in a token bitset returned from
// NextTokensAt or computeFollowSet. Out-of-range ids are treated as absent.
func tokenInSet(set []bool, id int) bool {
	return id >= 0 && id < len(set) && set[id]
}

func (BailErrorRecovery) RecoverInline(_ *ParserState, _ *core.TokenType) *core.Token {
	return nil
}

func (BailErrorRecovery) Recover(_ *ParserState) {}

func (BailErrorRecovery) Sync(_ *ParserState, _ int) {}

// makeSyntheticToken creates a zero-width token at the same position as near
func makeSyntheticToken(tokenType *core.TokenType, near *core.Token) core.Token {
	start := int(near.TextSegment.Indices.Start)
	line := int(near.TextSegment.Range.Start.Line)
	col := int(near.TextSegment.Range.Start.Column)
	return core.NewToken(tokenType, "", start, start, line, line, col, col)
}
