// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"github.com/TypeFox/langium-to-go/lexer"
)

type ParserState struct {
	Tokens []*lexer.Token
	Length int
	Index  int
	State  int

	next *lexer.Token
	// Indicates whether a parsing error has occurred
	// to prevent further processing
	//
	// TODO: implement proper error handling
	err bool
}

type LookaheadPath []int
type LookaheadOption []LookaheadPath
type LLkLookahead []LookaheadOption

func NewParserState(tokens []*lexer.Token) *ParserState {
	var next *lexer.Token
	if len(tokens) > 0 {
		next = tokens[0]
	}
	return &ParserState{
		Tokens: tokens,
		Length: len(tokens),
		Index:  0,
		State:  0,
		next:   next,
		err:    false,
	}
}

func (p *ParserState) LA(offset int) int {
	pos := p.Index + offset - 1
	if pos >= p.Length || p.err {
		return -1
	}
	return p.Tokens[pos].TypeId
}
func (p *ParserState) Consume(tokenType int) *lexer.Token {
	current := p.next
	if current == nil || p.err {
		// EOF reached
		return nil
	}
	if current.TypeId != tokenType {
		// Generate error
		p.err = true
		return nil
	}
	p.Index++
	if p.Index < p.Length {
		p.next = p.Tokens[p.Index]
	}
	return current
}
func (p *ParserState) Lookahead(value LLkLookahead) int {
	for i, option := range value {
	outer:
		for _, path := range option {
			for j, tokenType := range path {
				if p.LA(j+1) != tokenType {
					continue outer
				}
			}
			return i
		}
	}
	return -1
}
