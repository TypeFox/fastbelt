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

	// Test 1: Empty query returns all symbols
	fixture.AssertWorkspaceSymbols("", []string{
		"Test1", "Test2",
		"Person", "Address", "Company",
		"name", "age", "street", "title",
		"PersonRule",
	})

	// Test 2: Query filtering returns matching subset
	fixture.AssertWorkspaceSymbols("Per", []string{"Person", "PersonRule"})

	// Test 3: Verify correct symbol kinds for different node types
	fixture.AssertWorkspaceSymbolKind("Test1", lsp.File)          // Grammar
	fixture.AssertWorkspaceSymbolKind("Person", lsp.Interface)    // Interface
	fixture.AssertWorkspaceSymbolKind("name", lsp.Property)       // Field
	fixture.AssertWorkspaceSymbolKind("PersonRule", lsp.Function) // Parser rule

	// Test 4: Case-insensitive matching
	fixture.AssertWorkspaceSymbols("person", []string{"Person", "PersonRule"})

	// Test 5: No match
	fixture.AssertWorkspaceSymbols("XYZ", []string{})
}
