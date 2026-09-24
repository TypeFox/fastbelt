// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
)

// writeGrammarDir writes the given (name, source) pairs into a fresh temp
// directory and returns its path.
func writeGrammarDir(t *testing.T, nameSourcePairs ...string) string {
	t.Helper()
	dir := t.TempDir()
	for i := 0; i < len(nameSourcePairs); i += 2 {
		path := filepath.Join(dir, nameSourcePairs[i])
		require.NoError(t, os.WriteFile(path, []byte(nameSourcePairs[i+1]), 0644))
	}
	return dir
}

const multiTokens = `
	token ID: /[_a-zA-Z][\w_]*/
	hidden token WS: /\s+/
`

// A rule and the interface that declares its return type conventionally share
// a name, and the interface may live in another file.
func TestLoadGrammarDirMergesFiles(t *testing.T) {
	dir := writeGrammarDir(t,
		"a.fb", `
			grammar Multi
			interface Greeting { Name string }
			interface Farewell { Name string }
			entry Greeting: "hello" Name=ID
		`+multiTokens,
		"b.fb", `
			grammar Multi
			Farewell: "goodbye" Name=ID
		`,
	)

	g, err := LoadGrammarDir(dir)
	require.NoError(t, err)
	require.Len(t, g.Rules(), 2)
	require.Len(t, g.Interfaces(), 2)
	require.Len(t, g.Terminals(), 2)
	// The implicit return type of Farewell is found in the merged grammar.
	farewell := g.Rules()[1]
	require.Equal(t, "Farewell", farewell.Name())
	require.NotNil(t, grammar.FindReturnType(farewell, t.Context()))
}

func TestLoadGrammarDirKeepsTokenModesAndInfixRules(t *testing.T) {
	dir := writeGrammarDir(t,
		"a.fb", `
			grammar Multi
			interface Expr {}
			interface Primary extends Expr { Value string }
			interface Binary extends Expr { Left Expr Operator string Right Expr }
			entry Expr returns Expr: Binary
			Primary: Value=ID
			token ID: /[_a-zA-Z][\w_]*/
			hidden token WS: /\s+/
			token mode default {
				ID
				hidden WS
				"+" -> push(Inner)
			}
		`,
		"b.fb", `
			grammar Multi
			infix Binary on Primary returns Expr: "+"
			token mode Inner {
				ID -> pop
			}
		`,
	)

	g, err := LoadGrammarDir(dir)
	require.NoError(t, err)
	require.Len(t, g.TokenModes(), 2)
	require.Len(t, g.InfixRules(), 1)
}

func TestLoadGrammarDirRejectsDuplicateRuleAcrossFiles(t *testing.T) {
	dir := writeGrammarDir(t,
		"a.fb", `
			grammar Multi
			interface Greeting { Name string }
			entry Greeting: "hello" Name=ID
		`+multiTokens,
		"b.fb", `
			grammar Multi
			Greeting: "hi" Name=ID
		`,
	)

	_, err := LoadGrammarDir(dir)
	require.ErrorContains(t, err, "aborting code generation due to 2 errors")
}

func TestLoadGrammarDirRejectsMismatchedGrammarNames(t *testing.T) {
	dir := writeGrammarDir(t,
		"a.fb", `
			grammar Multi
			interface Greeting { Name string }
			entry Greeting: "hello" Name=ID
		`+multiTokens,
		"b.fb", `
			grammar Other
			interface Farewell { Name string }
			Farewell: "goodbye" Name=ID
		`,
	)

	_, err := LoadGrammarDir(dir)
	require.ErrorContains(t, err, "aborting code generation due to 2 errors")
}
