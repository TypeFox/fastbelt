// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/lsp"
)

// GrammarNodeKindProvider provides Fastbelt grammar-specific symbol kinds.
type GrammarNodeKindProvider struct {
	server.DefaultNodeKindProvider
}

func (p *GrammarNodeKindProvider) GetSymbolKind(node core.AstNode) lsp.SymbolKind {
	switch node.(type) {
	case *GrammarImpl:
		return lsp.File
	case *InterfaceImpl:
		return lsp.Interface
	case *ParserRuleImpl:
		return lsp.Function
	case *FieldImpl:
		return lsp.Property
	case *TokenImpl:
		return lsp.Constant
	default:
		return p.DefaultNodeKindProvider.GetSymbolKind(node)
	}
}
