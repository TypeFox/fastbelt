// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"github.com/TypeFox/go-lsp/protocol"
)

// TextdocServices contains the services for the textdoc package.
type TextdocServices struct {
	Store Store
}

// LoadDefaultServices creates the default services for the textdoc package.
// If the services are already set, they are not overwritten.
func LoadDefaultServices(s *TextdocServices) {
	if s.Store == nil {
		s.Store = NewDefaultStore()
	}
}

// Handle represents a reference to a text document's content and metadata.
type Handle interface {
	// URI returns the associated URI for this document.
	URI() protocol.DocumentURI
	// LanguageID returns the identifier of the language associated with this document.
	LanguageID() string
	// Version returns the version number of this document.
	Version() int32
	// Content returns the document content as a byte slice.
	Content() []byte
	// PositionAt converts a zero-based offset to a position.
	PositionAt(offset int) protocol.Position
	// OffsetAt converts a position to a zero-based offset.
	OffsetAt(position protocol.Position) int
}
