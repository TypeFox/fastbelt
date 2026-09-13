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
	target, sourceRange, ok := ResolveHoverTarget(ctx, s.sc, params)
	if !ok {
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

// ResolveHoverTarget resolves the AST node referenced at the hover position
// in params - the declaration itself, or any name/reference that resolves
// to it - along with the source range of that name/reference. ok is false
// when there's nothing to hover at that position (no document, no token
// there, or no resolvable name), in which case callers should return
// (nil, nil) from their own HandleHoverRequest.
//
// It's exported so a language that needs a custom HoverProvider - e.g. to
// show content beyond documentation comments - can reuse this resolution
// step (built on [NameFinder]) instead of duplicating it. See
// internal/grammarhover for an example: it composes documentation with a
// railroad syntax diagram, a concept specific to fastbelt's own grammar
// language that has no place in this generic package.
func ResolveHoverTarget(ctx context.Context, sc *service.Container, params *lsp.HoverParams) (target core.AstNode, sourceRange lsp.Range, ok bool) {
	documentManager := service.MustGet[workspace.DocumentManager](sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, lsp.Range{}, false
	}

	offset := doc.TextDoc.OffsetAt(params.Position)
	first, second := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, lsp.Range{}, false
	}

	nameFinder := service.MustGet[NameFinder](sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, lsp.Range{}, false
	}

	return foundName.Target.Owner(), foundName.Source.TextRange().LspRange(doc.TextDoc), true
}
