// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// TestSingleOperandPassThrough asserts that an expression without operators
// is returned unchanged - no BinaryExpression wrapper is created.
func TestSingleOperandPassThrough(t *testing.T) {
	expr := parseExpression(t, "42")
	_, isNumber := expr.(NumberLiteral)
	assert.True(t, isNumber, "expected a NumberLiteral, got %T", expr)
}

// TestBinaryExpressionTokens verifies that after the infix tree rewrite every
// BinaryExpression node owns exactly its operator token, with correct
// Element, Kind, concrete token type, and a text range spanning its operands.
func TestBinaryExpressionTokens(t *testing.T) {
	expr := parseExpression(t, "1 + 2 * 3")

	add, ok := expr.(BinaryExpression)
	require.True(t, ok)
	mul, ok := add.Right().(BinaryExpression)
	require.True(t, ok)

	cases := []struct {
		node      BinaryExpression
		image     string
		tokenType *core.TokenType
		text      string
	}{
		{add, "+", Keyword_Plus, "1 + 2 * 3"},
		{mul, "*", Keyword_Asterisk, "2 * 3"},
	}
	for _, c := range cases {
		tokens := c.node.Tokens()
		require.Len(t, tokens, 1, "a binary node owns exactly its operator token")
		token := tokens[0]
		assert.Equal(t, c.image, token.Image)
		assert.Same(t, c.node, token.Element, "operator token must be assigned to its binary node")
		assert.Equal(t, BinaryExpression_BinaryExpressionOperator, token.Kind)
		assert.Same(t, c.tokenType, token.Type, "the token keeps its concrete member type, not the group type")
		assert.Equal(t, c.text, c.node.Text(), "node range spans left operand start to right operand end")
	}
}

// TestExponentiationRightAssociative verifies that the "right" group of
// the infix rule produces a right-leaning tree, while the left-associative
// groups keep producing left-leaning trees.
func TestExponentiationRightAssociative(t *testing.T) {
	expr := parseExpression(t, "2 ^ 3 ^ 2 ^ 1")
	assert.Equal(t, "(2 ^ (3 ^ (2 ^ 1)))", printExpression(expr))

	expr = parseExpression(t, "1 - 2 - 3")
	assert.Equal(t, "((1 - 2) - 3)", printExpression(expr))

	// Mixed with tighter ("%") and looser ("*") levels.
	expr = parseExpression(t, "1 * 2 ^ 3 % 4 ^ 5")
	assert.Equal(t, "(1 * (2 ^ ((3 % 4) ^ 5)))", printExpression(expr))
}

// TestDanglingOperator checks error tolerance: a trailing operator without a
// right operand reports a parse error but still produces a binary node.
func TestDanglingOperator(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse("module test 1 + ;")
	assert.NotEmpty(t, doc.Document.ParserErrors, "dangling operator should produce a parse error")

	module := test.MustFindNode[Module](doc)
	expr := module.Statements()[0].(Evaluation).Expression()
	add, ok := expr.(BinaryExpression)
	require.True(t, ok, "expected a BinaryExpression, got %T", expr)
	assert.Equal(t, "1", printExpression(add.Left()))
	assert.Equal(t, "+", add.Operator())
}

type operatorCompletionContributor struct {
	server.DefaultCompletionContributor
}

func (c *operatorCompletionContributor) CompletionForToken(_ context.Context, tt *core.TokenType, _ int, _ server.ContributorContext, accept server.CompletionAcceptor) {
	if tt != nil {
		// Show all keywords, including those in the BinaryExpressionOperator group,
		// which the default contributor suppresses.
		for _, keyword := range tt.MatchingTokens {
			accept(lsp.CompletionItem{
				Label: keyword.Name,
			})
		}
	}
}

// TestOperatorCompletion verifies that after an operand the completion engine
// offers every operator of the infix rule.
func TestOperatorCompletion(t *testing.T) {
	sc := service.NewContainer()
	SetupServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
	var contributor server.CompletionContributor = &operatorCompletionContributor{}
	service.Override(sc, contributor)
	sc.Seal()
	doc := test.New(t, sc).Parse("module test 1 <|cursor>")

	items := doc.CompletionItems("cursor")
	labels := make(map[string]bool, len(items))
	for _, item := range items {
		labels[item.Label] = true
	}
	for _, operator := range []string{"+", "-", "*", "/", "^", "%"} {
		assert.True(t, labels[operator], "expected operator %q in completion items, got %v", operator, labels)
	}
}

// TestAfterOperatorCompletion verifies that after an operator, the completion
// engine will offer to complete an expression again.
func TestAfterOperatorCompletion(t *testing.T) {
	sc := service.NewContainer()
	SetupServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
	var contributor server.CompletionContributor = &operatorCompletionContributor{}
	service.Override(sc, contributor)
	sc.Seal()
	doc := test.New(t, sc).Parse(`
	module test
	def x: 1;
	1 + 2 * <|cursor>
	`)

	items := doc.CompletionItems("cursor")
	labels := make(map[string]bool, len(items))
	for _, item := range items {
		labels[item.Label] = true
	}
	assert.True(t, labels["x"], "expected reference to definition 'x' in completion items, got %v", labels)
}
