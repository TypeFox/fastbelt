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

func TestCodeLensProvider_CustomImplementation(t *testing.T) {
	// Create services and register custom provider before sealing
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)

	// Register custom provider (provider-only pattern - full implementation).
	// Use Put, not Override, since there is no default to override.
	service.Put[server.CodeLensProvider](sc, &grammarCodeLensProvider{sc: sc})
	sc.Seal()

	f := test.New(t, sc)
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string }
		Person: Name=ID;
	` + commonTokens)
	doc.AssertNoErrors()

	provider := service.MustGet[server.CodeLensProvider](f.Services())
	result, err := provider.HandleCodeLensRequest(
		context.Background(),
		&lsp.CodeLensParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, result, "Should find code lenses for grammar elements")

	// Verify we have code lenses for both the rule and interface
	var foundRuleLens, foundInterfaceLens bool
	for _, lens := range result {
		if lens.Command != nil {
			if lens.Command.Title == "Parser Rule" {
				foundRuleLens = true
				assert.Equal(t, "Person", lens.Command.Command)
			}
			if lens.Command.Title == "Interface Declaration" {
				foundInterfaceLens = true
				assert.Equal(t, "Person", lens.Command.Command)
			}
		}
	}
	assert.True(t, foundRuleLens, "Should have code lens for parser rule")
	assert.True(t, foundInterfaceLens, "Should have code lens for interface")
}

// grammarCodeLensProvider provides custom code lenses for grammar language elements.
// Provider-only pattern: adopter implements the full provider interface.
type grammarCodeLensProvider struct {
	sc *service.Container
}

func (p *grammarCodeLensProvider) HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil || doc.Root == nil {
		return []lsp.CodeLens{}, nil
	}

	var lenses []lsp.CodeLens

	// Iterate over all nodes to find rules and interfaces
	for node := range core.AllNodes(doc.Root) {
		switch n := node.(type) {
		case *grammar.ParserRuleImpl:
			if n.Name() != "" {
				lenses = append(lenses, lsp.CodeLens{
					Range: n.TextRange().LspRange(doc.TextDoc),
					Command: &lsp.Command{
						Title:   "Parser Rule",
						Command: n.Name(),
					},
				})
			}
		case *grammar.InterfaceImpl:
			if n.Name() != "" {
				lenses = append(lenses, lsp.CodeLens{
					Range: n.TextRange().LspRange(doc.TextDoc),
					Command: &lsp.Command{
						Title:   "Interface Declaration",
						Command: n.Name(),
					},
				})
			}
		}
	}

	return lenses, nil
}
