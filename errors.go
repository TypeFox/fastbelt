// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

// A ReferenceError reports a failure while resolving a cross-reference.
// It is converted to diagnostics during linking.
type ReferenceError struct {
	// Msg is the human-readable error message shown in diagnostics.
	Msg string
	// Severity is an LSP-compatible [DiagnosticSeverity] value.
	// Use [SeverityError], [SeverityWarning], [SeverityInfo], or [SeverityHint].
	Severity int
}

// NewReferenceError returns a [ReferenceError] with [SeverityError].
func NewReferenceError(msg string) *ReferenceError {
	return &ReferenceError{Msg: msg, Severity: 1}
}

// A ParserError describes a syntax error detected by the parser.
type ParserError struct {
	// Msg is the parser diagnostic message.
	Msg string
	// Token is the token this error is associated with.
	// It is nil when the parser error refers to unexpected end of input.
	Token *Token
}

// NewParserError returns a [ParserError] for msg at token.
func NewParserError(msg string, token *Token) *ParserError {
	return &ParserError{Msg: msg, Token: token}
}

// A LexerError describes invalid input encountered during tokenization.
type LexerError struct {
	// Msg is the lexer diagnostic message.
	Msg string
	// StartOffset is the byte offset where the invalid segment starts.
	StartOffset int
	// EndOffset is the exclusive byte offset where the invalid segment ends.
	EndOffset int
	// StartLine is the zero-based line where the invalid segment starts.
	StartLine int
	// EndLine is the zero-based line where the invalid segment ends.
	EndLine int
	// StartColumn is the zero-based column where the invalid segment starts.
	StartColumn int
	// EndColumn is the zero-based column where the invalid segment ends.
	EndColumn int
}

// NewLexerError returns a [LexerError] for msg and the offending text range.
func NewLexerError(msg string, startOffset, endOffset, startLine, endLine, startColumn, endColumn int) *LexerError {
	return &LexerError{
		Msg:         msg,
		StartOffset: startOffset,
		EndOffset:   endOffset,
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
	}
}
