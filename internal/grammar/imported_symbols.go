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
	// mu serializes building the folder view, since documents are imported in parallel.
	mu sync.Mutex
}

func newImportedSymbolsProviderImpl(sc *service.Container) linking.SymbolImporter {
	return &importedSymbolsProviderImpl{sc: sc}
}

// folderViewKey is the document data key under which the folderView is stored.
type folderViewKey struct{}

// folderView is the folder-wide grammar view shared by all .fb files of a folder.
// It is built once per folder state and stored on every document of the folder.
// Note: A view is stale once a document or its root pointer changed.
type folderView struct {
	// docs are the documents of the folder, sorted by file name.
	docs []*core.Document
	// roots are the roots of docs at build time.
	roots []core.AstNode
	// grammars are the roots that are grammars, in the same order.
	grammars []Grammar
	// folder is the aggregate of all grammars, or the single grammar.
	folder Grammar
	// interfaces indexes the interfaces of all grammars by name.
	interfaces map[string]Interface

	// mu guards diagnostics.
	mu sync.Mutex
	// diagnostics caches the result of the folder-wide checks per document.
	diagnostics map[*core.Document][]*core.Diagnostic
}

func (s *importedSymbolsProviderImpl) ImportSymbols(ctx context.Context, doc *core.Document, allDocs iter.Seq[*core.Document]) core.SymbolContainer {
	// Only grammar definitions in the same package are visible
	sameFolderDocs := slices.Collect(extiter.Filter(allDocs, func(other *core.Document) bool {
		return sameFolder(doc.URI, other.URI)
	}))
	// Sort by file name, like cmd.LoadGrammarDir does. The file order determines the terminal order.
	slices.SortFunc(sameFolderDocs, func(a, b *core.Document) int {
		return strings.Compare(a.URI.Path(), b.URI.Path())
	})
	allExportedSymbols := extiter.Map(slices.Values(sameFolderDocs), func(d *core.Document) core.SymbolContainer {
		return d.ExportedSymbols
	})
	imported := core.MergeSymbolContainers(allExportedSymbols)
	doc.ImportedSymbols = imported
	s.storeFolderView(doc, sameFolderDocs)
	return imported
}

// storeFolderView stores the folder view on doc. It reuses the view of a
// sibling if it still matches docs and builds a new one otherwise.
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
// all given grammars. The listed nodes keep their original container and document.
// Elements are sorted by name, except terminals, which keep the order of the
// grammars and their declaration order within each file.
// Note: Grammars with overlapping tokens across files should declare token modes,
// which fix the precedence independent of file names.
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
	// The getters return the backing slices, so sorting in place works.
	sortByName(folder.Rules())
	sortByName(folder.Composites())
	sortByName(folder.InfixRules())
	sortByName(folder.TokenGroups())
	sortByName(folder.TokenModes())
	sortByName(folder.Interfaces())
	return folder
}

// sortByName stably sorts items by name in place.
func sortByName[T interface{ Name() string }](items []T) {
	slices.SortStableFunc(items, func(a, b T) int { return strings.Compare(a.Name(), b.Name()) })
}

// viewOf returns the folderView of g's document, or nil if there is none.
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

// sameFolder reports whether two URIs point into the same folder.
func sameFolder(first, second core.URI) bool {
	return first.Scheme() == second.Scheme() &&
		first.Authority() == second.Authority() &&
		path.Dir(first.Path()) == path.Dir(second.Path())
}
