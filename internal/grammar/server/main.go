// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

//go:generate go run ../../../cmd/main.go -g ../grammar.fb -o ../generated -v

import (
	"context"
	"log"

	"typefox.dev/fastbelt/internal/grammar/generated"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

func main() {
	ctx := context.Background()
	srv := createServices()

	if err := server.StartLanguageServer(ctx, srv); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

type GrammarSrv struct {
	textdoc.TextdocSrvContBlock
	workspace.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	server.ServerSrvContBlock
}

func createServices() *GrammarSrv {
	srv := &GrammarSrv{}
	textdoc.CreateDefaultServices(srv)
	generated.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	server.CreateDefaultServices(srv)
	return srv
}
