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
	"typefox.dev/lsp"
)

const commonTokens = `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`

func TestInlayHintProvider_Default(t *testing.T) {
	sc := service.NewContainer()
	server.SetupDefaultServices(sc)
	sc.Seal()

	provider := service.MustGet[server.InlayHintProvider](sc)
	result, err := provider.HandleInlayHintRequest(context.Background(), &lsp.InlayHintParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: "file:///test.txt"},
		Range:        lsp.Range{},
	})

	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestInlayHintProvider_DefaultWithGrammar(t *testing.T) {
	// Test that the default provider returns empty hints
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)
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
	assert.Empty(t, result, "Default provider should return no hints")
}

func TestInlayHintProvider_CustomImplementation(t *testing.T) {
	// Create services and register provider with custom computer before sealing
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)
	
	// Register provider with custom computer (lightweight filter-like pattern)
	service.Override(sc, server.NewInlayHintProviderWithComputer(sc, &grammarInlayHintComputer{}))
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

// grammarInlayHintComputer provides custom inlay hints for grammar language nodes.
// This is a lightweight filter-like interface, not a separate service.
type grammarInlayHintComputer struct{}

func (c *grammarInlayHintComputer) ComputeInlayHint(ctx context.Context, node core.AstNode, accept func(lsp.InlayHint)) {
	switch n := node.(type) {
	case *grammar.ParserRuleImpl:
		if seg := n.Segment(); seg != nil && n.Name() != "" {
			accept(lsp.InlayHint{
				Position: seg.Range.End.LspPosition(),
				Label:    []lsp.InlayHintLabelPart{{Value: ": ParserRule"}},
				Kind:     lsp.Type,
			})
		}
	case *grammar.InterfaceImpl:
		if seg := n.Segment(); seg != nil && n.Name() != "" {
			accept(lsp.InlayHint{
				Position: seg.Range.End.LspPosition(),
				Label:    []lsp.InlayHintLabelPart{{Value: ": Interface"}},
				Kind:     lsp.Type,
			})
		}
	}
}
