// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import core "typefox.dev/fastbelt"

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
	// NoViableAlternative is used when the parser cannot decide between multiple alternatives.
	// TODO: This needs a more complex signature, but this depends on the LL(*) implementation.
	NoViableAlternative(found *core.Token) string
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

func (DefaultErrorMessageProvider) NoViableAlternative(found *core.Token) string {
	return "No viable alternative at input '" + tokenImage(found) + "'."
}

func tokenImage(t *core.Token) string {
	if t == nil {
		return "<EOF>"
	}
	return t.Image
}
