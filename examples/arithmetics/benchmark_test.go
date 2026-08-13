// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"fmt"
	"strings"
	"testing"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/service"
)

// generateArithmeticsContent produces an expression-heavy module
func generateArithmeticsContent(statements int) string {
	var sb strings.Builder
	sb.WriteString("module benchmark\n")
	sb.WriteString("def base: 42;\n")
	sb.WriteString("def calc(x, y): x * y + base % 2 - (x / y) ^ 2;\n")
	for i := range statements {
		fmt.Fprintf(&sb, "def val%d: %d + %d * 3 - base / (%d + 1);\n", i, i, i+1, i+2)
		fmt.Fprintf(&sb, "calc(val%d, %d) + val%d * (base - %d);\n", i, i, i, i)
	}
	return sb.String()
}

// BenchmarkParser benchmarks parsing a single expression-heavy arithmetics
// document, reusing the pre-lexed token slice every iteration.
func BenchmarkParser(b *testing.B) {
	content := generateArithmeticsContent(50)
	srv := CreateServices()
	lexerService := service.MustGet[lexer.Lexer](srv)
	parserService := service.MustGet[parser.Parser](srv)
	lexed := lexerService.Exec(content)
	if len(lexed.Errors) > 0 {
		b.Fatalf("lexer errors: %v", lexed.Errors)
	}
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/benchmark.arithmetics", "arithmetics", content)
	if err != nil {
		b.Fatal(err)
	}
	doc.Tokens = lexed.Tokens
	if result := parserService.Parse(doc); len(result.Errors) > 0 {
		b.Fatalf("parser errors: %v", result.Errors)
	}
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		result := parserService.Parse(doc)
		doc.Root = result.Node
	}
}
