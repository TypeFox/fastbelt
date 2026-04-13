// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"
	"iter"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/extiter"
)

// ImportedSymbolsProvider computes the symbols imported into a document from
// other documents, making them available for cross-document reference resolution.
type ImportedSymbolsProvider interface {
	// Provide creates a sequence of all symbols that are visible from other documents.
	// The result is stored in the document's ImportedSymbols field.
	// The caller must hold the document's write lock.
	Provide(ctx context.Context, document *core.Document, allDocuments iter.Seq[*core.Document])
}

// DefaultImportedSymbolsProvider is the default implementation of ImportedSymbolsProvider.
// It flat-maps the exported symbols of all documents into a single lazy sequence.
type DefaultImportedSymbolsProvider struct {
	srv LinkingSrvCont
}

func NewDefaultImportedSymbolsProvider(srv LinkingSrvCont) ImportedSymbolsProvider {
	return &DefaultImportedSymbolsProvider{
		srv: srv,
	}
}

func (p *DefaultImportedSymbolsProvider) Provide(ctx context.Context, doc *core.Document, allDocs iter.Seq[*core.Document]) {
	doc.ImportedSymbols = core.MergeSymbolContainers(extiter.Map(allDocs, func(d *core.Document) core.SymbolContainer {
		return d.ExportedSymbols
	}))
}
