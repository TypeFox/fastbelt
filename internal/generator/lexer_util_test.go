// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

func TestGetAllKeywords_InRule(t *testing.T) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Model {
			Greeting string
		}
		Model: Greeting="hello" "world" | Greeting="hi" "there"
	`).AssertNoErrors()
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	result := GetAllKeywords(grammr)
	require.Len(t, result.Keywords, 4)
	assert.Equal(t, "hello", grammar.KeywordValue(result.Keywords[0]))
	assert.Equal(t, "hi", grammar.KeywordValue(result.Keywords[1]))
	assert.Equal(t, "there", grammar.KeywordValue(result.Keywords[2]))
	assert.Equal(t, "world", grammar.KeywordValue(result.Keywords[3]))
}

func TestGetAllKeywords_InTokenDeclaration(t *testing.T) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Model {
			Greeting string
		}
		Model: Greeting=HELLO "world"
		token HELLO: "hello"
	`).AssertNoErrors()
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	result := GetAllKeywords(grammr)
	require.Len(t, result.Keywords, 2)
	assert.Equal(t, "hello", grammar.KeywordValue(result.Keywords[0]))
	assert.Equal(t, "world", grammar.KeywordValue(result.Keywords[1]))
}

func TestGetAllKeywords_InTokenGroup(t *testing.T) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Model {
			Greeting string
		}
		Model: Greeting=GREETING "world"
		token group GREETING {
			"hello"
			"hi"
		}
	`).AssertNoErrors()
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	result := GetAllKeywords(grammr)
	require.Len(t, result.Keywords, 3)
	assert.Equal(t, "hello", grammar.KeywordValue(result.Keywords[0]))
	assert.Equal(t, "hi", grammar.KeywordValue(result.Keywords[1]))
	assert.Equal(t, "world", grammar.KeywordValue(result.Keywords[2]))
}

func TestGetAllKeywords_InTokenMode(t *testing.T) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Model {
			Greeting string
		}
		Model: Greeting="hello" "world"
		token mode default {
			"hello"
			"world"
		}
	`).AssertNoErrors()
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	result := GetAllKeywords(grammr)
	require.Len(t, result.Keywords, 2)
	assert.Equal(t, "hello", grammar.KeywordValue(result.Keywords[0]))
	assert.Equal(t, "world", grammar.KeywordValue(result.Keywords[1]))
}
