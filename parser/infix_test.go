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

// parseInfixInput splits "1 + 2 * 3" into leaf nodes and operator tokens,
// assigning each distinct operator symbol a stable token type id.
func parseInfixInput(input string, typeIds map[string]int) (parts []*infixTestNode, operators []*core.Token) {
	for i, field := range strings.Fields(input) {
		if i%2 == 0 {
			parts = append(parts, &infixTestNode{name: field})
		} else {
			operators = append(operators, &core.Token{Image: field, Type: &core.TokenType{Id: typeIds[field], Name: field}})
		}
	}
	return parts, operators
}

func buildInfixTestNode(left *infixTestNode, operator *core.Token, right *infixTestNode) *infixTestNode {
	op := "?"
	if operator != nil {
		op = operator.Image
	}
	return &infixTestNode{left: left, right: right, op: op}
}

func TestBuildInfixTree(t *testing.T) {
	typeIds := map[string]int{"%": 1, "^": 2, "*": 3, "/": 4, "+": 5, "-": 6, "=": 7}
	// Mirrors the arithmetics grammar: "%" > "^" > "*" | "/" > "+" | "-",
	// plus a loosest right-associative "=" level.
	precedence := map[int]InfixPrecedence{
		typeIds["%"]: {Level: 0},
		typeIds["^"]: {Level: 1},
		typeIds["*"]: {Level: 2},
		typeIds["/"]: {Level: 2},
		typeIds["+"]: {Level: 3},
		typeIds["-"]: {Level: 3},
		typeIds["="]: {Level: 4, Right: true},
	}

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
	for input, expected := range cases {
		parts, operators := parseInfixInput(input, typeIds)
		tree := BuildInfixTree(parts, operators, precedence, buildInfixTestNode)
		assert.Equal(t, expected, tree.String(), "input %q", input)
	}
}

func TestBuildInfixTreeSinglePartIsReturnedUnchanged(t *testing.T) {
	part := &infixTestNode{name: "1"}
	tree := BuildInfixTree([]*infixTestNode{part}, nil, nil, buildInfixTestNode)
	assert.Same(t, part, tree)
}

func TestBuildInfixTreeEmptyParts(t *testing.T) {
	tree := BuildInfixTree(nil, nil, nil, buildInfixTestNode)
	assert.Nil(t, tree)
}

func TestBuildInfixTreeDanglingOperator(t *testing.T) {
	precedence := map[int]InfixPrecedence{1: {Level: 0}, 2: {Level: 0}}

	// "1 + 2 -" — the trailing operator has no right operand.
	parts := []*infixTestNode{{name: "1"}, {name: "2"}}
	operators := []*core.Token{
		{Image: "+", Type: &core.TokenType{Id: 1, Name: "+"}},
		{Image: "-", Type: &core.TokenType{Id: 2, Name: "-"}},
	}
	tree := BuildInfixTree(parts, operators, precedence, buildInfixTestNode)
	assert.Equal(t, "((1 + 2) - <nil>)", tree.String())
}

func TestBuildInfixTreeUnknownOperatorBindsLoosest(t *testing.T) {
	precedence := map[int]InfixPrecedence{1: {Level: 0}}
	parts := []*infixTestNode{{name: "1"}, {name: "2"}, {name: "3"}}
	operators := []*core.Token{
		{Image: "?", Type: &core.TokenType{Id: 99, Name: "?"}},
		{Image: "*", Type: &core.TokenType{Id: 1, Name: "*"}},
	}
	tree := BuildInfixTree(parts, operators, precedence, buildInfixTestNode)
	assert.Equal(t, "(1 ? (2 * 3))", tree.String())
}
