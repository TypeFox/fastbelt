// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"fmt"
	"io"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/examples/statemachine"
)

// Runner is the CLI harness: it wires stdin/stdout/stderr to parsing and interpretation.
//
// Parser-related methods live in parser.go; interpreter methods in interpreter.go; this file holds the type
// definition and the orchestration entrypoints ParseAndValidate and Run.
type Runner struct {
	// ProgName is the program basename for usage text (e.g. filepath.Base(os.Args[0])). Set before ParseArgs.
	ProgName string
	// InputPath is the model path from the CLI, set by ParseArgs. LoadSource reads it into Content and SourcePath.
	InputPath  string
	Content    []byte
	SourcePath string
	Stdout     io.Writer
	Stderr     io.Writer
	EventInput io.Reader

	doc *core.Document
	sm  statemachine.Statemachine
}

// ParseAndValidate runs the workspace build and stores the document and root AST on success.
func (r *Runner) ParseAndValidate() error {
	return r.ParseStatemachine(context.Background())
}

// Run prints diagnostics and a model summary, then runs the interactive interpreter on EventInput.
func (r *Runner) Run() error {
	if r.doc == nil || r.sm == nil {
		return fmt.Errorf("runner: call ParseAndValidate before Run")
	}
	if r.Stdout == nil || r.Stderr == nil {
		return fmt.Errorf("runner: Stdout and Stderr must be non-nil")
	}
	if r.EventInput == nil {
		return fmt.Errorf("runner: EventInput must be non-nil")
	}

	printDiagnostics(r.Stderr, r.doc)
	if err := r.PrintModelSummary(); err != nil {
		return err
	}

	_, _ =fmt.Fprintf(r.Stderr, "\nEnter event names, one per line (EOF to stop). Unknown events are reported.\n")
	return r.Interpret()
}
