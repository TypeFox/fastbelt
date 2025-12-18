package fastbelt

import (
	"context"
	"reflect"
	"sync"
)

type UntypedReference interface {
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
	Description *AstNodeDescription
	err         *ReferenceError
	ref         T
	mu          sync.Mutex
	getter      ReferenceGetter[T]
	resolved    bool
}

func (r *Reference[T]) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.resolved = false
	r.Description = nil
	r.err = nil
	var zero T
	r.ref = zero
}

func (r *Reference[T]) RefNode(ctx context.Context) AstNode {
	return r.Ref(ctx)
}

func (r *Reference[T]) Ref(ctx context.Context) T {
	r.Resolve(ctx)
	return r.ref
}

func (r *Reference[T]) Error() *ReferenceError {
	return r.err
}

func (r *Reference[T]) Segment() *TextSegment {
	if r != nil && r.Token != nil {
		return &r.Token.Segment
	}
	return nil
}

type getterResult struct {
	desc *AstNodeDescription
	err  *ReferenceError
}

func (r *Reference[T]) Resolve(ctx context.Context) {
	// We can use the context to detect cyclic reference resolution attempts
	// We are allowed to do this outside of the mutex lock because context is immutable
	if ctx.Value(r) != nil {
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
	ch := make(chan getterResult)
	go func() {
		desc, e := r.getter(newCtx, r)
		ch <- getterResult{desc, e}
	}()
	result := <-ch
	desc := result.desc
	r.Description = desc
	r.err = result.err
	if desc != nil {
		if node, ok := desc.Node.(T); ok {
			r.ref = node
		} else {
			expectedType := reflect.TypeOf(*new(T)).String()
			actualType := reflect.TypeOf(desc.Node).String()
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

type ReferenceError struct {
	Msg      string
	Severity int
}

func (e *ReferenceError) Error() string {
	return e.Msg
}

func NewReferenceError(msg string) *ReferenceError {
	return &ReferenceError{Msg: msg, Severity: 1}
}

type AstNodeDescription struct {
	Node        AstNode
	Name        string
	NameSegment *TextSegment
	FullSegment *TextSegment
}

func NewAstNodeDescription(node AstNode, name string, nameSegment, fullSegment *TextSegment) *AstNodeDescription {
	return &AstNodeDescription{
		Node:        node,
		Name:        name,
		NameSegment: nameSegment,
		FullSegment: fullSegment,
	}
}
