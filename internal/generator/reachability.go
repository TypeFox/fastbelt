// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"sort"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
)

// KeywordsReachableFrom returns all grammar keyword literals that can be
// reached from entryRule by following rule calls transitively. The result is
// sorted by keyword value and is stable — it can be passed directly to a
// per-language lexer constructor so that only keywords belonging to that
// language's grammar are recognised, preventing keyword-pollution between
// languages that share a merged grammar.
func KeywordsReachableFrom(entryRule grammar.ParserRule, grammr grammar.Grammar) []grammar.Keyword {
	visited := map[string]bool{}
	kwSeen := map[string]bool{}
	var keywords []grammar.Keyword

	var visitRule func(rule grammar.AbstractRule)
	visitRule = func(rule grammar.AbstractRule) {
		if rule == nil || visited[rule.Name()] {
			return
		}
		visited[rule.Name()] = true

		for node := range core.AllChildren(rule) {
			if kw, ok := node.(grammar.Keyword); ok {
				if !kwSeen[kw.Value()] {
					kwSeen[kw.Value()] = true
					keywords = append(keywords, kw)
				}
			}
			if rc, ok := node.(grammar.RuleCall); ok {
				ref := rc.Rule().Ref(context.Background())
				if ref != nil {
					visitRule(ref)
				}
			}
		}
	}

	visitRule(entryRule)

	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Value() < keywords[j].Value()
	})
	return keywords
}
