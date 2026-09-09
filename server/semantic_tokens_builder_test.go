// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"slices"
	"strings"
	"testing"

	core "typefox.dev/fastbelt"
)

func TestLspTokenDataPush(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		ranges   []core.TextRange
		expected []uint32
	}{
		{
			name:     "Single token at start",
			text:     "hello world",
			ranges:   []core.TextRange{core.NewTextRange(0, 5)},
			expected: []uint32{0, 0, 5, 1, 2},
		},
		{
			name:   "Two tokens on same line use char delta",
			text:   "hello world",
			ranges: []core.TextRange{core.NewTextRange(0, 5), core.NewTextRange(6, 11)},
			expected: []uint32{
				0, 0, 5, 1, 2,
				0, 6, 5, 1, 2,
			},
		},
		{
			name:   "Token on next line resets char delta",
			text:   "hello\nworld",
			ranges: []core.TextRange{core.NewTextRange(0, 5), core.NewTextRange(6, 11)},
			expected: []uint32{
				0, 0, 5, 1, 2,
				1, 0, 5, 1, 2,
			},
		},
		{
			name:   "Multi-line token emits one token per line",
			text:   "ab\ncdef\ngh",
			ranges: []core.TextRange{core.NewTextRange(1, 9)},
			expected: []uint32{
				0, 1, 1, 1, 2, // "b" on line 0
				1, 0, 4, 1, 2, // "cdef" on line 1
				1, 0, 1, 1, 2, // "g" on line 2
			},
		},
		{
			name:   "Token after multi-line token",
			text:   "ab\ncd ef",
			ranges: []core.TextRange{core.NewTextRange(0, 5), core.NewTextRange(6, 8)},
			expected: []uint32{
				0, 0, 2, 1, 2,
				1, 0, 2, 1, 2,
				0, 3, 2, 1, 2,
			},
		},
		{
			name: "Non-ASCII counts UTF-16 code units",
			// "😀" is 4 bytes but 2 UTF-16 code units
			text:   "😀ab",
			ranges: []core.TextRange{core.NewTextRange(4, 6)},
			expected: []uint32{
				0, 2, 2, 1, 2,
			},
		},
		{
			name:     "Range past end of text is clamped",
			text:     "ab",
			ranges:   []core.TextRange{core.NewTextRange(0, 10)},
			expected: []uint32{0, 0, 2, 1, 2},
		},
		{
			name:   "CRLF line breaks are not counted as columns",
			text:   "/* x\r\n y */\r\nab",
			ranges: []core.TextRange{core.NewTextRange(0, 11), core.NewTextRange(13, 15)},
			expected: []uint32{
				0, 0, 4, 1, 2,
				1, 0, 5, 1, 2,
				1, 0, 2, 1, 2,
			},
		},
		{
			name:   "Lone CR is a line break",
			text:   "ab\rcd",
			ranges: []core.TextRange{core.NewTextRange(0, 2), core.NewTextRange(3, 5)},
			expected: []uint32{
				0, 0, 2, 1, 2,
				1, 0, 2, 1, 2,
			},
		},
		{
			name:   "Empty lines inside a token emit no segments",
			text:   "/*\n\n*/\nx",
			ranges: []core.TextRange{core.NewTextRange(0, 7), core.NewTextRange(7, 8)},
			expected: []uint32{
				0, 0, 2, 1, 2,
				2, 0, 2, 1, 2,
				1, 0, 1, 1, 2,
			},
		},
		{
			name:   "Out-of-order push degrades to current position",
			text:   "a\n\nfoo",
			ranges: []core.TextRange{core.NewTextRange(3, 6), core.NewTextRange(3, 6)},
			expected: []uint32{
				2, 0, 3, 1, 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSemanticTokensBuilder(tt.text, len(tt.ranges))
			for _, rng := range tt.ranges {
				builder.Push(rng, 1, 2)
			}
			if !slices.Equal(builder.Data(), tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, builder.Data())
			}
		})
	}
}

func BenchmarkLspTokenDataPush(b *testing.B) {
	// Build a document of 1000 lines with 4 tokens each
	line := "foo bar baz qux\n"
	text := strings.Repeat(line, 1000)
	ranges := []core.TextRange{}
	for i := range 1000 {
		offset := i * len(line)
		for start := 0; start < 15; start += 4 {
			ranges = append(ranges, core.NewTextRange(offset+start, offset+start+3))
		}
	}

	for b.Loop() {
		builder := NewSemanticTokensBuilder(text, len(ranges))
		for _, rng := range ranges {
			builder.Push(rng, 1, 2)
		}
	}
}
