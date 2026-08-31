// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// SemanticTokensProvider defines the interface for handling semantic tokens requests in the LSP.
// Must be registered together with a [SemanticTokensLegendProvider] in the service container
// to enable semantic tokens support for the language server.
type SemanticTokensProvider interface {
	HandleSemanticTokensFullRequest(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error)
}

// TokenHighlightingStrategyAcceptor is a function type used in the [TokenHighlightingStrategy].
type TokenHighlightingStrategyAcceptor func(tokenType uint32, tokenModifier uint32)

// TokenHighlightingStrategy defines the interface for strategies that determine how individual tokens
// are highlighted by the [TokenBasedSemanticTokensProvider].
//
// Note that the "accept" function should only be called once per token.
// Calling it multiple times for the same token will result in an error being returned to the language client.
type TokenHighlightingStrategy interface {
	Highlight(ctx context.Context, token core.Token, accept TokenHighlightingStrategyAcceptor)
}

// CommentTokenHighlightingStrategy extends the [TokenHighlightingStrategy] interface to include a method for highlighting comment tokens.
// If a [TokenHighlightingStrategy] also implements this interface, the [TokenBasedSemanticTokensProvider] will use it to highlight
// comment tokens in addition to regular tokens.
// Otherwise, comment tokens will be highlighted using the comment token type from the legend with no modifiers.
type CommentTokenHighlightingStrategy interface {
	TokenHighlightingStrategy
	HighlightComment(ctx context.Context, commentToken core.Token, accept TokenHighlightingStrategyAcceptor)
}

// TokenBasedSemanticTokensProvider is an implementation of [SemanticTokensProvider] that generates semantic tokens
// for each individual token in the document, using a provided [TokenHighlightingStrategy] to determine the highlighting for each token.
// It also generates semantic tokens for comments in the document, if the "comment" token type is present in the legend.
type TokenBasedSemanticTokensProvider struct {
	sc                   *service.Container
	strategy             TokenHighlightingStrategy
	commentTypeIndexFunc func() int // Lazily initialized index of the comment token type in the legend
}

// NewTokenBasedSemanticTokensProvider creates a new instance of [TokenBasedSemanticTokensProvider] with the given [TokenHighlightingStrategy].
func NewTokenBasedSemanticTokensProvider(sc *service.Container, strategy TokenHighlightingStrategy) SemanticTokensProvider {
	return &TokenBasedSemanticTokensProvider{sc: sc, strategy: strategy, commentTypeIndexFunc: sync.OnceValue(func() int {
		tokenTypes := service.MustGet[SemanticTokensLegendProvider](sc).Legend().TokenTypes
		commentTypeIndex := slices.Index(tokenTypes, string(lsp.CommentType))
		return commentTypeIndex
	})}
}

func (p *TokenBasedSemanticTokensProvider) HandleSemanticTokensFullRequest(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil // Document not found
	}
	tokens := doc.Tokens
	comments := doc.Comments
	totalLen := len(tokens) + len(comments)
	if totalLen == 0 {
		return nil, nil // Document is empty, no tokens found
	}
	commentTypeIndex := p.commentTypeIndexFunc()
	tokenBuilder := NewSemanticTokensBuilder(doc.TextDoc.Text(nil), totalLen)
	highlightComment := func(commentToken core.Token) {}
	if commentStrategy, ok := p.strategy.(CommentTokenHighlightingStrategy); ok {
		// Adopter has supplied a comment highlighting strategy, use that one
		highlightComment = func(commentToken core.Token) {
			commentStrategy.HighlightComment(ctx, commentToken, func(tokenType uint32, tokenModifier uint32) {
				tokenBuilder.Push(commentToken.Range, tokenType, tokenModifier)
			})
		}
	} else if commentTypeIndex >= 0 {
		// Highlight comments using the comment token type from the legend with no modifiers
		highlightComment = func(commentToken core.Token) {
			tokenBuilder.Push(commentToken.Range, uint32(commentTypeIndex), 0)
		}
	}
	var errorRanges []core.TextRange
	commentIndex := 0
	for _, token := range tokens {
		for commentIndex < len(comments) &&
			comments[commentIndex].Range.Start < token.Range.Start {
			// Add all comments that precede the current token
			highlightComment(comments[commentIndex])
			commentIndex++
		}
		added := false
		p.strategy.Highlight(ctx, token, func(tokenType uint32, tokenModifier uint32) {
			if !added {
				tokenBuilder.Push(token.Range, tokenType, tokenModifier)
				added = true
			} else {
				errorRanges = append(errorRanges, token.Range)
			}
		})
	}
	// Report any tokens that were highlighted multiple times for the same range
	if len(errorRanges) > 0 {
		sb := strings.Builder{}
		sb.WriteString("Multiple semantic tokens returned for the same token ranges: ")
		for i, rng := range errorRanges {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(strconv.Itoa(int(rng.Start)))
			sb.WriteString("-")
			sb.WriteString(strconv.Itoa(int(rng.End)))
		}
		return nil, errors.New(sb.String())
	}
	for commentIndex < len(comments) {
		// Add remaining comments after the last token
		highlightComment(comments[commentIndex])
		commentIndex++
	}
	return &lsp.SemanticTokens{
		Data: tokenBuilder.Data(),
	}, nil
}
