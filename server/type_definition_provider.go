// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// TypeDefinitionProvider is a service for handling LSP type definition requests.
//
// Usage:
//
//	type MyTypeDefinitionProvider struct{ sc *service.Container }
//
//	func (p *MyTypeDefinitionProvider) HandleTypeDefinitionRequest(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error) {
//	    // Resolve the symbol under the cursor and return its type's declaration.
//	    return links, nil
//	}
//
//	// Register with the service container
//	service.Put[server.TypeDefinitionProvider](sc, &MyTypeDefinitionProvider{sc: sc})
type TypeDefinitionProvider interface {
	HandleTypeDefinitionRequest(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error)
}
