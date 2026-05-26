// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"typefox.dev/lsp"
)

// SymbolKind implementations for grammar AST nodes.
// These methods are called by server.SymbolKind() to determine
// the LSP symbol kind for each node in the document symbol outline.

func (g *GrammarImpl) SymbolKind() lsp.SymbolKind {
	return lsp.File
}

func (i *InterfaceImpl) SymbolKind() lsp.SymbolKind {
	return lsp.Interface
}

func (p *ParserRuleImpl) SymbolKind() lsp.SymbolKind {
	return lsp.Function
}

func (f *FieldImpl) SymbolKind() lsp.SymbolKind {
	return lsp.Property
}

func (t *TokenImpl) SymbolKind() lsp.SymbolKind {
	return lsp.Constant
}
