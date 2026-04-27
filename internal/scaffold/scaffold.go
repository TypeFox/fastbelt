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

// Scaffolder runs the language scaffold workflow (templates, optional go mod init, go get attempts,
// go generate, go mod tidy).
type Scaffolder struct {
	// CreateModule, when true, ensures ModuleRoot is empty or new, runs go mod init there. WriteRoot
	// must be ModuleRoot or a subdirectory of it (templates are written to WriteRoot). When false,
	// WriteRoot is the target package directory under the existing module at ModuleRoot.
	CreateModule bool
	// CreateVSCodeExtension, when true, creates the VS Code extension directory and files.
	CreateVSCodeExtension bool
	// ModuleRoot is the directory containing go.mod (or that will after init). All go commands run here.
	ModuleRoot string
	// WriteRoot is where embedded templates are written.
	WriteRoot string
	// ImportPath is the module path (new module) or full package import path (package mode) for naming.
	ImportPath string
	// Language is the human-readable language label for templates and derived identifiers.
	Language string
}

// Run executes the scaffold: prepare names, optionally go mod init, write files, try fastbelt go get,
// go generate, and go mod tidy.
func (s *Scaffolder) Run() error {
	if s.ModuleRoot == "" {
		return fmt.Errorf("scaffold: ModuleRoot is empty")
	}
	modRoot, absErr := filepath.Abs(filepath.Clean(s.ModuleRoot))
	if absErr != nil {
		return fmt.Errorf("scaffold: ModuleRoot: %w", absErr)
	}
	s.ModuleRoot = modRoot
	if s.WriteRoot != "" {
		writeRoot, writeAbsErr := filepath.Abs(filepath.Clean(s.WriteRoot))
		if writeAbsErr != nil {
			return fmt.Errorf("scaffold: WriteRoot: %w", writeAbsErr)
		}
		s.WriteRoot = writeRoot
	}
	names, err := newTemplateParams(s)
	if err != nil {
		return err
	}
	if s.CreateModule {
		relToMod, relErr := filepath.Rel(s.ModuleRoot, s.WriteRoot)
		if relErr != nil || (relToMod != "." && !filepath.IsLocal(relToMod)) {
			return fmt.Errorf("scaffold: WriteRoot must be inside ModuleRoot")
		}
		if err := ensureScaffoldDir(s.ModuleRoot); err != nil {
			return err
		}
		if err := runGo(s.ModuleRoot, "mod", "init", names.ModulePath); err != nil {
			return fmt.Errorf("go mod init: %w", err)
		}
		if s.WriteRoot != s.ModuleRoot {
			if err := ensureScaffoldDir(s.WriteRoot); err != nil {
				return err
			}
		}
	} else {
		if s.WriteRoot == "" {
			return fmt.Errorf("scaffold: WriteRoot is required when CreateModule is false")
		}
		if err := ensureScaffoldDir(s.WriteRoot); err != nil {
			return err
		}
	}
	if err := s.scaffoldTemplatedFiles(names); err != nil {
		return err
	}
	if localPath := fastbeltModuleLocalPath(); localPath != "" {
		// When running the scaffold in the dev environment, the module version should point to the local copy
		// This requires to edit the go mod file to point to the local path
		if err := runGo(s.ModuleRoot, "mod", "edit", "-replace", fastbeltModulePath+"="+localPath); err != nil {
			fmt.Fprintf(os.Stderr, "fastbelt scaffold: warning: go mod edit -replace: %v\n", err)
		}
	}
	tryGoGetFastbeltDependencies(s.ModuleRoot, names.FastbeltGoGet, names.FastbeltToolGoGet)

	genArg := "./..."
	if s.WriteRoot != s.ModuleRoot {
		var patternErr error
		genArg, patternErr = goGeneratePattern(s.ModuleRoot, s.WriteRoot)
		if patternErr != nil {
			return patternErr
		}
	}
	if err := runGo(s.ModuleRoot, "generate", genArg); err != nil {
		return fmt.Errorf("go generate: %w", err)
	}
	if err := runGo(s.ModuleRoot, "mod", "tidy"); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	return nil
}

// PopulateDirectorysFromWorkDir resolves the Go module that contains the directory
// filepath.Join(workDir, packageRel). packageRel is relative to workDir (for example "." or
// "examples/mylang"). It populates s.ModuleRoot and s.WriteRoot and s.ImportPath.
func (s *Scaffolder) PopulateDirectorysFromWorkDir(workDir, packageRel string) error {
	packageRel = strings.TrimSpace(packageRel)
	if packageRel == "" {
		packageRel = "."
	}
	if filepath.IsAbs(packageRel) {
		return fmt.Errorf("package path %q must be relative to the working directory, not an absolute path", packageRel)
	}
	absWd, err := filepath.Abs(workDir)
	if err != nil {
		return err
	}
	joined := filepath.Join(absWd, filepath.FromSlash(path.Clean(filepath.ToSlash(packageRel))))
	s.WriteRoot, err = filepath.Abs(joined)
	if err != nil {
		return err
	}
	s.ModuleRoot, err = findModuleRoot(s.WriteRoot)
	if err != nil {
		return err
	}
	modPath, readErr := readGoModulePath(filepath.Join(s.ModuleRoot, "go.mod"))
	if readErr != nil {
		return fmt.Errorf("read go.mod: %w", readErr)
	}
	rel, relErr := filepath.Rel(s.ModuleRoot, s.WriteRoot)
	if relErr != nil {
		return relErr
	}
	if rel != "." && !filepath.IsLocal(rel) {
		return fmt.Errorf("package directory %q is outside the module root %q", s.WriteRoot, s.ModuleRoot)
	}
	if rel == "." {
		s.ImportPath = modPath
	} else {
		s.ImportPath = modPath + "/" + filepath.ToSlash(rel)
	}
	return nil
}

// PopulateDirectoriesFromModuleDir maps packageSpec to a directory under moduleRootDir and the full
// package import path. packageSpec is relative to moduleRootDir (for example "." or "pkg/lang").
// moduleImportPath is the Go module path (as in go.mod). moduleRootDir need not exist yet.
// It populates s.ModuleRoot and s.WriteRoot and s.ImportPath.
func (s *Scaffolder) PopulateDirectoriesFromModuleDir(moduleRootDir, moduleImportPath, packageSpec string) error {
	s.ModuleRoot = moduleRootDir
	packageSpec = strings.TrimSpace(packageSpec)
	if packageSpec == "" {
		return fmt.Errorf("package path is empty")
	}
	if packageSpec == moduleImportPath {
		return fmt.Errorf("package path must be relative to the module root (use \".\" for the module root), not the module path %q", moduleImportPath)
	}
	if strings.HasPrefix(packageSpec, moduleImportPath+"/") {
		rel := strings.TrimPrefix(packageSpec, moduleImportPath+"/")
		return fmt.Errorf("package path must be relative to the module root, not a full import path (use %q instead of %q)", rel, packageSpec)
	}
	if filepath.IsAbs(packageSpec) {
		return fmt.Errorf("package path %q must be relative to the module root, not an absolute path", packageSpec)
	}
	relSlash := filepath.ToSlash(packageSpec)
	relSlash = strings.TrimPrefix(relSlash, "./")
	clean := path.Clean(relSlash)
	if clean == "." {
		s.WriteRoot = moduleRootDir
		s.ImportPath = moduleImportPath
		return nil
	}
	if strings.HasPrefix(clean, "../") || clean == ".." {
		return fmt.Errorf("package path %q must not escape the module root", packageSpec)
	}
	packageRoot := filepath.Join(moduleRootDir, filepath.FromSlash(clean))
	absModRoot, err := filepath.Abs(moduleRootDir)
	if err != nil {
		return err
	}
	absPkgRoot, err := filepath.Abs(packageRoot)
	if err != nil {
		return err
	}
	rel, relErr := filepath.Rel(absModRoot, absPkgRoot)
	if relErr != nil {
		return relErr
	}
	if rel != "." && !filepath.IsLocal(rel) {
		return fmt.Errorf("package path %q resolves outside the module root", packageSpec)
	}
	s.WriteRoot = packageRoot
	s.ImportPath = moduleImportPath + "/" + clean
	return nil
}

func goGeneratePattern(moduleRoot, packageRoot string) (string, error) {
	rel, err := filepath.Rel(moduleRoot, packageRoot)
	if err != nil {
		return "", err
	}
	if rel == "." {
		return "./...", nil
	}
	if !filepath.IsLocal(rel) {
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

// scaffoldOutputDir is the directory where generated template files are written.
func (s *Scaffolder) scaffoldOutputDir() string {
	if s.WriteRoot != "" {
		return s.WriteRoot
	}
	return s.ModuleRoot
}

// tryGoGetFastbeltDependencies runs go get for the fastbelt library and tool at librarySpec and toolSpec
// (for example from ModuleNames.FastbeltGoGet / FastbeltToolGoGet). Errors are printed to stderr and
// ignored so scaffolding can proceed when the module already provides those packages (for example
// developing fastbelt itself).
func tryGoGetFastbeltDependencies(moduleRoot, librarySpec, toolSpec string) {
	if err := runGo(moduleRoot, "get", librarySpec); err != nil {
		fmt.Fprintf(os.Stderr, "fastbelt scaffold: warning: go get %s: %v\n", librarySpec, err)
	}
	if err := runGo(moduleRoot, "get", "-tool", toolSpec); err != nil {
		fmt.Fprintf(os.Stderr, "fastbelt scaffold: warning: go get -tool %s: %v\n", toolSpec, err)
	}
}

func (s *Scaffolder) scaffoldTemplatedFiles(params templateParams) error {
	type job struct {
		templateRel string
		outRel      string
	}
	jobs := []job{
		{"README.md.tmpl", "README.md"},
		{"gitignore.tmpl", ".gitignore"},
		{"gen.go.tmpl", "gen.go"},
		{"services.go.tmpl", "services.go"},
		{"grammar.fb.tmpl", params.GrammarFile},
		{"cmd/main.go.tmpl", filepath.Join("cmd", params.LSPSlug, "main.go")},
	}
	if s.CreateVSCodeExtension {
		jobs = append(jobs, []job{
			{"package.root.json.tmpl", "package.json"},
			{"vscode-extension/package.json.tmpl", filepath.Join("vscode-extension", "package.json")},
			{"vscode-extension/tsconfig.json.tmpl", filepath.Join("vscode-extension", "tsconfig.json")},
			{"vscode-extension/src/extension.ts.tmpl", filepath.Join("vscode-extension", "src", "extension.ts")},
			{"vscode-extension/syntaxes/language.tmLanguage.json.tmpl", filepath.Join("vscode-extension", "syntaxes", params.SyntaxFile)},
			{"vscode-extension/vscodeignore.tmpl", filepath.Join("vscode-extension", ".vscodeignore")},
		}...)
	}
	for _, j := range jobs {
		outPath := filepath.Join(s.scaffoldOutputDir(), j.outRel)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		buf, execErr := s.renderTemplate(j.templateRel, params)
		if execErr != nil {
			return fmt.Errorf("template %s: %w", j.templateRel, execErr)
		}
		if writeErr := s.writeFile(j.outRel, buf); writeErr != nil {
			return writeErr
		}
	}
	return s.copyStaticScaffoldFiles()
}

func (s *Scaffolder) copyStaticScaffoldFiles() error {
	if !s.CreateVSCodeExtension {
		return nil
	}
	static := []string{
		"vscode-extension/esbuild.js",
		"vscode-extension/language-configuration.json",
	}
	for _, rel := range static {
		body, err := s.readTemplateFile(rel)
		if err != nil {
			return err
		}
		if err := s.writeFile(rel, body); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scaffolder) readTemplateFile(relativePath string) ([]byte, error) {
	// embed.FS paths always use forward slashes
	body, err := templateFS.ReadFile(path.Join("templates", filepath.ToSlash(relativePath)))
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", relativePath, err)
	}
	return body, nil
}

func (s *Scaffolder) renderTemplate(templatePath string, params templateParams) ([]byte, error) {
	body, err := s.readTemplateFile(templatePath)
	if err != nil {
		return nil, err
	}
	t, err := template.New(templatePath).Parse(string(body))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, params); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Scaffolder) writeFile(relativePath string, body []byte) error {
	outPath := filepath.Join(s.scaffoldOutputDir(), relativePath)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, body, 0644); err != nil {
		return err
	}
	return nil
}
