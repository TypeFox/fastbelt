// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"strconv"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

// Selector describes how documents are matched to one language.
type Selector struct {
	// LanguageID is the LSP language ID of the language.
	LanguageID string
	// Patterns is a list of glob patterns matching the documents of the language.
	Patterns []string
}

func GenerateServices(grammr grammar.Grammar, selectors []Selector, packageName string) string {
	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/lexer\"")
		n.AppendLine("\"typefox.dev/fastbelt/parser\"")
		n.AppendLine("\"typefox.dev/fastbelt/util/service\"")
		if len(selectors) > 0 {
			n.AppendLine("\"typefox.dev/fastbelt/workspace\"")
		}
	})
	node.AppendLine(")").AppendLine()

	node.AppendLine("// SetupGeneratedServices sets up the generated services for this grammar.")
	node.AppendLine("// If any service is already set, it's not overwritten.")
	node.AppendLine("func SetupGeneratedServices(sc *service.Container) {")
	node.Indent(func(n codegen.Node) {
		if len(selectors) > 0 {
			n.AppendLine("if !sc.Has[core.LanguageSelector]() {")
			n.Indent(func(n codegen.Node) {
				n.AppendLine("sc.Put[core.LanguageSelector](")
				n.Indent(func(n codegen.Node) {
					n.AppendLine("workspace.NewDefaultLanguageSelector(")
					n.Indent(func(n codegen.Node) {
						n.AppendLine("sc,")
						for _, s := range selectors {
							// Quote the literals, patterns might contain `"` or `\`.
							n.Append("workspace.NewDocumentMatcherWithPatterns(", strconv.Quote(s.LanguageID))
							for _, pattern := range s.Patterns {
								n.Append(", ", strconv.Quote(pattern))
							}
							n.AppendLine("),")
						}
					})
					n.AppendLine("),")
				})
				n.AppendLine(")")
			})
			n.AppendLine("}")
		}
		n.AppendLine("if !sc.Has[", grammr.Name(), "ScopeProvider]() {")
		n.AppendLine("    sc.Put(NewDefault", grammr.Name(), "ScopeProvider(sc))")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[", grammr.Name(), "ReferenceLinker]() {")
		n.AppendLine("    sc.Put(NewDefault", grammr.Name(), "ReferenceLinker(sc))")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[", grammr.Name(), "ReferencesConstructor]() {")
		n.AppendLine("    sc.Put(NewDefault", grammr.Name(), "ReferencesConstructor(sc))")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[", grammr.Name(), "ParserLookahead]() {")
		n.AppendLine("    sc.Put(NewDefault", grammr.Name(), "ParserLookahead())")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[lexer.Lexer]() {")
		n.AppendLine("    sc.Put(NewLexer(sc))")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[parser.Parser]() {")
		n.AppendLine("    sc.Put[parser.Parser](NewParser(sc))")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[parser.ErrorRecoveryStrategy]() {")
		n.AppendLine("    sc.Put[parser.ErrorRecoveryStrategy](parser.NewDefaultErrorRecovery())")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[parser.ErrorMessageProvider]() {")
		n.AppendLine("    sc.Put[parser.ErrorMessageProvider](parser.NewDefaultErrorMessageProvider())")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[core.SymbolContainers]() {")
		n.AppendLine("    sc.Put[core.SymbolContainers](NewSymbolContainers())")
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("// SetupGeneratedServerServices sets up the generated language server services for this grammar.")
	node.AppendLine("// If any service is already set, it's not overwritten.")
	node.AppendLine("func SetupGeneratedServerServices(sc *service.Container) {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("if !sc.Has[", grammr.Name(), "CompletionFilter]() {")
		n.AppendLine("    sc.Put(NewDefault", grammr.Name(), "CompletionFilter())")
		n.AppendLine("}")
		n.AppendLine("if !sc.Has[parser.LanguageCompletionAdapter]() {")
		n.AppendLine("    sc.Put[parser.LanguageCompletionAdapter](New", grammr.Name(), "CompletionAdapter(sc))")
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	return FormatIfPossible(node.String())
}
