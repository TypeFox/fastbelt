// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"typefox.dev/fastbelt/cmd"
	"typefox.dev/fastbelt/internal/grammar"
)

type generateOptions struct {
	grammarPath string
	outputPath  string
	packageName string
	atn         bool
	railroad    bool
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

	// A grammar consists of all .fb files in a directory, so a file argument
	// refers to its directory.
	grammarDir := grammarPath
	if info, err := os.Stat(grammarPath); err != nil {
		return err
	} else if !info.IsDir() {
		grammarDir = filepath.Dir(grammarPath)
	}
	g, err := cmd.LoadGrammarDir(grammarDir)
	if err != nil {
		return err
	}
	entryRule, err := validateEntryRule(g)
	if err != nil {
		return err
	}
	// The CLI only supports a single language. Use the cmd.BuildContext API
	// to generate code for multiple languages.
	return cmd.Generate(g, []grammar.ParserRule{entryRule}, nil, nil, outputPath, packageName, opts.atn, opts.railroad, verbose)
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
