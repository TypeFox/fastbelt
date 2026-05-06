// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"reflect"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/extiter"
)

// DefaultScopeOfType is the default way of creating a scope for a reference in the given AST node.
// It combines the global scope (imported symbols) and the local scope (symbols visible from a container).
func DefaultScopeOfType[T core.AstNode](node core.AstNode) core.Scope {
	globalScope := GlobalScopeOfType[T](node)
	return LocalScopeOfType[T](node, globalScope)
}

// GlobalScopeOfType is the default way of creating a global scope for a reference in the given AST node.
// It returns the scope of imported symbols for the given type.
func GlobalScopeOfType[T core.AstNode](node core.AstNode) core.Scope {
	imported := node.Document().ImportedSymbols
	if imported == nil {
		return core.EmptyScope
	}
	symbols := imported.ForType(reflect.TypeFor[T]())
	return core.NewSeqScope(symbols, nil)
}

// LocalScopeOfType is the default way of creating a local scope for a reference in the given AST node.
// It returns the scope of symbols visible from the given AST node or any of its containers.
func LocalScopeOfType[T core.AstNode](node core.AstNode, globalScope core.Scope) core.Scope {
	symbols := GetLocalSymbols[T](node)
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

// GetLocalSymbols retrieves the local symbols for the given AST node and type.
func GetLocalSymbols[T core.AstNode](node core.AstNode) core.SymbolSeq {
	doc := node.Document()
	localSymbols := doc.LocalSymbols
	container := localSymbols.For(node)
	symbols := container.ForType(reflect.TypeFor[T]())
	return symbols
}
