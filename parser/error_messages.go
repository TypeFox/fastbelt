// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"strings"

	core "typefox.dev/fastbelt"
)

// ErrorMessageProvider produces the user-facing strings attached to parser
// errors. Implementations can override individual messages to change wording
// or provide localised (i18n) variants.
type ErrorMessageProvider interface {
	// UnexpectedEndOfInput is used when the parser needs another token but
	// the token stream is already exhausted.
	UnexpectedEndOfInput(expected *core.TokenType) string
	// UnexpectedToken is used when the current token does not match what the
	// parser expected and recovery decided not to synthesise or skip a token.
	UnexpectedToken(found *core.Token) string
	// ExtraneousInput is used when a token is being skipped as part of
	// single-token deletion or sync-style resynchronisation.
	ExtraneousInput(found *core.Token) string
	// MissingToken is used when the parser synthesises a token of the
	// expected type because it appears to be missing at the current position.
	MissingToken(expected *core.TokenType, found *core.Token) string
	// NoViableAlternative is used when the parser cannot decide between multiple
	// alternatives. The failure carries the divergence token (where input
	// dead-ended) and, when adaptive prediction produced it, the set of token
	// types that would have allowed prediction to continue.
	NoViableAlternative(failure *PredictionFailure) string
}

// DefaultErrorMessageProvider produces English diagnostic messages. It is the
// default ErrorMessageProvider used when none is supplied.
type DefaultErrorMessageProvider struct{}

func NewDefaultErrorMessageProvider() DefaultErrorMessageProvider {
	return DefaultErrorMessageProvider{}
}

func (DefaultErrorMessageProvider) UnexpectedEndOfInput(expected *core.TokenType) string {
	return "Unexpected end of input, expected '" + expected.Name + "'."
}

func (DefaultErrorMessageProvider) UnexpectedToken(found *core.Token) string {
	return "Unexpected token '" + tokenImage(found) + "'."
}

func (DefaultErrorMessageProvider) ExtraneousInput(found *core.Token) string {
	return "Extraneous input '" + tokenImage(found) + "'."
}

func (DefaultErrorMessageProvider) MissingToken(expected *core.TokenType, found *core.Token) string {
	return "Missing '" + expected.Name + "', got '" + tokenImage(found) + "'."
}

func (DefaultErrorMessageProvider) NoViableAlternative(failure *PredictionFailure) string {
	var found *core.Token
	if failure != nil {
		found = failure.Token
	}
	if failure == nil || len(failure.Expected) == 0 {
		return "No viable alternative at input '" + tokenImage(found) + "'."
	}
	return "No viable alternative at input '" + tokenImage(found) + "', expected one of: " + joinExpected(failure.Expected) + "."
}

// joinExpected renders expected token types as a comma-separated list of quoted
// names, e.g. "'A', 'B'".
func joinExpected(expected []*core.TokenType) string {
	var b strings.Builder
	for i, tt := range expected {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteByte('\'')
		b.WriteString(tt.Name)
		b.WriteByte('\'')
	}
	return b.String()
}

func tokenImage(t *core.Token) string {
	if t == nil {
		return "<EOF>"
	}
	return t.Image
}
