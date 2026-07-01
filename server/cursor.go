// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

// NodeAtCursor returns the AST node associated with the token at the given
// cursor position in doc. If the cursor lies between two tokens, the second
// token's node is used as a fallback when the first token has no associated
// node. Returns nil if there is no token at the position, or if neither
// token has an associated node.
func NodeAtCursor(doc *core.Document, position lsp.Position) core.AstNode {
	offset := doc.TextDoc.OffsetAt(position)
	first, second := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil
	}
	if first.Element != nil {
		return first.Element
	}
	if second != nil && second.Element != nil {
		return second.Element
	}
	return nil
}
