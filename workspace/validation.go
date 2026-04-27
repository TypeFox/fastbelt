// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
)

// CreateLexerDiagnostics creates diagnostics from lexer errors.
func CreateLexerDiagnostics(doc *core.Document) []*core.Diagnostic {
	if len(doc.LexerErrors) == 0 {
		return []*core.Diagnostic{}
	}

	diagnostics := make([]*core.Diagnostic, 0, len(doc.LexerErrors))
	for _, lexErr := range doc.LexerErrors {
		diagnostics = append(diagnostics, &core.Diagnostic{
			Range: core.TextRange{
				Start: core.TextLocation{
					Line:   core.TextLine(lexErr.StartLine),
					Column: core.TextColumn(lexErr.StartColumn),
				},
				End: core.TextLocation{
					Line:   core.TextLine(lexErr.EndLine),
					Column: core.TextColumn(lexErr.EndColumn),
				},
			},
			Severity: core.SeverityError,
			Message:  lexErr.Msg,
		})
	}
	return diagnostics
}

// CreateParserDiagnostics creates diagnostics from parser errors.
func CreateParserDiagnostics(doc *core.Document) []*core.Diagnostic {
	if len(doc.ParserErrors) == 0 {
		return []*core.Diagnostic{}
	}
	lspEnd := doc.TextDoc.PositionAt(len(doc.TextDoc.Content()))
	end := core.TextLocation{
		Line:   core.TextLine(lspEnd.Line),
		Column: core.TextColumn(lspEnd.Character),
	}
	diagnostics := make([]*core.Diagnostic, 0, len(doc.ParserErrors))
	for _, err := range doc.ParserErrors {
		token := err.Token
		if token == nil {
			diagnostics = append(diagnostics, &core.Diagnostic{
				Range:    core.TextRange{Start: end, End: end},
				Severity: core.SeverityError,
				Message:  err.Msg,
			})
		} else {
			diagnostics = append(diagnostics, &core.Diagnostic{
				Range:    token.TextSegment.Range,
				Severity: core.SeverityError,
				Message:  err.Msg,
			})
		}
	}
	return diagnostics
}

// CreateLinkerDiagnostics creates diagnostics from linker errors (unresolved references).
func CreateLinkerDiagnostics(doc *core.Document) []*core.Diagnostic {
	diagnostics := []*core.Diagnostic{}
	for _, ref := range doc.References {
		err := ref.Error()
		segment := ref.Segment()
		if err != nil && segment != nil {
			diagnostics = append(diagnostics, &core.Diagnostic{
				Range:    segment.Range,
				Severity: core.DiagnosticSeverity(err.Severity),
				Message:  err.Msg,
			})
		}
	}
	return diagnostics
}
