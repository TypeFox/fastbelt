// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/extiter"
)

func TestDefinitionsAndCallsLinking(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`module test
		def two: 2;
		def root(x, y): x^(1/y);
		def sqrt(x): root(x, two);
		sqrt(16);
	`)
	doc.AssertState(core.DocStateLinked)
	doc.AssertNoParseErrors()
	doc.AssertNoLinkingErrors()

	require.Len(t, test.MustFindNode[Module](doc).Statements(), 4)

	defTwo := test.MustFindNamedNode[Definition](doc, "two")
	defRoot := test.MustFindNamedNode[Definition](doc, "root")
	defSqrt := test.MustFindNamedNode[Definition](doc, "sqrt")
	eval := test.MustFindNode[Evaluation](doc)

	evalCall, isCall := eval.Expression().(FunctionCall)
	require.True(t, isCall)

	assert.Same(
		t,
		defSqrt,
		evalCall.Callable().Ref(doc.Ctx()),
	)

	callsInSqrt := slices.Collect(
		extiter.FilterType[core.AstNode, FunctionCall](
			core.AllChildren(defSqrt),
		),
	)
	require.Len(t, callsInSqrt, 3)

	assert.Same(
		t,
		defRoot,
		callsInSqrt[0].Callable().Ref(doc.Ctx()),
	)

	assert.Same(
		t,
		defSqrt.Args()[0],
		callsInSqrt[1].Callable().Ref(doc.Ctx()),
	)

	assert.Same(
		t,
		defTwo,
		callsInSqrt[2].Callable().Ref(doc.Ctx()),
	)
}
