// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"encoding/json/v2"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
	utilJson "typefox.dev/fastbelt/util/json"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func TestJsonRoundtrip(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"Fastbelt Grammar language", "grammar.fb"},
		{"Completion test language", "../languages/completion/completion.fb"},
		{"Lookahead test language", "../languages/lookahead/lookahead.fb"},
		{"TokenGroups test language", "../languages/token_groups/token_groups.fb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grammmar, err := os.ReadFile(tt.path)
			require.NoError(t, err)

			exported, imported := parseExportImportGrammar(t, grammmar)
			test.AssertEqualAst(t, exported, imported) // expected := exported, actual := imported
		})
	}
}

func BenchmarkJsonRoundtrip(b *testing.B) {
	tests := []struct {
		name string
		path string
	}{
		{"Fastbelt Grammar language", "grammar.fb"},
		{"Completion test language", "../languages/completion/completion.fb"},
		{"Lookahead test language", "../languages/lookahead/lookahead.fb"},
		{"TokenGroups test language", "../languages/token_groups/token_groups.fb"},
	}
	services1 := CreateServices()
	f1 := test.New(b, services1)
	for _, tt := range tests {
		grammar, err := os.ReadFile(tt.path)
		require.NoError(b, err)

		doc1 := f1.Parse(string(grammar))
		doc1.AssertNoParseErrors()
		doc1.AssertNoLinkingErrors()
		grammar1 := doc1.Root()
		f1.Clear()
		inJson, err := json.Marshal(grammar1)
		require.NoError(b, err)

		b.Run(tt.name+"/marshal", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				if _, err := json.Marshal(grammar1); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(tt.name+"/unmarshal", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				grammar2 := NewGrammar()
				if err := json.Unmarshal(inJson, grammar2); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(tt.name+"/unmarshal-self", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				if _, err := Unmarshal[Grammar](inJson); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func parseExportImportGrammar(t testing.TB, grammar []byte) (core.AstNode, core.AstNode) {
	// Step 1: First container — parse grammar.fb from disk
	services1 := CreateServices()
	f1 := test.New(t, services1)
	defer f1.Clear()
	doc1 := f1.Parse(string(grammar))
	doc1.AssertNoParseErrors()
	doc1.AssertNoLinkingErrors()

	// Step 2: Marshal to JSON via encoding/json
	jsonData, err := json.Marshal(doc1.Root())
	require.NoError(t, err)

	// Step 3: Second container — independent services instance
	services2 := CreateServices()
	f2 := test.NewWithContext(t, services2, context.WithValue(
		context.Background(), core.JsonLinkingHelperKey(), utilJson.NewJsonLinkingHelper(
			service.MustGet[workspace.DocumentManager](services2),
		),
	))
	defer f2.Clear()

	language2 := service.MustGet[workspace.LanguageID](f2.Services())
	documents2 := service.MustGet[workspace.DocumentManager](f2.Services())
	builder2 := service.MustGet[workspace.Builder](f2.Services())

	// Step 4: Unmarshal via encoding/json and wire up the document
	grammar2 := NewGrammar()
	err = json.Unmarshal(jsonData, grammar2)
	require.NoError(t, err)

	doc2, err := core.NewDocumentFromString("inmemory:///grammar.fb", string(language2), "")
	require.NoError(t, err)
	doc2.Root = grammar2
	doc2.State = core.DocStateParsed
	core.AssignContainers(doc2)
	documents2.Set(doc2)
	f2.NewDoc(doc2, nil, nil)

	// Step 5: Run the full build pipeline on the unmarshaled document
	err = builder2.Build(f2.Ctx(), slices.Collect(documents2.All()), nil)
	require.NoError(t, err)

	return doc1.Root(), doc2.Root
}
