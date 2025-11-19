// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"context"
	"log"
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

// DocumentSyncher is the interface for handling LSP text document synchronization notifications.
// It processes document lifecycle events (open, change, close, save) and manages
// event handlers.
//
// Thread Safety:
// All methods are safe for concurrent use. Handler registration uses a mutex
// to ensure thread-safe handler management.
//
// Handler Execution:
// Handlers registered via OnDidOpen, OnDidChange, etc. are executed synchronously
// in the order they were registered. The context passed to handlers respects
// cancellation, and handlers should check ctx.Err() if they perform long operations.
type DocumentSyncher interface {
	OnDidOpen(handler TextDocumentChangeHandler)
	OnDidChange(handler TextDocumentChangeHandler)
	OnDidClose(handler TextDocumentChangeHandler)
	OnDidSave(handler TextDocumentChangeHandler)
	OnWillSave(handler TextDocumentWillSaveHandler)
	OnWillSaveWaitUntil(handler TextDocumentWillSaveWaitUntilHandler)
	DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams)
	DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams)
	DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams)
	WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams)
	WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error)
	DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams)
}

// DefaultDocumentSyncher is the default implementation of DocumentSyncher.
type DefaultDocumentSyncher struct {
	srv                 *LspServices
	handlersMu          sync.Mutex
	onDidOpen           []TextDocumentChangeHandler
	onDidChange         []TextDocumentChangeHandler
	onDidClose          []TextDocumentChangeHandler
	onDidSave           []TextDocumentChangeHandler
	onWillSave          []TextDocumentWillSaveHandler
	onWillSaveWaitUntil TextDocumentWillSaveWaitUntilHandler
}

// OnDidOpen registers a handler for document open events.
func (td *DefaultDocumentSyncher) OnDidOpen(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidOpen = append(td.onDidOpen, handler)
}

// OnDidChange registers a handler for document change events.
func (td *DefaultDocumentSyncher) OnDidChange(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidChange = append(td.onDidChange, handler)
}

// OnDidClose registers a handler for document close events.
func (td *DefaultDocumentSyncher) OnDidClose(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidClose = append(td.onDidClose, handler)
}

// OnDidSave registers a handler for document save events.
func (td *DefaultDocumentSyncher) OnDidSave(handler TextDocumentChangeHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onDidSave = append(td.onDidSave, handler)
}

// OnWillSave registers a handler for document will-save events.
func (td *DefaultDocumentSyncher) OnWillSave(handler TextDocumentWillSaveHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onWillSave = append(td.onWillSave, handler)
}

// OnWillSaveWaitUntil registers a handler that can provide edits during save.
// Only one handler can be registered; calling this multiple times will replace the previous handler.
func (td *DefaultDocumentSyncher) OnWillSaveWaitUntil(handler TextDocumentWillSaveWaitUntilHandler) {
	td.handlersMu.Lock()
	defer td.handlersMu.Unlock()
	td.onWillSaveWaitUntil = handler
}

// DidOpen processes a textDocument/didOpen notification.
func (td *DefaultDocumentSyncher) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) {
	doc, err := textdoc.NewOverlay(
		params.TextDocument.URI,
		string(params.TextDocument.LanguageID),
		params.TextDocument.Version,
		params.TextDocument.Text,
	)
	if err != nil {
		// Log error but continue - this is a notification, not a request
		log.Printf("failed to create document overlay: %v", err)
		return
	}

	if td.srv != nil && td.srv.TextdocServices != nil && td.srv.TextdocServices.Store != nil {
		td.srv.TextdocServices.Store.AddOverlay(doc)
	}

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
func (td *DefaultDocumentSyncher) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) {
	if len(params.ContentChanges) == 0 {
		return
	}

	if td.srv == nil || td.srv.TextdocServices == nil || td.srv.TextdocServices.Store == nil {
		return
	}

	doc := td.srv.TextdocServices.Store.GetOverlay(params.TextDocument.URI)
	if doc == nil {
		return
	}

	if err := doc.Update(params.ContentChanges, params.TextDocument.Version); err != nil {
		// Log error but continue - this is a notification, not a request
		log.Printf("failed to update document: %v", err)
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
func (td *DefaultDocumentSyncher) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) {
	if td.srv == nil || td.srv.TextdocServices == nil || td.srv.TextdocServices.Store == nil {
		return
	}

	doc := td.srv.TextdocServices.Store.GetOverlay(params.TextDocument.URI)
	if doc == nil {
		return
	}

	td.srv.TextdocServices.Store.RemoveOverlay(params.TextDocument.URI)

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
func (td *DefaultDocumentSyncher) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) {
	if td.srv == nil || td.srv.TextdocServices == nil || td.srv.TextdocServices.Store == nil {
		return
	}

	doc := td.srv.TextdocServices.Store.GetOverlay(params.TextDocument.URI)
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
func (td *DefaultDocumentSyncher) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	var doc *textdoc.Overlay
	if td.srv != nil && td.srv.TextdocServices != nil && td.srv.TextdocServices.Store != nil {
		doc = td.srv.TextdocServices.Store.GetOverlay(params.TextDocument.URI)
	}

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
func (td *DefaultDocumentSyncher) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) {
	if td.srv == nil || td.srv.TextdocServices == nil || td.srv.TextdocServices.Store == nil {
		return
	}

	doc := td.srv.TextdocServices.Store.GetOverlay(params.TextDocument.URI)
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
