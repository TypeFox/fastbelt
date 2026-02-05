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
	document.RLock()
	root := document.Root
	document.RUnlock()
	core.TraverseNode(root, func(node core.AstNode) {
		node.ForEachReference(func(ref core.UntypedReference) {
			references = append(references, ref)
			waitgroup.Go(func() {
				ref.Resolve(ctx)
			})
		})
	})
	document.Lock()
	document.References = references
	document.Unlock()
	waitgroup.Wait()
}
