// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Command gen generates the multilang example via the programmatic build API.
// It is invoked by the //go:generate directive in the parent package.
package main

import (
	"log"

	"typefox.dev/fastbelt/cmd"
)

func main() {
	ctx := &cmd.BuildContext{
		Languages: []cmd.Language{
			{
				Entry:      "Greeting",
				TokenMode:  "GreetingMode",
				LanguageID: "greeting",
				Patterns:   []string{"**/*.hello"},
			},
			{
				Entry:      "Farewell",
				TokenMode:  "FarewellMode",
				LanguageID: "farewell",
				Patterns:   []string{"**/*.bye"},
			},
		},
		Verbose: true,
	}
	if err := ctx.Build(); err != nil {
		log.Fatalf("multilang generation failed: %v", err)
	}
}
