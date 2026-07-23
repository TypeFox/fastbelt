// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func TestChangeImpactCoversSiblingsInSameFolder(t *testing.T) {
	sc := CreateServices()
	f := test.New(t, sc)
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface Foo { Name string }
			Foo: Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar Test;
			interface Bar { Name string }
			Bar: Name=ID;
		`,
		"file:///other/c.fb", `
			grammar Other;
			interface Baz { Name string }
			Baz: Name=ID;
		`+commonTokens,
		"file:///other/d.fb", `
			grammar Other;
			Qux returns Baz: Name=ID;
		`,
	)
	a, b, c, d := docs[0].Document, docs[1].Document, docs[2].Document, docs[3].Document
	impact := service.MustGet[workspace.DocumentChangeImpact](sc)
	changed := collections.NewSet(a.URI.StringUnencoded())

	// b shares a's folder and is affected although it holds no reference into a.
	assert.True(t, impact.Affected(b, changed))
	// c is in another folder without references into a.
	assert.False(t, impact.Affected(c, changed))
	// d references c, so a change to c affects it through the default rule.
	assert.True(t, impact.Affected(d, collections.NewSet(c.URI.StringUnencoded())))
}
