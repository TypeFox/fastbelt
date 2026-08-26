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

// ReferencesProvider is a service for handling LSP reference requests.
type ReferencesProvider interface {
	HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error)
}

// DefaultReferencesProvider is the default implementation of [ReferencesProvider].
type DefaultReferencesProvider struct {
	sc *service.Container
}

// Ensure DefaultReferencesProvider implements the expected interfaces.
var _ ReferencesProvider = (*DefaultReferencesProvider)(nil)
var _ DocumentStateRequirements = (*DefaultReferencesProvider)(nil)

// RequiredState indicates that the default references provider requires the workspace
// to have reference information collected.
func (rf *DefaultReferencesProvider) RequiredState() (core.DocumentState, bool) {
	return core.DocStateReferences, true
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
	first, second := tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil // No token at the given position
	}
	nameFinder := service.MustGet[NameFinder](s.sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, nil // Could not find a name
	}
	target := foundName.Target.Owner()
	referencesFinder := service.MustGet[ReferencesFinder](s.sc)
	locations := []lsp.Location{}
	for refDesc := range referencesFinder.Find(ctx, target, FindReferencesOptions{
		IncludeDeclaration: true,
	}) {
		textDoc := refDesc.SourceNode.Document().TextDoc
		location := lsp.Location{
			URI:   refDesc.SourceURI().DocumentURI(),
			Range: refDesc.Range.LspRange(textDoc),
		}
		locations = append(locations, location)
	}
	return locations, nil
}
