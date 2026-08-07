// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/lsp"
)

// SignatureHelpProvider is a service for handling LSP signature help requests.
//
// Usage:
//
//	type MySignatureHelpProvider struct{ sc *service.Container }
//
//	func (p *MySignatureHelpProvider) TriggerCharacters() []string {
//	    return []string{"(", ","}
//	}
//
//	func (p *MySignatureHelpProvider) RetriggerCharacters() []string {
//	    return []string{","}
//	}
//
//	func (p *MySignatureHelpProvider) HandleSignatureHelpRequest(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
//	    documentManager := service.MustGet[workspace.DocumentManager](p.sc)
//	    doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
//	    if doc == nil {
//	        return nil, nil
//	    }
//	    node := server.NodeAtCursor(doc, params.Position)
//	    if funcCall, ok := node.(*ast.FunctionCall); ok {
//	        return &lsp.SignatureHelp{
//	            Signatures: []lsp.SignatureInformation{
//	                {
//	                    Label: "myFunc(param1: string, param2: int)",
//	                    Parameters: []lsp.ParameterInformation{
//	                        {Label: "param1: string"},
//	                        {Label: "param2: int"},
//	                    },
//	                },
//	            },
//	        }, nil
//	    }
//	    return nil, nil
//	}
//
//	// Register with the service container
//	service.Put[server.SignatureHelpProvider](sc, &MySignatureHelpProvider{sc: sc})
type SignatureHelpProvider interface {
	// TriggerCharacters returns the characters that should trigger signature help.
	TriggerCharacters() []string

	// RetriggerCharacters returns the characters that should re-trigger signature help
	// when it's already active.
	RetriggerCharacters() []string

	HandleSignatureHelpRequest(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error)
}
