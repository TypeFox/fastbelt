// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// TODO Move this stuff to the core package?
package parser

import (
	core "typefox.dev/fastbelt"
)

// Parser defines the interface for parsing tokens (lexer output) into AST nodes.
type Parser interface {
	Parse(document *core.Document) *ParseResult
}

type ParseResult struct {
	Node   core.AstNode
	Errors []*core.ParserError
}

type ParserState struct {
	Tokens []*core.Token
	Length int
	Index  int
	State  int

	next *core.Token
	// Indicates whether a parsing error has occurred
	// and the parser is in error recovery mode.
	// The parser will not consume any more tokens while in error mode.
	//
	// TODO: implement proper error handling
	inError bool
	errors  []*core.ParserError
}

func (p *ParserState) Errors() []*core.ParserError {
	return p.errors
}

func (p *ParserState) appendError(msg string, token *core.Token) {
	p.errors = append(p.errors, core.NewParserError(msg, token))
	p.inError = true
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
		Tokens:  tokens,
		Length:  len(tokens),
		Index:   0,
		State:   0,
		next:    next,
		inError: false,
		errors:  []*core.ParserError{},
	}
}

func (p *ParserState) LA(offset int) *core.Token {
	pos := p.Index + offset - 1
	if pos < 0 || pos >= p.Length || p.inError {
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
	if p.inError {
		return nil
	}
	current := p.next
	if current == nil {
		// EOF reached
		p.appendError("Unexpected end of input.", nil)
		return nil
	}
	if current.TypeId != tokenType {
		// Generate error
		p.appendError("Unexpected token '"+current.Image+"'.", current)
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
