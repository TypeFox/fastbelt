// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.
package generator

import (
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser/allstar"
)

func GenerateATNMarkdown(grammr grammar.Grammar, packageName string) string {
	tokenTypes := GetTokenTypes(grammr)
	tokenTypeNames := make(map[int]string, len(tokenTypes))
	for name, id := range tokenTypes {
		tokenTypeNames[id.ID] = strings.ReplaceAll(name, "\"", "'")
	}
	rules, err := FromParserRules(grammr.Rules(), tokenTypes)
	if err != nil {
		panic(err)
	}
	atn := allstar.CreateATN(rules)
	rtn := allstar.BuildRuntimeATN(atn)
	source := allstar.EmitMarkdownSource(packageName, rtn, tokenTypeNames)
	return source.String()
}
