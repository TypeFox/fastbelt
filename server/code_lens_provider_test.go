// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// newCodeLensTestServer builds a sealed, fully set-up lsp.Server for
// resolve-dispatch tests.
func newCodeLensTestServer(t *testing.T, setup func(sc *service.Container)) lsp.Server {
	t.Helper()
	sc := service.NewContainer()
	workspace.SetupDefaultServices(sc)
	textdoc.SetupDefaultServices(sc)
	setup(sc)
	server.SetupDefaultServices(sc)
	service.Put[workspace.LanguageID](sc, "plaintext")
	service.Put[workspace.FileExtensions](sc, []string{".txt"})
	sc.Seal()
	return service.MustGet[lsp.Server](sc)
}

// plainCodeLensProvider implements only server.CodeLensProvider, not the resolving extension.
type plainCodeLensProvider struct{}

func (plainCodeLensProvider) HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	return nil, nil
}

// resolvingCodeLensProvider additionally implements server.ResolvingCodeLensProvider.
type resolvingCodeLensProvider struct{ plainCodeLensProvider }

func (resolvingCodeLensProvider) HandleCodeLensResolveRequest(ctx context.Context, lens *lsp.CodeLens) (*lsp.CodeLens, error) {
	lens.Command = &lsp.Command{Title: "resolved"}
	return lens, nil
}

var _ server.CodeLensProvider = plainCodeLensProvider{}
var _ server.ResolvingCodeLensProvider = resolvingCodeLensProvider{}

func TestCodeLensResolve_AdvertisedOnlyWhenSupported(t *testing.T) {
	ctx := context.Background()

	plainServer := newCodeLensTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeLensProvider](sc, plainCodeLensProvider{})
	})
	plainInit, err := plainServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if assert.NotNil(t, plainInit.Capabilities.CodeLensProvider) {
		assert.False(t, plainInit.Capabilities.CodeLensProvider.ResolveProvider, "a plain CodeLensProvider should not advertise resolve support")
	}

	resolvingServer := newCodeLensTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeLensProvider](sc, resolvingCodeLensProvider{})
	})
	resolvingInit, err := resolvingServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if assert.NotNil(t, resolvingInit.Capabilities.CodeLensProvider) {
		assert.True(t, resolvingInit.Capabilities.CodeLensProvider.ResolveProvider, "a ResolvingCodeLensProvider should advertise resolve support")
	}
}

func TestCodeLensResolve_Dispatch(t *testing.T) {
	ctx := context.Background()
	resolvingServer := newCodeLensTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeLensProvider](sc, resolvingCodeLensProvider{})
	})

	result, err := resolvingServer.ResolveCodeLens(ctx, &lsp.CodeLens{})
	assert.NoError(t, err)
	if assert.NotNil(t, result) && assert.NotNil(t, result.Command) {
		assert.Equal(t, "resolved", result.Command.Title)
	}

	plainServer := newCodeLensTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeLensProvider](sc, plainCodeLensProvider{})
	})
	original := &lsp.CodeLens{}
	result, err = plainServer.ResolveCodeLens(ctx, original)
	assert.NoError(t, err)
	assert.Same(t, original, result, "the lens should be returned unchanged when the provider doesn't support resolve")
}
