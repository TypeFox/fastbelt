// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

func TestWorkspaceSymbolIntegration(t *testing.T) {
	sc := service.NewContainer()
	SetupServices(sc)
	server.SetupDefaultServices(sc)
	sc.Seal()

	fixture := test.New(t, sc)

	// Create multiple documents with various symbol types
	grammar1 := `grammar Test1;
interface Person {
	name string
	age int
}

interface Address {
	street string
}
`

	grammar2 := `grammar Test2;
interface Company {
	title string
}

PersonRule: /[a-z]+/;
`

	fixture.ParseURI(grammar1, "file:///test1.fb")
	fixture.ParseURI(grammar2, "file:///test2.fb")

	t.Run("Empty query returns all symbols", func(t *testing.T) {
		fixture.ExpectWorkspaceSymbols("").ExactMatch(
			"Test1", "Test2",
			"Person", "Address", "Company",
			"name", "age", "street", "title",
			"PersonRule",
		)
	})

	t.Run("Query filtering returns matching subset", func(t *testing.T) {
		fixture.ExpectWorkspaceSymbols("Per").ExactMatch("Person", "PersonRule")
	})

	t.Run("Verify correct symbol kinds for different node types", func(t *testing.T) {
		fixture.ExpectWorkspaceSymbols("Test1").SymbolKind("Test1", lsp.File) // Grammar
		fixture.ExpectWorkspaceSymbols("Person").
			SymbolKind("Person", lsp.Interface).   // Interface
			SymbolKind("PersonRule", lsp.Function) // Parser rule (also matches "Person")
		fixture.ExpectWorkspaceSymbols("name").SymbolKind("name", lsp.Property) // Field
	})

	t.Run("Case-insensitive matching", func(t *testing.T) {
		fixture.ExpectWorkspaceSymbols("person").ExactMatch("Person", "PersonRule")
	})

	t.Run("No match", func(t *testing.T) {
		fixture.ExpectWorkspaceSymbols("XYZ").ExactMatch()
	})
}
