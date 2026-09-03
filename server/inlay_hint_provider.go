// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// InlayHintProvider is a service for handling LSP inlay hint requests.
//
// Usage:
//
//	type MyInlayHintProvider struct{ sc *service.Container }
//
//	func (p *MyInlayHintProvider) HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
//	    documentManager := service.MustGet[workspace.DocumentManager](p.sc)
//	    doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
//	    if doc == nil {
//	        return nil, nil
//	    }
//	    var hints []lsp.InlayHint
//	    for node := range server.NodesInRange(doc, params.Range) {
//	        if funcCall, ok := node.(*ast.FunctionCall); ok {
//	            hints = append(hints, lsp.InlayHint{
//	                Position: funcCall.TextRange().LspRange(doc.TextDoc).End,
//	                Label:    []lsp.InlayHintLabelPart{{Value: "paramName:"}},
//	                Kind:     lsp.Parameter,
//	            })
//	        }
//	    }
//	    return hints, nil
//	}
//
//	// Register with the service container
//	service.Put[server.InlayHintProvider](sc, &MyInlayHintProvider{sc: sc})
type InlayHintProvider interface {
	HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error)
}

// ResolvingInlayHintProvider is an optional extension of [InlayHintProvider]
// for providers that defer work until an inlay hint actually becomes
// visible, via the LSP inlayHint/resolve request.
//
// If a registered InlayHintProvider also implements this interface, the
// server advertises resolve support and routes inlayHint/resolve requests to
// HandleInlayHintResolveRequest. Providers that don't need lazy resolution
// can simply not implement it.
//
// Usage:
//
//	func (p *MyInlayHintProvider) HandleInlayHintResolveRequest(ctx context.Context, hint *lsp.InlayHint) (*lsp.InlayHint, error) {
//	    // Use hint.Data (set in HandleInlayHintRequest) to compute the
//	    // deferred TextEdits now that the hint is actually visible.
//	    hint.TextEdits = []lsp.TextEdit{...}
//	    return hint, nil
//	}
type ResolvingInlayHintProvider interface {
	InlayHintProvider
	HandleInlayHintResolveRequest(ctx context.Context, hint *lsp.InlayHint) (*lsp.InlayHint, error)
}
