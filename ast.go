// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import "iter"

type AstNodeBase struct {
	document  *Document
	container AstNode
	tokens    []*Token
	segment   TextSegment
}

func (node *AstNodeBase) Document() *Document {
	if node != nil {
		return node.document
	} else {
		return nil
	}
}

func (node *AstNodeBase) SetDocument(document *Document) {
	if node != nil {
		node.document = document
	}
}

func (node *AstNodeBase) Container() AstNode {
	if node != nil {
		return node.container
	} else {
		return nil
	}
}

// TODO: If concrete methods gain access to generics, refactor this into a method
// See https://github.com/golang/go/issues/77273
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

func (node *AstNodeBase) SetContainer(container AstNode) {
	if node != nil {
		node.container = container
	}
}

func (node *AstNodeBase) Tokens() []*Token {
	if node != nil {
		return node.tokens
	} else {
		return nil
	}
}

func (node *AstNodeBase) SetSegmentStartToken(token *Token) {
	if node != nil && token != nil {
		node.segment.Indices.Start = token.Segment.Indices.Start
		node.segment.Range.Start = token.Segment.Range.Start
	}
}

func (node *AstNodeBase) SetSegmentEndToken(token *Token) {
	if node != nil && token != nil {
		node.segment.Indices.End = token.Segment.Indices.End
		node.segment.Range.End = token.Segment.Range.End
	}
}

func (node *AstNodeBase) SetSegment(segment *TextSegment) {
	if node != nil {
		node.segment = *segment
	}
}

func (node *AstNodeBase) Segment() *TextSegment {
	if node != nil {
		return &node.segment
	} else {
		return nil
	}
}

func (node *AstNodeBase) SetToken(token *Token) {
	if node != nil && token != nil {
		node.tokens = append(node.tokens, token)
	}
}

func (node *AstNodeBase) SetTokens(tokens []*Token) {
	if node != nil {
		// The method is called to set all tokens of the node at once
		// The old node is discarded in the process
		// Therefore, we don't append but replace the token slice
		node.tokens = tokens
	}
}

func (node *AstNodeBase) Text() string {
	if node == nil || node.document == nil || node.document.TextDoc == nil {
		return ""
	} else {
		return node.document.TextDoc.Text(nil)[node.segment.Indices.Start:node.segment.Indices.End]
	}
}

func (node *AstNodeBase) ForEachNode(fn func(AstNode)) {
	// This base implementation does not have any contained nodes.
}

func (node *AstNodeBase) ForEachReference(fn func(UntypedReference)) {
	// This base implementation does not have any references.
}

// AstNode is the base interface for all AST nodes.
type AstNode interface {
	Document() *Document
	SetDocument(document *Document)
	Container() AstNode
	SetContainer(container AstNode)
	Tokens() []*Token
	SetToken(token *Token)
	SetTokens(tokens []*Token)
	Segment() *TextSegment
	SetSegment(segment *TextSegment)
	// Sets the start of the node's segment to the start of the given token's segment.
	// Should only be called by the parser. Use SetSegment to set both start and end manually.
	SetSegmentStartToken(token *Token)
	// Sets the end of the node's segment to the end of the given token's segment.
	// Should only be called by the parser. Use SetSegment to set both start and end manually.
	SetSegmentEndToken(token *Token)
	Text() string
	// ForEachNode calls the given function for each direct child node of this node.
	// Note that this does not traverse the entire subtree. Use [AllNodes] or [TraverseNode] for that.
	//
	// Calling this method directly is not recommended. Use [ChildNodes] instead for better readability.
	ForEachNode(fn func(AstNode))
	// ForEachReference calls the given function for each reference contained in this node.
	//
	// Calling this method directly is not recommended. Use [References] instead for better readability.
	ForEachReference(fn func(UntypedReference))
}

// Performance note about traversal function:
// Theoretically, we could have ChildNodes and References directly as methods on the AstNode interface.
// However, implementing the deep traversal on top of an iter.Seq is very inefficient.
// In benchmarks, it is roughly 5x slower than the current implementation.
// By using a callback-based approach, we can traverse the entire subtree with minimal overhead.
// But we lose the ability to short-circuit the traversal when we find what we're looking for.
// In practice, this is not a big issue, because most traversals will need to visit most of the nodes anyway.
// AllNodes and AllChildren are slightly less efficient than TraverseNode and TraverseContent,
// but only by roughly 10%, and they provide a much nicer API for most use cases, so the trade-off is worth it.

// Traverses the given node and all its children, calling the given function for each node.
//
// Calling this function directly is not recommended. Use [AllNodes] instead for better readability.
// Note that both [TraverseNode] and [AllNodes] will traverse the entire subtree, without short-circuiting.
func TraverseNode(node AstNode, fn func(AstNode)) {
	fn(node)
	TraverseContent(node, fn)
}

// Traverses all children of the given node, calling the specified function for each child.
// Does not call the function for the given node itself. Use [TraverseNode] for that.
//
// Calling this function directly is not recommended. Use [AllChildren] instead for better readability.
// Note that both [TraverseContent] and [AllChildren] will traverse the entire subtree, without short-circuiting.
func TraverseContent(node AstNode, fn func(AstNode)) {
	node.ForEachNode(func(child AstNode) {
		fn(child)
		TraverseContent(child, fn)
	})
}

// [AllNodes] creates an iterator over the given node and all its descendant nodes.
//
// This function wraps [TraverseNode] in an [iter.Seq].
// Early loop exit is honoured correctly, but does not short-circuit the traversal.
func AllNodes(node AstNode) iter.Seq[AstNode] {
	return func(yield func(AstNode) bool) {
		stopped := false
		TraverseNode(node, func(n AstNode) {
			if !stopped && !yield(n) {
				stopped = true
			}
		})
	}
}

// [AllChildren] creates an iterator over all descendant nodes of the given node, excluding the node itself.
//
// This function wraps [TraverseContent] in an [iter.Seq].
// Early loop exit is honoured correctly, but does not short-circuit the traversal.
func AllChildren(node AstNode) iter.Seq[AstNode] {
	return func(yield func(AstNode) bool) {
		stopped := false
		TraverseContent(node, func(n AstNode) {
			if !stopped && !yield(n) {
				stopped = true
			}
		})
	}
}

// [ChildNodes] creates an iterator over the direct child nodes of the given node.
//
// This function wraps [AstNode.ForEachNode] in an [iter.Seq].
// Early loop exit is honoured correctly, but does not short-circuit the traversal.
func ChildNodes(node AstNode) iter.Seq[AstNode] {
	return func(yield func(AstNode) bool) {
		stopped := false
		node.ForEachNode(func(child AstNode) {
			if !stopped && !yield(child) {
				stopped = true
			}
		})
	}
}

// [References] creates an iterator over all references of the given node.
//
// This function wraps [AstNode.ForEachReference] in an [iter.Seq].
// Early loop exit is honoured correctly, but does not short-circuit the traversal.
func References(node AstNode) iter.Seq[UntypedReference] {
	return func(yield func(UntypedReference) bool) {
		stopped := false
		node.ForEachReference(func(ref UntypedReference) {
			if !stopped && !yield(ref) {
				stopped = true
			}
		})
	}
}

func NewAstNode() AstNodeBase {
	return AstNodeBase{
		tokens: []*Token{},
	}
}

func AssignToken(node AstNode, token *Token, kind int) {
	if node != nil && token != nil {
		node.SetToken(token)
		token.Element = node
		token.Kind = kind
	}
}

func AssignTokens(node AstNode, tokens []*Token) {
	if node != nil {
		node.SetTokens(tokens)
		for _, token := range tokens {
			token.Element = node
		}
	}
}

func MergeTokens(newNode AstNode, oldTokens []*Token) {
	if newNode != nil && len(oldTokens) > 0 {
		// Prepend old tokens to the new node's tokens
		AssignTokens(newNode, append(oldTokens, newNode.Tokens()...))
	}
}

func AssignContainers(doc *Document, root AstNode) {
	root.SetDocument(doc)
	root.ForEachNode(func(child AstNode) {
		child.SetDocument(doc)
		child.SetContainer(root)
		AssignContainers(doc, child)
	})
}

type NamedNode interface {
	AstNode
	Name() string
	NameToken() *Token
}
