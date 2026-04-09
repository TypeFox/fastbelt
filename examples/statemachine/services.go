// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

//go:generate go run ../../cmd/fastbelt -g ./statemachine.fb -v

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
	StatemachineModelLinkingSrvContBlock
}

func CreateServices() *StatemachineSrv {
	srv := &StatemachineSrv{}

	// Default document validator: walks the AST and invokes fastbelt.Validator on each node that implements it.
	// Statemachine checks are implemented in validation.go on StatemachineImpl.
	srv.Workspace().DocumentValidator = workspace.NewDefaultDocumentValidator()

	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	linking.CreateDefaultServices(srv)
	CreateDefaultServices(srv)

	srv.Workspace().LanguageID = "statemachine"
	srv.Workspace().FileExtensions = []string{".statemachine"}

	return srv
}
