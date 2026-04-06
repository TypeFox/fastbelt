// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"fmt"
	"os"

	"text/tabwriter"
)

const fastbeltRepoURL = "https://github.com/TypeFox/fastbelt"

const globalHelpText = `Fastbelt generates Go lexer, parser, types, linker, and service glue from .fb grammars.

Usage:
  fastbelt [flags]		Run code generation
  fastbelt scaffold -module <path> -language <name>	Create a new Go module for a language
  fastbelt scaffold -package <dir> -language <name>	Create a package under the module (needs go.mod)
  fastbelt help		Show this help

The canonical module path for installs is typefox.dev/fastbelt (see %[1]s).

generate flags:
  -g path	Grammar file (default ./grammar.fb)
  -o path	Output directory (default ./)
  -p name	Go package name (default last segment of -o)
  -v		Verbose file writes

Scaffolding:
  Module mode: directory named after the last segment of -module, go mod init, then the
  same file layout, go get (library + tool), go generate, go mod tidy.
  Package mode (-package): finds go.mod from cwd, writes under a directory relative to
  that module (import path inferred from go.mod), skips go mod init, then go get / generate / tidy.
`

func printGlobalHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintf(w, globalHelpText, fastbeltRepoURL)
	_ = w.Flush()
}
