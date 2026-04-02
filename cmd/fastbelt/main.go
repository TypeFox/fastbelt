// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/generator"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/lsp"
)

func main() {
	if err := runCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCmd() error {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "help", "-h", "-help", "--help":
			printGlobalHelp()
			return nil
		case "scaffold":
			return runScaffoldCLI(args[1:])
		}
	}
	return runLegacyGenerate(args)
}

func runLegacyGenerate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	grammarPathFlag := fs.String("g", "./grammar.fb", "Path to the grammar file")
	outputPathFlag := fs.String("o", "./", "Path to the output directory")
	packageNameFlag := fs.String("p", "", "Package name for generated code (defaults to the last segment of the output path)")
	verboseFlag := fs.Bool("v", false, "Enable verbose output about written files")

	fs.Usage = func() {
		printGlobalHelp()
		fmt.Fprintf(os.Stderr, "\nGenerate flags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	grammarPath, err := filepath.Abs(*grammarPathFlag)
	if err != nil {
		return err
	}
	outputPath, err := filepath.Abs(*outputPathFlag)
	if err != nil {
		return err
	}
	verbose := *verboseFlag

	packageName := *packageNameFlag
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

	srv := grammar.CreateServices()
	file, _ := textdoc.NewFile(lsp.URIFromPath(grammarPath), "fb", 0, string(grammarText))

	document := core.NewDocument(file)
	srv.Workspace().DocumentManager.Set(document)
	if err := srv.Workspace().Builder.Build(context.Background(), []*core.Document{document}, func() {}); err != nil {
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

	grammar, ok := document.Root.(grammar.Grammar)
	if !ok {
		return fmt.Errorf("parser result is not a Grammar")
	}

	writeFile := func(name, path, content string) error {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", name, err)
		}
		if verbose {
			fmt.Printf("Written: %s\n", path)
		}
		return nil
	}

	if err := writeFile("linker", filepath.Join(outputPath, "linker_gen.go"),
		generator.GenerateLinker(grammar, packageName)); err != nil {
		return err
	}
	if err := writeFile("types", filepath.Join(outputPath, "types_gen.go"),
		generator.GenerateTypes(grammar, packageName)); err != nil {
		return err
	}
	if err := writeFile("parser", filepath.Join(outputPath, "parser_gen.go"),
		generator.GenerateParser(grammar, packageName)); err != nil {
		return err
	}
	if err := writeFile("lexer", filepath.Join(outputPath, "lexer_gen.go"),
		generator.GenerateLexer(grammar, packageName)); err != nil {
		return err
	}
	if err := writeFile("services", filepath.Join(outputPath, "services_gen.go"),
		generator.GenerateServices(grammar, packageName)); err != nil {
		return err
	}

	return nil
}
