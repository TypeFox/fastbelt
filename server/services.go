// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"log/slog"

	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// ServerSrvCont is an interface for service containers which include the server services.
type ServerSrvCont interface {
	textdoc.TextdocSrvCont
	workspace.WorkspaceSrvCont
	Server() *ServerSrv
}

// ServerSrvContBlock is used to define a service container satisfying ServerSrvCont.
type ServerSrvContBlock struct {
	server ServerSrv
}

func (b *ServerSrvContBlock) Server() *ServerSrv {
	return &b.server
}

// ServerSrv contains the LSP-related services for the server package.
type ServerSrv struct {
	SlogHandler        slog.Handler
	LanguageServer     LanguageServer
	DocumentSyncher    DocumentSyncher
	DefinitionProvider DefinitionProvider
	ReferencesProvider ReferencesProvider
	// WorkspaceFolders is populated during the LSP initialize request.
	WorkspaceFolders []lsp.WorkspaceFolder
	// Connection is assigned by ConnectionBinder when the language server is started.
	Connection       *jsonrpc2.Connection
	ConnectionBinder jsonrpc2.Binder
	ConnectionDialer jsonrpc2.Dialer
}

// CreateDefaultServices creates the default services for the language server.
// If the services are already set, they are not overwritten.
func CreateDefaultServices(c ServerSrvCont) {
	s := c.Server()
	if s.SlogHandler == nil {
		s.SlogHandler = NewSlogHandler(c)
	}
	if s.LanguageServer == nil {
		s.LanguageServer = NewDefaultLanguageServer(c)
	}
	if s.DocumentSyncher == nil {
		s.DocumentSyncher = NewDefaultDocumentSyncher(c)
	}
	if s.DefinitionProvider == nil {
		s.DefinitionProvider = NewDefaultDefinitionProvider(c)
	}
	if s.ReferencesProvider == nil {
		s.ReferencesProvider = NewDefaultReferencesProvider(c)
	}
	if s.ConnectionBinder == nil {
		s.ConnectionBinder = NewDefaultBinder(c)
	}
	if s.ConnectionDialer == nil {
		s.ConnectionDialer = &StdioDialer{}
	}
}
