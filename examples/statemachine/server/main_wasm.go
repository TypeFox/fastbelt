// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

//go:build js && wasm

package main

import (
	"context"

	"typefox.dev/fastbelt/examples/statemachine"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func main() {
	ctx := context.Background()

	sc := statemachine.CreateLspServices(func(sc *service.Container) {
		// Override the default services with browser-compatible implementations.
		server.SetupWasmServices(sc)
		service.Override(sc, workspace.NewNoopInitializer())
	})

	// StartLanguageServer blocks on the connection until it is closed, which
	// keeps the WASM instance alive while the worker is running.
	if err := server.StartLanguageServer(ctx, sc); err != nil {
		// Logging goes to the worker console via the WASM runtime.
		panic(err)
	}
}
