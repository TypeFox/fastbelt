// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"typefox.dev/fastbelt/internal/scaffold"
)

const scaffoldUsageText = `Usage:
  fastbelt scaffold -module <path> -language <name>
  fastbelt scaffold -package <dir-or-import> -language <name>

Module mode (-module): creates a new Go module under a directory named after the final
segment of -module (for example, -module=example.com/acme/foo creates ./foo/).

Package mode (-package): requires go.mod in the current directory or a parent. The
argument is usually a path relative to the module root (for example -package=examples/mylang
creates ./examples/mylang/); the import path is inferred from the module line in go.mod.
You may still pass a full import path (module path + suffix) if you prefer. Does not run go mod init.

Flags:
`

func runScaffoldCLI(args []string) error {
	fs := flag.NewFlagSet("scaffold", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	modulePath := fs.String("module", "", "module path for go mod init (use this or -package, not both)")
	packagePath := fs.String("package", "", "package directory relative to go.mod (or full import path under the module); requires go.mod; use this or -module, not both")
	language := fs.String("language", "", "human-readable language name (required)")

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
	if *modulePath != "" && *packagePath != "" {
		fs.Usage()
		return fmt.Errorf("use either -module or -package, not both")
	}
	if *modulePath == "" && *packagePath == "" {
		fs.Usage()
		return fmt.Errorf("one of -module or -package is required")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if *packagePath != "" {
		modRoot, pkgDir, pkgImport, resolveErr := scaffold.ResolvePackageScaffoldDir(wd, *packagePath)
		if resolveErr != nil {
			return resolveErr
		}
		if err := scaffold.RunPackage(modRoot, pkgDir, pkgImport, *language); err != nil {
			return err
		}
		fmt.Printf("Scaffolded package at %s\n", pkgDir)
		return nil
	}

	outDir := filepath.Join(wd, filepath.Base(*modulePath))
	if err := scaffold.RunModule(outDir, *modulePath, *language); err != nil {
		return err
	}
	fmt.Printf("Scaffolded module at %s\n", outDir)
	return nil
}
