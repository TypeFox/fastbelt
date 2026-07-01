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

// TypeHierarchyProvider is a service for handling LSP type hierarchy requests.
type TypeHierarchyProvider interface {
	HandlePrepareTypeHierarchyRequest(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error)
	HandleSupertypesRequest(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error)
	HandleSubtypesRequest(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error)
}

// DefaultTypeHierarchyProvider is the default implementation of [TypeHierarchyProvider].
// It delegates type relationship detection to [TypeHierarchyContributor].
type DefaultTypeHierarchyProvider struct {
	sc *service.Container
}

func NewDefaultTypeHierarchyProvider(sc *service.Container) TypeHierarchyProvider {
	return &DefaultTypeHierarchyProvider{sc: sc}
}

func (p *DefaultTypeHierarchyProvider) HandlePrepareTypeHierarchyRequest(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil
	}

	contributor, err := service.Get[TypeHierarchyContributor](p.sc)
	if err != nil {
		// No contributor available
		return nil, nil
	}

	// Find token at cursor position
	offset := doc.TextDoc.OffsetAt(params.Position)
	first, _ := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil
	}

	// Get the AST node for the token
	node := first.Element

	// Delegate to contributor
	item := contributor.PrepareItem(ctx, first, node)
	if item == nil {
		return nil, nil
	}

	return []lsp.TypeHierarchyItem{*item}, nil
}

func (p *DefaultTypeHierarchyProvider) HandleSupertypesRequest(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error) {
	contributor, err := service.Get[TypeHierarchyContributor](p.sc)
	if err != nil {
		// No contributor available
		return nil, nil
	}

	supertypes := contributor.FindSupertypes(ctx, &params.Item)

	return supertypes, nil
}

func (p *DefaultTypeHierarchyProvider) HandleSubtypesRequest(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error) {
	contributor, err := service.Get[TypeHierarchyContributor](p.sc)
	if err != nil {
		// No contributor available
		return nil, nil
	}

	subtypes := contributor.FindSubtypes(ctx, &params.Item)

	return subtypes, nil
}
