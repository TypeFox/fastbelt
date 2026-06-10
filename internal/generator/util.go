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

func GeneratedTokenName(t core.AstNode) string {
	switch t := t.(type) {
	case grammar.AbstractTokenRule:
		return "Token_" + t.Name()
	case grammar.Keyword:
		return "Keyword_" + grammar.KeywordName(t)
	default:
		panic("unexpected type")
	}
}

func GeneratedTokenIdxName(t core.AstNode) string {
	return GeneratedTokenName(t) + "_Idx"
}

func KeywordValue(k grammar.Keyword) string {
	return k.Value()[1 : len(k.Value())-1]
}

// topoSort returns uniqueTargets sorted so that child elements appear before
// their parent element.
func topoSort(uniqueTargets []string, parents map[string][]string, ascending bool) []string {
	targetSet := map[string]bool{}
	for _, t := range uniqueTargets {
		targetSet[t] = true
	}
	visited := map[string]bool{}
	result := []string{}
	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true
		for _, parent := range parents[name] {
			visit(parent)
		}
		if ascending {
			// Prepend so children end up before parents in the final slice.
			result = append([]string{name}, result...)
		} else {
			// Append so parents end up after children in the final slice.
			result = append(result, name)
		}
	}
	for _, t := range uniqueTargets {
		visit(t)
	}
	// Keep only types that were in uniqueTargets (visit may have walked ancestors
	// not in the set; filter them out while preserving order).
	filtered := result[:0]
	for _, t := range result {
		if targetSet[t] {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
