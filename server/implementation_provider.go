// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// ImplementationProvider is a service for handling LSP implementation requests.
type ImplementationProvider interface {
	HandleImplementationRequest(ctx context.Context, params *lsp.ImplementationParams) ([]lsp.DefinitionLink, error)
}
