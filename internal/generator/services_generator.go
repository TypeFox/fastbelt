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
		n.AppendLine("\"typefox.dev/fastbelt/linking\"")
		n.AppendLine("\"typefox.dev/fastbelt/workspace\"")
	})
	node.AppendLine(")").AppendLine()

	node.AppendLine("type ", grammar.Name(), "LinkingSrv struct {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("ScopeProvider       ", grammar.Name(), "ScopeProvider")
		n.AppendLine("ReferenceLinker     ", grammar.Name(), "ReferenceLinker")
		n.AppendLine("ReferencesConstructor ", grammar.Name(), "ReferencesConstructor")
	})
	node.AppendLine("}").AppendLine()
	node.AppendLine("type ", grammar.Name(), "LinkingSrvCont interface {")
	node.AppendLine("	linking.LinkingSrvCont")
	node.AppendLine("	", grammar.Name(), "Linking() *", grammar.Name(), "LinkingSrv")
	node.AppendLine("}").AppendLine()

	node.AppendLine("type ", grammar.Name(), "LinkingSrvContBlock struct {")
	node.AppendLine("	", "fastbeltLinking ", grammar.Name(), "LinkingSrv")
	node.AppendLine("}").AppendLine()
	node.AppendLine("func (b *", grammar.Name(), "LinkingSrvContBlock) ", grammar.Name(), "Linking() *", grammar.Name(), "LinkingSrv {")
	node.AppendLine("	return &b.fastbeltLinking")
	node.AppendLine("}").AppendLine()

	node.AppendLine("type ", grammar.Name(), "GeneratedSrvCont interface {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("workspace.GeneratedSrvCont")
		n.AppendLine(grammar.Name(), "LinkingSrvCont")
	})
	node.AppendLine("}").AppendLine()

	node.AppendLine("func CreateDefaultServices(srv ", grammar.Name(), "GeneratedSrvCont) {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("linking := srv.", grammar.Name(), "Linking()")
		n.AppendLine("if linking.ScopeProvider == nil {")
		n.AppendLine("    linking.ScopeProvider = NewDefault", grammar.Name(), "ScopeProvider(srv)")
		n.AppendLine("}")
		n.AppendLine("if linking.ReferenceLinker == nil {")
		n.AppendLine("    linking.ReferenceLinker = NewDefault", grammar.Name(), "ReferenceLinker(srv)")
		n.AppendLine("}")
		n.AppendLine("if linking.ReferencesConstructor == nil {")
		n.AppendLine("    linking.ReferencesConstructor = NewDefault", grammar.Name(), "ReferencesConstructor(srv)")
		n.AppendLine("}")
		n.AppendLine("generated := srv.Generated()")
		n.AppendLine("if generated.Lexer == nil {")
		n.Indent(func(n2 generator.Node) {
			n2.AppendLine("generated.Lexer = NewLexer()")
		})
		n.AppendLine("}")
		n.AppendLine("if generated.Parser == nil {")
		n.Indent(func(n2 generator.Node) {
			n2.AppendLine("generated.Parser = New", grammar.Name(), "Parser(srv)")
		})
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()
	return FormatIfPossible(node.String())
}
