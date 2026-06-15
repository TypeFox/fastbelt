// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// FoldingRangeProvider is a service for handling LSP folding range requests.
type FoldingRangeProvider interface {
	HandleFoldingRangeRequest(ctx context.Context, params *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error)
}

// FoldingRangeFilter defines the customization points for folding behavior.
type FoldingRangeFilter interface {
	// ShouldProcess determines whether the specified AstNode should be folded.
	ShouldProcess(node core.AstNode) bool

	// IncludeLastFoldingLine determines whether the folding range should include its last line for AST nodes.
	// Return false to exclude the last line (e.g., for nodes ending in closing braces).
	IncludeLastFoldingLine(node core.AstNode) bool

	// IncludeLastFoldingLineForComment determines whether comment folding ranges should include the last line.
	// Return false to exclude the last line (typically desired for multi-line comments).
	IncludeLastFoldingLineForComment() bool
}

// DefaultFoldingRangeFilter provides the standard folding behavior:
// Processes all nodes, excludes last line for nodes ending with closing brackets, and includes last line for comments.
type DefaultFoldingRangeFilter struct{}

func (f *DefaultFoldingRangeFilter) ShouldProcess(node core.AstNode) bool {
	return true
}

func (f *DefaultFoldingRangeFilter) IncludeLastFoldingLine(node core.AstNode) bool {
	if node == nil {
		return true
	}

	tokens := node.Tokens()
	if len(tokens) > 0 {
		lastTokenImage := tokens[len(tokens)-1].Image
		if lastTokenImage == "}" || lastTokenImage == ")" || lastTokenImage == "]" {
			return false
		}
	}

	return true
}

func (f *DefaultFoldingRangeFilter) IncludeLastFoldingLineForComment() bool {
	return true
}

type DefaultFoldingRangeProvider struct {
	sc     *service.Container
	filter FoldingRangeFilter
}

func NewDefaultFoldingRangeProvider(sc *service.Container) FoldingRangeProvider {
	return &DefaultFoldingRangeProvider{
		sc:     sc,
		filter: &DefaultFoldingRangeFilter{},
	}
}

func NewFoldingRangeProviderWithFilter(sc *service.Container, filter FoldingRangeFilter) FoldingRangeProvider {
	return &DefaultFoldingRangeProvider{
		sc:     sc,
		filter: filter,
	}
}

func (p *DefaultFoldingRangeProvider) HandleFoldingRangeRequest(ctx context.Context, params *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil // Document not found
	}
	return p.collectFolding(doc), nil
}

func (p *DefaultFoldingRangeProvider) collectFolding(document *core.Document) []lsp.FoldingRange {
	if document.Root == nil {
		return nil
	}

	foldings := []lsp.FoldingRange{}

	for node := range core.AllChildren(document.Root) {
		if p.filter.ShouldProcess(node) {
			includeLastLine := p.filter.IncludeLastFoldingLine(node)
			if foldingRange := p.toFoldingRange(node.Segment(), "", includeLastLine); foldingRange != nil {
				foldings = append(foldings, *foldingRange)
			}
		}
	}

	p.collectCommentFolding(document, &foldings)

	return foldings
}

func (p *DefaultFoldingRangeProvider) collectCommentFolding(document *core.Document, foldings *[]lsp.FoldingRange) {
	includeLastLine := p.filter.IncludeLastFoldingLineForComment()
	for _, comment := range document.Comments {
		if foldingRange := p.toFoldingRange(&comment.TextSegment, "comment", includeLastLine); foldingRange != nil {
			*foldings = append(*foldings, *foldingRange)
		}
	}
}

func (p *DefaultFoldingRangeProvider) toFoldingRange(segment *core.TextSegment, kind string, includeLastLine bool) *lsp.FoldingRange {
	if segment == nil {
		return nil
	}

	// Minimum 2-line difference required
	if segment.Range.End.Line-segment.Range.Start.Line < 2 {
		return nil
	}

	lspRange := segment.Range.LspRange()
	endLine := lspRange.End.Line
	endChar := lspRange.End.Character

	folding := &lsp.FoldingRange{
		StartLine:      &lspRange.Start.Line,
		StartCharacter: &lspRange.Start.Character,
		EndLine:        &endLine,
		Kind:           kind,
	}

	if !includeLastLine {
		endLine--
	} else {
		// Only set EndCharacter when including the last line.
		// When not set, LSP defaults to end of line.
		folding.EndCharacter = &endChar
	}

	return folding
}
