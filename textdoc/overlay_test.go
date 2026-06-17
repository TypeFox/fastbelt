// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"strings"
	"testing"

	"typefox.dev/lsp"
)

func TestUpdate(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Test incremental change
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 6},
				End:   lsp.Position{Line: 0, Character: 11},
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
	changes := []lsp.TextDocumentContentChangeEvent{
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

	edits := []lsp.TextEdit{
		{
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 6},
				End:   lsp.Position{Line: 0, Character: 11},
			},
			NewText: "Go",
		},
		{
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 0, Character: 5},
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

func TestErrorHandling(t *testing.T) {
	// Test version going backwards
	doc, err := NewOverlay("file:///test.txt", "plaintext", 5, "content")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	err = doc.Update([]lsp.TextDocumentContentChangeEvent{{Text: "new"}}, 3)
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 100},
				End:   lsp.Position{Line: 0, Character: 200},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 0, Character: 5},
			},
			Text: "goodbye",
		},
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 8}, // Position in "goodbye world"
				End:   lsp.Position{Line: 0, Character: 13},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 5, Character: 0},
				End:   lsp.Position{Line: 5, Character: 1},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 0, Character: 5},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 0, Character: 11},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 3},
				End:   lsp.Position{Line: 2, Character: 3},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 0, Character: 0},
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
	changes = []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 11},
				End:   lsp.Position{Line: 0, Character: 11},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 1, Character: 0},
				End:   lsp.Position{Line: 1, Character: 5},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 0, Character: 5},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 0, Character: 0},
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
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 0, Character: 11},
				End:   lsp.Position{Line: 0, Character: 6},
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

func TestUpdateUTF16(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		changes  []lsp.TextDocumentContentChangeEvent
		expected string
	}{
		{
			// é = U+00E9: 2 UTF-8 bytes, 1 UTF-16 code unit
			name:    "replace two-byte rune",
			content: "aéb",
			changes: []lsp.TextDocumentContentChangeEvent{
				{
					Range: &lsp.Range{
						Start: lsp.Position{Line: 0, Character: 1},
						End:   lsp.Position{Line: 0, Character: 2},
					},
					Text: "x",
				},
			},
			expected: "axb",
		},
		{
			// 😀 = U+1F600: 4 UTF-8 bytes, 2 UTF-16 code units (surrogate pair)
			name:    "replace surrogate-pair rune",
			content: "a😀b",
			changes: []lsp.TextDocumentContentChangeEvent{
				{
					Range: &lsp.Range{
						Start: lsp.Position{Line: 0, Character: 1},
						End:   lsp.Position{Line: 0, Character: 3},
					},
					Text: "x",
				},
			},
			expected: "axb",
		},
		{
			name:    "insert after two-byte rune",
			content: "aé\nbc",
			changes: []lsp.TextDocumentContentChangeEvent{
				{
					Range: &lsp.Range{
						Start: lsp.Position{Line: 0, Character: 2},
						End:   lsp.Position{Line: 0, Character: 2},
					},
					Text: "!",
				},
			},
			expected: "aé!\nbc",
		},
		{
			name:    "cross-line range with multi-byte rune",
			content: "aé\nbc",
			changes: []lsp.TextDocumentContentChangeEvent{
				{
					// Replace from after 'a' on line 0 to after 'b' on line 1
					Range: &lsp.Range{
						Start: lsp.Position{Line: 0, Character: 1},
						End:   lsp.Position{Line: 1, Character: 1},
					},
					Text: "X",
				},
			},
			expected: "aXc",
		},
		{
			name:    "replace surrogate-pair rune on second line",
			content: "ab\n😀x",
			changes: []lsp.TextDocumentContentChangeEvent{
				{
					Range: &lsp.Range{
						Start: lsp.Position{Line: 1, Character: 0},
						End:   lsp.Position{Line: 1, Character: 2},
					},
					Text: "y",
				},
			},
			expected: "ab\nyx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := NewOverlay("file:///test.txt", "plaintext", 1, tc.content)
			if err != nil {
				t.Fatalf("NewOverlay failed: %v", err)
			}
			err = doc.Update(tc.changes, 2)
			if err != nil {
				t.Fatalf("Update failed: %v", err)
			}
			if got := doc.Text(nil); got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestUpdateUTF16Validation(t *testing.T) {
	testCases := []struct {
		name        string
		content     string
		character   uint32
		shouldError bool
	}{
		{
			// "aé" has UTF-16 length 2; character 2 is valid (at end of line)
			name:        "character at end of two-byte rune line is valid",
			content:     "aé",
			character:   2,
			shouldError: false,
		},
		{
			// "aé" has UTF-16 length 2; character 3 exceeds the line
			name:        "character past end of two-byte rune line is invalid",
			content:     "aé",
			character:   3,
			shouldError: true,
		},
		{
			// "a😀" has UTF-16 length 3 (😀 = 2 units); character 3 is valid
			name:        "character at end of surrogate-pair line is valid",
			content:     "a😀",
			character:   3,
			shouldError: false,
		},
		{
			// "a😀" has UTF-16 length 3; character 4 exceeds the line
			name:        "character past end of surrogate-pair line is invalid",
			content:     "a😀",
			character:   4,
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := NewOverlay("file:///test.txt", "plaintext", 1, tc.content)
			if err != nil {
				t.Fatalf("NewOverlay failed: %v", err)
			}
			err = doc.Update([]lsp.TextDocumentContentChangeEvent{
				{
					Range: &lsp.Range{
						Start: lsp.Position{Line: 0, Character: tc.character},
						End:   lsp.Position{Line: 0, Character: tc.character},
					},
					Text: "",
				},
			}, 2)
			if tc.shouldError && err == nil {
				t.Error("expected error, got nil")
			} else if !tc.shouldError && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

func TestApplyEditsUTF16(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		edits    []lsp.TextEdit
		expected string
	}{
		{
			name:    "replace two-byte rune",
			content: "aéb",
			edits: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 0, Character: 1},
						End:   lsp.Position{Line: 0, Character: 2},
					},
					NewText: "x",
				},
			},
			expected: "axb",
		},
		{
			name:    "replace surrogate-pair rune",
			content: "a😀b",
			edits: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 0, Character: 1},
						End:   lsp.Position{Line: 0, Character: 3},
					},
					NewText: "x",
				},
			},
			expected: "axb",
		},
		{
			name:    "multiple edits with multi-byte runes",
			content: "aé\n😀b",
			edits: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 0, Character: 1},
						End:   lsp.Position{Line: 0, Character: 2},
					},
					NewText: "X",
				},
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 1, Character: 0},
						End:   lsp.Position{Line: 1, Character: 2},
					},
					NewText: "Y",
				},
			},
			expected: "aX\nYb",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := NewOverlay("file:///test.txt", "plaintext", 1, tc.content)
			if err != nil {
				t.Fatalf("NewOverlay failed: %v", err)
			}
			result, err := doc.ApplyEdits(tc.edits)
			if err != nil {
				t.Fatalf("ApplyEdits failed: %v", err)
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestUpdateConsecutiveUpdates(t *testing.T) {
	doc, err := NewOverlay("file:///test.txt", "plaintext", 1, "a\nb\nc")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// First update
	changes := []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 1, Character: 0},
				End:   lsp.Position{Line: 1, Character: 1},
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
	changes = []lsp.TextDocumentContentChangeEvent{
		{
			Range: &lsp.Range{
				Start: lsp.Position{Line: 2, Character: 0},
				End:   lsp.Position{Line: 2, Character: 1},
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
