// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
)

// SemanticTokensContributor provides language-specific logic for semantic token classification.
//
// Usage:
//
//	type MySemanticTokensContributor struct{}
//
//	const (
//	    TokenTypeClass = iota
//	    TokenTypeFunction
//	    TokenTypeVariable
//	    TokenTypeKeyword
//	)
//
//	const (
//	    TokenModifierReadonly = iota
//	    TokenModifierDeprecated
//	)
//
//	func (c *MySemanticTokensContributor) TokenTypes() []string {
//	    return []string{"class", "function", "variable", "keyword"}
//	}
//
//	func (c *MySemanticTokensContributor) TokenModifiers() []string {
//	    return []string{"readonly", "deprecated"}
//	}
//
//	func (c *MySemanticTokensContributor) ClassifyToken(ctx context.Context, token *core.Token, node core.AstNode) int {
//	    if _, ok := node.(*ast.ClassDeclaration); ok {
//	        return TokenTypeClass
//	    }
//	    return -1  // Not classified
//	}
//
//	func (c *MySemanticTokensContributor) GetModifiers(ctx context.Context, token *core.Token, node core.AstNode) []int {
//	    if varDecl, ok := node.(*ast.VarDeclaration); ok && varDecl.Const {
//	        return []int{TokenModifierReadonly}
//	    }
//	    return nil
//	}
type SemanticTokensContributor interface {
	// TokenTypes returns the semantic token types supported by this contributor.
	// The order matters: it defines the token type indices used in LSP encoding.
	//
	// Standard LSP token types include:
	//   - "namespace", "class", "enum", "interface", "struct", "typeParameter"
	//   - "type", "parameter", "variable", "property", "enumMember"
	//   - "function", "method", "macro", "keyword", "comment"
	//   - "string", "number", "regexp", "operator"
	//
	// Return an empty array to disable semantic tokens.
	TokenTypes() []string

	// TokenModifiers returns the semantic token modifiers supported by this contributor.
	// Modifiers are bit flags that can be combined.
	//
	// Standard LSP token modifiers include:
	//   - "declaration", "definition", "readonly", "static"
	//   - "deprecated", "abstract", "async", "modification"
	//   - "documentation", "defaultLibrary"
	//
	// Return an empty array if no modifiers are used.
	TokenModifiers() []string

	// ClassifyToken determines the semantic token type for a lexer token and its AST node.
	//
	// Return the token type index (position in TokenTypes() array), or -1 if
	// the token should not receive semantic highlighting.
	ClassifyToken(ctx context.Context, token *core.Token, node core.AstNode) int

	// GetModifiers determines which semantic token modifiers apply to a token.
	//
	// Return modifier indices (positions in TokenModifiers() array), or nil/empty if
	// no modifiers apply. Multiple modifiers can be returned and will be combined.
	GetModifiers(ctx context.Context, token *core.Token, node core.AstNode) []int
}

// DefaultSemanticTokensContributor is the default implementation of [SemanticTokensContributor].
// It returns empty token types and modifiers, effectively disabling semantic tokens.
type DefaultSemanticTokensContributor struct{}

func NewDefaultSemanticTokensContributor() SemanticTokensContributor {
	return &DefaultSemanticTokensContributor{}
}

func (c *DefaultSemanticTokensContributor) TokenTypes() []string {
	return []string{}
}

func (c *DefaultSemanticTokensContributor) TokenModifiers() []string {
	return []string{}
}

func (c *DefaultSemanticTokensContributor) ClassifyToken(ctx context.Context, token *core.Token, node core.AstNode) int {
	return -1
}

func (c *DefaultSemanticTokensContributor) GetModifiers(ctx context.Context, token *core.Token, node core.AstNode) []int {
	return nil
}
