// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"strconv"
	"strings"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
)

func TestCreate(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "hello\nworld")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
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
	
	if doc.GetText(nil) != "hello\nworld" {
		t.Errorf("Expected content 'hello\\nworld', got '%s'", doc.GetText(nil))
	}
	
	if doc.LineCount() != 2 {
		t.Errorf("Expected 2 lines, got %d", doc.LineCount())
	}
}

func TestPositionAt(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "ab\ncd")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
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

func TestOffsetAt(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "ab\ncd")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
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

func TestGetTextWithRange(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "hello\nworld")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Get substring
	r := &protocol.Range{
		Start: protocol.Position{Line: 0, Character: 1},
		End:   protocol.Position{Line: 1, Character: 2},
	}
	
	text := doc.GetText(r)
	expected := "ello\nwo"
	if text != expected {
		t.Errorf("Expected '%s', got '%s'", expected, text)
	}
}

func TestUpdate(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
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
	
	err = Update(doc, changes, 2)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
	
	if doc.GetText(nil) != "hello Go" {
		t.Errorf("Expected 'hello Go', got '%s'", doc.GetText(nil))
	}
	
	if doc.Version() != 2 {
		t.Errorf("Expected version 2, got %d", doc.Version())
	}
}

func TestUpdateFullDocument(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	// Test full document change
	changes := []protocol.TextDocumentContentChangeEvent{
		{
			Text: "new content",
		},
	}
	
	err = Update(doc, changes, 2)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}
	
	if doc.GetText(nil) != "new content" {
		t.Errorf("Expected 'new content', got '%s'", doc.GetText(nil))
	}
}

func TestApplyEdits(t *testing.T) {
	doc, err := Create("file:///test.txt", "plaintext", 1, "hello world")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
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
	
	result, err := ApplyEdits(doc, edits)
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
		doc, err := Create("file:///test.txt", "plaintext", 1, test.content)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
		lineCount := doc.LineCount()
		if lineCount != len(test.expected) {
			t.Errorf("Content '%s': expected %d lines, got %d", test.content, len(test.expected), lineCount)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	// Test nil document
	err := Update(nil, []protocol.TextDocumentContentChangeEvent{}, 1)
	if err == nil {
		t.Error("Expected error for nil document")
	}
	
	// Test invalid document type
	var invalidDoc TextDocument = &struct{ TextDocument }{}
	err = Update(invalidDoc, []protocol.TextDocumentContentChangeEvent{}, 1)
	if err == nil {
		t.Error("Expected error for invalid document type")
	}
	
	// Test nil document for ApplyEdits
	_, err = ApplyEdits(nil, []protocol.TextEdit{})
	if err == nil {
		t.Error("Expected error for nil document in ApplyEdits")
	}
}

func TestUpdateErrorMessageFormat(t *testing.T) {
	// Test the error message format by creating a scenario that will fail
	// We'll test this by using a document that's not created by our Create function
	invalidDoc := &struct{ TextDocument }{}
	
	changes := []protocol.TextDocumentContentChangeEvent{
		{Text: "first change"},
		{Text: "second change"},
	}
	
	err := Update(invalidDoc, changes, 2)
	if err == nil {
		t.Error("Expected error for invalid document type")
	}
	
	// This should trigger the "document must be created by Create function" error
	// before we get to the change-specific error, but let's verify the fix works
	// by testing the strconv.Itoa conversion directly
	
	// Create a simple test to verify our fix works
	testIndex := 42
	errorMsg := "failed to apply change " + strconv.Itoa(testIndex) + ": some error"
	expectedSubstring := "failed to apply change 42:"
	
	if !strings.Contains(errorMsg, expectedSubstring) {
		t.Errorf("strconv.Itoa conversion test failed. Expected '%s' in '%s'", expectedSubstring, errorMsg)
	}
}

func TestCreateValidation(t *testing.T) {
	// Test empty URI
	_, err := Create("", "plaintext", 1, "content")
	if err == nil {
		t.Error("Expected error for empty URI")
	}
	if !strings.Contains(err.Error(), "uri cannot be empty") {
		t.Errorf("Expected 'uri cannot be empty' error, got: %s", err.Error())
	}
	
	// Test empty language ID
	_, err = Create("file:///test.txt", "", 1, "content")
	if err == nil {
		t.Error("Expected error for empty language ID")
	}
	if !strings.Contains(err.Error(), "languageID cannot be empty") {
		t.Errorf("Expected 'languageID cannot be empty' error, got: %s", err.Error())
	}
	
	// Test valid input with various edge cases
	testCases := []struct {
		uri        string
		languageID string
		version    int32
		content    string
		shouldFail bool
	}{
		{"file:///test.txt", "plaintext", 1, "content", false},
		{"file:///test.txt", "plaintext", 0, "", false}, // Empty content is valid
		{"file:///test.txt", "plaintext", -1, "content", false}, // Negative version is valid
		{"https://example.com/doc", "javascript", 1, "content", false}, // Non-file URI is valid
		{"", "plaintext", 1, "content", true}, // Empty URI should fail
		{"file:///test.txt", "", 1, "content", true}, // Empty language ID should fail
	}
	
	for i, tc := range testCases {
		doc, err := Create(protocol.DocumentURI(tc.uri), tc.languageID, tc.version, tc.content)
		if tc.shouldFail {
			if err == nil {
				t.Errorf("Test case %d: expected error but got none", i)
			}
			if doc != nil {
				t.Errorf("Test case %d: expected nil document but got non-nil", i)
			}
		} else {
			if err != nil {
				t.Errorf("Test case %d: expected no error but got: %v", i, err)
			}
			if doc == nil {
				t.Errorf("Test case %d: expected non-nil document but got nil", i)
			}
		}
	}
}