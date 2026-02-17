// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"testing"

	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

type serverSrvContTest struct {
	textdoc.TextdocSrvContBlock
	workspace.WorkspaceSrvContBlock
	workspace.GeneratedSrvContBlock
	linking.LinkingSrvContBlock
	ServerSrvContBlock
}

// TestLanguageServerBasicLifecycle tests that the language server handles the basic LSP lifecycle.
func TestLanguageServerBasicLifecycle(t *testing.T) {
	srv := &serverSrvContTest{}
	workspace.CreateDefaultServices(srv)
	CreateDefaultServices(srv)

	server := srv.Server().LanguageServer
	ctx := context.Background()

	// Test Initialize
	initResult, err := server.Initialize(ctx, &lsp.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
	}

	// Test Initialized
	err = server.Initialized(ctx, &lsp.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed: %v", err)
	}

	// Test Shutdown
	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Test Exit
	err = server.Exit(ctx)
	if err != nil {
		t.Errorf("Exit failed: %v", err)
	}
}
