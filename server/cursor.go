// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"iter"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
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

// TargetAtCursor resolves the name or reference at the cursor position in
// params to the AST node it refers to - the declaration itself, or the
// declaration a reference resolves to - along with the source range of the
// name or reference. Unlike [NodeAtCursor], it resolves references via
// [NameFinder]. Returns a nil target if there is no document, no token at
// the position, or no resolvable name.
func TargetAtCursor(ctx context.Context, sc *service.Container, params *lsp.TextDocumentPositionParams) (target core.AstNode, sourceRange lsp.Range) {
	documentManager := service.MustGet[workspace.DocumentManager](sc)
	doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
	if doc == nil {
		return nil, lsp.Range{}
	}

	offset := doc.TextDoc.OffsetAt(params.Position)
	first, second := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, lsp.Range{}
	}

	nameFinder := service.MustGet[NameFinder](sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, lsp.Range{}
	}

	target = foundName.Target.Owner()
	if target == nil {
		return nil, lsp.Range{}
	}
	return target, foundName.Source.TextRange().LspRange(doc.TextDoc)
}

// NodesInRange iterates over every AST node in doc whose text range overlaps
// the given LSP range. Nodes are yielded in the same order as core.AllNodes.
func NodesInRange(doc *core.Document, r lsp.Range) iter.Seq[core.AstNode] {
	startOffset := doc.TextDoc.OffsetAt(r.Start)
	endOffset := doc.TextDoc.OffsetAt(r.End)
	return func(yield func(core.AstNode) bool) {
		for node := range core.AllNodes(doc.Root) {
			rng := node.TextRange()
			nodeStart := int(rng.Start)
			nodeEnd := int(rng.End)
			if nodeEnd <= startOffset || nodeStart >= endOffset {
				continue
			}
			if !yield(node) {
				return
			}
		}
	}
}
