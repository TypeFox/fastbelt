// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

const SkippedGroup = -1
const CommentGroup = -2

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

func (t *TokenType) IsComment() bool {
	return t.Group == CommentGroup
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
	Segment TextSegment
	// Semantic information
	Element AstNode
	Kind    int
}

func NewToken(tokenType *TokenType, image string, startOffset, endOffset, startLine, endLine, startColumn, endColumn int) *Token {
	return &Token{
		Type:  tokenType,
		Image: image,
		Segment: TextSegment{
			Indices: TextIndexRange{
				Start: TextIndex(startOffset),
				End:   TextIndex(endOffset),
			},
			Range: TextRange{
				Start: TextLocation{
					Line:   TextLine(startLine),
					Column: TextColumn(startColumn),
				},
				End: TextLocation{
					Line:   TextLine(endLine),
					Column: TextColumn(endColumn),
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

func (t *Token) Assign(element AstNode, kind int) {
	t.Element = element
	t.Kind = kind
}

type TokenSlice []*Token

// Searches for the token that contains the given offset.
// Expects that the tokens are sorted by their offsets to perform a binary search.
func (ts TokenSlice) SearchOffset(offset int) *Token {
	low, high := 0, len(ts)-1
	for low <= high {
		mid := (low + high) / 2
		token := ts[mid]
		if offset < int(token.Segment.Indices.Start) {
			high = mid - 1
		} else if offset >= int(token.Segment.Indices.End) {
			low = mid + 1
		} else {
			return token
		}
	}
	return nil
}
