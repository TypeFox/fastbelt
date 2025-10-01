// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"log"

	"github.com/TypeFox/go-lsp/protocol"
	"github.com/TypeFox/langium-to-go/inject"
	"github.com/TypeFox/langium-to-go/lsp"
)

func main() {
	ctx := context.Background()
	services := createServices()

	if err := lsp.StartLanguageServer(ctx, services); err != nil {
		log.Fatalf("Failed to start language server: %v", err)
	}
}

func createServices() *inject.ServiceContainer {
	services := inject.NewServiceContainer()

	// Create and register the language server handlers
	handlers := &lsp.LanguageServerHandlers{
		Completion: func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
			// Return dummy completions for testing
			return &protocol.CompletionList{
				IsIncomplete: false,
				Items: []protocol.CompletionItem{
					{
						Label:      "state",
						Kind:       protocol.KeywordCompletion,
						Detail:     "Define a state",
						InsertText: "state ${1:name} {\n\t$0\n}",
					},
					{
						Label:      "event",
						Kind:       protocol.KeywordCompletion,
						Detail:     "Define an event",
						InsertText: "event ${1:name}",
					},
					{
						Label:      "transition",
						Kind:       protocol.KeywordCompletion,
						Detail:     "Define a transition",
						InsertText: "${1:from} -> ${2:to} on ${3:event}",
					},
				},
			}, nil
		},
	}
	if err := inject.Register(lsp.LanguageServerHandlersKey, handlers, services); err != nil {
		log.Fatalf("Failed to register handlers: %v", err)
	}

	// Create and register the connection binder
	binder := &lsp.DefaultBinder{}
	if err := inject.Register(lsp.ConnectionBinderKey, lsp.ConnectionBinder(binder), services); err != nil {
		log.Fatalf("Failed to register connection binder: %v", err)
	}

	// Create and register the connection dialer
	dialer := &lsp.StdioDialer{}
	if err := inject.Register(lsp.ConnectionDialerKey, lsp.ConnectionDialer(dialer), services); err != nil {
		log.Fatalf("Failed to register connection dialer: %v", err)
	}

	// Create and register the language server
	server := &lsp.DefaultLanguageServer{}
	if err := inject.Register(lsp.LanguageServerKey, lsp.LanguageServer(server), services); err != nil {
		log.Fatalf("Failed to register language server: %v", err)
	}

	// Inject dependencies into all registered services
	if err := inject.InjectAll(services); err != nil {
		log.Fatalf("Failed to inject dependencies: %v", err)
	}

	return services
}
