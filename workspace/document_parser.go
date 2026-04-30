// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/service"
)

// DocumentParser is a service for parsing documents.
type DocumentParser interface {
	// Parses the document and stores the resulting Tokens and AST node Root (incl. potential errors) into the document.
	// The caller must hold the document's write lock.
	Parse(doc *core.Document)
}

// DefaultDocumentParser is the default implementation of [DocumentParser].
type DefaultDocumentParser struct {
	sc *service.Container
}

func NewDefaultDocumentParser(sc *service.Container) DocumentParser {
	return &DefaultDocumentParser{sc: sc}
}

func (s *DefaultDocumentParser) Parse(doc *core.Document) {
	text := doc.TextDoc.Text(nil)
	// Run the lexer
	lexer := service.MustGet[lexer.Lexer](s.sc)
	lexerRes := lexer.Lex(text)
	doc.LexerErrors = lexerRes.Errors
	doc.Tokens = lexerRes.Tokens
	doc.Comments = lexerRes.Comments
	// Run the parser
	parser := service.MustGet[parser.Parser](s.sc)
	parserRes := parser.Parse(doc)
	doc.ParserErrors = parserRes.Errors
	doc.Root = parserRes.Node
}
