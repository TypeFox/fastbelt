// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

// parseGrammars parses the given (uri, source) pairs in one shared grammar
// container and returns the resulting Grammar roots in order.
func parseGrammars(t *testing.T, uriSourcePairs ...string) []grammar.Grammar {
	t.Helper()
	docs := test.New(t, grammar.CreateServices()).ParseAll(uriSourcePairs...)
	grammars := make([]grammar.Grammar, 0, len(docs))
	for _, doc := range docs {
		g, ok := doc.Document.Root.(grammar.Grammar)
		require.True(t, ok, "%s: root is %T, want Grammar", doc.Document.URI, doc.Document.Root)
		grammars = append(grammars, g)
	}
	return grammars
}

// A rule and the interface that declares its return type conventionally share
// a name; the merge must not report that as a cross-file duplicate.
func TestMergeGrammarsAllowsInterfaceNamedLikeRule(t *testing.T) {
	grammars := parseGrammars(t,
		"file:///ws/a.fb", `
			grammar Multi
			interface Greeting { Name string }
			entry Greeting: "hello" Name=ID
			token ID: /[_a-zA-Z][\w_]*/
			hidden token WS: /\s+/
		`,
		"file:///ws/b.fb", `
			grammar Multi
			interface Farewell { Name string }
			entry Farewell: "goodbye" Name=ID
		`,
	)

	merged, err := mergeGrammars(grammars)
	require.NoError(t, err)
	require.Len(t, merged.Rules(), 2)
	require.Len(t, merged.Interfaces(), 2)
	require.Len(t, merged.Terminals(), 2)
}

func TestMergeGrammarsRejectsDuplicateRuleAcrossFiles(t *testing.T) {
	grammars := parseGrammars(t,
		"file:///ws/a.fb", `
			grammar Multi
			interface Greeting { Name string }
			entry Greeting: "hello" Name=ID
			token ID: /[_a-zA-Z][\w_]*/
		`,
		"file:///ws/b.fb", `
			grammar Multi
			Greeting: "hi" Name=ID
		`,
	)

	_, err := mergeGrammars(grammars)
	require.ErrorContains(t, err, `duplicate rule name "Greeting"`)
}

func TestMergeGrammarsRejectsDuplicateInterfaceAcrossFiles(t *testing.T) {
	grammars := parseGrammars(t,
		"file:///ws/a.fb", `
			grammar Multi
			interface Greeting { Name string }
			entry Greeting: "hello" Name=ID
			token ID: /[_a-zA-Z][\w_]*/
		`,
		"file:///ws/b.fb", `
			grammar Multi
			interface Greeting { Other string }
		`,
	)

	_, err := mergeGrammars(grammars)
	require.ErrorContains(t, err, `duplicate interface name "Greeting"`)
}
