// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/atn"
	"typefox.dev/fastbelt/internal/grammar"
)

func GenerateATN(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult) string {
	tokenTypeIds := tokenTypes.TokenTypeIds()
	tokenTypeNames := tokenTypes.TokenTypeVarNamesByTokenIndex()
	a, _ := atn.CreateATN(grammr, tokenTypeIds)
	source := atn.EmitGoSource(packageName, a, grammr, tokenTypeNames)
	return FormatIfPossible(source.String())
}
