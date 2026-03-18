// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

type ReferencesProvider interface {
	HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error)
}

type DefaultReferencesProvider struct {
	srv ServerSrvCont
}

func NewDefaultReferencesProvider(srv ServerSrvCont) ReferencesProvider {
	return &DefaultReferencesProvider{srv: srv}
}

func (rp *DefaultReferencesProvider) HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error) {
	uri := core.ParseURI(string(params.TextDocument.URI))
	targetDoc := rp.srv.Workspace().DocumentManager.Get(uri)
	if targetDoc == nil {
		return nil, nil // Document not found
	}
	offset := targetDoc.TextDoc.OffsetAt(params.Position)
	tokens := targetDoc.Tokens
	sourceToken := tokens.SearchOffset(offset)
	if sourceToken == nil {
		return nil, nil // No token at the given position
	}
	namer := rp.srv.Linking().Namer
	target := rp.findSourceAstNode(ctx, sourceToken)
	if target == nil {
		return nil, nil // No AST node associated with the token
	}
	_, nameToken := namer.Name(target)
	if nameToken == nil {
		return nil, nil // No name token for the target node
	}
	locations := []lsp.Location{
		// Include the definition location itself
		{
			URI:   target.Document().URI.DocumentURI(),
			Range: nameToken.Segment.Range.LspRange(),
		},
	}
	documentManager := rp.srv.Workspace().DocumentManager
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

func (rp *DefaultReferencesProvider) findSourceAstNode(ctx context.Context, token *core.Token) core.AstNode {
	ref := core.ReferenceOfToken(token)
	if ref != nil {
		return ref.RefNode(ctx)
	} else {
		node := token.Element
		if node == nil {
			return nil
		}
		namer := rp.srv.Linking().Namer
		_, nameToken := namer.Name(node)
		if nameToken == nil || nameToken != token {
			return nil // The token at the position is not the name token
		}
		return node
	}
}
