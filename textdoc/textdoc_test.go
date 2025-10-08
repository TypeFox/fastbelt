// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"strings"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
)

func TestNewOverlay(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello\nworld")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if doc.URI() != "file:///test.txt" {
		t.Errorf("Expected URI 'file:///test.txt', got '%s'", doc.URI())
	}

	if doc.LanguageID() != "plaintext" {
		t.Errorf("Expected language ID 'plaintext', got '%s'", doc.LanguageID())
	}

	if doc.Version() != 1 {
		t.Errorf("Expected version 1, got %d", doc.Version())
	}

	if doc.Text(nil) != "hello\nworld" {
		t.Errorf("Expected content 'hello\\nworld', got '%s'", doc.Text(nil))
	}

	if doc.LineCount() != 2 {
		t.Errorf("Expected 2 lines, got %d", doc.LineCount())
	}

	// Test Content() method from Handle interface
	content, err := doc.Content()
	if err != nil {
		t.Errorf("Content() failed: %v", err)
	}
	if string(content) != "hello\nworld" {
		t.Errorf("Expected content 'hello\\nworld', got '%s'", string(content))
	}
}

func TestPositionAt(t *testing.T) {
	// Test with multi-line document
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "ab\ncd")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	tests := []struct {
		offset   int
		expected protocol.Position
	}{
		{0, protocol.Position{Line: 0, Character: 0}},
		{1, protocol.Position{Line: 0, Character: 1}},
		{2, protocol.Position{Line: 0, Character: 2}}, // At \n
		{3, protocol.Position{Line: 1, Character: 0}},
		{4, protocol.Position{Line: 1, Character: 1}},
		{5, protocol.Position{Line: 1, Character: 2}}, // Beyond end
	}

	for _, test := range tests {
		pos := doc.PositionAt(test.offset)
		if pos.Line != test.expected.Line || pos.Character != test.expected.Character {
			t.Errorf("PositionAt(%d): expected {%d, %d}, got {%d, %d}",
				test.offset, test.expected.Line, test.expected.Character, pos.Line, pos.Character)
		}
	}
}

func TestPositionAtEdgeCases(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		tests   []struct {
			offset   int
			expected protocol.Position
		}
	}{
		{
			name:    "empty document",
			content: "",
			tests: []struct {
				offset   int
				expected protocol.Position
			}{
				{0, protocol.Position{Line: 0, Character: 0}},
				{1, protocol.Position{Line: 0, Character: 0}},  // Beyond end, should clamp
				{-1, protocol.Position{Line: 0, Character: 0}}, // Negative, should clamp
			},
		},
		{
			name:    "single character",
			content: "a",
			tests: []struct {
				offset   int
				expected protocol.Position
			}{
				{0, protocol.Position{Line: 0, Character: 0}},
				{1, protocol.Position{Line: 0, Character: 1}}, // At end
				{2, protocol.Position{Line: 0, Character: 1}}, // Beyond end, should clamp
			},
		},
		{
			name:    "single line with newline",
			content: "hello\n",
			tests: []struct {
				offset   int
				expected protocol.Position
			}{
				{0, protocol.Position{Line: 0, Character: 0}},
				{5, protocol.Position{Line: 0, Character: 5}}, // At \n
				{6, protocol.Position{Line: 1, Character: 0}}, // After \n
				{7, protocol.Position{Line: 1, Character: 0}}, // Beyond end
			},
		},
		{
			name:    "windows line endings",
			content: "a\r\nb",
			tests: []struct {
				offset   int
				expected protocol.Position
			}{
				{0, protocol.Position{Line: 0, Character: 0}},
				{1, protocol.Position{Line: 0, Character: 1}}, // At \r
				{2, protocol.Position{Line: 0, Character: 1}}, // At \n (should be before EOL)
				{3, protocol.Position{Line: 1, Character: 0}}, // After \r\n
				{4, protocol.Position{Line: 1, Character: 1}}, // At 'b'
			},
		},
		{
			name:    "multiple empty lines",
			content: "\n\n\n",
			tests: []struct {
				offset   int
				expected protocol.Position
			}{
				{0, protocol.Position{Line: 0, Character: 0}}, // At first \n
				{1, protocol.Position{Line: 1, Character: 0}}, // At second \n
				{2, protocol.Position{Line: 2, Character: 0}}, // At third \n
				{3, protocol.Position{Line: 3, Character: 0}}, // After all \n
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := NewOverlay("file:///test.txt", "plaintext", 1, tc.content)
			if err != nil {
				t.Fatalf("New failed: %v", err)
			}

			for _, test := range tc.tests {
				pos := doc.PositionAt(test.offset)
				if pos.Line != test.expected.Line || pos.Character != test.expected.Character {
					t.Errorf("PositionAt(%d): expected {%d, %d}, got {%d, %d}",
						test.offset, test.expected.Line, test.expected.Character, pos.Line, pos.Character)
				}
			}
		})
	}
}

func TestOffsetAt(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "ab\ncd")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	tests := []struct {
		position protocol.Position
		expected int
	}{
		{protocol.Position{Line: 0, Character: 0}, 0},
		{protocol.Position{Line: 0, Character: 1}, 1},
		{protocol.Position{Line: 0, Character: 2}, 2},
		{protocol.Position{Line: 1, Character: 0}, 3},
		{protocol.Position{Line: 1, Character: 1}, 4},
		{protocol.Position{Line: 1, Character: 2}, 5},
		{protocol.Position{Line: 2, Character: 0}, 5}, // Beyond end
	}

	for _, test := range tests {
		offset := doc.OffsetAt(test.position)
		if offset != test.expected {
			t.Errorf("OffsetAt({%d, %d}): expected %d, got %d",
				test.position.Line, test.position.Character, test.expected, offset)
		}
	}
}

func TestTextWithRange(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello\nworld")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Get substring
	r := &protocol.Range{
		Start: protocol.Position{Line: 0, Character: 1},
		End:   protocol.Position{Line: 1, Character: 2},
	}

	text := doc.Text(r)
	expected := "ello\nwo"
	if text != expected {
		t.Errorf("Expected '%s', got '%s'", expected, text)
	}
}

func TestUpdate(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Test incremental change
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 6},
				End:   protocol.Position{Line: 0, Character: 11},
			},
			Text: "Go",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if doc.Text(nil) != "hello Go" {
		t.Errorf("Expected 'hello Go', got '%s'", doc.Text(nil))
	}

	if doc.Version() != 2 {
		t.Errorf("Expected version 2, got %d", doc.Version())
	}
}

func TestUpdateFullDocument(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Test full document change
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Text: "new content",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if doc.Text(nil) != "new content" {
		t.Errorf("Expected 'new content', got '%s'", doc.Text(nil))
	}
}

func TestApplyEdits(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	edits := []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 6},
				End:   protocol.Position{Line: 0, Character: 11},
			},
			NewText: "Go",
		},
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 5},
			},
			NewText: "hi",
		},
	}

	result, err := doc.ApplyEdits(edits)
	if err != nil {
		t.Errorf("ApplyEdits failed: %v", err)
	}

	expected := "hi Go"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestLineOffsets(t *testing.T) {
	tests := []struct {
		content  string
		expected []int
	}{
		{"", []int{0}},
		{"a", []int{0}},
		{"a\n", []int{0, 2}},
		{"a\nb", []int{0, 2}},
		{"a\r\nb", []int{0, 3}},
		{"a\n\nb", []int{0, 2, 3}},
	}

	for _, test := range tests {
		doc, err := NewOverlay("file:///test.txt", "plaintext", 1, test.content)
		if err != nil {
			t.Fatalf("New failed: %v", err)
		}
		lineCount := doc.LineCount()
		if lineCount != len(test.expected) {
			t.Errorf("Content '%s': expected %d lines, got %d", test.content, len(test.expected), lineCount)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	// Test version going backwards
	doc, err := NewOverlay("file:///test.txt", "plaintext", 5, "content")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	err = doc.Update([]protocol.TextDocumentContentChangeEvent{{Text: "new"}}, 3)
	if err == nil {
		t.Error("Expected error for version going backwards")
	}
	if !strings.Contains(err.Error(), "version") {
		t.Errorf("Expected version error, got: %v", err)
	}
}

func TestUpdateErrorMessageFormat(t *testing.T) {
	// Test that error messages are properly formatted with package prefix
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Test with invalid range to trigger an error
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 100},
				End:   protocol.Position{Line: 0, Character: 200},
			},
			Text: "test",
		},
	}

	err = doc.Update(changes, 2)
	if err == nil {
		t.Error("Expected error for invalid range")
	}

	// Verify error has package prefix
	if !strings.Contains(err.Error(), "textdoc:") {
		t.Errorf("Expected error to contain 'textdoc:' prefix, got: %v", err)
	}
}

func TestUpdateMultipleChanges(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Apply multiple changes in one update.
	// Important: Changes are applied sequentially, so each change's positions
	// must account for the effects of previous changes.
	// Document starts as: "hello world" (11 chars)
	// After first change: "goodbye world" (13 chars)
	// After second change: "goodbye Go" (10 chars)
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 5},
			},
			Text: "goodbye",
		},
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 8}, // Position in "goodbye world"
				End:   protocol.Position{Line: 0, Character: 13},
			},
			Text: "Go",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update with multiple changes failed: %v", err)
	}

	expected := "goodbye Go"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestUpdateInvalidLine(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Try to change a line that doesn't exist
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 5, Character: 0},
				End:   protocol.Position{Line: 5, Character: 1},
			},
			Text: "test",
		},
	}

	err = doc.Update(changes, 2)
	if err == nil {
		t.Error("Expected error for invalid line number")
	}
	if !strings.Contains(err.Error(), "line out of bounds") {
		t.Errorf("Expected 'line out of bounds' error, got: %v", err)
	}
}

func TestUpdateEmptyRangeInsertion(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Insert text without replacing anything (empty range)
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 5},
				End:   protocol.Position{Line: 0, Character: 5},
			},
			Text: " world",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update with empty range insertion failed: %v", err)
	}

	expected := "hello world"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestUpdateDeletion(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Delete text (empty replacement)
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 5},
				End:   protocol.Position{Line: 0, Character: 11},
			},
			Text: "",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update with deletion failed: %v", err)
	}

	expected := "hello"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestUpdateMultiLine(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "line1\nline2\nline3")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Replace text spanning multiple lines
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 3},
				End:   protocol.Position{Line: 2, Character: 3},
			},
			Text: "A\nB",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update with multi-line range failed: %v", err)
	}

	expected := "linA\nBe3"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestUpdateAtDocumentBoundaries(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Insert at start
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Text: "start ",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update at document start failed: %v", err)
	}

	if doc.Text(nil) != "start hello" {
		t.Errorf("Expected 'start hello', got '%s'", doc.Text(nil))
	}

	// Insert at end
	changes = []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 11},
				End:   protocol.Position{Line: 0, Character: 11},
			},
			Text: " end",
		},
	}

	err = doc.Update(changes, 3)
	if err != nil {
		t.Errorf("Update at document end failed: %v", err)
	}

	if doc.Text(nil) != "start hello end" {
		t.Errorf("Expected 'start hello end', got '%s'", doc.Text(nil))
	}
}

func TestUpdateWithWindowsLineEndings(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "line1\r\nline2\r\nline3")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Replace text on second line
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 1, Character: 0},
				End:   protocol.Position{Line: 1, Character: 5},
			},
			Text: "modified",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update with Windows line endings failed: %v", err)
	}

	expected := "line1\r\nmodified\r\nline3"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}

	// Verify line count is correct (Windows line endings should be treated as single line breaks)
	if doc.LineCount() != 3 {
		t.Errorf("Expected 3 lines, got %d", doc.LineCount())
	}
}

func TestUpdatePositionAtEndOfLine(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello\nworld")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Insert at end of first line (at the newline position)
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 5},
				End:   protocol.Position{Line: 0, Character: 5},
			},
			Text: "!",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update at end of line failed: %v", err)
	}

	expected := "hello!\nworld"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestUpdateEmptyDocument(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Insert into empty document
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Text: "hello",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update on empty document failed: %v", err)
	}

	if doc.Text(nil) != "hello" {
		t.Errorf("Expected 'hello', got '%s'", doc.Text(nil))
	}
}

func TestUpdateBackwardsRange(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Provide a backwards range (end before start) - should be normalized by getWellFormedRange
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 0, Character: 11},
				End:   protocol.Position{Line: 0, Character: 6},
			},
			Text: "Go",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("Update with backwards range failed: %v", err)
	}

	expected := "hello Go"
	if doc.Text(nil) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestUpdateConsecutiveUpdates(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "a\nb\nc")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// First update
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 1, Character: 0},
				End:   protocol.Position{Line: 1, Character: 1},
			},
			Text: "modified",
		},
	}

	err = doc.Update(changes, 2)
	if err != nil {
		t.Errorf("First update failed: %v", err)
	}

	if doc.Text(nil) != "a\nmodified\nc" {
		t.Errorf("After first update: expected 'a\\nmodified\\nc', got '%s'", doc.Text(nil))
	}

	// Second update on line 2 (which should still be line 2 after cache invalidation)
	changes = []protocol.TextDocumentContentChangeEvent{
		{
			Range: &protocol.Range{
				Start: protocol.Position{Line: 2, Character: 0},
				End:   protocol.Position{Line: 2, Character: 1},
			},
			Text: "changed",
		},
	}

	err = doc.Update(changes, 3)
	if err != nil {
		t.Errorf("Second update failed: %v", err)
	}

	if doc.Text(nil) != "a\nmodified\nchanged" {
		t.Errorf("After second update: expected 'a\\nmodified\\nchanged', got '%s'", doc.Text(nil))
	}
}
