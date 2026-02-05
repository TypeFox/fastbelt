// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"iter"
	"maps"
	"sync"

	"github.com/TypeFox/go-lsp/protocol"
	core "typefox.dev/fastbelt"
)

type DocumentManager interface {
	// Checks if a document with the given URI exists.
	//
	// This method is thread-safe.
	Has(uri protocol.DocumentURI) bool
	// Retrieves the document for the given URI, or nil if it does not exist.
	//
	// This method is thread-safe.
	Get(uri protocol.DocumentURI) *core.Document
	// Adds or updates the given document.
	//
	// This method is thread-safe.
	Set(document *core.Document)
	// Returns a sequence of all managed documents.
	//
	// This method does not lock the document manager for the duration of the iteration.
	// Care should be taken when using the returned sequence in concurrent scenarios.
	All() iter.Seq[*core.Document]
	// Deletes the document with the given URI and returns it, or nil if it did not exist.
	//
	// This method is thread-safe.
	Delete(uri protocol.DocumentURI) *core.Document
}

type DefaultDocumentManager struct {
	mu        sync.RWMutex
	documents map[protocol.DocumentURI]*core.Document
}

func NewDefaultDocumentManager() DocumentManager {
	return &DefaultDocumentManager{
		documents: map[protocol.DocumentURI]*core.Document{},
	}
}

func (d *DefaultDocumentManager) Has(uri protocol.DocumentURI) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.documents[uri]
	return exists
}

func (d *DefaultDocumentManager) Get(uri protocol.DocumentURI) *core.Document {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.documents[uri]
}

func (d *DefaultDocumentManager) Set(document *core.Document) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.documents[document.URI()] = document
}

func (d *DefaultDocumentManager) All() iter.Seq[*core.Document] {
	return maps.Values(d.documents)
}

func (d *DefaultDocumentManager) Delete(uri protocol.DocumentURI) *core.Document {
	d.mu.Lock()
	defer d.mu.Unlock()
	document, exists := d.documents[uri]
	if exists {
		delete(d.documents, uri)
		return document
	} else {
		return nil
	}
}
