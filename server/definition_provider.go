// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

type DefinitionProvider interface {
	// TODO: Maybe add the document directly to the params to avoid looking it up again in the workspace?
	// Also, maybe add a separate params struct that doesn't directly depend on the lsp lib
	// TODO: Use `LocationLink` instead of `Location` to support more advanced scenarios
	// Requires a change in the LSP library
	HandleDefinitionRequest(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.Location, error)
}

type DefaultDefinitionProvider struct {
	srv ServerSrvCont
}

func NewDefaultDefinitionProvider(srv ServerSrvCont) DefinitionProvider {
	return &DefaultDefinitionProvider{srv: srv}
}

func (dp *DefaultDefinitionProvider) HandleDefinitionRequest(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.Location, error) {
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := dp.srv.Workspace().DocumentManager.Get(uri)
	if doc == nil {
		return nil, nil // Document not found
	}
	offset := doc.TextDoc.OffsetAt(params.Position)
	tokens := doc.Tokens
	sourceToken := tokens.SearchOffset(offset)
	if sourceToken == nil {
		return nil, nil // No token at the given position
	}
	ref := core.ReferenceOfToken(sourceToken)
	if ref != nil {
		// The token at the position is a reference
		// Try to resolve it and return the location of the target symbol
		return dp.fromReference(ref), nil
	} else {
		// The token might still be the name of a symbol
		// In this case, we want to return the location of the symbol itself
		return dp.fromName(sourceToken), nil
	}
}

func (dp *DefaultDefinitionProvider) fromReference(ref core.UntypedReference) []lsp.Location {
	target := ref.Description()
	if target == nil || target.NameSegment == nil {
		return nil // No target description
	}
	link := lsp.Location{
		URI:   target.URI.DocumentURI(),
		Range: target.NameSegment.Range.LspRange(),
	}
	return []lsp.Location{link}
}

func (dp *DefaultDefinitionProvider) fromName(token *core.Token) []lsp.Location {
	target := token.Owner()
	if target == nil {
		return nil
	}
	namer := dp.srv.Linking().Namer
	nameUnit := namer.Name(target)
	if nameUnit == nil {
		return nil
	}
	segment := nameUnit.Segment()
	if token.TextSegment.Indices.Start < segment.Indices.Start || token.TextSegment.Indices.End > segment.Indices.End {
		return nil // The token is not within the name segment
	}
	link := lsp.Location{
		URI:   target.Document().URI.DocumentURI(),
		Range: segment.Range.LspRange(),
	}
	return []lsp.Location{link}
}
