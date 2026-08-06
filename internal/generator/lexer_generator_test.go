// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

// generateLexerFor generates the lexer for src without asserting that src is
// free of diagnostics, so malformed grammars can be exercised too.
func generateLexerFor(t *testing.T, src string) string {
	t.Helper()
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(src)
	grammr, ok := doc.Document.Root.(grammar.Grammar)
	require.True(t, ok)
	return GenerateLexer(grammr, "test", GenerateTokenTypes(grammr))
}

func TestGenerateLexerModeCommands(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID -> push(Other)
			hidden WS
		}
		token mode Other {
			ID -> pop
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID).WithPushMode(TokenMode_Other)")
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID).WithPopMode()")
	assert.Contains(t, code, "lexer.NewDefaultLexer(TokenMode_default, modes...)")
}

func TestGenerateLexerSetModeCommand(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID -> push(Other)
			hidden WS
		}
		token mode Other {
			ID -> mode(default)
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID).WithSetMode(TokenMode_default)")
}

// A push command without a target mode is a validation error. Code generation
// must still terminate instead of panicking, because GenerateLexer is reachable
// for grammars that have not been validated.
func TestGenerateLexerPushCommandWithoutTargetMode(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID -> push
			hidden WS
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID)")
	assert.NotContains(t, code, "WithPushMode")
}

func TestGenerateLexerPushCommandWithUnknownTargetMode(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID -> push(Missing)
			hidden WS
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID)")
	assert.NotContains(t, code, "WithPushMode")
}

// A pop command keeps working even when a target mode was given, which is what
// the corresponding validation error tells the user.
func TestGenerateLexerPopCommandIgnoresTargetMode(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID -> push(Other)
			hidden WS
		}
		token mode Other {
			ID -> pop(default)
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID).WithPopMode()")
	assert.NotContains(t, code, "WithSetMode")
}

// Regression test for the keyword selector dropping its bookkeeping between
// members: a keyword that is both listed explicitly and matched by a selector
// must be registered once.
func TestGenerateLexerKeywordSelectorDoesNotDuplicateKeywords(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting="hello" "world";
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			"hello"
			keywords /^[a-z]+$/
			ID
			hidden WS
		}
	`)
	assert.Equal(t, 1, strings.Count(code, "lexer.UseTokenType(Keyword_hello)"),
		"keyword registered more than once:\n%s", code)
	assert.Equal(t, 1, strings.Count(code, "lexer.UseTokenType(Keyword_world)"))
}
