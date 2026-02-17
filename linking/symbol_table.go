// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"
	"slices"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/extiter"
)

func SymbolsOfType[T core.AstNode](s core.SymbolList) core.SymbolList {
	return extiter.Filter(s, func(desc *core.AstNodeDescription) bool {
		_, ok := desc.Node.(T)
		return ok
	})
}

func LocalScopeOfType[T core.AstNode](node core.AstNode, fn func(core.AstNode) core.SymbolList) core.Scope {
	symbols := fn(node)
	filtered := SymbolsOfType[T](symbols)
	var outer core.Scope = nil
	if container := node.Container(); container != nil {
		outer = LocalScopeOfType[T](container, fn)
	}
	if extiter.IsEmpty(filtered) {
		// Shortcut to generate fewer scopes
		if outer != nil {
			return outer
		} else {
			return core.EmptyScope
		}
	}
	return core.NewMapScopeFromSeq(filtered, outer)
}

// LocalSymbolTableProvider computes and provides local symbol information for documents.
type LocalSymbolTableProvider interface {
	// Compute traverses the document's AST and builds the local symbol table.
	// The caller must hold the document's write lock.
	Compute(ctx context.Context, document *core.Document)
	// LocalSymbols returns the symbols visible at the given node.
	LocalSymbols(node core.AstNode) core.SymbolList
}

type DefaultLocalSymbolTableProvider struct {
	srv LinkingSrvCont
}

func NewDefaultLocalSymbolTableProvider(srv LinkingSrvCont) LocalSymbolTableProvider {
	return &DefaultLocalSymbolTableProvider{
		srv: srv,
	}
}

func (s *DefaultLocalSymbolTableProvider) Compute(ctx context.Context, document *core.Document) {
	root := document.Root
	symbols := collections.NewMultiMap[core.AstNode, *core.AstNodeDescription]()

	core.TraverseContent(root, func(node core.AstNode) {
		item := s.srv.Linking().LocalSymbolTableItemProvider.Item(node)
		if item != nil {
			symbols.Put(item.Container, item.Description)
		}
	})
	document.LocalSymbols = NewDefaultLocalSymbolsFromMap(symbols)
	document.State = document.State.With(core.DocStateComputedSymbolTable)
}

func (s *DefaultLocalSymbolTableProvider) LocalSymbols(node core.AstNode) core.SymbolList {
	doc := node.Document()
	doc.RLock()
	defer doc.RUnlock()
	localSymbols := doc.LocalSymbols
	return localSymbols.Iter(node)
}

type DefaultLocalSymbols struct {
	Symbols collections.MultiMap[core.AstNode, *core.AstNodeDescription]
}

func NewDefaultLocalSymbols() *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: collections.NewMultiMap[core.AstNode, *core.AstNodeDescription](),
	}
}

func NewDefaultLocalSymbolsFromMap(m collections.MultiMap[core.AstNode, *core.AstNodeDescription]) *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: m,
	}
}

func (ls *DefaultLocalSymbols) Has(node core.AstNode) bool {
	return ls.Symbols.Has(node)
}

func (ls *DefaultLocalSymbols) Iter(node core.AstNode) core.SymbolList {
	symbols := ls.Symbols.Get(node)
	return slices.Values(symbols)
}

type LocalSymbolTableItem struct {
	Container   core.AstNode
	Description *core.AstNodeDescription
}

type LocalSymbolTableItemProvider interface {
	Item(node core.AstNode) *LocalSymbolTableItem
}

type DefaultLocalSymbolTableItemProvider struct {
	srv LinkingSrvCont
}

func NewDefaultLocalSymbolTableItemProvider(srv LinkingSrvCont) LocalSymbolTableItemProvider {
	return &DefaultLocalSymbolTableItemProvider{
		srv: srv,
	}
}

func (p *DefaultLocalSymbolTableItemProvider) Item(node core.AstNode) *LocalSymbolTableItem {
	container := node.Container()
	name, nameToken := p.srv.Linking().Namer.Name(node)
	if name != "" {
		var segment *core.TextSegment
		if nameToken != nil {
			segment = &nameToken.Segment
		}
		desc := core.NewAstNodeDescription(node, name, segment, node.Segment())
		return &LocalSymbolTableItem{
			Container:   container,
			Description: desc,
		}
	}
	return nil
}
