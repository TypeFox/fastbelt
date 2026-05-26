// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// DocumentSymbolProvider is a service for handling LSP document symbol requests.
type DocumentSymbolProvider interface {
	HandleDocumentSymbolRequest(ctx context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error)
}

// DocumentSymbolFilter defines document-symbol-specific customization points.
type DocumentSymbolFilter interface {
	// ShouldInclude determines whether the specified AstNode should appear as a symbol.
	// By default, nodes with names are included.
	ShouldInclude(node core.AstNode) bool
}

// DefaultDocumentSymbolFilter provides standard symbol filtering behavior.
type DefaultDocumentSymbolFilter struct{}

func (f *DefaultDocumentSymbolFilter) ShouldInclude(node core.AstNode) bool {
	return true
}

// DefaultDocumentSymbolProvider implements the DocumentSymbolProvider interface.
type DefaultDocumentSymbolProvider struct {
	sc     *service.Container
	filter DocumentSymbolFilter
}

// NewDefaultDocumentSymbolProvider creates a provider using services from the container.
func NewDefaultDocumentSymbolProvider(sc *service.Container) DocumentSymbolProvider {
	return &DefaultDocumentSymbolProvider{
		sc:     sc,
		filter: &DefaultDocumentSymbolFilter{},
	}
}

// NewDocumentSymbolProviderWithFilter creates a provider with a custom filter.
func NewDocumentSymbolProviderWithFilter(sc *service.Container, filter DocumentSymbolFilter) DocumentSymbolProvider {
	return &DefaultDocumentSymbolProvider{
		sc:     sc,
		filter: filter,
	}
}

func (p *DefaultDocumentSymbolProvider) HandleDocumentSymbolRequest(ctx context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil // Document not found
	}
	return p.collectSymbols(doc), nil
}

func (p *DefaultDocumentSymbolProvider) collectSymbols(document *core.Document) []lsp.DocumentSymbol {
	if document.Root == nil {
		return nil
	}

	return p.getSymbolsForNode(document.Root)
}

func (p *DefaultDocumentSymbolProvider) getSymbolsForNode(node core.AstNode) []lsp.DocumentSymbol {
	// Check if this node should be included as a symbol
	if p.filter.ShouldInclude(node) {
		nameUnit := linking.Name(node)
		if nameUnit != nil && node.Segment() != nil {
			symbol := p.createSymbol(node, nameUnit)
			return []lsp.DocumentSymbol{symbol}
		}
	}

	// If node is not a symbol itself, collect symbols from children
	return p.getChildSymbols(node)
}

func (p *DefaultDocumentSymbolProvider) createSymbol(node core.AstNode, nameUnit core.StringUnit) lsp.DocumentSymbol {
	segment := node.Segment()

	name := nameUnit.String()
	selectionRange := segment.Range.LspRange()
	if nameSegment := nameUnit.Segment(); nameSegment != nil {
		selectionRange = nameSegment.Range.LspRange()
	}

	symbol := lsp.DocumentSymbol{
		Name:           name,
		Kind:           SymbolKind(node),
		Range:          segment.Range.LspRange(),
		SelectionRange: selectionRange,
	}

	children := p.getChildSymbols(node)
	if len(children) > 0 {
		symbol.Children = children
	}

	return symbol
}

func (p *DefaultDocumentSymbolProvider) getChildSymbols(node core.AstNode) []lsp.DocumentSymbol {
	var symbols []lsp.DocumentSymbol

	for child := range core.ChildNodes(node) {
		childSymbols := p.getSymbolsForNode(child)
		symbols = append(symbols, childSymbols...)
	}

	return symbols
}
