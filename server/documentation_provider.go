// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"sort"
	"strings"

	core "typefox.dev/fastbelt"
)

// DocumentationProvider extracts and formats documentation comments for AST nodes.
type DocumentationProvider interface {
	Documentation(node core.AstNode) string
}

// DocumentationTrimmer strips comment markers from a raw comment token image.
// Implement this interface to support comment styles beyond the defaults (// and /* */).
type DocumentationTrimmer interface {
	TrimComment(content string) string
}

// DefaultDocumentationTrimmer strips // single-line and /* */ block comment markers.
type DefaultDocumentationTrimmer struct{}

// TrimComment removes comment markers from a raw comment token image.
// For // single-line comments the prefix is stripped and surrounding whitespace trimmed.
// For /* */ block comments the delimiters and leading * prefixes are removed per line.
func (t *DefaultDocumentationTrimmer) TrimComment(content string) string {
	trimmed := strings.TrimSpace(content)
	if strings.HasPrefix(trimmed, "/*") {
		inner := strings.TrimSuffix(strings.TrimPrefix(trimmed, "/*"), "*/")
		lines := strings.Split(inner, "\n")
		result := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
			result = append(result, line)
		}
		for len(result) > 0 && result[0] == "" {
			result = result[1:]
		}
		for len(result) > 0 && result[len(result)-1] == "" {
			result = result[:len(result)-1]
		}
		return strings.Join(result, "\n")
	}
	// Single-line // comment: strip prefix and trim surrounding whitespace.
	return strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))
}

// DefaultDocumentationProvider is the default implementation of [DocumentationProvider].
type DefaultDocumentationProvider struct {
	trimmer DocumentationTrimmer
}

// NewDefaultDocumentationProvider creates a DocumentationProvider with the default trimmer.
func NewDefaultDocumentationProvider() DocumentationProvider {
	return &DefaultDocumentationProvider{trimmer: &DefaultDocumentationTrimmer{}}
}

// NewDocumentationProviderWithTrimmer creates a DocumentationProvider with a custom trimmer.
func NewDocumentationProviderWithTrimmer(trimmer DocumentationTrimmer) DocumentationProvider {
	return &DefaultDocumentationProvider{trimmer: trimmer}
}

// Documentation returns the documentation comment block immediately preceding node,
// with comment markers stripped. Returns an empty string when no attached comment exists.
func (s *DefaultDocumentationProvider) Documentation(node core.AstNode) string {
	doc := node.Document()
	if doc == nil {
		return ""
	}

	nodeTokens := node.Tokens()
	if len(nodeTokens) == 0 {
		return ""
	}
	nodeStart := int(nodeTokens[0].TextSegment.Indices.Start)

	// Find the lower bound: end of the code token immediately before this node.
	// doc.Tokens is sorted by offset; comments live in doc.Comments, not doc.Tokens.
	lowerBound := 0
	i := sort.Search(len(doc.Tokens), func(k int) bool {
		return int(doc.Tokens[k].TextSegment.Indices.Start) >= nodeStart
	})
	if i > 0 {
		lowerBound = int(doc.Tokens[i-1].TextSegment.Indices.End)
	}

	// Collect comments in [lowerBound, nodeStart).
	first := sort.Search(len(doc.Comments), func(j int) bool {
		return int(doc.Comments[j].TextSegment.Indices.Start) >= lowerBound
	})
	last := sort.Search(len(doc.Comments), func(j int) bool {
		return int(doc.Comments[j].TextSegment.Indices.Start) >= nodeStart
	}) - 1
	if last < first {
		return ""
	}

	// Only include the last contiguous block — stop at any blank line between comments.
	blockStart := last
	for blockStart > first {
		prev := doc.Comments[blockStart-1]
		curr := doc.Comments[blockStart]
		if int(curr.TextSegment.Range.Start.Line)-int(prev.TextSegment.Range.End.Line) > 1 {
			break
		}
		blockStart--
	}

	var parts []string
	for k := blockStart; k <= last; k++ {
		parts = append(parts, s.trimmer.TrimComment(doc.Comments[k].Image))
	}
	return strings.Join(parts, "  \n")
}
