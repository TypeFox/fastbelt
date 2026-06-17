// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"
	"slices"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
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
	// Write cancels any previous pending or in-progress build and issues a
	// fresh context. The outer ctx is from jsonrpc2 and has a different lifetime than the build.
	go func() {
		docManager := service.MustGet[DocumentManager](s.sc)
		builder := service.MustGet[Builder](s.sc)
		lock := service.MustGet[Lock](s.sc)
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			for _, handle := range changed {
				doc := core.NewDocument(handle)
				docManager.Set(doc)
			}
			for _, uri := range deleted {
				docManager.Delete(uri)
			}

			// Collect all documents to be processed.
			// TODO implement this properly: determine the minimal set of documents to be processed
			docs := slices.Collect(docManager.All())

			// Reset documents so linking and validation are re-executed.
			keepState := core.DocStateParsed | core.DocStateExportedSymbols | core.DocStateLocalSymbols
			for _, doc := range docs {
				builder.Reset(doc, keepState)
			}

			if err := builder.Build(ctx, docs, downgrade); err != nil {
				if ctx.Err() == nil {
					log.Printf("build failed: %v", err)
				}
			}
		})
	}()
}
