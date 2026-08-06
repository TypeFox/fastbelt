// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"unicode/utf16"
	"unicode/utf8"

	core "typefox.dev/fastbelt"
)

type SemanticTokensBuilder interface {
	Data() []uint32
	Push(textRange core.TextRange, tokenType, tokenModifiers uint32)
}

func NewSemanticTokensBuilder(text string, tokenCount int) SemanticTokensBuilder {
	return &semanticTokensBuilder{
		// Preallocate the data with the maximum possible length
		// Each token potentially contributes 5 uint32 values
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
	data        []uint32
	text        string
	cursor      int
	prevLine    int
	prevChar    int
	currentLine int
	currentChar int
	// lineBreaks is reused across push calls to avoid per-token allocations
	lineBreaks []int
}

func (tokenData *semanticTokensBuilder) Data() []uint32 {
	return tokenData.data
}

func (tokenData *semanticTokensBuilder) Push(textRange core.TextRange, typeIndex, modifierIndex uint32) {
	textLen := len(tokenData.text)
	tokenStart := int(textRange.Start)
	tokenEnd := int(textRange.End)
	cursor := tokenData.cursor
	currentLine := tokenData.currentLine
	currentChar := tokenData.currentChar
	startLine, startChar := 0, 0
	// Count line breaks within the token range as necessary
	// We need to emit multiple tokens if the token spans multiple lines
	lineBreaks := tokenData.lineBreaks[:0]
	// Advance the cursor up to the end of the token
	for tokenEnd > cursor {
		if cursor >= textLen {
			break
		} else if cursor == tokenStart {
			startLine = currentLine
			startChar = currentChar
		}
		if c := tokenData.text[cursor]; c < utf8.RuneSelf {
			// ASCII fast path: one byte, one UTF-16 code unit
			if c == '\n' {
				// Record the line break character position
				if cursor >= tokenStart {
					lineBreaks = append(lineBreaks, currentChar)
				}
				// New line, reset currentChar and increment currentLine
				currentLine++
				currentChar = 0
			} else {
				currentChar++
			}
			cursor++
			continue
		}
		rune, size := utf8.DecodeRuneInString(tokenData.text[cursor:])
		// Advance column by the number of UTF-16 code units for the rune
		// (newlines are ASCII, so this rune can never be one)
		currentChar += utf16.RuneLen(rune)
		// Advance cursor by the byte size of the rune
		cursor += size
	}
	tokenData.lineBreaks = lineBreaks
	lineDelta := uint32(startLine - tokenData.prevLine)
	charDelta := uint32(startChar)
	if lineDelta == 0 {
		// If the token is on the same line as the previous token, calculate the character delta
		charDelta -= uint32(tokenData.prevChar)
	}
	if len(lineBreaks) == 0 {
		// Token is on a single line, emit it directly
		length := uint32(currentChar - startChar)
		tokenData.data = append(tokenData.data, lineDelta, charDelta, length, typeIndex, modifierIndex)
		// Update the previous character position for the next token
		tokenData.prevChar = startChar
	} else {
		// Token spans multiple lines, emit a token for each line segment
		// First segment: from startChar to the first line break
		length := uint32(lineBreaks[0] - startChar)
		tokenData.data = append(tokenData.data, lineDelta, charDelta, length, typeIndex, modifierIndex)
		// Subsequent segments: from each line break to the next line break
		for i := 1; i < len(lineBreaks); i++ {
			// always use the full length of the line
			length = uint32(lineBreaks[i])
			// Note: lineDelta is always 1, since each segment is on a new line
			// charDelta is always 0, since we are starting at the beginning of the line
			tokenData.data = append(tokenData.data, 1, 0, length, typeIndex, modifierIndex)
		}
		// Last segment: from the start of the last line to the end of the token
		length = uint32(currentChar)
		tokenData.data = append(tokenData.data, 1, 0, length, typeIndex, modifierIndex)
		tokenData.prevChar = 0
	}
	// Update the data for the next token
	tokenData.cursor = cursor
	tokenData.prevLine = currentLine
	tokenData.currentLine = currentLine
	tokenData.currentChar = currentChar
}
