// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
)

// ExportedSymbolsProvider computes the symbols to be exported from a document
// so they can be imported into other documents.
type ExportedSymbolsProvider interface {
	// Provide traverses the document's AST and computes the exported symbols.
	// The result is stored in the document's ExportedSymbols field.
	// The caller must hold the document's write lock.
	Provide(ctx context.Context, document *core.Document)
}

// DefaultExportedSymbolsProvider is the default implementation of ExportedSymbolsProvider.
// By default, it exports the root node and its direct children that have a name.
type DefaultExportedSymbolsProvider struct {
	srv LinkingSrvCont
}

func NewDefaultExportedSymbolsProvider(srv LinkingSrvCont) ExportedSymbolsProvider {
	return &DefaultExportedSymbolsProvider{
		srv: srv,
	}
}

func (s *DefaultExportedSymbolsProvider) Provide(ctx context.Context, document *core.Document) {
	root := document.Root
	describer := s.srv.Linking().ExportedSymbolDescriber
	exports := s.srv.Generated().SymbolContainers.New()

	// Describe the root node itself
	if desc := describer.Describe(root); desc != nil {
		exports.Put(desc)
	}
	// Describe direct children of the root (not nested deeper)
	for child := range core.ChildNodes(root) {
		if desc := describer.Describe(child); desc != nil {
			exports.Put(desc)
		}
	}

	document.ExportedSymbols = exports
}

// ExportedSymbolDescriber describes how symbols are exported from a document.
type ExportedSymbolDescriber interface {
	// Describe determines the name and other metadata of an exported symbol.
	// It returns the description, or nil if the node should not be exported.
	Describe(node core.AstNode) *core.SymbolDescription
}

// DefaultExportedSymbolDescriber is the default implementation of ExportedSymbolDescriber.
type DefaultExportedSymbolDescriber struct {
	srv LinkingSrvCont
}

func NewDefaultExportedSymbolDescriber(srv LinkingSrvCont) ExportedSymbolDescriber {
	return &DefaultExportedSymbolDescriber{
		srv: srv,
	}
}

func (p *DefaultExportedSymbolDescriber) Describe(node core.AstNode) *core.SymbolDescription {
	nameUnit := p.srv.Linking().Namer.Name(node)
	if nameUnit != nil {
		name := nameUnit.String()
		segment := nameUnit.Segment()
		fullSegment := node.Segment()
		return core.NewSymbolDescription(node, name, segment, fullSegment)
	}
	return nil
}
