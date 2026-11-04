// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"context"
	"testing"
	"time"
	"unsafe"
)

// A self-referencing getter must report a cyclic error instead of deadlocking,
// especially if the caller passes the correct chain.
func TestResolveSelfCycleWithChain(t *testing.T) {
	ref := NewReference(nil, nil, func(ctx context.Context, r *Reference[AstNode]) (*SymbolDescription, *ReferenceError) {
		r.Resolve(ctx)
		return nil, nil
	})
	done := make(chan struct{})
	go func() {
		ref.Resolve(WithResolutionChain(context.Background()))
		close(done)
	}()
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("Resolve deadlocked on self cycle")
	}
	if ref.Error() == nil {
		t.Fatal("expected cyclic reference error")
	}
}

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
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("Resolve deadlocked on self cycle")
	}
	if ref.Error() == nil {
		t.Fatal("expected cyclic reference error")
	}
}

// Reference is a large share of the heap, so its layout is pinned to the 48-byte size class.
// Note: If we ever get it smaller, this test will fail - pin it to the new size.
func TestReferenceSize(t *testing.T) {
	if size := unsafe.Sizeof(Reference[AstNode]{}); size != 48 {
		t.Fatalf("Reference grew to %d bytes", size)
	}
}

// A nil description and a resolved-but-unresolvable reference must both read as "no symbol",
// while Resolve still runs exactly once and Reset re-arms it.
func TestReferenceUnresolvedAndReset(t *testing.T) {
	calls := 0
	owner := &AstNodeBase{}
	ref := NewReference(owner, nil, func(ctx context.Context, r *Reference[AstNode]) (*SymbolDescription, *ReferenceError) {
		calls++
		return nil, NewReferenceError("nope")
	})
	if ref.Owner() != AstNode(owner) {
		t.Fatal("owner not derived from synthetic unit")
	}
	ref.Resolve(context.Background())
	ref.Resolve(context.Background())
	if calls != 1 || ref.Description() != nil || ref.Ref(context.Background()) != nil || ref.Error() == nil {
		t.Fatalf("calls=%d desc=%v err=%v", calls, ref.Description(), ref.Error())
	}
	ref.Reset()
	ref.Resolve(context.Background())
	if calls != 2 {
		t.Fatalf("Reset did not re-arm resolution, calls=%d", calls)
	}
}

// A description whose node has the wrong type yields the zero target and a mismatch error.
func TestReferenceTypeMismatch(t *testing.T) {
	ref := NewReference(nil, nil, func(ctx context.Context, r *Reference[CompositeNode]) (*SymbolDescription, *ReferenceError) {
		return &SymbolDescription{Node: &AstNodeBase{}}, nil
	})
	if ref.Ref(context.Background()) != nil || ref.Description() == nil || ref.Error() == nil {
		t.Fatalf("desc=%v err=%v", ref.Description(), ref.Error())
	}
}
