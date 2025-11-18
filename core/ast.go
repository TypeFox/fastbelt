// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package core

type Reference[T AstNode] struct {
	Token *Token
	Text  string
	ref   *T
}

func (r *Reference[T]) Get() *T {
	return r.ref
}

// TODO: implement this properly. This probably should point to `textdoc.Overlay`
type Document struct {
	Text string
}

type AstNodeBase struct {
	document  *Document
	container AstNode
	tokens    []*Token
	segment   TextSegment
}

func (node *AstNodeBase) ForEachNode(fn func(AstNode)) {}

func (node *AstNodeBase) Document() *Document {
	if node != nil {
		return node.document
	} else {
		return nil
	}
}

func (node *AstNodeBase) WithDocument(document *Document) {
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

func (node *AstNodeBase) WithContainer(container AstNode) {
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

func (node *AstNodeBase) WithSegmentStartToken(token *Token) {
	if node != nil && token != nil {
		node.segment.Indices.Start = token.Segment.Indices.Start
		node.segment.Range.Start = token.Segment.Range.Start
	}
}

func (node *AstNodeBase) WithSegmentEndToken(token *Token) {
	if node != nil && token != nil {
		node.segment.Indices.End = token.Segment.Indices.End
		node.segment.Range.End = token.Segment.Range.End
	}
}

func (node *AstNodeBase) Segment() *TextSegment {
	if node != nil {
		return &node.segment
	} else {
		return nil
	}
}

func (node *AstNodeBase) WithToken(token *Token) {
	if node != nil && token != nil {
		node.tokens = append(node.tokens, token)
	}
}

func (node *AstNodeBase) WithTokens(tokens []*Token) {
	if node != nil {
		// The method is called to set all tokens of the node at once
		// The old node is discarded in the process
		// Therefore, we don't append but replace the token slice
		node.tokens = tokens
	}
}

func (node *AstNodeBase) Text() string {
	if node == nil || node.document == nil {
		return ""
	} else {
		return node.document.Text[node.segment.Indices.Start:node.segment.Indices.End]
	}
}

type AstNode interface {
	// Getters and setters
	Document() *Document
	WithDocument(document *Document)
	Container() AstNode
	WithContainer(container AstNode)
	WithToken(token *Token)
	WithTokens(tokens []*Token)
	WithSegmentStartToken(token *Token)
	WithSegmentEndToken(token *Token)
	// Getter only
	Text() string
	Tokens() []*Token
	Segment() *TextSegment
	ForEachNode(func(AstNode))
}

func NewAstNode() AstNodeBase {
	return AstNodeBase{
		tokens: []*Token{},
	}
}

func AssignToken(node AstNode, token *Token, kind int) {
	if node != nil && token != nil {
		node.WithToken(token)
		token.Element = node
		token.Kind = kind
	}
}

func AssignTokens(node AstNode, tokens []*Token) {
	if node != nil {
		node.WithTokens(tokens)
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

func AssignContainers(root AstNode) {
	root.ForEachNode(func(child AstNode) {
		child.WithContainer(root)
		AssignContainers(child)
	})
}
