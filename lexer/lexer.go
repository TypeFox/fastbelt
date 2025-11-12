// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"unicode/utf8"
)

type LexerResult struct {
	Tokens []*Token
	Errors []*LexerError
	Groups map[int][]*Token
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

type Lexer interface {
	Lex(input string) *LexerResult
}

type lexer struct {
	tokenTypes []*TokenType
	tokenMap   [][]*TokenType
}

func (l *lexer) Lex(input string) *LexerResult {
	length := len(input)
	tokens := make([]*Token, 0, length/5) // rough estimate
	errors := make([]*LexerError, 0)
	groups := make(map[int][]*Token)

	var offset, line, column int
	line = 0
	column = 0
	for offset < length {
		r, size := rune(input[offset]), 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRuneInString(input[offset:])
		}
		mapIndex := int(r) % maxChar
		candidates := l.tokenMap[mapIndex]
		longestMatch := 0
		var longestType *TokenType
		for _, tokenType := range candidates {
			tokenMatch := tokenType.Match(input, offset)
			if tokenMatch > longestMatch {
				longestMatch = tokenMatch
				longestType = tokenType
			}
		}

		if longestMatch == 0 {
			longestMatch = size
		}

		end := offset + longestMatch
		startLine := line
		startColumn := column

		for i := offset; i < end; i++ {
			if input[i] == '\n' {
				line++
				column = 0
			} else {
				column++
			}
		}

		if longestType != nil {
			if !longestType.IsSkipped() {
				tokens = append(tokens, NewToken(
					longestType,
					input[offset:end],
					offset, end,
					startLine,
					line,
					startColumn,
					column,
				))
			}
		} else {
			errors = append(errors, NewLexerError(
				"No matching token",
				offset,
				end,
				startLine,
				line,
				startColumn,
				column,
			))
		}
		offset = end
	}

	return &LexerResult{
		Tokens: tokens,
		Errors: errors,
		Groups: groups,
	}
}

const maxChar = 256

func NewLexer(tokenTypes ...*TokenType) Lexer {
	tokenMap := make([][]*TokenType, maxChar)
	for i := range maxChar {
		tokenMap[i] = make([]*TokenType, 0)
	}
	for _, tokenType := range tokenTypes {
		for _, r := range tokenType.StartChars {
			if r < maxChar {
				tokenMap[int(r)] = append(tokenMap[int(r)], tokenType)
			}
		}
	}

	return &lexer{
		tokenTypes: tokenTypes,
		tokenMap:   tokenMap,
	}
}
