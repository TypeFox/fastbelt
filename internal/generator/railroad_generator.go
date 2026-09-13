// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"fmt"
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/internal/railroad"
)

// GenerateRailroadDiagrams renders a standalone SVG railroad diagram for
// every parser, composite, and infix rule of grammr, keyed by rule name.
// Rule kinds with no meaningful concrete-syntax diagram (terminal/token
// rules) are omitted.
func GenerateRailroadDiagrams(grammr grammar.Grammar) map[string]string {
	ctx := context.Background()
	out := make(map[string]string)
	for _, rule := range grammr.Rules() {
		if diagram, ok := railroad.BuildRuleDiagram(ctx, rule); ok {
			out[rule.Name()] = diagram.SVG()
		}
	}
	for _, rule := range grammr.Composites() {
		if diagram, ok := railroad.BuildRuleDiagram(ctx, rule); ok {
			out[rule.Name()] = diagram.SVG()
		}
	}
	for _, rule := range grammr.InfixRules() {
		if diagram, ok := railroad.BuildRuleDiagram(ctx, rule); ok {
			out[rule.Name()] = diagram.SVG()
		}
	}
	return out
}

// GenerateRailroadIndexMarkdown renders a Markdown page linking every named
// diagram (by relative <name>.svg file path) for convenient one-page
// browsing.
func GenerateRailroadIndexMarkdown(packageName string, ruleNames []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Railroad diagrams for %s\n\n", packageName)
	for _, name := range ruleNames {
		fmt.Fprintf(&b, "## %s\n\n![%s](%s.svg)\n\n", name, name, name)
	}
	return b.String()
}
