// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lookahead

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

func expectValue(t *testing.T, doc *test.Doc, expected string) Obj {
	root := doc.Document.Root.(Root)
	obj := root.Item()
	require.NotNil(t, obj)
	var actual string
	if v := obj.Value(); v != "" {
		actual = v
	} else if node := obj.Node(); node != "" {
		actual = node
	} else {
		t.Fatalf("expected value or node, got empty")
	}
	assert.Equal(t, expected, actual)
	return obj
}

func TestProperQualifiedPathLookahead(t *testing.T) {
	values := []string{
		"myId",
		"some/path:node",
		"pkg::node",
		"my.value.dot::node",
	}
	for _, v := range values {
		t.Run(v, func(t *testing.T) {
			sc := CreateServices()
			doc := test.New(t, sc).Parse("a " + v)
			doc.AssertNoParseErrors()
			expectValue(t, doc, v)
		})
	}
}

func TestPostIdAfterQualifiedPathLookahead(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("b my.value.dot::node.postId")
	doc.AssertNoParseErrors()
	obj := expectValue(t, doc, "my.value.dot::node")
	b := obj.(B)
	assert.Equal(t, "postId", b.Post())
}

func TestPostIdAfterIdLookahead(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("b first.second")
	doc.AssertNoParseErrors()
	obj := expectValue(t, doc, "first")
	b := obj.(B)
	assert.Equal(t, "second", b.Post())
}

func TestCLoopOptionalLookahead(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("c A::B::**")
	doc.AssertNoParseErrors()
	root := doc.Document.Root.(Root)
	obj := root.Item()
	require.NotNil(t, obj)
	assert.Equal(t, "A::B", obj.Node())
	assert.Equal(t, "**", obj.Value())
}

func TestCLoopOptionalLookaheadNegative(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("c A::B::C")
	doc.AssertNoParseErrors()
	root := doc.Document.Root.(Root)
	obj := root.Item()
	require.NotNil(t, obj)
	assert.Equal(t, "A::B::C", obj.Node())
	assert.Equal(t, "", obj.Value())
}

func TestDOptOptionalLookaheadPositive(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("d A::**")
	doc.AssertNoParseErrors()
	root := doc.Document.Root.(Root)
	obj := root.Item()
	require.NotNil(t, obj)
	assert.Equal(t, "A", obj.Node())
	assert.Equal(t, "**", obj.Value())
}

func TestDOptOptionalLookaheadNegative(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("d A::B::**")
	doc.AssertNoParseErrors()
	root := doc.Document.Root.(Root)
	obj := root.Item()
	require.NotNil(t, obj)
	assert.Equal(t, "A::B", obj.Node())
	assert.Equal(t, "**", obj.Value())
}

func TestEOptionalValueLookahead(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("e hello")
	doc.AssertNoParseErrors()
	expectValue(t, doc, "hello")
}

func TestEMissingValueLookahead(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("e world")
	doc.AssertNoParseErrors()
}

func TestEBothValueLookahead(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("e hello world")
	doc.AssertNoParseErrors()
	expectValue(t, doc, "hello")
}

func TestFUnlimitedLookaheadPositiveHello(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("f A B C D E F G hello")
	doc.AssertNoParseErrors()
	expectValue(t, doc, "hello")
}

func TestFUnlimitedLookaheadPositiveWorld(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("f A B C D E F G world")
	doc.AssertNoParseErrors()
	expectValue(t, doc, "world")
}

func TestFUnlimitedLookaheadNegative(t *testing.T) {
	sc := CreateServices()
	doc := test.New(t, sc).Parse("f A B C D E F G")
	assert.Greater(t, len(doc.Document.ParserErrors), 0, "expected parse errors")
}
