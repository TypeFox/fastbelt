// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"unicode/utf8"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/parallel"
)

// Lexer tokenizes a complete source string in one shot.
type Lexer interface {
	Lex(document *core.Document)
}

// Allocate a new token every ~5 characters on average
// This average is updated after lexing to adapt to the actual language
const defaultTokenRatio = 1.0 / 5.0

// DefaultLexer is the standard [Lexer] implementation. Generated `NewLexer`
// functions build one from the [core.TokenType] descriptors emitted for a
// grammar.
type DefaultLexer struct {
	tokenTypes []*core.TokenType
	tokenMap   [][]*core.TokenType
	// running exponential moving average of tokens-per-byte
	avgRatio *parallel.RunningAverage
}

// Lex scans input from left to right using longest-match disambiguation among
// token types registered at construction time.
func (l *DefaultLexer) Lex(document *core.Document) {
	input := document.TextDoc.Text(nil)
	length := len(input)
	tokens := make([]core.Token, 0, l.avgRatio.Capacity(length))
	comments := make([]core.Token, 0)
	errors := make([]*core.LexerError, 0)

	var offset int
	for offset < length {
		r, size := utf8.DecodeRuneInString(input[offset:])
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
			// No matching token, consume one rune to avoid infinite loop
			longestMatch = size
		}

		end := offset + longestMatch

		if longestType != nil {
			switch longestType.Group {
			case core.SkippedGroup:
				// do nothing
			case core.CommentGroup:
				comments = append(comments, core.NewToken(
					longestType,
					input[offset:end],
					offset, end,
				))
			case 0:
				tokens = append(tokens, core.NewToken(
					longestType,
					input[offset:end],
					offset, end,
				))
			}
		} else {
			errors = append(errors, core.NewLexerError(
				"No matching token",
				offset,
				end,
			))
		}
		offset = end
	}

	if length > 0 {
		// Update the average tokens-per-byte
		l.avgRatio.Update(float64(len(tokens)) / float64(length))
	}

	document.Tokens = tokens
	document.Comments = comments
	document.LexerErrors = errors
}

const maxChar = 256

// NewDefaultLexer returns a lexer that recognizes the given token types.
// At each position the longest match wins; among equal-length matches, the
// first argument wins.
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
		avgRatio:   parallel.NewRunningAverage(defaultTokenRatio),
	}
}
