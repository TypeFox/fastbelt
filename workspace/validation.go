// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"github.com/TypeFox/go-lsp/protocol"
	core "typefox.dev/fastbelt"
)

// CreateLexerDiagnostics creates diagnostics from lexer errors.
func CreateLexerDiagnostics(doc *core.Document) []protocol.Diagnostic {
	if len(doc.LexerErrors) == 0 {
		return []protocol.Diagnostic{}
	}

	diagnostics := make([]protocol.Diagnostic, 0, len(doc.LexerErrors))
	for _, lexErr := range doc.LexerErrors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(lexErr.StartLine),
					Character: uint32(lexErr.StartColumn),
				},
				End: protocol.Position{
					Line:      uint32(lexErr.EndLine),
					Character: uint32(lexErr.EndColumn),
				},
			},
			Severity: protocol.SeverityError,
			Message:  lexErr.Msg,
		})
	}
	return diagnostics
}

// CreateParserDiagnostics creates diagnostics from parser errors.
func CreateParserDiagnostics(doc *core.Document) []protocol.Diagnostic {
	if len(doc.ParserErrors) == 0 {
		return []protocol.Diagnostic{}
	}
	end := doc.TextDoc.PositionAt(len(doc.TextDoc.Content()))
	diagnostics := make([]protocol.Diagnostic, 0, len(doc.ParserErrors))
	for _, err := range doc.ParserErrors {
		token := err.Token
		if token == nil {
			eofDiagnostic := protocol.Diagnostic{
				Range: protocol.Range{
					Start: end,
					End:   end,
				},
				Severity: protocol.SeverityError,
				Message:  err.Msg,
			}
			diagnostics = append(diagnostics, eofDiagnostic)
		} else {
			tokenRange := token.Segment.Range.LspRange()
			diagnostic := protocol.Diagnostic{
				Range:    tokenRange,
				Severity: protocol.SeverityError,
				Message:  err.Msg,
			}
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	return diagnostics
}

func CreateLinkerDiagnostics(doc *core.Document) []protocol.Diagnostic {
	diagnostics := []protocol.Diagnostic{}
	for _, ref := range doc.References {
		err := ref.Error()
		segment := ref.Segment()
		if err != nil && segment != nil {
			diagnostics = append(diagnostics, protocol.Diagnostic{
				Range:    segment.Range.LspRange(),
				Severity: protocol.DiagnosticSeverity(err.Severity),
				Message:  err.Msg,
			})
		}
	}
	return diagnostics
}
