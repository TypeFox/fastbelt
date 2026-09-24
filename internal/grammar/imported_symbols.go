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

// siblingDocsKey is the [core.Document.Data] key under which [ImportSymbols]
// stores the live sequence of same-folder documents (the document itself
// included). The grammar language treats every .fb file of a folder as one
// grammar, and validation needs the sibling ASTs, not just their exported
// symbols: a file may contribute nothing linkable (e.g. only an unnamed default
// token mode) and still take part in folder-wide checks.
type siblingDocsKey struct{}

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
	// The filter is lazy over the document manager, so the stored value never
	// goes stale even though the builder does not clear Data on reset.
	doc.Data.Store(siblingDocsKey{}, sameFolderDocs)
	return imported
}

// siblingDocuments returns the documents of the folder doc belongs to, as
// recorded by [ImportSymbols], or nil when the document was never imported
// (e.g. the synthetic merged grammar of the code generator).
func siblingDocuments(doc *core.Document) iter.Seq[*core.Document] {
	if doc == nil {
		return nil
	}
	if seq, ok := doc.Data.Load(siblingDocsKey{}); ok {
		return seq.(iter.Seq[*core.Document])
	}
	return nil
}

func sameFolder(first, second core.URI) bool {
	return path.Dir(first.Path()) == path.Dir(second.Path())
}
