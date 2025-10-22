// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"context"
	"fmt"
	"sync"

	"github.com/TypeFox/go-lsp/protocol"
	"github.com/TypeFox/langium-to-go/textdoc"
)

// TextDocumentChangeEvent signals changes to a text document.
type TextDocumentChangeEvent struct {
	Document *textdoc.Overlay
}

// TextDocumentWillSaveEvent signals that a document will be saved.
type TextDocumentWillSaveEvent struct {
	Document *textdoc.Overlay
	Reason   protocol.TextDocumentSaveReason
}

// TextDocumentChangeHandler is called when a document changes.
type TextDocumentChangeHandler func(ctx context.Context, event *TextDocumentChangeEvent)

// TextDocumentWillSaveHandler is called when a document will be saved.
type TextDocumentWillSaveHandler func(ctx context.Context, event *TextDocumentWillSaveEvent)

// TextDocumentWillSaveWaitUntilHandler is called when a document will be saved and can return edits.
type TextDocumentWillSaveWaitUntilHandler func(ctx context.Context, event *TextDocumentWillSaveEvent) ([]protocol.TextEdit, error)

// TextDocuments manages a collection of text documents synchronized with the client.
// It handles document lifecycle events (open, change, close, save) and maintains
// in-memory representations using the textdoc.Overlay type.
type TextDocuments struct {
	mu              sync.RWMutex
	documents       map[protocol.DocumentURI]*textdoc.Overlay
	onDidOpen       []TextDocumentChangeHandler
	onDidChange     []TextDocumentChangeHandler
	onDidClose      []TextDocumentChangeHandler
	onDidSave       []TextDocumentChangeHandler
	onWillSave      []TextDocumentWillSaveHandler
	onWillSaveWaitUntil TextDocumentWillSaveWaitUntilHandler
}

// NewTextDocuments creates a new TextDocuments manager.
func NewTextDocuments() *TextDocuments {
	return &TextDocuments{
		documents: make(map[protocol.DocumentURI]*textdoc.Overlay),
	}
}

// Get retrieves a document by URI. Returns nil if the document is not managed.
func (td *TextDocuments) Get(uri protocol.DocumentURI) *textdoc.Overlay {
	td.mu.RLock()
	defer td.mu.RUnlock()
	return td.documents[uri]
}

// All returns all managed documents.
func (td *TextDocuments) All() []*textdoc.Overlay {
	td.mu.RLock()
	defer td.mu.RUnlock()
	
	docs := make([]*textdoc.Overlay, 0, len(td.documents))
	for _, doc := range td.documents {
		docs = append(docs, doc)
	}
	return docs
}

// Keys returns the URIs of all managed documents.
func (td *TextDocuments) Keys() []protocol.DocumentURI {
	td.mu.RLock()
	defer td.mu.RUnlock()
	
	keys := make([]protocol.DocumentURI, 0, len(td.documents))
	for uri := range td.documents {
		keys = append(keys, uri)
	}
	return keys
}

// OnDidOpen registers a handler for document open events.
func (td *TextDocuments) OnDidOpen(handler TextDocumentChangeHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.onDidOpen = append(td.onDidOpen, handler)
}

// OnDidChangeContent registers a handler for document change events.
func (td *TextDocuments) OnDidChangeContent(handler TextDocumentChangeHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.onDidChange = append(td.onDidChange, handler)
}

// OnDidClose registers a handler for document close events.
func (td *TextDocuments) OnDidClose(handler TextDocumentChangeHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.onDidClose = append(td.onDidClose, handler)
}

// OnDidSave registers a handler for document save events.
func (td *TextDocuments) OnDidSave(handler TextDocumentChangeHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.onDidSave = append(td.onDidSave, handler)
}

// OnWillSave registers a handler for document will-save events.
func (td *TextDocuments) OnWillSave(handler TextDocumentWillSaveHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.onWillSave = append(td.onWillSave, handler)
}

// OnWillSaveWaitUntil registers a handler that can provide edits during save.
func (td *TextDocuments) OnWillSaveWaitUntil(handler TextDocumentWillSaveWaitUntilHandler) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.onWillSaveWaitUntil = handler
}

// handleDidOpen processes a textDocument/didOpen notification.
func (td *TextDocuments) handleDidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	doc, err := textdoc.NewOverlay(
		params.TextDocument.URI,
		string(params.TextDocument.LanguageID),
		params.TextDocument.Version,
		params.TextDocument.Text,
	)
	if err != nil {
		return fmt.Errorf("failed to create document overlay: %w", err)
	}

	td.mu.Lock()
	td.documents[params.TextDocument.URI] = doc
	
	// Copy handlers while holding lock
	openHandlers := make([]TextDocumentChangeHandler, len(td.onDidOpen))
	copy(openHandlers, td.onDidOpen)
	changeHandlers := make([]TextDocumentChangeHandler, len(td.onDidChange))
	copy(changeHandlers, td.onDidChange)
	td.mu.Unlock()

	event := &TextDocumentChangeEvent{Document: doc}
	
	// Fire onDidOpen handlers
	for _, handler := range openHandlers {
		handler(ctx, event)
	}
	
	// Fire onDidChangeContent handlers (as per TypeScript implementation)
	for _, handler := range changeHandlers {
		handler(ctx, event)
	}

	return nil
}

// handleDidChange processes a textDocument/didChange notification.
func (td *TextDocuments) handleDidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) == 0 {
		return nil
	}

	td.mu.RLock()
	doc := td.documents[params.TextDocument.URI]
	td.mu.RUnlock()

	if doc == nil {
		return nil
	}

	if err := doc.Update(params.ContentChanges, params.TextDocument.Version); err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	td.mu.RLock()
	changeHandlers := make([]TextDocumentChangeHandler, len(td.onDidChange))
	copy(changeHandlers, td.onDidChange)
	td.mu.RUnlock()

	event := &TextDocumentChangeEvent{Document: doc}
	for _, handler := range changeHandlers {
		handler(ctx, event)
	}

	return nil
}

// handleDidClose processes a textDocument/didClose notification.
func (td *TextDocuments) handleDidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	td.mu.Lock()
	doc := td.documents[params.TextDocument.URI]
	if doc != nil {
		delete(td.documents, params.TextDocument.URI)
	}
	closeHandlers := make([]TextDocumentChangeHandler, len(td.onDidClose))
	copy(closeHandlers, td.onDidClose)
	td.mu.Unlock()

	if doc != nil {
		event := &TextDocumentChangeEvent{Document: doc}
		for _, handler := range closeHandlers {
			handler(ctx, event)
		}
	}

	return nil
}

// handleWillSave processes a textDocument/willSave notification.
func (td *TextDocuments) handleWillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) error {
	td.mu.RLock()
	doc := td.documents[params.TextDocument.URI]
	willSaveHandlers := make([]TextDocumentWillSaveHandler, len(td.onWillSave))
	copy(willSaveHandlers, td.onWillSave)
	td.mu.RUnlock()

	if doc != nil {
		event := &TextDocumentWillSaveEvent{
			Document: doc,
			Reason:   params.Reason,
		}
		for _, handler := range willSaveHandlers {
			handler(ctx, event)
		}
	}

	return nil
}

// handleWillSaveWaitUntil processes a textDocument/willSaveWaitUntil request.
func (td *TextDocuments) handleWillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	td.mu.RLock()
	doc := td.documents[params.TextDocument.URI]
	handler := td.onWillSaveWaitUntil
	td.mu.RUnlock()

	if doc != nil && handler != nil {
		event := &TextDocumentWillSaveEvent{
			Document: doc,
			Reason:   params.Reason,
		}
		return handler(ctx, event)
	}

	return []protocol.TextEdit{}, nil
}

// handleDidSave processes a textDocument/didSave notification.
func (td *TextDocuments) handleDidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	td.mu.RLock()
	doc := td.documents[params.TextDocument.URI]
	saveHandlers := make([]TextDocumentChangeHandler, len(td.onDidSave))
	copy(saveHandlers, td.onDidSave)
	td.mu.RUnlock()

	if doc != nil {
		event := &TextDocumentChangeEvent{Document: doc}
		for _, handler := range saveHandlers {
			handler(ctx, event)
		}
	}

	return nil
}

// Listen registers the TextDocuments manager with the language server handlers.
// This connects the document lifecycle events to the manager's internal handlers.
func (td *TextDocuments) Listen(handlers *LanguageServerHandlers) {
	handlers.DidOpen = td.handleDidOpen
	handlers.DidChange = td.handleDidChange
	handlers.DidClose = td.handleDidClose
}
