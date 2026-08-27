// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// DeclarationProvider is a service for handling LSP declaration requests.
// Implement this service if your language distinguishes between declaration
// and definition.
type DeclarationProvider interface {
	HandleDeclarationRequest(ctx context.Context, params *lsp.DeclarationParams) ([]lsp.DefinitionLink, error)
}
