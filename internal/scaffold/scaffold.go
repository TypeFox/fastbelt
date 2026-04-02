// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:templates
var templateFS embed.FS

// RunModule creates moduleRoot as a new empty directory, runs go mod init for modulePath,
// writes scaffold files from embedded templates, runs go get (library + tool), go generate, and go mod tidy.
func RunModule(moduleRoot, modulePath, language string) error {
	names, err := prepareNames(modulePath, language)
	if err != nil {
		return err
	}
	if err := ensureScaffoldDir(moduleRoot); err != nil {
		return err
	}
	if err := runGo(moduleRoot, "mod", "init", names.ModulePath); err != nil {
		return fmt.Errorf("go mod init: %w", err)
	}
	if err := writeScaffoldFiles(moduleRoot, names); err != nil {
		return err
	}
	if err := runGo(moduleRoot, "get", "typefox.dev/fastbelt@latest"); err != nil {
		return fmt.Errorf("go get typefox.dev/fastbelt: %w", err)
	}
	if err := runGo(moduleRoot, "get", "-tool", "typefox.dev/fastbelt/cmd@latest"); err != nil {
		return fmt.Errorf("go get -tool typefox.dev/fastbelt/cmd: %w", err)
	}
	if err := runGo(moduleRoot, "generate", "./..."); err != nil {
		return fmt.Errorf("go generate: %w", err)
	}
	if err := runGo(moduleRoot, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	return nil
}

func versionSuffixFromGoGet(fastbeltGoGet string) string {
	_, v, ok := strings.Cut(fastbeltGoGet, "@")
	if !ok || v == "" {
		return "latest"
	}
	return v
}

func ensureScaffoldDir(dir string) error {
	_, statErr := os.Stat(dir)
	if statErr == nil {
		entries, readErr := os.ReadDir(dir)
		if readErr != nil {
			return readErr
		}
		if len(entries) > 0 {
			return fmt.Errorf("directory %s already exists and is not empty", dir)
		}
		return nil
	}
	if !os.IsNotExist(statErr) {
		return statErr
	}
	return os.MkdirAll(dir, 0755)
}

func runGo(dir string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func writeScaffoldFiles(moduleRoot string, names ModuleNames) error {
	type job struct {
		templateRel string
		outRel      string
	}
	jobs := []job{
		{"README.md.tmpl", "README.md"},
		{"gitignore.tmpl", ".gitignore"},
		{"package.root.json.tmpl", "package.json"},
		{"gen.go.tmpl", "gen.go"},
		{"services.go.tmpl", "services.go"},
		{"grammar.fb.tmpl", names.GrammarFile},
		{"cmd/main.go.tmpl", filepath.Join("cmd", names.LSPSlug, "main.go")},
		{"vscode-extension/package.json.tmpl", filepath.Join("vscode-extension", "package.json")},
		{"vscode-extension/src/extension.ts.tmpl", filepath.Join("vscode-extension", "src", "extension.ts")},
		{"vscode-extension/syntaxes/language.tmLanguage.json.tmpl", filepath.Join("vscode-extension", "syntaxes", names.SyntaxFile)},
		{"vscode-extension/vscodeignore.tmpl", filepath.Join("vscode-extension", ".vscodeignore")},
	}
	for _, j := range jobs {
		outPath := filepath.Join(moduleRoot, j.outRel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		buf, execErr := renderTemplate(j.templateRel, names)
		if execErr != nil {
			return fmt.Errorf("template %s: %w", j.templateRel, execErr)
		}
		if writeErr := os.WriteFile(outPath, buf, 0644); writeErr != nil {
			return writeErr
		}
	}
	return copyStaticScaffoldFiles(moduleRoot)
}

func renderTemplate(rel string, names ModuleNames) ([]byte, error) {
	b, err := templateFS.ReadFile(path.Join("templates", filepath.ToSlash(rel)))
	if err != nil {
		return nil, err
	}
	t, err := template.New(rel).Parse(string(b))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, names); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func copyStaticScaffoldFiles(moduleRoot string) error {
	static := []string{
		"vscode-extension/esbuild.js",
		"vscode-extension/tsconfig.json",
		"vscode-extension/language-configuration.json",
	}
	for _, rel := range static {
		body, err := templateFS.ReadFile(path.Join("templates", rel))
		if err != nil {
			return fmt.Errorf("read static %s: %w", rel, err)
		}
		outPath := filepath.Join(moduleRoot, rel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, body, 0644); err != nil {
			return err
		}
	}
	return nil
}

// WriteScaffoldFilesOnly writes templated and static scaffold files into moduleRoot without running go commands.
// It is used by tests and does not create a go.mod file.
func WriteScaffoldFilesOnly(moduleRoot string, names ModuleNames) error {
	if err := os.MkdirAll(moduleRoot, 0755); err != nil {
		return err
	}
	return writeScaffoldFiles(moduleRoot, names)
}
