// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"context"
	"iter"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"

	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/extiter"
)

// An untyped representation of a [Reference] that can be used when the target type is not known at compile time.
// Used throughout the fastbelt codebase to generically deal with different references.
type UntypedReference interface {
	Owner() AstNode
	Description() *SymbolDescription
	RefNode(ctx context.Context) AstNode
	Resolve(ctx context.Context)
	Reset()
	Error() *ReferenceError
	Segment() *TextSegment
	Token() *Token
	Text() string
}

type ReferenceGetter[T AstNode] func(context.Context, *Reference[T]) (*SymbolDescription, *ReferenceError)

// Reference represents a reference to another AST node of type T.
// Resolving is thread safe and is done concurrently by default.
// The resolution is triggered when [Ref], [RefNode] or [Resolve] are called for the first time.
type Reference[T AstNode] struct {
	token       *Token
	owner       AstNode
	description *SymbolDescription
	err         *ReferenceError
	ref         T
	mu          sync.Mutex
	getter      ReferenceGetter[T]
	resolved    atomic.Bool
}

func (r *Reference[T]) Reset() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resolved.Store(false)
	r.description = nil
	r.err = nil
	var zero T
	r.ref = zero
}

func (r *Reference[T]) Token() *Token {
	if r == nil {
		return nil
	}
	return r.token
}

func (r *Reference[T]) Text() string {
	if r == nil || r.token == nil {
		return ""
	}
	return r.token.Image
}

func (r *Reference[T]) Owner() AstNode {
	if r == nil {
		return nil
	}
	return r.owner
}

func (r *Reference[T]) Description() *SymbolDescription {
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
	if r != nil && r.token != nil {
		return &r.token.Segment
	}
	return nil
}

func (r *Reference[T]) Resolve(ctx context.Context) {
	// Fast path: check if already resolved without locking
	if r == nil || r.resolved.Load() {
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	r.resolveSlow(ctx)
}

func (r *Reference[T]) resolveSlow(ctx context.Context) {
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
	if r.resolved.Load() {
		// Another goroutine might have resolved it while we were waiting for the lock
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
	r.resolved.Store(true)
}

func NewReference[T AstNode](owner AstNode, token *Token, getter ReferenceGetter[T]) *Reference[T] {
	return &Reference[T]{
		owner:  owner,
		token:  token,
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
	for ur := range References(owner) {
		// Simply compare the text indices to find the matching reference
		if ur.Segment().Indices == token.Segment.Indices {
			ref = ur
			break
		}
	}
	return ref
}

type ReferenceDescription struct {
	// SourceNode is the node that contains the reference.
	SourceNode AstNode
	// TargetNode is the node that is being referenced.
	TargetNode AstNode
	// Segment is the text segment of the reference in the source document.
	Segment *TextSegment
}

func NewReferenceDescription(source, target AstNode, segment *TextSegment) *ReferenceDescription {
	return &ReferenceDescription{
		SourceNode: source,
		TargetNode: target,
		Segment:    segment,
	}
}

// SourceURI returns the URI of the document containing the source node, or nil if not available.
func (d *ReferenceDescription) SourceURI() URI {
	if d.SourceNode != nil {
		doc := d.SourceNode.Document()
		if doc != nil {
			return doc.URI
		}
	}
	return nil
}

// TargetURI returns the URI of the document containing the target node, or nil if not available.
func (d *ReferenceDescription) TargetURI() URI {
	if d.TargetNode != nil {
		doc := d.TargetNode.Document()
		if doc != nil {
			return doc.URI
		}
	}
	return nil
}

// ReferenceDescriptions is a collection of reference descriptions for a document, indexed by target node.
// It allows efficient retrieval of all references to a given target node.
type ReferenceDescriptions interface {
	// All returns an iterator over all reference descriptions in the document.
	All() iter.Seq[*ReferenceDescription]
	// Returns an iterator over all reference descriptions that point to the given target node.
	ForTarget(target AstNode) iter.Seq[*ReferenceDescription]
}

type referenceDescriptions struct {
	descriptions collections.MultiMap[AstNode, *ReferenceDescription]
}

func (d *referenceDescriptions) All() iter.Seq[*ReferenceDescription] {
	if d == nil {
		return extiter.Empty[*ReferenceDescription]()
	}
	return d.descriptions.Values()
}

func (d *referenceDescriptions) ForTarget(target AstNode) iter.Seq[*ReferenceDescription] {
	if d == nil {
		return extiter.Empty[*ReferenceDescription]()
	}
	return slices.Values(d.descriptions.Get(target))
}

func NewReferenceDescriptionsFromMap(descriptions collections.MultiMap[AstNode, *ReferenceDescription]) ReferenceDescriptions {
	return &referenceDescriptions{
		descriptions: descriptions,
	}
}
