// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

//go:generate go run ../../cmd/fastbelt -g ./statemachine.fb -v

import (
	"typefox.dev/fastbelt/generated"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

type StatemachineSrv struct {
	textdoc.TextdocSrvContBlock
	generated.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	linking.LinkingSrvContBlock
	StatemachineModelLinkingSrvContBlock
}

func CreateServices() *StatemachineSrv {
	srv := &StatemachineSrv{}
	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	linking.CreateDefaultServices(srv)
	CreateDefaultServices(srv)

	srv.Workspace().LanguageID = "statemachine"
	srv.Workspace().FileExtensions = []string{".statemachine"}

	return srv
}
