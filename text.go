// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import "typefox.dev/lsp"

// TextIndex is a zero-based byte offset in source text.
type TextIndex int32

// TextLine is a zero-based line number in source text.
type TextLine int32

// TextColumn is a zero-based byte column within a line.
type TextColumn int32

// A TextIndexRange describes a half-open byte range [Start, End).
type TextIndexRange struct {
	// Start is the inclusive start byte offset.
	Start TextIndex
	// End is the exclusive end byte offset.
	End TextIndex
}

// A TextLocation identifies a position in source text.
type TextLocation struct {
	// Line is the zero-based line number.
	Line TextLine
	// Column is the zero-based byte column in Line.
	Column TextColumn
}

// LspPosition returns l as an [lsp.Position] using the same coordinates.
func (l TextLocation) LspPosition() lsp.Position {
	return lsp.Position{
		Line:      uint32(l.Line),
		Character: uint32(l.Column),
	}
}

// A TextRange describes a half-open range from Start to End.
type TextRange struct {
	// Start is the inclusive start location.
	Start TextLocation
	// End is the exclusive end location.
	End TextLocation
}

// LspRange returns r as an [lsp.Range] using the same boundaries.
func (r TextRange) LspRange() lsp.Range {
	return lsp.Range{
		Start: r.Start.LspPosition(),
		End:   r.End.LspPosition(),
	}
}

// A TextSegment combines byte offsets and line/column locations for one span.
type TextSegment struct {
	// Indices stores the span as byte offsets.
	Indices TextIndexRange
	// Range stores the same span as line and column locations.
	Range TextRange
}
