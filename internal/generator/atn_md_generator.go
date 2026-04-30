// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.
package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
)

func GenerateATNMarkdown(grammr grammar.Grammar, packageName string) string {
	atn, _, tokenTypes := CreateATN(grammr)
	tokenTypeNames := make(map[int]string, len(tokenTypes))
	for name, info := range tokenTypes {
		tokenTypeNames[info.ID] = name
	}
	rtn := BuildRuntimeATN(atn)
	source := EmitMarkdownSource(packageName, rtn, tokenTypeNames)
	return source.String()
}
