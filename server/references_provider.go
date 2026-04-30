// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// ReferencesProvider is a service for handling LSP reference requests.
type ReferencesProvider interface {
	HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error)
}

// DefaultReferencesProvider is the default implementation of [ReferencesProvider].
type DefaultReferencesProvider struct {
	sc *service.Container
}

func NewDefaultReferencesProvider(sc *service.Container) ReferencesProvider {
	return &DefaultReferencesProvider{sc: sc}
}

func (s *DefaultReferencesProvider) HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error) {
	documentManager := service.MustGet[workspace.DocumentManager](s.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	targetDoc := documentManager.Get(uri)
	if targetDoc == nil {
		return nil, nil // Document not found
	}
	offset := targetDoc.TextDoc.OffsetAt(params.Position)
	tokens := targetDoc.Tokens
	sourceToken := tokens.SearchOffset(offset)
	if sourceToken == nil {
		return nil, nil // No token at the given position
	}
	target := s.findSourceAstNode(ctx, sourceToken)
	if target == nil {
		return nil, nil // No AST node associated with the token
	}
	nameUnit := linking.Name(target)
	if nameUnit == nil {
		return nil, nil // No name token for the target node
	}
	locations := []lsp.Location{
		// Include the definition location itself
		{
			URI:   target.Document().URI.DocumentURI(),
			Range: nameUnit.Segment().Range.LspRange(),
		},
	}
	// Iterate through all documents and collect references to the symbol
	for doc := range documentManager.All() {
		refDescriptions := doc.ReferenceDescriptions.ForTarget(target)
		for refDesc := range refDescriptions {
			location := lsp.Location{
				URI:   refDesc.SourceURI().DocumentURI(),
				Range: refDesc.Segment.Range.LspRange(),
			}
			locations = append(locations, location)
		}
	}
	return locations, nil
}

func (s *DefaultReferencesProvider) findSourceAstNode(ctx context.Context, token *core.Token) core.AstNode {
	ref := core.ReferenceOfToken(token)
	if ref != nil {
		return ref.RefNode(ctx)
	} else {
		node := token.Owner()
		if node == nil {
			return nil
		}
		nameUnit := linking.Name(node)
		if nameUnit == nil {
			return nil
		}
		segment := nameUnit.Segment()
		if token.TextSegment.Indices.Start < segment.Indices.Start || token.TextSegment.Indices.End > segment.Indices.End {
			return nil // The token is not within the name segment
		}
		return node
	}
}
