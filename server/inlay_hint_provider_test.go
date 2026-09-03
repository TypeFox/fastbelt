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

// newInlayHintTestServer builds a sealed, fully set-up lsp.Server for
// resolve-dispatch tests.
func newInlayHintTestServer(t *testing.T, setup func(sc *service.Container)) lsp.Server {
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

// plainInlayHintProvider implements only server.InlayHintProvider, not the resolving extension.
type plainInlayHintProvider struct{}

func (plainInlayHintProvider) HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	return nil, nil
}

// resolvingInlayHintProvider additionally implements server.ResolvingInlayHintProvider.
type resolvingInlayHintProvider struct{ plainInlayHintProvider }

func (resolvingInlayHintProvider) HandleInlayHintResolveRequest(ctx context.Context, hint *lsp.InlayHint) (*lsp.InlayHint, error) {
	hint.Label = []lsp.InlayHintLabelPart{{Value: "resolved"}}
	return hint, nil
}

var _ server.InlayHintProvider = plainInlayHintProvider{}
var _ server.ResolvingInlayHintProvider = resolvingInlayHintProvider{}

func TestInlayHintResolve_AdvertisedOnlyWhenSupported(t *testing.T) {
	ctx := context.Background()

	plainServer := newInlayHintTestServer(t, func(sc *service.Container) {
		service.Put[server.InlayHintProvider](sc, plainInlayHintProvider{})
	})
	plainInit, err := plainServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if wrapper, ok := plainInit.Capabilities.InlayHintProvider.(*lsp.Or_ServerCapabilities_inlayHintProvider); assert.True(t, ok, "expected InlayHintProvider capability to be advertised") {
		opts, ok := wrapper.Value.(lsp.InlayHintOptions)
		if assert.True(t, ok, "expected InlayHintOptions value") {
			assert.False(t, opts.ResolveProvider, "a plain InlayHintProvider should not advertise resolve support")
		}
	}

	resolvingServer := newInlayHintTestServer(t, func(sc *service.Container) {
		service.Put[server.InlayHintProvider](sc, resolvingInlayHintProvider{})
	})
	resolvingInit, err := resolvingServer.Initialize(ctx, &lsp.ParamInitialize{})
	assert.NoError(t, err)
	if wrapper, ok := resolvingInit.Capabilities.InlayHintProvider.(*lsp.Or_ServerCapabilities_inlayHintProvider); assert.True(t, ok, "expected InlayHintProvider capability to be advertised") {
		opts, ok := wrapper.Value.(lsp.InlayHintOptions)
		if assert.True(t, ok, "expected InlayHintOptions value") {
			assert.True(t, opts.ResolveProvider, "a ResolvingInlayHintProvider should advertise resolve support")
		}
	}
}

func TestInlayHintResolve_Dispatch(t *testing.T) {
	ctx := context.Background()
	resolvingServer := newInlayHintTestServer(t, func(sc *service.Container) {
		service.Put[server.InlayHintProvider](sc, resolvingInlayHintProvider{})
	})

	result, err := resolvingServer.Resolve(ctx, &lsp.InlayHint{})
	assert.NoError(t, err)
	if assert.NotNil(t, result) && assert.Len(t, result.Label, 1) {
		assert.Equal(t, "resolved", result.Label[0].Value)
	}

	plainServer := newInlayHintTestServer(t, func(sc *service.Container) {
		service.Put[server.InlayHintProvider](sc, plainInlayHintProvider{})
	})
	original := &lsp.InlayHint{}
	result, err = plainServer.Resolve(ctx, original)
	assert.NoError(t, err)
	assert.Same(t, original, result, "the hint should be returned unchanged when the provider doesn't support resolve")
}
