// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"errors"

	"github.com/TypeFox/go-lsp/protocol"
)

// File represents an immutable text document read from the file system.
// It implements Handle but does not support modifications.
type File struct {
	uri         protocol.DocumentURI
	languageID  string
	version     int32
	content     []byte
	lineOffsets []int
}

// NewFile creates a new text document file.
// The content parameter is accepted as a string and stored internally as bytes.
func NewFile(uri protocol.DocumentURI, languageID string, version int32, content string) (*File, error) {
	if uri == "" {
		return nil, errors.New("textdoc: uri cannot be empty")
	}
	if languageID == "" {
		return nil, errors.New("textdoc: languageID cannot be empty")
	}
	f := &File{
		uri:        uri,
		languageID: languageID,
		version:    version,
		content:    []byte(content),
	}
	f.lineOffsets = computeLineOffsets(f.content, true, 0)
	return f, nil
}

// URI returns the document URI.
func (f *File) URI() protocol.DocumentURI {
	return f.uri
}

// LanguageID returns the language identifier.
func (f *File) LanguageID() string {
	return f.languageID
}

// Version returns the document version.
func (f *File) Version() int32 {
	return f.version
}

// Content returns a copy of the document content as a byte slice.
// Returning bytes enables efficient manipulation; a copy is returned to prevent
// external modification of internal state.
func (f *File) Content() []byte {
	return f.content
}

// Text returns the text content or a substring if range is provided.
// This is a convenience method that returns a string instead of []byte.
func (f *File) Text(r *protocol.Range) string {
	if r != nil {
		start := f.offsetAt(r.Start)
		end := f.offsetAt(r.End)
		return string(f.content[start:end])
	}
	return string(f.content)
}

// PositionAt converts a zero-based offset to a position.
func (f *File) PositionAt(offset int) protocol.Position {
	offset = max(min(offset, len(f.content)), 0)
	lineOffsets := f.getLineOffsets()

	// lineOffsets always has at least one element (offset 0), so len(lineOffsets) == 0 is impossible
	// Handle the case where we only have one line (empty document or single line)
	if len(lineOffsets) == 1 {
		// Only one line, so we're on line 0
		offset = f.ensureBeforeEOL(offset, lineOffsets[0])
		return protocol.Position{
			Line:      0,
			Character: uint32(offset - lineOffsets[0]),
		}
	}

	// Binary search for the line
	low, high := 0, len(lineOffsets)
	for low < high {
		mid := (low + high) / 2
		if lineOffsets[mid] > offset {
			high = mid
		} else {
			low = mid + 1
		}
	}

	// low is the least x for which lineOffsets[x] > offset
	// So the line we want is low - 1
	line := low - 1

	// Ensure line is valid (should never be negative, but be defensive)
	if line < 0 {
		line = 0
	}

	offset = f.ensureBeforeEOL(offset, lineOffsets[line])
	return protocol.Position{
		Line:      uint32(line),
		Character: uint32(offset - lineOffsets[line]),
	}
}

// OffsetAt converts a position to a zero-based offset.
func (f *File) OffsetAt(position protocol.Position) int {
	return f.offsetAt(position)
}

// offsetAt is the internal implementation of OffsetAt.
func (f *File) offsetAt(position protocol.Position) int {
	lineOffsets := f.getLineOffsets()

	if int(position.Line) >= len(lineOffsets) {
		return len(f.content)
	}

	lineOffset := lineOffsets[position.Line]
	if position.Character == 0 {
		return lineOffset
	}

	var nextLineOffset int
	if int(position.Line+1) < len(lineOffsets) {
		nextLineOffset = lineOffsets[position.Line+1]
	} else {
		nextLineOffset = len(f.content)
	}

	offset := min(lineOffset+int(position.Character), nextLineOffset)
	return f.ensureBeforeEOL(offset, lineOffset)
}

// LineCount returns the number of lines in the document.
func (f *File) LineCount() int {
	return len(f.getLineOffsets())
}

// getLineOffsets returns the cached line start offsets.
func (f *File) getLineOffsets() []int {
	if f.lineOffsets != nil {
		return f.lineOffsets
	}
	// This should not happen if NewFile is used, but handle it defensively
	f.lineOffsets = computeLineOffsets(f.content, true, 0)
	return f.lineOffsets
}

// ensureBeforeEOL ensures the offset is before any end-of-line characters.
func (f *File) ensureBeforeEOL(offset, lineOffset int) int {
	for offset > lineOffset && isEOL(f.content[offset-1]) {
		offset--
	}
	return offset
}
