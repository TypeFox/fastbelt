// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammarhover

import (
	"encoding/base64"
	"encoding/xml"
	"strings"
	"testing"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

const dataImagePrefix = "data:image/svg+xml;base64,"

func newHoverFixture(t *testing.T) (*service.Container, *test.Fixture) {
	t.Helper()
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	service.Put(sc, NewHoverProvider(sc))
	server.SetupDefaultServices(sc)
	sc.Seal()
	return sc, test.New(t, sc)
}

func hoverAt(t *testing.T, sc *service.Container, doc *test.Doc, label string) *lsp.Hover {
	t.Helper()
	rng, err := doc.MarkerRange(label)
	if err != nil {
		t.Fatalf("marker %q: %v", label, err)
	}
	hoverProvider := service.MustGet[server.HoverProvider](sc)
	result, err := hoverProvider.HandleHoverRequest(doc.Ctx(), &lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
			Position:     doc.Document.TextDoc.PositionAt(int(rng.Start)),
		},
	})
	if err != nil {
		t.Fatalf("hover request at %q failed: %v", label, err)
	}
	return result
}

func extractSVG(t *testing.T, markdown string) string {
	t.Helper()
	idx := strings.Index(markdown, dataImagePrefix)
	if idx < 0 {
		t.Fatalf("hover content does not contain %q: %s", dataImagePrefix, markdown)
	}
	encoded := markdown[idx+len(dataImagePrefix):]
	if end := strings.IndexByte(encoded, ')'); end >= 0 {
		encoded = encoded[:end]
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}
	var root struct{ XMLName xml.Name }
	if err := xml.Unmarshal(decoded, &root); err != nil {
		t.Fatalf("embedded SVG is not well-formed XML: %v\n%s", err, decoded)
	}
	if root.XMLName.Local != "svg" {
		t.Fatalf("embedded root element = %q, want svg", root.XMLName.Local)
	}
	return string(decoded)
}

const railroadHoverFixtureGrammar = `grammar Test

interface Foo {
    Name string
    Items []Item
}

interface Item {
    Name string
}

entry Foo returns Foo:
    "begin" Name=<|idUsage:ID|> Items+=<|itemUsage:Item|>*

<|itemDecl:Item|> returns Item:
    Name=ID

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/
`

func TestRailroadHover_RuleDeclarationEmbedsSVG(t *testing.T) {
	sc, f := newHoverFixture(t)
	doc := f.Parse(railroadHoverFixtureGrammar)
	doc.AssertNoParseErrors()

	hover := hoverAt(t, sc, doc, "itemDecl")
	if hover == nil {
		t.Fatal("expected a hover result at the rule declaration, got nil")
	}
	extractSVG(t, hover.Contents.Value)
}

func TestRailroadHover_UsageResolvesToSameContentAsDeclaration(t *testing.T) {
	sc, f := newHoverFixture(t)
	doc := f.Parse(railroadHoverFixtureGrammar)
	doc.AssertNoParseErrors()

	decl := hoverAt(t, sc, doc, "itemDecl")
	usage := hoverAt(t, sc, doc, "itemUsage")
	if decl == nil || usage == nil {
		t.Fatalf("expected hover results at both positions, got decl=%v usage=%v", decl, usage)
	}
	if decl.Contents.Value != usage.Contents.Value {
		t.Errorf("hover content differs between declaration and usage:\ndecl:  %s\nusage: %s", decl.Contents.Value, usage.Contents.Value)
	}
}

func TestRailroadHover_TokenRuleHasNoDiagram(t *testing.T) {
	sc, f := newHoverFixture(t)
	doc := f.Parse(railroadHoverFixtureGrammar)
	doc.AssertNoParseErrors()

	hover := hoverAt(t, sc, doc, "idUsage")
	if hover != nil && strings.Contains(hover.Contents.Value, "data:image") {
		t.Errorf("hovering a token rule usage should not embed a diagram, got: %s", hover.Contents.Value)
	}
}

const railroadHoverInfixFixtureGrammar = `grammar Test2

interface Expr {}

interface BinaryExpression extends Expr {
    Left Expr
    Operator string
    Right Expr
}

entry Expr:
    PrimaryExpression

infix <|infixDecl:BinaryExpression|> on PrimaryExpression:
    "+" | "-"
    > MulOp

token group MulOp {
    "*"
    "/"
}

PrimaryExpression returns Expr:
    Value=NUMBER

token NUMBER: /[0-9]+/
`

func TestRailroadHover_InfixRuleEmbedsSVG(t *testing.T) {
	sc, f := newHoverFixture(t)
	doc := f.Parse(railroadHoverInfixFixtureGrammar)
	doc.AssertNoParseErrors()

	hover := hoverAt(t, sc, doc, "infixDecl")
	if hover == nil {
		t.Fatal("expected a hover result at the infix rule declaration, got nil")
	}
	extractSVG(t, hover.Contents.Value)
}
