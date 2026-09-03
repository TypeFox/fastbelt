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

// newCodeActionTestServer builds a sealed, fully set-up lsp.Server for
// resolve-dispatch tests.
func newCodeActionTestServer(t *testing.T, setup func(sc *service.Container)) lsp.Server {
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

// plainCodeActionProvider implements only server.CodeActionProvider, not the resolving extension.
type plainCodeActionProvider struct{}

func (plainCodeActionProvider) HandleCodeActionRequest(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	return nil, nil
}

// resolvingCodeActionProvider additionally implements server.ResolvingCodeActionProvider.
type resolvingCodeActionProvider struct{ plainCodeActionProvider }

func (resolvingCodeActionProvider) HandleCodeActionResolveRequest(ctx context.Context, action *lsp.CodeAction) (*lsp.CodeAction, error) {
	action.Title = "resolved: " + action.Title
	return action, nil
}

var _ server.CodeActionProvider = plainCodeActionProvider{}
var _ server.ResolvingCodeActionProvider = resolvingCodeActionProvider{}

func TestCodeActionResolve_AdvertisedOnlyWhenSupported(t *testing.T) {
	ctx := context.Background()

	plainServer := newCodeActionTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeActionProvider](sc, plainCodeActionProvider{})
	})
	plainInit, err := plainServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if opts, ok := plainInit.Capabilities.CodeActionProvider.(*lsp.CodeActionOptions); assert.True(t, ok, "expected CodeActionProvider capability to be advertised") {
		assert.False(t, opts.ResolveProvider, "a plain CodeActionProvider should not advertise resolve support")
	}

	resolvingServer := newCodeActionTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeActionProvider](sc, resolvingCodeActionProvider{})
	})
	resolvingInit, err := resolvingServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if opts, ok := resolvingInit.Capabilities.CodeActionProvider.(*lsp.CodeActionOptions); assert.True(t, ok, "expected CodeActionProvider capability to be advertised") {
		assert.True(t, opts.ResolveProvider, "a ResolvingCodeActionProvider should advertise resolve support")
	}
}

func TestCodeActionResolve_Dispatch(t *testing.T) {
	ctx := context.Background()
	resolvingServer := newCodeActionTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeActionProvider](sc, resolvingCodeActionProvider{})
	})

	result, err := resolvingServer.ResolveCodeAction(ctx, &lsp.CodeAction{Title: "fix"})
	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "resolved: fix", result.Title)
	}

	plainServer := newCodeActionTestServer(t, func(sc *service.Container) {
		service.Put[server.CodeActionProvider](sc, plainCodeActionProvider{})
	})
	original := &lsp.CodeAction{Title: "fix"}
	result, err = plainServer.ResolveCodeAction(ctx, original)
	assert.NoError(t, err)
	assert.Same(t, original, result, "the action should be returned unchanged when the provider doesn't support resolve")
}
