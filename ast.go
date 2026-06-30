// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"errors"
	"iter"
	"strconv"
	"strings"
	"sync/atomic"
	"unique"
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
	// ContainmentData returns a [unique.Handle] denoting the containing property within it's [AstNode.Container],
	// defaults to a a [unique.Handle] of the empty string,
	// and the element index within the containing property, defaults to zero for single item fields
	ContainmentData() (unique.Handle[string], uint16)
	// SetContainer sets the direct parent node of the node.
	//
	// When constructing an AST programmatically, use [AssignContainers] to link the node in the AST.
	SetContainer(container AstNode, containerField unique.Handle[string], index uint16)
	// Tokens returns the tokens associated with the node.
	Tokens() []*Token
	// SetToken appends token to the node's token list.
	SetToken(token *Token)
	// SetTokens replaces the node's token list with tokens.
	SetTokens(tokens []*Token)
	// Segment returns the text segment metadata of the node.
	Segment() *TextSegment
	// SetSegment sets the full text segment metadata of the node.
	//
	// It is primarily used by generated parsers while constructing nodes incrementally.
	SetSegment(segment *TextSegment)
	// SetSegmentStartToken sets the start of the node's segment from token.
	//
	// It is primarily used by generated parsers while constructing nodes incrementally.
	SetSegmentStartToken(token *Token)
	// SetSegmentEndToken sets the end of the node's segment from token.
	//
	// It is primarily used by generated parsers while constructing nodes incrementally.
	SetSegmentEndToken(token *Token)
	// Text returns the source substring covered by the node's segment.
	Text() string
	// ForEachNode calls fn for each direct child node of node.
	//
	// Note that this does not traverse the entire subtree. Use [AllNodes] or [AllChildren] for that.
	//
	// Calling this method directly is not recommended. Use [ChildNodes] instead for better readability.
	ForEachNode(fn func(AstNode, unique.Handle[string], uint16))
	// ForEachReference calls fn for each reference field of node.
	//
	// Calling this method directly is not recommended. Use [References] instead for better readability.
	ForEachReference(fn func(UntypedReference, unique.Handle[string], uint16))
	// FieldInfos returns a descriptor for the denoted field of this AstNode.
	FieldInfos(field unique.Handle[string]) FieldInfos
	// NodePath returns a slash-separated path string that uniquely identifies this node
	// within its document tree, e.g. "rules@2/alternatives@0".
	// Returns "" for the root node (no container).
	NodePath() (string, error)
	// GetByPath returns a (nested) child node denoted by the given path
	GetByPath(path string) (AstNode, error)
}

// FieldInfos is a simple struct of meta data describing a field of an AstNode.
type FieldInfos struct {
	Multi     bool
	Reference bool
}

// AstNodeBase provides the default [AstNode] implementation used by generated AST node types.
type AstNodeBase struct {
	document       *Document
	container      AstNode
	containerField unique.Handle[string]
	containerIndex uint16
	tokens         []*Token
	segment        TextSegment
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

func (node *AstNodeBase) ContainmentData() (unique.Handle[string], uint16) {
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
func (node *AstNodeBase) SetContainer(container AstNode, field unique.Handle[string], index uint16) {
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

// SetSegmentStartToken sets the start of the node's segment from token.
func (node *AstNodeBase) SetSegmentStartToken(token *Token) {
	if node != nil && token != nil {
		node.segment.Indices.Start = token.TextSegment.Indices.Start
		node.segment.Range.Start = token.TextSegment.Range.Start
	}
}

// SetSegmentEndToken sets the end of the node's segment from token.
func (node *AstNodeBase) SetSegmentEndToken(token *Token) {
	if node != nil && token != nil {
		node.segment.Indices.End = token.TextSegment.Indices.End
		node.segment.Range.End = token.TextSegment.Range.End
	}
}

// SetSegment sets the full text segment metadata of the node.
func (node *AstNodeBase) SetSegment(segment *TextSegment) {
	if node != nil {
		node.segment = *segment
	}
}

// Segment returns the text segment metadata of the node.
func (node *AstNodeBase) Segment() *TextSegment {
	if node != nil {
		return &node.segment
	} else {
		return nil
	}
}

// SetToken appends token to the node's token list.
func (node *AstNodeBase) SetToken(token *Token) {
	if node != nil && token != nil {
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

// Text returns the source substring covered by the node's segment.
func (node *AstNodeBase) Text() string {
	if node == nil || node.document == nil || node.document.TextDoc == nil {
		return ""
	} else {
		fullText := node.document.TextDoc.Text(nil)
		return fullText[node.segment.Indices.Start:node.segment.Indices.End]
	}
}

// Base Implementation returning field meta data for an AstNode.
// The generator produces specific overwrites for each generated ...Impl type.
func (node *AstNodeBase) FieldInfos(field unique.Handle[string]) FieldInfos {
	return FieldInfos{}
}

// NodePath determines node's path with it's root container in a recursive manner,
// based on node's [AstNodeBase.ContainmentData] and [AstNode.FieldInfos]
func (node *AstNodeBase) NodePath() (string, error) {
	container := node.Container()
	containerField, index := node.ContainmentData()

	if container == nil {
		return "", nil
	} else if containerField.Value() == "" {
		return "", errors.New("Can't determine node path, found AstNode's field 'containerField' being empty.")
	}

	parentPath, err := container.NodePath()
	if err != nil {
		return "", err
	}

	fieldPath := parentPath + "/" + containerField.Value()
	if container.FieldInfos(containerField).Multi {
		return fieldPath + "@" + strconv.Itoa(int(index)), nil
	} else {
		return fieldPath, nil
	}
}

// Base Implementation for instances of [AstNodeBase].
// The generator produces specific overwrites for each generated ...Impl type.
func (node *AstNodeBase) GetByPath(path string) (AstNode, error) {
	if path != "" {
		return nil, errors.New("AstNodeBase.GetByPath: Cannot identify children of plain AstNodeBase instances.")
	} else {
		return node, nil
	}
}

// ForEachNode calls fn for each direct child node of node.
//
// ForEachNode on AstNodeBase is a no-op because the base type has no child fields.
func (node *AstNodeBase) ForEachNode(fn func(AstNode, unique.Handle[string], uint16)) {
	// This base implementation does not have any contained nodes.
}

// ForEachReference calls fn for each reference field of node.
//
// ForEachReference on AstNodeBase is a no-op because the base type has no reference fields.
func (node *AstNodeBase) ForEachReference(fn func(UntypedReference, unique.Handle[string], uint16)) {
	// This base implementation does not have any references.
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
	node.ForEachNode(func(child AstNode, containerField unique.Handle[string], index uint16) {
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
		node.ForEachNode(func(child AstNode, containerField unique.Handle[string], index uint16) {
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
		node.ForEachReference(func(ref UntypedReference, containerField unique.Handle[string], index uint16) {
			if !stopped && !yield(ref) {
				stopped = true
			}
		})
	}
}

// NewAstNode creates an [AstNodeBase] with initialized token storage.
//
// It is intended for generated node implementations that embed AstNodeBase.
// AstNodeBase carries framework metadata only and has no language-specific
// semantic fields on its own.
func NewAstNode() AstNodeBase {
	return AstNodeBase{
		tokens: []*Token{},
	}
}

// AssignToken appends token to node and records node and kind on the token.
//
// It is primarily used by generated parsers while constructing nodes incrementally.
func AssignToken(node AstNode, token *Token, kind int) {
	if node != nil && token != nil {
		node.SetToken(token)
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
		// Prepend old tokens to the new node's tokens
		AssignTokens(newNode, append(oldTokens, newNode.Tokens()...))
	}
}

// AssignContainers recursively assigns document and parent pointers for root and its subtree.
//
// It also assigns document and container on composite reference units reachable via references.
func AssignContainers(doc *Document, root AstNode) {
	root.SetDocument(doc)
	root.ForEachNode(func(child AstNode, containerField unique.Handle[string], index uint16) {
		child.SetDocument(doc)
		child.SetContainer(root, containerField, index)
		AssignContainers(doc, child)
	})
	root.ForEachReference(func(ur UntypedReference, containerField unique.Handle[string], index uint16) {
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
	// Segment returns the text segment metadata of this string unit.
	Segment() *TextSegment
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

// NewCompositeNode creates a [CompositeNode] backed by [CompositeNodeBase].
func NewCompositeNode() CompositeNode {
	return &CompositeNodeBase{
		AstNodeBase: NewAstNode(),
	}
}

// CompositeNodeBase provides the default [CompositeNode] implementation for generated composite rules.
type CompositeNodeBase struct {
	AstNodeBase
	// We could use a sync.Once here, but that would add some overhead
	// In benchmarks, using an atomic pointer here is much faster (roughly 2x)
	cache atomic.Pointer[string]
}

// IsCompositeNode marks CompositeNodeBase as implementing [CompositeNode].
func (node *CompositeNodeBase) IsCompositeNode() {}

// Owner returns the AST node that owns this string unit.
func (node *CompositeNodeBase) Owner() AstNode {
	return node.container
}

// String returns the concatenated token images of node, caching the computed value.
func (node *CompositeNodeBase) String() string {
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

func (node *CompositeNodeBase) stringSlow() string {
	// Construct the string value by concatenating the text of all tokens of the node
	// Only need to do this once, as the tokens are usually not modified after parsing
	var sb strings.Builder
	for _, token := range node.Tokens() {
		sb.WriteString(token.Image)
	}
	return sb.String()
}
