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
type TokenHighlightingStrategy interface {
	Highlight(ctx context.Context, token core.Token, accept TokenHighlightingStrategyAcceptor)
}

// TokenBasedSemanticTokensProvider is an implementation of [SemanticTokensProvider] that generates semantic tokens
// for each individual token in the document, using a provided [TokenHighlightingStrategy] to determine the highlighting for each token.
// It also generates semantic tokens for comments in the document, if the "comment" token type is present in the legend.
type TokenBasedSemanticTokensProvider struct {
	sc       *service.Container
	strategy TokenHighlightingStrategy
}

// NewTokenBasedSemanticTokensProvider creates a new instance of [TokenBasedSemanticTokensProvider] with the given [TokenHighlightingStrategy].
func NewTokenBasedSemanticTokensProvider(sc *service.Container, strategy TokenHighlightingStrategy) SemanticTokensProvider {
	return &TokenBasedSemanticTokensProvider{sc: sc, strategy: strategy}
}

func (p *TokenBasedSemanticTokensProvider) HandleSemanticTokensFullRequest(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	tokenTypes := service.MustGet[SemanticTokensLegendProvider](p.sc).Legend().TokenTypes
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
	commentTypeIndex := slices.Index(tokenTypes, string(lsp.CommentType))
	tokenBuilder := NewSemanticTokensBuilder(doc.TextDoc.Text(nil), totalLen)
	var errorRanges []core.TextRange
	commentIndex := 0
	for _, token := range tokens {
		for commentTypeIndex != -1 &&
			commentIndex < len(comments) &&
			comments[commentIndex].Range.Start < token.Range.Start {
			// Add all comments that precede the current token
			tokenBuilder.Push(comments[commentIndex].Range, uint32(commentTypeIndex), 0)
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
	if commentTypeIndex != -1 {
		for commentIndex < len(comments) {
			// Add remaining comments after the last token
			tokenBuilder.Push(comments[commentIndex].Range, uint32(commentTypeIndex), 0)
			commentIndex++
		}
	}
	return &lsp.SemanticTokens{
		Data: tokenBuilder.Data(),
	}, nil
}
