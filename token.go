// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

const SkippedGroup = -1
const CommentGroup = -2

// TokenKind names the grammar construct that produced a TokenType. It is
// distinct from Group: Group controls lexer-stream behaviour (skipped /
// comment), while Kind describes the grammar origin and is consumed by
// downstream features such as the completion engine, which by default
// only surfaces keyword-kind tokens as completion candidates.
type TokenKind int

const (
	// TokenKindToken is the default - a TokenType produced by a named
	// `token` rule in the .fb grammar (regex-matched). Hidden and comment
	// tokens are also TokenKindToken; their stream behaviour is encoded
	// separately in Group.
	TokenKindToken TokenKind = 0
	// TokenKindKeyword is a TokenType produced by a literal string in a
	// parser rule (e.g. `"statemachine"`). Matched by a string prefix.
	TokenKindKeyword TokenKind = 1
)

// TokenMatcher is a function that attempts to match a token at the given offset in the input string.
// Returns the length of the match if successful, or 0 if no match is found.
type TokenMatcher func(input string, offset int) int

// TokenTypeMatcher is a function that checks if the specified other TokenType can be matched by this TokenType.
// Used for optimizations in the lookahead and other parts of the parser.
type TokenTypeMatcher func(other *TokenType) bool

type TokenType struct {
	Id         int
	Name       string
	Label      string
	StartChars []rune
	Group      int
	Kind       TokenKind
	PushMode   int
	PopMode    bool
	Match      TokenMatcher
	Matches    TokenTypeMatcher
	bitset     *BitSet
}

func NewTokenType(id int, name, label string, group int, kind TokenKind, pushMode int, popMode bool, match TokenMatcher, startChars []rune) *TokenType {
	matching := NewBitset()
	matching.Insert(id)
	return &TokenType{
		Id:    id,
		Name:  name,
		Label: label,
		Group: group,
		Kind:  kind,
		Match: match,
		Matches: func(other *TokenType) bool {
			return other.Id == id
		},
		bitset:     matching,
		PushMode:   pushMode,
		PopMode:    popMode,
		StartChars: startChars,
	}
}

func NewTokenGroup(id int, name, label string, matchingTypes []*TokenType) *TokenType {
	bitsets := make([]*BitSet, len(matchingTypes))
	for _, mt := range matchingTypes {
		bitsets = append(bitsets, mt.bitset)
	}
	matching := MergeBitSets(bitsets)
	matching.Insert(id)
	tt := &TokenType{
		Id:    id,
		Name:  name,
		Label: label,
		Matches: func(other *TokenType) bool {
			return matching.At(other.Id)
		},
		bitset: matching,
	}
	return tt
}

func (t *TokenType) IsSkipped() bool {
	return t.Group == SkippedGroup
}

func (t *TokenType) IsComment() bool {
	return t.Group == CommentGroup
}

func (t *TokenType) IsKeyword() bool {
	return t.Kind == TokenKindKeyword
}

func (t *TokenType) Bitset() *BitSet {
	return t.bitset
}

var EOF = NewTokenType(
	0,
	"EOF",
	"EOF",
	0,
	TokenKindToken,
	0,
	false,
	nil,
	nil,
)

var EOFToken = NewToken(EOF, "", 0, 0, 0, 0, 0, 0)

type Token struct {
	Type        *TokenType
	Image       string
	TypeId      int
	TextSegment TextSegment
	// Semantic information
	Element AstNode
	Kind    int
}

func NewToken(tokenType *TokenType, image string, startOffset, endOffset, startLine, endLine, startColumn, endColumn int) Token {
	return Token{
		Type:  tokenType,
		Image: image,
		TextSegment: TextSegment{
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

func (t *Token) String() string {
	return t.Image
}

func (t *Token) Segment() *TextSegment {
	return &t.TextSegment
}

func (t *Token) Owner() AstNode {
	element := t.Element
	if composite, ok := element.(CompositeNode); ok {
		// If the token is part of a composite node, its owner is the container of that
		return composite.Container()
	}
	return element
}

type TokenSlice []Token

// Searches for the token that contains the given offset.
// Expects that the tokens are sorted by their offsets to perform a binary search.
func (ts TokenSlice) SearchOffset(offset int) *Token {
	low, high := 0, len(ts)-1
	for low <= high {
		mid := (low + high) / 2
		token := ts[mid]
		if offset < int(token.TextSegment.Indices.Start) {
			high = mid - 1
		} else if offset >= int(token.TextSegment.Indices.End) {
			low = mid + 1
		} else {
			return &token
		}
	}
	return nil
}
