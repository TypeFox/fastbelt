// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// DocumentSymbolProvider is a service for handling LSP document symbol requests.
type DocumentSymbolProvider interface {
	HandleDocumentSymbolRequest(ctx context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error)
}

// DocumentSymbolFilter defines document-symbol-specific customization points.
// For name and kind information, use NameProvider and NodeKindProvider respectively.
type DocumentSymbolFilter interface {
	// ShouldInclude determines whether the specified AstNode should appear as a symbol.
	// By default, nodes with names are included.
	ShouldInclude(node core.AstNode) bool
}

// DefaultDocumentSymbolFilter provides standard symbol filtering behavior.
type DefaultDocumentSymbolFilter struct {
	nameProvider NameProvider
}

func (f *DefaultDocumentSymbolFilter) ShouldInclude(node core.AstNode) bool {
	return f.nameProvider.GetName(node) != ""
}

// DefaultDocumentSymbolProvider implements the DocumentSymbolProvider interface.
type DefaultDocumentSymbolProvider struct {
	sc           *service.Container
	nameProvider NameProvider
	kindProvider NodeKindProvider
	filter       DocumentSymbolFilter
}

// NewDefaultDocumentSymbolProvider creates a provider using services from the container.
// Services are loaded lazily on first use, allowing for service overrides after construction.
func NewDefaultDocumentSymbolProvider(sc *service.Container) DocumentSymbolProvider {
	return &DefaultDocumentSymbolProvider{
		sc: sc,
	}
}

// NewDocumentSymbolProviderWithFilter creates a provider with a custom filter.
func NewDocumentSymbolProviderWithFilter(sc *service.Container, filter DocumentSymbolFilter) DocumentSymbolProvider {
	return &DefaultDocumentSymbolProvider{
		sc:     sc,
		filter: filter,
	}
}

// ReloadServices forces reloading of NameProvider and NodeKindProvider from the service container.
func (p *DefaultDocumentSymbolProvider) ReloadServices() {
	p.nameProvider = service.MustGet[NameProvider](p.sc)
	p.kindProvider = service.MustGet[NodeKindProvider](p.sc)
	if p.filter == nil {
		p.filter = &DefaultDocumentSymbolFilter{nameProvider: p.nameProvider}
	}
}

func (p *DefaultDocumentSymbolProvider) ensureInitialized() {
	if p.nameProvider == nil || p.kindProvider == nil {
		p.ReloadServices()
	}
}

func (p *DefaultDocumentSymbolProvider) HandleDocumentSymbolRequest(ctx context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error) {
	p.ensureInitialized()

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
	if p.filter.ShouldInclude(node) {
		name := p.nameProvider.GetName(node)
		if name != "" && node.Segment() != nil {
			symbol := p.createSymbol(node, name)
			return []lsp.DocumentSymbol{symbol}
		}
	}

	return p.getChildSymbols(node)
}

func (p *DefaultDocumentSymbolProvider) createSymbol(node core.AstNode, name string) lsp.DocumentSymbol {
	segment := node.Segment()

	nameSegment := p.nameProvider.GetNameSegment(node)
	selectionRange := segment.Range.LspRange()
	if nameSegment != nil {
		selectionRange = nameSegment.Range.LspRange()
	}

	symbol := lsp.DocumentSymbol{
		Name:           name,
		Kind:           p.kindProvider.GetSymbolKind(node),
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
