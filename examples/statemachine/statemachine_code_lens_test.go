// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// resolveAll simulates the client: request lenses, then resolve every one of
// them (as if all were visible in the viewport), returning the resolved titles.
func resolveAll(t *testing.T, provider server.CodeLensProvider, lenses []lsp.CodeLens) []string {
	t.Helper()
	resolving, ok := provider.(server.ResolvingCodeLensProvider)
	require.True(t, ok, "stateMachineCodeLensProvider should implement ResolvingCodeLensProvider")

	var titles []string
	for _, lens := range lenses {
		resolved, err := resolving.HandleCodeLensResolveRequest(context.Background(), &lens)
		require.NoError(t, err)
		require.NotNil(t, resolved.Command, "resolved lens should have a Command")
		titles = append(titles, resolved.Command.Title)
	}
	return titles
}

func TestCodeLens_ReferenceCounts(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))

	doc := f.Parse(`
statemachine Test
events flick
commands beep
initialState off

state off
  actions { beep }
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	provider := service.MustGet[server.CodeLensProvider](f.Services())
	result, err := provider.HandleCodeLensRequest(
		context.Background(),
		&lsp.CodeLensParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
		},
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// HandleCodeLensRequest defers counting - lenses carry Data, not a
	// Command, until HandleCodeLensResolveRequest is called for them.
	for _, lens := range result {
		assert.Nil(t, lens.Command)
		assert.NotNil(t, lens.Data)
	}

	titles := map[string]int{}
	for _, title := range resolveAll(t, provider, result) {
		titles[title]++
	}

	// States get no lens here: their transition count is already shown
	// inline via the InlayHint, so a duplicate lens would be redundant.
	assert.Equal(t, 0, titles["1 transition(s)"], "states should not get a code lens")

	// flick: referenced once per state's transition = 2 references
	// beep: referenced once (actions { beep }) = 1 reference
	assert.Equal(t, 1, titles["2 reference(s)"], "the flick event is referenced by both transitions")
	assert.Equal(t, 1, titles["1 reference(s)"], "the beep command is referenced once")

	// Total lenses: 1 event + 1 command = 2 (no state lenses)
	assert.Len(t, result, 2)
}

func TestCodeLens_EmptyForUnreferencedDeclarations(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))

	doc := f.Parse(`
statemachine Test
events unused
initialState off

state off
end
`).AssertNoErrors()

	provider := service.MustGet[server.CodeLensProvider](f.Services())
	result, err := provider.HandleCodeLensRequest(
		context.Background(),
		&lsp.CodeLensParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
		},
	)
	assert.NoError(t, err)
	require.Len(t, result, 1, "the unused event still gets a lens, just with a zero count once resolved")

	titles := resolveAll(t, provider, result)
	assert.Equal(t, "0 reference(s)", titles[0])
}
