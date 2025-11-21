// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
	"typefox.dev/fastbelt/textdoc"
)

// TestLanguageServerPartialHandlers tests that the language server works with some handlers nil
func TestLanguageServerPartialHandlers(t *testing.T) {
	var completionCalled bool
	services := &ServerSrv{}
	LoadDefaultServices(services, &textdoc.TextdocSrv{})

	// Create a test completion handler
	services.LanguageServerHandlers.Completion = func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
		completionCalled = true
		return &protocol.CompletionList{
			IsIncomplete: false,
			Items: []protocol.CompletionItem{
				{
					Label: "partial-test",
					Kind:  protocol.KeywordCompletion,
				},
			},
		}, nil
	}

	server := services.LanguageServer
	ctx := context.Background()

	// Test Initialize - should use default implementation
	initResult, err := server.Initialize(ctx, &protocol.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
	}

	// Test other methods - should use default implementations (no-op)
	err = server.Initialized(ctx, &protocol.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed: %v", err)
	}

	// Test Completion - should call our handler
	completionResult, err := server.Completion(ctx, &protocol.CompletionParams{})
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
