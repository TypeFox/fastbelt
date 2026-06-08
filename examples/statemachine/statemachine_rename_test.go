// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

func TestRenameOnDefinition(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(`
		statemachine Toggle

		events <|flick|>

		initialState a

		state a
		  flick => b
		end

		state b
		  flick => a
		end
	`).AssertNoErrors()

	doc.RunRename("flick", "flip").AssertNoErrors()

	flipTest(t, doc)
}

func TestRenameOnReference(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(`
		statemachine Toggle

		events flick

		initialState a

		state a
		  <|flick|> => b
		end

		state b
		  flick => a
		end
	`).AssertNoErrors()

	doc.RunRename("flick", "flip").AssertNoErrors()

	flipTest(t, doc)
}

func flipTest(t *testing.T, doc *test.Doc) {
	sm := test.MustFindNode[Statemachine](doc)
	// Still just one event
	require.Len(t, sm.Events(), 1)
	event := sm.Events()[0]
	assert.Equal(t, "flip", event.Name())
	require.Len(t, sm.States(), 2)
	for _, state := range sm.States() {
		require.Len(t, state.Transitions(), 1)
		transition := state.Transitions()[0]
		assert.Equal(t, "flip", transition.Event().Unit().String())
		assert.NotNil(t, transition.Event().Ref(doc.Ctx()))
	}
}
