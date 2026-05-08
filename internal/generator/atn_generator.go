// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
)

func GenerateATN(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult) string {
	atn, _, _ := CreateATN(grammr)
	source := EmitGoSource(packageName, "BuildATN", "typefox.dev/fastbelt/parser", atn, tokenTypes)
	return FormatIfPossible(source.String())
}
