package generator

import (
	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/grammar/generated"
)

func GenerateATN(grammar generated.Grammar) string {
	node := generator.NewNode()

	return formatIfPossible(node.String())
}
