// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
)

func newSelectorContainer(t *testing.T) *service.Container {
	t.Helper()
	sc := service.NewContainer()
	textdoc.SetupDefaultServices(sc)
	sc.Seal()
	return sc
}

func TestDefaultLanguageSelectorMatchesByGlob(t *testing.T) {
	sc := newSelectorContainer(t)
	selector := NewDefaultLanguageSelector(sc,
		NewDocumentSelectorWithPatterns("greeting", "**/*.hello"),
		NewDocumentSelectorWithPatterns("farewell", "**/*.bye"),
	)

	i, id := selector.Select(ParseURI("file:///ws/a.hello"))
	assert.Equal(t, 0, i)
	assert.Equal(t, "greeting", id)

	i, id = selector.Select(ParseURI("file:///ws/b.bye"))
	assert.Equal(t, 1, i)
	assert.Equal(t, "farewell", id)

	i, id = selector.Select(ParseURI("file:///ws/c.other"))
	assert.Equal(t, -1, i)
	assert.Equal(t, "", id)
}

// A client-supplied language id (from an open document) must win over the
// path glob of an earlier selector; globs are only consulted when no selector
// claims the id.
func TestDefaultLanguageSelectorPrefersStoredLanguageIDOverEarlierGlob(t *testing.T) {
	sc := newSelectorContainer(t)
	selector := NewDefaultLanguageSelector(sc,
		NewDocumentSelectorWithPatterns("a", "**/*.txt"),
		NewDocumentSelectorWithPatterns("b", "**/*.txt"),
	)
	uri := ParseURI("file:///ws/x.txt")

	// Not open yet: the first glob wins.
	i, _ := selector.Select(uri)
	assert.Equal(t, 0, i)

	store := service.MustGet[textdoc.Store](sc)
	store.AddOverlay(textdoc.NewOverlay(uri.DocumentURI(), "b", 1, ""))
	i, id := selector.Select(uri)
	assert.Equal(t, 1, i)
	assert.Equal(t, "b", id)

	// An id no selector claims falls back to the globs.
	store.AddOverlay(textdoc.NewOverlay(uri.DocumentURI(), "plaintext", 2, ""))
	i, _ = selector.Select(uri)
	assert.Equal(t, 0, i)
}
