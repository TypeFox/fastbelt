// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"strings"
	"unicode"

	core "typefox.dev/fastbelt"
)

// DocumentationProvider extracts and formats documentation comments for AST nodes.
type DocumentationProvider interface {
	Documentation(node core.AstNode) string
}

// DefaultDocumentationProvider is the default implementation of [DocumentationProvider].
type DefaultDocumentationProvider struct{}

func NewDefaultDocumentationProvider() DocumentationProvider {
	return &DefaultDocumentationProvider{}
}

// Documentation returns the documentation comment block immediately preceding node,
// with comment markers stripped. Returns an empty string when no attached comment exists.
func (s *DefaultDocumentationProvider) Documentation(node core.AstNode) string {
	doc := node.Document()
	if doc == nil {
		return ""
	}

	targetStart := int(node.Segment().Indices.Start)

	var block []core.Token
	prevEnd := targetStart
	for i := len(doc.Comments) - 1; i >= 0; i-- {
		c := doc.Comments[i]
		commentEnd := int(c.TextSegment.Indices.End)
		if commentEnd > prevEnd {
			continue
		}
		if hasTokenInRange(doc.Tokens, commentEnd, prevEnd) {
			break
		}
		block = append(block, c)
		prevEnd = int(c.TextSegment.Indices.Start)
	}

	if len(block) == 0 {
		return ""
	}

	// Reverse to restore chronological order.
	for i, j := 0, len(block)-1; i < j; i, j = i+1, j-1 {
		block[i], block[j] = block[j], block[i]
	}

	var lines []string
	for _, c := range block {
		if line := stripCommentMarkers(c.Image); line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "  \n")
}

// hasTokenInRange reports whether any token in ts starts in the half-open range [lo, hi).
// lo is the end of the preceding comment; hi is the start of the target node.
func hasTokenInRange(ts core.TokenSlice, lo, hi int) bool {
	left, right := 0, len(ts)-1
	for left <= right {
		mid := (left + right) / 2
		if int(ts[mid].TextSegment.Indices.Start) < lo {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left < len(ts) && int(ts[left].TextSegment.Indices.Start) < hi
}

// stripCommentMarkers removes comment markers from a raw comment token image.
// Leading and trailing non-alphanumeric runes are trimmed from the whole image
// and then again from each individual line, stripping interior markers (e.g. the
// * in /* */ blocks). Works for any comment style (// # -- /* */ etc.).
func stripCommentMarkers(raw string) string {
	isMarker := func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}
	raw = strings.TrimFunc(strings.TrimSpace(raw), isMarker)
	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		if line = strings.TrimFunc(line, isMarker); line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
