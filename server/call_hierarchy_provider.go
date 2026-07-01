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

// CallHierarchyProvider is a service for handling LSP call hierarchy requests, delegating
// the language-specific call detection to [CallHierarchyContributor].
type CallHierarchyProvider interface {
	HandlePrepareCallHierarchyRequest(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error)
	HandleIncomingCallsRequest(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error)
	HandleOutgoingCallsRequest(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error)
}

// DefaultCallHierarchyProvider is the default implementation of [CallHierarchyProvider].
// It delegates call detection to [CallHierarchyContributor].
type DefaultCallHierarchyProvider struct {
	sc *service.Container
}

func NewDefaultCallHierarchyProvider(sc *service.Container) CallHierarchyProvider {
	return &DefaultCallHierarchyProvider{sc: sc}
}

func (p *DefaultCallHierarchyProvider) HandlePrepareCallHierarchyRequest(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return nil, nil
	}

	contributor, err := service.Get[CallHierarchyContributor](p.sc)
	if err != nil {
		// No contributor available
		return nil, nil
	}

	// Find token at cursor position
	offset := doc.TextDoc.OffsetAt(params.Position)
	first, _ := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil
	}

	// Get the AST node for the token and delegate to the contributor
	node := first.Element
	item := contributor.PrepareItem(ctx, first, node)
	if item == nil {
		return nil, nil
	}

	return []lsp.CallHierarchyItem{*item}, nil
}

func (p *DefaultCallHierarchyProvider) HandleIncomingCallsRequest(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
	contributor, err := service.Get[CallHierarchyContributor](p.sc)
	if err != nil {
		// No contributor available
		return nil, nil
	}

	incomingCalls := contributor.FindIncomingCalls(ctx, &params.Item)

	return incomingCalls, nil
}

func (p *DefaultCallHierarchyProvider) HandleOutgoingCallsRequest(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
	contributor, err := service.Get[CallHierarchyContributor](p.sc)
	if err != nil {
		// No contributor available
		return nil, nil
	}

	outgoingCalls := contributor.FindOutgoingCalls(ctx, &params.Item)

	return outgoingCalls, nil
}
