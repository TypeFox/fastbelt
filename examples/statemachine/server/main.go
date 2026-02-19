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

func main() {
	ctx := context.Background()
	srv := statemachine.CreateServices()

	if err := server.StartLanguageServer(ctx, srv); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
