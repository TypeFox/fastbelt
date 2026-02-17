// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"sync"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// mockDocumentUpdater is a test implementation of DocumentUpdater that tracks calls.
type mockDocumentUpdater struct {
	mu          sync.Mutex
	updateCalls []updateCall
}

type updateCall struct {
	changed []textdoc.Handle
	deleted []core.URI
}

func (m *mockDocumentUpdater) Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateCalls = append(m.updateCalls, updateCall{changed: changed, deleted: deleted})
}

func (m *mockDocumentUpdater) getUpdateCalls() []updateCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]updateCall, len(m.updateCalls))
	copy(result, m.updateCalls)
	return result
}

func (m *mockDocumentUpdater) reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateCalls = nil
}

func createTestServices() ServerSrvCont {
	s := &serverSrvContTest{}
	textdoc.CreateDefaultServices(s)
	workspace.CreateDefaultServices(s)
	return s
}

func createTestServicesWithUpdater() (ServerSrvCont, *mockDocumentUpdater) {
	s := createTestServices()
	updater := &mockDocumentUpdater{}
	s.Workspace().DocumentUpdater = updater
	return s, updater
}

func TestTextDocuments_Lifecycle(t *testing.T) {
	s, updater := createTestServicesWithUpdater()
	ds := &DefaultDocumentSyncher{srv: s}
	ctx := context.Background()

	// Open a document
	uri := lsp.DocumentURI("file:///test.txt")
	ds.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	})

	// Verify DocumentUpdater.Update was called for DidOpen
	updateCalls := updater.getUpdateCalls()
	if len(updateCalls) != 1 {
		t.Fatalf("expected 1 update call, got %d", len(updateCalls))
	}
	if len(updateCalls[0].changed) != 1 {
		t.Fatalf("expected 1 changed document in update call, got %d", len(updateCalls[0].changed))
	}
	if updateCalls[0].changed[0].URI() != uri {
		t.Errorf("expected URI %s, got %s", uri, updateCalls[0].changed[0].URI())
	}

	// Verify document is in collection
	doc := s.Textdoc().Store.GetOverlay(uri)
	if doc == nil {
		t.Fatal("document not found in collection")
	}
	if doc.Version() != 1 {
		t.Errorf("expected version 1, got %d", doc.Version())
	}

	// Change the document
	updater.reset()
	ds.DidChange(ctx, &lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: uri},
			Version:                2,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{
			{Text: "Hello, Go!"},
		},
	})

	// Verify DocumentUpdater.Update was called for DidChange
	updateCalls = updater.getUpdateCalls()
	if len(updateCalls) != 1 {
		t.Fatalf("expected 1 update call after change, got %d", len(updateCalls))
	}
	if doc.Version() != 2 {
		t.Errorf("expected version 2, got %d", doc.Version())
	}
	if doc.Text(nil) != "Hello, Go!" {
		t.Errorf("expected text 'Hello, Go!', got '%s'", doc.Text(nil))
	}

	// Close the document
	updater.reset()
	ds.DidClose(ctx, &lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})

	// Verify overlay was removed from store
	if s.Textdoc().Store.GetOverlay(uri) != nil {
		t.Error("document should have been removed from collection")
	}

	// Verify DocumentUpdater.Update was called with deleted URI
	updateCalls = updater.getUpdateCalls()
	if len(updateCalls) != 1 {
		t.Fatalf("expected 1 update call after close, got %d", len(updateCalls))
	}
	if len(updateCalls[0].deleted) != 1 {
		t.Fatalf("expected 1 deleted URI in update call, got %d", len(updateCalls[0].deleted))
	}
}

func TestTextDocuments_MultipleDocuments(t *testing.T) {
	s, _ := createTestServicesWithUpdater()
	ds := &DefaultDocumentSyncher{srv: s}
	ctx := context.Background()

	// Open multiple documents
	uris := []lsp.DocumentURI{
		"file:///test1.txt",
		"file:///test2.txt",
		"file:///test3.txt",
	}

	for i, uri := range uris {
		ds.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{
			TextDocument: lsp.TextDocumentItem{
				URI:        uri,
				LanguageID: "plaintext",
				Version:    int32(i + 1),
				Text:       "content",
			},
		})
	}

	// Verify all documents are present
	all := s.Textdoc().Store.AllOverlays()
	if len(all) != len(uris) {
		t.Errorf("expected %d documents, got %d", len(uris), len(all))
	}

	keys := s.Textdoc().Store.KeysOverlays()
	if len(keys) != len(uris) {
		t.Errorf("expected %d keys, got %d", len(uris), len(keys))
	}

	// Verify each document
	for _, uri := range uris {
		doc := s.Textdoc().Store.GetOverlay(uri)
		if doc == nil {
			t.Errorf("document %s not found", uri)
		}
	}
}

func TestTextDocuments_IncrementalChanges(t *testing.T) {
	s, _ := createTestServicesWithUpdater()
	ds := &DefaultDocumentSyncher{srv: s}
	ctx := context.Background()

	uri := lsp.DocumentURI("file:///test.txt")
	ds.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	})

	// Apply incremental change
	ds.DidChange(ctx, &lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: uri},
			Version:                2,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{
			{
				Range: &lsp.Range{
					Start: lsp.Position{Line: 0, Character: 7},
					End:   lsp.Position{Line: 0, Character: 12},
				},
				Text: "Go",
			},
		},
	})

	doc := s.Textdoc().Store.GetOverlay(uri)
	if doc == nil {
		t.Fatal("document not found")
	}

	expected := "Hello, Go!"
	if doc.Text(nil) != expected {
		t.Errorf("expected text '%s', got '%s'", expected, doc.Text(nil))
	}
}

func TestTextDocuments_WillSave(t *testing.T) {
	s, _ := createTestServicesWithUpdater()
	ds := &DefaultDocumentSyncher{srv: s}
	ctx := context.Background()

	uri := lsp.DocumentURI("file:///test.txt")
	ds.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "content",
		},
	})

	// Trigger will-save
	ds.WillSave(ctx, &lsp.WillSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Reason:       lsp.Manual,
	})

	// Verify document still exists
	doc := s.Textdoc().Store.GetOverlay(uri)
	if doc == nil {
		t.Fatal("document should still exist after WillSave")
	}
}

func TestTextDocuments_WillSaveWaitUntil(t *testing.T) {
	s, _ := createTestServicesWithUpdater()
	ds := &DefaultDocumentSyncher{srv: s}
	ctx := context.Background()

	uri := lsp.DocumentURI("file:///test.txt")
	ds.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "content",
		},
	})

	// Trigger will-save-wait-until
	edits, err := ds.WillSaveWaitUntil(ctx, &lsp.WillSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Reason:       lsp.Manual,
	})
	if err != nil {
		t.Fatalf("WillSaveWaitUntil failed: %v", err)
	}

	// Verify empty edits are returned
	if len(edits) != 0 {
		t.Errorf("expected 0 edits, got %d", len(edits))
	}
}

func TestTextDocuments_DidSave(t *testing.T) {
	s, _ := createTestServicesWithUpdater()
	ds := &DefaultDocumentSyncher{srv: s}
	ctx := context.Background()

	uri := lsp.DocumentURI("file:///test.txt")
	ds.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "plaintext",
			Version:    1,
			Text:       "content",
		},
	})

	// Trigger save
	ds.DidSave(ctx, &lsp.DidSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})

	// Verify document still exists
	doc := s.Textdoc().Store.GetOverlay(uri)
	if doc == nil {
		t.Fatal("document should still exist after DidSave")
	}
	if doc.URI() != uri {
		t.Errorf("expected URI %s, got %s", uri, doc.URI())
	}
}
