// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import core "typefox.dev/fastbelt"

// ErrorRecoveryStrategy is the pluggable interface for parser error recovery.
type ErrorRecoveryStrategy interface {
	// RecoverInline is called by Consume when the next token does not match
	// the expected type. Implementations can attempt to find a token in the token stream
	// that would match the expectation and return it as if it were a successful match.
	// The boolean return value indicates whether the parser should treat the error as recovered (true) or not (false).
	// If false, the parser will set InError and unwind.
	RecoverInline(parserState *ParserState, expectedTokenType *core.TokenType) (*core.Token, bool)
	// Recover resynchronises the input stream after a hard-halt error.
	Recover(parserState *ParserState)
	// Sync is called before optional/loop guards to discard unexpected tokens.
	Sync(parserState *ParserState, decisionStateIdx int)
}

// DefaultErrorRecovery implements single-token deletion in
// RecoverInline plus a consume-until-FOLLOW strategy in Sync.
type DefaultErrorRecovery struct{}

// See [DefaultErrorRecovery].
func NewDefaultErrorRecovery() DefaultErrorRecovery {
	return DefaultErrorRecovery{}
}

// BailErrorRecovery stops on the first mismatch without attempting any
// recovery: its RecoverInline always returns nil, which sets inError and
// unwinds the parser via the LA-nil mechanism. Useful for two-stage parsing
// or scenarios where the parse result is discarded on any error.
type BailErrorRecovery struct{}

func NewBailErrorRecovery() BailErrorRecovery {
	return BailErrorRecovery{}
}

func (DefaultErrorRecovery) RecoverInline(parserState *ParserState, expectedTokenType *core.TokenType) (*core.Token, bool) {
	la1 := parserState.LARaw(1)
	if la1 == nil {
		parserState.AppendError(parserState.messages.UnexpectedEndOfInput(expectedTokenType), nil)
		return nil, false
	}
	la2 := parserState.LARaw(2)
	// Single-token deletion: if the next-next token matches, skip the current one.
	if la2 != nil && la2.TypeId == expectedTokenType.Id {
		parserState.ReportError(parserState.messages.ExtraneousInput(la1), la1)
		// Skip the bad token and return the next one as if it were a match.
		parserState.Index += 2
		parserState.ReportMatch()
		return la2, true
	}
	// Most other parser generators would attempt single-token insertion here
	// (i.e. return a fabricated token of the expected type and leave the input stream alone),
	// but our generated AST can deal with nil tokens, so we don't need to generate anything.
	// Note: Index is NOT advanced and errorRecoveryMode stays on, so subsequent
	// reportError calls at the same position remain suppressed.
	parserState.ReportError(parserState.messages.MissingToken(expectedTokenType, la1), la1)
	return nil, true
}

func (DefaultErrorRecovery) Recover(parserState *ParserState) {
	if parserState.ErrorMode != ErrorModeRecover {
		return
	}
	// Compute the set of tokens that could legally come next at this point in the parse.
	// Discards tokens until we find one that matches, or until we hit EOF.
	followSet := parserState.FollowSet()
	if len(followSet) > 0 {
		for {
			la := parserState.LARaw(1)
			if la == nil {
				// EOF reached, give up and let the parser unwind.
				return
			}
			if tokenInSet(followSet, la.TypeId) {
				break
			}
			parserState.Index++
		}
	}
	parserState.ErrorMode = ErrorModeNone
	parserState.ErrorRecoveryMode = false
}

// Attempts to recover from a parser error by calling the current stategy's Recover method.
// Immediately returns if the recovery fails.
// If not within in error, it tries to ensure that the upcoming token is valid for the current decision state.
// This ensures that the parser can continue to make progress and doesn't get stuck on a bad token.
func (DefaultErrorRecovery) Sync(parserState *ParserState, decisionStateIdx int) {
	if parserState.ErrorMode == ErrorModeRecover {
		parserState.RecoveryStrategy().Recover(parserState)
		if parserState.ErrorMode != ErrorModeNone {
			return
		}
	}
	tok := parserState.LARaw(1)
	validTokens := parserState.atn.NextTokensAt(decisionStateIdx)
	// This is the expected case: LA(1) is in the set of valid tokens.
	// We can immediately return, the parser is valid and ready to continue.
	// This is the hot path for error-free parsing, so it's important that this check is fast.
	if tok == nil || tokenInSet(validTokens, tok.TypeId) {
		return
	}
	followTokens := parserState.FollowSet()
	if tokenInSet(followTokens, tok.TypeId) {
		return
	}
	// Single-token deletion: if the *next* token is in the valid set, treat
	// the current token as extraneous and skip it.
	la2 := parserState.LARaw(2)
	if la2 != nil && tokenInSet(validTokens, la2.TypeId) {
		parserState.ReportError(parserState.messages.ExtraneousInput(tok), tok)
		parserState.Index++
		return
	}
	// Otherwise consume unexpected tokens until we land on something that
	// belongs to either the decision's valid set or any enclosing follow set.
	// Skipping only one token can leave the parser stuck when several
	// unexpected tokens cluster together (e.g. after a broken rule body).
	parserState.ReportError(parserState.messages.ExtraneousInput(tok), tok)
	for {
		parserState.Index++
		la := parserState.LARaw(1)
		if la == nil || tokenInSet(validTokens, la.TypeId) || tokenInSet(followTokens, la.TypeId) {
			return
		}
	}
}

// tokenInSet reports whether id is set in a token bitset returned from
// NextTokensAt or computeFollowSet. Out-of-range ids are treated as absent.
func tokenInSet(set []bool, id int) bool {
	return id >= 0 && id < len(set) && set[id]
}

func (BailErrorRecovery) RecoverInline(_ *ParserState, _ *core.TokenType) (*core.Token, bool) {
	return nil, false
}

func (BailErrorRecovery) Recover(_ *ParserState) {}

func (BailErrorRecovery) Sync(_ *ParserState, _ int) {}
