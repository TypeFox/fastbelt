// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

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
)

// Matcher reports how many bytes a token type matches at offset in input.
//
// It returns 0 when there is no match.
type Matcher func(input string, offset int) int

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
	Match Matcher
}

// NewTokenType creates a token type descriptor used by generated lexers and parsers.
func NewTokenType(id int, name, label string, group int, kind TokenKind, pushMode int, popMode bool, match Matcher, startChars []rune) *TokenType {
	return &TokenType{
		Id:         id,
		Name:       name,
		Label:      label,
		Group:      group,
		Kind:       kind,
		Match:      match,
		PushMode:   pushMode,
		PopMode:    popMode,
		StartChars: startChars,
	}
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

// EOF is the sentinel token type used to represent end of input.
var EOF = NewTokenType(
	-1,
	"EOF",
	"EOF",
	0,
	TokenKindToken,
	0,
	false,
	nil,
	[]rune{},
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

// SearchOffset returns the token that contains offset.
//
// It expects ts to be sorted by token offsets and uses binary search.
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
