// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package ast

import (
	"github.com/TypeFox/langium-to-go/core"
	"github.com/TypeFox/langium-to-go/lexer"
)

type Reference[T AstNode] struct {
	Token *lexer.Token
	Text  string
	ref   *T
}

func (r *Reference[T]) Get() *T {
	return r.ref
}

type AstNodeBase struct {
	document  *core.Document
	container AstNode
	tokens    []*lexer.Token
	segment   *core.TextSegment
}

func (node *AstNodeBase) Document() *core.Document {
	if node != nil {
		return node.document
	} else {
		return nil
	}
}

func (node *AstNodeBase) WithDocument(document *core.Document) {
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

func (node *AstNodeBase) Tokens() []*lexer.Token {
	if node != nil {
		return node.tokens
	} else {
		return nil
	}
}

func (node *AstNodeBase) Segment() *core.TextSegment {
	if node != nil {
		return node.segment
	} else {
		return nil
	}
}

func (node *AstNodeBase) WithToken(token *lexer.Token, kind int) {
	if node != nil && token != nil {
		node.tokens = append(node.tokens, token)
		token.Element = node
		token.Kind = kind
	}
}

func (node *AstNodeBase) WithTokens(tokens []*lexer.Token) {
	if node != nil {
		node.tokens = append(node.tokens, tokens...)
		for _, token := range tokens {
			token.Element = node
		}
	}
}

func (node *AstNodeBase) Text() string {
	return node.document.Text[node.segment.Indices.Start:node.segment.Indices.End]
}

type AstNodeCallback func(AstNode)

type AstNode interface {
	// Getters and setters
	Document() *core.Document
	WithDocument(document *core.Document)
	Container() AstNode
	WithContainer(container AstNode)
	WithToken(token *lexer.Token, kind int)
	WithTokens(tokens []*lexer.Token)
	// Getter only
	Text() string
	Tokens() []*lexer.Token
	Segment() *core.TextSegment
	ForEachNode(func(AstNode))
}

func NewAstNode() AstNodeBase {
	return AstNodeBase{
		tokens: []*lexer.Token{},
	}
}
