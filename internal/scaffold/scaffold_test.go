// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package scaffold

import (
	"bytes"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestPrepareNames(t *testing.T) {
	n, err := prepareNames("example.com/acme/foo", "My Lang")
	require.NoError(t, err)
	require.Equal(t, "example.com/acme/foo", n.ModulePath)
	require.Equal(t, "foo", n.ModuleDirBase)
	require.Equal(t, "my-lang", n.LanguageID)
	require.Equal(t, "MyLang", n.ExportedBase)
	require.Equal(t, "MyLangModel", n.GrammarName)
	require.Equal(t, "mylang", n.GoPackage)
	require.Equal(t, "my-lang.fb", n.GrammarFile)
	require.Equal(t, ".mylang", n.FileDotExt)
	require.Equal(t, "my-lang-lsp", n.LSPSlug)
	require.Contains(t, n.FastbeltGoGet, "typefox.dev/fastbelt@")
}

func TestPrepareNames_goKeywordPackage(t *testing.T) {
	n, err := prepareNames("x/y", "break")
	require.NoError(t, err)
	require.Equal(t, "breaklang", n.GoPackage)
}

func TestWriteScaffoldFilesOnly(t *testing.T) {
	dir := t.TempDir()
	n, err := prepareNames("example.com/demo/mylang", "Demo")
	require.NoError(t, err)
	require.NoError(t, WriteScaffoldFilesOnly(dir, n))

	require.FileExists(t, filepath.Join(dir, "README.md"))
	require.FileExists(t, filepath.Join(dir, n.GrammarFile))
	require.FileExists(t, filepath.Join(dir, "gen.go"))
	require.FileExists(t, filepath.Join(dir, "services.go"))
	require.FileExists(t, filepath.Join(dir, ".gitignore"))
	require.FileExists(t, filepath.Join(dir, "package.json"))
	require.FileExists(t, filepath.Join(dir, "cmd", n.LSPSlug, "main.go"))
	require.FileExists(t, filepath.Join(dir, "vscode-extension", "package.json"))
	require.FileExists(t, filepath.Join(dir, "vscode-extension", "esbuild.js"))
	require.FileExists(t, filepath.Join(dir, "vscode-extension", "src", "extension.ts"))
	require.FileExists(t, filepath.Join(dir, "vscode-extension", "syntaxes", n.SyntaxFile))
}

func TestEmbeddedTemplatesExecute(t *testing.T) {
	n, err := prepareNames("example.com/a/b", "Zed")
	require.NoError(t, err)
	require.NoError(t, fs.WalkDir(templateFS, "templates", func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".tmpl") {
			return nil
		}
		b, readErr := templateFS.ReadFile(p)
		require.NoError(t, readErr)
		tmpl, parseErr := template.New(p).Parse(string(b))
		require.NoError(t, parseErr, p)
		var buf bytes.Buffer
		require.NoError(t, tmpl.Execute(&buf, n), p)
		require.NotEmpty(t, buf.String(), p)
		return nil
	}))
}

func TestRunModule_rejectsNonemptyDir(t *testing.T) {
	nonEmpty := t.TempDir()
	require.NoError(t, WriteScaffoldFilesOnly(nonEmpty, mustNames(t, "x/y", "L")))
	runErr := RunModule(nonEmpty, "example.com/neo/mod", "Neo")
	require.Error(t, runErr)
}

func mustNames(t *testing.T, modulePath, lang string) ModuleNames {
	t.Helper()
	n, err := prepareNames(modulePath, lang)
	require.NoError(t, err)
	return n
}
