// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// DocumentValidator validates a document's AST and collects diagnostics.
type DocumentValidator interface {
	// Validate returns diagnostics for doc at the given validation level
	// (for example "on-save"). It walks doc.Root and calls [core.Validator]
	// methods on AST nodes that implement that interface.
	Validate(ctx context.Context, doc *core.Document, level string) []*core.Diagnostic
}

// DefaultDocumentValidator is the default implementation of [DocumentValidator].
type DefaultDocumentValidator struct {
	sc *service.Container
}

// NewDefaultDocumentValidator returns a [DocumentValidator] that includes
// lexer, parser, and linker diagnostics before running AST validators.
func NewDefaultDocumentValidator(sc *service.Container) DocumentValidator {
	return &DefaultDocumentValidator{sc: sc}
}

func (s *DefaultDocumentValidator) Validate(ctx context.Context, doc *core.Document, level string) []*core.Diagnostic {
	if doc.Root == nil {
		return nil
	}

	source := diagnosticSource(s.sc, doc)

	lexerErrors := CreateLexerDiagnostics(doc, source)
	parserErrors := CreateParserDiagnostics(doc, source)
	linkerErrors := CreateLinkerDiagnostics(doc, source)

	capacity := len(lexerErrors) + len(parserErrors) + len(linkerErrors) + 8
	diagnostics := make([]*core.Diagnostic, 0, capacity)
	diagnostics = append(diagnostics, lexerErrors...)
	diagnostics = append(diagnostics, parserErrors...)
	diagnostics = append(diagnostics, linkerErrors...)

	accept := func(d *core.Diagnostic) {
		if d.Source == "" {
			d.Source = source
		}
		diagnostics = append(diagnostics, d)
	}
	for node := range core.AllNodes(doc.Root) {
		if ctx.Err() != nil {
			break
		}
		if validator, ok := node.(core.Validator); ok {
			validator.Validate(ctx, level, accept)
		}
	}
	return diagnostics
}

func diagnosticSource(sc *service.Container, doc *core.Document) string {
	if doc.TextDoc != nil {
		if source := doc.TextDoc.LanguageID(); source != "" {
			return source
		}
	}
	if languageID, err := service.Get[LanguageID](sc); err == nil {
		return string(languageID)
	}
	return ""
}
