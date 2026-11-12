// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"bytes"
	"context"
	"log"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/parallel"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
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
//
// Handlers for didOpen, didChange, and didClose run while the
// [workspace.Lock] is held for writing, before the workspace is rebuilt. They
// may access the workspace directly but must not acquire the lock themselves,
// which would deadlock. Their ctx is the write's context: it is cancelled as
// soon as a newer write arrives. A handler that blocks delays subsequent
// handlers, the rebuild, and every later write and read.
//
// Handlers for didSave run without the lock.
type TextDocumentChangeHandler func(ctx context.Context, event *TextDocumentChangeEvent)

// TextDocumentWillSaveHandler is called when a document will be saved.
// Handlers are invoked synchronously in the order they were registered.
type TextDocumentWillSaveHandler func(ctx context.Context, event *TextDocumentWillSaveEvent)

// TextDocumentWillSaveWaitUntilHandler is called when a document will be saved and can return edits.
// Only one handler can be registered for this event to ensure deterministic edit behavior.
type TextDocumentWillSaveWaitUntilHandler func(ctx context.Context, event *TextDocumentWillSaveEvent) ([]lsp.TextEdit, error)

// DocumentSynchronizer is a service for handling LSP text document synchronization notifications.
//
// DidOpen, DidChange, and DidClose each apply their change to the
// [textdoc.Store] while holding the [workspace.Lock] for writing, so readers
// never observe a text change without the matching workspace reset. They
// return once the text change has been applied, which requires the running
// write to finish and all readers to drain; the registered handlers (for
// DidOpen and DidChange) and the rebuild through [workspace.DocumentUpdater]
// continue in the background under the same write. Changes take effect in the
// order the notifications were received.
//
// Thread Safety: All methods are safe for concurrent use.
type DocumentSynchronizer interface {
	DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams)
	DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams)
	DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams)
	WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams)
	WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error)
	DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams)

	OnDidOpen(handler TextDocumentChangeHandler)
	OnDidChange(handler TextDocumentChangeHandler)
	OnDidClose(handler TextDocumentChangeHandler)
	OnWillSave(handler TextDocumentWillSaveHandler)
	OnWillSaveWaitUntil(handler TextDocumentWillSaveWaitUntilHandler)
	OnDidSave(handler TextDocumentChangeHandler)
}

// DefaultDocumentSynchronizer is the default implementation of DocumentSynchronizer.
type DefaultDocumentSynchronizer struct {
	sc                        *service.Container
	didOpenHandlers           []TextDocumentChangeHandler
	didChangeHandlers         []TextDocumentChangeHandler
	didCloseHandlers          []TextDocumentChangeHandler
	willSaveHandlers          []TextDocumentWillSaveHandler
	willSaveWaitUntilHandlers []TextDocumentWillSaveWaitUntilHandler
	didSaveHandlers           []TextDocumentChangeHandler
	mu                        sync.RWMutex
}

// NewDefaultDocumentSynchronizer creates a new [DefaultDocumentSynchronizer].
func NewDefaultDocumentSynchronizer(sc *service.Container) DocumentSynchronizer {
	return &DefaultDocumentSynchronizer{sc: sc}
}

func (s *DefaultDocumentSynchronizer) DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) {
	lock := s.sc.MustGet[workspace.Lock]()
	waiter := parallel.Waiter{}

	// Run the write in a goroutine: the notification only waits until the
	// overlay has been added, not for the workspace build that follows.
	// The write gets a fresh context, because ctx ends with the notification.
	go func() {
		lock.Write(context.Background(), func(writeCtx context.Context) {
			textdocStore := s.sc.MustGet[textdoc.Store]()
			uri := core.NormalizeURI(params.TextDocument.URI)
			existing := textdocStore.Get(uri)

			doc := textdoc.NewOverlay(
				uri,
				string(params.TextDocument.LanguageID),
				params.TextDocument.Version,
				params.TextDocument.Text,
			)

			textdocStore.AddOverlay(doc)
			// Document has been added, notification can return
			waiter.Close()
			s.mu.RLock()
			for _, handler := range s.didOpenHandlers {
				handler(writeCtx, &TextDocumentChangeEvent{Document: doc})
			}
			s.mu.RUnlock()
			if existing != nil && existing.Text(nil) == params.TextDocument.Text && existing.LanguageID() == doc.LanguageID() {
				// The text editor content and language are the same as the file's.
				// Nothing changed, don't start another workspace build.
				// Note: Language clients will send didClose+didOpen with a different
				// language ID if the user changes it - this will correctly generate
				// a new document and start a new build.
				return
			}
			updater := s.sc.MustGet[workspace.DocumentUpdater]()
			updater.Update(writeCtx, []textdoc.Handle{doc}, nil)
		})
	}()

	// Wait for signal before returning from the notification
	waiter.Await()
}

func (s *DefaultDocumentSynchronizer) OnDidOpen(handler TextDocumentChangeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.didOpenHandlers = append(s.didOpenHandlers, handler)
}

func (s *DefaultDocumentSynchronizer) DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams) {
	if len(params.ContentChanges) == 0 {
		return
	}
	lock := s.sc.MustGet[workspace.Lock]()
	waiter := parallel.Waiter{}

	// Run the write in a goroutine: the notification only waits until the
	// overlay has been updated, not for the workspace build that follows.
	// The write gets a fresh context, because ctx ends with the notification.
	go func() {
		lock.Write(context.Background(), func(writeCtx context.Context) {

			textdocStore := s.sc.MustGet[textdoc.Store]()
			uri := core.NormalizeURI(params.TextDocument.URI)
			doc := textdocStore.GetOverlay(uri)
			if doc == nil {
				waiter.Close()
				return
			}

			oldContent := doc.Content()

			err := doc.Update(params.ContentChanges, params.TextDocument.Version)
			// Document has been updated, notification can return
			waiter.Close()
			if err != nil {
				// Log error and return - this is a notification, not a request
				log.Printf("failed to update document: %v", err)
				return
			}

			if bytes.Equal(oldContent, doc.Content()) {
				// Document hasn't actually changed, no need to notify handlers or update the workspace
				return
			}
			s.mu.RLock()
			for _, handler := range s.didChangeHandlers {
				handler(writeCtx, &TextDocumentChangeEvent{Document: doc})
			}
			s.mu.RUnlock()
			updater := s.sc.MustGet[workspace.DocumentUpdater]()
			updater.Update(writeCtx, []textdoc.Handle{doc}, nil)
		})
	}()

	// Wait for signal before returning from the notification
	waiter.Await()
}

func (s *DefaultDocumentSynchronizer) OnDidChange(handler TextDocumentChangeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.didChangeHandlers = append(s.didChangeHandlers, handler)
}

func (s *DefaultDocumentSynchronizer) DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) {

	lock := s.sc.MustGet[workspace.Lock]()
	waiter := parallel.Waiter{}

	// Run the write in a goroutine with a fresh context, because ctx ends
	// with the notification. The notification waits until the handlers have
	// run and the overlay has been removed.
	go func() {
		lock.Write(context.Background(), func(writeCtx context.Context) {
			textdocStore := s.sc.MustGet[textdoc.Store]()
			uri := core.NormalizeURI(params.TextDocument.URI)
			existing := textdocStore.GetOverlay(uri)
			if existing != nil {
				s.mu.Lock()
				for _, handler := range s.didCloseHandlers {
					handler(writeCtx, &TextDocumentChangeEvent{Document: existing})
				}
				s.mu.Unlock()
			}
			textdocStore.RemoveOverlay(uri)
			waiter.Close()
			// TODO: Maybe do some more work here - a closed file might trigger a new workspace build
			// see also https://github.com/eclipse-fastbelt/fastbelt/issues/185
		})
	}()
	waiter.Await()
}

func (s *DefaultDocumentSynchronizer) OnDidClose(handler TextDocumentChangeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.didCloseHandlers = append(s.didCloseHandlers, handler)
}

// WillSave does nothing by default.
func (s *DefaultDocumentSynchronizer) WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams) {
	store := s.sc.MustGet[textdoc.Store]()
	uri := core.NormalizeURI(params.TextDocument.URI)
	doc := store.GetOverlay(uri)
	if doc != nil {
		s.mu.RLock()
		for _, handler := range s.willSaveHandlers {
			handler(ctx, &TextDocumentWillSaveEvent{Document: doc, Reason: params.Reason})
		}
		s.mu.RUnlock()
	}
}

func (s *DefaultDocumentSynchronizer) OnWillSave(handler TextDocumentWillSaveHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.willSaveHandlers = append(s.willSaveHandlers, handler)
}

// WillSaveWaitUntil does nothing by default.
func (s *DefaultDocumentSynchronizer) WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error) {
	edits := []lsp.TextEdit{}
	store := s.sc.MustGet[textdoc.Store]()
	uri := core.NormalizeURI(params.TextDocument.URI)
	doc := store.GetOverlay(uri)
	if doc == nil {
		return edits, nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, handler := range s.willSaveWaitUntilHandlers {
		handlerEdits, err := handler(ctx, &TextDocumentWillSaveEvent{Document: doc, Reason: params.Reason})
		if err != nil {
			return nil, err
		}
		edits = append(edits, handlerEdits...)
	}
	return edits, nil
}

func (s *DefaultDocumentSynchronizer) OnWillSaveWaitUntil(handler TextDocumentWillSaveWaitUntilHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.willSaveWaitUntilHandlers = append(s.willSaveWaitUntilHandlers, handler)
}

// DidSave does nothing by default.
func (s *DefaultDocumentSynchronizer) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) {
	store := s.sc.MustGet[textdoc.Store]()
	uri := core.NormalizeURI(params.TextDocument.URI)
	doc := store.GetOverlay(uri)
	if doc != nil {
		s.mu.RLock()
		for _, handler := range s.didSaveHandlers {
			handler(ctx, &TextDocumentChangeEvent{Document: doc})
		}
		s.mu.RUnlock()
	}
}

func (s *DefaultDocumentSynchronizer) OnDidSave(handler TextDocumentChangeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.didSaveHandlers = append(s.didSaveHandlers, handler)
}
