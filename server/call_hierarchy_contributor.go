// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

// CallHierarchyContributor provides language-specific logic for call hierarchy navigation.
//
// Usage:
//
//	type MyCallHierarchyContributor struct{ sc *service.Container }
//
//	func (c *MyCallHierarchyContributor) PrepareItem(ctx context.Context, token *core.Token, node core.AstNode) *lsp.CallHierarchyItem {
//	    if funcDecl, ok := node.(*ast.FunctionDeclaration); ok {
//	        return &lsp.CallHierarchyItem{
//	            Name: funcDecl.Name,
//	            Kind: lsp.Function,
//	            URI:  node.Document().URI.DocumentURI(),
//	            Range: node.Segment().Range.LspRange(),
//	        }
//	    }
//	    return nil
//	}
//
//	func (c *MyCallHierarchyContributor) FindIncomingCalls(ctx context.Context, item *lsp.CallHierarchyItem) []lsp.CallHierarchyIncomingCall {
//	    // Use references to find call sites
//	    return calls
//	}
type CallHierarchyContributor interface {
	// PrepareItem creates a call hierarchy item for a symbol at the given position.
	// Return nil if the symbol at the position is not callable (not a function/method).
	// The returned item will be passed to incoming/outgoing call requests.
	PrepareItem(ctx context.Context, token *core.Token, node core.AstNode) *lsp.CallHierarchyItem

	// FindIncomingCalls finds all locations that call the given symbol.
	// Each incoming call should include the call site location and the item
	// representing the caller.
	FindIncomingCalls(ctx context.Context, item *lsp.CallHierarchyItem) []lsp.CallHierarchyIncomingCall

	// FindOutgoingCalls finds all locations called by the given symbol.
	// Each outgoing call should include the call site location and the item
	// representing the callee.
	FindOutgoingCalls(ctx context.Context, item *lsp.CallHierarchyItem) []lsp.CallHierarchyOutgoingCall
}

// DefaultCallHierarchyContributor is the default implementation of [CallHierarchyContributor].
// It returns nil for all operations, effectively disabling call hierarchy.
type DefaultCallHierarchyContributor struct{}

func NewDefaultCallHierarchyContributor() CallHierarchyContributor {
	return &DefaultCallHierarchyContributor{}
}

func (c *DefaultCallHierarchyContributor) PrepareItem(ctx context.Context, token *core.Token, node core.AstNode) *lsp.CallHierarchyItem {
	return nil
}

func (c *DefaultCallHierarchyContributor) FindIncomingCalls(ctx context.Context, item *lsp.CallHierarchyItem) []lsp.CallHierarchyIncomingCall {
	return nil
}

func (c *DefaultCallHierarchyContributor) FindOutgoingCalls(ctx context.Context, item *lsp.CallHierarchyItem) []lsp.CallHierarchyOutgoingCall {
	return nil
}
