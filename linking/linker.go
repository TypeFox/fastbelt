// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"
	"sync"

	core "typefox.dev/fastbelt"
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

func NewDefaultLinker(sc *service.Container) Linker {
	return &DefaultLinker{sc: sc}
}

func (s *DefaultLinker) Link(ctx context.Context, document *core.Document) {
	waitgroup := sync.WaitGroup{}
	references := []core.UntypedReference{}
	root := document.Root
	for node := range core.AllNodes(root) {
		for ref := range core.References(node) {
			references = append(references, ref)
			waitgroup.Go(func() {
				ref.Resolve(ctx)
			})
		}
	}
	waitgroup.Wait()
	document.References = references
}
