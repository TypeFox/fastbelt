// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"
	"sync"

	core "typefox.dev/fastbelt"
)

type Linker interface {
	Link(ctx context.Context, document *core.Document)
}

type DefaultLinker struct{}

func NewDefaultLinker() Linker {
	return &DefaultLinker{}
}

func (l *DefaultLinker) Link(ctx context.Context, document *core.Document) {
	waitgroup := sync.WaitGroup{}
	references := []core.UntypedReference{}
	root := document.Root
	core.TraverseNode(root, func(node core.AstNode) {
		node.ForEachReference(func(ref core.UntypedReference) {
			references = append(references, ref)
			waitgroup.Go(func() {
				ref.Resolve(ctx)
			})
		})
	})
	waitgroup.Wait()
	document.References = references
}
