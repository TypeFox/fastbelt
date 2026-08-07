// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
)

// Semantic token type indices. Order must match the slice returned by
// TokenTypes(); these are the indices the LSP client uses to look up colors.
const (
	semanticTokenKeyword     = iota // grammar keywords: statemachine, state, events, ...
	semanticTokenClass              // State names, both declared and referenced
	semanticTokenEnumMember         // Event names, both declared and referenced
	semanticTokenFunction           // Command names, both declared and referenced
	semanticTokenNamespace          // the statemachine's own name
)

// Semantic token modifier indices. Order must match TokenModifiers().
const (
	semanticModifierDeclaration = iota
)

// stateMachineSemanticTokensContributor classifies statemachine tokens for
// LSP semantic highlighting.
//
// Keywords are classified uniformly via TokenType.IsKeyword(). Identifiers
// are classified by the kind of symbol they name or reference (State, Event,
// or Command). Several AST nodes carry more than one identifier field (e.g.
// a Transition has both an Event reference and a State reference, both
// parsed as plain ID tokens owned by the same Transition node), so a single
// switch on the AST node type is not enough to tell them apart. Instead,
// each identifier-bearing field's own TextRange (from the generated
// Reference[T].TextRange() accessor) is compared against the current
// token's range to determine which field this specific token belongs to.
type stateMachineSemanticTokensContributor struct{}

var _ server.SemanticTokensContributor = (*stateMachineSemanticTokensContributor)(nil)

func (c *stateMachineSemanticTokensContributor) TokenTypes() []string {
	return []string{
		semanticTokenKeyword:    "keyword",
		semanticTokenClass:      "class",
		semanticTokenEnumMember: "enumMember",
		semanticTokenFunction:   "function",
		semanticTokenNamespace:  "namespace",
	}
}

func (c *stateMachineSemanticTokensContributor) TokenModifiers() []string {
	return []string{
		semanticModifierDeclaration: "declaration",
	}
}

func (c *stateMachineSemanticTokensContributor) ClassifyToken(ctx context.Context, token *core.Token, node core.AstNode) int {
	if token == nil || token.Type == nil || node == nil {
		return -1
	}
	if token.Type.IsKeyword() {
		return semanticTokenKeyword
	}

	switch n := node.(type) {
	case *EventImpl:
		// Event only has a Name field, so any ID token owned by it is that name.
		return semanticTokenEnumMember
	case *CommandImpl:
		// Command only has a Name field, so any ID token owned by it is that name.
		return semanticTokenFunction
	case *StateImpl:
		// State has both its own Name and a list of Command references
		// (Actions). Distinguish them by comparing token ranges.
		for _, action := range n.Actions() {
			if action != nil && action.TextRange() == token.Range {
				return semanticTokenFunction
			}
		}
		return semanticTokenClass
	case *TransitionImpl:
		// Transition has an Event reference and a State reference; neither
		// is "the transition's own name" since Transition has no Name field.
		if event := n.Event(); event != nil && event.TextRange() == token.Range {
			return semanticTokenEnumMember
		}
		if state := n.State(); state != nil && state.TextRange() == token.Range {
			return semanticTokenClass
		}
		return -1
	case *StatemachineImpl:
		// Statemachine has both its own Name and an Init reference to a State.
		if init := n.Init(); init != nil && init.TextRange() == token.Range {
			return semanticTokenClass
		}
		return semanticTokenNamespace
	}

	return -1
}

func (c *stateMachineSemanticTokensContributor) GetModifiers(ctx context.Context, token *core.Token, node core.AstNode) []int {
	if token == nil || node == nil {
		return nil
	}

	switch n := node.(type) {
	case *EventImpl, *CommandImpl:
		return []int{semanticModifierDeclaration}
	case *StateImpl:
		for _, action := range n.Actions() {
			if action != nil && action.TextRange() == token.Range {
				return nil // a reference to a Command, not a declaration
			}
		}
		return []int{semanticModifierDeclaration} // the state's own name
	case *StatemachineImpl:
		if init := n.Init(); init != nil && init.TextRange() == token.Range {
			return nil // a reference to a State, not a declaration
		}
		return []int{semanticModifierDeclaration} // the statemachine's own name
	}

	return nil
}
