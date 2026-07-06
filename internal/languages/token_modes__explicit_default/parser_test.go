// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package token_modes__explicit_default

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

func TestDeclaration(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := 123")
	doc.AssertNoErrors()
}

func TestSimpleString(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := <|string:`simple string`|>")
	doc.AssertNoErrors()
	node, ok := doc.FindAstNode("string")
	require.True(t, ok)
	str, ok := node.(StringLiteral)
	require.True(t, ok)
	content := str.Content()
	require.Len(t, content, 1)
	require.Equal(t, "simple string", content[0].(StringText).Value())
}

func TestEscapedString(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := <|string:`escaped\\nstring`|>")
	doc.AssertNoErrors()
	node, ok := doc.FindAstNode("string")
	require.True(t, ok)
	str, ok := node.(StringLiteral)
	require.True(t, ok)
	content := str.Content()
	require.Len(t, content, 1)
	require.Equal(t, "escaped\\nstring", content[0].(StringText).Value())
}

func TestInterpolatingString(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("NUM := 123\nVAR := <|string:`nested #{NUM}`|>")
	doc.AssertNoErrors()
	node, ok := doc.FindAstNode("string")
	require.True(t, ok)
	interpolation, ok := node.(StringLiteral)
	require.True(t, ok)
	content := interpolation.Content()
	require.Len(t, content, 2)
	content0, ok := content[0].(StringText)
	require.True(t, ok)
	require.Equal(t, "nested ", content0.Value())
	content1, ok := content[1].(Interpolation)
	require.True(t, ok)
	require.Equal(t, "NUM", content1.Expression().(VariableRef).Name().Ref(context.Background()).Name())

}

func TestNestedString(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("NUM := 123\nVAR := <|string:`this is #{`a number: #{NUM}`}`|>")
	doc.AssertNoErrors()
	node, ok := doc.FindAstNode("string")
	require.True(t, ok)
	interpolation, ok := node.(StringLiteral)
	require.True(t, ok)
	content := interpolation.Content()
	require.Len(t, content, 2)
	content0, ok := content[0].(StringText)
	require.True(t, ok)
	require.Equal(t, "this is ", content0.Value())
	content1, ok := content[1].(Interpolation)
	require.True(t, ok)
	content2, ok := content1.Expression().(StringLiteral)
	require.True(t, ok)
	require.Len(t, content2.Content(), 2)
	content20, ok := content2.Content()[0].(StringText)
	require.True(t, ok)
	require.Equal(t, "a number: ", content20.Value())
	content21, ok := content2.Content()[1].(Interpolation)
	require.True(t, ok)
	require.Equal(t, "NUM", content21.Expression().(VariableRef).Name().Ref(context.Background()).Name())
}
