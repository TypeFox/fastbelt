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
)

// Builder is the interface for building workspace-related structures.
type Builder interface {
	// Build processes the provided documents through all build phases (parse, compute
	// symbol table, link). It should regularly check ctx for cancellation between phases.
	Build(ctx context.Context, docs []*core.Document) error
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

// Build processes the provided documents through all build phases.
func (b *DefaultBuilder) Build(ctx context.Context, docs []*core.Document) error {
	if b.srv == nil {
		return nil
	}
	parser := b.srv.Workspace().DocumentParser
	exportedSymbols := b.srv.Linking().ExportedSymbolsProvider
	localSymbols := b.srv.Linking().LocalSymbolTableProvider
	linker := b.srv.Linking().Linker

	results := make([]ValidationResult, 0, len(docs))
	for _, document := range docs {
		if err := ctx.Err(); err != nil {
			return err
		}
		document.Lock()

		// PHASE 1: Parse the text content into an AST
		parser.Parse(document)
		if err := ctx.Err(); err != nil {
			document.Unlock()
			return err
		}

		// PHASE 2: Compute exported symbols
		exportedSymbols.Compute(ctx, document)
		if err := ctx.Err(); err != nil {
			document.Unlock()
			return err
		}

		// PHASE 3: Compute the local symbol table
		localSymbols.Compute(ctx, document)
		if err := ctx.Err(); err != nil {
			document.Unlock()
			return err
		}

		// PHASE 4: Link the AST to resolve cross-references
		linker.Link(ctx, document)
		document.Unlock()

		// PHASE 5: Validate the document
		// TODO: implement custom validation for specific languages (with read-only lock)
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
		if err := ctx.Err(); err != nil {
			return err
		}
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
