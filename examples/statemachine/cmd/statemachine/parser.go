// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/examples/statemachine"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

// ParseArgs reads exactly one CLI argument into InputPath. ProgName is used for usage messages when set;
// otherwise "statemachine" is used.
func (r *Runner) ParseArgs(args []string) error {
	if len(args) != 1 {
		prog := r.ProgName
		if prog == "" {
			prog = "statemachine"
		}
		return fmt.Errorf("usage: %s <path.statemachine>", prog)
	}
	r.InputPath = args[0]
	return nil
}

// LoadSource reads InputPath from disk, sets Content to the file bytes and SourcePath to an absolute path
// for stable document URIs.
func (r *Runner) LoadSource() error {
	if r.InputPath == "" {
		return fmt.Errorf("runner: InputPath is empty; call ParseArgs first")
	}
	abs, absErr := filepath.Abs(r.InputPath)
	if absErr != nil {
		return absErr
	}
	b, readErr := os.ReadFile(abs)
	if readErr != nil {
		return readErr
	}
	r.Content = b
	r.SourcePath = abs
	return nil
}

// ParseStatemachine runs the fastbelt workspace pipeline (lex, parse, link, validate) on Content and
// SourcePath. On blocking diagnostics it prints to Stderr when non-nil and returns an error. On success it
// stores the built document and typed AST root on the runner.
func (r *Runner) ParseStatemachine(ctx context.Context) error {
	r.doc = nil
	r.sm = nil

	docURI := core.FileURI(r.SourcePath)
	fileDoc, ferr := textdoc.NewFile(docURI.DocumentURI(), "statemachine", 0, string(r.Content))
	if ferr != nil {
		return ferr
	}
	document := core.NewDocument(fileDoc)

	srv := statemachine.CreateServices()
	srv.Workspace().DocumentManager.Set(document)

	var buildErr error
	srv.Workspace().Lock.Write(ctx, func(buildCtx context.Context, downgrade func()) {
		buildErr = srv.Workspace().Builder.Build(buildCtx, []*core.Document{document}, downgrade)
	})
	if buildErr != nil {
		return fmt.Errorf("build: %w", buildErr)
	}

	if fatalBuildIssues(document) {
		if r.Stderr != nil {
			printDiagnostics(r.Stderr, document)
		}
		return fmt.Errorf("document has errors (see diagnostics above)")
	}

	if document.Root == nil {
		return fmt.Errorf("no AST root after build")
	}
	sm, ok := document.Root.(statemachine.Statemachine)
	if !ok {
		return fmt.Errorf("expected a statemachine root AST, got %T", document.Root)
	}

	r.doc = document
	r.sm = sm
	return nil
}

func printDiagnostics(w io.Writer, doc *core.Document) {
	for _, d := range workspace.CreateLexerDiagnostics(doc) {
		_, _ = fmt.Fprintf(w, "%s\n", formatDiagnostic("lexer", d))
	}
	for _, d := range workspace.CreateParserDiagnostics(doc) {
		_, _ = fmt.Fprintf(w, "%s\n", formatDiagnostic("parser", d))
	}
	for _, d := range workspace.CreateLinkerDiagnostics(doc) {
		_, _ = fmt.Fprintf(w, "%s\n", formatDiagnostic("linker", d))
	}
	for _, d := range doc.Diagnostics {
		_, _ = fmt.Fprintf(w, "%s\n", formatDiagnostic("validate", d))
	}
}

func formatDiagnostic(source string, d *core.Diagnostic) string {
	loc := formatLocation(d.Range.Start)
	return fmt.Sprintf("%s %s: %s", source, loc, d.Message)
}

func formatLocation(l core.TextLocation) string {
	return fmt.Sprintf("line %d col %d", int(l.Line)+1, int(l.Column)+1)
}

func fatalBuildIssues(doc *core.Document) bool {
	if len(doc.LexerErrors) > 0 || len(doc.ParserErrors) > 0 {
		return true
	}
	for _, ref := range doc.References {
		if ref.Error() != nil {
			return true
		}
	}
	for _, d := range doc.Diagnostics {
		if d != nil && d.Severity == core.SeverityError {
			return true
		}
	}
	return false
}
