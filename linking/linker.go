// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/parallel"
	"typefox.dev/fastbelt/util/service"
)

// Linker resolves cross-references in a document's AST.
type Linker interface {
	// Link resolves all references in the document.
	// A list of all references is stored in the document's References field.
	// The caller must hold the document's write lock.
	Link(ctx context.Context, document *core.Document)
}

// DefaultLinker is the default implementation of [Linker].
// It resolves all references in the document.
type DefaultLinker struct {
	sc *service.Container
}

// NewDefaultLinker returns a [Linker] that resolves every reference in a
// document's AST in parallel.
func NewDefaultLinker(sc *service.Container) Linker {
	return &DefaultLinker{sc: sc}
}

func (s *DefaultLinker) Link(ctx context.Context, document *core.Document) {
	// Each worker resolves references with its own resolution chain context,
	// so cycle detection does not need to allocate per resolved reference.
	parallel.ForEachWithSetup(document.References,
		func() context.Context { return core.WithResolutionChain(ctx) },
		func(workerCtx context.Context, ref core.UntypedReference, _ int) {
			if workerCtx.Err() != nil {
				return
			}
			ref.Resolve(workerCtx)
		})
}
