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
// Access to the fields of Document should be synchronized using a [typefox.dev/fastbelt/workspace] Lock.
// The document struct should never be copied after creation.
type Document struct {
	// URI identifies the document in the workspace.
	URI URI
	// State tracks which build phases already ran for this document.
	State DocumentState
	// Root is the AST root produced by parsing.
	// It is nil until parsing succeeds.
	Root AstNode
	// Tokens is the token stream produced by lexing.
	// It is empty until lexing succeeds.
	Tokens TokenSlice
	// Comments contains comment tokens extracted by the lexer.
	// It is empty until lexing succeeds.
	Comments TokenSlice
	// LocalSymbols stores symbols declared in this document.
	// It is nil until the local symbols build phase succeeds.
	LocalSymbols LocalSymbols
	// ExportedSymbols stores symbols this document contributes to global linking.
	// It is nil until the exported symbols build phase succeeds.
	ExportedSymbols SymbolContainer
	// ImportedSymbols stores symbols imported from other workspace documents.
	// It is nil until the imported symbols build phase succeeds.
	ImportedSymbols SymbolContainer
	// ParserErrors contains syntax and parser-recovery errors.
	// It is empty until parsing succeeds.
	ParserErrors []*ParserError
	// LexerErrors contains lexical errors found while tokenizing the document text.
	// It is empty until lexing succeeds.
	LexerErrors []*LexerError
	// References contains unresolved and resolved cross-references found during linking.
	// It is empty until the linking build phase succeeds.
	References []UntypedReference
	// ReferenceDescriptions describes references for features such as "find references".
	// It is nil until the reference descriptions build phase succeeds.
	ReferenceDescriptions ReferenceDescriptions
	// TextDoc is the backing text document handle.
	TextDoc textdoc.Handle
	// Diagnostics contains validation and analysis diagnostics for editor clients.
	// It is empty until the validation build phase succeeds.
	Diagnostics []*Diagnostic
	// Data can be used to store arbitrary additional information related to the document.
	// This is not used by the framework itself, but adopters may find it useful.
	// The document builder does not clear this data during the build process.
	// It is the responsibility of the caller to manage it appropriately (e.g. clearing or
	// updating it when the document changes).
	//
	// The map is concurrent to allow storing data from different goroutines without
	// additional synchronization.
	Data sync.Map
}

// NewDocument creates a [Document] from a text document handle.
//
// It initializes all build-related collections to empty values and sets [Document.TextDoc]
// and [Document.URI], but leaves semantic data (such as [Document.Root]) unset until
// the corresponding build phases run.
func NewDocument(textDoc textdoc.Handle) *Document {
	uri := ParseURI(string(textDoc.URI()))
	return &Document{
		URI:                   uri,
		State:                 0,
		TextDoc:               textDoc,
		Root:                  nil,
		LocalSymbols:          nil,
		ExportedSymbols:       nil,
		ImportedSymbols:       nil,
		Data:                  sync.Map{},
		Tokens:                TokenSlice{},
		Comments:              TokenSlice{},
		ParserErrors:          []*ParserError{},
		LexerErrors:           []*LexerError{},
		References:            []UntypedReference{},
		ReferenceDescriptions: nil,
		Diagnostics:           []*Diagnostic{},
	}
}

// NewDocumentFromString creates a [Document] backed by an in-memory text document.
//
// It is useful in tests, benchmarks, and tooling code that needs a document instance
// without going through the workspace file-loading pipeline.
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
	// DocStateParsed marks that lexing and parsing completed and produced AST and token data.
	DocStateParsed DocumentState = 1 << iota // 0x0001
	// DocStateExportedSymbols marks that exported symbols were collected for cross-document linking.
	DocStateExportedSymbols // 0x0002
	// DocStateImportedSymbols marks that symbols from other documents were imported.
	DocStateImportedSymbols // 0x0004
	// DocStateLocalSymbols marks that local symbols for this document were collected.
	DocStateLocalSymbols // 0x0008
	// DocStateLinked marks that cross-references were linked.
	DocStateLinked // 0x0010
	// DocStateReferences marks that reference descriptions were collected.
	DocStateReferences // 0x0020
	// DocStateValidated marks that validators were executed and diagnostics were collected.
	DocStateValidated // 0x0040
)

// String returns a readable representation of the set state flags.
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

// Has reports whether flag is set in s.
func (s DocumentState) Has(flag DocumentState) bool {
	return s&flag != 0
}

// With returns s with flag set.
func (s DocumentState) With(flag DocumentState) DocumentState {
	return s | flag
}

// Without returns s with flag cleared.
func (s DocumentState) Without(flag DocumentState) DocumentState {
	return s &^ flag
}
