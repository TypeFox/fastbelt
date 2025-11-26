// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/textdoc"
)

// ParseResult contains the result of parsing a document, including errors from the lexer and parser.
type ParseResult struct {
	Root         core.AstNode
	LexerErrors  []*lexer.LexerError
	ParserErrors string
}

// DocumentParser defines the interface for parsing a document into an AST node.
type DocumentParser interface {
	Parse(doc textdoc.Handle) ParseResult
}

// DefaultDocumentParser is the default implementation of DocumentParser.
type DefaultDocumentParser struct {
	srv WorkspaceSrvCont
}

// NewDefaultDocumentParser creates a new default document parser.
func NewDefaultDocumentParser(srv WorkspaceSrvCont) DocumentParser {
	return &DefaultDocumentParser{srv: srv}
}

func (p *DefaultDocumentParser) Parse(doc textdoc.Handle) ParseResult {
	result := ParseResult{}
	text := doc.Text(nil)
	lexerRes := p.srv.Generated().Lexer.Lex(text)
	result.LexerErrors = lexerRes.Errors
	parserRes := p.srv.Generated().Parser.Parse(lexerRes.Tokens)
	if parserRes == nil {
		result.ParserErrors = "parser error"
	} else {
		result.Root = parserRes
	}
	return result
}
