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
	assert.Contains(t, result.ByValue, "\"hello\"")
	assert.Contains(t, result.ByValue, "\"hi\"")
	assert.Contains(t, result.ByValue, "\"there\"")
	assert.Contains(t, result.ByValue, "\"world\"")
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
	assert.Contains(t, result.ByValue, "\"hello\"")
	assert.Contains(t, result.ByValue, "\"world\"")
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
	assert.Contains(t, result.ByValue, "\"hello\"")
	assert.Contains(t, result.ByValue, "\"hi\"")
	assert.Contains(t, result.ByValue, "\"world\"")
}

func TestGetAllTokenDecls_SplitsTopLevelAndModeLevel(t *testing.T) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Model {
			Greeting string
		}
		Model: Greeting=ID INNER
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID
			hidden WS
		}
		token mode Inner {
			token INNER: /[A-Z]+/
		}
	`)
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	result := GetAllTokenDecls(grammr)
	assert.Equal(t, []string{"ID", "WS"}, tokenDeclNames(result.TopLevel))
	assert.Equal(t, []string{"INNER"}, tokenDeclNames(result.ModeLevel))
	// All lists top-level declarations first, then the mode-local ones.
	assert.Equal(t, []string{"ID", "WS", "INNER"}, tokenDeclNames(result.All))
}

func tokenDeclNames(decls []grammar.TokenDecl) []string {
	names := make([]string, len(decls))
	for i, decl := range decls {
		names[i] = decl.Name()
	}
	return names
}

func TestGetAllTokenGroups_SplitsTopLevelAndModeLevel(t *testing.T) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Model {
			Greeting string
		}
		Model: Greeting=Outer Closers
		token group Outer { ID }
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			Outer
			hidden WS
		}
		token mode Inner {
			token group Closers { ")" }
		}
	`)
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	result := GetAllTokenGroups(grammr)
	assert.Equal(t, []string{"Outer"}, tokenGroupNames(result.TopLevel))
	assert.Equal(t, []string{"Closers"}, tokenGroupNames(result.ModeLevel))
	assert.Equal(t, []string{"Outer", "Closers"}, tokenGroupNames(result.All))
}

func tokenGroupNames(groups []grammar.TokenGroup) []string {
	names := make([]string, len(groups))
	for i, group := range groups {
		names[i] = group.Name()
	}
	return names
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
	assert.Contains(t, result.ByValue, "\"hello\"")
	assert.Contains(t, result.ByValue, "\"world\"")
}
