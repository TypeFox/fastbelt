// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"encoding/json"
	"testing"

	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// TestLanguageServerBasicLifecycle tests that the language server handles the basic LSP lifecycle.
func TestLanguageServerBasicLifecycle(t *testing.T) {
	sc := service.NewContainer()
	workspace.SetupDefaultServices(sc)
	textdoc.SetupDefaultServices(sc)
	SetupDefaultServices(sc)
	service.Put[workspace.LanguageID](sc, "plaintext")
	service.Put[workspace.FileExtensions](sc, []string{".txt"})
	sc.Seal()

	server := service.MustGet[lsp.Server](sc)
	ctx := context.Background()

	// Test Initialize
	initResult, err := server.Initialize(ctx, &lsp.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
	}
	// jsonrpc2 marshals the result independently *after* the handler returns,
	// and a marshal failure there is reported via golang.org/x/exp/event
	if _, err := json.Marshal(initResult); err != nil {
		t.Errorf("Initialize result is not JSON-marshalable: %v", err)
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
