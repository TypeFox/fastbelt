// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// DeclarationProvider is a service for handling LSP declaration requests.

type DeclarationProvider interface {
	HandleDeclarationRequest(ctx context.Context, params *lsp.DeclarationParams) ([]lsp.DefinitionLink, error)
}

// DefaultDeclarationProvider is the default implementation of [DeclarationProvider].
// Language adopters should override this service if their language distinguishes
// between declaration and definition.
type DefaultDeclarationProvider struct {
	sc *service.Container
}

func NewDefaultDeclarationProvider(sc *service.Container) DeclarationProvider {
	return &DefaultDeclarationProvider{sc: sc}
}

func (p *DefaultDeclarationProvider) HandleDeclarationRequest(ctx context.Context, params *lsp.DeclarationParams) ([]lsp.DefinitionLink, error) {
	// Delegate to definition provider
	definitionProvider, err := service.Get[DefinitionProvider](p.sc)
	if err != nil {
		// No definition provider available
		return nil, nil
	}

	defParams := &lsp.DefinitionParams{
		TextDocumentPositionParams: params.TextDocumentPositionParams,
		WorkDoneProgressParams:     params.WorkDoneProgressParams,
		PartialResultParams:        params.PartialResultParams,
	}

	return definitionProvider.HandleDefinitionRequest(ctx, defParams)
}
