// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	"typefox.dev/fastbelt/test"
)

func TestDefinitionOnName(t *testing.T) {
	f := test.New(t, CreateLspServices())
	doc := f.Parse(`
		statemachine Toggle

		events <|flick><|flick|>

		initialState a

		state a
		  flick => b
		end

		state b
		  flick => a
		end
	`).AssertNoErrors()

	doc.AssertDefinition("flick")
}

func TestDefinitionOnReference(t *testing.T) {
	f := test.New(t, CreateLspServices())
	doc := f.Parse(`
		statemachine Toggle

		events <|flick|>

		initialState a

		state a
		  <|flick>flick => b
		end

		state b
		  flick => a
		end
	`).AssertNoErrors()

	doc.AssertDefinition("flick")
}
