// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import (
	"testing"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

func newFixture(t *testing.T) *test.Fixture {
	t.Helper()
	return test.New(t, grammar.CreateServices())
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
	f := newFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	foo := test.MustFindNamedNode[grammar.ParserRule](doc, "Foo")
	node := MapElement(doc.Ctx(), foo.Body())

	seq, ok := node.(*Sequence)
	if !ok {
		t.Fatalf("Foo's body = %T, want *Sequence", node)
	}
	if len(seq.Items) != 4 {
		t.Fatalf("Foo's body has %d items, want 4 (%v)", len(seq.Items), seq.Items)
	}

	if _, ok := seq.Items[0].(*Terminal); !ok {
		t.Errorf("item 0 (keyword \"begin\") = %T, want *Terminal", seq.Items[0])
	}
	if _, ok := seq.Items[1].(*Terminal); !ok {
		t.Errorf("item 1 (Name=ID, ID is a token) = %T, want *Terminal", seq.Items[1])
	}

	// Items+=Item*: Item is a parser rule, so it maps to NonTerminal, and
	// the '*' on the assignment wraps it in ZeroOrMore == Optional(OneOrMore).
	choice, ok := seq.Items[2].(*Choice)
	if !ok {
		t.Fatalf("item 2 (Items+=Item*) = %T, want *Choice (from ZeroOrMore)", seq.Items[2])
	}
	oneOrMore, ok := choice.Items[1].(*OneOrMore)
	if !ok {
		t.Fatalf("ZeroOrMore's wrapped item = %T, want *OneOrMore", choice.Items[1])
	}
	if _, ok := oneOrMore.Item.(*NonTerminal); !ok {
		t.Errorf("Item ruleCall = %T, want *NonTerminal", oneOrMore.Item)
	}

	// ("end")? : Optional(Terminal).
	optChoice, ok := seq.Items[3].(*Choice)
	if !ok {
		t.Fatalf("item 3 (\"end\")? = %T, want *Choice (from Optional)", seq.Items[3])
	}
	if _, ok := optChoice.Items[1].(*Terminal); !ok {
		t.Errorf("Optional's wrapped item = %T, want *Terminal", optChoice.Items[1])
	}
}

func TestMapElement_Alternatives(t *testing.T) {
	f := newFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	item := test.MustFindNamedNode[grammar.ParserRule](doc, "Item")
	node := MapElement(doc.Ctx(), item.Body())

	choice, ok := node.(*Choice)
	if !ok {
		t.Fatalf("Item's body = %T, want *Choice", node)
	}
	if len(choice.Items) != 2 {
		t.Fatalf("Item's body has %d alternatives, want 2", len(choice.Items))
	}
	for i, alt := range choice.Items {
		if _, ok := alt.(*Terminal); !ok {
			t.Errorf("alternative %d = %T, want *Terminal (both branches are token assignments)", i, alt)
		}
	}
}

func TestMapElement_Action(t *testing.T) {
	f := newFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[grammar.ParserRule](doc, "ActionRule")
	node := MapElement(doc.Ctx(), rule.Body())

	seq, ok := node.(*Sequence)
	if !ok {
		t.Fatalf("ActionRule's body = %T, want *Sequence", node)
	}
	if _, ok := seq.Items[0].(Skip); !ok {
		t.Errorf("action element = %T, want Skip (consumes no tokens)", seq.Items[0])
	}
}

func TestMapElement_UnresolvedRuleCallDoesNotPanic(t *testing.T) {
	f := newFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[grammar.ParserRule](doc, "BrokenRule")
	node := MapElement(doc.Ctx(), rule.Body())

	nt, ok := node.(*NonTerminal)
	if !ok {
		t.Fatalf("BrokenRule's body = %T, want *NonTerminal", node)
	}
	if !nt.Unresolved {
		t.Errorf("unresolved rule call should be flagged Unresolved, got %+v", nt)
	}
	if nt.Text != "UnknownRule" {
		t.Errorf("Text = %q, want %q", nt.Text, "UnknownRule")
	}
}

func TestBuildRuleDiagram_CompositeRule(t *testing.T) {
	f := newFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	combo := test.MustFindNamedNode[grammar.CompositeRule](doc, "Combo")
	diagram, ok := BuildRuleDiagram(doc.Ctx(), combo)
	if !ok {
		t.Fatal("BuildRuleDiagram returned ok=false for a CompositeRule")
	}
	if diagram == nil {
		t.Fatal("BuildRuleDiagram returned a nil diagram with ok=true")
	}
}

func TestBuildRuleDiagram_TerminalRuleHasNoDiagram(t *testing.T) {
	f := newFixture(t)
	doc := f.Parse(mappingFixtureGrammar)
	doc.AssertNoParseErrors()

	idToken := test.MustFindNamedNode[grammar.Token](doc, "ID")
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
	f := newFixture(t)
	doc := f.Parse(infixFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[grammar.InfixRule](doc, "BinaryExpression")
	if rule.Body() != nil {
		t.Fatal("expected the InfixRule's Body to be nil before ExpandInfixRules runs")
	}

	node := mapInfixRule(doc.Ctx(), rule)
	seq, ok := node.(*Sequence)
	if !ok {
		t.Fatalf("infix rule diagram = %T, want *Sequence (operand, (operator operand)*)", node)
	}
	if len(seq.Items) != 2 {
		t.Fatalf("infix sequence has %d items, want 2", len(seq.Items))
	}
	if _, ok := seq.Items[0].(*NonTerminal); !ok {
		t.Errorf("operand = %T, want *NonTerminal (PrimaryExpression)", seq.Items[0])
	}

	repeatChoice, ok := seq.Items[1].(*Choice)
	if !ok {
		t.Fatalf("repeat wrapper = %T, want *Choice (from ZeroOrMore)", seq.Items[1])
	}
	oneOrMore, ok := repeatChoice.Items[1].(*OneOrMore)
	if !ok {
		t.Fatalf("ZeroOrMore's wrapped item = %T, want *OneOrMore", repeatChoice.Items[1])
	}
	innerSeq, ok := oneOrMore.Item.(*Sequence)
	if !ok {
		t.Fatalf("repeated body = %T, want *Sequence (operator, operand)", oneOrMore.Item)
	}
	if len(innerSeq.Items) != 2 {
		t.Fatalf("repeated body has %d items, want 2 (operator, operand)", len(innerSeq.Items))
	}
	// Operators across all precedence groups: "+", "-" (keywords) and MulOp
	// (a token group ruleCall) - 3 operators, so they render as a Choice.
	operatorChoice, ok := innerSeq.Items[0].(*Choice)
	if !ok {
		t.Fatalf("operator = %T, want *Choice (3 operators across precedence groups)", innerSeq.Items[0])
	}
	if len(operatorChoice.Items) != 3 {
		t.Fatalf("operatorChoice has %d items, want 3", len(operatorChoice.Items))
	}
	if _, ok := innerSeq.Items[1].(*NonTerminal); !ok {
		t.Errorf("second operand = %T, want *NonTerminal", innerSeq.Items[1])
	}

	// Reading Call()/Groups() must not have mutated the AST: Body() is
	// still nil (ExpandInfixRules was never called).
	if rule.Body() != nil {
		t.Error("mapInfixRule must not mutate the InfixRule's Body")
	}
}

func TestBuildRuleDiagram_InfixRule(t *testing.T) {
	f := newFixture(t)
	doc := f.Parse(infixFixtureGrammar)
	doc.AssertNoParseErrors()

	rule := test.MustFindNamedNode[grammar.InfixRule](doc, "BinaryExpression")
	diagram, ok := BuildRuleDiagram(doc.Ctx(), rule)
	if !ok || diagram == nil {
		t.Fatalf("BuildRuleDiagram(InfixRule) = (%v, %v), want a non-nil diagram and ok=true", diagram, ok)
	}
	_ = diagram.SVG() // must render without panicking
}
