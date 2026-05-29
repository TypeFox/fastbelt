// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server/fuzzy"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// WorkspaceSymbolProvider is a service for handling LSP workspace symbol requests.
type WorkspaceSymbolProvider interface {
	HandleWorkspaceSymbolRequest(ctx context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error)
}

// WorkspaceSymbolFilter defines customization points for workspace symbol collection.
type WorkspaceSymbolFilter interface {
	// ShouldInclude determines whether the specified AstNode should appear as a workspace symbol.
	// By default, nodes with names are included.
	ShouldInclude(node core.AstNode) bool

	// MaxSymbolCount returns the maximum number of symbols to return.
	// Returns 0 for unlimited. Default is 1000 to avoid overwhelming clients.
	MaxSymbolCount() int
}

// DefaultWorkspaceSymbolFilter provides standard workspace symbol filtering.
type DefaultWorkspaceSymbolFilter struct{}

func (f *DefaultWorkspaceSymbolFilter) ShouldInclude(node core.AstNode) bool {
	return true
}

func (f *DefaultWorkspaceSymbolFilter) MaxSymbolCount() int {
	return 1000
}

// DefaultWorkspaceSymbolProvider implements WorkspaceSymbolProvider.
type DefaultWorkspaceSymbolProvider struct {
	sc     *service.Container
	filter WorkspaceSymbolFilter
}

// NewDefaultWorkspaceSymbolProvider creates a provider using services from the container.
func NewDefaultWorkspaceSymbolProvider(sc *service.Container) WorkspaceSymbolProvider {
	return &DefaultWorkspaceSymbolProvider{
		sc:     sc,
		filter: &DefaultWorkspaceSymbolFilter{},
	}
}

// NewWorkspaceSymbolProviderWithFilter creates a provider with a custom filter.
func NewWorkspaceSymbolProviderWithFilter(sc *service.Container, filter WorkspaceSymbolFilter) WorkspaceSymbolProvider {
	return &DefaultWorkspaceSymbolProvider{
		sc:     sc,
		filter: filter,
	}
}

func (p *DefaultWorkspaceSymbolProvider) HandleWorkspaceSymbolRequest(ctx context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	query := strings.ToLower(params.Query)
	maxCount := p.filter.MaxSymbolCount()

	var symbols []lsp.SymbolInformation

	// Iterate all documents in workspace
	for doc := range documentManager.All() {
		if doc.Root == nil {
			continue
		}

		// Collect symbols from this document, passing current count and max for early abort
		docSymbols := p.collectSymbolsFromDocument(doc, query, len(symbols), maxCount)
		symbols = append(symbols, docSymbols...)

		if maxCount > 0 && len(symbols) >= maxCount {
			return symbols[:maxCount], nil
		}
	}

	return symbols, nil
}

func (p *DefaultWorkspaceSymbolProvider) collectSymbolsFromDocument(doc *core.Document, query string, currentCount, maxCount int) []lsp.SymbolInformation {
	var symbols []lsp.SymbolInformation

	// Iterate all nodes in document
	for node := range core.AllNodes(doc.Root) {
		if maxCount > 0 && currentCount >= maxCount {
			break
		}

		if !p.filter.ShouldInclude(node) {
			continue
		}

		nameUnit := linking.Name(node)
		if nameUnit == nil {
			continue
		}

		name := nameUnit.String()

		if !fuzzy.Match(query, name) {
			continue
		}

		nameSegment := nameUnit.Segment()
		if nameSegment == nil {
			continue
		}

		// Create SymbolInformation (matching Langium's structure)
		symbol := lsp.SymbolInformation{
			Name: name,
			Kind: SymbolKind(node),
			Location: lsp.Location{
				URI:   lsp.DocumentURI(doc.URI.DocumentURI()),
				Range: nameSegment.Range.LspRange(),
			},
		}

		symbols = append(symbols, symbol)
		currentCount++
	}

	return symbols
}
