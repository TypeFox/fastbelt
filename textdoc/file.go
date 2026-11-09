// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"sync"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"

	"typefox.dev/lsp"
)

// File represents an immutable text document read from the file system.
// It implements Handle but does not support modifications.
// Its methods are thread-safe.
type File struct {
	uri        lsp.DocumentURI
	languageID string

	// mu guards the fields below. A File computes lineOffsets on first use:
	// many documents are only lexed and parsed, and never need positions.
	// Overlay, which embeds File, also changes version and content.
	mu          sync.RWMutex
	version     int32
	content     []byte
	lineOffsets []int // nil until first use
}

// NewFile creates a new text document file.
// The content parameter is accepted as a string and stored internally as bytes.
func NewFile(uri lsp.DocumentURI, languageID string, version int32, content string) *File {
	f := &File{}
	f.init(uri, languageID, version, content)
	return f
}

// init sets up f in place. File must not be copied (it holds a mutex), so
// Overlay embeds it initialized through init.
func (f *File) init(uri lsp.DocumentURI, languageID string, version int32, content string) {
	f.uri = uri
	f.languageID = languageID
	f.version = version
	// Convert in-place from string to []byte without copying.
	f.content = unsafe.Slice(unsafe.StringData(content), len(content))
}

// URI returns the document URI.
func (f *File) URI() lsp.DocumentURI {
	return f.uri
}

// LanguageID returns the language identifier.
func (f *File) LanguageID() string {
	return f.languageID
}

// Version returns the document version.
func (f *File) Version() int32 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.version
}

// Content returns the byte slice that represents the document content.
// This is not a copy. Do not modify the returned slice, as it may lead
// to unexpected behavior.
func (f *File) Content() []byte {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.content
}

// Text returns the text content or a substring if range is provided.
// This is a convenience method that returns a string instead of []byte.
//
// Note: This string is not a copy, and holding onto it will prevent
// the underlying byte slice from being garbage collected.
func (f *File) Text(r *lsp.Range) string {
	var content []byte
	if r == nil {
		content = f.Content()
	} else {
		var lineOffsets []int
		content, lineOffsets = f.lines()
		start := offsetAt(content, lineOffsets, r.Start)
		end := offsetAt(content, lineOffsets, r.End)
		content = content[start:end]
	}
	// We are creating an "unsafe" view on the byte slice to avoid copying it.
	// This is safe, since the content is never actually modified after creation.
	return unsafe.String(unsafe.SliceData(content), len(content))
}

// PositionAt converts a zero-based byte-offset to a UTF16 position.
func (f *File) PositionAt(offset int) lsp.Position {
	content, lineOffsets := f.lines()
	return positionAt(content, lineOffsets, offset)
}

// OffsetAt converts a UTF16 position to a zero-based byte-offset.
func (f *File) OffsetAt(position lsp.Position) int {
	content, lineOffsets := f.lines()
	return offsetAt(content, lineOffsets, position)
}

// LineCount returns the number of lines in the document.
func (f *File) LineCount() int {
	_, lineOffsets := f.lines()
	return len(lineOffsets)
}

// lines returns the content and its line start offsets, computing the
// offsets on first use. The pair stays consistent after the lock is
// released, since content is never modified in place.
func (f *File) lines() ([]byte, []int) {
	f.mu.RLock()
	content, lineOffsets := f.content, f.lineOffsets
	f.mu.RUnlock()
	if lineOffsets != nil {
		return content, lineOffsets
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.content, f.lineOffsetsLocked()
}

// lineOffsetsLocked returns the line start offsets of the current content,
// computing them if needed. The caller must hold the write lock.
func (f *File) lineOffsetsLocked() []int {
	if f.lineOffsets == nil {
		f.lineOffsets = computeLineOffsets(f.content, true, 0)
	}
	return f.lineOffsets
}

// positionAt converts a zero-based byte-offset in content to a UTF16 position.
func positionAt(content []byte, lineOffsets []int, offset int) lsp.Position {
	offset = max(min(offset, len(content)), 0)

	// lineOffsets always has at least one element (offset 0), so len(lineOffsets) == 0 is impossible
	// Handle the case where we only have one line (empty document or single line)
	if len(lineOffsets) == 1 {
		// Only one line, so we're on line 0
		offset = ensureBeforeEOL(content, offset, lineOffsets[0])
		return toUTF16Position(content, 0, offset, lineOffsets)
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
	line := max(low-1, 0)

	offset = ensureBeforeEOL(content, offset, lineOffsets[line])
	return toUTF16Position(content, line, offset, lineOffsets)
}

func toUTF16Position(content []byte, line, offset int, offsets []int) lsp.Position {
	lineOffset := offsets[line]
	return lsp.Position{
		Line:      uint32(line),
		Character: uint32(utf16LineLen(content, lineOffset, offset)),
	}
}

// offsetAt converts a UTF16 position to a zero-based byte-offset in content.
func offsetAt(content []byte, lineOffsets []int, position lsp.Position) int {
	if int(position.Line) >= len(lineOffsets) {
		return len(content)
	}

	lineOffset := lineOffsets[position.Line]
	if position.Character == 0 {
		return lineOffset
	}

	var nextLineOffset int
	if int(position.Line+1) < len(lineOffsets) {
		nextLineOffset = lineOffsets[position.Line+1]
	} else {
		nextLineOffset = len(content)
	}

	character := int(position.Character)
	offset := lineOffset
	for i := 0; i < character && offset < nextLineOffset; {
		r, size := utf8.DecodeRune(content[offset:nextLineOffset])
		offset += size
		i += utf16.RuneLen(r)
	}
	return ensureBeforeEOL(content, offset, lineOffset)
}

// ensureBeforeEOL ensures the offset is before any end-of-line characters.
func ensureBeforeEOL(content []byte, offset, lineOffset int) int {
	for offset > lineOffset && isEOL(content[offset-1]) {
		offset--
	}
	return offset
}

// utf16LineLen returns the number of UTF-16 code units in content[start:end].
func utf16LineLen(content []byte, start, end int) int {
	count := 0
	for i := start; i < end; {
		r, size := utf8.DecodeRune(content[i:end])
		count += utf16.RuneLen(r)
		i += size
	}
	return count
}
