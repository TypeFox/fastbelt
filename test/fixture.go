// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package test

import (
	"context"
	"strings"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// DefaultMarking is the marker configuration used by [New].
// Also used if a nil *[TestMarking] is passed to [Fixture.WithMarking].
var DefaultMarking = &TestMarking{
	StartRange: "<|",
	EndRange:   "|>",
	StartIndex: "",
	EndIndex:   ">",
	Delimiter:  ":",
}

// TestMarking configures the delimiter syntax for position markers in test content.
// See [DefaultMarking] for the default configuration.
type TestMarking struct {
	StartRange string // Opening delimiter for range markers.
	EndRange   string // Closing delimiter for range markers.
	StartIndex string // Opening delimiter for index markers, defaults to the current [TestMarking.StartRange] if empty.
	EndIndex   string // Closing delimiter for index markers.
	Delimiter  string // Separator between label and text in range markers.
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
	sc      *service.Container
	ctx     context.Context
	marking *TestMarking
}

// New creates a [Fixture] backed by the given service container.
// It uses [DefaultMarking] and [context.Background].
func New(t testing.TB, sc *service.Container) *Fixture {
	t.Helper()
	return &Fixture{t: t, sc: sc, ctx: context.Background(), marking: DefaultMarking}
}

// NewWithContext creates a [Fixture] with an explicit context.
func NewWithContext(t testing.TB, sc *service.Container, ctx context.Context) *Fixture {
	t.Helper()
	return &Fixture{t: t, sc: sc, ctx: ctx, marking: DefaultMarking}
}

// WithMarking returns a shallow copy of the [Fixture] using m as the marker configuration.
func (f *Fixture) WithMarking(m *TestMarking) *Fixture {
	f.t.Helper()
	if m == nil {
		// Fallback to default
		m = DefaultMarking
	} else if m.StartRange == "" {
		f.t.Fatal("Received empty TestMarking.StartRange")
	} else if m.EndRange == "" {
		f.t.Fatal("Received empty TestMarking.EndRange")
	} else if m.EndIndex == "" {
		f.t.Fatal("Received empty TestMarking.EndIndex")
	} else if (m.StartIndex == "" || m.StartIndex == m.StartRange) && m.EndIndex == m.EndRange {
		// Index and range markers are the same, which is an invalid configuration
		f.t.Fatal("Received TestMarking with identical index and range delimiters.")
	} else if m.Delimiter == "" {
		f.t.Fatal("Received empty TestMarking.Delimiter")
	}
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
	if exts := service.MustGet[workspace.FileExtensions](f.sc); len(exts) > 0 {
		uri = "inmemory://test" + exts[0]
	}
	return f.ParseURI(content, uri)
}

// ParseURI builds a single in-memory document with the given URI. Markers are
// extracted before parsing.
func (f *Fixture) ParseURI(content, uri string) *Doc {
	f.t.Helper()
	cleanText, ranges, indices := extractMarkers(content, f.marking)
	languageID := service.MustGet[workspace.LanguageID](f.sc)
	doc, err := core.NewDocumentFromString(uri, string(languageID), cleanText)
	if err != nil {
		f.t.Fatalf("fbtest: failed to create document: %v", err)
	}
	documents := service.MustGet[workspace.DocumentManager](f.sc)
	documents.Set(doc)
	builder := service.MustGet[workspace.Builder](f.sc)
	if err := builder.Build(f.ctx, []*core.Document{doc}, func() {}); err != nil {
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
		languageID := service.MustGet[workspace.LanguageID](f.sc)
		doc, err := core.NewDocumentFromString(uri, string(languageID), cleanText)
		if err != nil {
			f.t.Fatalf("fbtest: failed to create document %q: %v", uri, err)
		}
		documents := service.MustGet[workspace.DocumentManager](f.sc)
		documents.Set(doc)
		coreDocs = append(coreDocs, doc)
		results = append(results, f.newDoc(doc, ranges, indices))
	}
	builder := service.MustGet[workspace.Builder](f.sc)
	if err := builder.Build(f.ctx, coreDocs, func() {}); err != nil {
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
	if m == nil {
		m = DefaultMarking
	}
	var sb strings.Builder
	sb.Grow(len(content))
	offset := 0
	i := 0
	startIndexMarker := m.StartRange
	if m.StartIndex != "" {
		startIndexMarker = m.StartIndex
	}
	separateIndexPrefix := startIndexMarker != m.StartRange
	for i < len(content) {
		nextRange := strings.Index(content[i:], m.StartRange)
		nextIndex := -1
		if separateIndexPrefix {
			nextIndex = strings.Index(content[i:], startIndexMarker)
		}

		// Determine which delimiter comes first. On a tie (same position),
		// prefer the index delimiter when it is longer than the range delimiter
		// (e.g. index "[[" beats range "[").
		useIndex := false
		next := nextRange
		if nextIndex != -1 && (next == -1 || nextIndex < next ||
			(nextIndex == next && len(startIndexMarker) > len(m.StartRange))) {
			useIndex = true
			next = nextIndex
		}
		if next == -1 {
			sb.WriteString(content[i:])
			break
		}
		next += i

		// Write everything up to the marker opening.
		sb.WriteString(content[i:next])
		offset += next - i

		if useIndex {
			// Distinct index prefix: only match index markers.
			rest := content[next+len(startIndexMarker):]
			if idxEnd := strings.Index(rest, m.EndIndex); idxEnd != -1 {
				indices = append(indices, IndexMarker{Label: rest[:idxEnd], Offset: offset})
				i = next + len(startIndexMarker) + idxEnd + len(m.EndIndex)
				continue
			}
			// No closing delimiter; treat as literal.
			sb.WriteString(startIndexMarker)
			offset += len(startIndexMarker)
			i = next + len(startIndexMarker)
		} else {
			rest := content[next+len(m.StartRange):]

			// Try range marker first: look for EndRange in the remainder.
			if rangeEnd := strings.Index(rest, m.EndRange); rangeEnd != -1 {
				inner := rest[:rangeEnd]
				var label, text string
				if before, after, ok := strings.Cut(inner, m.Delimiter); ok {
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

			// Fall back to index marker: only when StartIndex == StartRange (both marker
			// types share the same opening delimiter, distinguished only by EndIndex).
			if !separateIndexPrefix {
				if idxEnd := strings.Index(rest, m.EndIndex); idxEnd != -1 {
					indices = append(indices, IndexMarker{Label: rest[:idxEnd], Offset: offset})
					i = next + len(m.StartRange) + idxEnd + len(m.EndIndex)
					continue
				}
			}

			// No closing delimiter found; treat the opening as literal text.
			sb.WriteString(m.StartRange)
			offset += len(m.StartRange)
			i = next + len(m.StartRange)
		}
	}
	cleanText = sb.String()
	return
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
