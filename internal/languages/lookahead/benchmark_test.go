// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lookahead

import (
	"strconv"
	"strings"
	"testing"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/test"
)

// deepJRangeDepth is the nesting depth of the generated JRange benchmark input.
const deepJRangeDepth = 400

// generateDeepJRange builds "j x0 is in-between x1 is in-between ... z and w and w ...",
// a JRange nested depth levels deep through its first operand. Every level
// pushes a rule frame and runs the adaptive JRangeLoop decision three times
// (on "is", on "and" and after the "and" operand), so a parse makes about
// 3*depth adaptive predictions under a rule stack up to depth+2 frames deep.
// The prediction DFA resolves all of them without the full-context fallback.
func generateDeepJRange(depth int) string {
	var sb strings.Builder
	sb.WriteString("j ")
	for i := range depth {
		sb.WriteString("x")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(" is in-between ")
	}
	sb.WriteString("z")
	for range depth {
		sb.WriteString(" and w")
	}
	return sb.String()
}

// TestDeepNestedJRangeContentParses guards BenchmarkDeepNestedJRange:
// a syntax error would silently benchmark the error-recovery path instead.
func TestDeepNestedJRangeContentParses(t *testing.T) {
	doc := test.New(t, CreateServices()).Parse(generateDeepJRange(deepJRangeDepth))
	doc.AssertNoParseErrors()
}

// BenchmarkDeepNestedJRange benchmarks adaptive prediction in LL mode
// (PredictionModeLL, see CreateServices) with a deep rule stack and many
// adaptive decisions.
func BenchmarkDeepNestedJRange(b *testing.B) {
	content := generateDeepJRange(deepJRangeDepth)
	srv := CreateServices()
	lexerService := srv.MustGet[lexer.Lexer]()
	parserService := srv.MustGet[parser.Parser]()
	doc := fastbelt.NewDocumentFromString("file:///workspace/deep.la", "lookahead", content)
	lexerService.Exec(b.Context(), doc)
	if len(doc.LexerErrors) > 0 {
		b.Fatalf("lexer errors: %v", doc.LexerErrors)
	}
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		parserService.Parse(b.Context(), doc)
	}
}
