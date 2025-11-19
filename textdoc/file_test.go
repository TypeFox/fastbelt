// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
)

func TestNewFile(t *testing.T) {
	doc, err := NewFile("file:///test.txt", "plaintext", 1, "hello\nworld")
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
	content := doc.Content()
	if string(content) != "hello\nworld" {
		t.Errorf("Expected content 'hello\\nworld', got '%s'", string(content))
	}
}

func TestPositionAt(t *testing.T) {
	// Test with multi-line document
	doc, err := NewFile("file:///test.txt", "plaintext", 1, "ab\ncd")
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
			doc, err := NewFile("file:///test.txt", "plaintext", 1, tc.content)
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
	doc, err := NewFile("file:///test.txt", "plaintext", 1, "ab\ncd")
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
	doc, err := NewFile("file:///test.txt", "plaintext", 1, "hello\nworld")
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
		doc, err := NewFile("file:///test.txt", "plaintext", 1, test.content)
		if err != nil {
			t.Fatalf("New failed: %v", err)
		}
		lineCount := doc.LineCount()
		if lineCount != len(test.expected) {
			t.Errorf("Content '%s': expected %d lines, got %d", test.content, len(test.expected), lineCount)
		}
	}
}
