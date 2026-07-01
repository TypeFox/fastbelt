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

// TypeDefinitionProvider is a service for handling LSP type definition requests.
type TypeDefinitionProvider interface {
	HandleTypeDefinitionRequest(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error)
}

// DefaultTypeDefinitionProvider is the default implementation of [TypeDefinitionProvider].
// The default implementation provides basic support by looking for type references
// in the AST node's cross-references. Language adopters may need to override this
// service to handle language-specific type inference or implicit type relationships.
type DefaultTypeDefinitionProvider struct {
	sc *service.Container
}

func NewDefaultTypeDefinitionProvider(sc *service.Container) TypeDefinitionProvider {
	return &DefaultTypeDefinitionProvider{sc: sc}
}

func (p *DefaultTypeDefinitionProvider) HandleTypeDefinitionRequest(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error) {
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
	if foundName.Source == nil {
		return nil, nil
	}

	sourceNode := foundName.Source.Owner()
	sourceRange := foundName.Source.Segment().Range.LspRange()

	var results []lsp.DefinitionLink

	// Look for type references in the node's cross-references
	for xref := range core.References(sourceNode) {
		// Check if this is a type reference (many languages have a "type" feature)
		// This is a heuristic and may need language-specific customization
		targetDesc := xref.Description()
		if targetDesc != nil && targetDesc.Node != nil {
			targetNode := targetDesc.Node

			// Check if the target looks like a type definition
			// (This is a simple heuristic; language-specific logic would be better)
			fullRange := targetNode.Segment().Range.LspRange()

			// Try to get a more specific range if available
			targetRange := fullRange
			if targetDesc.Name != nil {
				targetRange = targetDesc.Name.Segment().Range.LspRange()
			}

			link := lsp.DefinitionLink{
				OriginSelectionRange: &sourceRange,
				TargetURI:            targetNode.Document().URI.DocumentURI(),
				TargetRange:          fullRange,
				TargetSelectionRange: targetRange,
			}
			results = append(results, link)
		}
	}

	return results, nil
}
