// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
)

// DocumentUpdater manages document updates and coordinates builds.
// It sits between the document synchronization layer and the Builder,
// handling cancellation of in-progress builds when new changes arrive.
//
// Thread Safety:
// All methods are safe for concurrent use.
type DocumentUpdater interface {
	// Update processes document changes and triggers a new build.
	// Changed handles are used to create or update Documents in the DocumentManager.
	// Deleted URIs cause documents to be removed from the DocumentManager.
	// Any in-progress build is cancelled before starting a new one.
	Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI)
}

// DefaultDocumentUpdater is the default implementation of DocumentUpdater.
type DefaultDocumentUpdater struct {
	srv      WorkspaceSrvCont
	mu       sync.Mutex
	cancelFn context.CancelFunc
}

// NewDefaultDocumentUpdater creates a new default document updater.
func NewDefaultDocumentUpdater(srv WorkspaceSrvCont) DocumentUpdater {
	return &DefaultDocumentUpdater{srv: srv}
}

// Update processes document changes and triggers a new build.
func (u *DefaultDocumentUpdater) Update(ctx context.Context, changed []textdoc.Handle, deleted []core.URI) {
	docManager := u.srv.Workspace().DocumentManager

	for _, handle := range changed {
		doc := core.NewDocument(handle)
		docManager.Set(doc)
	}
	for _, uri := range deleted {
		docManager.Delete(uri)
	}

	// Cancel any in-progress build and create a new cancellable context.
	u.mu.Lock()
	if u.cancelFn != nil {
		u.cancelFn()
	}
	// Use a detached context: the incoming ctx is request-scoped and will be
	// cancelled by jsonrpc2 as soon as the notification handler returns.
	buildCtx, cancel := context.WithCancel(context.Background())
	u.cancelFn = cancel
	u.mu.Unlock()

	go func() {
		if buildCtx.Err() != nil {
			return
		}

		// TODO: Select which documents to include in the build; for now, rebuild all.
		var docs []*core.Document
		for doc := range docManager.All() {
			docs = append(docs, doc)
		}

		if err := u.srv.Workspace().Builder.Build(buildCtx, docs); err != nil {
			if buildCtx.Err() == nil {
				log.Printf("build failed: %v", err)
			}
		}
	}()
}
