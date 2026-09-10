// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// DocumentLinkProvider is a service for handling LSP document link requests.
//
// Usage:
//
//	type MyDocumentLinkProvider struct{ sc *service.Container }
//
//	func (p *MyDocumentLinkProvider) HandleDocumentLinkRequest(ctx context.Context, params *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error) {
//	    // Find URLs, file references, etc. in the document
//	    return []lsp.DocumentLink{
//	        {
//	            Range:  lsp.Range{...},
//	            Target: "file:///path/to/file.txt",
//	        },
//	    }, nil
//	}
type DocumentLinkProvider interface {
	HandleDocumentLinkRequest(ctx context.Context, params *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error)
}

// ResolvingDocumentLinkProvider is an optional extension of
// [DocumentLinkProvider] for providers that defer work until a document
// link actually becomes visible, via the LSP documentLink/resolve request.
//
// If a registered DocumentLinkProvider also implements this interface, the
// server advertises resolve support and routes documentLink/resolve
// requests to HandleDocumentLinkResolveRequest. Providers that don't need
// lazy resolution can simply not implement it.
//
// Usage:
//
//	func (p *MyDocumentLinkProvider) HandleDocumentLinkResolveRequest(ctx context.Context, link *lsp.DocumentLink) (*lsp.DocumentLink, error) {
//	    // Use link.Data (set in HandleDocumentLinkRequest) to compute the
//	    // deferred Target now that the link is actually visible.
//	    target := lsp.URI("file:///path/to/file.txt")
//	    link.Target = &target
//	    return link, nil
//	}
type ResolvingDocumentLinkProvider interface {
	DocumentLinkProvider
	HandleDocumentLinkResolveRequest(ctx context.Context, link *lsp.DocumentLink) (*lsp.DocumentLink, error)
}
