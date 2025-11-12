// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lexer

import (
	"github.com/TypeFox/langium-to-go/core"
)

const SkippedGroup = -1

type Matcher func(input string, offset int) int

type TokenType struct {
	Id         int
	Name       string
	Label      string
	StartChars []rune
	Group      int
	PushMode   int
	PopMode    bool
	Match      Matcher
}

func NewTokenType(id int, name, label string, group int, pushMode int, popMode bool, match Matcher, startChars []rune) *TokenType {
	return &TokenType{
		Id:         id,
		Name:       name,
		Label:      label,
		Group:      group,
		Match:      match,
		PushMode:   pushMode,
		PopMode:    popMode,
		StartChars: startChars,
	}
}

func (t *TokenType) IsSkipped() bool {
	return t.Group == SkippedGroup
}

var EOF = NewTokenType(
	0,
	"EOF",
	"EOF",
	0,
	0,
	false,
	nil,
	[]rune{},
)

var EOFToken = NewToken(EOF, "", 0, 0, 0, 0, 0, 0)

type Token struct {
	Type    *TokenType
	Image   string
	TypeId  int
	Segment core.TextSegment
	// Semantic information
	Element any
	Kind    int
}

func NewToken(tokenType *TokenType, image string, startOffset, endOffset, startLine, endLine, startColumn, endColumn int) *Token {
	return &Token{
		Type:  tokenType,
		Image: image,
		Segment: core.TextSegment{
			Indices: core.TextIndexRange{
				Start: core.TextIndex(startOffset),
				End:   core.TextIndex(endOffset),
			},
			Range: core.TextRange{
				Start: core.TextLocation{
					Line:   core.TextLine(startLine),
					Column: core.TextColumn(startColumn),
				},
				End: core.TextLocation{
					Line:   core.TextLine(endLine),
					Column: core.TextColumn(endColumn),
				},
			},
		},
		TypeId: tokenType.Id,
		Kind:   0,
	}
}

func (t *Token) IsSkipped() bool {
	return t.Type != nil && t.Type.Group == SkippedGroup
}

func (t *Token) IsEOF() bool {
	return t.Type == EOF
}

func (t *Token) Assign(element any, kind int) {
	t.Element = element
	t.Kind = kind
}
