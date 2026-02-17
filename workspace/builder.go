// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"
	"reflect"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
)

// Builder is the interface for building workspace-related structures.
type Builder interface {
	// Update updates the workspace based on the provided documents.
	Update(ctx context.Context, docs []textdoc.Handle) error
	// AddValidationListener registers a listener that will be called when validation completes.
	AddValidationListener(listener ValidationListener)
	// RemoveValidationListener unregisters a previously registered listener.
	RemoveValidationListener(listener ValidationListener)
}

// ValidationResult contains the result of validating a document.
type ValidationResult struct {
	Document *core.Document
}

// ValidationListener is a function that is called when validation completes for a set of documents.
// It receives the context and all validation results (one per document).
// If the listener returns an error, it will be logged but will not prevent other listeners from being called.
type ValidationListener func(ctx context.Context, results []ValidationResult) error

// DefaultBuilder is the default implementation of Builder.
type DefaultBuilder struct {
	srv         WorkspaceSrvCont
	listeners   []ValidationListener
	listenersMu sync.RWMutex
}

// NewDefaultBuilder creates a new default builder.
func NewDefaultBuilder(srv WorkspaceSrvCont) Builder {
	return &DefaultBuilder{
		srv:       srv,
		listeners: make([]ValidationListener, 0),
	}
}

// Update updates the workspace based on the provided documents.
func (b *DefaultBuilder) Update(ctx context.Context, textDocs []textdoc.Handle) error {
	if b.srv == nil {
		return nil
	}
	// TODO: Do we need to check whether all services are set?
	docManager := b.srv.Workspace().DocumentManager
	parser := b.srv.Workspace().DocumentParser
	symbolTableProvider := b.srv.Linking().LocalSymbolTableProvider
	linker := b.srv.Linking().Linker

	// Parse all documents and collect validation results
	results := make([]ValidationResult, 0, len(textDocs))
	for _, textDoc := range textDocs {
		document := core.NewDocument(textDoc)
		docManager.Set(document)
		parser.Parse(document)
		symbolTableProvider.Compute(ctx, document)
		linker.Link(ctx, document)
		results = append(results, ValidationResult{
			Document: document,
		})
	}

	// Notify all registered listeners
	b.listenersMu.RLock()
	listeners := make([]ValidationListener, len(b.listeners))
	copy(listeners, b.listeners)
	b.listenersMu.RUnlock()

	for _, listener := range listeners {
		if err := listener(ctx, results); err != nil {
			log.Printf("validation listener error: %v", err)
		}
	}

	return nil
}

// AddValidationListener registers a listener that will be called when validation completes.
func (b *DefaultBuilder) AddValidationListener(listener ValidationListener) {
	if listener == nil {
		return
	}
	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()
	b.listeners = append(b.listeners, listener)
}

// RemoveValidationListener unregisters a previously registered listener.
func (b *DefaultBuilder) RemoveValidationListener(listener ValidationListener) {
	if listener == nil {
		return
	}
	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()
	listenerPtr := reflect.ValueOf(listener).Pointer()
	for i, l := range b.listeners {
		if reflect.ValueOf(l).Pointer() == listenerPtr {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			return
		}
	}
}
