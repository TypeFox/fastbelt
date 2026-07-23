// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/cmd"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

type generateOptions struct {
	grammarPath string
	outputPath  string
	packageName string
	atn         bool
	verbose     bool
}

func runGenerateCLI(opts generateOptions) error {
	grammarPath, err := filepath.Abs(opts.grammarPath)
	if err != nil {
		return err
	}
	outputPath, err := filepath.Abs(opts.outputPath)
	if err != nil {
		return err
	}
	verbose := opts.verbose

	packageName := opts.packageName
	if packageName == "" {
		packageName = filepath.Base(outputPath)
	}

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	grammarText, err := os.ReadFile(grammarPath)
	if err != nil {
		return err
	}

	sc := grammar.CreateServices()
	file := textdoc.NewFile(lsp.URIFromPath(grammarPath), "fb", 0, string(grammarText))

	document := core.NewDocument(file)
	documents, err := service.Get[workspace.DocumentManager](sc)
	if err != nil {
		return err
	}
	documents.Set(document)
	builder, err := service.Get[workspace.Builder](sc)
	if err != nil {
		return err
	}
	if err := builder.Build(context.Background(), []*core.Document{document}, nil); err != nil {
		return err
	}

	diagnostics := document.Diagnostics
	errCount := 0

	sort.SliceStable(diagnostics, func(i, j int) bool {
		iStartLine := diagnostics[i].Range.Start.Line
		jStartLine := diagnostics[j].Range.Start.Line
		if iStartLine == jStartLine {
			return diagnostics[i].Range.Start.Column < diagnostics[j].Range.Start.Column
		} else {
			return iStartLine < jStartLine
		}
	})

	for _, diag := range diagnostics {
		if diag.Severity == core.SeverityError {
			errCount++
		}
		fmt.Printf(
			"%s - %d:%d %s\n",
			diag.Severity.String(),
			// For printing, convert to 1-based line and column numbers.
			diag.Range.Start.Line+1,
			diag.Range.Start.Column+1,
			diag.Message,
		)
	}

	if errCount > 0 {
		return fmt.Errorf("aborting code generation due to %d errors", errCount)
	}

	g, ok := document.Root.(grammar.Grammar)
	if !ok {
		return fmt.Errorf("parser result is not a Grammar")
	}
	entryRule, err := validateEntryRule(g)
	if err != nil {
		return err
	}

	// Delegate code generation to the shared build API. A single-language CLI
	// build passes exactly one entry rule, so the generated parser keeps its
	// direct (non-dispatching) Parse body.
	return cmd.Generate(g, []grammar.ParserRule{entryRule}, nil, outputPath, packageName, opts.atn, verbose)
}

func validateEntryRule(g grammar.Grammar) (grammar.ParserRule, error) {
	var entries []grammar.ParserRule
	for _, rule := range g.Rules() {
		if rule.IsEntry() {
			entries = append(entries, rule)
		}
	}
	switch len(entries) {
	case 1:
		return entries[0], nil
	case 0:
		return nil, fmt.Errorf("grammar must have exactly one parser rule marked as entry, but none were found")
	default:
		names := make([]string, len(entries))
		for i, rule := range entries {
			names[i] = rule.Name()
		}
		return nil, fmt.Errorf(
			"grammar must have exactly one parser rule marked as entry, but found %d: %s",
			len(entries),
			strings.Join(names, ", "),
		)
	}
}
