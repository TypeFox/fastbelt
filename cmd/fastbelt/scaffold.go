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

func runScaffoldCLI(args []string) error {
	fs := flag.NewFlagSet("scaffold", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	modulePath := fs.String("module", "", "module path for go mod init (required)")
	language := fs.String("language", "", "human-readable language name (required)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s scaffold -module <path> -language <name>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Creates a new Go module under a directory named after the final segment of -module\n"+
			"(for example, -module=example.com/acme/foo creates ./foo/).\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSee also: %s help\n", os.Args[0])
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *modulePath == "" || *language == "" {
		fs.Usage()
		return fmt.Errorf("-module and -language are required")
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	outDir := filepath.Join(wd, filepath.Base(*modulePath))
	if err := scaffold.RunModule(outDir, *modulePath, *language); err != nil {
		return err
	}
	fmt.Printf("Scaffolded module at %s\n", outDir)
	return nil
}
