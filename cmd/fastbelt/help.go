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

func printGlobalHelp() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Fastbelt generates Go lexer, parser, types, linker, and service glue from .fb grammars.\n\n")
	fmt.Fprintf(w, "Usage:\n")
	fmt.Fprintf(w, "  %s [flags]\t\tRun code generation (legacy CLI)\n", os.Args[0])
	fmt.Fprintf(w, "  %s scaffold -module <path> -language <name>\tCreate a new Go module for a language\n", os.Args[0])
	fmt.Fprintf(w, "  %s scaffold -package <dir> -language <name>\tCreate a package under the module (needs go.mod)\n", os.Args[0])
	fmt.Fprintf(w, "  %s help\t\tShow this help\n\n", os.Args[0])
	fmt.Fprintf(w, "The canonical module path for installs is typefox.dev/fastbelt (see %s).\n\n", fastbeltRepoURL)
	fmt.Fprintf(w, "Legacy generate flags:\n")
	fmt.Fprintf(w, "  -g path\tGrammar file (default ./grammar.fb)\n")
	fmt.Fprintf(w, "  -o path\tOutput directory (default ./)\n")
	fmt.Fprintf(w, "  -p name\tGo package name (default last segment of -o)\n")
	fmt.Fprintf(w, "  -v\t\tVerbose file writes\n\n")
	fmt.Fprintf(w, "Scaffolding:\n")
	fmt.Fprintf(w, "  Module mode: directory named after the last segment of -module, go mod init, then the\n")
	fmt.Fprintf(w, "  same file layout, go get (library + tool), go generate, go mod tidy.\n")
	fmt.Fprintf(w, "  Package mode (-package): finds go.mod from cwd, writes under a directory relative to\n")
	fmt.Fprintf(w, "  that module (import path inferred from go.mod), skips go mod init, then go get / generate / tidy.\n")
	_ = w.Flush()
}
