// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/grammar/generated"
)

func GenerateServices(grammar generated.Grammar) string {
	node := generator.NewNode()
	node.AppendLine("package generated")
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n generator.Node) {
		n.AppendLine("\"typefox.dev/fastbelt/workspace\"")
	})
	node.AppendLine(")")
	node.AppendLine()
	node.AppendLine("func CreateDefaultServices(c workspace.GeneratedSrvCont) {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("s := c.Generated()")
		n.AppendLine("if s.Lexer == nil {")
		n.Indent(func(n2 generator.Node) {
			n2.AppendLine("s.Lexer = NewLexer()")
		})
		n.AppendLine("}")
		n.AppendLine("if s.Parser == nil {")
		n.Indent(func(n2 generator.Node) {
			n2.AppendLine("s.Parser = NewParser()")
		})
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()
	return formatIfPossible(node.String())
}
