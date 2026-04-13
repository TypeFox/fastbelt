// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
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
	symbolContainers := s.srv.Generated().SymbolContainers
	localSymbols := NewDefaultLocalSymbols()
	containers := localSymbols.Symbols

	for node := range core.AllChildren(root) {
		desc, containerNode := s.srv.Linking().LocalSymbolDescriber.Describe(node)
		if desc != nil {
			if symbols, exists := containers[containerNode]; exists {
				// Fast path for existing container: just put the new description into it
				symbols.Put(desc)
			} else {
				// No container yet for this node, create a new one
				newContainer := symbolContainers.New()
				if newContainer.Put(desc) {
					// Only add the new container if the description was successfully added
					containers[containerNode] = newContainer
				}
			}
		}
	}
	document.LocalSymbols = localSymbols
}

// DefaultLocalSymbols is the default implementation of core.LocalSymbols.
type DefaultLocalSymbols struct {
	Symbols map[core.AstNode]core.SymbolContainer
}

func NewDefaultLocalSymbols() *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: make(map[core.AstNode]core.SymbolContainer),
	}
}

func NewDefaultLocalSymbolsFromMap(m map[core.AstNode]core.SymbolContainer) *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: m,
	}
}

func (ls *DefaultLocalSymbols) For(node core.AstNode) core.SymbolContainer {
	return ls.Symbols[node]
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
	name, nameToken := p.srv.Linking().Namer.Name(node)
	if name != "" {
		var segment *core.TextSegment
		if nameToken != nil {
			segment = &nameToken.Segment
		}
		desc := core.NewAstNodeDescription(node, name, segment, node.Segment())
		container := node.Container()
		return desc, container
	}
	return nil, nil
}
