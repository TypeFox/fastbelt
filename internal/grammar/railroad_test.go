// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"typefox.dev/fastbelt/internal/railroad"
	"typefox.dev/fastbelt/test"
)

func newRailroadFixture(t *testing.T) *test.Fixture {
	t.Helper()
	return test.New(t, CreateServices())
}

const mappingFixtureGrammar = `grammar Test

interface Foo {
    Name string
    Items []Item
    Broken string
}

interface Item {
    Name string
}

entry Foo returns Foo:
    "begin" Name=ID Items+=Item* ("end")?

Item returns Item:
    Name=ID | Name=STRING

ActionRule returns Foo:
    {Foo} Name=ID

BrokenRule returns Foo:
    Broken=UnknownRule

composite Combo: "prefix" Item;

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/
token STRING: /"[^"]*"/
`

func TestMapElement_SequenceWithCardinalityAndTokenVsRuleCall(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	foo := test.MustFindNamedNode[ParserRule](doc, "Foo")
	node := mapElement(doc.Ctx(), foo.Body())

	seq, ok := node.(*railroad.Sequence)
	if !ok {
		t.Fatalf("Foo's body = %T, want *railroad.Sequence", node)
	}
	if len(seq.Items) != 4 {
		t.Fatalf("Foo's body has %d items, want 4 (%v)", len(seq.Items), seq.Items)
	}

	if _, ok := seq.Items[0].(*railroad.Terminal); !ok {
		t.Errorf("item 0 (keyword \"begin\") = %T, want *railroad.Terminal", seq.Items[0])
	}
	if _, ok := seq.Items[1].(*railroad.Terminal); !ok {
		t.Errorf("item 1 (Name=ID, ID is a token) = %T, want *railroad.Terminal", seq.Items[1])
	}

	// Items+=Item*: Item is a parser rule, so it maps to NonTerminal, and
	// the '*' on the assignment wraps it in ZeroOrMore == Optional(OneOrMore).
	choice, ok := seq.Items[2].(*railroad.Choice)
	if !ok {
		t.Fatalf("item 2 (Items+=Item*) = %T, want *railroad.Choice (from ZeroOrMore)", seq.Items[2])
	}
	oneOrMore, ok := choice.Items[1].(*railroad.OneOrMore)
	if !ok {
		t.Fatalf("ZeroOrMore's wrapped item = %T, want *railroad.OneOrMore", choice.Items[1])
	}
	if _, ok := oneOrMore.Item.(*railroad.NonTerminal); !ok {
		t.Errorf("Item ruleCall = %T, want *railroad.NonTerminal", oneOrMore.Item)
	}

	// ("end")? : Optional(Terminal).
	optChoice, ok := seq.Items[3].(*railroad.Choice)
	if !ok {
		t.Fatalf("item 3 (\"end\")? = %T, want *railroad.Choice (from Optional)", seq.Items[3])
	}
	if _, ok := optChoice.Items[1].(*railroad.Terminal); !ok {
		t.Errorf("Optional's wrapped item = %T, want *railroad.Terminal", optChoice.Items[1])
	}
}

func TestMapElement_Alternatives(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	item := test.MustFindNamedNode[ParserRule](doc, "Item")
	node := mapElement(doc.Ctx(), item.Body())

	choice, ok := node.(*railroad.Choice)
	if !ok {
		t.Fatalf("Item's body = %T, want *railroad.Choice", node)
	}
	if len(choice.Items) != 2 {
		t.Fatalf("Item's body has %d alternatives, want 2", len(choice.Items))
	}
	for i, alt := range choice.Items {
		if _, ok := alt.(*railroad.Terminal); !ok {
			t.Errorf("alternative %d = %T, want *railroad.Terminal (both branches are token assignments)", i, alt)
		}
	}
}

func TestMapElement_Action(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[ParserRule](doc, "ActionRule")
	node := mapElement(doc.Ctx(), rule.Body())

	seq, ok := node.(*railroad.Sequence)
	if !ok {
		t.Fatalf("ActionRule's body = %T, want *railroad.Sequence", node)
	}
	if _, ok := seq.Items[0].(railroad.Skip); !ok {
		t.Errorf("action element = %T, want railroad.Skip (consumes no tokens)", seq.Items[0])
	}
}

func TestMapElement_UnresolvedRuleCallDoesNotPanic(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[ParserRule](doc, "BrokenRule")
	node := mapElement(doc.Ctx(), rule.Body())

	nt, ok := node.(*railroad.NonTerminal)
	if !ok {
		t.Fatalf("BrokenRule's body = %T, want *railroad.NonTerminal", node)
	}
	if !nt.Unresolved {
		t.Errorf("unresolved rule call should be flagged Unresolved, got %+v", nt)
	}
	if nt.Text != "UnknownRule" {
		t.Errorf("Text = %q, want %q", nt.Text, "UnknownRule")
	}
}

const crossRefFixtureGrammar = `grammar Test

interface Foo {
    Target *Foo
    Other *Foo
}

entry Foo returns Foo:
    "ref" Target=[Foo:ID] Other=[Foo]

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/
`

func TestMapElement_CrossRef(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(crossRefFixtureGrammar)
	doc.AssertNoParseErrors()

	foo := test.MustFindNamedNode[ParserRule](doc, "Foo")
	seq, ok := mapElement(doc.Ctx(), foo.Body()).(*railroad.Sequence)
	if !ok || len(seq.Items) != 3 {
		t.Fatalf("Foo's body = %#v, want a *railroad.Sequence of 3 items", seq)
	}

	// [Foo:ID]: the cross-reference renders as its token.
	if term, ok := seq.Items[1].(*railroad.Terminal); !ok || term.Text != "ID" {
		t.Errorf("[Foo:ID] = %#v, want *railroad.Terminal{Text: \"ID\"}", seq.Items[1])
	}

	// [Foo] (missing ":Rule", a validation error): still rendered, flagged.
	nt, ok := seq.Items[2].(*railroad.NonTerminal)
	if !ok {
		t.Fatalf("[Foo] = %T, want *railroad.NonTerminal", seq.Items[2])
	}
	if nt.Text != "Foo" || !nt.Unresolved {
		t.Errorf("[Foo] = %+v, want Text \"Foo\" and Unresolved", nt)
	}
}

func TestBuildRuleDiagram_CompositeRule(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	combo := test.MustFindNamedNode[CompositeRule](doc, "Combo")
	diagram, ok := BuildRuleDiagram(doc.Ctx(), combo)
	if !ok {
		t.Fatal("BuildRuleDiagram returned ok=false for a CompositeRule")
	}
	if diagram == nil {
		t.Fatal("BuildRuleDiagram returned a nil diagram with ok=true")
	}
}

func TestBuildRuleDiagram_TerminalRuleHasNoDiagram(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	idToken := test.MustFindNamedNode[TokenDecl](doc, "ID")
	_, ok := BuildRuleDiagram(doc.Ctx(), idToken)
	if ok {
		t.Error("BuildRuleDiagram should return ok=false for a token rule")
	}
}

const infixFixtureGrammar = `grammar Test2

interface Expr {}

interface BinaryExpression extends Expr {
    Left Expr
    Operator string
    Right Expr
}

entry Expr:
    PrimaryExpression

infix BinaryExpression on PrimaryExpression:
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

func TestMapInfixRule_SynthesizesOperandOperatorRepeatShape(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(infixFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[InfixRule](doc, "BinaryExpression")
	if rule.Body() != nil {
		t.Fatal("expected the InfixRule's Body to be nil before ExpandInfixRules runs")
	}

	node := mapInfixRule(doc.Ctx(), rule)
	seq, ok := node.(*railroad.Sequence)
	if !ok {
		t.Fatalf("infix rule diagram = %T, want *railroad.Sequence (operand, (operator operand)*)", node)
	}
	if len(seq.Items) != 2 {
		t.Fatalf("infix sequence has %d items, want 2", len(seq.Items))
	}
	if _, ok := seq.Items[0].(*railroad.NonTerminal); !ok {
		t.Errorf("operand = %T, want *railroad.NonTerminal (PrimaryExpression)", seq.Items[0])
	}

	repeatChoice, ok := seq.Items[1].(*railroad.Choice)
	if !ok {
		t.Fatalf("repeat wrapper = %T, want *railroad.Choice (from ZeroOrMore)", seq.Items[1])
	}
	oneOrMore, ok := repeatChoice.Items[1].(*railroad.OneOrMore)
	if !ok {
		t.Fatalf("ZeroOrMore's wrapped item = %T, want *railroad.OneOrMore", repeatChoice.Items[1])
	}
	innerSeq, ok := oneOrMore.Item.(*railroad.Sequence)
	if !ok {
		t.Fatalf("repeated body = %T, want *railroad.Sequence (operator, operand)", oneOrMore.Item)
	}
	if len(innerSeq.Items) != 2 {
		t.Fatalf("repeated body has %d items, want 2 (operator, operand)", len(innerSeq.Items))
	}
	// Operators across all precedence groups: "+", "-" (keywords) and MulOp
	// (a token group ruleCall) - 3 operators, so they render as a Choice.
	operatorChoice, ok := innerSeq.Items[0].(*railroad.Choice)
	if !ok {
		t.Fatalf("operator = %T, want *railroad.Choice (3 operators across precedence groups)", innerSeq.Items[0])
	}
	if len(operatorChoice.Items) != 3 {
		t.Fatalf("operatorChoice has %d items, want 3", len(operatorChoice.Items))
	}
	if _, ok := innerSeq.Items[1].(*railroad.NonTerminal); !ok {
		t.Errorf("second operand = %T, want *railroad.NonTerminal", innerSeq.Items[1])
	}

	// Reading Call()/Groups() must not have mutated the AST: Body() is
	// still nil (ExpandInfixRules was never called).
	if rule.Body() != nil {
		t.Error("mapInfixRule must not mutate the InfixRule's Body")
	}
}

func TestBuildRuleDiagram_InfixRule(t *testing.T) {
	f := newRailroadFixture(t)
	doc := f.Parse(infixFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[InfixRule](doc, "BinaryExpression")
	diagram, ok := BuildRuleDiagram(doc.Ctx(), rule)
	if !ok || diagram == nil {
		t.Fatalf("BuildRuleDiagram(InfixRule) = (%v, %v), want a non-nil diagram and ok=true", diagram, ok)
	}
	_ = diagram.SVG() // must render without panicking
}
