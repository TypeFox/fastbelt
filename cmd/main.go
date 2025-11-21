// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"flag"
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
	flag.Parse()

	grammarPath := *grammarPathFlag
	outputPath := *outputPathFlag

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
	types := generator.GenerateTypes(parserResult)
	err = os.WriteFile(filepath.Join(outputPath, "types_gen.go"), []byte(types), 0644)
	if err != nil {
		panic(err)
	}
	generatedParser := generator.GenerateParser(parserResult)
	err = os.WriteFile(filepath.Join(outputPath, "parser_gen.go"), []byte(generatedParser), 0644)
	if err != nil {
		panic(err)
	}
	lexer := generator.GenerateLexer(parserResult)
	err = os.WriteFile(filepath.Join(outputPath, "lexer_gen.go"), []byte(lexer), 0644)
	if err != nil {
		panic(err)
	}
}
