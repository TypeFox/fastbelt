// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package scaffold

import (
	"fmt"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
)

// ModuleNames holds derived filenames, identifiers, and metadata for a scaffolded language module or package.
type ModuleNames struct {
	// ModulePath is the Go module path for a new module (go mod init), or the import path of the
	// language package when scaffolding into an existing module.
	ModulePath string
	// ModuleDirBase is the last path segment of ModulePath (output directory name).
	ModuleDirBase string
	// LanguageLabel is the display name supplied by the user.
	LanguageLabel string
	// LanguageID is the VS Code language identifier (lowercase slug).
	LanguageID string
	// GoPackage is the root Go package name (language implementation package).
	GoPackage string
	// ExportedBase is the PascalCase prefix for user-visible Go types (for example MyLangSrv).
	ExportedBase string
	// GrammarName is the grammar declaration name in the .fb file (ExportedBase + "Model").
	GrammarName string
	// GrammarFile is the relative path to the .fb grammar file at the module root.
	GrammarFile string
	// FileDotExt is the document file extension configured for the workspace (including a leading dot).
	FileDotExt string
	// VSCodeFileExt is the file extension without a leading dot for contributes.languages.
	VSCodeFileExt string
	// ScopeName is the TextMate scope name for syntax highlighting.
	ScopeName string
	// LSPSlug is the subdirectory name cmd/<LSPSlug> for the LSP main package.
	LSPSlug string
	// SyntaxFile is the *.tmLanguage.json basename under vscode-extension/syntaxes.
	SyntaxFile string
	// NPMExtension is the vscode-extension npm package name (alnum + "-vscode").
	NPMExtension string
	// VSIXFilename is the output .vsix file name produced by vsce package.
	VSIXFilename string
	// RootNPMName is the npm package name for the root workspace package.json.
	RootNPMName string
	// FastbeltGoGet is the library module path with version for go get (M@V).
	FastbeltGoGet string
	// FastbeltToolGoGet is the fastbelt CLI package path with version for go get -tool (M@V).
	FastbeltToolGoGet string
}

func prepareNames(modulePath, languageLabel string) (ModuleNames, error) {
	modulePath = strings.TrimSpace(modulePath)
	if modulePath == "" {
		return ModuleNames{}, fmt.Errorf("module path is empty")
	}
	base := filepath.Base(modulePath)
	if base == "." || base == string(filepath.Separator) {
		return ModuleNames{}, fmt.Errorf("invalid module path %q", modulePath)
	}
	languageLabel = strings.TrimSpace(languageLabel)
	if languageLabel == "" {
		return ModuleNames{}, fmt.Errorf("language name is empty")
	}

	languageID := deriveLanguageID(languageLabel)
	exported := exportedBase(languageLabel)
	if exported == "" {
		exported = "Lang"
	}
	grammarName := exported + "Model"
	goPkg := goPackageName(languageID)
	grammarFile := languageID + ".fb"
	alnumExt := strings.ReplaceAll(languageID, "-", "")
	if alnumExt == "" {
		alnumExt = "lang"
	}
	dotExt := "." + alnumExt
	lspSlug := languageID + "-lsp"
	syntaxFile := alnumExt + ".tmLanguage.json"
	npmExt := alnumExt + "-vscode"
	vsix := npmExt + ".vsix"

	ver := fastbeltModuleVersion()
	fastbeltGoGet := fastbeltModulePath + "@" + ver
	fastbeltToolGoGet := fastbeltModulePath + "/cmd/fastbelt@" + ver

	return ModuleNames{
		ModulePath:        modulePath,
		ModuleDirBase:     base,
		LanguageLabel:     languageLabel,
		LanguageID:        languageID,
		GoPackage:         goPkg,
		ExportedBase:      exported,
		GrammarName:       grammarName,
		GrammarFile:       grammarFile,
		FileDotExt:        dotExt,
		VSCodeFileExt:     alnumExt,
		ScopeName:         "source." + alnumExt,
		LSPSlug:           lspSlug,
		SyntaxFile:        syntaxFile,
		NPMExtension:      npmExt,
		VSIXFilename:      vsix,
		RootNPMName:       alnumExt + "-workspace",
		FastbeltGoGet:     fastbeltGoGet,
		FastbeltToolGoGet: fastbeltToolGoGet,
	}, nil
}

func deriveLanguageID(languageLabel string) string {
	var b strings.Builder
	prevSep := true
	for _, r := range strings.ToLower(languageLabel) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevSep = false
			continue
		}
		if r == '-' || r == '_' || unicode.IsSpace(r) {
			if !prevSep {
				b.WriteRune('-')
			}
			prevSep = true
		}
	}
	s := strings.Trim(b.String(), "-")
	if s == "" {
		return "language"
	}
	return s
}

func exportedBase(languageLabel string) string {
	parts := strings.FieldsFunc(languageLabel, func(r rune) bool {
		return r == '-' || r == '_' || unicode.IsSpace(r)
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		lower := strings.ToLower(p)
		r := []rune(lower)
		if len(r) == 0 {
			continue
		}
		r[0] = unicode.ToUpper(r[0])
		b.WriteString(string(r))
	}
	return b.String()
}

func goPackageName(languageID string) string {
	s := strings.ReplaceAll(languageID, "-", "")
	if s == "" {
		s = "language"
	}
	if !token.IsKeyword(s) && s != "true" && s != "false" && s != "iota" {
		return s
	}
	return s + "lang"
}
