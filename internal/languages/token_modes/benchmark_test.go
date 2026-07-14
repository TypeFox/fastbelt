package token_modes

import (
	"testing"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/service"
)

func BenchmarkNestedString(b *testing.B) {
	content, _ := generateNestedString()
	srv := CreateServices()
	lexerService := service.MustGet[lexer.Lexer](srv)
	parserService := service.MustGet[parser.Parser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/nested_string.mode", "modes", content)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		lexerResult := lexerService.Exec(content)
		doc.Tokens = lexerResult.Tokens
		parserResult := parserService.Parse(doc)
		doc.Root = parserResult.Node
	}
}

func generateNestedString() (string, error) {
	content := ""

	var recursion func(depth int)
	recursion = func(depth int) {
		if depth == 0 {
			content += "`EXIT!!!`"
			return
		}
		content += "`Hello, World! #{"
		recursion(depth - 1)
		content += "} How are you?`"
	}

	recursion(10000)

	return "VAR := " + content, nil
}
