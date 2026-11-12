// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/service"
)

// DocumentUpdater applies text-document changes to the workspace and triggers
// rebuilds. It sits between the LSP document sync layer and [Builder].
//
// The updater does not acquire the [Lock] itself: callers serialize updates
// by calling it from within a [Lock.Write] task, which also cancels
// in-progress builds when new changes arrive.
type DocumentUpdater interface {
	// Update creates or updates [core.Document] values for changed handles,
	// removes deleted URIs from [DocumentManager], and rebuilds the workspace.
	// It runs synchronously and returns once the build has finished or ctx
	// was cancelled.
	//
	// Update must be called under [Lock.Write], with the context passed to
	// the write task as ctx. The document mutations are always applied; only
	// the build is cut short when ctx is cancelled by a newer write, and that
	// write is then responsible for building the documents left incomplete.
	Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI)
}

// DefaultDocumentUpdater is the default implementation of [DocumentUpdater].
type DefaultDocumentUpdater struct {
	sc *service.Container
}

// NewDefaultDocumentUpdater returns a [DocumentUpdater] that coordinates
// [DocumentManager], [DocumentChangeImpact], and [Builder].
func NewDefaultDocumentUpdater(sc *service.Container) DocumentUpdater {
	return &DefaultDocumentUpdater{sc: sc}
}

func (s *DefaultDocumentUpdater) Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI) {
	docManager := s.sc.MustGet[DocumentManager]()
	changeImpact := s.sc.MustGet[DocumentChangeImpact]()
	builder := s.sc.MustGet[Builder]()
	// The caller holds the write lock, so this is the mutation phase: documents
	// may be replaced, deleted, and reset until the build starts below.
	changedURIs := make(collections.Set[string], len(changed)+len(deleted))
	for _, handle := range changed {
		doc := core.NewDocument(handle)
		docManager.Set(doc)
		changedURIs.Add(doc.URI.StringUnencoded())
	}
	for _, uri := range deleted {
		docManager.Delete(uri)
		changedURIs.Add(uri.StringUnencoded())
	}

	// Reset documents that depend on a changed document so their references
	// are resolved again, then collect every document that has not yet
	// completed all build phases. This covers the changed documents, the
	// documents reset above, documents left incomplete by a previously
	// cancelled build, and documents that have not been built at all
	// (e.g. after the initial workspace load). Relinked documents keep
	// their parsed AST, exported symbols, and local symbols because the
	// documents themselves did not change.
	keepState := core.DocStateLexed |
		core.DocStateParsed |
		core.DocStateExportedSymbols |
		core.DocStateLocalSymbols
	docs := make([]*core.Document, 0, len(changed)+10)
	for doc := range docManager.All() {
		if !changedURIs.Has(doc.URI.StringUnencoded()) && changeImpact.Affected(doc, changedURIs) {
			builder.Reset(doc, keepState)
		}
		if !doc.State().IsComplete() {
			docs = append(docs, doc)
		}
	}

	if err := builder.Build(ctx, docs); err != nil {
		if ctx.Err() == nil {
			log.Printf("build failed: %v", err)
		}
	}
}

// DocumentChangeImpact reports whether a document should be relinked in response
// to changes in other documents.
type DocumentChangeImpact interface {
	// Affected reports whether doc is affected by changedURIs and thus should be
	// relinked.
	Affected(doc *core.Document, changedURIs collections.Set[string]) bool
}

// DefaultDocumentChangeImpact is the default implementation of [DocumentChangeImpact].
type DefaultDocumentChangeImpact struct {
	sc *service.Container
}

// NewDefaultDocumentChangeImpact returns a [DocumentChangeImpact] that inspects
// [core.Document.References].
func NewDefaultDocumentChangeImpact(sc *service.Container) DocumentChangeImpact {
	return &DefaultDocumentChangeImpact{sc: sc}
}

// Affected reports whether doc contains a cross-document reference to any of the
// documents identified by changedURIs. URIs are compared using
// [core.URI.StringUnencoded].
//
// References to nodes within doc itself are ignored: such local references cannot
// be invalidated by a change to another document.
//
// Documents with unresolved references are also considered affected, since they
// might resolve against a symbol introduced by the change.
func (s *DefaultDocumentChangeImpact) Affected(doc *core.Document, changedURIs collections.Set[string]) bool {
	docURI := doc.URI.StringUnencoded()
	for _, ref := range doc.References {
		targetDescription := ref.Description()
		if targetDescription == nil || ref.Error() != nil {
			// Failed to resolve, or an error, so relink document
			return true
		}
		targetURI := targetDescription.URI
		if targetURI == nil {
			continue
		}
		target := targetURI.StringUnencoded()
		if target == docURI {
			continue
		}
		if changedURIs.Has(target) {
			return true
		}
	}
	return false
}
