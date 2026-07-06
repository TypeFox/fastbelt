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

func mustConvert[T core.AstNode](t *testing.T, node core.AstNode) T {
	t.Helper()
	converted, ok := node.(T)
	require.True(t, ok)
	return converted
}

// TestNodePath_PriceCalc verifies NodePath() for every FunctionCall (cross-reference)
// node in price-calc.calc. Each subtest covers one reference site and checks both
// the referencer (FunctionCall) and the referenced (AbstractDefinition) node paths.
func TestNodePath_PriceCalc(t *testing.T) {
	doc := loadPriceCalcDoc(t)

	module := test.MustFindNode[Module](doc)
	stmts := module.Statements()
	require.Len(t, stmts, 11)

	mustPathOf := func(node core.AstNode) string {
		t.Helper()
		path, err := core.PathOf(node)
		require.NoError(t, err)
		return path.String()
	}

	// ── def costPerUnit: materialPerUnit + laborPerUnit ───────────────────────

	t.Run("costPerUnit/materialPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[3])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Left())
		assert.Equal(t, "/statements@3/expression/left", mustPathOf(fc))
		assert.Equal(t, "materialPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@0", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("costPerUnit/laborPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[3])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Right())
		assert.Equal(t, "/statements@3/expression/right", mustPathOf(fc))
		assert.Equal(t, "laborPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@1", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	// ── def costOfGoodsSold: expectedNoOfSales * costPerUnit ─────────────────

	t.Run("costOfGoodsSold/expectedNoOfSales", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[4])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Left())
		assert.Equal(t, "/statements@4/expression/left", mustPathOf(fc))
		assert.Equal(t, "expectedNoOfSales", fc.Callable().Text())
		assert.Equal(t, "/statements@2", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("costOfGoodsSold/costPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[4])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Right())
		assert.Equal(t, "/statements@4/expression/right", mustPathOf(fc))
		assert.Equal(t, "costPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@3", mustPathOf(fc.Callable().Ref(doc.Ctx())))
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
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		div := mustConvert[BinaryExpression](t, outerAdd.Left())
		innerAdd := mustConvert[BinaryExpression](t, div.Left())
		fc := mustConvert[FunctionCall](t, innerAdd.Left())
		assert.Equal(t, "/statements@7/expression/left/left/left", mustPathOf(fc))
		assert.Equal(t, "costOfGoodsSold", fc.Callable().Text())
		assert.Equal(t, "/statements@4", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("netPrice/generalExpensesAndSales", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		div := mustConvert[BinaryExpression](t, outerAdd.Left())
		innerAdd := mustConvert[BinaryExpression](t, div.Left())
		fc := mustConvert[FunctionCall](t, innerAdd.Right())
		assert.Equal(t, "/statements@7/expression/left/left/right", mustPathOf(fc))
		assert.Equal(t, "generalExpensesAndSales", fc.Callable().Text())
		assert.Equal(t, "/statements@5", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("netPrice/expectedNoOfSales", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		div := mustConvert[BinaryExpression](t, outerAdd.Left())
		fc := mustConvert[FunctionCall](t, div.Right())
		assert.Equal(t, "/statements@7/expression/left/right", mustPathOf(fc))
		assert.Equal(t, "expectedNoOfSales", fc.Callable().Text())
		assert.Equal(t, "/statements@2", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("netPrice/desiredProfitPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, outerAdd.Right())
		assert.Equal(t, "/statements@7/expression/right", mustPathOf(fc))
		assert.Equal(t, "desiredProfitPerUnit", fc.Callable().Text())
		assert.Equal(t, "/statements@6", mustPathOf(fc.Callable().Ref(doc.Ctx())))
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
		def := mustConvert[Definition](t, stmts[9])
		div := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, div.Left())
		assert.Equal(t, "/statements@9/expression/left", mustPathOf(fc))
		assert.Equal(t, "net", fc.Callable().Text())
		assert.Equal(t, "/statements@9/args@0", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("calcGrossListPrice/tax", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[9])
		div := mustConvert[BinaryExpression](t, def.Expression())
		sub := mustConvert[BinaryExpression](t, div.Right())
		fc := mustConvert[FunctionCall](t, sub.Right())
		assert.Equal(t, "/statements@9/expression/right/right", mustPathOf(fc))
		assert.Equal(t, "tax", fc.Callable().Text())
		assert.Equal(t, "/statements@9/args@1", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	// ── calcGrossListPrice(netPrice, vat) ────────────────────────────────────

	t.Run("evaluation/calcGrossListPrice", func(t *testing.T) {
		eval := mustConvert[Evaluation](t, stmts[10])
		fc := mustConvert[FunctionCall](t, eval.Expression())
		assert.Equal(t, "/statements@10/expression", mustPathOf(fc))
		assert.Equal(t, "calcGrossListPrice", fc.Callable().Text())
		assert.Equal(t, "/statements@9", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("evaluation/netPrice", func(t *testing.T) {
		eval := mustConvert[Evaluation](t, stmts[10])
		outerFC := mustConvert[FunctionCall](t, eval.Expression())
		require.Len(t, outerFC.Args(), 2)
		fc := mustConvert[FunctionCall](t, outerFC.Args()[0])
		assert.Equal(t, "/statements@10/expression/args@0", mustPathOf(fc))
		assert.Equal(t, "netPrice", fc.Callable().Text())
		assert.Equal(t, "/statements@7", mustPathOf(fc.Callable().Ref(doc.Ctx())))
	})

	t.Run("evaluation/vat", func(t *testing.T) {
		eval := mustConvert[Evaluation](t, stmts[10])
		outerFC := mustConvert[FunctionCall](t, eval.Expression())
		require.Len(t, outerFC.Args(), 2)
		fc := mustConvert[FunctionCall](t, outerFC.Args()[1])
		assert.Equal(t, "/statements@10/expression/args@1", mustPathOf(fc))
		assert.Equal(t, "vat", fc.Callable().Text())
		assert.Equal(t, "/statements@8", mustPathOf(fc.Callable().Ref(doc.Ctx())))
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

		path, err := core.PathOf(&expr)
		assert.Nil(t, path)
		fmt.Println(err)
		assert.ErrorContains(
			t, err,
			"PathOf: error within node of type *arithmetics.ModuleImpl: cannot determine node path, 'containerField' is empty")
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

		path, err := core.PathOf(&expr)
		assert.Nil(t, path)
		fmt.Println(err)
		assert.ErrorContains(
			t, err,
			"PathOf: error within node of type *arithmetics.DefinitionImpl: cannot determine node path, 'containerField' is empty")
	})
}

// TestResolve_PriceCalc is the inverse of TestNodePath_PriceCalc: for each of the 13
// FunctionCall cross-reference sites in price-calc.calc, it resolves the hardcoded path
// string back to an AST node via Resolve and asserts pointer identity with the node
// obtained by direct AST navigation. Both the referencer (FunctionCall) and the
// referenced (AbstractDefinition) are covered per subtest.
func TestResolve_PriceCalc(t *testing.T) {
	doc := loadPriceCalcDoc(t)
	root := doc.Root()
	module := test.MustFindNode[Module](doc)
	stmts := module.Statements()
	require.Len(t, stmts, 11)

	mustResolve := func(path string, node core.AstNode) core.AstNode {
		t.Helper()
		child, err := core.Resolve(path, node)
		require.NoError(t, err)
		return child
	}
	// ── def costPerUnit: materialPerUnit + laborPerUnit ───────────────────────

	t.Run("costPerUnit/materialPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[3])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Left())

		left := mustResolve("/statements@3/expression/left", root)
		assert.Same(t, fc, left)

		statement := mustResolve("/statements@0", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("costPerUnit/laborPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[3])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Right())

		right := mustResolve("/statements@3/expression/right", root)
		assert.Same(t, fc, right)

		statement := mustResolve("/statements@1", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	// ── def costOfGoodsSold: expectedNoOfSales * costPerUnit ─────────────────

	t.Run("costOfGoodsSold/expectedNoOfSales", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[4])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Left())

		left := mustResolve("/statements@4/expression/left", root)
		assert.Same(t, fc, left)

		statement := mustResolve("/statements@2", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("costOfGoodsSold/costPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[4])
		binExpr := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, binExpr.Right())

		right := mustResolve("/statements@4/expression/right", root)
		assert.Same(t, fc, right)

		statement := mustResolve("/statements@3", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
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
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		div := mustConvert[BinaryExpression](t, outerAdd.Left())
		innerAdd := mustConvert[BinaryExpression](t, div.Left())
		fc := mustConvert[FunctionCall](t, innerAdd.Left())

		left := mustResolve("/statements@7/expression/left/left/left", root)
		assert.Same(t, fc, left)

		statement := mustResolve("/statements@4", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("netPrice/generalExpensesAndSales", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		div := mustConvert[BinaryExpression](t, outerAdd.Left())
		innerAdd := mustConvert[BinaryExpression](t, div.Left())
		fc := mustConvert[FunctionCall](t, innerAdd.Right())

		right := mustResolve("/statements@7/expression/left/left/right", root)
		assert.Same(t, fc, right)

		statement := mustResolve("/statements@5", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("netPrice/expectedNoOfSales", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		div := mustConvert[BinaryExpression](t, outerAdd.Left())
		fc := mustConvert[FunctionCall](t, div.Right())

		right := mustResolve("/statements@7/expression/left/right", root)
		assert.Same(t, fc, right)

		statement := mustResolve("/statements@2", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("netPrice/desiredProfitPerUnit", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[7])
		outerAdd := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, outerAdd.Right())

		right := mustResolve("/statements@7/expression/right", root)
		assert.Same(t, fc, right)

		statement := mustResolve("/statements@6", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
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
		def := mustConvert[Definition](t, stmts[9])
		div := mustConvert[BinaryExpression](t, def.Expression())
		fc := mustConvert[FunctionCall](t, div.Left())

		left := mustResolve("/statements@9/expression/left", root)
		assert.Same(t, fc, left)

		arg := mustResolve("/statements@9/args@0", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), arg)
	})

	t.Run("calcGrossListPrice/tax", func(t *testing.T) {
		def := mustConvert[Definition](t, stmts[9])
		div := mustConvert[BinaryExpression](t, def.Expression())
		sub := mustConvert[BinaryExpression](t, div.Right())
		fc := mustConvert[FunctionCall](t, sub.Right())

		right := mustResolve("/statements@9/expression/right/right", root)
		assert.Same(t, fc, right)

		arg := mustResolve("/statements@9/args@1", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), arg)
	})

	// ── calcGrossListPrice(netPrice, vat) ────────────────────────────────────

	t.Run("evaluation/calcGrossListPrice", func(t *testing.T) {
		eval := mustConvert[Evaluation](t, stmts[10])
		fc := mustConvert[FunctionCall](t, eval.Expression())

		expression := mustResolve("/statements@10/expression", root)
		assert.Same(t, fc, expression)

		statement := mustResolve("/statements@9", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("evaluation/netPrice", func(t *testing.T) {
		eval := mustConvert[Evaluation](t, stmts[10])
		outerFC := mustConvert[FunctionCall](t, eval.Expression())
		require.Len(t, outerFC.Args(), 2)
		fc := mustConvert[FunctionCall](t, outerFC.Args()[0])

		arg := mustResolve("/statements@10/expression/args@0", root)
		assert.Same(t, fc, arg)

		statement := mustResolve("/statements@7", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("evaluation/vat", func(t *testing.T) {
		eval := mustConvert[Evaluation](t, stmts[10])
		outerFC := mustConvert[FunctionCall](t, eval.Expression())
		require.Len(t, outerFC.Args(), 2)
		fc := mustConvert[FunctionCall](t, outerFC.Args()[1])

		arg := mustResolve("/statements@10/expression/args@1", root)
		assert.Same(t, fc, arg)

		statement := mustResolve("/statements@8", root)
		assert.Same(t, fc.Callable().Ref(doc.Ctx()), statement)
	})

	t.Run("errorReporting/no-such-field", func(t *testing.T) {
		_, err := core.Resolve("/statements@3/expressions/left", root)
		assert.ErrorContains(t, err, "DefinitionImpl.Resolve: field 'expressions' does not exist in node '/statements@3' of type 'Definition'")
	})

	t.Run("errorReporting/field-is-primitive", func(t *testing.T) {
		_, err := core.Resolve("/statements@0/name", root)
		assert.ErrorContains(t, err, "DefinitionImpl.Resolve: field 'name' holds a primitive value instead of an ast node")
	})

	t.Run("errorReporting/field-is-reference", func(t *testing.T) {
		_, err := core.Resolve("/statements@10/expression/callable", root)
		assert.ErrorContains(t, err, "FunctionCallImpl.Resolve: field 'callable' is a cross-reference instead of a container field")
	})

	t.Run("errorReporting/field-is-empty", func(t *testing.T) {
		// create local copies of the relevant ast nodes first to avoid manipulating the shared ones
		module := *module.(*ModuleImpl)
		def := *stmts[0].(*DefinitionImpl)

		module.statements = []Statement{&def}

		field, index := def.ContainmentData()
		def.SetContainer(&module, field, index)

		// set the tested field to 'nil'
		def.expression = nil

		_, err := core.Resolve("/statements@0/expression/left", &module)
		assert.ErrorContains(t, err, "DefinitionImpl.Resolve: field 'expression' is nil in node '/statements@0'")
	})

	t.Run("errorReporting/slice-index-out-of-bound-1", func(t *testing.T) {
		_, err := core.Resolve("/statements@15/expression", root)
		assert.ErrorContains(t, err, "ModuleImpl.Resolve: index 15 exceeds length of slice in 'statements' (length=11) in node ''")
	})

	t.Run("errorReporting/slice-index-out-of-bound-2", func(t *testing.T) {
		_, err := core.Resolve("/statements@10/expression/args@7/expression", root)
		assert.ErrorContains(t, err, "FunctionCallImpl.Resolve: index 7 exceeds length of slice in 'args' (length=2) in node '/statements@10/expression'")
	})

	t.Run("errorReporting/slice-item-is-nil", func(t *testing.T) {
		expr, err := core.Resolve("/statements@10/expression/", root)
		require.NoError(t, err)
		fc := mustConvert[FunctionCall](t, expr)

		// shamelessly manipulate the shared ast and add a 'nil' item, don't want to copy everything right now; shouldn't hurt
		fc.SetArgsItem(nil)
		_, err = core.Resolve("/statements@10/expression/args@2/expression", root)

		assert.ErrorContains(t, err, "FunctionCallImpl.Resolve: item 2 of slice in field 'args' is nil in node '/statements@10/expression'")
	})

	t.Run("errorReporting/slice-index-typo", func(t *testing.T) {
		_, err := core.Resolve("/statements@1a/expression", root)
		assert.ErrorContains(t, err, "parsePath: index '1a' is not a valid uint: strconv.Atoi: parsing \"1a\": invalid syntax")
	})
}
