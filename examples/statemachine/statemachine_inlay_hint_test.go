// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

func TestInlayHint_TransitionCounts(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))

	doc := f.Parse(`
statemachine Test
events flick
initialState off

state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	provider := service.MustGet[server.InlayHintProvider](f.Services())
	result, err := provider.HandleInlayHintRequest(
		context.Background(),
		&lsp.InlayHintParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 20, Character: 0},
			},
		},
	)
	assert.NoError(t, err)
	assert.Len(t, result, 2, "one hint per state")

	for _, hint := range result {
		assert.Len(t, hint.Label, 1)
		assert.Equal(t, ": 1 transition(s)", hint.Label[0].Value)
		assert.Equal(t, lsp.Type, hint.Kind)
	}
}

func TestInlayHint_RangeFiltering(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))

	doc := f.Parse(`
statemachine Test
events flick
initialState off

state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	provider := service.MustGet[server.InlayHintProvider](f.Services())
	result, err := provider.HandleInlayHintRequest(
		context.Background(),
		&lsp.InlayHintParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
			// Only cover the first couple of lines - neither state's name
			// falls in this range, so no hints should be returned.
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 1, Character: 0},
			},
		},
	)
	assert.NoError(t, err)
	assert.Empty(t, result, "server.NodesInRange should exclude both states")
}
