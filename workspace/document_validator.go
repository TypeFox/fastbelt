// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"

	core "typefox.dev/fastbelt"
)

// DocumentValidator validates a document's AST by traversing all nodes and calling
// any [core.Validator] implementations found on them.
type DocumentValidator interface {
	Validate(ctx context.Context, doc *core.Document, level string) []*core.Diagnostic
}

// DefaultDocumentValidator is the default implementation of [DocumentValidator].
type DefaultDocumentValidator struct{}

// NewDefaultDocumentValidator creates a new default document validator.
func NewDefaultDocumentValidator() DocumentValidator {
	return &DefaultDocumentValidator{}
}

// Validate traverses the AST of the given document, calling Validate on each node
// that implements the [core.Validator] interface.
func (v *DefaultDocumentValidator) Validate(ctx context.Context, doc *core.Document, level string) []*core.Diagnostic {
	if doc.Root == nil {
		return nil
	}
	lexerErrors := CreateLexerDiagnostics(doc)
	parserErrors := CreateParserDiagnostics(doc)
	linkerErrors := CreateLinkerDiagnostics(doc)
	capacity := len(lexerErrors) + len(parserErrors) + len(linkerErrors) + 8
	diagnostics := make([]*core.Diagnostic, 0, capacity)
	diagnostics = append(diagnostics, lexerErrors...)
	diagnostics = append(diagnostics, parserErrors...)
	diagnostics = append(diagnostics, linkerErrors...)
	accept := func(d *core.Diagnostic) {
		diagnostics = append(diagnostics, d)
	}
	core.TraverseNode(doc.Root, func(node core.AstNode) {
		if ctx.Err() != nil {
			return
		}
		if validator, ok := node.(core.Validator); ok {
			validator.Validate(ctx, level, accept)
		}
	})
	return diagnostics
}
