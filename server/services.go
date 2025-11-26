// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
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
	LanguageServerHandlers *LanguageServerHandlers
	LanguageServer         LanguageServer
	DocumentSyncher        DocumentSyncher
	// Connection is assigned by ConnectionBinder when the language server is started
	Connection       *jsonrpc2.Connection
	ConnectionBinder jsonrpc2.Binder
	ConnectionDialer jsonrpc2.Dialer
}

// CreateDefaultServices creates the default services for the language server.
// If the services are already set, they are not overwritten.
func CreateDefaultServices(c ServerSrvCont) {
	s := c.Server()
	if s.LanguageServerHandlers == nil {
		s.LanguageServerHandlers = &LanguageServerHandlers{}
	}
	if s.LanguageServer == nil {
		s.LanguageServer = NewDefaultLanguageServer(c)
	}
	if s.DocumentSyncher == nil {
		s.DocumentSyncher = NewDefaultDocumentSyncher(c)
	}
	if s.ConnectionBinder == nil {
		s.ConnectionBinder = NewDefaultBinder(c)
	}
	if s.ConnectionDialer == nil {
		s.ConnectionDialer = &StdioDialer{}
	}
}
