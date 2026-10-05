// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"iter"
	"path"
	"slices"
	"strings"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/extiter"
	"typefox.dev/fastbelt/util/service"
)

type importedSymbolsProviderImpl struct {
	sc *service.Container
	// Documents are imported in parallel, so the check-and-build of a folder
	// view is serialized to hand every sibling the same view.
	mu sync.Mutex
}

func newImportedSymbolsProviderImpl(sc *service.Container) linking.SymbolImporter {
	return &importedSymbolsProviderImpl{sc: sc}
}

// folderViewKey is the [core.Document.Data] key under which [ImportSymbols]
// stores the document's folderView.
type folderViewKey struct{}

// folderView is the "one grammar per folder" view shared by all .fb files of a
// folder: their grammar roots, an aggregate grammar listing every declaration,
// and an interface index. It is built once per folder state and reused by every
// document of the folder, so validation and return-type lookups do not scan the
// document manager or rebuild the aggregate per document.
//
// A view is valid for a given set of documents and roots. A changed file gets a
// new [core.Document], a reparsed file a new root, and a change in one file
// re-imports its siblings (see changeImpactImpl), so comparing the document and
// root pointers on import is enough to detect a stale view.
type folderView struct {
	docs       []*core.Document // sorted by URI
	roots      []core.AstNode   // docs[i].Root at build time
	grammars   []Grammar        // roots that are grammars, same order
	folder     Grammar          // aggregate of grammars, or the single grammar
	interfaces map[string]Interface
}

func (s *importedSymbolsProviderImpl) ImportSymbols(ctx context.Context, doc *core.Document, allDocs iter.Seq[*core.Document]) core.SymbolContainer {
	// Only grammar definitions in the same package are visible
	sameFolderDocs := slices.Collect(extiter.Filter(allDocs, func(other *core.Document) bool {
		return sameFolder(doc.URI, other.URI)
	}))
	slices.SortFunc(sameFolderDocs, func(a, b *core.Document) int {
		return strings.Compare(a.URI.String(), b.URI.String())
	})
	allExportedSymbols := extiter.Map(slices.Values(sameFolderDocs), func(d *core.Document) core.SymbolContainer {
		return d.ExportedSymbols
	})
	imported := core.MergeSymbolContainers(allExportedSymbols)
	doc.ImportedSymbols = imported
	s.storeFolderView(doc, sameFolderDocs)
	return imported
}

// storeFolderView stores the folder view on doc: the view a sibling already
// holds when it still matches docs, and a freshly built one otherwise. Views
// are referenced only by the documents they describe, so they are released
// together with them when a folder's files are deleted.
func (s *importedSymbolsProviderImpl) storeFolderView(doc *core.Document, docs []*core.Document) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, sibling := range docs {
		if value, ok := sibling.Data.Load(folderViewKey{}); ok && value.(*folderView).matches(docs) {
			doc.Data.Store(folderViewKey{}, value)
			return
		}
	}
	doc.Data.Store(folderViewKey{}, newFolderView(docs))
}

func (v *folderView) matches(docs []*core.Document) bool {
	if len(docs) != len(v.docs) {
		return false
	}
	for i, doc := range docs {
		if doc != v.docs[i] || doc.Root != v.roots[i] {
			return false
		}
	}
	return true
}

func newFolderView(docs []*core.Document) *folderView {
	view := &folderView{
		docs:       docs,
		roots:      make([]core.AstNode, len(docs)),
		interfaces: map[string]Interface{},
	}
	for i, doc := range docs {
		view.roots[i] = doc.Root
		if g, ok := doc.Root.(Grammar); ok {
			view.grammars = append(view.grammars, g)
		}
	}
	for _, g := range view.grammars {
		for _, iface := range g.Interfaces() {
			if _, seen := view.interfaces[iface.Name()]; !seen {
				view.interfaces[iface.Name()] = iface
			}
		}
	}
	switch len(view.grammars) {
	case 0:
	case 1:
		view.folder = view.grammars[0]
	default:
		view.folder = AggregateGrammar(view.grammars)
	}
	return view
}

// AggregateGrammar returns a detached grammar node listing the declarations of
// every given grammar. The generated setters only append to slices, so the
// listed nodes keep their real container and document; the aggregate is a
// read-only view for checks that must consider the whole folder. The code
// generator reparents it to merge the files of a folder.
//
// Elements are sorted by name so the generated code does not depend on the
// order of declarations across files. Terminals keep file order: it is
// semantically relevant and critical for correct lexing. The resulting
// terminal order is therefore the order of the grammars (sorted by URI by the
// callers, see ImportSymbols and cmd.LoadGrammarDir) and declaration order
// within each file. Nothing checks that token rules of different files do not
// overlap, so grammars with overlapping tokens across files should declare
// token modes, whose member order fixes the precedence independent of file
// names.
func AggregateGrammar(grammars []Grammar) Grammar {
	folder := NewGrammar()
	folder.SetName(grammars[0].NameToken())
	for _, g := range grammars {
		for _, item := range g.Rules() {
			folder.SetRulesItem(item)
		}
		for _, item := range g.Composites() {
			folder.SetCompositesItem(item)
		}
		for _, item := range g.InfixRules() {
			folder.SetInfixRulesItem(item)
		}
		for _, item := range g.Terminals() {
			folder.SetTerminalsItem(item)
		}
		for _, item := range g.TokenGroups() {
			folder.SetTokenGroupsItem(item)
		}
		for _, item := range g.TokenModes() {
			folder.SetTokenModesItem(item)
		}
		for _, item := range g.Interfaces() {
			folder.SetInterfacesItem(item)
		}
	}
	// The getters return the backing slices, so sorting in place reorders the
	// aggregate.
	sortByName(folder.Rules())
	sortByName(folder.Composites())
	sortByName(folder.InfixRules())
	sortByName(folder.TokenGroups())
	sortByName(folder.TokenModes())
	sortByName(folder.Interfaces())
	return folder
}

// sortByName sorts items by name in place. The sort is stable, so elements
// with the same name (rejected by validation anyway) keep file order.
func sortByName[T interface{ Name() string }](items []T) {
	slices.SortStableFunc(items, func(a, b T) int { return strings.Compare(a.Name(), b.Name()) })
}

// viewOf returns the folderView of the document g belongs to, or nil when the
// document was never imported (e.g. the synthetic merged grammar of the code
// generator) or g is not one of the view's grammars.
func viewOf(g Grammar) *folderView {
	doc := g.Document()
	if doc == nil {
		return nil
	}
	if value, ok := doc.Data.Load(folderViewKey{}); ok {
		if view := value.(*folderView); slices.Contains(view.grammars, g) {
			return view
		}
	}
	return nil
}

func sameFolder(first, second core.URI) bool {
	return path.Dir(first.Path()) == path.Dir(second.Path())
}
