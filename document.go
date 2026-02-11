// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"iter"
	"sync"

	"github.com/TypeFox/go-lsp/protocol"
	"typefox.dev/fastbelt/textdoc"
)

// Document represents a document in the workspace.
// The data stored does not have to be complete during the whole lifecycle of the document.
// For example, the Root node may be nil if the document has not been parsed yet.
//
// Access to the fields of Document should be synchronized using the embedded [sync.RWMutex].
// The document struct should never be copied after creation.
type Document struct {
	sync.RWMutex
	URI          URI
	Root         AstNode
	Tokens       TokenSlice
	LocalSymbols LocalSymbols
	ParserErrors []*ParserError
	LexerErrors  []*LexerError
	References   []UntypedReference
	TextDoc      textdoc.Handle
	Diagnostics  []*protocol.Diagnostic
	Data         map[any]any
}

func NewDocument(textDoc textdoc.Handle) (*Document, error) {
	uri, err := ParseURI(string(textDoc.URI()))
	if err != nil {
		return nil, err
	}
	return &Document{
		RWMutex:      sync.RWMutex{},
		URI:          uri,
		TextDoc:      textDoc,
		Root:         nil,
		LocalSymbols: nil,
		Data:         map[any]any{},
		Tokens:       TokenSlice{},
		ParserErrors: []*ParserError{},
		LexerErrors:  []*LexerError{},
		References:   []UntypedReference{},
		Diagnostics:  []*protocol.Diagnostic{},
	}, nil
}

func NewDocumentFromString(uri, languageId, content string) (*Document, error) {
	textDoc, err := textdoc.NewFile(protocol.DocumentURI(uri), languageId, 1, content)
	if err != nil {
		return nil, err
	}
	doc, err := NewDocument(textDoc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

type SymbolList = iter.Seq[*AstNodeDescription]

type LocalSymbols interface {
	Has(node AstNode) bool
	Iter(node AstNode) SymbolList
}
