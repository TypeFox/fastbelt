// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractMarkers_RangeMarker(t *testing.T) {
	clean, ranges, indices := extractMarkers("hello <|foo:world|> end", nil)
	assert.Equal(t, "hello world end", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "foo", Start: 6, End: 11}, ranges[0])
	assert.Empty(t, indices)
}

func TestExtractMarkers_RangeShorthand(t *testing.T) {
	// <|label|> is shorthand for <|label:label|>.
	clean, ranges, indices := extractMarkers("<|abc|>", nil)
	assert.Equal(t, "abc", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "abc", Start: 0, End: 3}, ranges[0])
	assert.Empty(t, indices)
}

func TestExtractMarkers_EmptyRange(t *testing.T) {
	// <|label:|> — empty range at the marker position.
	clean, ranges, indices := extractMarkers("ab<|x:|>cd", nil)
	assert.Equal(t, "abcd", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "x", Start: 2, End: 2}, ranges[0])
	assert.Empty(t, indices)
}

func TestExtractMarkers_IndexMarker(t *testing.T) {
	// <|label> — index marker (no EndRange found, falls back to EndIndex ">").
	clean, ranges, indices := extractMarkers("ab<|cursor>cd", nil)
	assert.Equal(t, "abcd", clean)
	assert.Empty(t, ranges)
	require.Len(t, indices, 1)
	assert.Equal(t, IndexMarker{Label: "cursor", Offset: 2}, indices[0])
}

func TestExtractMarkers_MultipleMarkers(t *testing.T) {
	clean, ranges, indices := extractMarkers("<|a:foo|> <|b>bar", nil)
	assert.Equal(t, "foo bar", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "a", Start: 0, End: 3}, ranges[0])
	require.Len(t, indices, 1)
	// "bar" is 4 bytes after "foo " (offset 4).
	assert.Equal(t, IndexMarker{Label: "b", Offset: 4}, indices[0])
}

func TestExtractMarkers_NoMarkers(t *testing.T) {
	clean, ranges, indices := extractMarkers("plain text", nil)
	assert.Equal(t, "plain text", clean)
	assert.Empty(t, ranges)
	assert.Empty(t, indices)
}

func TestExtractMarkers_UnclosedOpening(t *testing.T) {
	// No closing delimiter: treat the opening as literal text.
	clean, ranges, indices := extractMarkers("a <| b", nil)
	assert.Equal(t, "a <| b", clean)
	assert.Empty(t, ranges)
	assert.Empty(t, indices)
}

func TestExtractMarkers_AdjacentMarkers(t *testing.T) {
	// Two range markers with no gap between them.
	clean, ranges, indices := extractMarkers("<|a:x|><|b:y|>", nil)
	assert.Equal(t, "xy", clean)
	require.Len(t, ranges, 2)
	assert.Equal(t, RangeMarker{Label: "a", Start: 0, End: 1}, ranges[0])
	assert.Equal(t, RangeMarker{Label: "b", Start: 1, End: 2}, ranges[1])
	assert.Empty(t, indices)
}

func TestExtractMarkers_OffsetAccountsForRemovedMarkers(t *testing.T) {
	// Marker at position 3 in clean text ("ab" then "cd" after the marker).
	clean, ranges, _ := extractMarkers("ab<|mark:cd|>ef", nil)
	assert.Equal(t, "abcdef", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "mark", Start: 2, End: 4}, ranges[0])
}

// Ensures that we can provide different StartRange and StartIndex delimiters
func TestExtractMarkers_CustomMarkingDistinctIndexPrefix(t *testing.T) {
	m := &TestMarking{
		StartRange: "{{",
		EndRange:   "}}",
		StartIndex: "[[",
		EndIndex:   "]]",
		Delimiter:  ":",
	}

	t.Run("range only", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("a {{lbl:foo}} b", m)
		assert.Equal(t, "a foo b", clean)
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "lbl", Start: 2, End: 5}, ranges[0])
		assert.Empty(t, indices)
	})

	t.Run("index only", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("a [[pos]] b", m)
		assert.Equal(t, "a  b", clean)
		assert.Empty(t, ranges)
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "pos", Offset: 2}, indices[0])
	})

	t.Run("range before index", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("{{r:x}} [[i]]", m)
		assert.Equal(t, "x ", clean)
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "r", Start: 0, End: 1}, ranges[0])
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "i", Offset: 2}, indices[0])
	})

	t.Run("index before range", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[[i]] {{r:x}}", m)
		assert.Equal(t, " x", clean)
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "i", Offset: 0}, indices[0])
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "r", Start: 1, End: 2}, ranges[0])
	})

	t.Run("unclosed index prefix treated as literal", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[[nope", m)
		assert.Equal(t, "[[nope", clean)
		assert.Empty(t, ranges)
		assert.Empty(t, indices)
	})
}

// Ensures that we deal correctly with the case where StartIndex is a substring of StartRange.
func TestExtractMarkers_CustomMarkingSubstringIndex(t *testing.T) {
	m := &TestMarking{
		StartRange: "[[",
		EndRange:   "]]",
		StartIndex: "[",
		EndIndex:   "]",
		Delimiter:  ":",
	}

	t.Run("range only", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("a [[lbl:foo]] b", m)
		assert.Equal(t, "a foo b", clean)
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "lbl", Start: 2, End: 5}, ranges[0])
		assert.Empty(t, indices)
	})

	t.Run("index only", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("a [pos] b", m)
		assert.Equal(t, "a  b", clean)
		assert.Empty(t, ranges)
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "pos", Offset: 2}, indices[0])
	})

	t.Run("range before index", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[[r:x]] [i]", m)
		assert.Equal(t, "x ", clean)
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "r", Start: 0, End: 1}, ranges[0])
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "i", Offset: 2}, indices[0])
	})

	t.Run("index before range", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[i] [[r:x]]", m)
		assert.Equal(t, " x", clean)
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "i", Offset: 0}, indices[0])
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "r", Start: 1, End: 2}, ranges[0])
	})

	t.Run("unclosed index prefix treated as literal", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[[nope", m)
		assert.Equal(t, "[[nope", clean)
		assert.Empty(t, ranges)
		assert.Empty(t, indices)
	})
}

// Ensures that we deal correctly with the case where StartRange is a substring of StartIndex.
func TestExtractMarkers_CustomMarkingSubstringRange(t *testing.T) {
	m := &TestMarking{
		StartRange: "[",
		EndRange:   "]",
		StartIndex: "[[",
		EndIndex:   "]]",
		Delimiter:  ":",
	}

	t.Run("range only", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("a [lbl:foo] b", m)
		assert.Equal(t, "a foo b", clean)
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "lbl", Start: 2, End: 5}, ranges[0])
		assert.Empty(t, indices)
	})

	t.Run("index only", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("a [[pos]] b", m)
		assert.Equal(t, "a  b", clean)
		assert.Empty(t, ranges)
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "pos", Offset: 2}, indices[0])
	})

	t.Run("range before index", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[r:x] [[i]]", m)
		assert.Equal(t, "x ", clean)
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "r", Start: 0, End: 1}, ranges[0])
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "i", Offset: 2}, indices[0])
	})

	t.Run("index before range", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[[i]] [r:x]", m)
		assert.Equal(t, " x", clean)
		require.Len(t, indices, 1)
		assert.Equal(t, IndexMarker{Label: "i", Offset: 0}, indices[0])
		require.Len(t, ranges, 1)
		assert.Equal(t, RangeMarker{Label: "r", Start: 1, End: 2}, ranges[0])
	})

	t.Run("unclosed index prefix treated as literal", func(t *testing.T) {
		clean, ranges, indices := extractMarkers("[[nope", m)
		assert.Equal(t, "[[nope", clean)
		assert.Empty(t, ranges)
		assert.Empty(t, indices)
	})
}

func TestExtractMarkers_CustomMarkingDefaultIndexFallback(t *testing.T) {
	// When StartIndex is empty it falls back to StartRange, so both marker types
	// share the same opening delimiter and are distinguished by EndRange vs EndIndex.
	m := &TestMarking{
		StartRange: "![",
		EndRange:   "]!",
		StartIndex: "", // same as StartRange
		EndIndex:   "]",
		Delimiter:  ":",
	}

	clean, ranges, indices := extractMarkers("![r:foo]! ![idx]", m)
	assert.Equal(t, "foo ", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "r", Start: 0, End: 3}, ranges[0])
	require.Len(t, indices, 1)
	assert.Equal(t, IndexMarker{Label: "idx", Offset: 4}, indices[0])
}

func TestExtractMarkers_CustomDelimiter(t *testing.T) {
	m := &TestMarking{
		StartRange: "<|",
		EndRange:   "|>",
		StartIndex: "<|",
		EndIndex:   ">",
		Delimiter:  ";",
	}

	clean, ranges, indices := extractMarkers("a <|r;foo|> <|i;idx>", m)
	assert.Equal(t, "a foo ", clean)
	require.Len(t, ranges, 1)
	assert.Equal(t, RangeMarker{Label: "r", Start: 2, End: 5}, ranges[0])
	require.Len(t, indices, 1)
	assert.Equal(t, IndexMarker{Label: "i;idx", Offset: 6}, indices[0])
}
