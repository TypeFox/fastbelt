// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"log"

	"github.com/TypeFox/go-lsp/protocol"
	"github.com/TypeFox/langium-to-go/lsp"
)

func main() {
	ctx := context.Background()
	
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
	
	if err := lsp.StartLanguageServer(ctx, handlers); err != nil {
		log.Fatalf("Failed to start language server: %v", err)
	}
}
