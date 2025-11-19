// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"log"

	"github.com/TypeFox/go-lsp/protocol"
	"github.com/TypeFox/langium-to-go/lsp"
	"github.com/TypeFox/langium-to-go/textdoc"
)

func main() {
	ctx := context.Background()
	services := createServices()

	if err := lsp.StartLanguageServer(ctx, &services.LspServices); err != nil {
		log.Fatalf("Failed to start language server: %v", err)
	}
}

type StatemachineServices struct {
	textdoc.TextdocServices
	lsp.LspServices
	DummyServiceA string
	DummyServiceB string
}

func createServices() *StatemachineServices {
	services := &StatemachineServices{}
	textdoc.LoadDefaultServices(&services.TextdocServices)
	lsp.LoadDefaultServices(&services.LspServices, &services.TextdocServices)

	// Create a dummy completion handler for testing
	services.LanguageServerHandlers.Completion = func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
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
	}

	return services
}
