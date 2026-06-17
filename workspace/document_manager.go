// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"iter"
	"sync"
	"sync/atomic"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// DocumentManager holds all known documents of a workspace in memory.
// All methods are safe for concurrent use.
type DocumentManager interface {
	// Has reports whether a document with the given URI is managed.
	Has(uri core.URI) bool
	// Get returns the document for uri, or nil if it is not managed.
	Get(uri core.URI) *core.Document
	// Set adds or replaces a document keyed by document.URI.
	Set(document *core.Document)
	// All returns every managed document. The sequence may be iterated only once.
	All() iter.Seq[*core.Document]
	// Delete removes the document for uri and returns it, or nil if it was not managed.
	Delete(uri core.URI) *core.Document
	// Clear removes every document from the manager.
	Clear()
}

// DefaultDocumentManager is the default implementation of [DocumentManager].
type DefaultDocumentManager struct {
	sc        *service.Container
	documents sync.Map
}

// NewDefaultDocumentManager returns a [DocumentManager] backed by a concurrent map.
func NewDefaultDocumentManager(sc *service.Container) DocumentManager {
	return &DefaultDocumentManager{
		sc:        sc,
		documents: sync.Map{},
	}
}

func (d *DefaultDocumentManager) Has(uri core.URI) bool {
	_, exists := d.documents.Load(uri.StringUnencoded())
	return exists
}

func (d *DefaultDocumentManager) Get(uri core.URI) *core.Document {
	value, _ := d.documents.Load(uri.StringUnencoded())
	if doc, ok := value.(*core.Document); ok {
		return doc
	}
	return nil
}

func (d *DefaultDocumentManager) Set(document *core.Document) {
	d.documents.Store(document.URI.StringUnencoded(), document)
}

func (d *DefaultDocumentManager) All() iter.Seq[*core.Document] {
	return func(yield func(*core.Document) bool) {
		stop := atomic.Bool{}
		d.documents.Range(func(key, value any) bool {
			if stop.Load() {
				return false
			}
			if doc, ok := value.(*core.Document); ok {
				if !yield(doc) {
					stop.Store(true)
					return false
				}
			}
			return true
		})
	}
}

func (d *DefaultDocumentManager) Delete(uri core.URI) *core.Document {
	document, _ := d.documents.LoadAndDelete(uri.StringUnencoded())
	if doc, ok := document.(*core.Document); ok {
		return doc
	}
	return nil
}

func (d *DefaultDocumentManager) Clear() {
	d.documents.Clear()
}
