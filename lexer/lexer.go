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
	Exec(input string) *LexerResult
}

// LexerResult holds everything produced by a single [Lexer.Lex] pass over
// source text.
type LexerResult struct {
	// Tokens is the main token stream passed to the parser.
	Tokens []core.Token
	// Comments holds tokens whose [core.TokenType] uses [core.CommentGroup].
	Comments []core.Token
	// Errors lists recoverable lexing problems (unrecognized input).
	Errors []*core.LexerError
	// Groups collects tokens routed to custom [core.TokenType.Group] values
	// other than the default, skipped, or comment groups. Nil when empty.
	Groups map[int][]core.Token
}

const DefaultMode = "DefaultMode"

type Mode struct {
	Name       string
	TokenTypes []*core.TokenType
}

func NewMode(name string, tokenTypes ...*core.TokenType) *Mode {
	return &Mode{
		Name:       name,
		TokenTypes: tokenTypes,
	}
}

func NewDefaultMode(tokenTypes ...*core.TokenType) *Mode {
	return NewMode(DefaultMode, tokenTypes...)
}

// Allocate a new token every ~5 characters on average
// This average is updated after lexing to adapt to the actual language
const defaultTokenRatio = 1.0 / 5.0

// DefaultLexer is the standard [Lexer] implementation. Generated `NewLexer`
// functions build one from the [core.TokenType] descriptors emitted for a
// grammar.
type DefaultLexer struct {
	tokenModes []*TokenMode
	stack      TokenModeStack
	// running exponential moving average of tokens-per-byte
	avgRatio *parallel.RunningAverage
}

// Exec scans input from left to right using longest-match disambiguation among
// token types registered at construction time.
func (l *DefaultLexer) Exec(input string) *LexerResult {
	length := len(input)
	tokens := make([]core.Token, 0, l.avgRatio.Capacity(length))
	comments := make([]core.Token, 0)
	errors := make([]*core.LexerError, 0)
	var groups map[int][]core.Token

	var offset int
	for offset < length {
		r, size := utf8.DecodeRuneInString(input[offset:])
		mapIndex := int(r) % maxChar
		candidates := l.stack.Peek().TokenMap[mapIndex]
		longestMatch := 0
		var longestType *TokenTypeUsage
		for _, tokenTypeUsage := range candidates {
			tokenType := tokenTypeUsage.TokenType
			tokenMatch := tokenType.Match(input, offset)
			if tokenMatch > longestMatch {
				longestMatch = tokenMatch
				longestType = tokenTypeUsage
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
					longestType.TokenType,
					input[offset:end],
					offset, end,
				))
			case 0:
				tokens = append(tokens, core.NewToken(
					longestType.TokenType,
					input[offset:end],
					offset, end,
				))
			default:
				if groups == nil {
					groups = make(map[int][]core.Token)
				}
				groups[longestType.Group] = append(groups[longestType.Group], core.NewToken(
					longestType.TokenType,
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

	return &LexerResult{
		Tokens:   tokens,
		Comments: comments,
		Errors:   errors,
		Groups:   groups,
	}
}

const maxChar = 256

func NewDefaultLexer(tokenTypes ...*TokenMode) *DefaultLexer {
	return &DefaultLexer{
		tokenModes: tokenTypes,
		stack:      *NewTokenModeStack(tokenTypes[0]),
		avgRatio:   parallel.NewRunningAverage(defaultTokenRatio),
	}
}
