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
// soon as a newer write arrives, which can be before the handler starts.
// Work that must not be skipped, such as notifying the client, should
// detach from it with [context.WithoutCancel]. A handler that blocks delays
// subsequent handlers, the rebuild, and every later write and read.
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
// order the notifications were received. Every write ends with a rebuild,
// including DidClose and notifications that leave the text unchanged, because
// acquiring the write lock cancels the build that was running.
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

// write runs apply under the [workspace.Lock] write lock and returns as soon
// as apply has returned; the rest of the write continues in the background.
//
// apply performs the text change. It reports the handles it changed and,
// optionally, a notify function that runs the registered handlers once the
// notification has returned. Every write ends with a call to
// [workspace.DocumentUpdater.Update], even when apply changed nothing:
// acquiring the write lock cancels the running build and resets the workspace
// state, so a write that skipped the update would leave the documents of the
// cancelled build incomplete.
func (s *DefaultDocumentSynchronizer) write(apply func(ctx context.Context) (changed []textdoc.Handle, notify func())) {
	lock := s.sc.MustGet[workspace.Lock]()
	applied := make(chan struct{})
	// The write gets a fresh context, because the notification's context ends
	// when the notification returns.
	go lock.Write(context.Background(), func(ctx context.Context) {
		changed, notify := func() ([]textdoc.Handle, func()) {
			defer close(applied)
			return apply(ctx)
		}()
		if notify != nil {
			notify()
		}
		s.sc.MustGet[workspace.DocumentUpdater]().Update(ctx, changed, nil)
	})
	<-applied
}

// notify calls the handlers with a change event for doc. handlers points to
// one of the handler lists of s, which is read under s.mu.
func (s *DefaultDocumentSynchronizer) notify(ctx context.Context, handlers *[]TextDocumentChangeHandler, doc *textdoc.Overlay) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, handler := range *handlers {
		handler(ctx, &TextDocumentChangeEvent{Document: doc})
	}
}

func (s *DefaultDocumentSynchronizer) DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) {
	s.write(func(writeCtx context.Context) ([]textdoc.Handle, func()) {
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

		notify := func() { s.notify(writeCtx, &s.didOpenHandlers, doc) }
		if existing != nil && existing.Text(nil) == params.TextDocument.Text && existing.LanguageID() == doc.LanguageID() {
			// The text editor content and language are the same as the file's.
			// Nothing changed, keep the document that is already built.
			// Note: Language clients will send didClose+didOpen with a different
			// language ID if the user changes it - this will correctly generate
			// a new document and start a new build.
			return nil, notify
		}
		return []textdoc.Handle{doc}, notify
	})
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
	s.write(func(writeCtx context.Context) ([]textdoc.Handle, func()) {
		textdocStore := s.sc.MustGet[textdoc.Store]()
		uri := core.NormalizeURI(params.TextDocument.URI)
		doc := textdocStore.GetOverlay(uri)
		if doc == nil {
			return nil, nil
		}

		oldContent := doc.Content()
		if err := doc.Update(params.ContentChanges, params.TextDocument.Version); err != nil {
			// Log the error and carry on - this is a notification, not a request
			// Note: At this point, client and server are desynchronized - but the client is at fault here
			log.Printf("failed to update document: %v", err)
			return nil, nil
		}
		if bytes.Equal(oldContent, doc.Content()) {
			// Document hasn't actually changed, no need to notify handlers
			return nil, nil
		}
		return []textdoc.Handle{doc}, func() { s.notify(writeCtx, &s.didChangeHandlers, doc) }
	})
}

func (s *DefaultDocumentSynchronizer) OnDidChange(handler TextDocumentChangeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.didChangeHandlers = append(s.didChangeHandlers, handler)
}

func (s *DefaultDocumentSynchronizer) DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) {
	s.write(func(writeCtx context.Context) ([]textdoc.Handle, func()) {
		textdocStore := s.sc.MustGet[textdoc.Store]()
		uri := core.NormalizeURI(params.TextDocument.URI)
		// The handlers run before the overlay is removed, so the notification
		// waits for them.
		if existing := textdocStore.GetOverlay(uri); existing != nil {
			s.notify(writeCtx, &s.didCloseHandlers, existing)
		}
		textdocStore.RemoveOverlay(uri)
		// TODO: Maybe do some more work here - a closed file might trigger a new workspace build
		// see also https://github.com/eclipse-fastbelt/fastbelt/issues/185
		return nil, nil
	})
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
