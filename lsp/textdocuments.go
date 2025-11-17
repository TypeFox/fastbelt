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
// Handlers are invoked synchronously in the order they were registered.
// If a handler blocks, it will delay subsequent handlers and the notification response.
type TextDocumentChangeHandler func(ctx context.Context, event *TextDocumentChangeEvent)

// TextDocumentWillSaveHandler is called when a document will be saved.
// Handlers are invoked synchronously in the order they were registered.
type TextDocumentWillSaveHandler func(ctx context.Context, event *TextDocumentWillSaveEvent)

// TextDocumentWillSaveWaitUntilHandler is called when a document will be saved and can return edits.
// Only one handler can be registered for this event to ensure deterministic edit behavior.
type TextDocumentWillSaveWaitUntilHandler func(ctx context.Context, event *TextDocumentWillSaveEvent) ([]protocol.TextEdit, error)

// TextDocuments manages a collection of text documents synchronized with the client.
// It handles document lifecycle events (open, change, close, save) and maintains
// in-memory representations using the textdoc.Overlay type.
//
// Thread Safety:
// All methods are safe for concurrent use. Document access uses a read-write mutex
// for efficient concurrent reads, while handler registration uses a separate mutex
// to avoid contention.
//
// Handler Execution:
// Handlers registered via OnDidOpen, OnDidChange, etc. are executed synchronously
// in the order they were registered. The context passed to handlers respects
// cancellation, and handlers should check ctx.Err() if they perform long operations.
type TextDocuments struct {
	documentsMu         sync.RWMutex
	documents           map[protocol.DocumentURI]*textdoc.Overlay
	handlersMu          sync.Mutex
	onDidOpen           []TextDocumentChangeHandler
	onDidChange         []TextDocumentChangeHandler
	onDidClose          []TextDocumentChangeHandler
	onDidSave           []TextDocumentChangeHandler
	onWillSave          []TextDocumentWillSaveHandler
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
	td.documentsMu.RLock()
	defer td.documentsMu.RUnlock()
	return td.documents[uri]
}

// All returns all managed documents.
func (td *TextDocuments) All() []*textdoc.Overlay {
	td.documentsMu.RLock()
	defer td.documentsMu.RUnlock()

	docs := make([]*textdoc.Overlay, 0, len(td.documents))
	for _, doc := range td.documents {
		docs = append(docs, doc)
	}
	return docs
}

// Keys returns the URIs of all managed documents.
func (td *TextDocuments) Keys() []protocol.DocumentURI {
	td.documentsMu.RLock()
	defer td.documentsMu.RUnlock()

	keys := make([]protocol.DocumentURI, 0, len(td.documents))
	for uri := range td.documents {
		keys = append(keys, uri)
	}
	return keys
}

// OnDidOpen registers a handler for document open events.
func (td *TextDocuments) OnDidOpen(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidOpen = append(td.onDidOpen, handler)
}

// OnDidChange registers a handler for document change events.
func (td *TextDocuments) OnDidChange(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidChange = append(td.onDidChange, handler)
}

// OnDidClose registers a handler for document close events.
func (td *TextDocuments) OnDidClose(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidClose = append(td.onDidClose, handler)
}

// OnDidSave registers a handler for document save events.
func (td *TextDocuments) OnDidSave(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidSave = append(td.onDidSave, handler)
}

// OnWillSave registers a handler for document will-save events.
func (td *TextDocuments) OnWillSave(handler TextDocumentWillSaveHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onWillSave = append(td.onWillSave, handler)
}

// OnWillSaveWaitUntil registers a handler that can provide edits during save.
// Only one handler can be registered; calling this multiple times will replace the previous handler.
func (td *TextDocuments) OnWillSaveWaitUntil(handler TextDocumentWillSaveWaitUntilHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onWillSaveWaitUntil = handler
}

// DidOpen processes a textDocument/didOpen notification.
func (td *TextDocuments) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) {
	doc, err := textdoc.NewOverlay(
		params.TextDocument.URI,
		string(params.TextDocument.LanguageID),
		params.TextDocument.Version,
		params.TextDocument.Text,
	)
	if err != nil {
		// Log error but continue - this is a notification, not a request
		fmt.Printf("failed to create document overlay: %v\n", err)
		return
	}

	td.documentsMu.Lock()
	td.documents[params.TextDocument.URI] = doc
	td.documentsMu.Unlock()

	// Copy handlers while holding lock
	td.handlersMu.Lock()
	openHandlers := make([]TextDocumentChangeHandler, len(td.onDidOpen))
	copy(openHandlers, td.onDidOpen)
	changeHandlers := make([]TextDocumentChangeHandler, len(td.onDidChange))
	copy(changeHandlers, td.onDidChange)
	td.handlersMu.Unlock()

	event := &TextDocumentChangeEvent{Document: doc}

	// Fire onDidOpen handlers
	for _, handler := range openHandlers {
		if ctx.Err() != nil {
			return
		}
		handler(ctx, event)
	}

	// Fire onDidChange handlers (as per TypeScript implementation)
	for _, handler := range changeHandlers {
		if ctx.Err() != nil {
			return
		}
		handler(ctx, event)
	}
}

// DidChange processes a textDocument/didChange notification.
func (td *TextDocuments) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) {
	if len(params.ContentChanges) == 0 {
		return
	}

	td.documentsMu.RLock()
	doc := td.documents[params.TextDocument.URI]
	td.documentsMu.RUnlock()

	if doc == nil {
		return
	}

	if err := doc.Update(params.ContentChanges, params.TextDocument.Version); err != nil {
		// Log error but continue - this is a notification, not a request
		fmt.Printf("failed to update document: %v\n", err)
		return
	}

	td.handlersMu.Lock()
	changeHandlers := make([]TextDocumentChangeHandler, len(td.onDidChange))
	copy(changeHandlers, td.onDidChange)
	td.handlersMu.Unlock()

	event := &TextDocumentChangeEvent{Document: doc}
	for _, handler := range changeHandlers {
		if ctx.Err() != nil {
			return
		}
		handler(ctx, event)
	}
}

// DidClose processes a textDocument/didClose notification.
func (td *TextDocuments) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) {
	td.documentsMu.Lock()
	doc := td.documents[params.TextDocument.URI]
	if doc != nil {
		delete(td.documents, params.TextDocument.URI)
	}
	td.documentsMu.Unlock()

	if doc == nil {
		return
	}

	td.handlersMu.Lock()
	closeHandlers := make([]TextDocumentChangeHandler, len(td.onDidClose))
	copy(closeHandlers, td.onDidClose)
	td.handlersMu.Unlock()

	event := &TextDocumentChangeEvent{Document: doc}
	for _, handler := range closeHandlers {
		if ctx.Err() != nil {
			return
		}
		handler(ctx, event)
	}
}

// WillSave processes a textDocument/willSave notification.
func (td *TextDocuments) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) {
	td.documentsMu.RLock()
	doc := td.documents[params.TextDocument.URI]
	td.documentsMu.RUnlock()

	if doc == nil {
		return
	}

	td.handlersMu.Lock()
	willSaveHandlers := make([]TextDocumentWillSaveHandler, len(td.onWillSave))
	copy(willSaveHandlers, td.onWillSave)
	td.handlersMu.Unlock()

	event := &TextDocumentWillSaveEvent{
		Document: doc,
		Reason:   params.Reason,
	}
	for _, handler := range willSaveHandlers {
		if ctx.Err() != nil {
			return
		}
		handler(ctx, event)
	}
}

// WillSaveWaitUntil processes a textDocument/willSaveWaitUntil request.
func (td *TextDocuments) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	td.documentsMu.RLock()
	doc := td.documents[params.TextDocument.URI]
	td.documentsMu.RUnlock()

	td.handlersMu.Lock()
	handler := td.onWillSaveWaitUntil
	td.handlersMu.Unlock()

	if doc != nil && handler != nil {
		event := &TextDocumentWillSaveEvent{
			Document: doc,
			Reason:   params.Reason,
		}
		return handler(ctx, event)
	}

	return []protocol.TextEdit{}, nil
}

// DidSave processes a textDocument/didSave notification.
func (td *TextDocuments) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) {
	td.documentsMu.RLock()
	doc := td.documents[params.TextDocument.URI]
	td.documentsMu.RUnlock()

	if doc == nil {
		return
	}

	td.handlersMu.Lock()
	saveHandlers := make([]TextDocumentChangeHandler, len(td.onDidSave))
	copy(saveHandlers, td.onDidSave)
	td.handlersMu.Unlock()

	event := &TextDocumentChangeEvent{Document: doc}
	for _, handler := range saveHandlers {
		if ctx.Err() != nil {
			return
		}
		handler(ctx, event)
	}
}
