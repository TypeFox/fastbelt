// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package fbtest provides utilities for testing Fastbelt language implementations.
//
// Create a [Fixture] from your language's service container, then call [Fixture.Parse]
// or [Fixture.ParseAll] to build documents. Use [Doc] assertion methods and the
// generic Find/MustFind functions to inspect results.
//
// # Markers
//
// Content strings may embed position markers that are stripped before parsing.
// Positions are recorded relative to the cleaned text.
//
//   - Range marker:  <|label:text|>  — spans "text"; label identifies the range.
//   - Range shorthand: <|label|>     — equivalent to <|label:label|>.
//   - Empty range:   <|label:|>      — zero-length span at the marker position.
//   - Index marker:  <|label>        — single position (no text content).
//
// Range markers are tried before index markers when both share the same opening
// delimiter (the default). Marker delimiters are configurable via [TestMarking].
package fbtest

import (
	"context"
	"strings"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/workspace"
)

// DefaultMarking is the marker configuration used by [New].
var DefaultMarking = &TestMarking{
	StartRange: "<|",
	EndRange:   "|>",
	StartIndex: "<|",
	EndIndex:   ">",
}

// TestMarking configures the delimiter syntax for position markers in test content.
type TestMarking struct {
	StartRange string // opening delimiter for range markers, default "<|"
	EndRange   string // closing delimiter for range markers, default "|>"
	StartIndex string // opening delimiter for index markers, default "<|"
	EndIndex   string // closing delimiter for index markers, default ">"
}

// RangeMarker records a labeled span extracted from test content.
// Start and End are byte offsets into the clean text (markers removed).
type RangeMarker struct {
	Label string
	Start int // inclusive
	End   int // exclusive
}

// IndexMarker records a labeled position extracted from test content.
// Offset is a byte offset into the clean text (markers removed).
type IndexMarker struct {
	Label  string
	Offset int
}

// Fixture is the entry point for language tests. Create one per test with [New].
type Fixture struct {
	t       testing.TB
	srv     workspace.WorkspaceSrvCont
	ctx     context.Context
	marking *TestMarking
}

// New creates a Fixture backed by the given service container.
// It uses [DefaultMarking] and [context.Background].
func New(t testing.TB, srv workspace.WorkspaceSrvCont) *Fixture {
	t.Helper()
	return &Fixture{t: t, srv: srv, ctx: context.Background(), marking: DefaultMarking}
}

// NewWithContext creates a Fixture with an explicit context.
func NewWithContext(t testing.TB, srv workspace.WorkspaceSrvCont, ctx context.Context) *Fixture {
	t.Helper()
	return &Fixture{t: t, srv: srv, ctx: ctx, marking: DefaultMarking}
}

// WithMarking returns a shallow copy of the Fixture using m as the marker configuration.
func (f *Fixture) WithMarking(m *TestMarking) *Fixture {
	f.t.Helper()
	c := *f
	c.marking = m
	return &c
}

// Parse builds a single in-memory document from content. Markers are extracted before
// parsing. The URI is "inmemory://test<ext>" where <ext> is the first entry in
// WorkspaceSrv.FileExtensions (e.g. ".statemachine"), or "inmemory://test" if none
// are configured.
func (f *Fixture) Parse(content string) *Doc {
	f.t.Helper()
	uri := "inmemory://test"
	if exts := f.srv.Workspace().FileExtensions; len(exts) > 0 {
		uri = "inmemory://test" + exts[0]
	}
	return f.ParseURI(content, uri)
}

// ParseURI builds a single in-memory document with the given URI. Markers are
// extracted before parsing.
func (f *Fixture) ParseURI(content, uri string) *Doc {
	f.t.Helper()
	cleanText, ranges, indices := extractMarkers(content, f.marking)
	doc, err := core.NewDocumentFromString(uri, f.srv.Workspace().LanguageID, cleanText)
	if err != nil {
		f.t.Fatalf("fbtest: failed to create document: %v", err)
	}
	f.srv.Workspace().DocumentManager.Set(doc)
	if err := f.srv.Workspace().Builder.Build(f.ctx, []*core.Document{doc}, func() {}); err != nil {
		f.t.Fatalf("fbtest: build failed: %v", err)
	}
	return f.newDoc(doc, ranges, indices)
}

// ParseAll builds multiple documents together, enabling cross-document reference
// resolution. Arguments are alternating URI/content pairs: uri1, content1, uri2,
// content2, ... All documents are registered in the DocumentManager before building.
// Results are returned in the same order as the pairs.
func (f *Fixture) ParseAll(uriContentPairs ...string) []*Doc {
	f.t.Helper()
	if len(uriContentPairs)%2 != 0 {
		f.t.Fatalf("fbtest: ParseAll requires an even number of arguments (uri, content pairs)")
	}
	n := len(uriContentPairs) / 2
	coreDocs := make([]*core.Document, 0, n)
	results := make([]*Doc, 0, n)
	for i := 0; i < len(uriContentPairs); i += 2 {
		uri, content := uriContentPairs[i], uriContentPairs[i+1]
		cleanText, ranges, indices := extractMarkers(content, f.marking)
		doc, err := core.NewDocumentFromString(uri, f.srv.Workspace().LanguageID, cleanText)
		if err != nil {
			f.t.Fatalf("fbtest: failed to create document %q: %v", uri, err)
		}
		f.srv.Workspace().DocumentManager.Set(doc)
		coreDocs = append(coreDocs, doc)
		results = append(results, f.newDoc(doc, ranges, indices))
	}
	if err := f.srv.Workspace().Builder.Build(f.ctx, coreDocs, func() {}); err != nil {
		f.t.Fatalf("fbtest: build failed: %v", err)
	}
	return results
}

func (f *Fixture) newDoc(doc *core.Document, ranges []RangeMarker, indices []IndexMarker) *Doc {
	return &Doc{Document: doc, Ranges: ranges, Indices: indices, ctx: f.ctx, t: f.t}
}

// extractMarkers scans content for embedded position markers, removes them, and
// records their locations relative to the cleaned text. Range markers are matched
// before index markers when they share the same opening delimiter.
//
// If no closing delimiter is found after an opening, the opening is treated as
// literal text.
func extractMarkers(content string, m *TestMarking) (cleanText string, ranges []RangeMarker, indices []IndexMarker) {
	var sb strings.Builder
	sb.Grow(len(content))
	offset := 0
	i := 0
	for i < len(content) {
		next := strings.Index(content[i:], m.StartRange)
		if next == -1 {
			sb.WriteString(content[i:])
			break
		}
		next += i
		// Write everything up to the marker opening.
		sb.WriteString(content[i:next])
		offset += next - i
		rest := content[next+len(m.StartRange):]

		// Try range marker first: look for EndRange in the remainder.
		if rangeEnd := strings.Index(rest, m.EndRange); rangeEnd != -1 {
			inner := rest[:rangeEnd]
			var label, text string
			if before, after, ok := strings.Cut(inner, ":"); ok {
				label, text = before, after
			} else {
				// Shorthand: <|label|> — label and text are the same string.
				label, text = inner, inner
			}
			ranges = append(ranges, RangeMarker{Label: label, Start: offset, End: offset + len(text)})
			sb.WriteString(text)
			offset += len(text)
			i = next + len(m.StartRange) + rangeEnd + len(m.EndRange)
			continue
		}

		// Fall back to index marker: look for EndIndex.
		if idxEnd := strings.Index(rest, m.EndIndex); idxEnd != -1 {
			indices = append(indices, IndexMarker{Label: rest[:idxEnd], Offset: offset})
			i = next + len(m.StartIndex) + idxEnd + len(m.EndIndex)
			continue
		}

		// No closing delimiter found; treat the opening as literal text.
		sb.WriteString(m.StartRange)
		offset += len(m.StartRange)
		i = next + len(m.StartRange)
	}
	cleanText = sb.String()
	return
}

// offsetToLocation converts a byte offset into text to a zero-based line/column
// [core.TextLocation]. Multi-byte characters are counted by byte.
func offsetToLocation(text string, offset int) core.TextLocation {
	line := core.TextLine(0)
	col := core.TextColumn(0)
	for i := 0; i < offset && i < len(text); i++ {
		if text[i] == '\n' {
			line++
			col = 0
		} else {
			col++
		}
	}
	return core.TextLocation{Line: line, Column: col}
}

// locationInRange reports whether loc falls within r.
// For point ranges (start == end), loc must equal start exactly.
// For non-empty ranges, the interval is [start, end).
func locationInRange(loc core.TextLocation, r core.TextRange) bool {
	if r.Start == r.End {
		return loc == r.Start
	}
	if loc.Line < r.Start.Line || loc.Line > r.End.Line {
		return false
	}
	if loc.Line == r.Start.Line && loc.Column < r.Start.Column {
		return false
	}
	if loc.Line == r.End.Line && loc.Column >= r.End.Column {
		return false
	}
	return true
}
