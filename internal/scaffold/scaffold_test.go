// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package scaffold

import (
	"bytes"
	"io/fs"
	"os"
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

func TestReadGoModulePath(t *testing.T) {
	dir := t.TempDir()
	goMod := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(goMod, []byte("module example.com/foo/bar\n\ngo 1.23\n"), 0644))
	p, err := readGoModulePath(goMod)
	require.NoError(t, err)
	require.Equal(t, "example.com/foo/bar", p)
}

func TestReadGoModulePath_quoted(t *testing.T) {
	dir := t.TempDir()
	goMod := filepath.Join(dir, "go.mod")
	require.NoError(t, os.WriteFile(goMod, []byte("module \"example.com/q\"\n"), 0644))
	p, err := readGoModulePath(goMod)
	require.NoError(t, err)
	require.Equal(t, "example.com/q", p)
}

func TestFindModuleRoot(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	require.NoError(t, os.MkdirAll(nested, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module m\n"), 0644))
	found, err := findModuleRoot(nested)
	require.NoError(t, err)
	require.Equal(t, root, found)
}

func TestFindModuleRoot_missing(t *testing.T) {
	dir := t.TempDir()
	_, err := findModuleRoot(dir)
	require.Error(t, err)
}

func TestPackageDirForImport(t *testing.T) {
	root := t.TempDir()
	d, err := packageDirForImport(root, "example.com/proj", "example.com/proj/pkg/mylang")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, "pkg", "mylang"), d)

	d2, err := packageDirForImport(root, "example.com/proj", "example.com/proj")
	require.NoError(t, err)
	require.Equal(t, root, d2)
}

func TestPackageDirForImport_rejectsForeign(t *testing.T) {
	_, err := packageDirForImport(t.TempDir(), "example.com/proj", "other.com/x")
	require.Error(t, err)
}

func TestGoGeneratePattern(t *testing.T) {
	mod := t.TempDir()
	pkg := filepath.Join(mod, "pkg", "lang")
	require.NoError(t, os.MkdirAll(pkg, 0755))
	pat, err := goGeneratePattern(mod, pkg)
	require.NoError(t, err)
	require.Equal(t, "./pkg/lang/...", pat)

	patRoot, err := goGeneratePattern(mod, mod)
	require.NoError(t, err)
	require.Equal(t, "./...", patRoot)
}

func TestResolvePackageScaffoldDir_fullImportPath(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "inner")
	require.NoError(t, os.MkdirAll(sub, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/demo\n"), 0644))
	modRoot, pkgRoot, pkgImport, err := ResolvePackageScaffoldDir(sub, "example.com/demo/lang")
	require.NoError(t, err)
	require.Equal(t, root, modRoot)
	require.Equal(t, filepath.Join(root, "lang"), pkgRoot)
	require.Equal(t, "example.com/demo/lang", pkgImport)
}

func TestResolvePackageScaffoldDir_relativeToModuleRoot(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module typefox.dev/fastbelt\n"), 0644))
	modRoot, pkgRoot, pkgImport, err := ResolvePackageScaffoldDir(root, "examples/statemachine")
	require.NoError(t, err)
	require.Equal(t, root, modRoot)
	require.Equal(t, filepath.Join(root, "examples", "statemachine"), pkgRoot)
	require.Equal(t, "typefox.dev/fastbelt/examples/statemachine", pkgImport)
}

func TestResolvePackageScaffoldDir_rejectsAbsolutePackagePath(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module m\n"), 0644))
	_, _, _, err := ResolvePackageScaffoldDir(root, "/tmp/nope")
	require.Error(t, err)
}

func TestResolvePackageScaffoldDir_rejectsParentEscape(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module m\n"), 0644))
	_, _, _, err := ResolvePackageScaffoldDir(root, "../outside")
	require.Error(t, err)
}

func mustNames(t *testing.T, modulePath, lang string) ModuleNames {
	t.Helper()
	n, err := prepareNames(modulePath, lang)
	require.NoError(t, err)
	return n
}
