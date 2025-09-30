// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/TypeFox/go-lsp/protocol"
)



// TextDocument represents a simple text document that keeps content as string
type TextDocument interface {
	// URI returns the associated URI for this document
	URI() protocol.DocumentURI
	// LanguageID returns the identifier of the language associated with this document
	LanguageID() string
	// Version returns the version number of this document
	Version() int32
	// GetText returns the text of this document or a substring if range is provided
	GetText(r *protocol.Range) string
	// PositionAt converts a zero-based offset to a position
	PositionAt(offset int) protocol.Position
	// OffsetAt converts a position to a zero-based offset
	OffsetAt(position protocol.Position) int
	// LineCount returns the number of lines in this document
	LineCount() int
}

// fullTextDocument implements TextDocument interface
type fullTextDocument struct {
	uri         protocol.DocumentURI
	languageID  string
	version     int32
	content     string
	lineOffsets []int
}

// Create creates a new text document
func Create(uri protocol.DocumentURI, languageID string, version int32, content string) TextDocument {
	if uri == "" {
		panic("uri cannot be empty")
	}
	if languageID == "" {
		panic("languageID cannot be empty")
	}
	return &fullTextDocument{
		uri:        uri,
		languageID: languageID,
		version:    version,
		content:    content,
	}
}

// Update updates a TextDocument by modifying its content
func Update(document TextDocument, changes []protocol.TextDocumentContentChangeEvent, version int32) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}
	
	doc, ok := document.(*fullTextDocument)
	if !ok {
		return errors.New("document must be created by Create function")
	}

	if version < doc.version {
		return errors.New("version cannot go backwards")
	}

	for i, change := range changes {
		if err := doc.applyChange(change); err != nil {
			return errors.New("failed to apply change " + strconv.Itoa(i) + ": " + err.Error())
		}
	}
	doc.version = version
	return nil
}

// ApplyEdits applies a list of text edits to a document and returns the resulting text
func ApplyEdits(document TextDocument, edits []protocol.TextEdit) (string, error) {
	if document == nil {
		return "", errors.New("document cannot be nil")
	}
	
	text := document.GetText(nil)
	
	// Sort edits by position (start line, then start character)
	sortedEdits := make([]protocol.TextEdit, len(edits))
	copy(sortedEdits, edits)
	sort.Slice(sortedEdits, func(i, j int) bool {
		a, b := sortedEdits[i], sortedEdits[j]
		if a.Range.Start.Line != b.Range.Start.Line {
			return a.Range.Start.Line < b.Range.Start.Line
		}
		return a.Range.Start.Character < b.Range.Start.Character
	})

	var spans []string
	lastModifiedOffset := 0

	for _, edit := range sortedEdits {
		wellFormedEdit := getWellFormedEdit(edit)
		startOffset := document.OffsetAt(wellFormedEdit.Range.Start)
		
		if startOffset < lastModifiedOffset {
			return "", errors.New("overlapping edit")
		}
		
		if startOffset > lastModifiedOffset {
			spans = append(spans, text[lastModifiedOffset:startOffset])
		}
		
		if len(wellFormedEdit.NewText) > 0 {
			spans = append(spans, wellFormedEdit.NewText)
		}
		
		lastModifiedOffset = document.OffsetAt(wellFormedEdit.Range.End)
	}
	
	spans = append(spans, text[lastModifiedOffset:])
	return strings.Join(spans, ""), nil
}

// URI returns the document URI
func (d *fullTextDocument) URI() protocol.DocumentURI {
	return d.uri
}

// LanguageID returns the language identifier
func (d *fullTextDocument) LanguageID() string {
	return d.languageID
}

// Version returns the document version
func (d *fullTextDocument) Version() int32 {
	return d.version
}

// GetText returns the text content or a substring if range is provided
func (d *fullTextDocument) GetText(r *protocol.Range) string {
	if r != nil {
		start := d.OffsetAt(r.Start)
		end := d.OffsetAt(r.End)
		return d.content[start:end]
	}
	return d.content
}

// PositionAt converts a zero-based offset to a position
func (d *fullTextDocument) PositionAt(offset int) protocol.Position {
	offset = max(min(offset, len(d.content)), 0)
	lineOffsets := d.getLineOffsets()
	
	if len(lineOffsets) == 0 {
		return protocol.Position{Line: 0, Character: uint32(offset)}
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
	
	line := low - 1
	offset = d.ensureBeforeEOL(offset, lineOffsets[line])
	return protocol.Position{
		Line:      uint32(line),
		Character: uint32(offset - lineOffsets[line]),
	}
}

// OffsetAt converts a position to a zero-based offset
func (d *fullTextDocument) OffsetAt(position protocol.Position) int {
	lineOffsets := d.getLineOffsets()
	
	if int(position.Line) >= len(lineOffsets) {
		return len(d.content)
	}
	
	lineOffset := lineOffsets[position.Line]
	if position.Character == 0 {
		return lineOffset
	}
	
	var nextLineOffset int
	if int(position.Line+1) < len(lineOffsets) {
		nextLineOffset = lineOffsets[position.Line+1]
	} else {
		nextLineOffset = len(d.content)
	}
	
	offset := min(lineOffset+int(position.Character), nextLineOffset)
	return d.ensureBeforeEOL(offset, lineOffset)
}

// LineCount returns the number of lines in the document
func (d *fullTextDocument) LineCount() int {
	return len(d.getLineOffsets())
}

// applyChange applies a single content change to the document
func (d *fullTextDocument) applyChange(change protocol.TextDocumentContentChangeEvent) error {
	if change.Range != nil {
		// Incremental change
		wellFormedRange := getWellFormedRange(*change.Range)
		startOffset := d.OffsetAt(wellFormedRange.Start)
		endOffset := d.OffsetAt(wellFormedRange.End)
		
		// Validate offsets
		if startOffset < 0 || endOffset < 0 || startOffset > len(d.content) || endOffset > len(d.content) {
			return errors.New("invalid range: offsets out of bounds")
		}
		if startOffset > endOffset {
			return errors.New("invalid range: start offset greater than end offset")
		}
		
		// Update content
		d.content = d.content[:startOffset] + change.Text + d.content[endOffset:]
		
		// Invalidate line offsets cache
		d.lineOffsets = nil
	} else {
		// Full document change
		d.content = change.Text
		d.lineOffsets = nil
	}
	return nil
}

// getLineOffsets computes and caches line offsets
func (d *fullTextDocument) getLineOffsets() []int {
	if d.lineOffsets == nil {
		d.lineOffsets = computeLineOffsets(d.content, true, 0)
	}
	return d.lineOffsets
}

// ensureBeforeEOL ensures the offset is before any end-of-line characters
func (d *fullTextDocument) ensureBeforeEOL(offset, lineOffset int) int {
	for offset > lineOffset && isEOL(d.content[offset-1]) {
		offset--
	}
	return offset
}

// computeLineOffsets computes line offsets for the given text
func computeLineOffsets(text string, isAtLineStart bool, textOffset int) []int {
	var result []int
	if isAtLineStart {
		result = append(result, textOffset)
	}
	
	for i := 0; i < len(text); i++ {
		ch := text[i]
		if isEOL(ch) {
			if ch == '\r' && i+1 < len(text) && text[i+1] == '\n' {
				i++ // Skip the \n in \r\n
			}
			result = append(result, textOffset+i+1)
		}
	}
	
	return result
}

// isEOL checks if a character is an end-of-line character
func isEOL(ch byte) bool {
	return ch == '\r' || ch == '\n'
}

// getWellFormedRange ensures start is before end in a range
func getWellFormedRange(r protocol.Range) protocol.Range {
	start, end := r.Start, r.End
	if start.Line > end.Line || (start.Line == end.Line && start.Character > end.Character) {
		return protocol.Range{Start: end, End: start}
	}
	return r
}

// getWellFormedEdit ensures the edit has a well-formed range
func getWellFormedEdit(edit protocol.TextEdit) protocol.TextEdit {
	wellFormedRange := getWellFormedRange(edit.Range)
	if wellFormedRange.Start != edit.Range.Start || wellFormedRange.End != edit.Range.End {
		return protocol.TextEdit{
			Range:   wellFormedRange,
			NewText: edit.NewText,
		}
	}
	return edit
}

// Helper functions for min/max since Go 1.21+ has these in the standard library
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}