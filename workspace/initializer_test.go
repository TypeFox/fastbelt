// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

func TestDefaultInitializerReadyClosesAfterInitialize(t *testing.T) {
	init := NewDefaultInitializer(service.NewContainer())
	select {
	case <-init.Ready():
		t.Fatal("Ready closed before Initialize")
	default:
	}
	if err := init.Initialize(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	select {
	case <-init.Ready():
	default:
		t.Fatal("Ready not closed after Initialize")
	}
}

// Every document loaded by the initializer must be backed by a file in the
// [textdoc.Store]: later document changes resolve their text document through
// the store and would otherwise not find the files loaded from disk.
func TestDefaultInitializerAddsLoadedFilesToTextStore(t *testing.T) {
	fs := NewVirtualFileSystem()
	loaded := core.ParseURI("file:///ws/a.hello")
	nested := core.ParseURI("file:///ws/sub/b.hello")
	ignored := core.ParseURI("file:///ws/c.other")
	require.NoError(t, fs.WriteFile(loaded, []byte("hello a")))
	require.NoError(t, fs.WriteFile(nested, []byte("hello b")))
	require.NoError(t, fs.WriteFile(ignored, []byte("other")))

	sc := service.NewContainer()
	textdoc.SetupDefaultServices(sc)
	sc.Put[FileSystem](fs)
	sc.Put[core.LanguageSelector](NewDefaultLanguageSelector(sc,
		NewDocumentMatcherWithPatterns("greeting", "**/*.hello"),
	))
	SetupDefaultServices(sc)
	sc.Seal()

	init := sc.MustGet[Initializer]()
	folders := []lsp.WorkspaceFolder{{URI: "file:///ws", Name: "ws"}}
	require.NoError(t, init.Initialize(context.Background(), folders))

	store := sc.MustGet[textdoc.Store]()
	documents := sc.MustGet[DocumentManager]()
	for uri, content := range map[core.URI]string{loaded: "hello a", nested: "hello b"} {
		file := store.GetFile(uri.DocumentURI())
		require.NotNil(t, file, "no file in the text store for %s", uri)
		assert.Equal(t, "greeting", file.LanguageID())
		assert.Equal(t, content, file.Text(nil))

		doc := documents.Get(uri)
		require.NotNil(t, doc, "no document for %s", uri)
		// The document and the store must share the same text document.
		assert.Same(t, file, doc.TextDoc)
		assert.Same(t, file, store.Get(uri.DocumentURI()))
	}

	assert.Nil(t, store.GetFile(ignored.DocumentURI()))
	assert.False(t, documents.Has(ignored))
	assert.Len(t, store.AllFiles(), 2)
}

func TestNoopInitializerReadyIsClosed(t *testing.T) {
	select {
	case <-NewNoopInitializer().Ready():
	default:
		t.Fatal("Ready not closed")
	}
}
