// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"context"
	"testing"
)

// A self-referencing getter must report a cyclic error instead of deadlocking,
// even when the caller passes a context without a resolution chain.
func TestResolveSelfCycleWithoutChain(t *testing.T) {
	ref := NewReference(nil, nil, func(ctx context.Context, r *Reference[AstNode]) (*SymbolDescription, *ReferenceError) {
		r.Resolve(ctx)
		return nil, nil
	})
	done := make(chan struct{})
	go func() {
		ref.Resolve(context.Background())
		close(done)
	}()
	select {
	case <-done:
	case <-t.Context().Done():
		t.Fatal("Resolve deadlocked on self cycle")
	}
	if ref.Error() == nil {
		t.Fatal("expected cyclic reference error")
	}
}
