// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"
	"iter"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/extiter"
	"typefox.dev/fastbelt/util/service"
)

// ImportedSymbolsProvider is a service that computes the symbols imported into a document from
// other documents, making them available for cross-document reference resolution.
type ImportedSymbolsProvider interface {
	// ImportedSymbols creates a sequence of all symbols that are visible from other documents.
	// The result is stored in the document's ImportedSymbols field.
	// The caller must hold the document's write lock.
	ImportedSymbols(ctx context.Context, document *core.Document, allDocuments iter.Seq[*core.Document]) core.SymbolContainer
}

// DefaultImportedSymbolsProvider is the default implementation of [ImportedSymbolsProvider].
// It flat-maps the exported symbols of all documents into a single lazy sequence.
type DefaultImportedSymbolsProvider struct {
	sc *service.Container
}

func NewDefaultImportedSymbolsProvider(sc *service.Container) ImportedSymbolsProvider {
	return &DefaultImportedSymbolsProvider{sc: sc}
}

func (s *DefaultImportedSymbolsProvider) ImportedSymbols(ctx context.Context, doc *core.Document, allDocs iter.Seq[*core.Document]) core.SymbolContainer {
	allExportedSymbols := extiter.Map(allDocs, func(d *core.Document) core.SymbolContainer {
		return d.ExportedSymbols
	})
	imported := core.MergeSymbolContainers(allExportedSymbols)
	doc.ImportedSymbols = imported
	return imported
}
