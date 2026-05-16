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
	UnexpectedEndOfInput() string
	// UnexpectedToken is used when the current token does not match what the
	// parser expected and recovery decided not to synthesise or skip a token.
	UnexpectedToken(found *core.Token) string
	// ExtraneousInput is used when a token is being skipped as part of
	// single-token deletion or sync-style resynchronisation.
	ExtraneousInput(found *core.Token) string
	// MissingToken is used when the parser synthesises a token of the
	// expected type because it appears to be missing at the current position.
	MissingToken(expected *core.TokenType, found *core.Token) string
}

// DefaultErrorMessageProvider produces English diagnostic messages. It is the
// default ErrorMessageProvider used when none is supplied to NewParserState.
type DefaultErrorMessageProvider struct{}

func (DefaultErrorMessageProvider) UnexpectedEndOfInput() string {
	return "Unexpected end of input."
}

func (DefaultErrorMessageProvider) UnexpectedToken(found *core.Token) string {
	return "Unexpected token '" + found.Image + "'."
}

func (DefaultErrorMessageProvider) ExtraneousInput(found *core.Token) string {
	return "Extraneous input '" + found.Image + "'."
}

func (DefaultErrorMessageProvider) MissingToken(expected *core.TokenType, found *core.Token) string {
	return "Missing '" + expected.Name + "', got '" + found.Image + "'."
}
