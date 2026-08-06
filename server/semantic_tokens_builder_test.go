// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"slices"
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
	text := ""
	ranges := []core.TextRange{}
	for range 1000 {
		offset := len(text)
		for start := 0; start < 15; start += 4 {
			ranges = append(ranges, core.NewTextRange(offset+start, offset+start+3))
		}
		text += line
	}

	for b.Loop() {
		builder := NewSemanticTokensBuilder(text, len(ranges))
		for _, rng := range ranges {
			builder.Push(rng, 1, 2)
		}
	}
}
