// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

//go:generate go run ../../../cmd/main.go -g grammar.fb -o internal/generated

import (
	"context"
	"log"

	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
)

func main() {
	ctx := context.Background()
	services := createServices()

	if err := server.StartLanguageServer(ctx, &services.ServerSrv); err != nil {
		log.Fatalf("Failed to start language server: %v", err)
	}
}

type GrammarServices struct {
	textdoc.TextdocSrv
	server.ServerSrv
}

func createServices() *GrammarServices {
	services := &GrammarServices{}
	textdoc.LoadDefaultServices(&services.TextdocSrv)
	server.LoadDefaultServices(&services.ServerSrv, &services.TextdocSrv)
	return services
}
