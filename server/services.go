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
	service.Put(sc, &WorkspaceFolders{})
	service.Put(sc, &Connection{})
	if !service.Has[slog.Handler](sc) {
		service.Put(sc, NewSlogHandler(sc))
	}
	if !service.Has[jsonrpc2.Binder](sc) {
		service.Put(sc, NewDefaultBinder(sc))
	}
	if !service.Has[jsonrpc2.Dialer](sc) {
		service.Put[jsonrpc2.Dialer](sc, &StdioDialer{})
	}
	if !service.Has[lsp.Server](sc) {
		service.Put(sc, NewDefaultLanguageServer(sc))
	}
	if !service.Has[DocumentSyncher](sc) {
		service.Put(sc, NewDefaultDocumentSyncher(sc))
	}
	if !service.Has[DefinitionProvider](sc) {
		service.Put(sc, NewDefaultDefinitionProvider(sc))
	}
	if !service.Has[ReferencesProvider](sc) {
		service.Put(sc, NewDefaultReferencesProvider(sc))
	}
}
