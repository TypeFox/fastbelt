// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package token_modes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

// stringOf parses src and returns the string literal marked with the "string" label.
func stringOf(t *testing.T, src string) StringLiteral {
	t.Helper()
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse(src)
	doc.AssertNoErrors()
	return test.MustFindNodeAtLabel[StringLiteral](doc, "string")
}

// --- String content shapes ---

func TestEmptyString(t *testing.T) {
	str := stringOf(t, "VAR := <|string:``|>")
	assert.Empty(t, str.Content())
}

func TestStringWithOnlyInterpolation(t *testing.T) {
	str := stringOf(t, "NUM := 1\nVAR := <|string:`#{NUM}`|>")
	content := str.Content()
	require.Len(t, content, 1)
	interpolation, ok := content[0].(Interpolation)
	require.True(t, ok)
	assert.Equal(t, "NUM", interpolation.Expression().(VariableRef).Name().Ref(context.Background()).Name())
}

func TestStringPreservesWhitespace(t *testing.T) {
	// WS is skipped in the default mode but is not registered in IN_STRING, so
	// whitespace inside string text must survive as part of STRING_CONTENT.
	str := stringOf(t, "VAR := <|string:`  spaced   out  `|>")
	content := str.Content()
	require.Len(t, content, 1)
	assert.Equal(t, "  spaced   out  ", content[0].(StringText).Value())
}

func TestStringPreservesCommentSyntax(t *testing.T) {
	// Comments are only recognized in the default and interpolation modes, so
	// `//` and `/* */` inside a string are plain text.
	str := stringOf(t, "VAR := <|string:`// not a comment /* nor this */`|>")
	content := str.Content()
	require.Len(t, content, 1)
	assert.Equal(t, "// not a comment /* nor this */", content[0].(StringText).Value())
}

func TestStringTextAroundInterpolation(t *testing.T) {
	str := stringOf(t, "NUM := 1\nVAR := <|string:`a #{NUM} b`|>")
	content := str.Content()
	require.Len(t, content, 3)
	assert.Equal(t, "a ", content[0].(StringText).Value())
	_, isInterpolation := content[1].(Interpolation)
	assert.True(t, isInterpolation)
	assert.Equal(t, " b", content[2].(StringText).Value())
}

// --- Interpolation mode ---

func TestWhitespaceInsideInterpolationIsSkipped(t *testing.T) {
	// WS is hidden in IN_INTERPOLATION, so padding around the expression does
	// not become part of the AST.
	str := stringOf(t, "NUM := 1\nVAR := <|string:`#{  NUM  }`|>")
	content := str.Content()
	require.Len(t, content, 1)
	interpolation := content[0].(Interpolation)
	assert.Equal(t, "NUM", interpolation.Expression().(VariableRef).Name().Ref(context.Background()).Name())
}

func TestCommentInsideInterpolationIsIgnored(t *testing.T) {
	str := stringOf(t, "NUM := 1\nVAR := <|string:`#{NUM /* here */}`|>")
	content := str.Content()
	require.Len(t, content, 1)
	interpolation := content[0].(Interpolation)
	assert.Equal(t, "NUM", interpolation.Expression().(VariableRef).Name().Ref(context.Background()).Name())
}

func TestExpressionInsideInterpolation(t *testing.T) {
	str := stringOf(t, "NUM := 1\nVAR := <|string:`#{(NUM + 2) + 3}`|>")
	content := str.Content()
	require.Len(t, content, 1)
	interpolation := content[0].(Interpolation)
	_, isBinary := interpolation.Expression().(BinaryExpression)
	assert.True(t, isBinary)
}

func TestDeeplyNestedStringInterpolation(t *testing.T) {
	// default -> IN_STRING -> IN_INTERPOLATION -> IN_STRING -> IN_INTERPOLATION
	// -> IN_STRING, i.e. a mode stack six levels deep.
	str := stringOf(t, "NUM := 1\nVAR := <|string:`a#{`b#{`c#{NUM}c`}b`}a`|>")
	content := str.Content()
	require.Len(t, content, 3)
	assert.Equal(t, "a", content[0].(StringText).Value())
	assert.Equal(t, "a", content[2].(StringText).Value())

	level2 := content[1].(Interpolation).Expression().(StringLiteral).Content()
	require.Len(t, level2, 3)
	assert.Equal(t, "b", level2[0].(StringText).Value())

	level3 := level2[1].(Interpolation).Expression().(StringLiteral).Content()
	require.Len(t, level3, 3)
	assert.Equal(t, "c", level3[0].(StringText).Value())
	assert.Equal(t,
		"NUM",
		level3[1].(Interpolation).Expression().(VariableRef).Name().Ref(context.Background()).Name(),
	)
}

func TestConsecutiveStringsUseIndependentModeStacks(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("A := `one`\nB := `two`\nC := `three`")
	doc.AssertNoErrors()
	strings := test.FindAll[StringLiteral](doc)
	require.Len(t, strings, 3)
	for i, expected := range []string{"one", "two", "three"} {
		content := strings[i].Content()
		require.Len(t, content, 1)
		assert.Equal(t, expected, content[0].(StringText).Value())
	}
}

// --- Cross references out of an interpolation ---

func TestInterpolationReferencesVariableDeclaredLater(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := `#{LATER}`\nLATER := 42")
	doc.AssertNoErrors()
	ref := test.MustFindNode[Interpolation](doc)
	target := ref.Expression().(VariableRef).Name().Ref(doc.Ctx())
	require.NotNil(t, target)
	assert.Equal(t, "LATER", target.Name())
}

func TestInterpolationReferencesUnknownVariable(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := `#{MISSING}`")
	// The lexer handles the mode switches; the unresolved name is a linking
	// error, not a lexer error.
	assert.Empty(t, doc.Document.LexerErrors)
	assert.Empty(t, doc.Document.ParserErrors)
	ref := test.MustFindNode[Interpolation](doc)
	assert.NotNil(t, ref.Expression().(VariableRef).Name().Error())
}

// --- Malformed input ---

func TestUnterminatedString(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := `unterminated")
	// The lexer stays in IN_STRING to the end of input; the missing backtick is
	// reported by the parser.
	assert.Empty(t, doc.Document.LexerErrors)
	assert.NotEmpty(t, doc.Document.ParserErrors)
}

func TestUnterminatedInterpolation(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("NUM := 1\nVAR := `text #{NUM")
	assert.Empty(t, doc.Document.LexerErrors)
	assert.NotEmpty(t, doc.Document.ParserErrors)
}

func TestUnterminatedStringDoesNotAffectFollowingDocument(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// Both documents share one lexer instance. The first one ends inside
	// IN_STRING, which must not carry over.
	docs := fixture.ParseAll(
		"inmemory://broken.tm", "VAR := `unterminated",
		"inmemory://valid.tm", "OTHER := 123",
	)
	valid := docs[1]
	assert.Empty(t, valid.Document.LexerErrors)
	assert.Empty(t, valid.Document.ParserErrors)
	decl := test.MustFindNode[VariableDecl](valid)
	assert.Equal(t, "OTHER", decl.Name())
}

func TestStrayClosingBraceOutsideInterpolation(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := 1\n}")
	// "}" is only registered in IN_INTERPOLATION, so in the default mode it is
	// unrecognized input.
	require.NotEmpty(t, doc.Document.LexerErrors)
	assert.Equal(t, "No matching token", doc.Document.LexerErrors[0].Msg)
}

func TestBareHashInsideStringIsUnrecognized(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("VAR := `a # b`")
	// STRING_CONTENT excludes '#', and only "#{" enters the interpolation mode,
	// so a lone '#' cannot be matched in IN_STRING.
	require.NotEmpty(t, doc.Document.LexerErrors)
	assert.Equal(t, "No matching token", doc.Document.LexerErrors[0].Msg)
}

// bt is a single backtick, which delimits strings in this language.
const bt = "`"

func TestEscapesInsideString(t *testing.T) {
	// An escaped backtick or hash stays inside the string instead of switching
	// modes.
	str := stringOf(t, "VAR := <|string:"+bt+`esc\`+bt+` and \#{ ok`+bt+"|>")
	content := str.Content()
	require.Len(t, content, 1)
	assert.Equal(t, `esc\`+bt+` and \#{ ok`, content[0].(StringText).Value())
}

func TestEscapedBackslashBeforeStringTerminator(t *testing.T) {
	// STRING_CONTENT excludes the backslash from its literal alternative
	// (/([^`#\\]|\\[`#tnr\\])+/), so a backslash can only ever start an escape.
	// That keeps `\\` from consuming the closing backtick as `\` + `\` + "`",
	// which would lose the terminator.
	str := stringOf(t, "VAR := <|string:"+bt+`a\\`+bt+"|>")
	content := str.Content()
	require.Len(t, content, 1)
	assert.Equal(t, `a\\`, content[0].(StringText).Value())
}

func TestEscapedBackslashFollowedByInterpolation(t *testing.T) {
	// The same ambiguity would exist for "#{": `\\#{...}` must escape the
	// backslash and then enter the interpolation mode.
	str := stringOf(t, "NUM := 1\nVAR := <|string:"+bt+`a\\#{NUM}`+bt+"|>")
	content := str.Content()
	require.Len(t, content, 2)
	assert.Equal(t, `a\\`, content[0].(StringText).Value())
	interpolation, ok := content[1].(Interpolation)
	require.True(t, ok)
	assert.Equal(t, "NUM", interpolation.Expression().(VariableRef).Name().Ref(context.Background()).Name())
}

func TestAllEscapeSequencesInsideString(t *testing.T) {
	str := stringOf(t, "VAR := <|string:"+bt+`\`+bt+`\#\t\n\r\\`+bt+"|>")
	content := str.Content()
	require.Len(t, content, 1)
	assert.Equal(t, `\`+bt+`\#\t\n\r\\`, content[0].(StringText).Value())
}

func TestUnknownEscapeInsideStringIsUnrecognized(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// Now that a lone backslash is not literal text, an escape outside
	// [`#tnr\] cannot be matched at all.
	doc := fixture.Parse("VAR := " + bt + `a\xb` + bt)
	require.NotEmpty(t, doc.Document.LexerErrors)
	assert.Equal(t, "No matching token", doc.Document.LexerErrors[0].Msg)
	// The lexer recovers by dropping the offending rune and stays in IN_STRING,
	// so the closing backtick still terminates the string.
	assert.Empty(t, doc.Document.ParserErrors)
}

func TestEscapedTerminatorLeavesStringUnterminated(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// A single backslash before the closing backtick escapes it, so the whole
	// input is one unterminated string rather than a lexer error.
	doc := fixture.Parse("VAR := " + bt + `a\` + bt)
	assert.Empty(t, doc.Document.LexerErrors)
	assert.NotEmpty(t, doc.Document.ParserErrors)

	// Adding a real terminator closes the string, with the escape kept as text.
	closed := test.New(t, CreateServices()).Parse("VAR := <|string:" + bt + `a\` + bt + bt + "|>")
	closed.AssertNoErrors()
	node := test.MustFindNodeAtLabel[StringLiteral](closed, "string")
	content := node.Content()
	require.Len(t, content, 1)
	assert.Equal(t, `a\`+bt, content[0].(StringText).Value())
}
