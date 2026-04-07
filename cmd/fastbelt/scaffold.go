// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"typefox.dev/fastbelt/internal/scaffold"
)

const scaffoldUsageText = `Usage:
  fastbelt scaffold -module <path> [-package <dir>] -language <name>
  fastbelt scaffold [-package <dir>] -language <name>

With -module: creates a new Go module in a directory named after the final segment of -module
(for example -module=example.com/acme/foo creates ./foo/), runs go mod init, then writes templates
into -package relative to that new module directory (default "." = module root).

Without -module: requires go.mod in the working directory or a parent. Writes templates into
-package relative to the current working directory (default "."). The module path and import path
are read from go.mod. Does not run go mod init.

Flags:
`

func runScaffoldCLI(args []string) error {
	fs := flag.NewFlagSet("scaffold", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	modulePath := fs.String("module", "", "module path for go mod init (optional; omit to scaffold into an existing module)")
	packagePath := fs.String("package", ".", "template output directory: with -module, relative to the new module root; without -module, relative to the working directory")
	language := fs.String("language", "", "human-readable language name (required)")
	// Generate a VS Code extension by default; if the user doesn't want it, they can delete the generated code,
	// but if they decide they want it later, there's no simply way to re-create it.
	noVSCodeExtension := fs.Bool("no-vscode", false, "do not generate a VS Code extension")

	fs.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr, scaffoldUsageText)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *language == "" {
		fs.Usage()
		return fmt.Errorf("-language is required")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	packageRel := strings.TrimSpace(*packagePath)
	if packageRel == "" {
		packageRel = "."
	}

	scaffolder := &scaffold.Scaffolder{
		CreateVSCodeExtension: !(*noVSCodeExtension),
		Language:              *language,
		CreateModule:          *modulePath != "",
	}

	moduleArg := strings.TrimSpace(*modulePath)
	if moduleArg != "" {
		importPath := path.Clean(moduleArg)
		dirBase := path.Base(importPath)
		if dirBase == "." || dirBase == "/" {
			return fmt.Errorf("invalid -module path %q (cannot determine output directory)", *modulePath)
		}
		moduleDir := filepath.Join(wd, dirBase)
		err = scaffolder.PopulateDirectoriesFromModuleDir(moduleDir, importPath, packageRel)
		if err != nil {
			return err
		}
	} else {
		err = scaffolder.PopulateDirectorysFromWorkDir(wd, packageRel)
		if err != nil {
			return err
		}
	}

	err = scaffolder.Run()
	if err != nil {
		debug, _ := json.MarshalIndent(scaffolder, "", "  ")
		return fmt.Errorf("scaffold: %w; scaffolder: %s", err, string(debug))
	}
	return nil
}
