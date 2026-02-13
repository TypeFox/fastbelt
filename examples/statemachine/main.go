// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

//go:generate go run ../../cmd/main.go -g ./statemachine.fb -o ./generated -v

import (
	"context"
	"log"

	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

func main() {
	ctx := context.Background()
	srv := createServices()

	if err := server.StartLanguageServer(ctx, srv); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type StatemachineSrv struct {
	textdoc.TextdocSrvContBlock
	workspace.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	linking.LinkingSrvContBlock
	server.ServerSrvContBlock
}

func createServices() *StatemachineSrv {
	srv := &StatemachineSrv{}
	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	server.CreateDefaultServices(srv)

	// Create a dummy completion handler for testing
	srv.Server().LanguageServerHandlers.Completion = func(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
		return &lsp.CompletionList{
			IsIncomplete: false,
			Items: []lsp.CompletionItem{
				{
					Label:      "state",
					Kind:       lsp.KeywordCompletion,
					Detail:     "Define a state",
					InsertText: "state ${1:name} {\n\t$0\n}",
				},
				{
					Label:      "event",
					Kind:       lsp.KeywordCompletion,
					Detail:     "Define an event",
					InsertText: "event ${1:name}",
				},
				{
					Label:      "transition",
					Kind:       lsp.KeywordCompletion,
					Detail:     "Define a transition",
					InsertText: "${1:from} -> ${2:to} on ${3:event}",
				},
			},
		}, nil
	}

	return srv
}
