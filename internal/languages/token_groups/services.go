// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package token_groups

//go:generate go run ../../../cmd/fastbelt generate ./token_groups.fb -v

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func SetupServices(sc *service.Container) {
	service.Put[workspace.LanguageID](sc, "completion")
	service.Put[workspace.FileExtensions](sc, []string{".cmpl"})
	textdoc.SetupDefaultServices(sc)
	linking.SetupDefaultServices(sc)
	workspace.SetupDefaultServices(sc)
	SetupGeneratedServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
}

func CreateServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	sc.Seal()
	return sc
}
