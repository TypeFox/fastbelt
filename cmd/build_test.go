// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package cmd

import (
	"os"
	"path/filepath"
	"strings"
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
	// Rules are sorted by name regardless of file order, and the implicit
	// return type of Farewell is found in the merged grammar.
	farewell := g.Rules()[0]
	require.Equal(t, "Farewell", farewell.Name())
	require.Equal(t, "Greeting", g.Rules()[1].Name())
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

func TestLoadGrammarDirRejectsTwoDefaultTokenModes(t *testing.T) {
	dir := writeGrammarDir(t,
		"a.fb", `
			grammar Multi
			interface Greeting { Name string }
			entry Greeting: "hello" Name=ID
			token mode default {
				"hello"
				ID
				hidden WS
			}
		`+multiTokens,
		"b.fb", `
			grammar Multi
			token mode default {
				ID
			}
		`,
	)

	_, err := LoadGrammarDir(dir)
	require.ErrorContains(t, err, "aborting code generation due to")
}

const modeGrammar = `
	grammar Multi
	interface Greeting { Name string }
	interface Farewell { Name string }
	entry Greeting: "hello" Name=ID
	entry Farewell: "goodbye" Name=ID
	token ID: /[_a-zA-Z][\w_]*/
	hidden token WS: /\s+/
	token mode default {
		"hello" -> push(Inner)
		ID
		hidden WS
	}
	token mode Inner {
		"goodbye"
		ID -> pop
		hidden WS
	}
`

func TestBuildRejectsUnknownTokenMode(t *testing.T) {
	dir := writeGrammarDir(t, "a.fb", modeGrammar)
	ctx := &BuildContext{Input: dir, Output: t.TempDir(), Package: "multi", Languages: []Language{
		{Entry: "Farewell", TokenMode: "Nope"},
	}}
	require.ErrorContains(t, ctx.Build(), `token mode "Nope" not found`)
}

func TestBuildUsesConfiguredStartTokenMode(t *testing.T) {
	dir := writeGrammarDir(t, "a.fb", modeGrammar)
	out := t.TempDir()
	ctx := &BuildContext{Input: dir, Output: out, Package: "multi", Languages: []Language{
		{Entry: "Greeting", LanguageID: "greeting", Patterns: []string{"**/*.hello"}},
		{Entry: "Farewell", TokenMode: "Inner", LanguageID: "farewell", Patterns: []string{"**/*.bye"}},
	}}
	require.NoError(t, ctx.Build())
	code, err := os.ReadFile(filepath.Join(out, "lexer_gen.go"))
	require.NoError(t, err)
	require.Contains(t, string(code), "lexer.NewMultiLanguageLexer(sc, []int{TokenMode_default, TokenMode_Inner}, modes...)")
}

func TestBuildRequiresStartModeWithoutDefaultTokenMode(t *testing.T) {
	dir := writeGrammarDir(t, "a.fb", `
		grammar Multi
		interface Greeting { Name string }
		interface Farewell { Name string }
		entry Greeting: "hello" Name=ID
		entry Farewell: "goodbye" Name=ID
		token ID: /[_a-zA-Z][\w_]*/
		hidden token WS: /\s+/
		token mode GreetingMode {
			"hello"
			ID
			hidden WS
		}
		token mode FarewellMode {
			"goodbye"
			ID
			hidden WS
		}
	`)
	ctx := &BuildContext{Input: dir, Output: t.TempDir(), Package: "multi", Languages: []Language{
		{Entry: "Greeting", TokenMode: "GreetingMode", LanguageID: "greeting"},
		{Entry: "Farewell", LanguageID: "farewell"},
	}}
	require.ErrorContains(t, ctx.Build(), "declares no default token mode")
}

// Operators of an infix rule are only reachable through that rule, so a
// multi-language lexer must keep them in the modes of the language that calls
// the rule.
func TestBuildMultiLanguageLexerKeepsInfixOperators(t *testing.T) {
	dir := writeGrammarDir(t, "a.fb", `
		grammar Multi
		interface Greeting { Name string }
		interface Calc { Value Expression }
		interface Expression {}
		interface BinaryExpression extends Expression { Left Expression Operator string Right Expression }
		interface NumberLiteral extends Expression { Value string }
		entry Greeting: "hello" Name=ID
		entry Calc: "calc" Value=Expression
		Expression returns Expression: BinaryExpression
		infix BinaryExpression on PrimaryExpression:
			"+" | "-"
		PrimaryExpression returns Expression: NumberLiteral
		NumberLiteral: Value=NUMBER
		token NUMBER: /[0-9]+/
	`+multiTokens)
	out := t.TempDir()
	ctx := &BuildContext{Input: dir, Output: out, Package: "multi", Languages: []Language{
		{Entry: "Greeting", LanguageID: "greeting", Patterns: []string{"**/*.hello"}},
		{Entry: "Calc", LanguageID: "calc", Patterns: []string{"**/*.calc"}},
	}}
	require.NoError(t, ctx.Build())
	code, err := os.ReadFile(filepath.Join(out, "lexer_gen.go"))
	require.NoError(t, err)
	// The grammar declares no token modes, so each language gets a synthetic
	// one. Everything from the Calc mode on belongs to the second language.
	greetingMode, calcMode, found := strings.Cut(string(code), "modes[TokenMode_Calc] =")
	require.True(t, found)
	_, greetingMode, found = strings.Cut(greetingMode, "modes[TokenMode_Greeting] =")
	require.True(t, found)
	require.Contains(t, calcMode, "Keyword_Plus")
	require.Contains(t, calcMode, "Keyword_Dash")
	require.Contains(t, calcMode, "Token_NUMBER")
	require.NotContains(t, greetingMode, "Keyword_Plus")
	require.NotContains(t, greetingMode, "Token_NUMBER")
}
