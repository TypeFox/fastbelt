// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"strings"
	"sync"

	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/lsp"
)

// Document represents a document in the workspace.
// The data stored does not have to be complete during the whole lifecycle of the document.
// For example, the Root node may be nil if the document has not been parsed yet.
//
// Access to the fields of Document should be synchronized using the embedded [sync.RWMutex].
// The document struct should never be copied after creation.
type Document struct {
	sync.RWMutex
	URI             URI
	State           DocumentState
	Root            AstNode
	Tokens          TokenSlice
	LocalSymbols    LocalSymbols
	ExportedSymbols []*AstNodeDescription
	ImportedSymbols SymbolList
	ParserErrors    []*ParserError
	LexerErrors     []*LexerError
	References      []UntypedReference
	TextDoc         textdoc.Handle
	Diagnostics     []*lsp.Diagnostic
	Data            map[any]any
}

func NewDocument(textDoc textdoc.Handle) *Document {
	uri := ParseURI(string(textDoc.URI()))
	return &Document{
		RWMutex:         sync.RWMutex{},
		URI:             uri,
		State:           0,
		TextDoc:         textDoc,
		Root:            nil,
		LocalSymbols:    nil,
		ExportedSymbols: nil,
		Data:            map[any]any{},
		Tokens:          TokenSlice{},
		ParserErrors:    []*ParserError{},
		LexerErrors:     []*LexerError{},
		References:      []UntypedReference{},
		Diagnostics:     []*lsp.Diagnostic{},
	}
}

func NewDocumentFromString(uri, languageId, content string) (*Document, error) {
	textDoc, err := textdoc.NewFile(lsp.DocumentURI(uri), languageId, 1, content)
	if err != nil {
		return nil, err
	}
	doc := NewDocument(textDoc)
	return doc, nil
}

// DocumentState is a bitmask capturing the already completed build phases of a document.
type DocumentState uint32

const (
	DocStateParsed          DocumentState = 1 << iota // 0x0001
	DocStateExportedSymbols                           // 0x0002
	DocStateImportedSymbols                           // 0x0004
	DocStateLocalSymbols                              // 0x0008
	DocStateLinked                                    // 0x0010
	DocStateValidated                                 // 0x0020
)

func (s DocumentState) String() string {
	var flags []string
	if s.Has(DocStateParsed) {
		flags = append(flags, "Parsed")
	}
	if s.Has(DocStateExportedSymbols) {
		flags = append(flags, "ExportedSymbols")
	}
	if s.Has(DocStateImportedSymbols) {
		flags = append(flags, "ImportedSymbols")
	}
	if s.Has(DocStateLocalSymbols) {
		flags = append(flags, "LocalSymbols")
	}
	if s.Has(DocStateLinked) {
		flags = append(flags, "Linked")
	}
	if s.Has(DocStateValidated) {
		flags = append(flags, "Validated")
	}
	return strings.Join(flags, " | ")
}

func (s DocumentState) Has(flag DocumentState) bool {
	return s&flag != 0
}

func (s DocumentState) With(flag DocumentState) DocumentState {
	return s | flag
}

func (s DocumentState) Without(flag DocumentState) DocumentState {
	return s &^ flag
}
