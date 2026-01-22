// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"github.com/TypeFox/go-lsp/protocol"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
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
		})
	}
	return diagnostics
}

// CreateParserDiagnostics creates diagnostics from parser errors.
func CreateParserDiagnostics(doc textdoc.Handle, parserErrors []*parser.ParserError) []protocol.Diagnostic {
	if len(parserErrors) == 0 {
		return []protocol.Diagnostic{}
	}
	diagnostics := make([]protocol.Diagnostic, 0, len(parserErrors))
	for _, err := range parserErrors {
		token := err.Token
		if token == nil {
			endPosition := doc.PositionAt(len(doc.Content()))
			eofDiagnostic := protocol.Diagnostic{
				Range: protocol.Range{
					Start: endPosition,
					End:   endPosition,
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

func CreateLinkerDiagnostics(doc textdoc.Handle, root core.AstNode) []protocol.Diagnostic {
	diagnostics := []protocol.Diagnostic{}
	core.TraverseNode(root, func(node core.AstNode) {
		node.ForEachReference(func(ur core.UntypedReference) {
			err := ur.Error()
			segment := ur.Segment()
			if err != nil && segment != nil {
				diagnostics = append(diagnostics, protocol.Diagnostic{
					Range:    segment.Range.LspRange(),
					Severity: protocol.DiagnosticSeverity(err.Severity),
					Message:  err.Msg,
				})
			}
		})
	})
	return diagnostics
}
