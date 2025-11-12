// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"os"
	"path/filepath"

	"github.com/TypeFox/langium-to-go/internal/generated"
	"github.com/TypeFox/langium-to-go/internal/generator"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	grammarPath := filepath.Join(cwd, "..", "grammar.fb")
	grammarText, err := os.ReadFile(grammarPath)
	if err != nil {
		panic(err)
	}
	lexer_test := generated.NewLexer()
	lexerResult := lexer_test.Lex(string(grammarText))
	parser_test := generated.NewParser()
	parserResult := parser_test.Parse(lexerResult.Tokens)
	types := generator.GenerateTypes(parserResult)
	err = os.WriteFile(filepath.Join(cwd, "..", "internal", "generated", "types_gen.go"), []byte(types), 0644)
	if err != nil {
		panic(err)
	}
	generatedParser := generator.GenerateParser(parserResult)
	err = os.WriteFile(filepath.Join(cwd, "..", "internal", "generated", "parser_gen.go"), []byte(generatedParser), 0644)
	if err != nil {
		panic(err)
	}
	lexer := generator.GenerateLexer(parserResult)
	err = os.WriteFile(filepath.Join(cwd, "..", "internal", "generated", "lexer_gen.go"), []byte(lexer), 0644)
	if err != nil {
		panic(err)
	}
}
