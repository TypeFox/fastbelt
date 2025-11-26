// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"github.com/TypeFox/go-lsp/protocol"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/textdoc"
)

// CreateLexerDiagnostics creates diagnostics from lexer errors.
func CreateLexerDiagnostics(errors []*lexer.LexerError) []protocol.Diagnostic {
	if len(errors) == 0 {
		return []protocol.Diagnostic{}
	}

	diagnostics := make([]protocol.Diagnostic, 0, len(errors))
	for _, lexErr := range errors {
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
			Source:   "fastbelt",
		})
	}
	return diagnostics
}

// CreateParserDiagnostics creates diagnostics from parser errors.
// Currently, parser errors are represented by a placeholder string.
// More detailed information will be added later.
func CreateParserDiagnostics(doc textdoc.Handle, parserError string) []protocol.Diagnostic {
	if parserError == "" {
		return []protocol.Diagnostic{}
	}

	content := doc.Content()
	endLine := uint32(0)
	endChar := uint32(0)
	if len(content) > 0 {
		endPos := doc.PositionAt(len(content))
		endLine = endPos.Line
		endChar = endPos.Character
	}

	return []protocol.Diagnostic{
		{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 0,
				},
				End: protocol.Position{
					Line:      endLine,
					Character: endChar,
				},
			},
			Severity: protocol.SeverityError,
			Message:  parserError,
			Source:   "fastbelt",
		},
	}
}
