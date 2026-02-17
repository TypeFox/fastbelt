// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"context"
	"iter"
	"reflect"
	"sync"

	"typefox.dev/fastbelt/util/extiter"
)

type UntypedReference interface {
	Description() *AstNodeDescription
	RefNode(ctx context.Context) AstNode
	Resolve(ctx context.Context)
	Reset()
	Error() *ReferenceError
	Segment() *TextSegment
}

type ReferenceGetter[T AstNode] func(context.Context, *Reference[T]) (*AstNodeDescription, *ReferenceError)

// Reference represents a reference to another AST node of type T.
// Resolving is thread safe and is done concurrently by default.
// The resolution is triggered when Ref(), RefNode() or Resolve() is called for the first time.
type Reference[T AstNode] struct {
	Token       *Token
	Text        string
	Owner       AstNode
	description *AstNodeDescription
	err         *ReferenceError
	ref         T
	mu          sync.Mutex
	getter      ReferenceGetter[T]
	resolved    bool
}

func (r *Reference[T]) Reset() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resolved = false
	r.description = nil
	r.err = nil
	var zero T
	r.ref = zero
}

func (r *Reference[T]) Description() *AstNodeDescription {
	if r == nil {
		return nil
	}
	return r.description
}

func (r *Reference[T]) RefNode(ctx context.Context) AstNode {
	return r.Ref(ctx)
}

func (r *Reference[T]) Ref(ctx context.Context) T {
	var zero T
	if r == nil {
		return zero
	}
	r.Resolve(ctx)
	return r.ref
}

func (r *Reference[T]) Error() *ReferenceError {
	if r == nil {
		return nil
	}
	return r.err
}

func (r *Reference[T]) Segment() *TextSegment {
	if r != nil && r.Token != nil {
		return &r.Token.Segment
	}
	return nil
}

func (r *Reference[T]) Resolve(ctx context.Context) {
	// We can use the context to detect cyclic reference resolution attempts
	// We are allowed to do this outside of the mutex lock because context is immutable
	if ctx.Value(r) != nil {
		// Note that we write directly to r.err without locking here
		// This is safe, because the reference is already locked by the caller
		// Attempting to lock it again would cause a deadlock anyway
		r.err = NewReferenceError("Cyclic reference resolution detected")
		// Return directly, do not set the resolved flag
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.resolved {
		// Already resolved
		return
	}
	newCtx := context.WithValue(ctx, r, true)
	desc, e := r.getter(newCtx, r)
	r.description = desc
	if r.err == nil {
		// Do not overwrite existing errors
		r.err = e
	}
	if desc != nil {
		if node, ok := desc.Node.(T); ok {
			r.ref = node
		} else if r.err == nil {
			expectedType := reflect.TypeFor[T]().String()
			actualType := "nil"
			if desc.Node != nil {
				actualType = reflect.TypeOf(desc.Node).String()
			}
			r.err = NewReferenceError("Reference resolution type mismatch: expected " + expectedType + ", got " + actualType)
		}
	}
	r.resolved = true
}

func NewReference[T AstNode](owner AstNode, token *Token, getter ReferenceGetter[T]) *Reference[T] {
	return &Reference[T]{
		Owner:  owner,
		Token:  token,
		Text:   token.Image,
		getter: getter,
	}
}

// Computes the reference that this token represents.
// Returns nil if the token does not represent a reference.
func ReferenceOfToken(token *Token) UntypedReference {
	owner := token.Element
	if owner == nil {
		return nil
	}
	var ref UntypedReference = nil
	// We don't have a direct reference from token -> reference, so we need to search for it
	// We have to iterate over all references of the owner node
	// This might seem inefficient, but in practice the number of references per node is usually very small
	// Also, we only do this in select LSP requests, so the performance impact is negligible
	owner.ForEachReference(func(ur UntypedReference) {
		// Simply compare the text indices to find the matching reference
		if ur.Segment().Indices == token.Segment.Indices {
			ref = ur
		}
	})
	return ref
}

type AstNodeDescription struct {
	URI         URI
	Node        AstNode
	Name        string
	NameSegment *TextSegment
	FullSegment *TextSegment
}

func NewAstNodeDescription(node AstNode, name string, nameSegment, fullSegment *TextSegment) *AstNodeDescription {
	doc := node.Document()
	return &AstNodeDescription{
		URI:         doc.URI,
		Node:        node,
		Name:        name,
		NameSegment: nameSegment,
		FullSegment: fullSegment,
	}
}

var EmptyAstNodeDescriptions = extiter.Empty[*AstNodeDescription]()

type SymbolList = iter.Seq[*AstNodeDescription]

type LocalSymbols interface {
	Has(node AstNode) bool
	Iter(node AstNode) SymbolList
}
