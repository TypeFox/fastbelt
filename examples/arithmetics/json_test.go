// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func TestJsonExport(t *testing.T) {
	services := CreateServices()
	f := test.NewWithContext(t, services, context.WithValue(
		context.Background(), core.JsonLinkingHelperKey(), util.NewJsonLinkingHelper(
			service.MustGet[workspace.DocumentManager](services),
		),
	))

	cleanup := func() {
		f.Clear()
	}

	language := service.MustGet[workspace.LanguageID](f.Services())
	documents := service.MustGet[workspace.DocumentManager](f.Services())
	builder := service.MustGet[workspace.Builder](f.Services())

	t.Run("selfContained", func(t *testing.T) {
		t.Cleanup(cleanup)

		doc := f.Parse(`
			module test
			def two: 2;
			def root(x, y): x^(1/y);
			def sqrt(x): root(x, two);
			sqrt(16);
		`)
		doc.AssertState(core.DocStateLinked)
		doc.AssertNoParseErrors()
		doc.AssertNoLinkingErrors()

		res, err := json.Marshal(doc.Root(), jsontext.WithIndent("	"))
		require.NoError(t, err)

		assert.Equal(t, selfContainedJson, string(res))
	})

	t.Run("circular-refs", func(t *testing.T) {
		t.Cleanup(cleanup)

		docs := f.ParseAll(
			"inmemory:///testA", `
				module A
				def root(x, y): two * x^(1/y);
			`,
			"inmemory:///testB", `
				module B
				def two: 2;
				def sqrt(x): root(x, two);
				sqrt(16);
			`,
		)
		for _, doc := range docs {
			doc.AssertState(core.DocStateLinked)
			doc.AssertNoParseErrors()
			doc.AssertNoLinkingErrors()
		}

		for _, doc := range docs {
			res, err := json.Marshal(doc.Root(), jsontext.WithIndent("	"))
			require.NoError(t, err)

			switch path := doc.Document.URI.Path(); path {
			case "/testA":
				assert.Equal(t, circularAJson, string(res))
			case "/testB":
				assert.Equal(t, circularBJson, string(res))
			default:
				assert.Failf(t, "Unexpected document %s", path)
			}
		}
	})

	loadJsonDoc := func(uri string, jsonInput string) {
		doc, err := core.NewDocumentFromString(uri, string(language), "")
		require.NoError(t, err)
		mod, err := util.UnmarshalDecode[Module](
			jsontext.NewDecoder(bytes.NewReader([]byte(jsonInput))),
			ArithmeticsSyntheticFactories,
		)
		require.NoError(t, err)
		doc.Root = mod
		doc.State = core.DocStateParsed
		core.AssignContainers(doc)
		documents.Set(doc)
		f.NewDoc(doc, nil, nil)
	}

	t.Run("circular-refs-json-text", func(t *testing.T) {
		t.Cleanup(cleanup)

		loadJsonDoc("inmemory:///testA.json", circularAJson)

		f.ParseURI(`
				module B
				def two: 2;
				def sqrt(x): root(x, two);
				sqrt(16);
			`, "inmemory:///testB",
		)

		for doc := range documents.All() {
			builder.Reset(doc, core.DocStateParsed)
		}
		err := builder.Build(
			f.Ctx(), slices.Collect(documents.All()), nil,
		)
		require.NoError(t, err)

		for _, doc := range f.Documents() {
			doc.AssertState(core.DocStateParsed)
			doc.AssertState(core.DocStateExportedSymbols)
			doc.AssertState(core.DocStateImportedSymbols)
			doc.AssertState(core.DocStateLinked)
			doc.AssertNoParseErrors()
			doc.AssertNoLinkingErrors()
		}

		for doc := range documents.All() {
			res, err := json.Marshal(doc.Root, jsontext.WithIndent("	"))
			require.NoError(t, err)

			switch path := doc.URI.Path(); path {
			case "/testA.json":
				assert.Equal(t, circularAJson, string(res))
			case "/testB":
				assert.Equal(t, strings.ReplaceAll(circularBJson, "inmemory:/testA#", "inmemory:/testA.json#"), string(res))
			default:
				assert.Failf(t, "Unexpected document %s", path)
			}
		}
	})

	t.Run("circular-refs-json-json", func(t *testing.T) {
		t.Cleanup(cleanup)

		circularAJsonAdj := strings.ReplaceAll(circularAJson, "inmemory:/testB#", "inmemory:/testB.json#")
		loadJsonDoc("inmemory:///testA.json", circularAJsonAdj)

		circularBJsonAdj := strings.ReplaceAll(circularBJson, "inmemory:/testA#", "inmemory:/testA.json#")
		loadJsonDoc("inmemory:///testB.json", circularBJsonAdj)

		err := builder.Build(
			f.Ctx(), slices.Collect(documents.All()), nil,
		)
		require.NoError(t, err)

		for _, doc := range f.Documents() {
			doc.AssertState(core.DocStateParsed)
			doc.AssertState(core.DocStateExportedSymbols)
			doc.AssertState(core.DocStateImportedSymbols)
			doc.AssertState(core.DocStateLinked)
			doc.AssertNoParseErrors()
			doc.AssertNoLinkingErrors()
		}

		for doc := range documents.All() {
			res, err := json.Marshal(doc.Root, jsontext.WithIndent("	"))
			require.NoError(t, err)

			switch path := doc.URI.Path(); path {
			case "/testA.json":
				assert.Equal(t, circularAJsonAdj, string(res))
			case "/testB.json":
				assert.Equal(t, circularBJsonAdj, string(res))
			default:
				assert.Failf(t, "Unexpected document %s", path)
			}
		}
	})
}

const selfContainedJson = `{
	"$type": "Module",
	"name": "test",
	"statements": [
		{
			"$type": "Definition",
			"name": "two",
			"expression": {
				"$type": "NumberLiteral",
				"value": "2"
			}
		},
		{
			"$type": "Definition",
			"name": "root",
			"args": [
				{
					"$type": "DeclaredParameter",
					"name": "x"
				},
				{
					"$type": "DeclaredParameter",
					"name": "y"
				}
			],
			"expression": {
				"$type": "BinaryExpression",
				"left": {
					"$type": "FunctionCall",
					"callable": {
						"$refText": "x",
						"$ref": "#/statements@1/args@0"
					}
				},
				"operator": "^",
				"right": {
					"$type": "BinaryExpression",
					"left": {
						"$type": "NumberLiteral",
						"value": "1"
					},
					"operator": "/",
					"right": {
						"$type": "FunctionCall",
						"callable": {
							"$refText": "y",
							"$ref": "#/statements@1/args@1"
						}
					}
				}
			}
		},
		{
			"$type": "Definition",
			"name": "sqrt",
			"args": [
				{
					"$type": "DeclaredParameter",
					"name": "x"
				}
			],
			"expression": {
				"$type": "FunctionCall",
				"args": [
					{
						"$type": "FunctionCall",
						"callable": {
							"$refText": "x",
							"$ref": "#/statements@2/args@0"
						}
					},
					{
						"$type": "FunctionCall",
						"callable": {
							"$refText": "two",
							"$ref": "#/statements@0"
						}
					}
				],
				"callable": {
					"$refText": "root",
					"$ref": "#/statements@1"
				}
			}
		},
		{
			"$type": "Evaluation",
			"expression": {
				"$type": "FunctionCall",
				"args": [
					{
						"$type": "NumberLiteral",
						"value": "16"
					}
				],
				"callable": {
					"$refText": "sqrt",
					"$ref": "#/statements@2"
				}
			}
		}
	]
}`

const circularAJson = `{
	"$type": "Module",
	"name": "A",
	"statements": [
		{
			"$type": "Definition",
			"name": "root",
			"args": [
				{
					"$type": "DeclaredParameter",
					"name": "x"
				},
				{
					"$type": "DeclaredParameter",
					"name": "y"
				}
			],
			"expression": {
				"$type": "BinaryExpression",
				"left": {
					"$type": "FunctionCall",
					"callable": {
						"$refText": "two",
						"$ref": "inmemory:/testB#/statements@0"
					}
				},
				"operator": "*",
				"right": {
					"$type": "BinaryExpression",
					"left": {
						"$type": "FunctionCall",
						"callable": {
							"$refText": "x",
							"$ref": "#/statements@0/args@0"
						}
					},
					"operator": "^",
					"right": {
						"$type": "BinaryExpression",
						"left": {
							"$type": "NumberLiteral",
							"value": "1"
						},
						"operator": "/",
						"right": {
							"$type": "FunctionCall",
							"callable": {
								"$refText": "y",
								"$ref": "#/statements@0/args@1"
							}
						}
					}
				}
			}
		}
	]
}`

const circularBJson = `{
	"$type": "Module",
	"name": "B",
	"statements": [
		{
			"$type": "Definition",
			"name": "two",
			"expression": {
				"$type": "NumberLiteral",
				"value": "2"
			}
		},
		{
			"$type": "Definition",
			"name": "sqrt",
			"args": [
				{
					"$type": "DeclaredParameter",
					"name": "x"
				}
			],
			"expression": {
				"$type": "FunctionCall",
				"args": [
					{
						"$type": "FunctionCall",
						"callable": {
							"$refText": "x",
							"$ref": "#/statements@1/args@0"
						}
					},
					{
						"$type": "FunctionCall",
						"callable": {
							"$refText": "two",
							"$ref": "#/statements@0"
						}
					}
				],
				"callable": {
					"$refText": "root",
					"$ref": "inmemory:/testA#/statements@0"
				}
			}
		},
		{
			"$type": "Evaluation",
			"expression": {
				"$type": "FunctionCall",
				"args": [
					{
						"$type": "NumberLiteral",
						"value": "16"
					}
				],
				"callable": {
					"$refText": "sqrt",
					"$ref": "#/statements@1"
				}
			}
		}
	]
}`
