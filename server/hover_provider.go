// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/fastbelt/util/service"
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
	target, sourceRange := TargetAtCursor(ctx, s.sc, &params.TextDocumentPositionParams)
	if target == nil {
		return nil, nil
	}

	docProvider, err := service.Get[DocumentationProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	content := docProvider.Documentation(target)
	if content == "" {
		return nil, nil
	}

	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: content,
		},
		Range: sourceRange,
	}, nil
}
