// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.
package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser/allstar"
)

func GenerateATN(grammr grammar.Grammar, packageName string) string {
	rules, err := FromParserRules(grammr.Rules(), getTokenTypes(grammr))
	if err != nil {
		panic(err)
	}
	atn := allstar.CreateATN(rules)
	rtn := allstar.BuildRuntimeATN(atn)
	source := allstar.EmitGoSource(packageName, "BuildATN", "typefox.dev/fastbelt/parser/allstar", rtn)
	return FormatIfPossible(source)
}

func getTokenTypes(grammr grammar.Grammar) map[string]TokenInfo {
	tokens := grammr.Terminals()
	keywords := GetAllKeywords(grammr)
	nodes := make(map[string]TokenInfo, len(tokens)+len(keywords))
	id := 1
	for _, keyword := range keywords {
		nodes[keyword.Text()] = TokenInfo{ID: id}
		id++
	}
	for _, token := range tokens {
		nodes[token.Name()] = TokenInfo{ID: id}
		id++
	}
	return nodes
}
