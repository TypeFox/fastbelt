// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"errors"
	"iter"
	"strings"
	"sync"
	"sync/atomic"
	"unique"

	"typefox.dev/fastbelt/util/parallel"
)

// AstNode is the base interface for all AST nodes.
//
// Every language-specific AST node type which is generated from a grammar definition embeds
// this interface.
type AstNode interface {
	// Document returns the owning document of the node.
	Document() *Document
	// SetDocument sets the owning document of the node.
	//
	// When constructing an AST programmatically, use [AssignContainers] to link the node in the AST.
	SetDocument(document *Document)
	// Container returns the direct parent node of the node in the AST.
	// It returns nil if this is the root node.
	Container() AstNode
	// ContainmentData returns a [unique.Handle] denoting the containing property within its [AstNode.Container],
	// defaults to a [unique.Handle] of the empty string,
	// and the element index within the containing property, defaults to -1 for single item fields.
	ContainmentData() (unique.Handle[string], int)
	// SetContainer sets the direct parent node of the node.
	//
	// When constructing an AST programmatically, use [AssignContainers] to link the node in the AST.
	SetContainer(container AstNode, containerField unique.Handle[string], index int)
	// Tokens returns the tokens associated with the node.
	Tokens() []*Token
	// AppendToken appends token to the node's token list.
	AppendToken(token *Token)
	// SetTokens replaces the node's token list with another list of tokens.
	SetTokens(tokens []*Token)
	// TextRange returns the text range of the node.
	TextRange() TextRange
	// SetTextRange sets the full text range metadata of the node.
	//
	// It is primarily used by generated parsers while constructing nodes incrementally.
	SetTextRange(r TextRange)
	// SetTextRangeStart sets the start of the node's range.
	//
	// It is primarily used by generated parsers while constructing nodes incrementally.
	SetTextRangeStart(start int32)
	// SetTextRangeEnd sets the end of the node's range.
	//
	// It is primarily used by generated parsers while constructing nodes incrementally.
	SetTextRangeEnd(end int32)
	// Text returns the source substring covered by the node's range.
	Text() string
	// ForEachNode calls fn for each direct child node of node.
	//
	// Note that this does not traverse the entire subtree. Use [AllNodes] or [AllChildren] for that.
	//
	// Calling this method directly is not recommended. Use [ChildNodes] instead for better readability.
	ForEachNode(fn func(AstNode, unique.Handle[string], int))
	// ForEachReference calls fn for each reference field of node.
	//
	// Calling this method directly is not recommended. Use [References] instead for better readability.
	ForEachReference(fn func(UntypedReference, unique.Handle[string], int))
	// Resolve returns a (nested) child node denoted by the given (relative) fragment path descriptor.
	//
	// Calling this method directly is not recommended. Use [Resolve] instead.
	Resolve(path FragmentPath) (AstNode, error)
}

// AstNodeBase provides the default [AstNode] implementation used by generated AST node types.
type AstNodeBase struct {
	document       *Document
	container      AstNode
	containerField unique.Handle[string]
	containerIndex int
	tokens         []*Token
	// tokenBuf is a small preallocated buffer to avoid heap allocations for nodes with few tokens.
	// It is used as the backing array for the tokens slice when the node has 4 or fewer tokens.
	// Massively reduces the amount of allocations required for most languages, which increases
	// the parsing speed by roughly 25% in benchmarks.
	tokenBuf [4]*Token
	rng      TextRange
}

// Document returns the owning document of the node.
func (node *AstNodeBase) Document() *Document {
	if node != nil {
		return node.document
	} else {
		return nil
	}
}

// SetDocument sets the owning document of the node.
func (node *AstNodeBase) SetDocument(document *Document) {
	if node != nil {
		node.document = document
	}
}

// Container returns the direct parent node of the node in the AST.
// It returns nil if this is the root node.
func (node *AstNodeBase) Container() AstNode {
	if node != nil {
		return node.container
	} else {
		return nil
	}
}

func (node *AstNodeBase) ContainmentData() (unique.Handle[string], int) {
	return node.containerField, node.containerIndex
}

// TODO: If concrete methods gain access to generics, refactor this into a method
// See https://github.com/golang/go/issues/77273

// ContainerOfType walks up node's container chain and returns the first ancestor assignable to T.
func ContainerOfType[T AstNode](node AstNode) T {
	var zero T
	if node == nil {
		return zero
	}
	current := node.Container()
	for current != nil {
		if casted, ok := current.(T); ok {
			return casted
		}
		current = current.Container()
	}
	return zero
}

// SetContainer sets the direct parent node of the node.
func (node *AstNodeBase) SetContainer(container AstNode, field unique.Handle[string], index int) {
	if node != nil {
		node.container = container
		node.containerField = field
		node.containerIndex = index
	}
}

// Tokens returns the tokens associated with the node.
func (node *AstNodeBase) Tokens() []*Token {
	if node != nil {
		return node.tokens
	} else {
		return nil
	}
}

// SetRangeStartToken sets the start of the node's range from token.
func (node *AstNodeBase) SetTextRangeStart(start int32) {
	node.rng.Start = start
}

// SetRangeEndToken sets the end of the node's range from token.
func (node *AstNodeBase) SetTextRangeEnd(end int32) {
	node.rng.End = end
}

// SetRange sets the full text range of the node.
func (node *AstNodeBase) SetTextRange(rng TextRange) {
	if node != nil {
		node.rng = rng
	}
}

// Range returns the text range of the node.
func (node *AstNodeBase) TextRange() TextRange {
	if node != nil {
		return node.rng
	} else {
		return TextRange{}
	}
}

// AppendToken appends token to the node's token list.
func (node *AstNodeBase) AppendToken(token *Token) {
	if node != nil && token != nil {
		if node.tokens == nil {
			node.tokens = node.tokenBuf[:0]
		}
		node.tokens = append(node.tokens, token)
	}
}

// SetTokens replaces the node's token list with tokens.
func (node *AstNodeBase) SetTokens(tokens []*Token) {
	if node != nil {
		// The method is called to set all tokens of the node at once
		// The old node is discarded in the process
		// Therefore, we don't append but replace the token slice
		node.tokens = tokens
	}
}

// Text returns the source substring covered by the node's range.
func (node *AstNodeBase) Text() string {
	if node == nil || node.document == nil || node.document.TextDoc == nil {
		return ""
	} else {
		fullText := node.document.TextDoc.Text(nil)
		return fullText[node.rng.Start:node.rng.End]
	}
}

// ForEachNode calls fn for each direct child node of node.
//
// ForEachNode on AstNodeBase is a no-op because the base type has no child fields.
func (node *AstNodeBase) ForEachNode(fn func(AstNode, unique.Handle[string], int)) {
	// This base implementation does not have any contained nodes.
}

// ForEachReference calls fn for each reference field of node.
//
// ForEachReference on AstNodeBase is a no-op because the base type has no reference fields.
func (node *AstNodeBase) ForEachReference(fn func(UntypedReference, unique.Handle[string], int)) {
	// This base implementation does not have any references.
}

// Base Implementation for instances of [AstNodeBase].
// The generator produces specific override methods for each generated ...Impl type.
func (node *AstNodeBase) Resolve(path FragmentPath) (AstNode, error) {
	return nil, errors.New("AstNodeBase.Resolve: Cannot identify children of plain AstNodeBase instances")
}

// Performance note about traversal function:
// Theoretically, we could have ChildNodes and References directly as methods on the AstNode interface.
// However, implementing the deep traversal on top of an iter.Seq is very inefficient.
// In benchmarks, it is roughly 5x slower than the current implementation.
// By using a callback-based approach, we can traverse the entire subtree with minimal overhead.
// But we lose the ability to short-circuit the traversal when we find what we're looking for.
// In practice, this is not a big issue, because most traversals will need to visit most of the nodes anyway.
// AllNodes and AllChildren are slightly less efficient than traverseContent,
// but only by roughly 10%, and they provide a much nicer API for most use cases, so the trade-off is worth it.

// Traverses all children of the given node, calling the specified function for each child.
// Does not call the function for the given node itself.
//
// Note that this function will traverse the entire subtree, without short-circuiting.
func traverseContent(node AstNode, fn func(AstNode)) {
	node.ForEachNode(func(child AstNode, containerField unique.Handle[string], index int) {
		fn(child)
		traverseContent(child, fn)
	})
}

// AllNodes creates an iterator over the given node and all its descendant nodes.
//
// Early loop exit is honored correctly, but does not short-circuit the traversal.
func AllNodes(node AstNode) iter.Seq[AstNode] {
	return func(yield func(AstNode) bool) {
		if !yield(node) {
			return
		}
		stopped := false
		traverseContent(node, func(n AstNode) {
			if !stopped && !yield(n) {
				stopped = true
			}
		})
	}
}

// AllChildren creates an iterator over all descendant nodes of the given node, excluding the node itself.
//
// Early loop exit is honored correctly, but does not short-circuit the traversal.
func AllChildren(node AstNode) iter.Seq[AstNode] {
	return func(yield func(AstNode) bool) {
		stopped := false
		traverseContent(node, func(n AstNode) {
			if !stopped && !yield(n) {
				stopped = true
			}
		})
	}
}

// ChildNodes creates an iterator over the direct child nodes of the given node.
//
// This function wraps [AstNode.ForEachNode] in an [iter.Seq].
// Early loop exit is honored correctly, but does not short-circuit the traversal.
func ChildNodes(node AstNode) iter.Seq[AstNode] {
	return func(yield func(AstNode) bool) {
		stopped := false
		node.ForEachNode(func(child AstNode, containerField unique.Handle[string], index int) {
			if !stopped && !yield(child) {
				stopped = true
			}
		})
	}
}

// References creates an iterator over all references of the given node.
//
// This function wraps [AstNode.ForEachReference] in an [iter.Seq].
// Early loop exit is honored correctly, but does not short-circuit the traversal.
func References(node AstNode) iter.Seq[UntypedReference] {
	return func(yield func(UntypedReference) bool) {
		stopped := false
		node.ForEachReference(func(ref UntypedReference, containerField unique.Handle[string], index int) {
			if !stopped && !yield(ref) {
				stopped = true
			}
		})
	}
}

// AssignToken appends token to node and records node and kind on the token.
//
// It is primarily used by generated parsers while constructing nodes incrementally.
func AssignToken(node AstNode, token *Token, kind int) {
	if node != nil && token != nil {
		node.AppendToken(token)
		token.Element = node
		token.Kind = kind
	}
}

// AssignTokens replaces node tokens and records node as owner for each token.
//
// It is primarily used by generated parsers while constructing nodes incrementally.
func AssignTokens(node AstNode, tokens []*Token) {
	if node != nil {
		node.SetTokens(tokens)
		for _, token := range tokens {
			token.Element = node
		}
	}
}

// MergeTokens prepends oldTokens to newNode's existing token list.
//
// It is used when parser actions replace the current node while preserving already consumed text.
func MergeTokens(newNode AstNode, oldTokens []*Token) {
	if newNode != nil && len(oldTokens) > 0 {
		// Prepend old tokens to the new node's tokens. The full slice expression
		// forces append to copy, so the caller's backing array is never mutated.
		AssignTokens(newNode, append(oldTokens[:len(oldTokens):len(oldTokens)], newNode.Tokens()...))
	}
}

// Allocate a new reference slot for every ~10 tokens on average.
// This average is updated after each traversal to adapt to the actual language.
const defaultReferenceRatio = 1.0 / 10.0

// running exponential moving average of references-per-token for each language ID
// different languages may have different average reference ratios
var avgReferenceRatioMap sync.Map

func getAvgReferenceRatio(languageId string) *parallel.RunningAverage {
	if avg, ok := avgReferenceRatioMap.Load(languageId); ok {
		return avg.(*parallel.RunningAverage)
	} else {
		avg := parallel.NewRunningAverage(defaultReferenceRatio)
		avgReferenceRatioMap.Store(languageId, avg)
		return avg
	}
}

// AssignContainers recursively assigns document and parent pointers for the root node and its subtree.
//
// It also assigns document and container on composite reference units reachable via references.
// It will also fill the [Document.References] field with all references found in the subtree.
func AssignContainers(doc *Document) {
	languageId := doc.TextDoc.LanguageID()
	avgReferenceRatio := getAvgReferenceRatio(languageId)
	// Continually increasing the capacity of the references slice is expensive
	// So we use a running average of references-per-token to preallocate the slice
	references := make([]UntypedReference, 0, avgReferenceRatio.Capacity(len(doc.Tokens)))
	doAssignContainers(doc, doc.Root, &references)
	doc.References = references

	if len(doc.Tokens) > 0 {
		avgReferenceRatio.Update(float64(len(references)) / float64(len(doc.Tokens)))
	}
}

func doAssignContainers(doc *Document, root AstNode, references *[]UntypedReference) {
	root.SetDocument(doc)
	root.ForEachNode(func(child AstNode, containerField unique.Handle[string], index int) {
		child.SetDocument(doc)
		child.SetContainer(root, containerField, index)
		doAssignContainers(doc, child, references)
	})
	root.ForEachReference(func(ur UntypedReference, containerField unique.Handle[string], index int) {
		*references = append(*references, ur)
		unit := ur.Unit()
		if stringNode, ok := unit.(CompositeNode); ok {
			stringNode.SetDocument(doc)
			stringNode.SetContainer(root, containerField, index)
		}
	})
}

// NamedNode represents an [AstNode] whose name is accessible as a string in the Name field.
type NamedNode interface {
	AstNode
	// Name returns the name of this node as a string.
	Name() string
}

// NamedTokenNode represents a [NamedNode] whose name is represented by a [Token], stored in
// the "Name" field of the node.
type NamedTokenNode interface {
	NamedNode
	// NameToken returns the token stored in the node's "Name" field.
	NameToken() *Token
}

// NamedCompositeNode represents a [NamedNode] whose name is represented by a [CompositeNode],
// stored in the "Name" field of the node.
type NamedCompositeNode interface {
	NamedNode
	// NameNode returns the composite node stored in the node's "Name" field.
	NameNode() CompositeNode
}

// StringUnit is a common interface for both [Token] and [CompositeNode].
type StringUnit interface {
	// Owner returns the AST node that owns this string unit.
	Owner() AstNode
	// TextRange returns the text range of this string unit.
	TextRange() TextRange
	// String returns the string representation of this string unit.
	String() string
}

// CompositeNode represents a composed string value that is made up of multiple tokens.
//
// A common example for this is a fully qualified name that consists of multiple identifiers
// and dots, e.g. "a.b.c". Every "composite" rule of a grammar will be represented as a
// [CompositeNode] in the AST, even if it only consists of a single token.
type CompositeNode interface {
	AstNode
	StringUnit
	// IsCompositeNode marks a type as implementing [CompositeNode].
	IsCompositeNode()
}

// NewCompositeNode creates a [CompositeNode] backed by [compositeNode].
func NewCompositeNode() CompositeNode {
	return &compositeNode{}
}

// compositeNode is the default implementation of [CompositeNode].
// It is private, no adopter code should use it directly.
type compositeNode struct {
	AstNodeBase
	// We could use a sync.Once here, but that would add some overhead
	// In benchmarks, using an atomic pointer here is much faster (roughly 2x)
	cache atomic.Pointer[string]
}

func (node *compositeNode) IsCompositeNode() {}

// Owner returns the AST node that owns this string unit.
func (node *compositeNode) Owner() AstNode {
	return node.container
}

func (node *compositeNode) String() string {
	// Cache the string value, as it is accessed frequently
	// Since this operation can be done in parallel, we need an atomic pointer here
	if p := node.cache.Load(); p != nil {
		return *p
	} else {
		s := node.stringSlow()
		node.cache.Store(&s)
		return s
	}
}

func (node *compositeNode) stringSlow() string {
	// Construct the string value by concatenating the text of all tokens of the node
	// Only need to do this once, as the tokens are usually not modified after parsing
	var sb strings.Builder
	for _, token := range node.Tokens() {
		sb.WriteString(token.Image)
	}
	return sb.String()
}
