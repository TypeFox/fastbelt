// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"unicode/utf8"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/parallel"
	"typefox.dev/fastbelt/util/service"
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
	sc         *service.Container
	tokenTypes [][]*core.TokenType
	tokenMaps  [][][]*core.TokenType
	// running exponential moving average of tokens-per-byte (for each language)
	avgRatio []*parallel.RunningAverage
}

// Lex scans input from left to right using longest-match disambiguation among
// the token types visible to the document's language.
func (l *DefaultLexer) Lex(document *core.Document) {
	tokenMap := l.tokenMaps[0]
	avgRatio := l.avgRatio[0]
	if len(l.tokenMaps) > 1 {
		selector := service.MustGet[core.LanguageSelector](l.sc)
		if i, _ := selector.Select(document.URI); i > 0 && i < len(l.tokenMaps) {
			tokenMap = l.tokenMaps[i]
			avgRatio = l.avgRatio[i]
		}
	}
	input := document.TextDoc.Text(nil)
	length := len(input)
	tokens := make([]core.Token, 0, avgRatio.Capacity(length))
	comments := make([]core.Token, 0)
	errors := make([]*core.LexerError, 0)

	var offset int
	for offset < length {
		r, size := utf8.DecodeRuneInString(input[offset:])
		mapIndex := int(r) % maxChar
		candidates := tokenMap[mapIndex]
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
		avgRatio.Update(float64(len(tokens)) / float64(length))
	}

	document.Tokens = tokens
	document.Comments = comments
	document.LexerErrors = errors
}

const maxChar = 256

// NewDefaultLexer returns a lexer that recognizes the given token types.
// At each position the longest match wins; among equal-length matches, the
// first argument wins.
func NewDefaultLexer(sc *service.Container, tokenTypes ...*core.TokenType) *DefaultLexer {
	return NewMultiLanguageLexer(sc, tokenTypes)
}

// NewMultiLanguageLexer returns a lexer with one token type list per language.
// The document's language is resolved via [core.LanguageSelector], mirroring
// the generated parser's entry dispatch; index 0 is the fallback for documents
// that match no language.
func NewMultiLanguageLexer(sc *service.Container, languages ...[]*core.TokenType) *DefaultLexer {
	tokenMaps := make([][][]*core.TokenType, len(languages))
	avgRatios := make([]*parallel.RunningAverage, len(languages))
	for li, tokenTypes := range languages {
		tokenMap := make([][]*core.TokenType, maxChar)
		for i := range maxChar {
			tokenMap[i] = []*core.TokenType{}
		}
		for _, tokenType := range tokenTypes {
			for _, r := range tokenType.StartChars {
				index := int(r) % maxChar
				tokenMap[index] = append(tokenMap[index], tokenType)
			}
		}
		tokenMaps[li] = tokenMap
		avgRatios[li] = parallel.NewRunningAverage(defaultTokenRatio)
	}

	return &DefaultLexer{
		tokenTypes: languages,
		tokenMaps:  tokenMaps,
		avgRatio:   avgRatios,
		sc:         sc,
	}
}
