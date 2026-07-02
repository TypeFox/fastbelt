// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"fmt"
	"os"
	"testing"
	"unique"

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
		assert.Equal(t, "/statements@3/expression/left", mustNodePath(t, fc))
		assert.Equal(t, "materialPerUnit", fc.Callable().Text())
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

	t.Run("errorReporting/containerField-has-empty-string", func(t *testing.T) {
		// this test creates local copies of module, its first stmt, and its contained number literal
		//  and updates their container references accordingly

		// in order to test the error logging across some levels of nesting an additional root object is created,
		// which serves a container for the copied module,
		// while module's 'containerField' is set to the interned value of "" --> error!
		root := core.NewAstNode()

		module := *module.(*ModuleImpl)
		module.SetContainer(&root, unique.Make(""), 0)

		def := *stmts[0].(*DefinitionImpl)
		var field, index = def.ContainmentData()
		def.SetContainer(&module, field, index)

		expr := *def.expression.(*NumberLiteralImpl)
		field, index = expr.ContainmentData()
		expr.SetContainer(&def, field, index)

		path, err := expr.NodePath()
		assert.Zero(t, path)
		fmt.Println(err)
		assert.ErrorContains(
			t, err,
			"AstNodeBase.NodePath: error within node of type *arithmetics.ModuleImpl: cannot determine node path, 'containerField' is empty")
	})

	t.Run("errorReporting/containerField-has-zero-handle", func(t *testing.T) {
		// this test creates local copies of module, its first stmt, and its contained number literal
		//  and updates their container references accordingly

		module := *module.(*ModuleImpl)
		var fieldZero unique.Handle[string]

		def := *stmts[0].(*DefinitionImpl)
		def.SetContainer(&module, fieldZero, 0)

		expr := *def.expression.(*NumberLiteralImpl)
		field, index := expr.ContainmentData()
		expr.SetContainer(&def, field, index)

		path, err := expr.NodePath()
		assert.Zero(t, path)
		fmt.Println(err)
		assert.ErrorContains(
			t, err,
			"AstNodeBase.NodePath: error within node of type *arithmetics.DefinitionImpl: cannot determine node path, 'containerField' is empty")
	})
}

// TestGetByPath_PriceCalc is the inverse of TestNodePath_PriceCalc: for each of the 13
// FunctionCall cross-reference sites in price-calc.calc, it resolves the hardcoded path
// string back to an AST node via GetByPath and asserts pointer identity with the node
// obtained by direct AST navigation. Both the referencer (FunctionCall) and the
// referenced (AbstractDefinition) are covered per subtest.
func TestGetByPath_PriceCalc(t *testing.T) {
	doc := loadPriceCalcDoc(t)
	root := doc.Root()
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

		got, err := root.GetByPath("/statements@3/expression/left")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@0")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	t.Run("costPerUnit/laborPerUnit", func(t *testing.T) {
		def, ok := stmts[3].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Right().(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@3/expression/right")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@1")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	// ── def costOfGoodsSold: expectedNoOfSales * costPerUnit ─────────────────

	t.Run("costOfGoodsSold/expectedNoOfSales", func(t *testing.T) {
		def, ok := stmts[4].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Left().(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@4/expression/left")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@2")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	t.Run("costOfGoodsSold/costPerUnit", func(t *testing.T) {
		def, ok := stmts[4].(Definition)
		require.True(t, ok)
		binExpr, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := binExpr.Right().(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@4/expression/right")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@3")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
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

		got, err := root.GetByPath("/statements@7/expression/left/left/left")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@4")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
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

		got, err := root.GetByPath("/statements@7/expression/left/left/right")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@5")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
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

		got, err := root.GetByPath("/statements@7/expression/left/right")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@2")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	t.Run("netPrice/desiredProfitPerUnit", func(t *testing.T) {
		def, ok := stmts[7].(Definition)
		require.True(t, ok)
		outerAdd, ok := def.Expression().(BinaryExpression)
		require.True(t, ok)
		fc, ok := outerAdd.Right().(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@7/expression/right")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@6")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
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

		got, err := root.GetByPath("/statements@9/expression/left")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@9/args@0")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
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

		got, err := root.GetByPath("/statements@9/expression/right/right")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@9/args@1")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	// ── calcGrossListPrice(netPrice, vat) ────────────────────────────────────

	t.Run("evaluation/calcGrossListPrice", func(t *testing.T) {
		eval, ok := stmts[10].(Evaluation)
		require.True(t, ok)
		fc, ok := eval.Expression().(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@10/expression")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@9")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	t.Run("evaluation/netPrice", func(t *testing.T) {
		eval, ok := stmts[10].(Evaluation)
		require.True(t, ok)
		outerFC, ok := eval.Expression().(FunctionCall)
		require.True(t, ok)
		require.Len(t, outerFC.Args(), 2)
		fc, ok := outerFC.Args()[0].(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@10/expression/args@0")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@7")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	t.Run("evaluation/vat", func(t *testing.T) {
		eval, ok := stmts[10].(Evaluation)
		require.True(t, ok)
		outerFC, ok := eval.Expression().(FunctionCall)
		require.True(t, ok)
		require.Len(t, outerFC.Args(), 2)
		fc, ok := outerFC.Args()[1].(FunctionCall)
		require.True(t, ok)

		got, err := root.GetByPath("/statements@10/expression/args@1")
		require.NoError(t, err)
		assert.Same(t, fc, got)

		got, err = root.GetByPath("/statements@8")
		require.NoError(t, err)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), got)
	})

	t.Run("errorReporting/no-such-field", func(t *testing.T) {
		_, err := root.GetByPath("/statements@3/expressions/left")
		assert.ErrorContains(t, err, "DefinitionImpl.GetByPath: field 'expressions' does not exist in node '/statements@3' of type 'Definition'")
	})

	t.Run("errorReporting/field-is-primitive", func(t *testing.T) {
		_, err := root.GetByPath("/statements@0/name")
		assert.ErrorContains(t, err, "DefinitionImpl.GetByPath: field 'name' holds a primitive value instead of an ast node")
	})

	t.Run("errorReporting/field-is-reference", func(t *testing.T) {
		_, err := root.GetByPath("/statements@10/expression/callable")
		assert.ErrorContains(t, err, "FunctionCallImpl.GetByPath: field 'callable' is a cross-reference instead of a container field")
	})

	t.Run("errorReporting/field-is-empty", func(t *testing.T) {
		// create local copies of the relevant ast nodes first to avoid manipulating the shared ones
		module := *module.(*ModuleImpl)
		def := *stmts[0].(*DefinitionImpl)

		module.statements = make([]Statement, 1)
		module.statements[0] = &def

		field, index := def.ContainmentData()
		def.SetContainer(&module, field, index)

		// set the tested field to 'nil'
		def.expression = nil

		_, err := module.GetByPath("/statements@0/expression/left")
		assert.ErrorContains(t, err, "DefinitionImpl.GetByPath: field 'expression' is nil in node '/statements@0'")
	})

	t.Run("errorReporting/slice-index-out-of-bound-1", func(t *testing.T) {
		_, err := root.GetByPath("/statements@15/expression")
		assert.ErrorContains(t, err, "ModuleImpl.GetByPath: index 15 exceeds length of slice in 'statements' (length=11) in node ''")
	})

	t.Run("errorReporting/slice-index-out-of-bound-2", func(t *testing.T) {
		_, err := root.GetByPath("/statements@10/expression/args@7/expression")
		assert.ErrorContains(t, err, "FunctionCallImpl.GetByPath: index 7 exceeds length of slice in 'args' (length=2) in node '/statements@10/expression'")
	})

	t.Run("errorReporting/slice-item-is-nil", func(t *testing.T) {
		expr, err := root.GetByPath("/statements@10/expression/")
		require.NoError(t, err)
		fc, ok := expr.(FunctionCall)
		require.True(t, ok)

		// shamelessly manipulate the shared ast and add a 'nil' item, don't want to copy everything right now; shouldn't hurt
		fc.SetArgsItem(nil)
		_, err = root.GetByPath("/statements@10/expression/args@2/expression")

		assert.ErrorContains(t, err, "FunctionCallImpl.GetByPath: item 2 of slice in field 'args' is nil in node '/statements@10/expression'")
	})

	t.Run("errorReporting/slice-index-typo", func(t *testing.T) {
		_, err := root.GetByPath("/statements@1a/expression")
		assert.ErrorContains(t, err, "ModuleImpl.GetByPath: index '1a' is not a valid uint: strconv.Atoi: parsing \"1a\": invalid syntax")
	})
}
