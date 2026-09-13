// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"encoding/xml"
	"strings"
	"testing"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

const railroadGeneratorFixtureGrammar = `grammar Test

interface Foo {
    Name string
    Items []Item
}

interface Item {
    Name string
}

entry Foo returns Foo:
    "begin" Name=ID Items+=Item*

Item returns Item:
    Name=ID

composite Combo: "prefix" Item;

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/
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

	diagrams := GenerateRailroadDiagrams(g)

	for _, name := range []string{"Foo", "Item", "Combo"} {
		svg, ok := diagrams[name]
		if !ok {
			t.Errorf("no diagram generated for rule %q (got: %v)", name, diagramNames(diagrams))
			continue
		}
		var root struct{ XMLName xml.Name }
		if err := xml.Unmarshal([]byte(svg), &root); err != nil {
			t.Errorf("diagram for %q is not well-formed XML: %v\n%s", name, err, svg)
		}
	}

	if _, ok := diagrams["ID"]; ok {
		t.Error("a token rule should not have a diagram generated for it")
	}
}

func TestGenerateRailroadIndexMarkdown(t *testing.T) {
	md := GenerateRailroadIndexMarkdown("mylang", []string{"Foo", "Item"})
	for _, want := range []string{"# Railroad diagrams for mylang", "## Foo", "![Foo](Foo.svg)", "## Item", "![Item](Item.svg)"} {
		if !strings.Contains(md, want) {
			t.Errorf("index markdown missing %q, got:\n%s", want, md)
		}
	}
}

func diagramNames(diagrams map[string]string) []string {
	names := make([]string, 0, len(diagrams))
	for name := range diagrams {
		names = append(names, name)
	}
	return names
}
