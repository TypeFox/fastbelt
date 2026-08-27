// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// CodeLensProvider is a service for handling LSP code lens requests.
//
// Usage:
//
//	type MyCodeLensProvider struct{ sc *service.Container }
//
//	func (p *MyCodeLensProvider) HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
//	    // Analyze document and return inline commands/info
//	    return []lsp.CodeLens{
//	        {
//	            Range: lsp.Range{...},
//	            Command: &lsp.Command{
//	                Title:   "3 references",
//	                Command: "showReferences",
//	            },
//	        },
//	    }, nil
//	}
type CodeLensProvider interface {
	HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error)
}

// ResolvingCodeLensProvider is an optional extension of [CodeLensProvider]
// for providers that defer work until a code lens actually becomes visible,
// via the LSP codeLens/resolve request.
//
// If a registered CodeLensProvider also implements this interface, the
// server advertises resolve support and routes codeLens/resolve requests to
// HandleCodeLensResolveRequest. Providers that don't need lazy resolution
// can simply not implement it.
//
// Usage:
//
//	func (p *MyCodeLensProvider) HandleCodeLensResolveRequest(ctx context.Context, lens *lsp.CodeLens) (*lsp.CodeLens, error) {
//	    // Use lens.Data (set in HandleCodeLensRequest) to compute the
//	    // deferred Command now that the lens is actually visible.
//	    lens.Command = &lsp.Command{Title: "..."}
//	    return lens, nil
//	}
type ResolvingCodeLensProvider interface {
	CodeLensProvider
	HandleCodeLensResolveRequest(ctx context.Context, lens *lsp.CodeLens) (*lsp.CodeLens, error)
}
