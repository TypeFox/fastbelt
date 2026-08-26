// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package cmd is the importable build API for Fastbelt. It generates a
// parser/lexer/LSP server from one or more .fb grammar files and, when several
// languages are configured, wires them for dispatch inside a single language
// server. The fastbelt CLI (cmd/fastbelt) is a thin wrapper over this package
// for use cases that only need a single language.
package cmd

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/generator"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// Language specifies a language that can be parsed by the generated lexer/parser.
type Language struct {
	// Entry is the name of the parser rule that is used to enter the language
	Entry string
	// TokenMode is the entry token mode for that language. If not specified,
	// the generator will use the token mode called "default", or generate a custom
	// token mode based on the tokens used by the language.
	TokenMode string
	// LanguageID is the ID used for LSP documents for that language.
	LanguageID string
	// Patterns is a list of glob patterns that should be used to match documents
	// of this language on the file system.
	Patterns []string
}

type BuildContext struct {
	// Input points at a directory of .fb files, not a singular file.
	Input string
	// Output directory, defaults to Input.
	Output string
	// Package name for the generated Go code, defaults to filepath.Base(Output)
	Package string
	// Languages defines the set of languages for this build. See [Language].
	Languages []Language
	// ATN specifies whether to generate the markdown documentation for the
	// augmented transitioning network (ATN), which might help with debugging grammars.
	ATN bool
	// Railroad specifies whether to generate railroad SVG diagrams for the grammars.
	Railroad bool
	// Verbose, if true, will add additional logging information.
	Verbose bool
}

// Build performs the code generation step for the fastbelt build process.
// Any errors are reported via the error return value.
func (c *BuildContext) Build() error {
	if len(c.Languages) == 0 {
		return fmt.Errorf("no languages configured")
	}
	fullInput, err := filepath.Abs(c.Input)
	if err != nil {
		return err
	}
	g, err := LoadGrammarDir(fullInput)
	if err != nil {
		return err
	}

	entries := make([]grammar.ParserRule, len(c.Languages))
	entryModes := make([]string, len(c.Languages))
	selectors := make([]generator.Selector, len(c.Languages))
	languageIDs := map[string]bool{}
	entryNames := map[string]bool{}
	for i, lang := range c.Languages {
		// The language selector resolves open documents by language id first,
		// so an empty or shared id would route documents to the wrong language.
		if lang.LanguageID == "" {
			return fmt.Errorf("language with entry %s has no language id", lang.Entry)
		}
		if languageIDs[lang.LanguageID] {
			return fmt.Errorf("language id %q is used by more than one language", lang.LanguageID)
		}
		languageIDs[lang.LanguageID] = true
		// The parser dispatches and the lexer names its modes per entry rule,
		// so two languages cannot share one. Give one language several
		// patterns instead.
		if entryNames[lang.Entry] {
			return fmt.Errorf("entry rule %s is used by more than one language", lang.Entry)
		}
		entryNames[lang.Entry] = true
		for _, pattern := range lang.Patterns {
			if _, err := path.Match(pattern, ""); err != nil {
				return fmt.Errorf("language %q: invalid pattern %q: %w", lang.LanguageID, pattern, err)
			}
		}
		rule, err := findEntryRule(g, lang.Entry)
		if err != nil {
			return err
		}
		entries[i] = rule
		mode, err := findTokenMode(g, lang)
		if err != nil {
			return err
		}
		if mode != nil {
			// The folder-wide validation only checks that some mode lexes each
			// keyword, not that the modes this language can reach do.
			keywords, tokens := grammar.MissingLanguageTokens(g, rule, mode)
			if len(keywords) > 0 {
				return fmt.Errorf("language %q (entry %s) uses keywords that no token mode reachable from its entry mode registers: %q",
					lang.LanguageID, lang.Entry, keywords)
			}
			if len(tokens) > 0 {
				fmt.Printf("Warning: language %q (entry %s) uses tokens that no token mode reachable from its entry mode registers: %s\n",
					lang.LanguageID, lang.Entry, strings.Join(tokens, ", "))
			}
		}
		entryModes[i] = lang.TokenMode
		selectors[i] = generator.Selector{
			LanguageID: lang.LanguageID,
			Patterns:   lang.Patterns,
		}
		if len(lang.Patterns) == 0 {
			fmt.Printf("Info: language %q (entry %s) declares no patterns. "+
				"Register a custom core.LanguageSelector before SetupGeneratedServices "+
				"to ensure correct behavior of the language.\n", lang.LanguageID, lang.Entry)
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
	return Generate(g, entries, entryModes, selectors, out, pkg, c.ATN, c.Railroad, c.Verbose)
}

// LoadGrammarDir parses, links and validates every top-level .fb file in dir as
// one grammar and returns the combined result. Diagnostics are printed; any
// error diagnostic aborts with an error.
func LoadGrammarDir(dir string) (grammar.Grammar, error) {
	// Glob returns the matches sorted, so the merged grammar is deterministic.
	files, err := filepath.Glob(filepath.Join(dir, "*.fb"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .fb grammar files found in %s", dir)
	}
	return parseAndMerge(files)
}

// parseAndMerge parses every grammar file in a single shared container (so
// cross-file references link), aborts on any error diagnostic, and returns the
// combined grammar. A single file is returned as-is. Multiple files are merged
// into one grammar.
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
	if err := builder.Build(context.Background(), docs); err != nil {
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
	return mergeGrammars(grammars), nil
}

// mergeGrammars concatenates several grammars into a single one (elements
// sorted by name, see [grammar.AggregateGrammar]) and reparents every element
// so container-based lookups in the generator span all files.
func mergeGrammars(grammars []grammar.Grammar) grammar.Grammar {
	merged := grammar.AggregateGrammar(grammars)
	file := textdoc.NewFile(lsp.URIFromPath("merged.fb"), "fb", 0, "")
	doc := core.NewDocument(file)
	doc.Root = merged
	core.AssignContainers(doc)
	return merged
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

// findTokenMode resolves the language's entry token mode. An empty name means
// the default mode, which is implicit (and returned as nil) when the grammar
// declares no token modes at all and must be declared otherwise.
func findTokenMode(g grammar.Grammar, lang Language) (grammar.TokenMode, error) {
	name := cmp.Or(lang.TokenMode, "default")
	for _, mode := range g.TokenModes() {
		if mode.Name() == name {
			return mode, nil
		}
	}
	if name != "default" {
		return nil, fmt.Errorf("token mode %q not found in grammar", name)
	}
	if len(g.TokenModes()) > 0 {
		return nil, fmt.Errorf("language %q (entry %s) starts in the default token mode, "+
			"but the grammar declares no default token mode; set Language.TokenMode", lang.LanguageID, lang.Entry)
	}
	return nil, nil
}

// reportDiagnostics prints sorted diagnostics for the given documents and
// returns an error when any are of error severity.
func reportDiagnostics(docs []*core.Document) error {
	type located struct {
		file string
		pos  lsp.Position
		diag *core.Diagnostic
	}
	var diagnostics []located
	for _, doc := range docs {
		for _, diag := range doc.Diagnostics {
			diagnostics = append(diagnostics, located{
				file: doc.URI.FilePath(),
				pos:  doc.TextDoc.PositionAt(int(diag.Range.Start)),
				diag: diag,
			})
		}
	}
	slices.SortStableFunc(diagnostics, func(a, b located) int {
		return cmp.Or(strings.Compare(a.file, b.file), cmp.Compare(a.pos.Line, b.pos.Line), cmp.Compare(a.pos.Character, b.pos.Character))
	})
	errCount := 0
	for _, d := range diagnostics {
		if d.diag.Severity == core.SeverityError {
			errCount++
		}
		// For printing, convert to 1-based line and column numbers.
		fmt.Printf("%s - %s:%d:%d %s\n", d.diag.Severity.String(),
			d.file, d.pos.Line+1, d.pos.Character+1, d.diag.Message)
	}
	if errCount > 0 {
		return fmt.Errorf("aborting code generation due to %d errors", errCount)
	}
	return nil
}

// Generate writes the generated Go files for grammar g into outDir. entries
// lists the entry rules and entryModes the entry token modes (both
// index-aligned to the configured languages; a nil entryModes or empty name
// means the default mode); with a single entry the output matches the
// single-language CLI.
func Generate(g grammar.Grammar, entries []grammar.ParserRule, entryModes []string, selectors []generator.Selector, outDir, pkg string, atn, railroad, verbose bool) error {
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

	// Desugar infix rules once, so every generator below sees the synthesized
	// operator token groups and flat rule bodies.
	if err := grammar.ExpandInfixRules(g); err != nil {
		return err
	}
	tokenTypes := generator.GenerateTokenTypes(g)
	atnData := generator.BuildParserATNData(g, tokenTypes)

	files := []struct{ name, file, content string }{
		{"linker", "linker_gen.go", generator.GenerateLinker(g, pkg)},
		{"types", "types_gen.go", generator.GenerateTypes(g, pkg)},
		{"json", "json_gen.go", generator.GenerateJSON(g, pkg)},
		{"parser", "parser_gen.go", generator.GenerateParser(g, entries, pkg, tokenTypes, atnData)},
		{"completion-parser", "completion_parser_gen.go", generator.GenerateCompletionParser(g, entries, pkg, tokenTypes, atnData)},
		{"parser-lookahead", "parser_lookahead_gen.go", generator.GenerateParserLookahead(g, pkg, tokenTypes, atnData)},
		{"completion", "completion_gen.go", generator.GenerateCompletion(g, pkg)},
		{"lexer", "lexer_gen.go", generator.GenerateLexer(g, entries, entryModes, pkg, tokenTypes)},
		{"services", "services_gen.go", generator.GenerateServices(g, selectors, pkg)},
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
	if railroad {
		if err := os.MkdirAll(filepath.Join(outDir, "railroad"), 0755); err != nil {
			return err
		}
		diagrams := generator.GenerateRailroadDiagrams(g)
		for _, d := range diagrams {
			if err := write("railroad-"+d.Name, filepath.Join("railroad", d.Name+".svg"), d.SVG); err != nil {
				return err
			}
		}
		if err := write("railroad-index", filepath.Join("railroad", "index.md"),
			generator.GenerateRailroadIndexMarkdown(pkg, diagrams)); err != nil {
			return err
		}
	}
	return nil
}
