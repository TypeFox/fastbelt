// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
)

func printExpression(expr Expression) string {
	if be, ok := expr.(BinaryExpression); ok {
		return "(" + printExpression(be.Left()) + " " + be.Operator() + " " + printExpression(be.Right()) + ")"
	} else if nl, ok := expr.(NumberLiteral); ok {
		return nl.Value()
	} else if fc, ok := expr.(FunctionCall); ok {
		return fc.Callable().Text()
	}
	return ""
}

func parseExpression(t testing.TB, text string) Expression {
	t.Helper()
	f := test.New(t, CreateServices())
	doc := f.Parse("module test " + text + ";").AssertNoParseErrors()
	doc.AssertState(core.DocStateParsed)

	module := test.MustFindNode[Module](doc)
	return module.Statements()[0].(Evaluation).Expression()
}

func TestSingleExpression(t *testing.T) {
	expr := parseExpression(t, "1")
	if got := printExpression(expr); got != "1" {
		t.Errorf("expected %q, got %q", "1", got)
	}
}

func TestBinaryExpression(t *testing.T) {
	expr := parseExpression(t, "1 + 2 ^ 3 * 4 % 5")
	expected := "(1 + ((2 ^ 3) * (4 % 5)))"
	if got := printExpression(expr); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestNestedExpression(t *testing.T) {
	expr := parseExpression(t, "(1 + 2) ^ 3")
	expected := "((1 + 2) ^ 3)"
	if got := printExpression(expr); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDefinitionsAndCalls(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		module test
		def two: 2;
		def root(x, y): x^(1/y);
		def sqrt(x): root(x, two);
		sqrt(16);
	`)
	doc.AssertState(core.DocStateParsed)
	doc.AssertNoParseErrors()

	allStmts := test.MustFindNode[Module](doc).Statements()
	require.Len(t, allStmts, 4)

	defTwo := test.MustFindNamedNode[Definition](doc, "two")
	defRoot := test.MustFindNamedNode[Definition](doc, "root")
	defSqrt := test.MustFindNamedNode[Definition](doc, "sqrt")
	eval := test.MustFindNode[Evaluation](doc)

	assert.Same(t, allStmts[0], defTwo)
	assert.Same(t, allStmts[1], defRoot)
	assert.Same(t, allStmts[2], defSqrt)
	assert.Same(t, allStmts[3], eval)
}
