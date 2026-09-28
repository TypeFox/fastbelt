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

// Lexer tokenizes a document in one shot and stores the result on it.
type Lexer interface {
	// Exec tokenizes document.TextDoc and sets [core.Document.Tokens],
	// [core.Document.Comments], and [core.Document.LexerErrors].
	Exec(document *core.Document)
}

// lexerResult holds everything produced by a single [DefaultLexer.Lex] pass
// over source text.
type lexerResult struct {
	// Tokens is the main token stream passed to the parser.
	Tokens []core.Token
	// Comments holds tokens whose Type ([core.TokenType]) is marked with the modifier [core.CommentModifier] via the type [TokenTypeUsage].
	// A [TokenMode] holds all relevant token types. Each [core.TokenType] can be put into a [TokenTypeUsage] containing a modifier and a Token command.
	Comments []core.Token
	// Errors lists recoverable lexing problems (unrecognized input).
	Errors []*core.LexerError
}

// Allocate a new token every ~5 characters on average
// This average is updated after lexing to adapt to the actual language
const defaultTokenRatio = 1.0 / 5.0

// DefaultLexer is the standard [Lexer] implementation. Generated `NewLexer`
// functions build one from the [core.TokenType] descriptors emitted for a
// grammar.
type DefaultLexer struct {
	sc *service.Container
	// one token mode list per language; index 0 is the fallback
	languages [][]*TokenMode
	// per language: index into its token modes of the mode every run starts in
	startModes []int
	// running exponential moving average of tokens-per-byte (per language)
	avgRatio []*parallel.RunningAverage
}

// Exec tokenizes the document with the token modes of its language. The
// language is resolved via [core.LanguageSelector] when more than one was
// registered, mirroring the generated parser's entry dispatch.
func (l *DefaultLexer) Exec(document *core.Document) {
	language := 0
	if len(l.languages) > 1 {
		selector := service.MustGet[core.LanguageSelector](l.sc)
		if i, _ := selector.Select(document.URI); i > 0 && i < len(l.languages) {
			language = i
		}
	}
	result := l.exec(document.TextDoc.Text(nil), language)
	document.Tokens = result.Tokens
	document.Comments = result.Comments
	document.LexerErrors = result.Errors
}

// exec scans input from left to right using longest-match disambiguation among
// the token types of the active mode of the given language.
func (l *DefaultLexer) exec(input string, language int) *lexerResult {
	tokenModes := l.languages[language]
	avgRatio := l.avgRatio[language]
	length := len(input)
	tokens := make([]core.Token, 0, avgRatio.Capacity(length))
	comments := make([]core.Token, 0)
	errors := make([]*core.LexerError, 0)

	// The mode stack is local to this call: a DefaultLexer is shared between
	// documents and Exec may run concurrently, so input that ends inside a
	// pushed mode must not leak into the next run.
	stack := NewTokenModeStack(tokenModes[l.startModes[language]])
	currentTokenMode := stack.Peek()

	var offset int
	for offset < length && currentTokenMode != nil {
		r, size := utf8.DecodeRuneInString(input[offset:])
		mapIndex := int(r) % maxChar
		candidates := currentTokenMode.TokenMap[mapIndex]
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
			switch longestType.Modifier {
			case core.SkippedModifier:
				// do nothing
			case core.CommentModifier:
				comments = append(comments, core.NewToken(
					longestType.TokenType,
					input[offset:end],
					offset, end,
				))
			case core.DefaultTokenModifier:
				tokens = append(tokens, core.NewToken(
					longestType.TokenType,
					input[offset:end],
					offset, end,
				))
			}

			switch {
			case longestType.PopMode && longestType.PushMode > -1:
				// `mode(X)` replaces the active mode without deepening the
				// stack, so a later pop returns to whatever was below it rather
				// than to the mode that was replaced. This mirrors ANTLR's
				// `mode` and is intentional: only `push` can be undone by `pop`.
				stack.SetMode(tokenModes[longestType.PushMode])
				currentTokenMode = stack.Peek()
			case longestType.PopMode:
				stack.Pop()
				currentTokenMode = stack.Peek()
			case longestType.PushMode > -1:
				stack.Push(tokenModes[longestType.PushMode])
				currentTokenMode = stack.Peek()
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

	return &lexerResult{
		Tokens:   tokens,
		Comments: comments,
		Errors:   errors,
	}
}

const maxChar = 256

// NewDefaultLexer returns a [DefaultLexer] that starts every run in
// tokenModes[defaultMode]. The returned lexer is safe for concurrent use.
func NewDefaultLexer(sc *service.Container, defaultMode int, tokenModes ...*TokenMode) *DefaultLexer {
	return NewMultiLanguageLexer(sc, []int{defaultMode}, tokenModes)
}

// NewMultiLanguageLexer returns a lexer with one token mode list per language.
// startModes holds, per language, the index of the token mode a run starts in.
// Mode indices are shared across languages, so each list must have the same
// length. The document's language is resolved via [core.LanguageSelector];
// index 0 is the fallback for documents that match no language.
func NewMultiLanguageLexer(sc *service.Container, startModes []int, languages ...[]*TokenMode) *DefaultLexer {
	if len(languages) == 0 {
		panic("lexer: at least one language is required")
	}
	avgRatios := make([]*parallel.RunningAverage, len(languages))
	for i, tokenModes := range languages {
		if startModes[i] < 0 || startModes[i] >= len(tokenModes) {
			panic("lexer: start token mode index out of range")
		}
		avgRatios[i] = parallel.NewRunningAverage(defaultTokenRatio)
	}
	return &DefaultLexer{
		sc:         sc,
		languages:  languages,
		startModes: startModes,
		avgRatio:   avgRatios,
	}
}
