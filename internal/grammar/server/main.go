// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"log"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/server"
)

type GrammarLspSrv struct {
	*grammar.GrammarSrv
	server.ServerSrvContBlock
}

func main() {
	ctx := context.Background()
	grammarSrv := grammar.CreateServices()
	srv := &GrammarLspSrv{GrammarSrv: grammarSrv}
	server.CreateDefaultServices(srv)

	if err := server.StartLanguageServer(ctx, srv); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
