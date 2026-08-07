// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

func TestSignatureHelpProvider_CustomImplementation(t *testing.T) {
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)

	// Register custom provider (provider-only pattern - full implementation).
	// Use Put, not Override, since there is no default to override.
	service.Put[server.SignatureHelpProvider](sc, &grammarSignatureHelpProvider{sc: sc})
	sc.Seal()

	f := test.New(t, sc)
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string }
		Person: Name=ID;
	` + commonTokens)
	doc.AssertNoErrors()

	provider := service.MustGet[server.SignatureHelpProvider](f.Services())

	assert.Equal(t, []string{"("}, provider.TriggerCharacters())
	assert.Equal(t, []string{","}, provider.RetriggerCharacters())

	// "Person" appears twice: once as the interface name, once as the rule
	// name ("Person: Name=ID;"). Find the rule occurrence specifically by
	// matching "Person:" (only the rule declaration is followed by a colon).
	text := doc.Document.TextDoc.Text(nil)
	needle := "Person:"
	offset := 0
	for i := 0; i+len(needle) <= len(text); i++ {
		if text[i:i+len(needle)] == needle {
			offset = i
			break
		}
	}
	position := doc.Document.TextDoc.PositionAt(offset)

	result, err := provider.HandleSignatureHelpRequest(
		context.Background(),
		&lsp.SignatureHelpParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
				Position:     position,
			},
		},
	)
	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Len(t, result.Signatures, 1)
		assert.Equal(t, "rule returns Type: body", result.Signatures[0].Label)
	}
}

// grammarSignatureHelpProvider provides custom signature help for grammar
// language nodes. Provider-only pattern: adopter implements the full
// provider interface, using the shared server.NodeAtCursor helper for
// cursor-to-node resolution.
type grammarSignatureHelpProvider struct {
	sc *service.Container
}

func (p *grammarSignatureHelpProvider) TriggerCharacters() []string {
	return []string{"("}
}

func (p *grammarSignatureHelpProvider) RetriggerCharacters() []string {
	return []string{","}
}

func (p *grammarSignatureHelpProvider) HandleSignatureHelpRequest(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	node := server.NodeAtCursor(doc, params.Position)
	if _, ok := node.(*grammar.ParserRuleImpl); ok {
		return &lsp.SignatureHelp{
			Signatures: []lsp.SignatureInformation{
				{
					Label: "rule returns Type: body",
					Parameters: []lsp.ParameterInformation{
						{Label: "Type"},
						{Label: "body"},
					},
				},
			},
		}, nil
	}
	return nil, nil
}
