// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
)

// CreateLexerDiagnostics converts [core.Document.LexerErrors] into
// [core.Diagnostic] values with error severity. It returns a non-nil empty
// slice when there are no lexer errors.
func CreateLexerDiagnostics(doc *core.Document, source string) []*core.Diagnostic {
	if len(doc.LexerErrors) == 0 {
		return []*core.Diagnostic{}
	}

	diagnostics := make([]*core.Diagnostic, 0, len(doc.LexerErrors))
	for _, lexErr := range doc.LexerErrors {
		diagnostics = append(diagnostics, &core.Diagnostic{
			Range:    lexErr.Range,
			Severity: core.SeverityError,
			Message:  lexErr.Msg,
			Source:   source,
		})
	}
	return diagnostics
}

// CreateParserDiagnostics converts [core.Document.ParserErrors] into
// [core.Diagnostic] values with error severity. When a parser error has no
// associated token, its range is the end of the document text.
func CreateParserDiagnostics(doc *core.Document, source string) []*core.Diagnostic {
	if len(doc.ParserErrors) == 0 {
		return []*core.Diagnostic{}
	}
	textDoc := doc.TextDoc
	endIndex := int32(len(textDoc.Content()))
	diagnostics := make([]*core.Diagnostic, 0, len(doc.ParserErrors))
	for _, err := range doc.ParserErrors {
		token := err.Token
		if token == nil || token.Type == core.EOF {
			diagnostics = append(diagnostics, &core.Diagnostic{
				Range:    core.Range{Start: endIndex, End: endIndex},
				Severity: core.SeverityError,
				Message:  err.Msg,
				Source:   source,
			})
		} else {
			diagnostics = append(diagnostics, &core.Diagnostic{
				Range:    token.Range(),
				Severity: core.SeverityError,
				Message:  err.Msg,
				Source:   source,
			})
		}
	}
	return diagnostics
}

// CreateLinkerDiagnostics converts unresolved [core.Document.References] into
// [core.Diagnostic] values. A reference contributes a diagnostic only when
// [core.UntypedReference.Error] is non-nil; severity comes from the reference error.
func CreateLinkerDiagnostics(doc *core.Document, source string) []*core.Diagnostic {
	diagnostics := []*core.Diagnostic{}
	for _, ref := range doc.References {
		err := ref.Error()
		rng := ref.Range()
		if err != nil {
			diagnostics = append(diagnostics, &core.Diagnostic{
				Range:    rng,
				Severity: core.DiagnosticSeverity(err.Severity),
				Message:  err.Msg,
				Source:   source,
			})
		}
	}
	return diagnostics
}
