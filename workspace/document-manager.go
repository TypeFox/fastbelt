// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"iter"
	"sync"
	"sync/atomic"

	core "typefox.dev/fastbelt"
)

// All methods of DocumentManager are thread-safe and can be called concurrently.
type DocumentManager interface {
	// Checks if a document with the given URI exists.
	Has(uri core.URI) bool
	// Retrieves the document for the given URI, or nil if it does not exist.
	Get(uri core.URI) *core.Document
	// Adds or updates the given document.
	Set(document *core.Document)
	// Returns a sequence of all managed documents.
	All() iter.Seq[*core.Document]
	// Deletes the document with the given URI and returns it, or nil if it did not exist.
	Delete(uri core.URI) *core.Document
}

type DefaultDocumentManager struct {
	documents sync.Map
}

func NewDefaultDocumentManager() DocumentManager {
	return &DefaultDocumentManager{
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
