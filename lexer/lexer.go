// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"math"
	"sync/atomic"
	"unicode/utf8"

	core "typefox.dev/fastbelt"
)

type LexerResult struct {
	Tokens   []core.Token
	Comments []core.Token
	Errors   []*core.LexerError
	Groups   map[int][]core.Token
}

type Lexer interface {
	Lex(input string) *LexerResult
}

// Allocate a new token every ~5 characters on average
// This average is updated after lexing to adapt to the actual language
const defaultTokenRatio = 1.0 / 5.0

type DefaultLexer struct {
	tokenTypes []*core.TokenType
	tokenMap   [][]*core.TokenType
	// running exponential moving average of tokens-per-byte
	// stores a float64 as uint64 bits for atomic access
	avgRatio atomic.Uint64
}

func (l *DefaultLexer) Lex(input string) *LexerResult {
	length := len(input)
	ratio := math.Float64frombits(l.avgRatio.Load())
	if ratio == 0 {
		ratio = defaultTokenRatio
	}
	tokens := make([]core.Token, 0, int(float64(length)*ratio*1.1))
	comments := make([]core.Token, 0)
	errors := make([]*core.LexerError, 0)
	var groups map[int][]core.Token

	var offset, line, column int
	line = 0
	column = 0
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
			switch longestType.Group {
			case core.SkippedGroup:
				// do nothing
			case core.CommentGroup:
				comments = append(comments, core.NewToken(
					longestType,
					input[offset:end],
					offset, end,
					startLine,
					line,
					startColumn,
					column,
				))
			case 0:
				tokens = append(tokens, core.NewToken(
					longestType,
					input[offset:end],
					offset, end,
					startLine,
					line,
					startColumn,
					column,
				))
			default:
				if groups == nil {
					groups = make(map[int][]core.Token)
				}
				groups[longestType.Group] = append(groups[longestType.Group], core.NewToken(
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

	if length > 0 {
		// Update the average tokens-per-byte
		actual := float64(len(tokens)) / float64(length)
		prev := math.Float64frombits(l.avgRatio.Load())
		if prev == 0 {
			prev = defaultTokenRatio
		}
		l.avgRatio.Store(math.Float64bits(prev*0.9 + actual*0.1))
	}

	return &LexerResult{
		Tokens:   tokens,
		Comments: comments,
		Errors:   errors,
		Groups:   groups,
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
