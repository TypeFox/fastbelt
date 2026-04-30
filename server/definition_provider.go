// Copyright 2026 TypeFox GmbH
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

// DefinitionProvider is a service for handling LSP definition requests.
type DefinitionProvider interface {
	// TODO: Maybe add the document directly to the params to avoid looking it up again in the workspace?
	// Also, maybe add a separate params struct that doesn't directly depend on the lsp lib
	HandleDefinitionRequest(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.DefinitionLink, error)
}

// DefaultDefinitionProvider is the default implementation of [DefinitionProvider].
type DefaultDefinitionProvider struct {
	sc *service.Container
}

func NewDefaultDefinitionProvider(sc *service.Container) DefinitionProvider {
	return &DefaultDefinitionProvider{sc: sc}
}

func (s *DefaultDefinitionProvider) HandleDefinitionRequest(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.DefinitionLink, error) {
	documentManager := service.MustGet[workspace.DocumentManager](s.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
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
		return s.fromReference(ref), nil
	} else {
		// The token might still be the name of a symbol
		// In this case, we want to return the location of the symbol itself
		return s.fromName(sourceToken), nil
	}
}

func (s *DefaultDefinitionProvider) fromReference(ref core.UntypedReference) []lsp.DefinitionLink {
	target := ref.Description()
	if target == nil || target.NameSegment == nil {
		return nil // No target description
	}
	var originRange *lsp.Range
	if originSegment := ref.Segment(); originSegment != nil {
		lspRange := originSegment.Range.LspRange()
		originRange = &lspRange
	}
	fullRange := target.FullSegment.Range.LspRange()
	nameRange := target.NameSegment.Range.LspRange()
	link := lsp.DefinitionLink{
		OriginSelectionRange: originRange,
		TargetURI:            target.URI.DocumentURI(),
		TargetRange:          fullRange,
		TargetSelectionRange: nameRange,
	}
	return []lsp.DefinitionLink{link}
}

func (s *DefaultDefinitionProvider) fromName(token *core.Token) []lsp.DefinitionLink {
	target := token.Element
	if target == nil {
		return nil
	}
	sourceSegment := &token.TextSegment
	if stringNode, ok := token.Element.(core.CompositeNode); ok {
		// If the token is part of a string node, the actual name segment is the parent node
		target = stringNode.Container()
		sourceSegment = stringNode.Segment()
	}
	targetSegment := target.Segment()
	if targetSegment == nil {
		return nil
	}
	nameUnit := linking.Name(target)
	if nameUnit == nil {
		return nil
	}
	segment := nameUnit.Segment()
	if segment == nil || token.TextSegment.Indices.Start < segment.Indices.Start || token.TextSegment.Indices.End > segment.Indices.End {
		return nil // The token is not within the name segment
	}
	fullRange := targetSegment.Range.LspRange()
	targetRange := segment.Range.LspRange()
	var sourceRange *lsp.Range
	if sourceSegment != nil {
		lspRange := sourceSegment.Range.LspRange()
		sourceRange = &lspRange
	}
	link := lsp.DefinitionLink{
		OriginSelectionRange: sourceRange,
		TargetURI:            target.Document().URI.DocumentURI(),
		TargetRange:          fullRange,
		TargetSelectionRange: targetRange,
	}
	return []lsp.DefinitionLink{link}
}
