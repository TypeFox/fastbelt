// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// CodeActionProvider is a service for handling LSP code action requests.
//
// Usage:
//
//	type MyCodeActionProvider struct{ sc *service.Container }
//
//	func (p *MyCodeActionProvider) HandleCodeActionRequest(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
//	    // Analyze diagnostics and context, return quick fixes
//	    return []lsp.CodeAction{
//	        {
//	            Title: "Add missing import",
//	            Kind:  lsp.QuickFix,
//	            Edit:  &lsp.WorkspaceEdit{...},
//	        },
//	    }, nil
//	}
type CodeActionProvider interface {
	HandleCodeActionRequest(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error)
}

// DefaultCodeActionProvider returns no code actions.
type DefaultCodeActionProvider struct {
	sc *service.Container
}

func NewDefaultCodeActionProvider(sc *service.Container) CodeActionProvider {
	return &DefaultCodeActionProvider{sc: sc}
}

func (p *DefaultCodeActionProvider) HandleCodeActionRequest(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	return []lsp.CodeAction{}, nil
}
