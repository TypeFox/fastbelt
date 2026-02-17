// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/extiter"
)

func LocalScopeOfType[T core.AstNode](node core.AstNode) core.Scope {
	symbols := GetSymbolList(node)
	filtered := SymbolsOfType[T](symbols)
	var outer core.Scope = nil
	if container := node.Container(); container != nil {
		outer = LocalScopeOfType[T](container)
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
	return extiter.Filter(s, func(desc *core.AstNodeDescription) bool {
		_, ok := desc.Node.(T)
		return ok
	})
}
