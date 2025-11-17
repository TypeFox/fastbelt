// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"context"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
)

func TestTextDocuments_Lifecycle(t *testing.T) {
	td := NewTextDocuments()
	ctx := context.Background()

	// Track events
	var openedDoc, changedDoc, closedDoc *TextDocumentChangeEvent
	td.OnDidOpen(func(ctx context.Context, event *TextDocumentChangeEvent) {
		openedDoc = event
	})
	td.OnDidChange(func(ctx context.Context, event *TextDocumentChangeEvent) {
		changedDoc = event
	})
	td.OnDidClose(func(ctx context.Context, event *TextDocumentChangeEvent) {
		closedDoc = event
	})

	// Open a document
	uri := protocol.DocumentURI("file:///test.txt")
	td.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	})

	// Verify document was opened
	if openedDoc == nil {
		t.Fatal("onDidOpen handler was not called")
	}
	if openedDoc.Document.URI() != uri {
		t.Errorf("expected URI %s, got %s", uri, openedDoc.Document.URI())
	}
	if changedDoc == nil {
		t.Fatal("onDidChangeContent handler was not called after open")
	}

	// Verify document is in collection
	doc := td.Get(uri)
	if doc == nil {
		t.Fatal("document not found in collection")
	}
	if doc.Version() != 1 {
		t.Errorf("expected version 1, got %d", doc.Version())
	}

	// Change the document
	changedDoc = nil
	td.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: uri},
			Version:                2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{Text: "Hello, Go!"},
		},
	})

	// Verify change event
	if changedDoc == nil {
		t.Fatal("onDidChangeContent handler was not called after change")
	}
	if doc.Version() != 2 {
		t.Errorf("expected version 2, got %d", doc.Version())
	}
	if doc.Text(nil) != "Hello, Go!" {
		t.Errorf("expected text 'Hello, Go!', got '%s'", doc.Text(nil))
	}

	// Close the document
	td.DidClose(ctx, &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
	})

	// Verify close event
	if closedDoc == nil {
		t.Fatal("onDidClose handler was not called")
	}
	if closedDoc.Document.URI() != uri {
		t.Errorf("expected URI %s, got %s", uri, closedDoc.Document.URI())
	}

	// Verify document was removed
	if td.Get(uri) != nil {
		t.Error("document should have been removed from collection")
	}
}

func TestTextDocuments_MultipleDocuments(t *testing.T) {
	td := NewTextDocuments()
	ctx := context.Background()

	// Open multiple documents
	uris := []protocol.DocumentURI{
		"file:///test1.txt",
		"file:///test2.txt",
		"file:///test3.txt",
	}

	for i, uri := range uris {
		td.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
			TextDocument: protocol.TextDocumentItem{
				URI:        uri,
				LanguageID: "plaintext",
				Version:    int32(i + 1),
				Text:       "content",
			},
		})
	}

	// Verify all documents are present
	all := td.All()
	if len(all) != len(uris) {
		t.Errorf("expected %d documents, got %d", len(uris), len(all))
	}

	keys := td.Keys()
	if len(keys) != len(uris) {
		t.Errorf("expected %d keys, got %d", len(uris), len(keys))
	}

	// Verify each document
	for _, uri := range uris {
		doc := td.Get(uri)
		if doc == nil {
			t.Errorf("document %s not found", uri)
		}
	}
}

func TestTextDocuments_IncrementalChanges(t *testing.T) {
	td := NewTextDocuments()
	ctx := context.Background()

	uri := protocol.DocumentURI("file:///test.txt")
	td.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	})

	// Apply incremental change
	td.DidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: uri},
			Version:                2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Range: &protocol.Range{
					Start: protocol.Position{Line: 0, Character: 7},
					End:   protocol.Position{Line: 0, Character: 12},
				},
			Text: "Go",
		},
	},
	})

	doc := td.Get(uri)
	if doc == nil {
		t.Fatal("document not found")
	}

	expected := "Hello, Go!"
	if doc.Text(nil) != expected {
		t.Errorf("expected text '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestTextDocuments_WillSave(t *testing.T) {
	td := NewTextDocuments()
	ctx := context.Background()

	uri := protocol.DocumentURI("file:///test.txt")
	td.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "content",
		},
	})

	// Track will-save event
	var willSaveEvent *TextDocumentWillSaveEvent
	td.OnWillSave(func(ctx context.Context, event *TextDocumentWillSaveEvent) {
		willSaveEvent = event
	})

	// Trigger will-save
	td.WillSave(ctx, &protocol.WillSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Reason:       protocol.Manual,
	})

	// Verify event
	if willSaveEvent == nil {
		t.Fatal("onWillSave handler was not called")
	}
	if willSaveEvent.Document.URI() != uri {
		t.Errorf("expected URI %s, got %s", uri, willSaveEvent.Document.URI())
	}
	if willSaveEvent.Reason != protocol.Manual {
		t.Errorf("expected reason Manual, got %v", willSaveEvent.Reason)
	}
}

func TestTextDocuments_WillSaveWaitUntil(t *testing.T) {
	td := NewTextDocuments()
	ctx := context.Background()

	uri := protocol.DocumentURI("file:///test.txt")
	td.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "content",
		},
	})

	// Register handler that returns edits
	expectedEdits := []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 7},
			},
			NewText: "modified",
		},
	}
	td.OnWillSaveWaitUntil(func(ctx context.Context, event *TextDocumentWillSaveEvent) ([]protocol.TextEdit, error) {
		return expectedEdits, nil
	})

	// Trigger will-save-wait-until
	edits, err := td.WillSaveWaitUntil(ctx, &protocol.WillSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
		Reason:       protocol.Manual,
	})
	if err != nil {
		t.Fatalf("WillSaveWaitUntil failed: %v", err)
	}

	// Verify edits
	if len(edits) != len(expectedEdits) {
		t.Errorf("expected %d edits, got %d", len(expectedEdits), len(edits))
	}
}

func TestTextDocuments_DidSave(t *testing.T) {
	td := NewTextDocuments()
	ctx := context.Background()

	uri := protocol.DocumentURI("file:///test.txt")
	td.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "content",
		},
	})

	// Track save event
	var savedDoc *TextDocumentChangeEvent
	td.OnDidSave(func(ctx context.Context, event *TextDocumentChangeEvent) {
		savedDoc = event
	})

	// Trigger save
	td.DidSave(ctx, &protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: uri},
	})

	// Verify event
	if savedDoc == nil {
		t.Fatal("onDidSave handler was not called")
	}
	if savedDoc.Document.URI() != uri {
		t.Errorf("expected URI %s, got %s", uri, savedDoc.Document.URI())
	}
}


