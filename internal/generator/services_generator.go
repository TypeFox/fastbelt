// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"fmt"
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

type Selector struct {
	LanguageID string
	Patterns   []string
}

// GenerateServices emits services_gen.go for grammr. When entryRules has more
// than one element, it generates a language-aware DocumentParser that selects
// the per-language lexer (NewLexerFor<i>) based on the document URI, and
// overrides the default workspace.DocumentParser with it. This prevents
// keywords from one language polluting the lexer of another language.
func GenerateServices(grammr grammar.Grammar, selectors []Selector, entryRules []grammar.ParserRule, packageName string) string {
	multiLang := len(entryRules) > 1

	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/lexer\"")
		n.AppendLine("\"typefox.dev/fastbelt/parser\"")
		n.AppendLine("\"typefox.dev/fastbelt/util/service\"")
		if multiLang {
			n.AppendLine("\"typefox.dev/fastbelt/workspace\"")
		}
	})
	node.AppendLine(")").AppendLine()

	node.AppendLine("// SetupGeneratedServices sets up the generated services for this grammar.")
	node.AppendLine("// If any service is already set, it's not overwritten.")
	node.AppendLine("func SetupGeneratedServices(sc *service.Container) {")
	node.Indent(func(n codegen.Node) {
		if len(selectors) > 0 {
			n.AppendLine("if !service.Has[core.LanguageSelector](sc) {")
			n.Indent(func(n codegen.Node) {
				n.AppendLine("service.Put[core.LanguageSelector](")
				n.Indent(func(n codegen.Node) {
					n.AppendLine("sc,")
					n.AppendLine("core.NewDefaultLanguageSelector(")
					n.Indent(func(n codegen.Node) {
						n.AppendLine("sc,")
						for _, s := range selectors {
							patterns := strings.Join(s.Patterns, "\", \"")
							n.AppendLine("core.NewDocumentSelectorWithPatterns(\"", s.LanguageID, "\", \"", patterns, "\"),")
						}
					})
					n.AppendLine("),")
				})
				n.AppendLine(")")
			})
			n.AppendLine("}")
		}
		n.AppendLine("if !service.Has[", grammr.Name(), "ScopeProvider](sc) {")
		n.AppendLine("    service.Put(sc, NewDefault", grammr.Name(), "ScopeProvider(sc))")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[", grammr.Name(), "ReferenceLinker](sc) {")
		n.AppendLine("    service.Put(sc, NewDefault", grammr.Name(), "ReferenceLinker(sc))")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[", grammr.Name(), "ReferencesConstructor](sc) {")
		n.AppendLine("    service.Put(sc, NewDefault", grammr.Name(), "ReferencesConstructor(sc))")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[", grammr.Name(), "ParserLookahead](sc) {")
		n.AppendLine("    service.Put(sc, NewDefault", grammr.Name(), "ParserLookahead())")
		n.AppendLine("}")
		if multiLang {
			// For multi-language grammars, override the DocumentParser with a
			// language-aware version that selects the per-language lexer based
			// on the document URI. This runs after workspace.SetupDefaultServices
			// has already put the default DocumentParser, so we use Override.
			n.AppendLine("if service.Has[workspace.DocumentParser](sc) {")
			n.Indent(func(n codegen.Node) {
				n.AppendLine("service.Override[workspace.DocumentParser](sc, new", grammr.Name(), "DocumentParser(sc))")
			})
			n.AppendLine("} else {")
			n.Indent(func(n codegen.Node) {
				n.AppendLine("service.Put[workspace.DocumentParser](sc, new", grammr.Name(), "DocumentParser(sc))")
			})
			n.AppendLine("}")
		} else {
			n.AppendLine("if !service.Has[lexer.Lexer](sc) {")
			n.AppendLine("    service.Put(sc, NewLexer())")
			n.AppendLine("}")
		}
		n.AppendLine("if !service.Has[parser.Parser](sc) {")
		n.AppendLine("    service.Put[parser.Parser](sc, NewParser(sc))")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[parser.ErrorRecoveryStrategy](sc) {")
		n.AppendLine("    service.Put[parser.ErrorRecoveryStrategy](sc, parser.NewDefaultErrorRecovery())")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[parser.ErrorMessageProvider](sc) {")
		n.AppendLine("    service.Put[parser.ErrorMessageProvider](sc, parser.NewDefaultErrorMessageProvider())")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[core.SymbolContainers](sc) {")
		n.AppendLine("    service.Put[core.SymbolContainers](sc, NewSymbolContainers())")
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	if multiLang {
		generateLanguageAwareDocumentParser(node, grammr.Name(), len(entryRules))
	}

	node.AppendLine("// SetupGeneratedServerServices sets up the generated language server services for this grammar.")
	node.AppendLine("// If any service is already set, it's not overwritten.")
	node.AppendLine("func SetupGeneratedServerServices(sc *service.Container) {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("if !service.Has[", grammr.Name(), "CompletionFilter](sc) {")
		n.AppendLine("    service.Put(sc, NewDefault", grammr.Name(), "CompletionFilter())")
		n.AppendLine("}")
		n.AppendLine("if !service.Has[parser.LanguageCompletionAdapter](sc) {")
		n.AppendLine("    service.Put[parser.LanguageCompletionAdapter](sc, New", grammr.Name(), "CompletionAdapter(sc))")
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	return FormatIfPossible(node.String())
}

// generateLanguageAwareDocumentParser emits a DocumentParser implementation
// that dispatches to a per-language lexer based on the document URI.
func generateLanguageAwareDocumentParser(node codegen.Node, grammarName string, numLanguages int) {
	typeName := strings.ToLower(grammarName[:1]) + grammarName[1:] + "DocumentParser"

	node.AppendLine("// ", typeName, " is a [workspace.DocumentParser] that selects the per-language")
	node.AppendLine("// lexer based on the document URI, preventing cross-language keyword pollution.")
	node.AppendLine("type ", typeName, " struct {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("sc     *service.Container")
		n.AppendLine("lexers []*lexer.DefaultLexer")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func new", grammarName, "DocumentParser(sc *service.Container) *", typeName, " {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return &", typeName, "{")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("sc: sc,")
			n.AppendLine("lexers: []*lexer.DefaultLexer{")
			n.Indent(func(n codegen.Node) {
				for i := range numLanguages {
					n.AppendLine("NewLexerFor", fmt.Sprintf("%d", i), "(),")
				}
			})
			n.AppendLine("},")
		})
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (p *", typeName, ") Parse(doc *core.Document) {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("text := doc.TextDoc.Text(nil)")
		n.AppendLine("selector := service.MustGet[core.LanguageSelector](p.sc)")
		n.AppendLine("langIdx, _ := selector.Select(core.ParseURI(string(doc.TextDoc.URI())))")
		n.AppendLine("var lex lexer.Lexer")
		n.AppendLine("if langIdx >= 0 && langIdx < len(p.lexers) {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("lex = p.lexers[langIdx]")
		})
		n.AppendLine("} else {")
		n.Indent(func(n codegen.Node) {
			n.AppendLine("lex = p.lexers[0]")
		})
		n.AppendLine("}")
		n.AppendLine("lexerRes := lex.Lex(text)")
		n.AppendLine("doc.LexerErrors = lexerRes.Errors")
		n.AppendLine("doc.Tokens = lexerRes.Tokens")
		n.AppendLine("doc.Comments = lexerRes.Comments")
		n.AppendLine("prs := service.MustGet[parser.Parser](p.sc)")
		n.AppendLine("prsRes := prs.Parse(doc)")
		n.AppendLine("doc.ParserErrors = prsRes.Errors")
		n.AppendLine("doc.Root = prsRes.Node")
		n.AppendLine("core.AssignContainers(doc)")
	})
	node.AppendLine("}")
	node.AppendLine()
}
