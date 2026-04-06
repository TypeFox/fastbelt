// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

//go:generate go run ../../cmd/fastbelt -g ./grammar.fb -v

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

type GrammarSrv struct {
	textdoc.TextdocSrvContBlock
	workspace.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	linking.LinkingSrvContBlock
	FastbeltLinkingSrvContBlock
}

func CreateServices() *GrammarSrv {
	srv := &GrammarSrv{}
	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	linking.CreateDefaultServices(srv)
	CreateDefaultServices(srv)

	srv.Workspace().LanguageID = "fastbelt"
	srv.Workspace().FileExtensions = []string{".fb"}

	// Override the default scope provider
	linkingSrv := srv.FastbeltLinking()
	linkingSrv.ScopeProvider = newScopeProviderImpl(srv)

	return srv
}
