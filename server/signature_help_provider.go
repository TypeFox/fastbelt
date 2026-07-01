// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// SignatureHelpComputer provides per-node logic for computing signature help.
// Language adopters implement this lightweight interface to provide signatures for specific AST nodes.
//
// Usage:
//
//	type MySignatureHelpComputer struct{}
//
//	func (c *MySignatureHelpComputer) TriggerCharacters() []string {
//	    return []string{"(", ","}
//	}
//
//	func (c *MySignatureHelpComputer) RetriggerCharacters() []string {
//	    return []string{","}
//	}
//
//	func (c *MySignatureHelpComputer) GetSignatureHelp(ctx context.Context, node core.AstNode) (*lsp.SignatureHelp, error) {
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
//	// Register provider with custom computer
//	service.Put(sc, server.NewSignatureHelpProviderWithComputer(sc, &MySignatureHelpComputer{}))
type SignatureHelpComputer interface {
	// TriggerCharacters returns the characters that should trigger signature help.
	// For example: []string{"(", ","} for function calls.
	TriggerCharacters() []string

	// RetriggerCharacters returns the characters that should re-trigger signature help
	// when it's already active. For example: []string{","} to update active parameter.
	RetriggerCharacters() []string

	// GetSignatureHelp is called for the AST node at the cursor position.
	// Implementations should check the node type and return signature information,
	// or nil if no signature help is available for this node.
	GetSignatureHelp(ctx context.Context, node core.AstNode) (*lsp.SignatureHelp, error)
}

// DefaultSignatureHelpComputer is a no-op computer that provides no signatures
// but returns common default trigger characters.
type DefaultSignatureHelpComputer struct{}

func (c *DefaultSignatureHelpComputer) TriggerCharacters() []string {
	return []string{"("}
}

func (c *DefaultSignatureHelpComputer) RetriggerCharacters() []string {
	return []string{","}
}

func (c *DefaultSignatureHelpComputer) GetSignatureHelp(ctx context.Context, node core.AstNode) (*lsp.SignatureHelp, error) {
	return nil, nil
}

type SignatureHelpProvider interface {
	// TriggerCharacters returns the characters that should trigger signature help.
	TriggerCharacters() []string

	// RetriggerCharacters returns the characters that should re-trigger signature help
	// when it's already active.
	RetriggerCharacters() []string

	HandleSignatureHelpRequest(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error)
}

// DefaultSignatureHelpProvider finds the AST node at the cursor position and
// delegates to SignatureHelpComputer for signature computation.
// Language adopters implement SignatureHelpComputer and pass it to
// NewSignatureHelpProviderWithComputer.
type DefaultSignatureHelpProvider struct {
	sc       *service.Container
	computer SignatureHelpComputer
}

// NewDefaultSignatureHelpProvider creates a signature help provider with the default computer (no signatures).
func NewDefaultSignatureHelpProvider(sc *service.Container) SignatureHelpProvider {
	return &DefaultSignatureHelpProvider{
		sc:       sc,
		computer: &DefaultSignatureHelpComputer{},
	}
}

// NewSignatureHelpProviderWithComputer creates a signature help provider with a custom computer.
// The computer determines which signatures to provide for each AST node and which characters trigger help.
func NewSignatureHelpProviderWithComputer(sc *service.Container, computer SignatureHelpComputer) SignatureHelpProvider {
	return &DefaultSignatureHelpProvider{
		sc:       sc,
		computer: computer,
	}
}

func (p *DefaultSignatureHelpProvider) TriggerCharacters() []string {
	return p.computer.TriggerCharacters()
}

func (p *DefaultSignatureHelpProvider) RetriggerCharacters() []string {
	return p.computer.RetriggerCharacters()
}

func (p *DefaultSignatureHelpProvider) HandleSignatureHelpRequest(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	documentManager, err := service.Get[workspace.DocumentManager](p.sc)
	if err != nil {
		return nil, nil
	}

	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	offset := doc.TextDoc.OffsetAt(params.Position)
	first, second := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil
	}

	var node core.AstNode
	if first.Element != nil {
		node = first.Element
	} else if second != nil && second.Element != nil {
		node = second.Element
	}

	if node == nil {
		return nil, nil
	}

	return p.computer.GetSignatureHelp(ctx, node)
}
