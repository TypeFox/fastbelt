// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"context"
	"path/filepath"
	"slices"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// initialize runs the registered Initializer over the given folders and
// returns the sorted URIs of all documents it loaded.
func initialize(t *testing.T, sc *service.Container, folders ...core.URI) []string {
	t.Helper()
	init := sc.MustGet[workspace.Initializer]()
	var wsFolders []lsp.WorkspaceFolder
	for _, f := range folders {
		wsFolders = append(wsFolders, lsp.WorkspaceFolder{URI: lsp.URI(f.String())})
	}
	require.NoError(t, init.Initialize(context.Background(), wsFolders))
	select {
	case <-init.Ready():
	default:
		t.Fatal("Ready not closed after Initialize")
	}
	var uris []string
	for doc := range sc.MustGet[workspace.DocumentManager]().All() {
		uris = append(uris, doc.URI.StringUnencoded())
		assert.Equal(t, "arithmetics", doc.TextDoc.LanguageID(), doc.URI.String())
	}
	sort.Strings(uris)
	return uris
}

func TestInitializerLoadsExamplesFolder(t *testing.T) {
	dir, err := filepath.Abs("examples")
	require.NoError(t, err)
	root := core.FileURI(dir)

	uris := initialize(t, CreateServices(), root)

	assert.Equal(t, []string{
		root.JoinPath("example.calc").StringUnencoded(),
		root.JoinPath("price-calc.calc").StringUnencoded(),
		root.JoinPath("recursion.calc").StringUnencoded(),
	}, uris, "only .calc files are loaded, price-calc.calc.json is skipped")
}

func TestInitializerWalksVirtualFileSystem(t *testing.T) {
	vfs := workspace.NewVirtualFileSystem()
	write := func(path string) {
		require.NoError(t, vfs.WriteFile(core.FileURI(path), []byte("module m")))
	}
	write("/ws/a.calc")
	write("/ws/nested/deep/b.calc")
	write("/ws/notes.txt")          // extension not mapped to a language
	write("/ws/.hidden/c.calc")     // hidden directories are skipped
	write("/ws/nested/.git/d.calc") // ...also when nested
	write("/other/e.calc")          // outside the workspace folder

	sc := service.NewContainer()
	sc.Put[workspace.FileSystem](vfs)
	SetupServices(sc)
	sc.Seal()

	// A document registered before initialization must not be replaced.
	docs := sc.MustGet[workspace.DocumentManager]()
	preloaded := core.NewDocumentFromString("file:///ws/a.calc", "arithmetics", "module preloaded")
	docs.Set(preloaded)

	uris := initialize(t, sc, core.FileURI("/ws"))

	assert.Equal(t, []string{
		"file:///ws/a.calc",
		"file:///ws/nested/deep/b.calc",
	}, uris)
	assert.Same(t, preloaded, docs.Get(core.FileURI("/ws/a.calc")))
	assert.Equal(t, "module m", docs.Get(core.FileURI("/ws/nested/deep/b.calc")).TextDoc.Text(nil))
}

func TestInitializerContinuesAfterMissingFolder(t *testing.T) {
	vfs := workspace.NewVirtualFileSystem()
	require.NoError(t, vfs.WriteFile(core.FileURI("/ws/a.calc"), []byte("module m")))

	sc := service.NewContainer()
	sc.Put[workspace.FileSystem](vfs)
	SetupServices(sc)
	sc.Seal()

	// The missing folder is logged and skipped; the next folder is still walked.
	uris := initialize(t, sc, core.FileURI("/missing"), core.FileURI("/ws"))
	assert.True(t, slices.Equal([]string{"file:///ws/a.calc"}, uris), uris)
}
