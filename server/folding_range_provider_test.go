// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

func assertHasFoldingRange(t *testing.T, result []lsp.FoldingRange, markerRange core.TextRange, label string) {
	t.Helper()
	for _, fr := range result {
		if fr.StartLine == nil || fr.EndLine == nil {
			continue
		}

		if *fr.StartLine == uint32(markerRange.Start.Line) && *fr.EndLine == uint32(markerRange.End.Line) {
			if label == "comment" {
				assert.Equal(t, "comment", fr.Kind, "Comment folding should have kind='comment'")
			}
			return
		}
	}
	t.Errorf("Should have folding range for marker '%s' (lines %d-%d)", label, markerRange.Start.Line, markerRange.End.Line)
}

func TestFoldingRangeIntegration(t *testing.T) {
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	SetupDefaultServices(sc)
	sc.Seal()

	fixture := test.New(t, sc)

	grammarText := `grammar Test;

<|first:interface First {
	name string
	value string
	another string|>
}

<|second:interface Second extends First {
	extra string
	more string
	evenMore bool|>
}

<|comment:/* Multi-line comment
   that spans multiple lines
   and should be foldable */|>
<|third:interface Third {
	data []string
	items []string
	flags []bool|>
}`

	doc := fixture.ParseURI(grammarText, "file:///test.fb")
	doc.AssertNoParseErrors()

	provider := service.MustGet[FoldingRangeProvider](sc)
	params := &lsp.FoldingRangeParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: lsp.DocumentURI(doc.Document.URI.DocumentURI()),
		},
	}

	result, err := provider.HandleFoldingRangeRequest(context.Background(), params)
	require.NoError(t, err)
	require.NotNil(t, result)

	// We expect folding ranges for the marked regions plus the grammar root
	expectedLabels := []string{"first", "second", "comment", "third"}
	assert.Equal(t, len(result), len(expectedLabels)+1, "Should have at least %d folding ranges", len(expectedLabels))

	for _, label := range expectedLabels {
		markerRange, ok := doc.MarkerRange(label)
		require.True(t, ok, "Marker '%s' should exist", label)
		assertHasFoldingRange(t, result, markerRange, label)
	}
}
