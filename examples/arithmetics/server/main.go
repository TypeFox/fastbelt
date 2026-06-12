// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"context"
	"log"

	lang "typefox.dev/fastbelt/examples/arithmetics"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
)

// CreateServices creates a service container for the arithmetics language to be used in the language server.
func CreateServices() *service.Container {
	sc := service.NewContainer()
	lang.SetupServices(sc)
	server.SetupDefaultServices(sc)
	server.SetupStdioServices(sc)
	sc.Seal()
	return sc
}

func main() {
	ctx := context.Background()
	sc := CreateServices()

	if err := server.StartLanguageServer(ctx, sc); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
