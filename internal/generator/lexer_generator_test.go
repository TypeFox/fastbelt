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

func TestGenerateLexerPushDefaultModeCommand(t *testing.T) {
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
			ID -> push(default)
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID).WithPushMode(TokenMode_default)")
}

// --- Mode ordering and identifiers ---

func TestGenerateLexerModeIdsFollowDeclarationOrder(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode First {
			ID -> push(default)
		}
		token mode default {
			ID -> push(First)
			hidden WS
		}
	`)
	// The default mode is declared second, so it gets id 1 and the lexer must
	// still start in it.
	assert.Contains(t, code, "TokenMode_First   = 0")
	assert.Contains(t, code, "TokenMode_default = 1")
	assert.Contains(t, code, "modes := make([]*lexer.TokenMode, 2, 2)")
	assert.Contains(t, code, "lexer.NewDefaultLexer(TokenMode_default, modes...)")
}

func TestGenerateLexerImplicitDefaultModeRegistersEverything(t *testing.T) {
	// Without any token mode declaration the generator synthesizes a default
	// mode from all top-level keywords, tokens and token groups.
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world" Punct;
		token group Punct { "!" }
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		comment token SL_COMMENT: /\/\/[^\n]*/
	`)
	assert.Contains(t, code, "TokenMode_default = 0")
	assert.Contains(t, code, "lexer.UseTokenType(Keyword_world)")
	assert.Contains(t, code, "lexer.UseTokenType(Token_ID)")
	assert.Contains(t, code, "lexer.UseTokenType(TokenGroup_Punct)")
	assert.Contains(t, code, "lexer.UseTokenType(Token_WS).WithGroup(core.SkippedGroup)")
	assert.Contains(t, code, "lexer.UseTokenType(Token_SL_COMMENT).WithGroup(core.CommentGroup)")
	assert.Contains(t, code, "lexer.NewDefaultLexer(TokenMode_default, modes...)")
}

// --- Modifier and command precedence ---

func TestGenerateLexerModeUsageInheritsGlobalModifier(t *testing.T) {
	// WS is hidden globally and listed without a modifier in the mode, so the
	// global modifier carries over - `hidden` does not have to be repeated.
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID
			WS
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_WS).WithGroup(core.SkippedGroup)")
}

func TestGenerateLexerModeUsageOverridesGlobalModifier(t *testing.T) {
	// A modifier on the mode entry replaces the global one.
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token NOTE: /\/\/[^\n]*/
		token mode default {
			ID
			comment NOTE
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_NOTE).WithGroup(core.CommentGroup)")
	assert.NotContains(t, code, "Token_NOTE).WithGroup(core.SkippedGroup)")
}

func TestGenerateLexerModeUsageAddsModifierToGlobalToken(t *testing.T) {
	// The reverse direction: a plain global token is hidden inside the mode.
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		token WS: /\s+/
		token mode default {
			ID
			hidden WS
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_WS).WithGroup(core.SkippedGroup)")
}

func TestGenerateLexerModeUsageInheritsGlobalCommentModifier(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		comment token SL_COMMENT: /\/\/[^\n]*/
		token mode default {
			ID
			SL_COMMENT
		}
	`)
	assert.Contains(t, code, "lexer.UseTokenType(Token_SL_COMMENT).WithGroup(core.CommentGroup)")
}

// --- Mode-local declarations ---

func TestGenerateLexerModeLocalTokenDeclaration(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string Content string }
		Foo: Greeting=ID "(" Content=INNER ")";
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID
			"(" -> push(Inner)
			hidden WS
		}
		token mode Inner {
			token INNER: /[A-Z]+/
			")" -> pop
		}
	`)
	// The mode-local declaration produces its own token type and is registered
	// only in the mode that declares it.
	assert.Contains(t, code, "var Token_INNER =")
	inner := code[strings.Index(code, `NewTokenMode("Inner"`):]
	assert.Contains(t, inner, "lexer.UseTokenType(Token_INNER)")
	outer := code[strings.Index(code, `NewTokenMode("default"`):strings.Index(code, `NewTokenMode("Inner"`)]
	assert.NotContains(t, outer, "Token_INNER")
}

// A mode-local declaration is visible to parser rules but not to a token usage
// in a different mode, which the linker reports. The generated lexer must still
// be well-formed - the unresolvable entry is dropped rather than emitted broken.
func TestGenerateLexerModeLocalTokenIsNotSharedAcrossModes(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string Content string }
		Foo: Greeting=ID "(" Content=INNER ")";
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID
			"(" -> push(Inner)
			hidden WS
		}
		token mode Inner {
			token INNER: /[A-Z]+/
			")" -> pop
		}
		token mode Other {
			INNER -> pop
		}
	`)
	assert.Equal(t, 1, strings.Count(code, "var Token_INNER ="))
	other := code[strings.Index(code, `NewTokenMode("Other"`):]
	assert.NotContains(t, other, "Token_INNER")
}

func TestGenerateLexerModeLocalTokenGroupWithCommand(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID Closers;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID
			"(" -> push(Inner)
			hidden WS
		}
		token mode Inner {
			token group Closers {
				")"
				"]"
			} -> pop
		}
	`)
	assert.Contains(t, code, "var TokenGroup_Closers =")
	inner := code[strings.Index(code, `NewTokenMode("Inner"`):]
	assert.Contains(t, inner, "lexer.UseTokenType(TokenGroup_Closers).WithPopMode()")
}

func TestGenerateLexerTokenListedTwiceInSameMode(t *testing.T) {
	code := generateLexerFor(t, `
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;
		token ID: /[a-z]+/
		hidden token WS: /\s+/
		token mode default {
			ID
			ID
			hidden WS
		}
	`)
	// Listing a token twice registers it twice. Harmless at runtime - the
	// second candidate can never win - but worth pinning down.
	assert.Equal(t, 2, strings.Count(code, "lexer.UseTokenType(Token_ID)"))
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
