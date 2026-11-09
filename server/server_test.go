// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// TestLanguageServerBasicLifecycle tests that the language server handles the basic LSP lifecycle.
func TestLanguageServerBasicLifecycle(t *testing.T) {
	sc := service.NewContainer()
	workspace.SetupDefaultServices(sc)
	textdoc.SetupDefaultServices(sc)
	SetupDefaultServices(sc)
	sc.Seal()

	server := service.MustGet[lsp.Server](sc)
	ctx := context.Background()

	// Test Initialize
	initResult, err := server.Initialize(ctx, &lsp.ParamInitialize{})
	if err != nil {
		t.Errorf("Initialize failed: %v", err)
	}
	if initResult == nil {
		t.Error("Initialize returned nil result")
	}
	// jsonrpc2 marshals the result independently *after* the handler returns,
	// and a marshal failure there is reported via golang.org/x/exp/event
	if _, err := json.Marshal(initResult); err != nil {
		t.Errorf("Initialize result is not JSON-marshalable: %v", err)
	}

	// Test Initialized
	err = server.Initialized(ctx, &lsp.InitializedParams{})
	if err != nil {
		t.Errorf("Initialized failed: %v", err)
	}

	// Test Shutdown
	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Test Exit
	err = server.Exit(ctx)
	if err != nil {
		t.Errorf("Exit failed: %v", err)
	}
}

// TestStdioServerStopsOnExit checks that the exit notification ends StartLanguageServer
// although stdin stays open.
func TestStdioServerStopsOnExit(t *testing.T) {
	stdin, client, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStdin := os.Stdin
	os.Stdin = stdin
	t.Cleanup(func() {
		os.Stdin = oldStdin
		_ = client.Close() // ends the stdin pump of the dialer
	})

	sc := service.NewContainer()
	workspace.SetupDefaultServices(sc)
	textdoc.SetupDefaultServices(sc)
	SetupDefaultServices(sc)
	SetupStdioServices(sc)
	sc.Seal()

	done := make(chan error, 1)
	go func() { done <- StartLanguageServer(context.Background(), sc) }()

	body := `{"jsonrpc":"2.0","method":"exit"}`
	if _, err := fmt.Fprintf(client, "Content-Length: %d\r\n\r\n%s", len(body), body); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("StartLanguageServer failed: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("language server did not stop after the exit notification")
	}
}
