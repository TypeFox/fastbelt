// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
)

// DocumentParser defines the interface for parsing a document.
type DocumentParser interface {
	// Parses the document and stores the resulting data back into the document.
	Parse(doc *core.Document)
}

// DefaultDocumentParser is the default implementation of DocumentParser.
type DefaultDocumentParser struct {
	srv WorkspaceSrvCont
}

// NewDefaultDocumentParser creates a new default document parser.
func NewDefaultDocumentParser(srv WorkspaceSrvCont) DocumentParser {
	return &DefaultDocumentParser{srv: srv}
}

func (p *DefaultDocumentParser) Parse(doc *core.Document) {
	text := doc.TextDoc.Text(nil)
	lexerRes := p.srv.Generated().Lexer.Lex(text)
	doc.LexerErrors = lexerRes.Errors
	doc.Tokens = lexerRes.Tokens
	parserRes := p.srv.Generated().Parser.Parse(doc)
	doc.ParserErrors = parserRes.Errors
	doc.Root = parserRes.Node
}
