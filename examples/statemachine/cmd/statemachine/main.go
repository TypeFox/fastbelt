// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Command statemachine loads a .statemachine file through the same workspace and document pipeline used in
// production (see parser.go), prints diagnostics and a model summary, then reads event names from stdin to step
// the machine (see interpreter.go). runner.go defines Runner and ties the phases together.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	r := &Runner{
		ProgName:   filepath.Base(os.Args[0]),
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		EventInput: os.Stdin,
	}
	if err := r.ParseArgs(os.Args[1:]); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}
	if err := r.LoadSource(); err != nil {
		return fmt.Errorf("load source: %w", err)
	}
	if err := r.ParseAndValidate(); err != nil {
		return fmt.Errorf("parse and validate: %w", err)
	}
	if err := r.Run(); err != nil {
		return fmt.Errorf("run: %w", err)
	}
	return nil
}
