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

// ImplementationFilter provides language-specific filtering for implementation detection.
//
// Usage:
//
//	type MyImplementationFilter struct{}
//
//	func (f *MyImplementationFilter) ShouldInclude(refDesc *core.ReferenceDescription, targetNode core.AstNode) bool {
//	    // Check if reference is in an "implements" or "extends" clause
//	    return refDesc.Field == "implements" || refDesc.Field == "extends"
//	}
//
//	// Register with custom filter
//	service.Put(sc, server.NewImplementationProviderWithFilter(sc, &MyImplementationFilter{}))
type ImplementationFilter interface {
	ShouldInclude(refDesc *core.ReferenceDescription, targetNode core.AstNode) bool
}

// DefaultImplementationFilter is the default implementation of [ImplementationFilter].
// It includes all references as potential implementations, providing a conservative
// approximation that may include false positives.
type DefaultImplementationFilter struct{}

func (f *DefaultImplementationFilter) ShouldInclude(refDesc *core.ReferenceDescription, targetNode core.AstNode) bool {
	return true
}

// ImplementationProvider is a service for handling LSP implementation requests.
type ImplementationProvider interface {
	HandleImplementationRequest(ctx context.Context, params *lsp.ImplementationParams) ([]lsp.DefinitionLink, error)
}

// DefaultImplementationProvider is the default implementation of [ImplementationProvider].
//
// The default implementation finds all references to the target symbol and uses an
// [ImplementationFilter] to determine which references represent actual implementations.
type DefaultImplementationProvider struct {
	sc     *service.Container
	filter ImplementationFilter
}

// NewDefaultImplementationProvider creates an implementation provider with the default filter.
func NewDefaultImplementationProvider(sc *service.Container) ImplementationProvider {
	return &DefaultImplementationProvider{
		sc:     sc,
		filter: &DefaultImplementationFilter{},
	}
}

// NewImplementationProviderWithFilter creates an implementation provider with a custom filter.
// The filter determines which references should be considered implementations.
func NewImplementationProviderWithFilter(sc *service.Container, filter ImplementationFilter) ImplementationProvider {
	return &DefaultImplementationProvider{
		sc:     sc,
		filter: filter,
	}
}

func (p *DefaultImplementationProvider) HandleImplementationRequest(ctx context.Context, params *lsp.ImplementationParams) ([]lsp.DefinitionLink, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil
	}

	offset := doc.TextDoc.OffsetAt(params.Position)
	first, second := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil
	}

	nameFinder := service.MustGet[NameFinder](p.sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil {
		return nil, nil
	}

	targetNode := foundName.Target.Owner()

	referencesFinder, err := service.Get[ReferencesFinder](p.sc)
	if err != nil {
		// No references finder available
		return nil, nil
	}

	var results []lsp.DefinitionLink
	sourceRange := foundName.Source.Segment().Range.LspRange()

	for refDesc := range referencesFinder.Find(ctx, targetNode, FindReferencesOptions{
		IncludeDeclaration: false,
	}) {
		if !p.filter.ShouldInclude(refDesc, targetNode) {
			continue
		}

		link := lsp.DefinitionLink{
			OriginSelectionRange: &sourceRange,
			TargetURI:            refDesc.SourceURI().DocumentURI(),
			TargetRange:          refDesc.Segment.Range.LspRange(),
			TargetSelectionRange: refDesc.Segment.Range.LspRange(),
		}
		results = append(results, link)
	}

	return results, nil
}
