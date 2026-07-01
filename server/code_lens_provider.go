// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/fastbelt/util/service"
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

// DefaultCodeLensProvider returns no code lenses.
type DefaultCodeLensProvider struct {
	sc *service.Container
}

func NewDefaultCodeLensProvider(sc *service.Container) CodeLensProvider {
	return &DefaultCodeLensProvider{sc: sc}
}

func (p *DefaultCodeLensProvider) HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	return []lsp.CodeLens{}, nil
}
