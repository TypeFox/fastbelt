// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/lsp"
)

// TextRange represents a text range in a document. The indices are zero-based and represent byte offsets.
// It is defined by a start and end offset, where the start is inclusive and the end is exclusive.
type TextRange struct {
	Start int32
	End   int32
}

func NewTextRange(start, end int) TextRange {
	return TextRange{Start: int32(start), End: int32(end)}
}

// LspRange converts the TextRange to an LSP range using the provided text document handle.
// The [lsp.Range] returned is suitable for use in LSP responses without further conversion.
func (r TextRange) LspRange(doc textdoc.Handle) lsp.Range {
	return lsp.Range{
		Start: doc.PositionAt(int(r.Start)),
		End:   doc.PositionAt(int(r.End)),
	}
}
