// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

// TypeHierarchyContributor provides language-specific logic for type hierarchy navigation.
//
// Usage:
//
//	type MyTypeHierarchyContributor struct{ sc *service.Container }
//
//	func (c *MyTypeHierarchyContributor) PrepareItem(ctx context.Context, token *core.Token, node core.AstNode) *lsp.TypeHierarchyItem {
//	    if classDecl, ok := node.(*ast.ClassDeclaration); ok {
//	        return &lsp.TypeHierarchyItem{
//	            Name: classDecl.Name,
//	            Kind: lsp.Class,
//	            URI:  node.Document().URI.DocumentURI(),
//	            Range: node.Segment().Range.LspRange(),
//	        }
//	    }
//	    return nil
//	}
//
//	func (c *MyTypeHierarchyContributor) FindSupertypes(ctx context.Context, item *lsp.TypeHierarchyItem) []lsp.TypeHierarchyItem {
//	    // Find base classes and implemented interfaces
//	    return supertypes
//	}
type TypeHierarchyContributor interface {
	// PrepareItem creates a type hierarchy item for a symbol at the given position.
	// Return nil if the symbol at the position is not a type (not a class/interface/struct).
	// The returned item will be passed to supertype/subtype requests.
	PrepareItem(ctx context.Context, token *core.Token, node core.AstNode) *lsp.TypeHierarchyItem

	// FindSupertypes finds all supertypes of the given type.
	// Supertypes include base classes, extended interfaces, implemented interfaces,
	// or parent types depending on the language semantics.
	FindSupertypes(ctx context.Context, item *lsp.TypeHierarchyItem) []lsp.TypeHierarchyItem

	// FindSubtypes finds all subtypes of the given type.
	// Subtypes include derived classes, implementing classes, extended interfaces,
	// or child types depending on the language semantics.
	FindSubtypes(ctx context.Context, item *lsp.TypeHierarchyItem) []lsp.TypeHierarchyItem
}

// DefaultTypeHierarchyContributor is the default implementation of [TypeHierarchyContributor].
// It returns nil for all operations, effectively disabling type hierarchy.
type DefaultTypeHierarchyContributor struct{}

func NewDefaultTypeHierarchyContributor() TypeHierarchyContributor {
	return &DefaultTypeHierarchyContributor{}
}

func (c *DefaultTypeHierarchyContributor) PrepareItem(ctx context.Context, token *core.Token, node core.AstNode) *lsp.TypeHierarchyItem {
	return nil
}

func (c *DefaultTypeHierarchyContributor) FindSupertypes(ctx context.Context, item *lsp.TypeHierarchyItem) []lsp.TypeHierarchyItem {
	return nil
}

func (c *DefaultTypeHierarchyContributor) FindSubtypes(ctx context.Context, item *lsp.TypeHierarchyItem) []lsp.TypeHierarchyItem {
	return nil
}
