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

// DefinitionProvider is a service for handling LSP definition requests.
type DefinitionProvider interface {
	// TODO: Maybe add the document directly to the params to avoid looking it up in the workspace?
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
	first, second := tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil // No token at the given position
	}
	nameFinder := service.MustGet[NameFinder](s.sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, nil // Could not find a name
	}
	target := foundName.Target
	targetNode := target.Owner()
	sourceRange := foundName.Source.Segment().Range.LspRange()
	fullRange := targetNode.Segment().Range.LspRange()
	targetRange := target.Segment().Range.LspRange()
	link := lsp.DefinitionLink{
		OriginSelectionRange: &sourceRange,
		TargetURI:            targetNode.Document().URI.DocumentURI(),
		TargetRange:          fullRange,
		TargetSelectionRange: targetRange,
	}
	return []lsp.DefinitionLink{link}, nil
}
