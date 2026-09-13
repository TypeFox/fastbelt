// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestDiagramSVGWellFormed(t *testing.T) {
	d := &Diagram{
		Item: &Sequence{Items: []Node{
			&Terminal{Text: "for"},
			&NonTerminal{Text: "Expr"},
			ZeroOrMore(&Terminal{Text: ","}),
		}},
	}
	svg := d.SVG()

	if !strings.HasPrefix(svg, "<svg") {
		t.Fatalf("expected SVG to start with <svg, got: %q", svg[:min(50, len(svg))])
	}

	var root struct{ XMLName xml.Name }
	if err := xml.Unmarshal([]byte(svg), &root); err != nil {
		t.Fatalf("SVG is not well-formed XML: %v\n%s", err, svg)
	}
	if root.XMLName.Local != "svg" {
		t.Fatalf("root element = %q, want svg", root.XMLName.Local)
	}
}

func TestDiagramSVGEmptyBody(t *testing.T) {
	// A rule with an entirely empty body (e.g. mid-edit) must still render
	// without panicking.
	d := &Diagram{Item: &Sequence{}}
	svg := d.SVG()

	var root struct{ XMLName xml.Name }
	if err := xml.Unmarshal([]byte(svg), &root); err != nil {
		t.Fatalf("SVG is not well-formed XML: %v\n%s", err, svg)
	}
}

func TestChoiceWithNoItemsDoesNotPanic(t *testing.T) {
	c := &Choice{Normal: 0, Items: nil}
	d := &Diagram{Item: c}
	_ = d.SVG() // must not panic
}
