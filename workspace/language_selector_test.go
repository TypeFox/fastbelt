// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
)

func newSelectorContainer() *service.Container {
	sc := service.NewContainer()
	textdoc.SetupDefaultServices(sc)
	service.Put(sc, NewDefaultDocumentManager(sc))
	sc.Seal()
	return sc
}

func TestDefaultLanguageSelectorMatchesByGlob(t *testing.T) {
	selector := NewDefaultLanguageSelector(newSelectorContainer(),
		core.NewDocumentSelectorWithPatterns("greeting", "**/*.hello"),
		core.NewDocumentSelectorWithPatterns("farewell", "**/*.bye"),
	)

	i, id := selector.Select(core.ParseURI("file:///ws/a.hello"))
	assert.Equal(t, 0, i)
	assert.Equal(t, "greeting", id)

	i, id = selector.Select(core.ParseURI("file:///ws/b.bye"))
	assert.Equal(t, 1, i)
	assert.Equal(t, "farewell", id)

	i, id = selector.Select(core.ParseURI("file:///ws/c.other"))
	assert.Equal(t, -1, i)
	assert.Equal(t, "", id)
}

func TestDefaultLanguageSelectorMatchesBySimpleGlob(t *testing.T) {
	selector := NewDefaultLanguageSelector(newSelectorContainer(),
		core.NewDocumentSelectorWithPatterns("greeting", "*.hello"),
		core.NewDocumentSelectorWithPatterns("farewell", "*.bye"),
	)

	i, id := selector.Select(core.ParseURI("file:///ws/a.hello"))
	assert.Equal(t, 0, i)
	assert.Equal(t, "greeting", id)

	i, id = selector.Select(core.ParseURI("file:///ws/b.bye"))
	assert.Equal(t, 1, i)
	assert.Equal(t, "farewell", id)

	i, id = selector.Select(core.ParseURI("file:///ws/c.other"))
	assert.Equal(t, -1, i)
	assert.Equal(t, "", id)
}

// The language id of the managed document must win over the path glob of an
// earlier selector; globs are only consulted when no selector claims the id.
func TestDefaultLanguageSelectorPrefersManagedDocumentLanguageID(t *testing.T) {
	sc := newSelectorContainer()
	selector := NewDefaultLanguageSelector(sc,
		core.NewDocumentSelectorWithPatterns("a", "**/*.txt"),
		core.NewDocumentSelectorWithPatterns("b", "**/*.txt"),
	)
	uri := core.ParseURI("file:///ws/x.txt")
	documents := service.MustGet[DocumentManager](sc)

	// Not managed yet: the first glob wins.
	i, _ := selector.Select(uri)
	assert.Equal(t, 0, i)

	documents.Set(core.NewDocument(textdoc.NewOverlay(uri.DocumentURI(), "b", 1, "")))
	i, id := selector.Select(uri)
	assert.Equal(t, 1, i)
	assert.Equal(t, "b", id)

	// A newer handle in the store does not change the language until the
	// document is rebuilt from it.
	service.MustGet[textdoc.Store](sc).AddOverlay(textdoc.NewOverlay(uri.DocumentURI(), "a", 2, ""))
	i, _ = selector.Select(uri)
	assert.Equal(t, 1, i)

	// An id no selector claims falls back to the globs.
	documents.Set(core.NewDocument(textdoc.NewOverlay(uri.DocumentURI(), "plaintext", 3, "")))
	i, _ = selector.Select(uri)
	assert.Equal(t, 0, i)
}
