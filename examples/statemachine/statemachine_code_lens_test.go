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

	titles := map[string]int{}
	for _, lens := range result {
		assert.NotNil(t, lens.Command)
		titles[lens.Command.Title]++
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
	assert.Len(t, result, 1, "the unused event still gets a lens, just with a zero count")
	assert.Equal(t, "0 reference(s)", result[0].Command.Title)
}
