// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package completion

//go:generate go run ../../../cmd/fastbelt -g ./completion.fb -v

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// SetupServices sets up the base services for the completion test language.
func SetupServices(sc *service.Container) {
	service.Put[workspace.LanguageID](sc, "completion")
	service.Put[workspace.FileExtensions](sc, []string{".cmpl"})
	textdoc.SetupDefaultServices(sc)
	linking.SetupDefaultServices(sc)
	workspace.SetupDefaultServices(sc)
	service.Put[CompletionScopeProvider](sc, NewCompletionScopeProviderImpl(sc))
	SetupGeneratedServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
}

// CreateServices builds a fully-wired service container for the completion
// test language. The language exists solely to exercise the completion
// engine, so CreateServices accepts an optional CompletionContributor that
// tests pass in to override the framework default. Passing nil leaves the
// framework default in place.
func CreateServices(contributor server.CompletionContributor) *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	if contributor != nil {
		service.Override(sc, contributor)
	}
	sc.Seal()
	return sc
}
