// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	core "typefox.dev/fastbelt"
)

// Parser defines the interface for parsing tokens (lexer output) into AST nodes.
type Parser interface {
	Parse(tokens []*core.Token) core.AstNode
}

type ParserState struct {
	Tokens []*core.Token
	Length int
	Index  int
	State  int

	next *core.Token
	// Indicates whether a parsing error has occurred
	// to prevent further processing
	//
	// TODO: implement proper error handling
	err bool
}

type LookaheadPath []int
type LookaheadOption []LookaheadPath
type LLkLookahead []LookaheadOption

func NewParserState(tokens []*core.Token) *ParserState {
	var next *core.Token
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

func (p *ParserState) LA(offset int) *core.Token {
	pos := p.Index + offset - 1
	if pos < 0 || pos >= p.Length || p.err {
		return nil
	}
	return p.Tokens[pos]
}

func (p *ParserState) LAId(offset int) int {
	la := p.LA(offset)
	if la == nil {
		return core.EOF.Id
	}
	return la.TypeId
}

func (p *ParserState) Consume(tokenType int) *core.Token {
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
	} else {
		p.next = nil
	}
	return current
}

func (p *ParserState) Lookahead(value LLkLookahead) int {
	for i, option := range value {
	outer:
		for _, path := range option {
			for j, tokenType := range path {
				if p.LAId(j+1) != tokenType {
					continue outer
				}
			}
			return i
		}
	}
	return -1
}
