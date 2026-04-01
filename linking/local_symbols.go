// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"
	"slices"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

// LocalSymbolsProvider computes local symbol information for documents.
type LocalSymbolsProvider interface {
	// Provide traverses the document's AST and computes the local symbol table.
	// The result is stored in the document's LocalSymbols field.
	// The caller must hold the document's write lock.
	Provide(ctx context.Context, document *core.Document)
}

// DefaultLocalSymbolsProvider is the default implementation of LocalSymbolsProvider.
type DefaultLocalSymbolsProvider struct {
	srv LinkingSrvCont
}

func NewDefaultLocalSymbolsProvider(srv LinkingSrvCont) LocalSymbolsProvider {
	return &DefaultLocalSymbolsProvider{
		srv: srv,
	}
}

func (s *DefaultLocalSymbolsProvider) Provide(ctx context.Context, document *core.Document) {
	root := document.Root
	symbols := collections.NewMultiMap[core.AstNode, *core.AstNodeDescription]()

	for node := range core.AllChildren(root) {
		desc, container := s.srv.Linking().LocalSymbolDescriber.Describe(node)
		if desc != nil {
			symbols.Put(container, desc)
		}
	}
	document.LocalSymbols = NewDefaultLocalSymbolsFromMap(symbols)
}

// DefaultLocalSymbols is the default implementation of core.LocalSymbols.
type DefaultLocalSymbols struct {
	Symbols collections.MultiMap[core.AstNode, *core.AstNodeDescription]
}

func NewDefaultLocalSymbols() *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: collections.NewMultiMap[core.AstNode, *core.AstNodeDescription](),
	}
}

func NewDefaultLocalSymbolsFromMap(m collections.MultiMap[core.AstNode, *core.AstNodeDescription]) *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: m,
	}
}

func (ls *DefaultLocalSymbols) Has(node core.AstNode) bool {
	return ls.Symbols.Has(node)
}

func (ls *DefaultLocalSymbols) Iter(node core.AstNode) core.SymbolList {
	symbols := ls.Symbols.Get(node)
	return slices.Values(symbols)
}

// LocalSymbolDescriber describes how symbols can be referenced in the same document.
type LocalSymbolDescriber interface {
	// Describe determines the name and other metadata of a locally visible symbol.
	// It returns the description and the container in which the symbol is visible, or nil if the symbol is not locally visible.
	Describe(node core.AstNode) (*core.AstNodeDescription, core.AstNode)
}

// DefaultLocalSymbolDescriber is the default implementation of LocalSymbolTableItemProvider.
type DefaultLocalSymbolDescriber struct {
	srv LinkingSrvCont
}

func NewDefaultLocalSymbolDescriber(srv LinkingSrvCont) LocalSymbolDescriber {
	return &DefaultLocalSymbolDescriber{
		srv: srv,
	}
}

func (p *DefaultLocalSymbolDescriber) Describe(node core.AstNode) (*core.AstNodeDescription, core.AstNode) {
	container := node.Container()
	name, nameToken := p.srv.Linking().Namer.Name(node)
	if name != "" {
		var segment *core.TextSegment
		if nameToken != nil {
			segment = &nameToken.Segment
		}
		desc := core.NewAstNodeDescription(node, name, segment, node.Segment())
		return desc, container
	}
	return nil, nil
}
