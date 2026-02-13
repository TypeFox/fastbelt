// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import "typefox.dev/lsp"

type TextIndex int32
type TextLine int32
type TextColumn int32

type TextIndexRange struct {
	Start TextIndex
	End   TextIndex
}

type TextLocation struct {
	Line   TextLine
	Column TextColumn
}

func (l TextLocation) LspPosition() lsp.Position {
	return lsp.Position{
		Line:      uint32(l.Line),
		Character: uint32(l.Column),
	}
}

type TextRange struct {
	Start TextLocation
	End   TextLocation
}

func (r TextRange) LspRange() lsp.Range {
	return lsp.Range{
		Start: r.Start.LspPosition(),
		End:   r.End.LspPosition(),
	}
}

type TextSegment struct {
	Indices TextIndexRange
	Range   TextRange
}
