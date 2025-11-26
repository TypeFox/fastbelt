// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"typefox.dev/fastbelt/internal/generator"
	"typefox.dev/fastbelt/internal/grammar/generated"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	grammarPathFlag := flag.String("g", filepath.Join(cwd, "..", "grammar.fb"), "Path to the grammar file")
	outputPathFlag := flag.String("o", filepath.Join(cwd, "..", "internal", "generated"), "Path to the output directory")
	verboseFlag := flag.Bool("v", false, "Enable verbose output about written files")
	flag.Parse()

	grammarPath, err := filepath.Abs(*grammarPathFlag)
	if err != nil {
		panic(err)
	}
	outputPath, err := filepath.Abs(*outputPathFlag)
	if err != nil {
		panic(err)
	}
	verbose := *verboseFlag

	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		panic(err)
	}

	grammarText, err := os.ReadFile(grammarPath)
	if err != nil {
		panic(err)
	}

	lexer_test := generated.NewLexer()
	lexerResult := lexer_test.Lex(string(grammarText))

	parser_test := generated.NewParser()
	parserResult := parser_test.Parse(lexerResult.Tokens)
	if grammar, ok := parserResult.(generated.Grammar); !ok {
		panic("Parser result is not a Grammar")
	} else {
		types := generator.GenerateTypes(grammar)
		typesPath := filepath.Join(outputPath, "types_gen.go")
		err = os.WriteFile(typesPath, []byte(types), 0644)
		if err != nil {
			panic(err)
		}
		if verbose {
			fmt.Printf("Written: %s\n", typesPath)
		}
		generatedParser := generator.GenerateParser(grammar)
		parserPath := filepath.Join(outputPath, "parser_gen.go")
		err = os.WriteFile(parserPath, []byte(generatedParser), 0644)
		if err != nil {
			panic(err)
		}
		if verbose {
			fmt.Printf("Written: %s\n", parserPath)
		}
		lexer := generator.GenerateLexer(grammar)
		lexerPath := filepath.Join(outputPath, "lexer_gen.go")
		err = os.WriteFile(lexerPath, []byte(lexer), 0644)
		if err != nil {
			panic(err)
		}
		if verbose {
			fmt.Printf("Written: %s\n", lexerPath)
		}
		services := generator.GenerateServices(grammar)
		servicesPath := filepath.Join(outputPath, "services_gen.go")
		err = os.WriteFile(servicesPath, []byte(services), 0644)
		if err != nil {
			panic(err)
		}
		if verbose {
			fmt.Printf("Written: %s\n", servicesPath)
		}
	}
}
