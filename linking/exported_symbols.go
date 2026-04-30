// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// ExportedSymbolsProvider is a service that computes the symbols to be exported from a document,
// so they can be imported into other documents.
type ExportedSymbolsProvider interface {
	// ExportedSymbols traverses the document's AST and computes the exported symbols.
	// The result is stored in the document's ExportedSymbols field.
	ExportedSymbols(ctx context.Context, document *core.Document) core.SymbolContainer
}

// DefaultExportedSymbolsProvider is the default implementation of [ExportedSymbolsProvider].
// By default, it exports the root node and its direct children that have a name.
type DefaultExportedSymbolsProvider struct {
	sc *service.Container
}

func NewDefaultExportedSymbolsProvider(sc *service.Container) ExportedSymbolsProvider {
	return &DefaultExportedSymbolsProvider{sc: sc}
}

func (s *DefaultExportedSymbolsProvider) ExportedSymbols(ctx context.Context, document *core.Document) core.SymbolContainer {
	root := document.Root
	exports := service.MustGet[core.SymbolContainers](s.sc).New()

	// Describe the root node itself
	if desc := DescribeExport(root); desc != nil {
		exports.Put(desc)
	}
	// Describe direct children of the root (not nested deeper)
	for child := range core.ChildNodes(root) {
		if desc := DescribeExport(child); desc != nil {
			exports.Put(desc)
		}
	}

	document.ExportedSymbols = exports
	return exports
}

// ExportedSymbolDescriber can be implemented by AST node Impl structs to provide custom exported symbol description logic.
type ExportedSymbolDescriber interface {
	// DescribeExport determines the name and other metadata of the receiver node as an exported symbol.
	// It returns the description, or nil if the node should not be exported.
	DescribeExport() *core.SymbolDescription
}

// DescribeExport checks whether the given node is a symbol (i.e. has a name) and returns a description for it.
//
// Language-specific implementations can be provided by implementing the [ExportedSymbolDescriber] interface.
func DescribeExport(node core.AstNode) *core.SymbolDescription {
	// Use the language-specific implementation associated with the node type if available
	if custom, ok := node.(ExportedSymbolDescriber); ok {
		return custom.DescribeExport()
	}

	nameUnit := Name(node)
	if nameUnit != nil {
		name := nameUnit.String()
		segment := nameUnit.Segment()
		fullSegment := node.Segment()
		return core.NewSymbolDescription(node, name, segment, fullSegment)
	}
	return nil
}
