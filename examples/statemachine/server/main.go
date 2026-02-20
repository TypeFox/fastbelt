// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"log"

	"typefox.dev/fastbelt/examples/statemachine"
	"typefox.dev/fastbelt/server"
)

type StatemachineLspSrv struct {
	*statemachine.StatemachineSrv
	server.ServerSrvContBlock
}

func main() {
	ctx := context.Background()
	statemachineSrv := statemachine.CreateServices()
	srv := &StatemachineLspSrv{StatemachineSrv: statemachineSrv}
	server.CreateDefaultServices(srv)

	if err := server.StartLanguageServer(ctx, srv); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
