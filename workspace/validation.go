// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

// CreateLexerDiagnostics creates diagnostics from lexer errors.
func CreateLexerDiagnostics(doc *core.Document) []lsp.Diagnostic {
	if len(doc.LexerErrors) == 0 {
		return []lsp.Diagnostic{}
	}

	diagnostics := make([]lsp.Diagnostic, 0, len(doc.LexerErrors))
	for _, lexErr := range doc.LexerErrors {
		diagnostics = append(diagnostics, lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      uint32(lexErr.StartLine),
					Character: uint32(lexErr.StartColumn),
				},
				End: lsp.Position{
					Line:      uint32(lexErr.EndLine),
					Character: uint32(lexErr.EndColumn),
				},
			},
			Severity: lsp.SeverityError,
			Message:  lexErr.Msg,
		})
	}
	return diagnostics
}

// CreateParserDiagnostics creates diagnostics from parser errors.
func CreateParserDiagnostics(doc *core.Document) []lsp.Diagnostic {
	if len(doc.ParserErrors) == 0 {
		return []lsp.Diagnostic{}
	}
	end := doc.TextDoc.PositionAt(len(doc.TextDoc.Content()))
	diagnostics := make([]lsp.Diagnostic, 0, len(doc.ParserErrors))
	for _, err := range doc.ParserErrors {
		token := err.Token
		if token == nil {
			eofDiagnostic := lsp.Diagnostic{
				Range: lsp.Range{
					Start: end,
					End:   end,
				},
				Severity: lsp.SeverityError,
				Message:  err.Msg,
			}
			diagnostics = append(diagnostics, eofDiagnostic)
		} else {
			tokenRange := token.Segment.Range.LspRange()
			diagnostic := lsp.Diagnostic{
				Range:    tokenRange,
				Severity: lsp.SeverityError,
				Message:  err.Msg,
			}
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	return diagnostics
}

func CreateLinkerDiagnostics(doc *core.Document) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for _, ref := range doc.References {
		err := ref.Error()
		segment := ref.Segment()
		if err != nil && segment != nil {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    segment.Range.LspRange(),
				Severity: lsp.DiagnosticSeverity(err.Severity),
				Message:  err.Msg,
			})
		}
	}
	return diagnostics
}
