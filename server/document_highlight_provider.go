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

// DocumentHighlightProvider is a service for handling LSP document highlight requests.
type DocumentHighlightProvider interface {
	HandleDocumentHighlightRequest(ctx context.Context, params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error)
}

// Ensure DefaultDocumentHighlightProvider implements the expected interfaces.
var _ DocumentHighlightProvider = (*DefaultDocumentHighlightProvider)(nil)
var _ DocumentStateRequirements = (*DefaultDocumentHighlightProvider)(nil)

// DefaultDocumentHighlightProvider is the default implementation of [DocumentHighlightProvider].
// It highlights all references to the symbol at the cursor position that are
// located in the same document, using [NameFinder] and [ReferencesFinder].
type DefaultDocumentHighlightProvider struct {
	sc *service.Container
}

// NewDefaultDocumentHighlightProvider creates a new [DefaultDocumentHighlightProvider].
func NewDefaultDocumentHighlightProvider(sc *service.Container) DocumentHighlightProvider {
	return &DefaultDocumentHighlightProvider{sc: sc}
}

// The default document highlight provider requires the document in question have its references indexed.
func (s *DefaultDocumentHighlightProvider) RequiredState() (core.DocumentState, bool) {
	return core.DocStateReferences, false
}

func (s *DefaultDocumentHighlightProvider) HandleDocumentHighlightRequest(ctx context.Context, params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
	documentManager := s.sc.MustGet[workspace.DocumentManager]()
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
	nameFinder := s.sc.MustGet[NameFinder]()
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, nil // Could not find a name
	}
	target := foundName.Target.Owner()
	referencesFinder := s.sc.MustGet[ReferencesFinder]()
	highlights := []lsp.DocumentHighlight{}
	for refDesc := range referencesFinder.Find(ctx, target, FindReferencesOptions{
		IncludeDeclaration: true,
		TargetURI:          uri,
	}) {
		textDoc := refDesc.SourceNode.Document().TextDoc
		highlight := lsp.DocumentHighlight{
			Range: refDesc.Range.LspRange(textDoc),
		}
		highlights = append(highlights, highlight)
	}
	return highlights, nil
}
