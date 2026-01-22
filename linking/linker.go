package linking

import (
	"context"
	"sync"

	core "typefox.dev/fastbelt"
)

type Linker interface {
	Link(ctx context.Context, root core.AstNode)
}

type DefaultLinker struct{}

func NewDefaultLinker() Linker {
	return &DefaultLinker{}
}

func (l *DefaultLinker) Link(ctx context.Context, root core.AstNode) {
	waitgroup := sync.WaitGroup{}
	core.TraverseNode(root, func(node core.AstNode) {
		node.ForEachReference(func(ref core.UntypedReference) {
			waitgroup.Go(func() {
				ref.Resolve(ctx)
			})
		})
	})
	waitgroup.Wait()
}
