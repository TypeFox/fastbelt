// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"log/slog"

	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// WorkspaceFolders is populated during the LSP initialize request.
type WorkspaceFolders struct {
	Value []lsp.WorkspaceFolder
}

// Connection is assigned by ConnectionBinder when the language server is started.
type Connection struct {
	Value *jsonrpc2.Connection
}

// SetupDefaultServices sets up the default services for the language server.
// If any service is already set, it's not overwritten.
func SetupDefaultServices(sc *service.Container) {
	service.MustPut(sc, &WorkspaceFolders{})
	service.MustPut(sc, &Connection{})
	if !service.Has[slog.Handler](sc) {
		service.MustPut(sc, NewSlogHandler(sc))
	}
	if !service.Has[jsonrpc2.Binder](sc) {
		service.MustPut(sc, NewDefaultBinder(sc))
	}
	if !service.Has[jsonrpc2.Dialer](sc) {
		service.MustPut[jsonrpc2.Dialer](sc, &StdioDialer{})
	}
	if !service.Has[lsp.Server](sc) {
		service.MustPut(sc, NewDefaultLanguageServer(sc))
	}
	if !service.Has[DocumentSyncher](sc) {
		service.MustPut(sc, NewDefaultDocumentSyncher(sc))
	}
	if !service.Has[DefinitionProvider](sc) {
		service.MustPut(sc, NewDefaultDefinitionProvider(sc))
	}
	if !service.Has[ReferencesProvider](sc) {
		service.MustPut(sc, NewDefaultReferencesProvider(sc))
	}
}
