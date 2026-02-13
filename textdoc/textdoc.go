// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import "typefox.dev/lsp"

// Handle represents a reference to a text document's content and metadata.
type Handle interface {
	// URI returns the associated URI for this document.
	URI() lsp.DocumentURI
	// LanguageID returns the identifier of the language associated with this document.
	LanguageID() string
	// Version returns the version number of this document.
	Version() int32
	// Content returns the document content as a byte slice.
	Content() []byte
	// Text returns the text content or a substring if range is provided.
	Text(r *lsp.Range) string
	// PositionAt converts a zero-based offset to a position.
	PositionAt(offset int) lsp.Position
	// OffsetAt converts a position to a zero-based offset.
	OffsetAt(position lsp.Position) int
}
