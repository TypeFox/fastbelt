// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"iter"
	"path"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/extiter"
	"typefox.dev/fastbelt/util/service"
)

type importedSymbolsProviderImpl struct {
	sc *service.Container
}

func newImportedSymbolsProviderImpl(sc *service.Container) linking.SymbolImporter {
	return &importedSymbolsProviderImpl{sc: sc}
}

func (s *importedSymbolsProviderImpl) ImportSymbols(ctx context.Context, doc *core.Document, allDocs iter.Seq[*core.Document]) core.SymbolContainer {
	// Only grammar definitions in the same package are visible
	sameFolderDocs := extiter.Filter(allDocs, func(other *core.Document) bool {
		return sameFolder(doc.URI, other.URI)
	})
	allExportedSymbols := extiter.Map(sameFolderDocs, func(d *core.Document) core.SymbolContainer {
		return d.ExportedSymbols
	})
	imported := core.MergeSymbolContainers(allExportedSymbols)
	doc.ImportedSymbols = imported
	return imported
}

func sameFolder(first, second core.URI) bool {
	return path.Dir(first.Path()) == path.Dir(second.Path())
}
