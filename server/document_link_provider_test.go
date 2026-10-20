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

// newDocumentLinkTestServer builds a sealed, fully set-up lsp.Server for
// resolve-dispatch tests.
func newDocumentLinkTestServer(t *testing.T, setup func(sc *service.Container)) lsp.Server {
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

// plainDocumentLinkProvider implements only server.DocumentLinkProvider, not the resolving extension.
type plainDocumentLinkProvider struct{}

func (plainDocumentLinkProvider) HandleDocumentLinkRequest(ctx context.Context, params *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error) {
	return nil, nil
}

// resolvingDocumentLinkProvider additionally implements server.ResolvingDocumentLinkProvider.
type resolvingDocumentLinkProvider struct{ plainDocumentLinkProvider }

func (resolvingDocumentLinkProvider) HandleDocumentLinkResolveRequest(ctx context.Context, link *lsp.DocumentLink) (*lsp.DocumentLink, error) {
	link.Tooltip = "resolved"
	return link, nil
}

var _ server.DocumentLinkProvider = plainDocumentLinkProvider{}
var _ server.ResolvingDocumentLinkProvider = resolvingDocumentLinkProvider{}

func TestDocumentLinkResolve_AdvertisedOnlyWhenSupported(t *testing.T) {
	ctx := context.Background()

	plainServer := newDocumentLinkTestServer(t, func(sc *service.Container) {
		service.Put[server.DocumentLinkProvider](sc, plainDocumentLinkProvider{})
	})
	plainInit, err := plainServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if assert.NotNil(t, plainInit.Capabilities.DocumentLinkProvider) {
		assert.False(t, plainInit.Capabilities.DocumentLinkProvider.ResolveProvider, "a plain DocumentLinkProvider should not advertise resolve support")
	}

	resolvingServer := newDocumentLinkTestServer(t, func(sc *service.Container) {
		service.Put[server.DocumentLinkProvider](sc, resolvingDocumentLinkProvider{})
	})
	resolvingInit, err := resolvingServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if assert.NotNil(t, resolvingInit.Capabilities.DocumentLinkProvider) {
		assert.True(t, resolvingInit.Capabilities.DocumentLinkProvider.ResolveProvider, "a ResolvingDocumentLinkProvider should advertise resolve support")
	}
}

func TestDocumentLinkResolve_Dispatch(t *testing.T) {
	ctx := context.Background()
	resolvingServer := newDocumentLinkTestServer(t, func(sc *service.Container) {
		service.Put[server.DocumentLinkProvider](sc, resolvingDocumentLinkProvider{})
	})

	result, err := resolvingServer.ResolveDocumentLink(ctx, &lsp.DocumentLink{})
	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, "resolved", result.Tooltip)
	}

	plainServer := newDocumentLinkTestServer(t, func(sc *service.Container) {
		service.Put[server.DocumentLinkProvider](sc, plainDocumentLinkProvider{})
	})
	original := &lsp.DocumentLink{}
	result, err = plainServer.ResolveDocumentLink(ctx, original)
	assert.NoError(t, err)
	assert.Same(t, original, result, "the link should be returned unchanged when the provider doesn't support resolve")
}
