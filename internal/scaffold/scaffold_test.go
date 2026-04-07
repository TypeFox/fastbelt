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

func TestFastbeltModuleVersion_envOverride(t *testing.T) {
	t.Setenv("FASTBELT_SCAFFOLD_FASTBELT_GO_VERSION", "latest")
	require.Equal(t, "latest", fastbeltModuleVersion())
}

func TestPrepareNames(t *testing.T) {
	n, err := newTemplateParams(&Scaffolder{
		ImportPath: "example.com/acme/foo",
		Language:   "My Lang",
	})
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
	require.Contains(t, n.FastbeltToolGoGet, "typefox.dev/fastbelt/cmd/fastbelt@")
}

func TestPrepareNames_goKeywordPackage(t *testing.T) {
	n, err := newTemplateParams(&Scaffolder{
		ImportPath: "x/y",
		Language:   "break",
	})
	require.NoError(t, err)
	require.Equal(t, "breaklang", n.GoPackage)
}

func TestEmbeddedTemplatesExecute(t *testing.T) {
	n, err := newTemplateParams(&Scaffolder{
		ImportPath: "example.com/a/b",
		Language:   "Zed",
	})
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

func TestGoGeneratePattern(t *testing.T) {
	mod := t.TempDir()
	pkg := filepath.Join(mod, "pkg", "lang")
	require.NoError(t, os.MkdirAll(pkg, 0755))
	modAbs, err := filepath.Abs(mod)
	require.NoError(t, err)
	pkgAbs, err := filepath.Abs(pkg)
	require.NoError(t, err)
	pat, err := goGeneratePattern(modAbs, pkgAbs)
	require.NoError(t, err)
	require.Equal(t, "./pkg/lang/...", pat)

	patRoot, err := goGeneratePattern(modAbs, modAbs)
	require.NoError(t, err)
	require.Equal(t, "./...", patRoot)
}

func TestResolveScaffoldFromWorkDir_nestedCwdUsesDotDot(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "inner")
	require.NoError(t, os.MkdirAll(sub, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/demo\n"), 0644))
	scaffolder := &Scaffolder{}
	err := scaffolder.PopulateDirectorysFromWorkDir(sub, "../lang")
	require.NoError(t, err)
	require.Equal(t, root, scaffolder.ModuleRoot)
	require.Equal(t, filepath.Join(root, "lang"), scaffolder.WriteRoot)
	require.Equal(t, "example.com/demo/lang", scaffolder.ImportPath)
}

func TestResolveWriteRootUnderModule_rejectsFullImportPath(t *testing.T) {
	modDir := filepath.Join(t.TempDir(), "m")
	scaffolder := &Scaffolder{}
	err := scaffolder.PopulateDirectoriesFromModuleDir(modDir, "example.com/demo", "example.com/demo/lang")
	require.Error(t, err)
	require.Contains(t, err.Error(), `"lang"`)
}

func TestResolveScaffoldFromWorkDir_relativeToWorkingDir(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module typefox.dev/fastbelt\n"), 0644))
	scaffolder := &Scaffolder{}
	err := scaffolder.PopulateDirectorysFromWorkDir(root, "examples/statemachine")
	require.NoError(t, err)
	require.Equal(t, root, scaffolder.ModuleRoot)
	require.Equal(t, filepath.Join(root, "examples", "statemachine"), scaffolder.WriteRoot)
	require.Equal(t, "typefox.dev/fastbelt/examples/statemachine", scaffolder.ImportPath)
}

func TestResolveScaffoldFromWorkDir_rejectsAbsolutePackagePath(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module m\n"), 0644))
	absPkg, absErr := filepath.Abs(filepath.Join(t.TempDir(), "nope"))
	require.NoError(t, absErr)
	require.True(t, filepath.IsAbs(absPkg), "test requires an absolute path on this OS")
	scaffolder := &Scaffolder{}
	err := scaffolder.PopulateDirectorysFromWorkDir(root, absPkg)
	require.Error(t, err)
}

func TestResolveScaffoldFromWorkDir_rejectsOutsideModule(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module m\n"), 0644))
	scaffolder := &Scaffolder{}
	err := scaffolder.PopulateDirectorysFromWorkDir(root, "../outside")
	require.Error(t, err)
}

func TestResolveWriteRootUnderModule_subdirOfNewModule(t *testing.T) {
	modDir := filepath.Join(t.TempDir(), "foo")
	scaffolder := &Scaffolder{}
	err := scaffolder.PopulateDirectoriesFromModuleDir(modDir, "example.com/acme/foo", "pkg/lang")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(modDir, "pkg", "lang"), scaffolder.WriteRoot)
	require.Equal(t, "example.com/acme/foo/pkg/lang", scaffolder.ImportPath)
}
