// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

// SymbolKindNode can be implemented by AST node Impl structs to provide
// custom LSP symbol kinds for the document symbol outline.
type SymbolKindNode interface {
	// SymbolKind returns the LSP symbol kind for this node.
	SymbolKind() lsp.SymbolKind
}

// SymbolKind returns the LSP symbol kind for a node.
// If the node implements SymbolKindNode, uses that.
// Otherwise returns lsp.Field as default.
//
// Language-specific implementations can be provided by implementing the [SymbolKindNode] interface.
func SymbolKind(node core.AstNode) lsp.SymbolKind {
	if sk, ok := node.(SymbolKindNode); ok {
		return sk.SymbolKind()
	}
	return lsp.Field // Default fallback
}
