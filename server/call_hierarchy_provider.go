// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// CallHierarchyProvider is a service for handling LSP call hierarchy requests.
//
// Usage:
//
//	type MyCallHierarchyProvider struct{ sc *service.Container }
//
//	func (p *MyCallHierarchyProvider) HandlePrepareCallHierarchyRequest(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error) {
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
//	    if funcDecl, ok := foundName.Target.Owner().(*ast.FunctionDeclaration); ok {
//	        return []lsp.CallHierarchyItem{{
//	            Name:  funcDecl.Name,
//	            Kind:  lsp.Function,
//	            URI:   doc.TextDoc.URI(),
//	            Range: funcDecl.TextRange().LspRange(doc.TextDoc),
//	        }}, nil
//	    }
//	    return nil, nil
//	}
//
//	func (p *MyCallHierarchyProvider) HandleIncomingCallsRequest(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
//	    // Use references to find call sites for params.Item
//	    return calls, nil
//	}
//
//	func (p *MyCallHierarchyProvider) HandleOutgoingCallsRequest(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
//	    // Find calls made from within params.Item
//	    return calls, nil
//	}
//
//	// Register with the service container
//	service.Put[server.CallHierarchyProvider](sc, &MyCallHierarchyProvider{sc: sc})
type CallHierarchyProvider interface {
	HandlePrepareCallHierarchyRequest(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error)
	HandleIncomingCallsRequest(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error)
	HandleOutgoingCallsRequest(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error)
}
