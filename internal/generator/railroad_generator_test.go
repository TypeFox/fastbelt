// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"encoding/xml"
	"slices"
	"strings"
	"testing"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

// Rule kinds are interleaved so source order differs from both alphabetical
// and grouped-by-kind order.
const railroadGeneratorFixtureGrammar = `grammar Test

interface Expr {}

interface BinaryExpression extends Expr {
    Left Expr
    Operator string
    Right Expr
}

interface Foo {
    Name string
    Items []Item
}

interface Item {
    Name string
}

entry Foo returns Foo:
    "begin" Name=ID Items+=Item*

composite Combo: "prefix" Item;

infix BinaryExpression on Primary:
    "+" | "-"

Item returns Item:
    Name=ID

Primary returns Expr:
    Value=NUMBER

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/
token NUMBER: /[0-9]+/
`

func TestGenerateRailroadDiagrams(t *testing.T) {
	sc := grammar.CreateServices()
	f := test.New(t, sc)
	doc := f.Parse(railroadGeneratorFixtureGrammar)
	doc.AssertNoParseErrors()

	g, ok := doc.Root().(grammar.Grammar)
	if !ok {
		t.Fatalf("document root = %T, want grammar.Grammar", doc.Root())
	}
	// Mirror the CLI, which desugars infix rules before running generators.
	if err := grammar.ExpandInfixRules(g); err != nil {
		t.Fatalf("ExpandInfixRules: %v", err)
	}

	diagrams := GenerateRailroadDiagrams(g)

	names := make([]string, 0, len(diagrams))
	for _, d := range diagrams {
		names = append(names, d.Name)
	}
	want := []string{"Foo", "Combo", "BinaryExpression", "Item", "Primary"}
	if !slices.Equal(names, want) {
		t.Fatalf("diagram names = %v, want %v (source order, no token rules)", names, want)
	}

	for _, d := range diagrams {
		var root struct{ XMLName xml.Name }
		if err := xml.Unmarshal([]byte(d.SVG), &root); err != nil {
			t.Errorf("diagram for %q is not well-formed XML: %v\n%s", d.Name, err, d.SVG)
			continue
		}
		if root.XMLName.Local != "svg" {
			t.Errorf("diagram for %q has root element %q, want svg", d.Name, root.XMLName.Local)
		}
	}
}

func TestGenerateRailroadIndexMarkdown(t *testing.T) {
	md := GenerateRailroadIndexMarkdown("mylang", []RailroadDiagram{{Name: "Foo"}, {Name: "Item"}})
	for _, want := range []string{"# Railroad diagrams for mylang", "## Foo", "![Foo](Foo.svg)", "## Item", "![Item](Item.svg)"} {
		if !strings.Contains(md, want) {
			t.Errorf("index markdown missing %q, got:\n%s", want, md)
		}
	}
	if strings.Index(md, "## Foo") > strings.Index(md, "## Item") {
		t.Errorf("index markdown should keep the given diagram order, got:\n%s", md)
	}
}
