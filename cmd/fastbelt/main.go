// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const rootLongHelp = `Fastbelt is a language engineering toolkit for Go.

It covers lexing, parsing, AST creation, cross-reference linking, workspace
processing, and Language Server Protocol (LSP) support. Fastbelt is inspired by
Xtext and Langium and uses a .fb grammar definition as the entry point.

Use the generate command to produce parser, lexer, types, linker, and service
code from your grammar, or use scaffold to create a full language project
template.`

const generateLongHelp = `Generate code artifacts from a .fb grammar file.

The generated files include parser, lexer, types, linker, and service wiring.
By default, output is written to the current directory and the package name is
derived from the output directory name.

A typical workflow is to iterate on grammar changes and rerun generation after
each step.`

const generateExamples = `  fastbelt generate
  fastbelt generate -g ./grammar.fb
  fastbelt generate -g ./grammar.fb -o ./internal/lang -p lang
  fastbelt generate -g ./grammar.fb --atn -v`

const scaffoldLongHelp = `Scaffold a new language project from templates.

With -module, Fastbelt creates a new module directory (derived from the last
segment of the module path), initializes go.mod, writes templates, then runs
go generate and go mod tidy.

Without -module, Fastbelt scaffolds into an existing module discovered from the
working directory.`

const scaffoldExamples = `  fastbelt scaffold -module example.com/acme/mylang -language "MyLanguage"
  fastbelt scaffold -package internal/lang -language "MyLanguage"
  fastbelt scaffold -module example.com/acme/mylang -language "MyLanguage" -no-vscode`

func main() {
	if err := runCmd(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCmd() error {
	rootCmd := &cobra.Command{
		Use:          "fastbelt",
		Short:        "Generate code from a grammar definition",
		Long:         rootLongHelp,
		SilenceUsage: true,
	}
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.AddCommand(newGenerateCmd(), newScaffoldCmd())
	return rootCmd.Execute()
}

func newGenerateCmd() *cobra.Command {
	opts := generateOptions{}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate code from a grammar definition",
		Long:    generateLongHelp,
		Example: generateExamples,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runGenerateCLI(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.grammarPath, "grammar", "g", "./grammar.fb", "Path to the grammar file")
	cmd.Flags().StringVarP(&opts.outputPath, "output", "o", "./", "Path to the output directory")
	cmd.Flags().StringVarP(&opts.packageName, "package", "p", "", "Package name for generated code (defaults to the last segment of the output path)")
	cmd.Flags().BoolVar(&opts.atn, "atn", false, "Enable markdown output about ATN construction")
	cmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Enable verbose output about written files")
	return cmd
}

func newScaffoldCmd() *cobra.Command {
	opts := scaffoldOptions{
		packagePath: ".",
	}
	cmd := &cobra.Command{
		Use:     "scaffold",
		Short:   "Scaffold a new language project",
		Long:    scaffoldLongHelp,
		Example: scaffoldExamples,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runScaffoldCLI(opts)
		},
	}
	cmd.Flags().StringVar(&opts.modulePath, "module", "", "Module path for go mod init (optional; omit to scaffold into an existing module)")
	cmd.Flags().StringVar(&opts.packagePath, "package", ".", "Template output directory: with --module, relative to the new module root; without --module, relative to the working directory")
	cmd.Flags().StringVar(&opts.language, "language", "", "Human-readable language name")
	cmd.Flags().BoolVar(&opts.noVSCodeExtension, "no-vscode", false, "Do not generate a VS Code extension")
	_ = cmd.MarkFlagRequired("language")
	return cmd
}
