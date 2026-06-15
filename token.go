// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import "typefox.dev/fastbelt/util/collections"

// SkippedGroup marks token types that the lexer should drop from all output streams.
const SkippedGroup = -1

// CommentGroup marks token types that the lexer should collect in Document.Comments.
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
	// TokenKindGroup is a TokenType that represents a named group of other
	// TokenTypes. It is not directly matched against the input, but instead
	// serves as a convenient (and fast!) way to refer to multiple TokenTypes
	// in the grammar and downstream features.
	TokenKindGroup TokenKind = 2
)

// TokenMatcher is a function that attempts to match a token at the given offset in the input string.
// Returns the byte-length of the match if successful, or 0 if no match is found.
type TokenMatcher func(input string, offset int) int

// TokenTypeMatcher is a function that checks if the specified other TokenType can be matched by this TokenType.
// Used for optimizations in the lookahead and other parts of the parser.
type TokenTypeMatcher func(other *TokenType) bool

// TokenType describes one lexer token category generated from a grammar.
type TokenType struct {
	// Id is the stable numeric token identifier used in parser tables and token bitsets.
	Id int
	// Name is the token name from generated code.
	Name string
	// Label is the user-facing token label, for example in completion output.
	Label string
	// StartChars contains candidate start runes used for lexer preselection.
	StartChars []rune
	// Group controls lexer output routing (default stream, skipped, comments, or custom groups).
	Group int
	// Kind records whether the token comes from a keyword literal or a token rule.
	Kind TokenKind
	// PushMode selects the next lexer mode after this token is matched.
	PushMode int
	// PopMode reports whether matching this token pops one lexer mode.
	PopMode bool
	// Match performs the actual token match at a given input offset.
	Match TokenMatcher
	// Matches returns whether the token type matches another, given type
	Matches TokenTypeMatcher
	// All TokenTypes that are matched by this TokenType.
	MatchingTokens []*TokenType
	bitset         *collections.BitSet
}

// NewTokenType creates a token type descriptor used by generated lexers and parsers.
func NewTokenType(id int, name, label string, group int, kind TokenKind, pushMode int, popMode bool, match TokenMatcher, startChars []rune) *TokenType {
	matching := collections.NewBitset()
	matching.Insert(id)
	tt := &TokenType{
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
	tt.MatchingTokens = []*TokenType{tt}
	return tt
}

// NewTokenGroup creates a token type that represents a named group of other token types.
func NewTokenGroup(id int, name, label string, matchingTypes []*TokenType) *TokenType {
	bitsets := make([]*collections.BitSet, len(matchingTypes))
	for _, mt := range matchingTypes {
		bitsets = append(bitsets, mt.bitset)
	}
	matching := collections.MergeBitSets(bitsets)
	matching.Insert(id)
	tt := &TokenType{
		Id:    id,
		Name:  name,
		Label: label,
		Kind:  TokenKindGroup,
		Matches: func(other *TokenType) bool {
			return matching.At(other.Id)
		},
		bitset: matching,
	}
	tt.MatchingTokens = unrollMatchingTokens(matchingTypes)
	return tt
}

func unrollMatchingTokens(matchingTypes []*TokenType) []*TokenType {
	result := map[int]*TokenType{}
	for _, mt := range matchingTypes {
		if mt.Kind == TokenKindGroup {
			for _, t := range unrollMatchingTokens(mt.MatchingTokens) {
				result[t.Id] = t
			}
		} else {
			result[mt.Id] = mt
		}
	}
	unrolled := make([]*TokenType, 0, len(result))
	for _, t := range result {
		unrolled = append(unrolled, t)
	}
	return unrolled
}

// IsSkipped reports whether t is routed to the skipped-token group.
func (t *TokenType) IsSkipped() bool {
	return t.Group == SkippedGroup
}

// IsComment reports whether t is routed to the comment-token group.
func (t *TokenType) IsComment() bool {
	return t.Group == CommentGroup
}

// IsKeyword reports whether t originates from a grammar keyword literal.
func (t *TokenType) IsKeyword() bool {
	return t.Kind == TokenKindKeyword
}

// Bitset returns a bitset that contains all token type IDs matched by t.
func (t *TokenType) Bitset() *collections.BitSet {
	return t.bitset
}

// EOF is the sentinel token type used to represent end of input.
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

// EOFToken is the reusable sentinel token value for end of input.
var EOFToken = NewToken(EOF, "", 0, 0, 0, 0, 0, 0)

// Token represents one lexed source slice with positional metadata.
type Token struct {
	// Type points to the matched token type metadata.
	Type *TokenType
	// Image is the exact text matched for this token.
	Image string
	// TypeId caches Type.Id for parser hot paths.
	TypeId int
	// TextSegment stores byte offsets and line/column ranges for this token.
	TextSegment TextSegment
	// Element points to the AST node this token was assigned to during parsing.
	Element AstNode
	// Kind stores the generated assignment slot identifier within Element.
	Kind int
}

// NewToken creates a token with image text and half-open source coordinates.
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

// IsSkipped reports whether t belongs to the skipped lexer group.
func (t *Token) IsSkipped() bool {
	return t.Type != nil && t.Type.Group == SkippedGroup
}

// IsEOF reports whether t uses the end-of-input token type.
func (t *Token) IsEOF() bool {
	return t.Type == EOF
}

// Assign records the owning AST element and assignment slot kind for t.
func (t *Token) Assign(element AstNode, kind int) {
	t.Element = element
	t.Kind = kind
}

// String returns the token image.
func (t *Token) String() string {
	return t.Image
}

// Segment returns the token text segment.
func (t *Token) Segment() *TextSegment {
	return &t.TextSegment
}

// Owner returns the AST node that owns t as a string unit.
//
// For tokens attached to a CompositeNode, Owner returns the composite's container.
func (t *Token) Owner() AstNode {
	element := t.Element
	if composite, ok := element.(CompositeNode); ok {
		// If the token is part of a composite node, its owner is the container of that
		return composite.Container()
	}
	return element
}

// TokenSlice is a token sequence sorted by source offsets.
type TokenSlice []Token

// SearchOffset returns the token that contains the given offset. If the
// offset is exactly between two tokens, it returns the token after the
// offset (i.e. the one that starts at that offset). If no token contains
// the offset, it returns nil.
//
// It expects the tokens to be sorted by token offsets and uses binary search.
func (ts TokenSlice) SearchOffset(offset int) *Token {
	prev, next := ts.SearchOffset2(offset)
	if next != nil {
		return next
	}
	return prev
}

// SearchOffset2 returns the tokens at the given offset. If the offset is
// inside of a token, the first return value is that token and the second
// is nil. If the offset is exactly between two tokens, the first return
// value is the previous token and the second is the next token.
//
// It expects the tokens to be sorted by token offsets and uses binary search.
func (ts TokenSlice) SearchOffset2(offset int) (*Token, *Token) {
	low, high := 0, len(ts)-1
	for low <= high {
		mid := (low + high) / 2
		token := &ts[mid]
		if offset < int(token.TextSegment.Indices.Start) {
			high = mid - 1
		} else if offset >= int(token.TextSegment.Indices.End) {
			low = mid + 1
		} else {
			// Offset sits exactly on the boundary between this token and the
			// previous one (prev.End == offset == token.Start): return both.
			if offset == int(token.TextSegment.Indices.Start) && mid > 0 {
				if prev := &ts[mid-1]; int(prev.TextSegment.Indices.End) == offset {
					return prev, token
				}
			}
			return token, nil
		}
	}
	return nil, nil
}
