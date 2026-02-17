// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"log"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/lsp"
)

// TextDocumentChangeEvent signals changes to a text document.
type TextDocumentChangeEvent struct {
	Document *textdoc.Overlay
}

// TextDocumentWillSaveEvent signals that a document will be saved.
type TextDocumentWillSaveEvent struct {
	Document *textdoc.Overlay
	Reason   lsp.TextDocumentSaveReason
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
type TextDocumentWillSaveWaitUntilHandler func(ctx context.Context, event *TextDocumentWillSaveEvent) ([]lsp.TextEdit, error)

// DocumentSyncher is the interface for handling LSP text document synchronization notifications.
// It processes document lifecycle events (open, change, close, save).
//
// Thread Safety:
// All methods are safe for concurrent use.
type DocumentSyncher interface {
	DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams)
	DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams)
	DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams)
	WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams)
	WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error)
	DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams)
}

// DefaultDocumentSyncher is the default implementation of DocumentSyncher.
type DefaultDocumentSyncher struct {
	srv ServerSrvCont
}

// NewDefaultDocumentSyncher creates a new default document syncher.
func NewDefaultDocumentSyncher(srv ServerSrvCont) DocumentSyncher {
	return &DefaultDocumentSyncher{srv: srv}
}

// DidOpen processes a textDocument/didOpen notification.
func (ds *DefaultDocumentSyncher) DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) {
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

	ds.srv.Textdoc().Store.AddOverlay(doc)
	ds.srv.Workspace().DocumentUpdater.Update(ctx, []textdoc.Handle{doc}, nil)
}

// DidChange processes a textDocument/didChange notification.
func (ds *DefaultDocumentSyncher) DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams) {
	if len(params.ContentChanges) == 0 {
		return
	}

	doc := ds.srv.Textdoc().Store.GetOverlay(params.TextDocument.URI)
	if doc == nil {
		return
	}

	if err := doc.Update(params.ContentChanges, params.TextDocument.Version); err != nil {
		// Log error but continue - this is a notification, not a request
		log.Printf("failed to update document: %v", err)
		return
	}

	ds.srv.Workspace().DocumentUpdater.Update(ctx, []textdoc.Handle{doc}, nil)
}

// DidClose processes a textDocument/didClose notification.
func (ds *DefaultDocumentSyncher) DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) {
	ds.srv.Textdoc().Store.RemoveOverlay(params.TextDocument.URI)
	connection := ds.srv.Server().Connection
	if connection != nil {
		// Ensure we clear diagnostics on close
		// TODO: Make this configurable - some adopters might want to keep diagnostics for closed documents
		client := lsp.ClientDispatcher(connection)
		err := client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: []lsp.Diagnostic{},
		})
		if err != nil {
			log.Printf("failed to publish diagnostics after document close: %v", err)
		}
	}
	// TODO msujew: Once we start handling cross-file references, we shouldn't delete the document.
	uri := core.ParseURI(string(params.TextDocument.URI))
	ds.srv.Workspace().DocumentUpdater.Update(ctx, nil, []core.URI{uri})
}

// WillSave processes a textDocument/willSave notification.
func (ds *DefaultDocumentSyncher) WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams) {
}

// WillSaveWaitUntil processes a textDocument/willSaveWaitUntil request.
func (ds *DefaultDocumentSyncher) WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error) {
	return []lsp.TextEdit{}, nil
}

// DidSave processes a textDocument/didSave notification.
func (ds *DefaultDocumentSyncher) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) {
}
