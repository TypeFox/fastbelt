// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

//go:generate go run ../../cmd/main.go -g ./statemachine.fb -v

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

type StatemachineSrv struct {
	textdoc.TextdocSrvContBlock
	workspace.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	linking.LinkingSrvContBlock
}

func CreateServices() *StatemachineSrv {
	srv := &StatemachineSrv{}
	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)

	srv.Workspace().LanguageID = "statemachine"
	srv.Workspace().FileExtensions = []string{".statemachine"}

	return srv
}
