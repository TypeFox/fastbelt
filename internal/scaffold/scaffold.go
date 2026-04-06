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
// writes scaffold files from embedded templates, attempts go get (library + tool), go generate, and go mod tidy.
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
	tryGoGetFastbeltDependencies(moduleRoot)
	if err := runGo(moduleRoot, "generate", "./..."); err != nil {
		return fmt.Errorf("go generate: %w", err)
	}
	if err := runGo(moduleRoot, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	return nil
}

// RunPackage creates packageRoot as a new empty directory inside an existing Go module, writes the
// same scaffold files as RunModule, attempts go get (library + tool), go generate for that package, and
// go mod tidy. It does not run go mod init; moduleRoot must contain go.mod.
func RunPackage(moduleRoot, packageRoot, packageImport, language string) error {
	names, err := prepareNames(packageImport, language)
	if err != nil {
		return err
	}
	if err := ensureScaffoldDir(packageRoot); err != nil {
		return err
	}
	if err := writeScaffoldFiles(packageRoot, names); err != nil {
		return err
	}
	tryGoGetFastbeltDependencies(moduleRoot)
	genArg, relErr := goGeneratePattern(moduleRoot, packageRoot)
	if relErr != nil {
		return relErr
	}
	if err := runGo(moduleRoot, "generate", genArg); err != nil {
		return fmt.Errorf("go generate: %w", err)
	}
	if err := runGo(moduleRoot, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	return nil
}

// ResolvePackageScaffoldDir finds the module root starting from workDir, reads its module path from
// go.mod, and returns the module root, the package directory, and the full package import path.
//
// packageSpec is normally a path relative to the module root (for example "examples/statemachine");
// the import path is modulePath + "/" + that path (using slash-separated segments). If packageSpec
// equals modulePath or starts with modulePath+"/", it is treated as a full import path instead.
func ResolvePackageScaffoldDir(workDir, packageSpec string) (moduleRoot, packageRoot, packageImport string, err error) {
	moduleRoot, err = findModuleRoot(workDir)
	if err != nil {
		return "", "", "", err
	}
	modPath, readErr := readGoModulePath(filepath.Join(moduleRoot, "go.mod"))
	if readErr != nil {
		return "", "", "", fmt.Errorf("read go.mod: %w", readErr)
	}
	packageRoot, packageImport, err = resolvePackageFromSpec(moduleRoot, modPath, packageSpec)
	if err != nil {
		return "", "", "", err
	}
	return moduleRoot, packageRoot, packageImport, nil
}

func resolvePackageFromSpec(moduleRoot, modulePath, spec string) (packageRoot, packageImport string, err error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return "", "", fmt.Errorf("package path is empty")
	}
	if spec == modulePath {
		return moduleRoot, modulePath, nil
	}
	if strings.HasPrefix(spec, modulePath+"/") {
		dir, impErr := packageDirForImport(moduleRoot, modulePath, spec)
		if impErr != nil {
			return "", "", impErr
		}
		return dir, spec, nil
	}
	if filepath.IsAbs(spec) {
		return "", "", fmt.Errorf("package path %q must be relative to the module root, not an absolute path", spec)
	}
	relSlash := filepath.ToSlash(spec)
	relSlash = strings.TrimPrefix(relSlash, "./")
	clean := path.Clean(relSlash)
	if clean == "." {
		return moduleRoot, modulePath, nil
	}
	if strings.HasPrefix(clean, "../") || clean == ".." {
		return "", "", fmt.Errorf("package path %q must not escape the module root", spec)
	}
	packageRoot = filepath.Join(moduleRoot, filepath.FromSlash(clean))
	absModRoot, absErr := filepath.Abs(moduleRoot)
	if absErr != nil {
		return "", "", absErr
	}
	absPkgRoot, absErr := filepath.Abs(packageRoot)
	if absErr != nil {
		return "", "", absErr
	}
	rel, relErr := filepath.Rel(absModRoot, absPkgRoot)
	if relErr != nil {
		return "", "", relErr
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("package path %q resolves outside the module root", spec)
	}
	return packageRoot, modulePath + "/" + clean, nil
}

func goGeneratePattern(moduleRoot, packageRoot string) (string, error) {
	rel, err := filepath.Rel(moduleRoot, packageRoot)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return "./...", nil
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("package directory %q is not inside module root %q", packageRoot, moduleRoot)
	}
	return "./" + filepath.ToSlash(rel) + "/...", nil
}

func findModuleRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		goMod := filepath.Join(dir, "go.mod")
		if _, statErr := os.Stat(goMod); statErr == nil {
			return dir, nil
		} else if !os.IsNotExist(statErr) {
			return "", statErr
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod found in %q or any parent directory; use an existing Go module when scaffolding a package", start)
		}
		dir = parent
	}
}

func readGoModulePath(goModPath string) (string, error) {
	b, readErr := os.ReadFile(goModPath)
	if readErr != nil {
		return "", readErr
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "module ") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			if i := strings.Index(rest, "//"); i >= 0 {
				rest = strings.TrimSpace(rest[:i])
			}
			rest = strings.Trim(rest, `"`)
			if rest == "" {
				return "", fmt.Errorf("empty module directive in %s", goModPath)
			}
			return rest, nil
		}
	}
	return "", fmt.Errorf("no module directive in %s", goModPath)
}

func packageDirForImport(moduleRoot, modulePath, packageImport string) (string, error) {
	packageImport = strings.TrimSpace(packageImport)
	if packageImport == "" {
		return "", fmt.Errorf("package import path is empty")
	}
	if packageImport == modulePath {
		return moduleRoot, nil
	}
	prefix := modulePath + "/"
	if !strings.HasPrefix(packageImport, prefix) {
		return "", fmt.Errorf("package import path %q must be equal to the module path %q or start with %q", packageImport, modulePath, prefix)
	}
	rel := strings.TrimPrefix(packageImport, prefix)
	if rel == "" {
		return "", fmt.Errorf("invalid package import path %q", packageImport)
	}
	return filepath.Join(moduleRoot, filepath.FromSlash(rel)), nil
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

// tryGoGetFastbeltDependencies runs go get for the fastbelt library and tool. Errors are printed
// to stderr and ignored so scaffolding can proceed when the module already provides those packages
// (for example developing fastbelt itself).
func tryGoGetFastbeltDependencies(moduleRoot string) {
	if err := runGo(moduleRoot, "get", "typefox.dev/fastbelt@latest"); err != nil {
		fmt.Fprintf(os.Stderr, "fastbelt scaffold: warning: go get typefox.dev/fastbelt@latest: %v\n", err)
	}
	if err := runGo(moduleRoot, "get", "-tool", "typefox.dev/fastbelt/cmd@latest"); err != nil {
		fmt.Fprintf(os.Stderr, "fastbelt scaffold: warning: go get -tool typefox.dev/fastbelt/cmd@latest: %v\n", err)
	}
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
