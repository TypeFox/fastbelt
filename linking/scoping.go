// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"reflect"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/extiter"
)

func DefaultScopeOfType[T core.AstNode](node core.AstNode) core.Scope {
	globalScope := GlobalScopeOfType[T](node)
	return LocalScopeOfType[T](node, globalScope)
}

func GlobalScopeOfType[T core.AstNode](node core.AstNode) core.Scope {
	imported := node.Document().ImportedSymbols
	if imported == nil {
		return core.EmptyScope
	}
	symbols := imported.Type(reflect.TypeFor[T]())
	return core.NewSeqScope(symbols, nil)
}

func LocalScopeOfType[T core.AstNode](node core.AstNode, globalScope core.Scope) core.Scope {
	t := reflect.TypeFor[T]()
	symbols := GetLocalSymbols(node, t)
	var outer core.Scope
	if container := node.Container(); container != nil {
		// The container node (or one of its ancestors) defines the outer scope
		outer = LocalScopeOfType[T](container, globalScope)
	} else {
		// We're at the root node, use the given global scope
		outer = globalScope
	}
	if extiter.IsEmpty(symbols) {
		// Shortcut to generate fewer scopes
		if outer != nil {
			return outer
		} else {
			return core.EmptyScope
		}
	}
	return core.NewSeqScope(symbols, outer)
}

func GetLocalSymbols(node core.AstNode, t reflect.Type) core.SymbolSeq {
	doc := node.Document()
	localSymbols := doc.LocalSymbols
	return localSymbols.For(node).Type(t)
}
