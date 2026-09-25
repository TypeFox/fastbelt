// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
)

// RailroadDiagram is the rendered SVG railroad diagram of one grammar rule.
type RailroadDiagram struct {
	Name string
	SVG  string
}

// GenerateRailroadDiagrams renders a standalone SVG railroad diagram for
// every parser, composite, and infix rule of grammr, in source order.
// Rule kinds with no meaningful concrete-syntax diagram (terminal/token
// rules) are omitted.
func GenerateRailroadDiagrams(grammr grammar.Grammar) []RailroadDiagram {
	var rules []grammar.AbstractRule
	for _, rule := range grammr.Rules() {
		rules = append(rules, rule)
	}
	for _, rule := range grammr.Composites() {
		rules = append(rules, rule)
	}
	for _, rule := range grammr.InfixRules() {
		rules = append(rules, rule)
	}
	slices.SortFunc(rules, func(a, b grammar.AbstractRule) int {
		return cmp.Compare(a.TextRange().Start, b.TextRange().Start)
	})

	ctx := context.Background()
	out := make([]RailroadDiagram, 0, len(rules))
	for _, rule := range rules {
		if diagram, ok := grammar.BuildRuleDiagram(ctx, rule); ok {
			out = append(out, RailroadDiagram{Name: rule.Name(), SVG: diagram.SVG()})
		}
	}
	return out
}

// GenerateRailroadIndexMarkdown renders a Markdown page linking every
// diagram (by relative <name>.svg file path) for convenient one-page
// browsing.
func GenerateRailroadIndexMarkdown(packageName string, diagrams []RailroadDiagram) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Railroad diagrams for %s\n\n", packageName)
	for _, d := range diagrams {
		fmt.Fprintf(&b, "## %s\n\n![%s](%s.svg)\n\n", d.Name, d.Name, d.Name)
	}
	return b.String()
}
