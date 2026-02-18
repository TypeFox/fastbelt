// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
)

// ExportedSymbolsProvider computes the symbols to be exported from a document,
// potentially making them visible from other documents.
type ExportedSymbolsProvider interface {
	// Compute traverses the document's AST and computes the exported symbols.
	// The result is stored in the document's ExportedSymbols field.
	// The caller must hold the document's write lock.
	Compute(ctx context.Context, document *core.Document)
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

func (s *DefaultExportedSymbolsProvider) Compute(ctx context.Context, document *core.Document) {
	root := document.Root
	describer := s.srv.Linking().ExportedSymbolDescriber
	var exports []*core.AstNodeDescription

	// Describe the root node itself
	if desc := describer.Describe(root); desc != nil {
		exports = append(exports, desc)
	}
	// Describe direct children of the root (not nested deeper)
	root.ForEachNode(func(node core.AstNode) {
		if desc := describer.Describe(node); desc != nil {
			exports = append(exports, desc)
		}
	})

	document.ExportedSymbols = exports
}

// ExportedSymbolDescriber describes how symbols are exported from a document.
type ExportedSymbolDescriber interface {
	// Describe determines the name and other metadata of an exported symbol.
	// It returns the description, or nil if the node should not be exported.
	Describe(node core.AstNode) *core.AstNodeDescription
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

func (p *DefaultExportedSymbolDescriber) Describe(node core.AstNode) *core.AstNodeDescription {
	name, nameToken := p.srv.Linking().Namer.Name(node)
	if name != "" {
		var segment *core.TextSegment
		if nameToken != nil {
			segment = &nameToken.Segment
		}
		return core.NewAstNodeDescription(node, name, segment, node.Segment())
	}
	return nil
}
