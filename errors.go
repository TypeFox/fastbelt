// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

type ReferenceError struct {
	Msg      string
	Severity int
}

func (e *ReferenceError) Error() string {
	return e.Msg
}

func NewReferenceError(msg string) *ReferenceError {
	return &ReferenceError{Msg: msg, Severity: 1}
}

type ParserError struct {
	Msg string
	// The token this error is associated with.
	// `nil` if the error is due to EOF.
	Token *Token
}

func NewParserError(msg string, token *Token) *ParserError {
	return &ParserError{Msg: msg, Token: token}
}

type LexerError struct {
	Msg         string
	StartOffset int
	EndOffset   int
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
}

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
