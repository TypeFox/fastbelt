// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/atn"
	"typefox.dev/fastbelt/internal/grammar"
)

func GenerateATN(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult) string {
	a, _ := atn.CreateATN(grammr, tokenTypes.TokenTypeIds)
	source := atn.EmitGoSource(packageName, "BuildATN", "typefox.dev/fastbelt/parser", a, tokenTypes.TokenTypeVarNames)
	return FormatIfPossible(source.String())
}
