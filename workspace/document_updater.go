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
// rebuilds. It sits between the LSP document sync layer and [Builder],
// serializing mutations through [Lock.Write] and cancelling in-progress builds
// when new changes arrive.
//
// All methods are safe for concurrent use.
type DocumentUpdater interface {
	// Update creates or updates [core.Document] values for changed handles,
	// removes deleted URIs from [DocumentManager], and starts a new build.
	// It returns immediately; the build runs asynchronously in a background
	// goroutine. Any in-progress build is cancelled before the new one starts.
	Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI)
}

// DefaultDocumentUpdater is the default implementation of [DocumentUpdater].
type DefaultDocumentUpdater struct {
	sc *service.Container
}

// NewDefaultDocumentUpdater returns a [DocumentUpdater] that coordinates
// [DocumentManager], [Builder], and [Lock].
func NewDefaultDocumentUpdater(sc *service.Container) DocumentUpdater {
	return &DefaultDocumentUpdater{sc: sc}
}

func (s *DefaultDocumentUpdater) Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI) {
	go func() {
		docManager := service.MustGet[DocumentManager](s.sc)
		changeImpact := service.MustGet[DocumentChangeImpact](s.sc)
		builder := service.MustGet[Builder](s.sc)
		lock := service.MustGet[Lock](s.sc)
		// Write cancels any previous pending or in-progress build and issues a
		// fresh context. The outer ctx is from jsonrpc2 and has a different
		// lifetime than the build.
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
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
			keepState := core.DocStateParsed | core.DocStateExportedSymbols | core.DocStateLocalSymbols
			docs := make([]*core.Document, 0, len(changed)+10)
			for doc := range docManager.All() {
				if !changedURIs.Has(doc.URI.StringUnencoded()) && changeImpact.Affected(doc, changedURIs) {
					builder.Reset(doc, keepState)
				}
				if !doc.State.IsComplete() {
					docs = append(docs, doc)
				}
			}

			if err := builder.Build(ctx, docs, downgrade); err != nil {
				if ctx.Err() == nil {
					log.Printf("build failed: %v", err)
				}
			}
		})
	}()
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
// [core.Document.ReferenceDescriptions].
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
	for _, ref := range doc.References {
		if ref.Error() != nil {
			return true
		}
	}

	if doc.ReferenceDescriptions == nil {
		return false
	}
	docURI := doc.URI.StringUnencoded()
	for desc := range doc.ReferenceDescriptions.All() {
		targetURI := desc.TargetURI()
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
