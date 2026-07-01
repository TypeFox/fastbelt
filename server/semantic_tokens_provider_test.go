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

func TestSemanticTokensProvider_Default(t *testing.T) {
	sc := service.NewContainer()
	server.SetupDefaultServices(sc)
	sc.Seal()

	// Just verify the service exists
	provider := service.MustGet[server.SemanticTokensProvider](sc)
	assert.NotNil(t, provider)

	// Also verify the default contributor
	contributor := service.MustGet[server.SemanticTokensContributor](sc)
	assert.NotNil(t, contributor)
	assert.Empty(t, contributor.TokenTypes(), "Default contributor should return empty token types")
}

func TestSemanticTokensProvider_DefaultWithGrammar(t *testing.T) {
	// Test that the default contributor returns empty tokens
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

	provider := service.MustGet[server.SemanticTokensProvider](f.Services())
	result, err := provider.HandleSemanticTokensFullRequest(
		context.Background(),
		&lsp.SemanticTokensParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Data, "Default contributor should return no tokens")
}

func TestSemanticTokensProvider_CustomContributor(t *testing.T) {
	// Create services and register custom contributor before sealing
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)

	// Register custom contributor (contributor pattern - separate service)
	service.Override[server.SemanticTokensContributor](sc, &grammarSemanticTokensContributor{})
	sc.Seal()

	f := test.New(t, sc)
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string }
		Person: Name=ID;
	` + commonTokens)
	doc.AssertNoErrors()

	provider := service.MustGet[server.SemanticTokensProvider](f.Services())
	result, err := provider.HandleSemanticTokensFullRequest(
		context.Background(),
		&lsp.SemanticTokensParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
		},
	)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Data, "Should find semantic tokens for grammar elements")

	// Semantic tokens are encoded as 5-integer tuples: [deltaLine, deltaStart, length, tokenType, tokenModifiers]
	assert.Equal(t, 0, len(result.Data)%5, "Data should be 5-integer tuples")
	assert.Greater(t, len(result.Data), 0, "Should have at least one token")
}

func TestSemanticTokensProvider_TokenTypesAndModifiers(t *testing.T) {
	// Test that contributor provides token types and modifiers
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)

	contributor := &grammarSemanticTokensContributor{}
	service.Override[server.SemanticTokensContributor](sc, contributor)
	sc.Seal()

	// Verify token types and modifiers are exposed
	tokenTypes := contributor.TokenTypes()
	assert.NotEmpty(t, tokenTypes)
	assert.Contains(t, tokenTypes, "keyword")
	assert.Contains(t, tokenTypes, "type")

	tokenModifiers := contributor.TokenModifiers()
	assert.NotEmpty(t, tokenModifiers)
	assert.Contains(t, tokenModifiers, "declaration")
}

// grammarSemanticTokensContributor provides custom semantic tokens for grammar language.
type grammarSemanticTokensContributor struct{}

// Token type indices (must match TokenTypes() order)
const (
	TokenTypeKeyword = iota
	TokenTypeType
	TokenTypeProperty
)

// Token modifier indices (must match TokenModifiers() order)
const (
	TokenModifierDeclaration = iota
)

func (c *grammarSemanticTokensContributor) TokenTypes() []string {
	return []string{
		"keyword",  // 0
		"type",     // 1
		"property", // 2
	}
}

func (c *grammarSemanticTokensContributor) TokenModifiers() []string {
	return []string{
		"declaration", // 0
	}
}

func (c *grammarSemanticTokensContributor) ClassifyToken(ctx context.Context, token *core.Token, node core.AstNode) int {
	if token == nil || token.Type == nil {
		return -1
	}

	// Classify keywords
	switch token.Type.Name {
	case "grammar", "interface", "returns", "token", "hidden":
		return TokenTypeKeyword
	}

	// Classify based on AST node type
	switch n := node.(type) {
	case *grammar.InterfaceImpl:
		// Interface name is a type
		if token.Element == n {
			return TokenTypeType
		}
	case *grammar.ParserRuleImpl:
		// Rule name is a type
		if token.Element == n {
			return TokenTypeType
		}
	case *grammar.FieldImpl:
		// Field name is a property
		if token.Element == n {
			return TokenTypeProperty
		}
	}

	return -1
}

func (c *grammarSemanticTokensContributor) GetModifiers(ctx context.Context, token *core.Token, node core.AstNode) []int {
	if token == nil {
		return nil
	}

	// Mark declarations
	switch node.(type) {
	case *grammar.InterfaceImpl, *grammar.ParserRuleImpl, *grammar.FieldImpl:
		return []int{TokenModifierDeclaration}
	}

	return nil
}
