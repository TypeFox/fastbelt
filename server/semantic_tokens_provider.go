// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// SemanticTokensProvider is a service for handling LSP semantic tokens requests.
type SemanticTokensProvider interface {
	HandleSemanticTokensFullRequest(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error)
	HandleSemanticTokensRangeRequest(ctx context.Context, params *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error)
	HandleSemanticTokensFullDeltaRequest(ctx context.Context, params *lsp.SemanticTokensDeltaParams) (any, error)
}

// DefaultSemanticTokensProvider is the default implementation of [SemanticTokensProvider].
// It delegates token classification to [SemanticTokensContributor] and handles LSP encoding.
type DefaultSemanticTokensProvider struct {
	sc *service.Container
}

func NewDefaultSemanticTokensProvider(sc *service.Container) SemanticTokensProvider {
	return &DefaultSemanticTokensProvider{sc: sc}
}

func (p *DefaultSemanticTokensProvider) HandleSemanticTokensFullRequest(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	return p.encodeSemanticTokens(ctx, string(params.TextDocument.URI), nil)
}

func (p *DefaultSemanticTokensProvider) HandleSemanticTokensRangeRequest(ctx context.Context, params *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error) {
	return p.encodeSemanticTokens(ctx, string(params.TextDocument.URI), &params.Range)
}

func (p *DefaultSemanticTokensProvider) HandleSemanticTokensFullDeltaRequest(ctx context.Context, params *lsp.SemanticTokensDeltaParams) (any, error) {
	// Delta support is not implemented. Since we recompute the entire document
	// on every change, delta computation would not provide meaningful benefits.
	// Return full tokens instead (same approach as Langium's default).
	// Future enhancement: cache token builders per document and compute diffs.
	return p.encodeSemanticTokens(ctx, string(params.TextDocument.URI), nil)
}

func (p *DefaultSemanticTokensProvider) encodeSemanticTokens(ctx context.Context, uriStr string, lspRange *lsp.Range) (*lsp.SemanticTokens, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(uriStr)
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil
	}

	contributor, err := service.Get[SemanticTokensContributor](p.sc)
	if err != nil || len(contributor.TokenTypes()) == 0 {
		return &lsp.SemanticTokens{Data: []uint32{}}, nil
	}

	data := p.encodeTokens(ctx, doc, lspRange, contributor)
	return &lsp.SemanticTokens{Data: data}, nil
}

// encodeTokens traverses document tokens and encodes them in LSP format.
func (p *DefaultSemanticTokensProvider) encodeTokens(
	ctx context.Context,
	doc *core.Document,
	lspRange *lsp.Range,
	contributor SemanticTokensContributor,
) []uint32 {
	var data []uint32
	var prevLine, prevChar uint32

	// Determine token range
	var startOffset, endOffset int
	if lspRange != nil {
		startOffset = doc.TextDoc.OffsetAt(lspRange.Start)
		endOffset = doc.TextDoc.OffsetAt(lspRange.End)
	} else {
		startOffset = 0
		endOffset = len(doc.TextDoc.Content())
	}

	// Traverse all tokens
	for i := 0; i < len(doc.Tokens); i++ {
		token := &doc.Tokens[i]
		tokenStart := token.TextSegment.Indices.Start
		tokenEnd := token.TextSegment.Indices.End

		// Skip tokens outside the range
		if int(tokenEnd) <= startOffset || int(tokenStart) >= endOffset {
			continue
		}

		// Get the AST node for this token
		node := token.Element

		// Classify the token
		typeIndex := contributor.ClassifyToken(ctx, token, node)
		if typeIndex < 0 {
			continue
		}

		// Get modifiers
		modifiers := contributor.GetModifiers(ctx, token, node)
		modifierBits := encodeModifiers(modifiers)

		// Get token position
		tokenRange := token.TextSegment.Range
		line := uint32(tokenRange.Start.Line)
		char := uint32(tokenRange.Start.Column)
		length := uint32(tokenEnd - tokenStart)

		// Compute deltas
		deltaLine := line - prevLine
		var deltaStart uint32
		if deltaLine == 0 {
			deltaStart = char - prevChar
		} else {
			deltaStart = char
		}

		// Append encoded token (5 integers)
		data = append(data,
			deltaLine,
			deltaStart,
			length,
			uint32(typeIndex),
			modifierBits,
		)

		// Update previous position
		prevLine = line
		prevChar = char
	}

	return data
}

// encodeModifiers combines modifier indices into a bit field.
func encodeModifiers(modifiers []int) uint32 {
	var bits uint32
	for _, index := range modifiers {
		if index >= 0 && index < 32 {
			bits |= (1 << index)
		}
	}
	return bits
}
