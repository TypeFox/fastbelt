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

// SetupServices sets up the base services for the multilang example. The
// LanguageSelector's entries are index-aligned with the languages configured in
// gen/main.go (0 = Greeting/*.hello, 1 = Farewell/*.bye), which is the order the
// generated parser's dispatch switch expects.
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
	SetupGeneratedServices(sc)
}

// CreateServices creates a service container for the multilang example to be
// used in the CLI and tests.
func CreateServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	sc.Seal()
	return sc
}

// CreateLspServices creates a service container for the multilang example to be
// used in the language server.
func CreateLspServices(setup func(*service.Container)) *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
	if setup != nil {
		setup(sc)
	}
	sc.Seal()
	return sc
}
