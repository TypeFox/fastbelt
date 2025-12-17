package main

import (
	"os"
	"time"

	"typefox.dev/fastbelt/internal/grammar/generated"
)

func main() {
	bytes, err := os.ReadFile("C:\\Users\\markh\\Desktop\\fastbelt\\input-text.txt")
	if err != nil {
		panic(err)
	}
	text := string(bytes)

	lexer := generated.NewLexer()
	startTime := time.Now()
	result := lexer.Lex(text)
	elapsed := time.Since(startTime)
	println("Lexing took:", elapsed.Milliseconds(), "ms and has an output length of", len(result.Tokens))
}
