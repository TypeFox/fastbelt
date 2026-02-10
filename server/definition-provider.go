// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"github.com/TypeFox/go-lsp/protocol"
	core "typefox.dev/fastbelt"
)

type DefinitionProvider interface {
	// TODO: Maybe add the document directly to the params to avoid looking it up again in the workspace?
	// Also, maybe add a separate params struct that doesn't directly depend on the lsp lib
	// TODO: Use `LocationLink` instead of `Location` to support more advanced scenarios
	// Requires a change in the LSP library
	HandleDefinitionRequest(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error)
}

type DefaultDefinitionProvider struct {
	srv ServerSrvCont
}

func NewDefaultDefinitionProvider(srv ServerSrvCont) DefinitionProvider {
	return &DefaultDefinitionProvider{srv: srv}
}

func (dp *DefaultDefinitionProvider) HandleDefinitionRequest(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	uri := params.TextDocument.URI
	doc := dp.srv.Workspace().DocumentManager.Get(uri)
	if doc == nil {
		return nil, nil // Document not found
	}
	offset := doc.TextDoc.OffsetAt(params.Position)
	doc.RLock()
	defer doc.RUnlock()
	tokens := doc.Tokens
	sourceToken := tokens.SearchOffset(offset)
	if sourceToken == nil {
		return nil, nil // No token at the given position
	}
	ref := core.ReferenceOfToken(sourceToken)
	if ref == nil {
		return nil, nil // No reference for the token
	}
	target := ref.Description()
	if target == nil || target.NameSegment == nil {
		return nil, nil // No target description
	}
	link := protocol.Location{
		URI:   target.URI,
		Range: target.NameSegment.Range.LspRange(),
	}
	return []protocol.Location{link}, nil
}
