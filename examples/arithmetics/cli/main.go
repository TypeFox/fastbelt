// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	core "typefox.dev/fastbelt"
	arithmetics "typefox.dev/fastbelt/examples/arithmetics"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

const rootLongHelp = `Arithmetics CLI is a tool for working with .arithmetics documents.

It supports exporting an arithmetics document to JSON and importing a JSON
description back into an arithmetics document.`

const exportLongHelp = `Export an .arithmetics document as JSON.

By default the JSON is written to stdout. Use --output to write to a file
instead.`

const exportExamples = `  arithmetics export ./example.arithmetics
  arithmetics export ./example.arithmetics -o ./example.json`

const importLongHelp = `Import an arithmetics description from JSON.

By default the result is written to stdout. Use --output to write to a file
instead.`

const importExamples = `  arithmetics import ./example.json
  arithmetics import ./example.json -o ./example.arithmetics`

func main() {
	if err := runCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCmd() error {
	rootCmd := &cobra.Command{
		Use:          "arithmetics",
		Short:        "Work with .arithmetics documents",
		Long:         rootLongHelp,
		SilenceUsage: true,
	}
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.AddCommand(newExportCmd(), newImportCmd())
	return rootCmd.Execute()
}

type exportOptions struct {
	inputPath  string
	outputPath string
}

func newExportCmd() *cobra.Command {
	opts := exportOptions{}
	cmd := &cobra.Command{
		Use:     "export <input.arithmetics>",
		Short:   "Export an .arithmetics document as JSON",
		Long:    exportLongHelp,
		Example: exportExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			opts.inputPath = args[0]
			return runExportCLI(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.outputPath, "output", "o", "", "Output file path (defaults to stdout)")
	return cmd
}

type importOptions struct {
	inputPath  string
	outputPath string
}

func newImportCmd() *cobra.Command {
	opts := importOptions{}
	cmd := &cobra.Command{
		Use:     "import <input.json>",
		Short:   "Import an arithmetics description from JSON",
		Long:    importLongHelp,
		Example: importExamples,
		Args:    cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			opts.inputPath = args[0]
			return runImportCLI(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.outputPath, "output", "o", "", "Output file path (defaults to stdout)")
	return cmd
}

func runExportCLI(opts exportOptions) error {
	inputPath, err := filepath.Abs(opts.inputPath)
	if err != nil {
		return err
	}

	inputText, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	sc := arithmetics.CreateServices()
	file, _ := textdoc.NewFile(lsp.URIFromPath(inputPath), "arithmetics", 0, string(inputText))

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
		return diagnostics[i].Range.Start < diagnostics[j].Range.Start
	})

	for _, diag := range diagnostics {
		if diag.Severity == core.SeverityError {
			errCount++
		}
		tRange := diag.Range.LspRange(document.TextDoc)
		fmt.Fprintf(os.Stderr, "%s - %d:%d %s\n",
			diag.Severity.String(),
			tRange.Start, tRange.End,
			diag.Message,
		)
	}

	if errCount > 0 {
		return fmt.Errorf("aborting export due to %d errors", errCount)
	}

	module, ok := document.Root.(arithmetics.Module)
	if !ok {
		return fmt.Errorf("parser result is not a Module")
	}

	out, err := json.Marshal(module, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}

	if opts.outputPath == "" {
		_, err = fmt.Fprintf(os.Stdout, "%s\n", out)
		return err
	}
	return os.WriteFile(opts.outputPath, append(out, '\n'), 0644)
}

func runImportCLI(opts importOptions) error {
	inputPath, err := filepath.Abs(opts.inputPath)
	if err != nil {
		return err
	}
	file, err := os.Open(inputPath)

	if err != nil {
		return err
	}
	document, err := core.NewDocumentFromString("file:///"+filepath.Base(inputPath), "arithmetics", "")
	if err != nil {
		return err
	}
	sc := arithmetics.CreateServices()

	if err := util.UnmarshalAndBuildDocument(context.Background(), sc, document, file, arithmetics.ArithmeticsSyntheticFactories); err != nil {
		return err
	}

	// re-marshal for the sake of verification
	out, err := json.Marshal(document.Root, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stdout, "%s\n", out)
	return err
}
