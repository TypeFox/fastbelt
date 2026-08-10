// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package multilang

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func SetupServices(sc *service.Container) {
	textdoc.SetupDefaultServices(sc)
	service.Put[core.LanguageSelector](
		sc,
		core.NewDefaultLanguageSelector(
			sc,
			core.NewDocumentSelectorWithPatterns("greeting", "**/*.hello"),
			core.NewDocumentSelectorWithPatterns("farewell", "**/*.bye"),
		),
	)
	linking.SetupDefaultServices(sc)
	workspace.SetupDefaultServices(sc)
	server.SetupDefaultServices(sc)

	SetupGeneratedServices(sc)
	SetupGeneratedServerServices(sc)
}

// CreateServices creates a service container for the multilang example to be
// used in the CLI and tests.
func CreateServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	sc.Seal()
	return sc
}
