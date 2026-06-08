// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// SetupServices sets up the base services for the arithmetics language.
func SetupServices(sc *service.Container) {
	service.Put[workspace.LanguageID](sc, "arithmetics")
	service.Put[workspace.FileExtensions](sc, []string{".arithmetics"})
	textdoc.SetupDefaultServices(sc)
	linking.SetupDefaultServices(sc)
	workspace.SetupDefaultServices(sc)
	SetupGeneratedServices(sc)
}

// CreateServices creates a service container for the arithmetics language to be used in the CLI and tests.
func CreateServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	sc.Seal()
	return sc
}
