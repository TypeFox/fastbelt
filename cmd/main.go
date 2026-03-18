// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	grammarPathFlag := flag.String("g", "./grammar.fb", "Path to the grammar file")
	outputPathFlag := flag.String("o", "./", "Path to the output directory")
	packageNameFlag := flag.String("p", "", "Package name for generated code (defaults to the last segment of the output path)")
	verboseFlag := flag.Bool("v", false, "Enable verbose output about written files")
	flag.Parse()

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

	grammr, ok := document.Root.(grammar.Grammar)
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
		generator.GenerateLinker(grammr, packageName)); err != nil {
		return err
	}
	if err := writeFile("types", filepath.Join(outputPath, "types_gen.go"),
		generator.GenerateTypes(grammr, packageName)); err != nil {
		return err
	}
	if err := writeFile("parser", filepath.Join(outputPath, "parser_gen.go"),
		generator.GenerateParser(grammr, packageName)); err != nil {
		return err
	}
	if err := writeFile("lexer", filepath.Join(outputPath, "lexer_gen.go"),
		generator.GenerateLexer(grammr, packageName)); err != nil {
		return err
	}
	if err := writeFile("services", filepath.Join(outputPath, "services_gen.go"),
		generator.GenerateServices(grammr, packageName)); err != nil {
		return err
	}

	return nil
}
