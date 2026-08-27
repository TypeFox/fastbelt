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

const commonTokens = `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`

func TestInlayHintProvider_CustomImplementation(t *testing.T) {
	// Create services and register custom provider before sealing
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)

	// Register custom provider:
	// Use Put, not Override, since there is no default to override.
	service.Put[server.InlayHintProvider](sc, &grammarInlayHintProvider{sc: sc})
	sc.Seal()

	f := test.New(t, sc)
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string }
		Person: Name=ID;
	` + commonTokens)
	doc.AssertNoErrors()

	provider := service.MustGet[server.InlayHintProvider](f.Services())
	result, err := provider.HandleInlayHintRequest(
		context.Background(),
		&lsp.InlayHintParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 10, Character: 0},
			},
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, result, "Should find inlay hints for grammar rules")

	// Verify we have hints for both the rule and interface
	var foundRuleTypeHint, foundInterfaceTypeHint bool
	for _, hint := range result {
		for _, part := range hint.Label {
			if part.Value == ": ParserRule" {
				foundRuleTypeHint = true
				assert.Equal(t, lsp.Type, hint.Kind)
			}
			if part.Value == ": Interface" {
				foundInterfaceTypeHint = true
				assert.Equal(t, lsp.Type, hint.Kind)
			}
		}
	}
	assert.True(t, foundRuleTypeHint, "Should have type hint for parser rule")
	assert.True(t, foundInterfaceTypeHint, "Should have type hint for interface")
}

func TestInlayHintProvider_RangeFiltering(t *testing.T) {
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)
	service.Put[server.InlayHintProvider](sc, &grammarInlayHintProvider{sc: sc})
	sc.Seal()

	f := test.New(t, sc)
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string }
		Person: Name=ID;
	` + commonTokens)
	doc.AssertNoErrors()

	provider := service.MustGet[server.InlayHintProvider](f.Services())
	result, err := provider.HandleInlayHintRequest(
		context.Background(),
		&lsp.InlayHintParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
			// Only cover the first line, before any rule/interface declaration.
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 1, Character: 0},
			},
		},
	)

	assert.NoError(t, err)
	assert.Empty(t, result, "server.NodesInRange should exclude declarations outside the requested range")
}

// grammarInlayHintProvider provides custom inlay hints for grammar language nodes.
type grammarInlayHintProvider struct {
	sc *service.Container
}

func (p *grammarInlayHintProvider) HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	var hints []lsp.InlayHint
	for node := range server.NodesInRange(doc, params.Range) {
		switch n := node.(type) {
		case *grammar.ParserRuleImpl:
			if n.Name() != "" {
				end := n.TextRange().LspRange(doc.TextDoc).End
				hints = append(hints, lsp.InlayHint{
					Position: end,
					Label:    []lsp.InlayHintLabelPart{{Value: ": ParserRule"}},
					Kind:     lsp.Type,
				})
			}
		case *grammar.InterfaceImpl:
			if n.Name() != "" {
				end := n.TextRange().LspRange(doc.TextDoc).End
				hints = append(hints, lsp.InlayHint{
					Position: end,
					Label:    []lsp.InlayHintLabelPart{{Value: ": Interface"}},
					Kind:     lsp.Type,
				})
			}
		}
	}
	return hints, nil
}
