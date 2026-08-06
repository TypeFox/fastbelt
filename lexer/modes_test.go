// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
)

// literal returns a token type that matches the exact string value.
func literal(id int, name, value string) *core.TokenType {
	return core.NewTokenType(id, name, name, core.TokenKindKeyword,
		func(input string, offset int) int {
			if strings.HasPrefix(input[offset:], value) {
				return len(value)
			}
			return 0
		},
		[]rune(value)[:1],
	)
}

// matchRunes returns a token type that matches a run of any of the given runes.
func matchRunes(id int, name string, runes string) *core.TokenType {
	return core.NewTokenType(id, name, name, core.TokenKindToken,
		func(input string, offset int) int {
			length := 0
			for _, r := range input[offset:] {
				if !strings.ContainsRune(runes, r) {
					break
				}
				length += len(string(r))
			}
			return length
		},
		[]rune(runes),
	)
}

// --- NewTokenMode: start character bucketing ---

func TestNewTokenModeBucketsByStartChar(t *testing.T) {
	a := literal(1, "A", "a")
	b := literal(2, "B", "b")
	mode := NewTokenMode("test", UseTokenType(a), UseTokenType(b))

	require.Len(t, mode.TokenMap, maxChar)
	require.Len(t, mode.TokenMap['a'], 1)
	assert.Same(t, a, mode.TokenMap['a'][0].TokenType)
	require.Len(t, mode.TokenMap['b'], 1)
	assert.Same(t, b, mode.TokenMap['b'][0].TokenType)
	assert.Empty(t, mode.TokenMap['c'])
}

func TestNewTokenModeSharedStartChar(t *testing.T) {
	short := literal(1, "Short", "a")
	long := literal(2, "Long", "ab")
	mode := NewTokenMode("test", UseTokenType(short), UseTokenType(long))

	// Both are candidates for 'a' and stay in registration order.
	require.Len(t, mode.TokenMap['a'], 2)
	assert.Same(t, short, mode.TokenMap['a'][0].TokenType)
	assert.Same(t, long, mode.TokenMap['a'][1].TokenType)
}

func TestNewTokenModeMultipleStartChars(t *testing.T) {
	digits := matchRunes(1, "DIGITS", "0123456789")
	mode := NewTokenMode("test", UseTokenType(digits))

	for _, r := range "0123456789" {
		assert.Len(t, mode.TokenMap[r], 1, "expected candidate for %q", r)
	}
	assert.Empty(t, mode.TokenMap['a'])
}

func TestNewTokenModeWithoutStartCharsIsNeverACandidate(t *testing.T) {
	// A token type without start characters cannot be preselected, so it is
	// unreachable no matter what the input is.
	unreachable := core.NewTokenType(1, "NONE", "NONE", core.TokenKindToken,
		func(string, int) int { return 1 }, nil)
	mode := NewTokenMode("test", UseTokenType(unreachable))

	for index := range maxChar {
		assert.Empty(t, mode.TokenMap[index])
	}
	// It is still registered on the mode itself.
	require.Len(t, mode.TokenTypes, 1)
}

func TestNewTokenModeNonAsciiStartCharWrapsIntoBucket(t *testing.T) {
	// Start characters are bucketed by rune value modulo maxChar, so a non-ASCII
	// rune shares a bucket with the ASCII rune it collides with. Both must
	// remain candidates for their own input.
	umlaut := literal(1, "UMLAUT", "ä") // U+00E4, 0xE4 % 256 == 0xE4
	mode := NewTokenMode("test", UseTokenType(umlaut))

	require.Len(t, mode.TokenMap[0xE4], 1)
	assert.Same(t, umlaut, mode.TokenMap[0xE4][0].TokenType)
}

func TestNewTokenModeCollidingStartCharsKeepBothCandidates(t *testing.T) {
	// U+0161 % 256 == 0x61 == 'a', so these two share a bucket.
	ascii := literal(1, "ASCII", "a")
	wrapped := literal(2, "WRAPPED", "š")
	mode := NewTokenMode("test", UseTokenType(ascii), UseTokenType(wrapped))

	require.Len(t, mode.TokenMap['a'], 2)
}

// --- TokenTypeUsage builders ---

func TestUseTokenTypeDefaults(t *testing.T) {
	usage := UseTokenType(literal(1, "A", "a"))
	assert.Equal(t, -1, usage.PushMode)
	assert.False(t, usage.PopMode)
	assert.Equal(t, 0, usage.Group)
	assert.False(t, usage.IsSkipped())
	assert.False(t, usage.IsComment())
}

func TestTokenTypeUsageBuilders(t *testing.T) {
	tokenType := literal(1, "A", "a")

	push := UseTokenType(tokenType).WithPushMode(3)
	assert.Equal(t, 3, push.PushMode)
	assert.False(t, push.PopMode)

	pop := UseTokenType(tokenType).WithPopMode()
	assert.True(t, pop.PopMode)
	assert.Equal(t, -1, pop.PushMode)

	set := UseTokenType(tokenType).WithSetMode(2)
	assert.True(t, set.PopMode)
	assert.Equal(t, 2, set.PushMode)

	assert.True(t, UseTokenType(tokenType).WithGroup(core.SkippedGroup).IsSkipped())
	assert.True(t, UseTokenType(tokenType).WithGroup(core.CommentGroup).IsComment())
	assert.False(t, UseTokenType(tokenType).WithGroup(7).IsSkipped())
}

// --- TokenModeStack ---

func TestTokenModeStackPushPopPeek(t *testing.T) {
	first := NewTokenMode("first")
	second := NewTokenMode("second")
	stack := NewTokenModeStack(first)

	assert.Same(t, first, stack.Peek())
	stack.Push(second)
	assert.Same(t, second, stack.Peek())
	assert.Same(t, second, stack.Pop())
	assert.Same(t, first, stack.Peek())
}

func TestTokenModeStackPopAtBottomIsNoOp(t *testing.T) {
	first := NewTokenMode("first")
	stack := NewTokenModeStack(first)

	// An unbalanced pop keeps the start mode active rather than emptying the
	// stack, so lexing can continue.
	assert.Same(t, first, stack.Pop())
	assert.Same(t, first, stack.Peek())
	assert.Same(t, first, stack.Pop())
	assert.Same(t, first, stack.Peek())
}

func TestTokenModeStackSetModeReplacesActiveMode(t *testing.T) {
	first := NewTokenMode("first")
	second := NewTokenMode("second")
	third := NewTokenMode("third")
	stack := NewTokenModeStack(first)

	stack.Push(second)
	stack.SetMode(third)
	assert.Same(t, third, stack.Peek())
	// SetMode replaced `second` instead of stacking on top of it, so popping
	// returns to `first`.
	assert.Same(t, third, stack.Pop())
	assert.Same(t, first, stack.Peek())
}

func TestTokenModeStackSetModeAtBottomKeepsDepth(t *testing.T) {
	first := NewTokenMode("first")
	second := NewTokenMode("second")
	stack := NewTokenModeStack(first)

	stack.SetMode(second)
	assert.Same(t, second, stack.Peek())
	// The stack is still one deep, so a pop cannot get back to `first`.
	assert.Same(t, second, stack.Pop())
	assert.Same(t, second, stack.Peek())
}
