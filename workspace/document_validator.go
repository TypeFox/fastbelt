// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// DocumentValidator validates a document's AST.
type DocumentValidator interface {
	// Validate validates the given document by traversing all nodes and calling
	// any [core.Validator] implementations found on them.
	// The level parameter identifies when validation runs (e.g. "on-type", "on-save").
	// The accept callback is used to collect diagnostics.
	Validate(ctx context.Context, doc *core.Document, level string) []*core.Diagnostic
}

// DefaultDocumentValidator is the default implementation of [DocumentValidator].
type DefaultDocumentValidator struct {
	sc *service.Container
}

func NewDefaultDocumentValidator(sc *service.Container) DocumentValidator {
	return &DefaultDocumentValidator{sc: sc}
}

func (s *DefaultDocumentValidator) Validate(ctx context.Context, doc *core.Document, level string) []*core.Diagnostic {
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
