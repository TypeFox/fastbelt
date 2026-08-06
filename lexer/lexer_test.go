// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	spaces    = " \t\r\n"
)

// images returns the token images of tokens, for compact assertions.
func images(tokens []core.Token) []string {
	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = token.Image
	}
	return result
}

// names returns the token type names of tokens.
func names(tokens []core.Token) []string {
	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = token.Type.Name
	}
	return result
}

// stringLexer builds a two-mode lexer: the default mode reads identifiers and
// enters IN_STRING on a quote, IN_STRING reads text and leaves on a quote.
// customGroup is the group assigned to IN_STRING's text token.
func stringLexer(customGroup int) *DefaultLexer {
	const (
		modeDefault = 0
		modeString  = 1
	)
	id := matchRunes(1, "ID", lowercase)
	quote := literal(2, "QUOTE", `"`)
	ws := matchRunes(3, "WS", spaces)
	text := matchRunes(4, "TEXT", lowercase+spaces)

	modes := make([]*TokenMode, 2)
	modes[modeDefault] = NewTokenMode("default",
		UseTokenType(quote).WithPushMode(modeString),
		UseTokenType(id),
		UseTokenType(ws).WithGroup(core.SkippedGroup),
	)
	modes[modeString] = NewTokenMode("IN_STRING",
		UseTokenType(quote).WithPopMode(),
		UseTokenType(text).WithGroup(customGroup),
	)
	return NewDefaultLexer(modeDefault, modes...)
}

// --- Longest match and tie-breaking ---

func TestExecPrefersLongestMatch(t *testing.T) {
	short := literal(1, "SHORT", "a")
	long := literal(2, "LONG", "abc")
	lexer := NewDefaultLexer(0, NewTokenMode("default",
		UseTokenType(short),
		UseTokenType(long),
	))

	result := lexer.Exec("abc")
	require.Len(t, result.Tokens, 1)
	assert.Equal(t, "LONG", result.Tokens[0].Type.Name)
}

func TestExecFirstRegisteredWinsEqualLengthMatch(t *testing.T) {
	first := literal(1, "FIRST", "a")
	second := literal(2, "SECOND", "a")
	lexer := NewDefaultLexer(0, NewTokenMode("default",
		UseTokenType(first),
		UseTokenType(second),
	))

	result := lexer.Exec("a")
	require.Len(t, result.Tokens, 1)
	assert.Equal(t, "FIRST", result.Tokens[0].Type.Name,
		"among equal-length matches the first registered token type must win")
}

// --- Groups ---

func TestExecRoutesTokensToGroups(t *testing.T) {
	word := matchRunes(1, "WORD", lowercase)
	ws := matchRunes(2, "WS", spaces)
	comment := literal(3, "COMMENT", "#")
	other := literal(4, "OTHER", "!")
	lexer := NewDefaultLexer(0, NewTokenMode("default",
		UseTokenType(word),
		UseTokenType(ws).WithGroup(core.SkippedGroup),
		UseTokenType(comment).WithGroup(core.CommentGroup),
		UseTokenType(other).WithGroup(7),
	))

	result := lexer.Exec("ab #ff!")
	assert.Equal(t, []string{"ab", "ff"}, images(result.Tokens))
	assert.Equal(t, []string{"#"}, images(result.Comments))
	require.Contains(t, result.Groups, 7)
	assert.Equal(t, []string{"!"}, images(result.Groups[7]))
	assert.Empty(t, result.Errors)
}

func TestExecGroupsIsNilWithoutCustomGroups(t *testing.T) {
	word := matchRunes(1, "WORD", lowercase)
	lexer := NewDefaultLexer(0, NewTokenMode("default", UseTokenType(word)))

	assert.Nil(t, lexer.Exec("ab").Groups)
}

func TestExecCollectsCustomGroupFromNonDefaultMode(t *testing.T) {
	lexer := stringLexer(7)

	result := lexer.Exec(`ab "inside" cd`)
	assert.Equal(t, []string{"ab", `"`, `"`, "cd"}, images(result.Tokens))
	require.Contains(t, result.Groups, 7)
	assert.Equal(t, []string{"inside"}, images(result.Groups[7]))
}

func TestExecSkipsHiddenTokenFromNonDefaultMode(t *testing.T) {
	lexer := stringLexer(core.SkippedGroup)

	result := lexer.Exec(`ab "inside" cd`)
	assert.Equal(t, []string{"ab", `"`, `"`, "cd"}, images(result.Tokens))
	assert.Nil(t, result.Groups)
}

func TestExecCollectsCommentFromNonDefaultMode(t *testing.T) {
	lexer := stringLexer(core.CommentGroup)

	result := lexer.Exec(`ab "inside" cd`)
	assert.Equal(t, []string{"inside"}, images(result.Comments))
}

// --- Mode switching ---

func TestExecPushAndPopMode(t *testing.T) {
	lexer := stringLexer(0)

	result := lexer.Exec(`ab "in string" cd`)
	assert.Equal(t, []string{"ID", "QUOTE", "TEXT", "QUOTE", "ID"}, names(result.Tokens))
	assert.Equal(t, []string{"ab", `"`, "in string", `"`, "cd"}, images(result.Tokens))
	assert.Empty(t, result.Errors)
}

func TestExecNestedPushMode(t *testing.T) {
	// outer -> middle -> inner, then two pops back to outer
	outer := literal(1, "OUTER", "o")
	enterMiddle := literal(2, "ENTER_MIDDLE", "(")
	enterInner := literal(3, "ENTER_INNER", "[")
	leave := literal(4, "LEAVE", ")")
	middleBody := literal(5, "MIDDLE", "m")
	innerBody := literal(6, "INNER", "i")

	modes := make([]*TokenMode, 3)
	modes[0] = NewTokenMode("default",
		UseTokenType(outer),
		UseTokenType(enterMiddle).WithPushMode(1),
	)
	modes[1] = NewTokenMode("middle",
		UseTokenType(middleBody),
		UseTokenType(enterInner).WithPushMode(2),
		UseTokenType(leave).WithPopMode(),
	)
	modes[2] = NewTokenMode("inner",
		UseTokenType(innerBody),
		UseTokenType(leave).WithPopMode(),
	)
	lexer := NewDefaultLexer(0, modes...)

	result := lexer.Exec("o(m[i))o")
	assert.Equal(t, []string{
		"OUTER", "ENTER_MIDDLE", "MIDDLE", "ENTER_INNER", "INNER", "LEAVE", "LEAVE", "OUTER",
	}, names(result.Tokens))
	assert.Empty(t, result.Errors)
}

func TestExecSetModeDoesNotGrowTheModeStack(t *testing.T) {
	// `a` switches from the start mode to `second` with a set-mode command.
	// Because set-mode replaces the active mode rather than pushing onto the
	// stack, the following pop has nothing to return to and `second` stays
	// active - which is why `y` from the start mode no longer matches.
	a := literal(1, "A", "a")
	b := literal(2, "B", "b")
	inSecond := literal(3, "IN_SECOND", "x")
	inDefault := literal(4, "IN_DEFAULT", "y")

	modes := make([]*TokenMode, 2)
	modes[0] = NewTokenMode("default",
		UseTokenType(a).WithSetMode(1),
		UseTokenType(inDefault),
	)
	modes[1] = NewTokenMode("second",
		UseTokenType(inSecond),
		UseTokenType(b).WithPopMode(),
	)
	lexer := NewDefaultLexer(0, modes...)

	result := lexer.Exec("axbxy")
	assert.Equal(t, []string{"A", "IN_SECOND", "B", "IN_SECOND"}, names(result.Tokens))
	require.Len(t, result.Errors, 1, "y belongs to the start mode, which is no longer active")
	assert.Equal(t, int32(4), result.Errors[0].Range.Start)
}

func TestExecSetModeKeepsOuterModeReachable(t *testing.T) {
	// default -pushes-> middle -sets-> other -pops-> default
	enter := literal(1, "ENTER", "(")
	set := literal(2, "SET", "s")
	leave := literal(3, "LEAVE", ")")
	inMiddle := literal(4, "IN_MIDDLE", "m")
	inOther := literal(5, "IN_OTHER", "o")
	inDefault := literal(6, "IN_DEFAULT", "d")

	modes := make([]*TokenMode, 3)
	modes[0] = NewTokenMode("default",
		UseTokenType(enter).WithPushMode(1),
		UseTokenType(inDefault),
	)
	modes[1] = NewTokenMode("middle",
		UseTokenType(inMiddle),
		UseTokenType(set).WithSetMode(2),
	)
	modes[2] = NewTokenMode("other",
		UseTokenType(inOther),
		UseTokenType(leave).WithPopMode(),
	)
	lexer := NewDefaultLexer(0, modes...)

	result := lexer.Exec("d(mso)d")
	assert.Equal(t, []string{
		"IN_DEFAULT", "ENTER", "IN_MIDDLE", "SET", "IN_OTHER", "LEAVE", "IN_DEFAULT",
	}, names(result.Tokens))
	assert.Empty(t, result.Errors)
}

func TestExecUnbalancedPopStaysInStartMode(t *testing.T) {
	// A pop in the start mode has nothing to return to and must not break
	// lexing of the remaining input.
	pop := literal(1, "POP", ")")
	word := matchRunes(2, "WORD", lowercase)
	lexer := NewDefaultLexer(0, NewTokenMode("default",
		UseTokenType(pop).WithPopMode(),
		UseTokenType(word),
	))

	result := lexer.Exec(")ab)cd")
	assert.Equal(t, []string{"POP", "WORD", "POP", "WORD"}, names(result.Tokens))
	assert.Empty(t, result.Errors)
}

func TestExecPopAndGroupApplyTogether(t *testing.T) {
	// The token that leaves the mode is itself routed to the comment group.
	enter := literal(1, "ENTER", "(")
	leave := literal(2, "LEAVE", ")")
	inside := matchRunes(3, "INSIDE", lowercase)
	outside := literal(4, "OUTSIDE", "!")

	modes := make([]*TokenMode, 2)
	modes[0] = NewTokenMode("default",
		UseTokenType(enter).WithPushMode(1),
		UseTokenType(outside),
	)
	modes[1] = NewTokenMode("inner",
		UseTokenType(leave).WithPopMode().WithGroup(core.CommentGroup),
		UseTokenType(inside),
	)
	lexer := NewDefaultLexer(0, modes...)

	result := lexer.Exec("(ab)!")
	assert.Equal(t, []string{"ENTER", "INSIDE", "OUTSIDE"}, names(result.Tokens))
	assert.Equal(t, []string{")"}, images(result.Comments))
}

// --- Error handling ---

func TestExecReportsUnmatchedInput(t *testing.T) {
	word := matchRunes(1, "WORD", lowercase)
	lexer := NewDefaultLexer(0, NewTokenMode("default", UseTokenType(word)))

	result := lexer.Exec("ab!cd")
	assert.Equal(t, []string{"ab", "cd"}, images(result.Tokens))
	require.Len(t, result.Errors, 1)
	assert.Equal(t, int32(2), result.Errors[0].Range.Start)
	assert.Equal(t, int32(3), result.Errors[0].Range.End)
}

func TestExecUnmatchedInputConsumesWholeRune(t *testing.T) {
	// A multi-byte rune that matches nothing must be skipped as one unit, so
	// the next offset stays on a rune boundary.
	word := matchRunes(1, "WORD", lowercase)
	lexer := NewDefaultLexer(0, NewTokenMode("default", UseTokenType(word)))

	result := lexer.Exec("aä b")
	require.Len(t, result.Errors, 2)
	assert.Equal(t, int32(1), result.Errors[0].Range.Start)
	assert.Equal(t, int32(3), result.Errors[0].Range.End, "the two-byte rune must be consumed as a whole")
	assert.Equal(t, []string{"a", "b"}, images(result.Tokens))
}

func TestExecUnmatchedInputKeepsActiveMode(t *testing.T) {
	// An unrecognized character inside a pushed mode must not change the mode:
	// the closing quote still pops.
	lexer := stringLexer(0)

	result := lexer.Exec(`"in#side" ab`)
	assert.Equal(t, []string{"QUOTE", "TEXT", "TEXT", "QUOTE", "ID"}, names(result.Tokens))
	require.Len(t, result.Errors, 1)
	assert.Equal(t, int32(3), result.Errors[0].Range.Start)
}

func TestExecUnterminatedModeAtEndOfInput(t *testing.T) {
	lexer := stringLexer(0)

	result := lexer.Exec(`ab "unterminated`)
	assert.Equal(t, []string{"ID", "QUOTE", "TEXT"}, names(result.Tokens))
	assert.Equal(t, []string{"ab", `"`, "unterminated"}, images(result.Tokens))
	// Leaving a mode open is not a lexer error; the parser reports the
	// unfinished construct.
	assert.Empty(t, result.Errors)
}

func TestExecEmptyInput(t *testing.T) {
	lexer := stringLexer(0)

	result := lexer.Exec("")
	assert.Empty(t, result.Tokens)
	assert.Empty(t, result.Comments)
	assert.Empty(t, result.Errors)
	assert.Nil(t, result.Groups)
}

// --- Reuse of a single lexer instance ---

func TestExecResetsModeStackBetweenRuns(t *testing.T) {
	// A lexer is registered once per language and reused for every document, so
	// input that ends inside a pushed mode must not affect the next run.
	lexer := stringLexer(0)

	first := lexer.Exec(`"unterminated`)
	require.Equal(t, []string{"QUOTE", "TEXT"}, names(first.Tokens))

	second := lexer.Exec(`ab "in string" cd`)
	assert.Equal(t, []string{"ID", "QUOTE", "TEXT", "QUOTE", "ID"}, names(second.Tokens))
	assert.Empty(t, second.Errors)
}

func TestExecRepeatedRunsAreIdentical(t *testing.T) {
	lexer := stringLexer(0)
	const input = `ab "in string" cd`

	first := names(lexer.Exec(input).Tokens)
	for range 5 {
		assert.Equal(t, first, names(lexer.Exec(input).Tokens))
	}
}

// TestExecIsSafeForConcurrentUse guards the shared lexer instance that the
// workspace builder uses to parse documents in parallel. Run with -race.
func TestExecIsSafeForConcurrentUse(t *testing.T) {
	lexer := stringLexer(0)
	inputs := []string{
		`ab "in string" cd`,
		`"unterminated`,
		`ab cd`,
		`"a" "b" "c"`,
	}
	expected := make([][]string, len(inputs))
	for i, input := range inputs {
		expected[i] = names(lexer.Exec(input).Tokens)
	}

	var wg sync.WaitGroup
	for round := range 50 {
		for i, input := range inputs {
			wg.Add(1)
			go func(i int, input string) {
				defer wg.Done()
				assert.Equal(t, expected[i], names(lexer.Exec(input).Tokens),
					"round %d, input %q", round, input)
			}(i, input)
		}
	}
	wg.Wait()
}
