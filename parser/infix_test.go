// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
)

type infixTestNode struct {
	core.AstNodeBase
	name        string
	left, right *infixTestNode
	op          string
}

func (n *infixTestNode) String() string {
	if n == nil {
		return "<nil>"
	}
	if n.op == "" {
		return n.name
	}
	return "(" + n.left.String() + " " + n.op + " " + n.right.String() + ")"
}

// Mirrors the arithmetics grammar: "%" > "^" > "*" | "/" > "+" | "-",
// plus a loosest right-associative "=" level.
var infixTestTypeIds = map[string]int{"%": 1, "^": 2, "*": 3, "/": 4, "+": 5, "-": 6, "=": 7}

var infixTestTable = []InfixPrecedence{
	1: {Level: 0},
	2: {Level: 1},
	3: {Level: 2},
	4: {Level: 2},
	5: {Level: 3},
	6: {Level: 3},
	7: {Level: 4, Right: true},
}

func infixTestOperator(image string, typeId int) *core.Token {
	return &core.Token{Image: image, Type: &core.TokenType{Id: typeId, Name: image}}
}

func buildInfixTestNode(left *infixTestNode, operator *core.Token, right *infixTestNode) *infixTestNode {
	op := "?"
	if operator != nil {
		op = operator.Image
	}
	return &infixTestNode{left: left, right: right, op: op}
}

// foldInfixInput feeds "1 + 2 * 3" into fold: even fields are operands, odd
// fields are operators looked up in infixTestTypeIds.
func foldInfixInput(fold *InfixFold[*infixTestNode], input string) *infixTestNode {
	frame := fold.Begin()
	for i, field := range strings.Fields(input) {
		if i%2 == 0 {
			fold.Operand(&infixTestNode{name: field})
		} else {
			fold.Operator(frame, infixTestOperator(field, infixTestTypeIds[field]), infixTestTable, buildInfixTestNode)
		}
	}
	return fold.End(frame, buildInfixTestNode)
}

func TestInfixFold(t *testing.T) {
	cases := map[string]string{
		"1":                 "1",
		"1 + 2":             "(1 + 2)",
		"1 + 2 * 3":         "(1 + (2 * 3))",
		"1 * 2 + 3":         "((1 * 2) + 3)",
		"1 + 2 + 3":         "((1 + 2) + 3)",
		"1 - 2 - 3":         "((1 - 2) - 3)",
		"1 + 2 * 3 / 4":     "(1 + ((2 * 3) / 4))",
		"1 + 2 ^ 3 * 4 % 5": "(1 + ((2 ^ 3) * (4 % 5)))",
		"a = b = c":         "(a = (b = c))",
		"a = b + c = d":     "(a = ((b + c) = d))",
	}
	var fold InfixFold[*infixTestNode]
	for input, expected := range cases {
		tree := foldInfixInput(&fold, input)
		assert.Equal(t, expected, tree.String(), "input %q", input)
		assert.Empty(t, fold.operands, "input %q", input)
		assert.Empty(t, fold.ops, "input %q", input)
	}
}

func TestInfixFoldSingleOperandIsReturnedUnchanged(t *testing.T) {
	var fold InfixFold[*infixTestNode]
	part := &infixTestNode{name: "1"}
	frame := fold.Begin()
	fold.Operand(part)
	assert.Same(t, part, fold.End(frame, buildInfixTestNode))
}

func TestInfixFoldNoOperands(t *testing.T) {
	var fold InfixFold[*infixTestNode]
	frame := fold.Begin()
	assert.Nil(t, fold.End(frame, buildInfixTestNode))
}

func TestInfixFoldDanglingOperator(t *testing.T) {
	// "1 + 2 -" — the trailing operator has no right operand.
	var fold InfixFold[*infixTestNode]
	tree := foldInfixInput(&fold, "1 + 2 -")
	assert.Equal(t, "((1 + 2) - <nil>)", tree.String())
}

func TestInfixFoldMissingLeftOperand(t *testing.T) {
	var fold InfixFold[*infixTestNode]
	frame := fold.Begin()
	fold.Operator(frame, infixTestOperator("+", 5), infixTestTable, buildInfixTestNode)
	fold.Operand(&infixTestNode{name: "1"})
	assert.Equal(t, "(<nil> + 1)", fold.End(frame, buildInfixTestNode).String())
}

func TestInfixFoldUnknownOperatorBindsLoosest(t *testing.T) {
	var fold InfixFold[*infixTestNode]
	frame := fold.Begin()
	fold.Operand(&infixTestNode{name: "1"})
	fold.Operator(frame, infixTestOperator("?", 99), infixTestTable, buildInfixTestNode)
	fold.Operand(&infixTestNode{name: "2"})
	fold.Operator(frame, infixTestOperator("*", 3), infixTestTable, buildInfixTestNode)
	fold.Operand(&infixTestNode{name: "3"})
	assert.Equal(t, "(1 ? (2 * 3))", fold.End(frame, buildInfixTestNode).String())
}

func TestInfixFoldNilOperatorBindsLoosest(t *testing.T) {
	var fold InfixFold[*infixTestNode]
	frame := fold.Begin()
	fold.Operand(&infixTestNode{name: "1"})
	fold.Operator(frame, nil, infixTestTable, buildInfixTestNode)
	fold.Operand(&infixTestNode{name: "2"})
	fold.Operator(frame, infixTestOperator("*", 3), infixTestTable, buildInfixTestNode)
	fold.Operand(&infixTestNode{name: "3"})
	assert.Equal(t, "(1 ? (2 * 3))", fold.End(frame, buildInfixTestNode).String())
}

// TestInfixFoldNestedFrame mirrors a parenthesized operand: the inner
// invocation runs on the same fold while the outer one is still open.
func TestInfixFoldNestedFrame(t *testing.T) {
	var fold InfixFold[*infixTestNode]
	outer := fold.Begin()
	fold.Operand(&infixTestNode{name: "1"})
	fold.Operator(outer, infixTestOperator("*", 3), infixTestTable, buildInfixTestNode)
	inner := foldInfixInput(&fold, "2 + 3")
	assert.Equal(t, "(2 + 3)", inner.String())
	fold.Operand(inner)
	fold.Operator(outer, infixTestOperator("+", 5), infixTestTable, buildInfixTestNode)
	fold.Operand(&infixTestNode{name: "4"})
	tree := fold.End(outer, buildInfixTestNode)
	assert.Equal(t, "((1 * (2 + 3)) + 4)", tree.String())
	assert.Empty(t, fold.operands)
	assert.Empty(t, fold.ops)
}
