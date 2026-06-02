// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	"typefox.dev/fastbelt/test"
)

func TestHighlightOnEventName(t *testing.T) {
	f := test.New(t, CreateLspServices())
	doc := f.Parse(`
		statemachine Toggle

		events <|flick><|flick|>

		initialState a

		state a
		  <|flick|> => b
		end

		state b
		  <|flick|> => a
		end
	`).AssertNoErrors()

	doc.AssertHighlights("flick")
}

func TestHighlightOnEventReference(t *testing.T) {
	f := test.New(t, CreateLspServices())
	doc := f.Parse(`
		statemachine Toggle

		events <|flick|>

		initialState a

		state a
		  <|flick><|flick|> => b
		end

		state b
		  <|flick|> => a
		end
	`).AssertNoErrors()

	doc.AssertHighlights("flick")
}

func TestHighlightOnStateName(t *testing.T) {
	f := test.New(t, CreateLspServices())
	doc := f.Parse(`
		statemachine Toggle

		events flick

		initialState <|a|>

		state <|a><|a|>
		  flick => b
		end

		state b
		  flick => <|a|>
		end
	`).AssertNoErrors()

	doc.AssertHighlights("a")
}

// An event and a command can share the same name. Highlighting one kind of
// element must not bleed into the equally named element of the other kind.
func TestHighlightDistinguishesCommandFromEvent(t *testing.T) {
	f := test.New(t, CreateLspServices())
	doc := f.Parse(`
		statemachine Toggle

		events <|ev><|ev:flick|>

		commands <|cmd><|cmd:flick|>

		initialState a

		state a
		  actions { <|cmd:flick|> }
		  <|ev:flick|> => b
		end

		state b
		  <|ev:flick|> => a
		end
	`).AssertNoErrors()

	// Cursor on the command declaration: only command occurrences are highlighted.
	doc.AssertHighlights("cmd")
	// Cursor on the event declaration: only event occurrences are highlighted.
	doc.AssertHighlights("ev")
}
