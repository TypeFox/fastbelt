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

// TestLanguageServerPartialHandlers tests that the language server works with some handlers nil
func TestLanguageServerPartialHandlers(t *testing.T) {
	var completionCalled bool
	srv := &serverSrvContTest{}
	CreateDefaultServices(srv)

	// Create a test completion handler
	srv.Server().LanguageServerHandlers.Completion = func(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
		completionCalled = true
		return &lsp.CompletionList{
			IsIncomplete: false,
			Items: []lsp.CompletionItem{
				{
					Label: "partial-test",
					Kind:  lsp.KeywordCompletion,
				},
			},
		}, nil
	}

	server := srv.Server().LanguageServer
	ctx := context.Background()

	// Test Initialize - should use default implementation
	initResult, err := server.Initialize(ctx, &lsp.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
	}

	// Test other methods - should use default implementations (no-op)
	err = server.Initialized(ctx, &lsp.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed: %v", err)
	}

	// Test Completion - should call our handler
	completionResult, err := server.Completion(ctx, &lsp.CompletionParams{})
	if err != nil {
		t.Errorf("Completion failed: %v", err)
	}
	if !completionCalled {
		t.Error("Completion handler was not called")
	}
	if completionResult.Items[0].Label != "partial-test" {
		t.Errorf("Expected completion label 'partial-test', got %v", completionResult.Items[0].Label)
	}

	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	err = server.Exit(ctx)
	if err != nil {
		t.Errorf("Exit failed: %v", err)
	}
}
