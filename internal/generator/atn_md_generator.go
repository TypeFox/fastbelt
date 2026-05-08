// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
)

func GenerateATNMarkdown(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult) string {
	atn, _ := CreateATN(grammr, tokenTypes.TokenTypeIds)
	source := EmitMarkdownSource(packageName, atn, tokenTypes.TokenTypeNames)
	return source.String()
}
