// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"typefox.dev/lsp"
)

// Overlay represents an open text document in the editor.
// It may have unsaved edits and implements both Handle and Mapper.
// Conceptually, an overlay models in-memory edits layered on top of any on-disk file content.
type Overlay struct {
	File
	mu sync.RWMutex
}

// NewOverlay creates a new text document overlay.
// The content parameter is accepted as a string and stored internally as bytes.
func NewOverlay(uri lsp.DocumentURI, languageID string, version int32, content string) (*Overlay, error) {
	file, err := NewFile(uri, languageID, version, content)
	if err != nil {
		return nil, err
	}
	return &Overlay{
		File: *file,
	}, nil
}

// Update applies content changes to the overlay and updates its version.
// Both incremental and full-document change events are supported. This method is thread-safe.
func (o *Overlay) Update(changes []lsp.TextDocumentContentChangeEvent, version int32) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if version < o.version {
		return fmt.Errorf("textdoc: version %d < current version %d", version, o.version)
	}

	for i, change := range changes {
		if err := o.applyChangeLocked(change); err != nil {
			return fmt.Errorf("textdoc: change %d: %w", i, err)
		}
	}
	o.version = version
	return nil
}

// ApplyEdits applies a list of text edits and returns the resulting text.
// The edits are automatically sorted by position, and overlapping edits return
// an error. This method is thread-safe.
func (o *Overlay) ApplyEdits(edits []lsp.TextEdit) (string, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	content := o.content

	// Sort edits by position (start line, then start character)
	sortedEdits := make([]lsp.TextEdit, len(edits))
	copy(sortedEdits, edits)
	sort.Slice(sortedEdits, func(i, j int) bool {
		a, b := sortedEdits[i], sortedEdits[j]
		if a.Range.Start.Line != b.Range.Start.Line {
			return a.Range.Start.Line < b.Range.Start.Line
		}
		return a.Range.Start.Character < b.Range.Start.Character
	})

	// Pre-calculate approximate result size to reduce allocations
	estimatedLen := len(content)
	for _, edit := range sortedEdits {
		wellFormedEdit := getWellFormedEdit(edit)
		startOffset := o.offsetAtLocked(wellFormedEdit.Range.Start)
		endOffset := o.offsetAtLocked(wellFormedEdit.Range.End)
		estimatedLen += len(wellFormedEdit.NewText) - (endOffset - startOffset)
	}

	var result strings.Builder
	result.Grow(estimatedLen)

	lastModifiedOffset := 0

	for _, edit := range sortedEdits {
		wellFormedEdit := getWellFormedEdit(edit)
		startOffset := o.offsetAtLocked(wellFormedEdit.Range.Start)

		if startOffset < lastModifiedOffset {
			return "", errors.New("textdoc: overlapping edit")
		}

		if startOffset > lastModifiedOffset {
			result.Write(content[lastModifiedOffset:startOffset])
		}

		if len(wellFormedEdit.NewText) > 0 {
			result.WriteString(wellFormedEdit.NewText)
		}

		lastModifiedOffset = o.offsetAtLocked(wellFormedEdit.Range.End)
	}

	result.Write(content[lastModifiedOffset:])
	return result.String(), nil
}

// Version returns the current document version.
// This method is thread-safe.
func (o *Overlay) Version() int32 {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.version
}

// Content returns a copy of the document content as a byte slice.
// Returning bytes enables efficient manipulation; a copy is returned to prevent
// external modification of internal state. This method is thread-safe.
func (o *Overlay) Content() []byte {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.File.Content()
}

// Text returns the text content or a substring if range is provided.
// This is a convenience method that returns a string instead of []byte.
// This method is thread-safe.
func (o *Overlay) Text(r *lsp.Range) string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.File.Text(r)
}

// PositionAt converts a zero-based offset to a position.
// This method is thread-safe.
func (o *Overlay) PositionAt(offset int) lsp.Position {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.File.PositionAt(offset)
}

// OffsetAt converts a position to a zero-based offset.
// This method is thread-safe.
func (o *Overlay) OffsetAt(position lsp.Position) int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.File.OffsetAt(position)
}

// offsetAtLocked is the internal implementation of OffsetAt that assumes the lock is held.
func (o *Overlay) offsetAtLocked(position lsp.Position) int {
	return o.offsetAt(position)
}

// LineCount returns the number of lines in the document.
// This method is thread-safe.
func (o *Overlay) LineCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.File.LineCount()
}

// applyChangeLocked applies a single content change to the overlay.
// This method assumes the write lock is already held.
func (o *Overlay) applyChangeLocked(change lsp.TextDocumentContentChangeEvent) error {
	if change.Range != nil {
		// Incremental change
		wellFormedRange := getWellFormedRange(*change.Range)

		// Validate that line numbers are within document bounds
		lineOffsets := o.getLineOffsetsLocked()
		if int(wellFormedRange.Start.Line) >= len(lineOffsets) || int(wellFormedRange.End.Line) >= len(lineOffsets) {
			return errors.New("textdoc: invalid range: line out of bounds")
		}

		// Validate that character positions are reasonable for their respective lines.
		// We check against the maximum possible character position for each line.
		// Note: offsetAtLocked will clamp positions and handle EOL positioning,
		// but we want to reject positions that are clearly out of bounds.
		for _, pos := range []lsp.Position{wellFormedRange.Start, wellFormedRange.End} {
			lineStart := lineOffsets[pos.Line]
			var lineEnd int
			if int(pos.Line+1) < len(lineOffsets) {
				lineEnd = lineOffsets[pos.Line+1]
			} else {
				lineEnd = len(o.content)
			}
			maxChar := lineEnd - lineStart
			if int(pos.Character) > maxChar {
				return fmt.Errorf("textdoc: invalid range: character %d exceeds line %d length %d",
					pos.Character, pos.Line, maxChar)
			}
		}

		startOffset := o.offsetAtLocked(wellFormedRange.Start)
		endOffset := o.offsetAtLocked(wellFormedRange.End)

		// Update content using []byte operations for efficiency
		newContent := make([]byte, 0, startOffset+len(change.Text)+(len(o.content)-endOffset))
		newContent = append(newContent, o.content[:startOffset]...)
		newContent = append(newContent, change.Text...)
		newContent = append(newContent, o.content[endOffset:]...)
		o.content = newContent

		// Invalidate line offsets cache
		o.lineOffsets = nil
	} else {
		// Full document change
		o.content = []byte(change.Text)
		o.lineOffsets = nil
	}
	return nil
}

// getLineOffsetsLocked computes and caches line start offsets for the current
// content. This method assumes the lock is already held (at least for reading),
// and fills the cache on first use.
func (o *Overlay) getLineOffsetsLocked() []int {
	if o.lineOffsets != nil {
		return o.lineOffsets
	}

	// Need to compute line offsets
	// Note: This is called from methods that already hold at least a read lock.
	// For simplicity, we compute inline. A more sophisticated approach would
	// upgrade to a write lock, but that requires releasing the read lock first.
	o.lineOffsets = computeLineOffsets(o.content, true, 0)
	return o.lineOffsets
}

// computeLineOffsets computes line offsets for the given content.
func computeLineOffsets(content []byte, isAtLineStart bool, textOffset int) []int {
	var result []int
	if isAtLineStart {
		result = append(result, textOffset)
	}

	for i := 0; i < len(content); i++ {
		ch := content[i]
		if isEOL(ch) {
			if ch == '\r' && i+1 < len(content) && content[i+1] == '\n' {
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
func getWellFormedRange(r lsp.Range) lsp.Range {
	start, end := r.Start, r.End
	if start.Line > end.Line || (start.Line == end.Line && start.Character > end.Character) {
		return lsp.Range{Start: end, End: start}
	}
	return r
}

// getWellFormedEdit ensures the edit has a well-formed range
func getWellFormedEdit(edit lsp.TextEdit) lsp.TextEdit {
	wellFormedRange := getWellFormedRange(edit.Range)
	if wellFormedRange.Start != edit.Range.Start || wellFormedRange.End != edit.Range.End {
		return lsp.TextEdit{
			Range:   wellFormedRange,
			NewText: edit.NewText,
		}
	}
	return edit
}
