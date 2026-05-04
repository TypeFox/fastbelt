// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

type RenameProvider interface {
	HandleRenameRequest(ctx context.Context, params *lsp.RenameParams) (*lsp.WorkspaceEdit, error)
	PrepareRenameRequest(ctx context.Context, params *lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error)
}

type DefaultRenameProvider struct {
	srv ServerSrvCont
}

func NewDefaultRenameProvider(srv ServerSrvCont) RenameProvider {
	return &DefaultRenameProvider{srv: srv}
}

func (rp *DefaultRenameProvider) HandleRenameRequest(ctx context.Context, params *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	foundName := rp.findTargetNode(ctx, &params.TextDocumentPositionParams)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, nil // Could not find a name
	}
	target := foundName.Target.Owner()
	referencesFinder := rp.srv.Server().ReferencesFinder
	workspaceEdit := &lsp.WorkspaceEdit{
		Changes: map[lsp.DocumentURI][]lsp.TextEdit{},
	}
	for refDesc := range referencesFinder.Find(ctx, target, FindReferencesOptions{
		IncludeDeclaration: true,
	}) {
		edit := lsp.TextEdit{
			Range:   refDesc.Segment.Range.LspRange(),
			NewText: params.NewName,
		}
		uri := refDesc.SourceURI().DocumentURI()
		workspaceEdit.Changes[uri] = append(workspaceEdit.Changes[uri], edit)
	}
	return workspaceEdit, nil
}

func (rp *DefaultRenameProvider) PrepareRenameRequest(ctx context.Context, params *lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error) {
	foundName := rp.findTargetNode(ctx, &params.TextDocumentPositionParams)
	if foundName.Target == nil || foundName.Source == nil {
		return nil, nil // Could not find a name
	}
	target := foundName.Target
	targetRange := target.Segment().Range.LspRange()
	return &lsp.PrepareRenameResult{
		Range: targetRange,
	}, nil
}

func (rp *DefaultRenameProvider) findTargetNode(ctx context.Context, params *lsp.TextDocumentPositionParams) FoundName {
	uri := core.ParseURI(string(params.TextDocument.URI))
	targetDoc := rp.srv.Workspace().DocumentManager.Get(uri)
	if targetDoc == nil {
		return FoundName{} // Document not found
	}
	offset := targetDoc.TextDoc.OffsetAt(params.Position)
	tokens := targetDoc.Tokens
	sourceToken := tokens.SearchOffset(offset)
	if sourceToken == nil {
		return FoundName{} // No token at the given position
	}
	nameFinder := rp.srv.Server().NameFinder
	foundName := nameFinder.Find(ctx, sourceToken)
	return foundName
}
