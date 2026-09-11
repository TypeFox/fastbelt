// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
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
		mod, err := UnmarshalValue[Module]([]byte(jsonInput))
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
						"$ref": "inmemory:/testB#%2Fstatements%400"
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
					"$ref": "inmemory:/testA#%2Fstatements%400"
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

// TestJsonWithAbsentFields programmatically composes an AST containing token fields that were
// never set as opposed to set to an empty string, and asserts that MarshalJSONTo /
// UnmarshalJSONFrom preserve that distinction through a roundtrip: an absent token field must
// marshal without its JSON property (not as ""), a token explicitly set to "" must marshal as
// "" (not be dropped), and both must revive as such - not as one another - after unmarshaling.
func TestJsonWithAbsentFields(t *testing.T) {
	mod := NewModule()
	mod.SetName(core.NewSyntheticToken("test", mod))

	def := NewDefinition()
	def.SetName(core.NewSyntheticToken("f", def))

	// named: Name is set to a non-empty string
	named := NewDeclaredParameter()
	named.SetName(core.NewSyntheticToken("x", named))
	def.SetArgsItem(named)

	// unnamed: Name is intentionally left unset (nil token), as opposed to an empty string
	unnamed := NewDeclaredParameter()
	def.SetArgsItem(unnamed)
	require.Nil(t, unnamed.NameToken(), "test setup: Name must be unset")

	// emptyNamed: Name is explicitly set to an empty string, as opposed to being left unset
	emptyNamed := NewDeclaredParameter()
	emptyNamed.SetName(core.NewSyntheticToken("", emptyNamed))
	def.SetArgsItem(emptyNamed)

	oneLiteral := NewNumberLiteral()
	oneLiteral.SetValue(core.NewSyntheticToken("1", oneLiteral))

	// emptyLiteral: Value is explicitly set to an empty string, as opposed to being left unset
	emptyLiteral := NewNumberLiteral()
	emptyLiteral.SetValue(core.NewSyntheticToken("", emptyLiteral))

	sum := NewBinaryExpression()
	sum.SetLeft(oneLiteral)
	sum.SetOperator(core.NewSyntheticToken("+", sum))
	sum.SetRight(emptyLiteral)
	def.SetExpression(sum)

	mod.SetStatementsItem(def)

	res, err := json.Marshal(mod, jsontext.WithIndent("	"))
	require.NoError(t, err)
	assert.Equal(t, absentFieldsJson, string(res))

	revived, err := UnmarshalValue[Module](res)
	require.NoError(t, err)

	revivedArgs := revived.Statements()[0].(Definition).Args()
	assert.Equal(t, "x", revivedArgs[0].Name())
	assert.Nil(t, revivedArgs[1].NameToken(), "an absent 'name' property must revive as a nil token, not a token with an empty image")
	assert.NotNil(t, revivedArgs[2].NameToken(), "a 'name' property explicitly set to \"\" must revive as a non-nil token with an empty image")
	assert.Equal(t, "", revivedArgs[2].Name())

	revivedLeft := revived.Statements()[0].(Definition).Expression().(BinaryExpression).Left().(NumberLiteral)
	assert.Equal(t, "1", revivedLeft.Value())

	revivedRight := revived.Statements()[0].(Definition).Expression().(BinaryExpression).Right().(NumberLiteral)
	assert.NotNil(t, revivedRight.ValueToken(), "a 'value' property explicitly set to \"\" must revive as a non-nil token with an empty image")
	assert.Equal(t, "", revivedRight.Value())
}

const absentFieldsJson = `{
	"$type": "Module",
	"name": "test",
	"statements": [
		{
			"$type": "Definition",
			"name": "f",
			"args": [
				{
					"$type": "DeclaredParameter",
					"name": "x"
				},
				{
					"$type": "DeclaredParameter"
				},
				{
					"$type": "DeclaredParameter",
					"name": ""
				}
			],
			"expression": {
				"$type": "BinaryExpression",
				"left": {
					"$type": "NumberLiteral",
					"value": "1"
				},
				"operator": "+",
				"right": {
					"$type": "NumberLiteral",
					"value": ""
				}
			}
		}
	]
}`
