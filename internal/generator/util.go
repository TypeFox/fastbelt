// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"go/format"
	"runtime"
	"sort"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
)

const CardinalityOne = ""
const CardinalityOptional = "?"
const CardinalityZeroOrMore = "*"
const CardinalityOneOrMore = "+"

func FormatIfPossible(text string) string {
	formatted, err := format.Source([]byte(text))
	if err != nil {
		return text
	}
	return string(formatted)
}

func eol() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

var EOL = eol()

const Indent = "    "

func EOLIndent(count int) string {
	return EOL + strings.Repeat(Indent, count)
}

func WriteSB(sb *strings.Builder, texts ...string) {
	for _, text := range texts {
		sb.WriteString(text)
	}
}

func GetAllKeywords(grammr grammar.Grammar) []grammar.Keyword {
	keywords := map[string]grammar.Keyword{}
	for node := range core.AllChildren(grammr) {
		if keyword, ok := node.(grammar.Keyword); ok {
			keywords[keyword.Value()] = keyword
		}
	}
	return keysFromMap(keywords)
}

func keysFromMap(m map[string]grammar.Keyword) []grammar.Keyword {
	keys := []grammar.Keyword{}
	for _, v := range m {
		keys = append(keys, v)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Value() < keys[j].Value()
	})
	return keys
}

func GeneratedTokenName(t grammar.Token) string {
	return "Token_" + t.Name()
}

func GeneratedTokenIdxName(t grammar.Token) string {
	return GeneratedTokenName(t) + "_Idx"
}

func GeneratedKeywordName(k grammar.Keyword) string {
	return "Keyword_" + grammar.KeywordName(k)
}

func GeneratedKeywordIdxName(k grammar.Keyword) string {
	return GeneratedKeywordName(k) + "_Idx"
}
