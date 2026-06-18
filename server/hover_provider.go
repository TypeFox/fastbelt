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

// HoverProvider is a service for handling LSP hover requests.
type HoverProvider interface {
	HandleHoverRequest(ctx context.Context, params *lsp.HoverParams) (*lsp.Hover, error)
}

// DefaultHoverProvider is the default implementation of [HoverProvider].
type DefaultHoverProvider struct {
	sc *service.Container
}

func NewDefaultHoverProvider(sc *service.Container) HoverProvider {
	return &DefaultHoverProvider{sc: sc}
}

func (s *DefaultHoverProvider) HandleHoverRequest(ctx context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	documentManager := service.MustGet[workspace.DocumentManager](s.sc)
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

	nameFinder := service.MustGet[NameFinder](s.sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, nil
	}

	targetNode := foundName.Target.Owner()
	docProvider, err := service.Get[DocumentationProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	content := docProvider.Documentation(targetNode)
	if content == "" {
		return nil, nil
	}

	sourceRange := foundName.Source.Segment().Range.LspRange()
	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: content,
		},
		Range: sourceRange,
	}, nil
}
