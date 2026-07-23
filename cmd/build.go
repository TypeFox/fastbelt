// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package cmd is the importable build API for Fastbelt. It generates a
// parser/lexer/LSP server from one or more .fb grammar files and, when several
// languages are configured, wires them for dispatch inside a single language
// server. The fastbelt CLI (cmd/fastbelt) is a thin wrapper over this package.
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/generator"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// Language configures one language served by the generated language server: the
// grammar entry rule it parses from and the [Selector] that claims its
// documents.
type Language struct {
	Entry      string
	LanguageID string
	Patterns   []string
}

// Plugin can mutate a [BuildContext] before its languages are resolved. It is
// intentionally minimal; richer hooks are deferred until a concrete need.
type Plugin func(*BuildContext)

// BuildContext drives a programmatic build. Input points at a directory of .fb
// files (all sharing the same grammar name); Languages selects one entry rule
// each. With a single language the result is equivalent to the fastbelt
// generate CLI.
type BuildContext struct {
	Input     string
	Output    string // defaults to Input
	Package   string // defaults to filepath.Base(Output)
	Languages []Language
	Plugins   []Plugin
	ATN       bool
	Verbose   bool
}

// Build parses and links every .fb file under Input, validates that each
// language's entry rule exists and is marked entry, and generates the combined
// language package into Output.
func (c *BuildContext) Build() error {
	for _, p := range c.Plugins {
		p(c)
	}
	if len(c.Languages) == 0 {
		return fmt.Errorf("no languages configured")
	}
	fullInput, err := filepath.Abs(c.Input)
	if err != nil {
		return err
	}
	files, err := collectGrammarFiles(fullInput)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no .fb grammar files found in %s", fullInput)
	}

	g, err := parseAndMerge(files)
	if err != nil {
		return err
	}

	entries := make([]grammar.ParserRule, len(c.Languages))
	selectors := make([]generator.Selector, len(c.Languages))
	for i, lang := range c.Languages {
		rule, err := findEntryRule(g, lang.Entry)
		if err != nil {
			return err
		}
		entries[i] = rule
		selectors[i] = generator.Selector{
			LanguageID: lang.LanguageID,
			Patterns:   lang.Patterns,
		}
	}

	out := c.Output
	if out == "" {
		out = fullInput
	} else {
		out, err = filepath.Abs(out)
		if err != nil {
			return err
		}
	}
	pkg := c.Package
	if pkg == "" {
		pkg = filepath.Base(out)
	}
	return Generate(g, entries, selectors, out, pkg, c.ATN, c.Verbose)
}

// collectGrammarFiles returns the absolute paths of top-level *.fb files in dir.
func collectGrammarFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".fb" {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)
	return files, nil
}

// parseAndMerge parses every grammar file in a single shared container (so
// cross-file references link), aborts on any error diagnostic, and returns the
// combined grammar. A single file is returned as-is; multiple files are merged
// into one grammar and reparented so return-type resolution spans all files.
func parseAndMerge(files []string) (grammar.Grammar, error) {
	sc := grammar.CreateServices()
	documents := service.MustGet[workspace.DocumentManager](sc)

	docs := make([]*core.Document, 0, len(files))
	for _, path := range files {
		text, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		file := textdoc.NewFile(lsp.URIFromPath(path), "fb", 0, string(text))
		doc := core.NewDocument(file)
		documents.Set(doc)
		docs = append(docs, doc)
	}

	builder := service.MustGet[workspace.Builder](sc)
	if err := builder.Build(context.Background(), docs, nil); err != nil {
		return nil, err
	}
	if err := reportDiagnostics(docs); err != nil {
		return nil, err
	}

	grammars := make([]grammar.Grammar, 0, len(docs))
	for _, doc := range docs {
		g, ok := doc.Root.(grammar.Grammar)
		if !ok {
			return nil, fmt.Errorf("parse result is not a Grammar")
		}
		grammars = append(grammars, g)
	}
	if len(grammars) == 1 {
		return grammars[0], nil
	}
	return mergeGrammars(grammars)
}

// mergeGrammars combines several grammars (which must share one grammar name)
// into a single grammar and reparents every element so container-based lookups
// (e.g. grammar.FindReturnType) span all files. Duplicate rule/interface names
// across files are a hard error.
func mergeGrammars(grammars []grammar.Grammar) (grammar.Grammar, error) {
	merged := grammar.NewGrammar()
	merged.SetName(grammars[0].NameToken())

	names := map[string]string{} // element name -> kind, for duplicate detection
	claim := func(kind, n string) error {
		if prev, ok := names[n]; ok {
			return fmt.Errorf("duplicate %s name %q across grammar files (already declared as %s)", kind, n, prev)
		}
		names[n] = kind
		return nil
	}

	for _, g := range grammars {
		if g.Name() != merged.Name() {
			return nil, fmt.Errorf("all grammar files must declare the same grammar name; found %q and %q", merged.Name(), g.Name())
		}
		for _, r := range g.Rules() {
			if err := claim("rule", r.Name()); err != nil {
				return nil, err
			}
			merged.SetRulesItem(r)
		}
		for _, ci := range g.Composites() {
			if err := claim("composite", ci.Name()); err != nil {
				return nil, err
			}
			merged.SetCompositesItem(ci)
		}
		for _, t := range g.Terminals() {
			if err := claim("token", t.Name()); err != nil {
				return nil, err
			}
			merged.SetTerminalsItem(t)
		}
		for _, tg := range g.TokenGroups() {
			if err := claim("token group", tg.Name()); err != nil {
				return nil, err
			}
			merged.SetTokenGroupsItem(tg)
		}
		for _, iface := range g.Interfaces() {
			if err := claim("interface", iface.Name()); err != nil {
				return nil, err
			}
			merged.SetInterfacesItem(iface)
		}
	}

	// Reparent every appended element onto the merged grammar so that
	// container walks (return-type resolution) see the combined element set.
	file := textdoc.NewFile(lsp.URIFromPath("merged.fb"), "fb", 0, "")
	doc := core.NewDocument(file)
	doc.Root = merged
	core.AssignContainers(doc)
	return merged, nil
}

// findEntryRule resolves a language's Entry to a parser rule that exists and is
// marked entry.
func findEntryRule(g grammar.Grammar, name string) (grammar.ParserRule, error) {
	for _, rule := range g.Rules() {
		if rule.Name() == name {
			if !rule.IsEntry() {
				return nil, fmt.Errorf("entry rule %q must be marked with the 'entry' keyword", name)
			}
			return rule, nil
		}
	}
	return nil, fmt.Errorf("entry rule %q not found in grammar", name)
}

// reportDiagnostics prints sorted diagnostics for the given documents and
// returns an error when any are of error severity.
func reportDiagnostics(docs []*core.Document) error {
	var diagnostics []*core.Diagnostic
	for _, doc := range docs {
		diagnostics = append(diagnostics, doc.Diagnostics...)
	}
	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].Range.Start.Line == diagnostics[j].Range.Start.Line {
			return diagnostics[i].Range.Start.Column < diagnostics[j].Range.Start.Column
		}
		return diagnostics[i].Range.Start.Line < diagnostics[j].Range.Start.Line
	})
	errCount := 0
	for _, diag := range diagnostics {
		if diag.Severity == core.SeverityError {
			errCount++
		}
		fmt.Printf("%s - %d:%d %s\n", diag.Severity.String(),
			diag.Range.Start.Line+1, diag.Range.Start.Column+1, diag.Message)
	}
	if errCount > 0 {
		return fmt.Errorf("aborting code generation due to %d errors", errCount)
	}
	return nil
}

// Generate writes the generated Go files for grammar g into outDir. entries
// lists the entry rules (index-aligned to the configured languages); with a
// single entry the output matches the single-language CLI.
func Generate(g grammar.Grammar, entries []grammar.ParserRule, selectors []generator.Selector, outDir, pkg string, atn, verbose bool) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	write := func(name, file, content string) error {
		path := filepath.Join(outDir, file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", name, err)
		}
		if verbose {
			fmt.Printf("Written: %s\n", path)
		}
		return nil
	}

	tokenTypes := generator.GenerateTokenTypes(g)
	atnData := generator.BuildParserATNData(g, tokenTypes)

	files := []struct{ name, file, content string }{
		{"linker", "linker_gen.go", generator.GenerateLinker(g, pkg)},
		{"types", "types_gen.go", generator.GenerateTypes(g, pkg)},
		{"parser", "parser_gen.go", generator.GenerateParser(g, entries, pkg, tokenTypes, atnData)},
		{"completion-parser", "completion_parser_gen.go", generator.GenerateCompletionParser(g, entries, pkg, tokenTypes, atnData)},
		{"parser-lookahead", "parser_lookahead_gen.go", generator.GenerateParserLookahead(g, pkg, tokenTypes, atnData)},
		{"completion", "completion_gen.go", generator.GenerateCompletion(g, pkg)},
		{"lexer", "lexer_gen.go", generator.GenerateLexer(g, pkg, tokenTypes, entries)},
		{"services", "services_gen.go", generator.GenerateServices(g, selectors, entries, pkg)},
		{"atn", "atn_gen.go", generator.GenerateATN(g, pkg, tokenTypes)},
	}
	for _, f := range files {
		if err := write(f.name, f.file, f.content); err != nil {
			return err
		}
	}
	if atn {
		if err := write("atn-md", "atn.md", generator.GenerateATNMarkdown(g, pkg, tokenTypes)); err != nil {
			return err
		}
	}
	return nil
}
