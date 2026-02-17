// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"iter"
	"slices"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/extiter"
)

func DefaultScopeOfType[T core.AstNode](node core.AstNode, allDocuments iter.Seq[*core.Document]) core.Scope {
	globalScope := GlobalScopeOfType[T](node, allDocuments)
	return LocalScopeOfType[T](node, globalScope)
}

func GlobalScopeOfType[T core.AstNode](node core.AstNode, allDocuments iter.Seq[*core.Document]) core.Scope {
	symbols := extiter.FlatMap(allDocuments, func(doc *core.Document) iter.Seq[*core.AstNodeDescription] {
		if len(doc.ExportedSymbols) == 0 {
			return core.EmptyAstNodeDescriptions
		}
		exported := slices.Values(doc.ExportedSymbols)
		return SymbolsOfType[T](exported)
	})
	return core.NewMapScopeFromSeq(symbols, nil)
}

func LocalScopeOfType[T core.AstNode](node core.AstNode, globalScope core.Scope) core.Scope {
	symbols := GetSymbolList(node)
	filtered := SymbolsOfType[T](symbols)
	var outer core.Scope
	if container := node.Container(); container != nil {
		// The container node (or one of its ancestors) defines the outer scope
		outer = LocalScopeOfType[T](container, globalScope)
	} else {
		// We're at the root node, use the given global scope
		outer = globalScope
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

func GetSymbolList(node core.AstNode) core.SymbolList {
	doc := node.Document()
	localSymbols := doc.LocalSymbols
	return localSymbols.Iter(node)
}

func SymbolsOfType[T core.AstNode](s core.SymbolList) core.SymbolList {
	return func(yield func(*core.AstNodeDescription) bool) {
		for desc := range s {
			if _, ok := desc.Node.(T); ok {
				if !yield(desc) {
					return
				}
			}
		}
	}
}
