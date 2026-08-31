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
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
)

// TestBenchmarkContentParses guards the benchmarks: a syntax error in the
// generated corpus would silently benchmark the error-recovery path instead.
func TestBenchmarkContentParses(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(generateArithmeticsContent()).AssertNoParseErrors()
}

// BenchmarkParser benchmarks parsing a single generated arithmetics document,
// reusing the pre-lexed token slice every iteration.
func BenchmarkParser(b *testing.B) {
	content := generateArithmeticsContent()
	srv := CreateServices()
	lexerService := service.MustGet[lexer.Lexer](srv)
	parserService := service.MustGet[parser.Parser](srv)
	tokens := lexerService.Lex(content).Tokens
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/bench.calc", "arithmetics", content)
	if err != nil {
		b.Fatal(err)
	}
	doc.Tokens = tokens
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		result := parserService.Parse(doc)
		doc.Root = result.Node
	}
}

// BenchmarkLexerAndParser benchmarks a full lex + parse pass per iteration.
func BenchmarkLexerAndParser(b *testing.B) {
	content := generateArithmeticsContent()
	srv := CreateServices()
	lexerService := service.MustGet[lexer.Lexer](srv)
	parserService := service.MustGet[parser.Parser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/bench.calc", "arithmetics", content)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		doc.Tokens = lexerService.Lex(content).Tokens
		result := parserService.Parse(doc)
		doc.Root = result.Node
	}
}

// generateArithmeticsContent generates a syntactically valid arithmetics
// document dominated by binary expressions: 100 definitions and 100
// evaluations, each holding a 40-operator expression mixing all six operators
// with parenthesized sub-expressions and calls to earlier definitions.
func generateArithmeticsContent() string {
	const numDefs = 100
	const opsPerExpr = 40
	operators := []string{"+", "-", "*", "/", "^", "%"}

	var sb strings.Builder
	sb.WriteString("module bench\n\n")

	writeExpr := func(defIdx int) {
		for op := range opsPerExpr {
			operand := fmt.Sprintf("%d", (defIdx+op)%97+1)
			if op%7 == 3 && defIdx > 0 {
				// Reference an earlier definition to include cross-references.
				operand = fmt.Sprintf("d%d", defIdx%op)
			}
			if op%5 == 2 {
				operand = "(" + operand + " " + operators[(defIdx+op)%len(operators)] + " 2)"
			}
			sb.WriteString(operand)
			if op < opsPerExpr-1 {
				sb.WriteString(" ")
				sb.WriteString(operators[(defIdx*7+op)%len(operators)])
				sb.WriteString(" ")
			}
		}
	}

	for d := range numDefs {
		fmt.Fprintf(&sb, "def d%d: ", d)
		writeExpr(d)
		sb.WriteString(";\n")
	}
	for d := range numDefs {
		writeExpr(d)
		sb.WriteString(";\n")
	}
	return sb.String()
}
