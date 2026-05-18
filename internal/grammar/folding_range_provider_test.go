// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
)

func TestFoldingRangeIntegration(t *testing.T) {
	sc := service.NewContainer()
	SetupServices(sc)
	server.SetupDefaultServices(sc)
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
	doc.AssertNoParseErrors().AssertFoldingRanges("first", "second", "comment", "third")
}
