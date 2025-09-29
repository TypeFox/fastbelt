// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"context"
	"testing"

	"github.com/TypeFox/go-lsp/protocol"
)

// TestLanguageServerHandlers tests that the language server properly calls handlers
func TestLanguageServerHandlers(t *testing.T) {
	// Track which handlers were called
	var handlersCalled []string
	
	handlers := &LanguageServerHandlers{
		Initialize: func(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
			handlersCalled = append(handlersCalled, "initialize")
			return &protocol.InitializeResult{
				Capabilities: protocol.ServerCapabilities{
					TextDocumentSync: protocol.Incremental,
					CompletionProvider: &protocol.CompletionOptions{
						ResolveProvider: false,
					},
				},
			}, nil
		},
		Initialized: func(ctx context.Context, params *protocol.InitializedParams) error {
			handlersCalled = append(handlersCalled, "initialized")
			return nil
		},
		DidOpen: func(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
			handlersCalled = append(handlersCalled, "didOpen")
			return nil
		},
		DidChange: func(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
			handlersCalled = append(handlersCalled, "didChange")
			return nil
		},
		DidClose: func(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
			handlersCalled = append(handlersCalled, "didClose")
			return nil
		},
		Completion: func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
			handlersCalled = append(handlersCalled, "completion")
			return &protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Label: "test-completion",
						Kind:  protocol.KeywordCompletion,
					},
				},
			}, nil
		},
		Shutdown: func(ctx context.Context) error {
			handlersCalled = append(handlersCalled, "shutdown")
			return nil
		},
		Exit: func(ctx context.Context) error {
			handlersCalled = append(handlersCalled, "exit")
			return nil
		},
	}
	
	// Create language server instance
	server := &languageServer{handlers: handlers}
	ctx := context.Background()
	
	// Test Initialize
	initResult, err := server.Initialize(ctx, &protocol.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
		return
	}
	if initResult.Capabilities.TextDocumentSync != protocol.Incremental {
		t.Errorf("Expected TextDocumentSync to be Incremental, got %v", initResult.Capabilities.TextDocumentSync)
	}
	
	// Test Initialized
	err = server.Initialized(ctx, &protocol.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed: %v", err)
	}
	
	// Test DidOpen
	err = server.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{})
	if err != nil {
		t.Errorf("DidOpen failed: %v", err)
	}
	
	// Test DidChange
	err = server.DidChange(ctx, &protocol.DidChangeTextDocumentParams{})
	if err != nil {
		t.Errorf("DidChange failed: %v", err)
	}
	
	// Test DidClose
	err = server.DidClose(ctx, &protocol.DidCloseTextDocumentParams{})
	if err != nil {
		t.Errorf("DidClose failed: %v", err)
	}
	
	// Test Completion
	completionResult, err := server.Completion(ctx, &protocol.CompletionParams{})
	if err != nil {
		t.Errorf("Completion failed: %v", err)
	}
	if completionResult == nil {
		t.Error("Completion returned nil result")
		return
	}
	if len(completionResult.Items) != 1 {
		t.Errorf("Expected 1 completion item, got %d", len(completionResult.Items))
	}
	if completionResult.Items[0].Label != "test-completion" {
		t.Errorf("Expected completion label 'test-completion', got %v", completionResult.Items[0].Label)
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
	
	// Verify all handlers were called in the expected order
	expectedCalls := []string{"initialize", "initialized", "didOpen", "didChange", "didClose", "completion", "shutdown", "exit"}
	if len(handlersCalled) != len(expectedCalls) {
		t.Errorf("Expected %d handler calls, got %d", len(expectedCalls), len(handlersCalled))
	}
	
	for i, expected := range expectedCalls {
		if i >= len(handlersCalled) {
			t.Errorf("Expected handler call %d to be '%s', got <missing>", i, expected)
		} else if handlersCalled[i] != expected {
			t.Errorf("Expected handler call %d to be '%s', got '%s'", i, expected, handlersCalled[i])
		}
	}
}

// TestLanguageServerNilHandlers tests that the language server works with nil handlers
func TestLanguageServerNilHandlers(t *testing.T) {
	// Create language server instance with nil handlers
	server := &languageServer{handlers: nil}
	ctx := context.Background()
	
	// Test Initialize with nil handlers - should return default capabilities
	initResult, err := server.Initialize(ctx, &protocol.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed with nil handlers: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result with nil handlers")
		return
	}
	if initResult.Capabilities.TextDocumentSync != protocol.Incremental {
		t.Errorf("Expected default TextDocumentSync to be Incremental, got %v", initResult.Capabilities.TextDocumentSync)
	}
	
	// Test other methods with nil handlers - should not panic
	err = server.Initialized(ctx, &protocol.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed with nil handlers: %v", err)
	}
	
	err = server.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{})
	if err != nil {
		t.Errorf("DidOpen failed with nil handlers: %v", err)
	}
	
	err = server.DidChange(ctx, &protocol.DidChangeTextDocumentParams{})
	if err != nil {
		t.Errorf("DidChange failed with nil handlers: %v", err)
	}
	
	err = server.DidClose(ctx, &protocol.DidCloseTextDocumentParams{})
	if err != nil {
		t.Errorf("DidClose failed with nil handlers: %v", err)
	}
	
	// Test Completion with nil handlers - should return empty list
	completionResult, err := server.Completion(ctx, &protocol.CompletionParams{})
	if err != nil {
		t.Errorf("Completion failed with nil handlers: %v", err)
	}
	if completionResult == nil {
		t.Error("Completion returned nil result with nil handlers")
		return
	}
	if len(completionResult.Items) != 0 {
		t.Errorf("Expected 0 completion items with nil handlers, got %d", len(completionResult.Items))
	}
	
	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed with nil handlers: %v", err)
	}
	
	err = server.Exit(ctx)
	if err != nil {
		t.Errorf("Exit failed with nil handlers: %v", err)
	}
}

// TestLanguageServerPartialHandlers tests that the language server works with some handlers nil
func TestLanguageServerPartialHandlers(t *testing.T) {
	var completionCalled bool
	
	// Create handlers with only completion handler set
	handlers := &LanguageServerHandlers{
		Completion: func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
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
		},
		// All other handlers are nil
	}
	
	server := &languageServer{handlers: handlers}
	ctx := context.Background()
	
	// Test Initialize - should use default implementation
	initResult, err := server.Initialize(ctx, &protocol.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
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
	
	// Test other methods - should use default implementations (no-op)
	err = server.Initialized(ctx, &protocol.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed: %v", err)
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