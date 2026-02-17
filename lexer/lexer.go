// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// TODO Move this stuff to the core package?
package lexer

import (
	"unicode/utf8"

	core "typefox.dev/fastbelt"
)

type LexerResult struct {
	Tokens []*core.Token
	Errors []*core.LexerError
	Groups map[int][]*core.Token
}

type Lexer interface {
	Lex(input string) *LexerResult
}

type DefaultLexer struct {
	tokenTypes []*core.TokenType
	tokenMap   [][]*core.TokenType
}

func (l *DefaultLexer) Lex(input string) *LexerResult {
	length := len(input)
	tokens := make([]*core.Token, 0, length/5) // rough estimate
	errors := make([]*core.LexerError, 0)
	groups := make(map[int][]*core.Token)

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
		var longestType *core.TokenType
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
				tokens = append(tokens, core.NewToken(
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
			errors = append(errors, core.NewLexerError(
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

func NewDefaultLexer(tokenTypes ...*core.TokenType) *DefaultLexer {
	tokenMap := make([][]*core.TokenType, maxChar)
	for i := range maxChar {
		tokenMap[i] = make([]*core.TokenType, 0)
	}
	for _, tokenType := range tokenTypes {
		for _, r := range tokenType.StartChars {
			index := int(r) % maxChar
			tokenMap[index] = append(tokenMap[index], tokenType)
		}
	}

	return &DefaultLexer{
		tokenTypes: tokenTypes,
		tokenMap:   tokenMap,
	}
}
