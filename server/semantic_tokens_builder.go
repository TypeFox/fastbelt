// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"unicode/utf16"
	"unicode/utf8"

	core "typefox.dev/fastbelt"
)

// SemanticTokensBuilder defines the interface for building semantic tokens data in the LSP format.
// It provides methods to push individual tokens and retrieve the final data slice.
//
// It is recommended to build the semantic tokens data by using the [TokenBasedSemanticTokensProvider] and its
// associated [TokenHighlightingStrategy] implementations.
type SemanticTokensBuilder interface {
	// Data returns the final semantic tokens data slice in the LSP format.
	// The LSP data slice is a flat array of uint32 values, where each token is represented by five consecutive values:
	// (1) deltaLine: token line number, relative to the previous token,
	// (2) deltaStart: token start character, relative to the previous token,
	// (3) length: the length of the token,
	// (4) tokenType: the token type index in the legend, and
	// (5) tokenModifiers: the token modifiers bitset in the legend.
	Data() []uint32
	// Push adds a new token to the semantic tokens data.
	// Tokens should be pushed in the order they appear in the document, as the LSP
	// data is based on the offset of the previous token.
	// Due to this, this method is not thread-safe.
	Push(textRange core.TextRange, tokenType, tokenModifiers uint32)
}

// NewSemanticTokensBuilder creates a new instance of [SemanticTokensBuilder].
func NewSemanticTokensBuilder(text string, tokenCount int) SemanticTokensBuilder {
	return &semanticTokensBuilder{
		// Preallocate the data with a reasonable length
		// Each token potentially contributes 5 uint32 values
		// Multiline tokens can contribute more, but we can resize the slice if needed
		data: make([]uint32, 0, tokenCount*5),
		text: text,
	}
}

type semanticTokensBuilder struct {
	// data is a slice of uint32 values representing the semantic tokens data in the LSP format.
	// Each token is represented by five consecutive values:
	// - deltaLine, token line number, relative to the previous token,
	// - deltaStart, token start character, relative to the previous token,
	// - length, the length of the token,
	// - tokenType, the token type index,
	// - tokenModifiers, the token modifiers bitset.
	data []uint32
	text string
	// cursor is the byte offset in text, line/char the corresponding LSP (line, UTF-16 column) position
	cursor, line, char int
	// prevLine/prevChar is the start position of the last emitted segment
	prevLine, prevChar int
}

func (b *semanticTokensBuilder) Data() []uint32 {
	return b.data
}

func (b *semanticTokensBuilder) Push(textRange core.TextRange, typeIndex, modifierIndex uint32) {
	// Out-of-order or out-of-range pushes violate the contract; clamp them instead of wrapping deltas
	tokenStart := min(max(int(textRange.Start), b.cursor), len(b.text))
	tokenEnd := min(max(int(textRange.End), tokenStart), len(b.text))
	for b.cursor < tokenStart {
		b.step()
	}
	segLine, segStart := b.line, b.char
	for b.cursor < tokenEnd {
		lineEnd := b.char
		if b.step() {
			// Token spans multiple lines, emit one segment per line
			b.emit(segLine, segStart, lineEnd, typeIndex, modifierIndex)
			segLine, segStart = b.line, 0
		}
	}
	b.emit(segLine, segStart, b.char, typeIndex, modifierIndex)
}

// emit appends a single-line segment [start, end) on the given line, skipping empty segments.
func (b *semanticTokensBuilder) emit(line, start, end int, typeIndex, modifierIndex uint32) {
	if end <= start {
		return
	}
	charDelta := start
	if line == b.prevLine {
		charDelta -= b.prevChar
	}
	b.data = append(b.data, uint32(line-b.prevLine), uint32(charDelta), uint32(end-start), typeIndex, modifierIndex)
	b.prevLine, b.prevChar = line, start
}

// step advances the cursor by one character and reports whether it crossed a line break.
// Line breaks follow the same rules as [textdoc]: "\r\n", "\r" and "\n".
func (b *semanticTokensBuilder) step() bool {
	c := b.text[b.cursor]
	switch {
	case c == '\r':
		b.cursor++
		if b.cursor < len(b.text) && b.text[b.cursor] == '\n' {
			b.cursor++
		}
	case c == '\n':
		b.cursor++
	case c < utf8.RuneSelf:
		// ASCII fast path: one byte, one UTF-16 code unit
		b.cursor++
		b.char++
		return false
	default:
		r, size := utf8.DecodeRuneInString(b.text[b.cursor:])
		b.cursor += size
		b.char += utf16.RuneLen(r)
		return false
	}
	b.line++
	b.char = 0
	return true
}
