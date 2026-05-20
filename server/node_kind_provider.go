// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

// NodeKindProvider consolidates the logic for gathering LSP kind information
// based on AST nodes.
type NodeKindProvider interface {
	// GetSymbolKind returns a SymbolKind as used by WorkspaceSymbolProvider
	// or DocumentSymbolProvider.
	GetSymbolKind(node core.AstNode) lsp.SymbolKind

	// GetCompletionItemKind returns a CompletionItemKind as used by the
	// CompletionProvider.
	// TODO: Implement when completion provider is added
	// GetCompletionItemKind(node core.AstNode) lsp.CompletionItemKind
}

// DefaultNodeKindProvider provides standard kind mapping behavior.
// This implementation returns lsp.Field for all nodes.
// Languages should extend this to provide specific kinds based on node types.
type DefaultNodeKindProvider struct{}

func (p *DefaultNodeKindProvider) GetSymbolKind(node core.AstNode) lsp.SymbolKind {
	// Default to Field for all nodes
	// Languages should override this to provide specific kinds based on node type
	return lsp.Field
}
