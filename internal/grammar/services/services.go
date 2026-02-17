// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package services

import (
	"typefox.dev/fastbelt/internal/grammar/generated"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

type GrammarSrv struct {
	textdoc.TextdocSrvContBlock
	workspace.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	linking.LinkingSrvContBlock
	server.ServerSrvContBlock
	generated.FastbeltLinkingSrvContBlock
}

func CreateServices() *GrammarSrv {
	srv := &GrammarSrv{}
	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	server.CreateDefaultServices(srv)
	linking.CreateDefaultServices(srv)
	generated.CreateDefaultServices(srv)

	// Set the file extensions to scan on workspace initialization
	workspaceIninitializer := srv.Workspace().Initializer.(*workspace.DefaultInitializer)
	workspaceIninitializer.FileExtensions = []string{".fb"}

	// Override the default scope provider
	linkingSrv := srv.FastbeltLinking()
	linkingSrv.ScopeProvider = NewFastbeltScopeProvider(srv)

	return srv
}
