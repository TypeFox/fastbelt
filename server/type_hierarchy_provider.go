// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// TypeHierarchyProvider is a service for handling LSP type hierarchy requests.
//
// Usage:
//
//	type MyTypeHierarchyProvider struct{ sc *service.Container }
//
//	func (p *MyTypeHierarchyProvider) HandlePrepareTypeHierarchyRequest(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error) {
//	    documentManager := service.MustGet[workspace.DocumentManager](p.sc)
//	    doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
//	    if doc == nil {
//	        return nil, nil
//	    }
//	    offset := doc.TextDoc.OffsetAt(params.Position)
//	    first, second := doc.Tokens.SearchOffset2(offset)
//	    nameFinder := service.MustGet[server.NameFinder](p.sc)
//	    foundName := nameFinder.Find(ctx, first, second)
//	    if foundName.Target == nil {
//	        return nil, nil
//	    }
//	    if classDecl, ok := foundName.Target.Owner().(*ast.ClassDeclaration); ok {
//	        return []lsp.TypeHierarchyItem{{
//	            Name:  classDecl.Name,
//	            Kind:  lsp.Class,
//	            URI:   doc.TextDoc.URI(),
//	            Range: classDecl.TextRange().LspRange(doc.TextDoc),
//	        }}, nil
//	    }
//	    return nil, nil
//	}
//
//	func (p *MyTypeHierarchyProvider) HandleSupertypesRequest(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error) {
//	    // Find base classes and implemented interfaces of params.Item
//	    return supertypes, nil
//	}
//
//	func (p *MyTypeHierarchyProvider) HandleSubtypesRequest(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error) {
//	    // Find derived classes and implementing types of params.Item
//	    return subtypes, nil
//	}
//
//	// Register with the service container
//	service.Put[server.TypeHierarchyProvider](sc, &MyTypeHierarchyProvider{sc: sc})
type TypeHierarchyProvider interface {
	HandlePrepareTypeHierarchyRequest(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error)
	HandleSupertypesRequest(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error)
	HandleSubtypesRequest(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error)
}
