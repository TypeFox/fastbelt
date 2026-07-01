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
