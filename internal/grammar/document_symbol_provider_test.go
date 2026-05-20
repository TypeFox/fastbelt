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

func TestDocumentSymbolIntegration(t *testing.T) {
	sc := service.NewContainer()
	SetupServices(sc)
	server.SetupDefaultServices(sc)
	sc.Seal()

	fixture := test.New(t, sc)

	grammarText := `grammar Test;

<|grammar:interface First {
	name string
	value string
}|>

<|second:interface Second extends First {
	extra string
}|>

interface Third {
	<|third:data []string|>
}`

	doc := fixture.ParseURI(grammarText, "file:///test.fb")
	doc.AssertNoParseErrors().
		AssertDocumentSymbol("grammar", "First", lsp.Interface).
		AssertDocumentSymbol("second", "Second", lsp.Interface).
		AssertDocumentSymbol("third", "data", lsp.Property)
}
