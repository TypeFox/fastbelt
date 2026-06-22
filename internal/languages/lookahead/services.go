// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lookahead

//go:generate go run ../../../cmd/fastbelt generate ./lookahead.fb -v

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// SetupServices sets up the base services for the lookahead test language.
func SetupServices(sc *service.Container) {
	service.Put[workspace.LanguageID](sc, "lookahead")
	service.Put[workspace.FileExtensions](sc, []string{".la"})
	textdoc.SetupDefaultServices(sc)
	linking.SetupDefaultServices(sc)
	workspace.SetupDefaultServices(sc)
	SetupGeneratedServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
}

// CreateServices builds a fully-wired service container for the lookahead
// test language. The language exists solely to exercise the lookahead
// engine, so CreateServices accepts an optional CompletionContributor that
// tests pass in to override the framework default. Passing nil leaves the
// framework default in place.
func CreateServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	sc.Seal()
	// Some of the adaptive prediction decisions in this grammar need full context LL to disambiguate
	lookahead := service.MustGet[LookaheadParserLookahead](sc)
	lookahead.SetPredictionMode(parser.PredictionModeLL)
	return sc
}
