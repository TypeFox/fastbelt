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

// LookaheadStrategy abstracts OR-decision prediction for generated parsers.
// key is the ATN decision key ("RuleName_ProdType_N", 1-based).
// Predict returns the chosen alternative index (0-based), or -1.
// PredictOpt returns true when the optional / loop body should be entered.
type LookaheadStrategy interface {
	Predict(src *ParserState, key string) int
	PredictOpt(src *ParserState, key string) bool
}

// LLkStrategy implements LookaheadStrategy using pre-built LL(k) tables.
// tables maps ATN decision key → LLkLookahead table.
// It is the default strategy installed by generated NewParser constructors.
type LLkStrategy struct {
	tables map[string]LLkLookahead
}

// NewLLkStrategy creates an LLkStrategy from a key→table map.
func NewLLkStrategy(tables map[string]LLkLookahead) *LLkStrategy {
	return &LLkStrategy{tables: tables}
}

func (s *LLkStrategy) Predict(src *ParserState, key string) int {
	t, ok := s.tables[key]
	if !ok {
		return -1
	}
	for i, option := range t {
	outer:
		for _, path := range option {
			for j, tokenType := range path {
				if src.LAId(j+1) != tokenType {
					continue outer
				}
			}
			return i
		}
	}
	return -1
}

func (s *LLkStrategy) PredictOpt(src *ParserState, key string) bool {
	return s.Predict(src, key) == 0
}

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


