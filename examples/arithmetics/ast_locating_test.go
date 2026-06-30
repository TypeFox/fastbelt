// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
)

func loadPriceCalcDoc(t *testing.T) *test.Doc {
	t.Helper()
	content, err := os.ReadFile("examples/price-calc.calc")
	require.NoError(t, err)
	f := test.New(t, CreateServices())
	doc := f.Parse(string(content))
	doc.AssertNoParseErrors()
	doc.AssertNoLinkingErrors()
	return doc
}

func mustNodePath(t *testing.T, node core.AstNode) string {
	t.Helper()
	path, err := node.NodePath()
	require.NoError(t, err)
	return path
}

// TestNodePath_PriceCalc verifies NodePath() for every FunctionCall (cross-reference)
// node in price-calc.calc. Each subtest covers one reference site and checks both
// the referencer (FunctionCall) and the referenced (AbstractDefinition) node paths.
func TestNodePath_PriceCalc(t *testing.T) {
	doc := loadPriceCalcDoc(t)

	module := test.MustFindNode[Module](doc)
	stmts := module.Statements()
	require.Len(t, stmts, 11)

	// ── def costPerUnit: materialPerUnit + laborPerUnit ───────────────────────

	t.Run("costPerUnit/materialPerUnit", func(t *testing.T) {
		def, ok := stmts[3].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Left().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "materialPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@3/expression/left", mustNodePath(t, fc))
		assert.Equal(t, "/statements@0", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("costPerUnit/laborPerUnit", func(t *testing.T) {
		def, ok := stmts[3].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Right().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@3/expression/right", mustNodePath(t, fc))
		assert.Equal(t, "laborPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@1", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	// ── def costOfGoodsSold: expectedNoOfSales * costPerUnit ─────────────────

	t.Run("costOfGoodsSold/expectedNoOfSales", func(t *testing.T) {
		def, ok := stmts[4].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Left().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@4/expression/left", mustNodePath(t, fc))
		assert.Equal(t, "expectedNoOfSales", fc.Callable().Text())
		assert.Equal(t, "/statements@2", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("costOfGoodsSold/costPerUnit", func(t *testing.T) {
		def, ok := stmts[4].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Right().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@4/expression/right", mustNodePath(t, fc))
		assert.Equal(t, "costPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@3", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	// ── def netPrice: (costOfGoodsSold + generalExpensesAndSales) / expectedNoOfSales + desiredProfitPerUnit
	//
	// AST (precedence: / binds tighter than +):
	//   Addition(+)
	//     left: Multiplication(/)
	//       left: Addition(+)  ← parenthesized
	//         left:  FC "costOfGoodsSold"
	//         right: FC "generalExpensesAndSales"
	//       right: FC "expectedNoOfSales"
	//     right: FC "desiredProfitPerUnit"

	t.Run("netPrice/costOfGoodsSold", func(t *testing.T) {
		def, ok := stmts[7].(Definition)
		require.True(t, ok)
		outerAdd, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		div, ok := outerAdd.Left().(BinaryExpression)
		require.True(t, ok)
		innerAdd, ok := div.Left().(BinaryExpression)
		require.True(t, ok)
		fc, ok := innerAdd.Left().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@7/expression/left/left/left", mustNodePath(t, fc))
		assert.Equal(t, "costOfGoodsSold", fc.Callable().Text())
		assert.Equal(t, "/statements@4", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("netPrice/generalExpensesAndSales", func(t *testing.T) {
		def, ok := stmts[7].(Definition)
		require.True(t, ok)
		outerAdd, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		div, ok := outerAdd.Left().(BinaryExpression)
		require.True(t, ok)
		innerAdd, ok := div.Left().(BinaryExpression)
		require.True(t, ok)
		fc, ok := innerAdd.Right().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@7/expression/left/left/right", mustNodePath(t, fc))
		assert.Equal(t, "generalExpensesAndSales", fc.Callable().Text())
		assert.Equal(t, "/statements@5", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("netPrice/expectedNoOfSales", func(t *testing.T) {
		def, ok := stmts[7].(Definition)
		require.True(t, ok)
		outerAdd, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		div, ok := outerAdd.Left().(BinaryExpression)
		require.True(t, ok)
		fc, ok := div.Right().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@7/expression/left/right", mustNodePath(t, fc))
		assert.Equal(t, "expectedNoOfSales", fc.Callable().Text())
		assert.Equal(t, "/statements@2", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("netPrice/desiredProfitPerUnit", func(t *testing.T) {
		def, ok := stmts[7].(Definition)
		require.True(t, ok)
		outerAdd, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := outerAdd.Right().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@7/expression/right", mustNodePath(t, fc))
		assert.Equal(t, "desiredProfitPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@6", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	// ── def calcGrossListPrice(net, tax): net / (1 - tax)
	//
	// AST:
	//   Multiplication(/)
	//     left:  FC "net"
	//     right: Addition(-)  ← parenthesized
	//       left:  NumberLiteral "1"
	//       right: FC "tax"

	t.Run("calcGrossListPrice/net", func(t *testing.T) {
		def, ok := stmts[9].(Definition)
		require.True(t, ok)
		div, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := div.Left().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@9/expression/left", mustNodePath(t, fc))
		assert.Equal(t, "net", fc.Callable().Text())
		assert.Equal(t, "/statements@9/args@0", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("calcGrossListPrice/tax", func(t *testing.T) {
		def, ok := stmts[9].(Definition)
		require.True(t, ok)
		div, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		sub, ok := div.Right().(BinaryExpression)
		require.True(t, ok)
		fc, ok := sub.Right().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@9/expression/right/right", mustNodePath(t, fc))
		assert.Equal(t, "tax", fc.Callable().Text())
		assert.Equal(t, "/statements@9/args@1", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	// ── calcGrossListPrice(netPrice, vat) ────────────────────────────────────

	t.Run("evaluation/calcGrossListPrice", func(t *testing.T) {
		eval, ok := stmts[10].(Evaluation)
		require.True(t, ok)
		fc, ok := eval.Expression().(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@10/expression", mustNodePath(t, fc))
		assert.Equal(t, "calcGrossListPrice", fc.Callable().Text())
		assert.Equal(t, "/statements@9", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("evaluation/netPrice", func(t *testing.T) {
		eval, ok := stmts[10].(Evaluation)
		require.True(t, ok)
		outerFC, ok := eval.Expression().(FunctionCall)
		require.True(t, ok)
		require.Len(t, outerFC.Args(), 2)
		fc, ok := outerFC.Args()[0].(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@10/expression/args@0", mustNodePath(t, fc))
		assert.Equal(t, "netPrice", fc.Callable().Text())
		assert.Equal(t, "/statements@7", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("evaluation/vat", func(t *testing.T) {
		eval, ok := stmts[10].(Evaluation)
		require.True(t, ok)
		outerFC, ok := eval.Expression().(FunctionCall)
		require.True(t, ok)
		require.Len(t, outerFC.Args(), 2)
		fc, ok := outerFC.Args()[1].(FunctionCall)
		require.True(t, ok)
		assert.Equal(t, "/statements@10/expression/args@1", mustNodePath(t, fc))
		assert.Equal(t, "vat", fc.Callable().Text())
		assert.Equal(t, "/statements@8", mustNodePath(t, fc.Callable().Ref(doc.Ctx())))
	})
}
