// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

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

// ResolvingCodeActionProvider is an optional extension of [CodeActionProvider]
// for providers that defer work until a code action is actually selected,
// via the LSP codeAction/resolve request.
//
// If a registered CodeActionProvider also implements this interface, the
// server advertises resolve support and routes codeAction/resolve requests
// to HandleCodeActionResolveRequest. Providers that don't need lazy
// resolution can simply not implement it.
//
// Usage:
//
//	func (p *MyCodeActionProvider) HandleCodeActionResolveRequest(ctx context.Context, action *lsp.CodeAction) (*lsp.CodeAction, error) {
//	    // Use action.Data (set in HandleCodeActionRequest) to compute the
//	    // deferred Edit now that the action was actually selected.
//	    action.Edit = &lsp.WorkspaceEdit{...}
//	    return action, nil
//	}
type ResolvingCodeActionProvider interface {
	CodeActionProvider
	HandleCodeActionResolveRequest(ctx context.Context, action *lsp.CodeAction) (*lsp.CodeAction, error)
}
