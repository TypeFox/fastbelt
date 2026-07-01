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

// InlayHintComputer provides per-node logic for computing inlay hints.
// Language adopters implement this lightweight interface to emit hints for specific AST nodes.
//
// Usage:
//
//	type MyInlayHintComputer struct{}
//
//	func (c *MyInlayHintComputer) ComputeInlayHint(ctx context.Context, node core.AstNode, accept func(lsp.InlayHint)) {
//	    if funcCall, ok := node.(*ast.FunctionCall); ok {
//	        accept(lsp.InlayHint{
//	            Position: funcCall.Segment().Range.End.LspPosition(),
//	            Label:    []lsp.InlayHintLabelPart{{Value: "paramName:"}},
//	            Kind:     lsp.Parameter,
//	        })
//	    }
//	}
//
//	// Register provider with custom computer
//	service.Put(sc, server.NewInlayHintProviderWithComputer(sc, &MyInlayHintComputer{}))
type InlayHintComputer interface {
	ComputeInlayHint(ctx context.Context, node core.AstNode, accept func(lsp.InlayHint))
}

// DefaultInlayHintComputer is a no-op computer that emits no hints.
type DefaultInlayHintComputer struct{}

func (c *DefaultInlayHintComputer) ComputeInlayHint(ctx context.Context, node core.AstNode, accept func(lsp.InlayHint)) {
}

type InlayHintProvider interface {
	HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error)
}

// DefaultInlayHintProvider iterates over all AST nodes in the requested
// range and delegates to InlayHintComputer for per-node hint computation.
type DefaultInlayHintProvider struct {
	sc       *service.Container
	computer InlayHintComputer
}

// NewDefaultInlayHintProvider creates an inlay hint provider with the default computer (no hints).
func NewDefaultInlayHintProvider(sc *service.Container) InlayHintProvider {
	return &DefaultInlayHintProvider{
		sc:       sc,
		computer: &DefaultInlayHintComputer{},
	}
}

// NewInlayHintProviderWithComputer creates an inlay hint provider with a custom computer.
// The computer determines which hints to emit for each AST node.
func NewInlayHintProviderWithComputer(sc *service.Container, computer InlayHintComputer) InlayHintProvider {
	return &DefaultInlayHintProvider{
		sc:       sc,
		computer: computer,
	}
}

func (p *DefaultInlayHintProvider) HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	documentManager, err := service.Get[workspace.DocumentManager](p.sc)
	if err != nil {
		return []lsp.InlayHint{}, nil
	}

	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil || doc.Root == nil {
		return []lsp.InlayHint{}, nil
	}

	startOffset := doc.TextDoc.OffsetAt(params.Range.Start)
	endOffset := doc.TextDoc.OffsetAt(params.Range.End)

	var hints []lsp.InlayHint
	acceptor := func(hint lsp.InlayHint) {
		hints = append(hints, hint)
	}

	for node := range core.AllNodes(doc.Root) {
		seg := node.Segment()
		if seg == nil {
			continue
		}
		nodeStart := int(seg.Indices.Start)
		nodeEnd := int(seg.Indices.End)

		if nodeEnd <= startOffset || nodeStart >= endOffset {
			continue
		}

		p.computer.ComputeInlayHint(ctx, node, acceptor)
	}

	return hints, nil
}
