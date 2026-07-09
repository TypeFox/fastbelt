// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"encoding/json"
	"iter"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"unique"

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
			assertEqualAst(t, exported, imported) // expected := exported, actual := imported
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
	core.AssignContainers(doc2, grammar2)
	documents2.Set(doc2)
	f2.NewDoc(doc2, nil, nil)

	// Step 5: Run the full build pipeline on the unmarshaled document
	err = builder2.Build(f2.Ctx(), slices.Collect(documents2.All()), nil)
	require.NoError(t, err)

	return doc1.Root(), doc2.Root
}

var REF_TOKEN_TYPE = reflect.TypeFor[*core.Token]()

// assertEqualAst recursively compares two AST subtrees node by node. Diffs are
// reported via t.Errorf with the node's document path from core.PathOf.
// Returns true when at least one difference was detected.
func assertEqualAst(t testing.TB, expected, actual core.AstNode) bool {
	t.Helper()

	// 1. Concrete type must match.
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		p, _ := core.PathOf(expected)
		t.Errorf("at %q: type mismatch: expected %T, got %T", p, expected, actual)
		return false
	}

	// 2. If named names must match.
	type named interface{ Name() string }
	expectedNamed, expectedIsNamed := expected.(named)
	actualNamed, actualIsNamed := actual.(named)
	if expectedIsNamed && actualIsNamed && expectedNamed.Name() != actualNamed.Name() {
		p, _ := core.PathOf(expected)
		t.Errorf("at %q: Name mismatch: expected %s, got %s", p, expectedNamed.Name(), actualNamed.Name())
	}

	// 3. Primitive field values must match
	//  create a reflect.Value of the AstNode:
	//   * fetch child (1) that is the specific '...Data' struct, child (0) is the 'AstNodeBase' struct
	//   * the derive a reference value via '.Addr()', and
	//   * wrap it into an array being used as argument while calling the getters, they're defined with pointer receivers
	expectedValue := reflect.ValueOf(expected).Elem().FieldByIndex([]int{1})
	expectedMethodArg := []reflect.Value{expectedValue.Addr()}

	actualValue := reflect.ValueOf(actual).Elem().FieldByIndex([]int{1})
	actualMethodArg := []reflect.Value{actualValue.Addr()}

	// * iterate the '...Data' type fields and consider those of the types 'bool' and '*core.Token'

	for _, field := range slices.Collect(expectedValue.Type().Fields()) {
		kind := field.Type.Kind()
		switch {
		case kind == reflect.Bool:
			p, _ := core.PathOf(expected)
			getter, exists := expectedValue.Addr().Type().MethodByName(strings.ToUpper(field.Name[0:1]) + field.Name[1:])
			if !exists {
				t.Errorf("at %q, type %T: string value getter for field %s missing", p, expectedValue, field.Name)
			}
			exp := getter.Func.Call(expectedMethodArg)[0].Bool()
			act := getter.Func.Call(actualMethodArg)[0].Bool()
			if exp != act {
				t.Errorf("at %q: primitive bool field '%s' mismatch\n  expected: %t\n  actual: %t", p, field.Name, exp, act)
			}
		case kind == reflect.Pointer && field.Type == REF_TOKEN_TYPE:
			p, _ := core.PathOf(expected)
			getter, exists := expectedValue.Addr().Type().MethodByName(strings.ToUpper(field.Name[0:1]) + field.Name[1:])
			if !exists {
				t.Errorf("at %q, type %T: string value getter for field %s missing", p, expectedValue, field.Name)
			}
			exp := getter.Func.Call(expectedMethodArg)[0].String()
			act := getter.Func.Call(actualMethodArg)[0].String()
			if exp != act {
				t.Errorf("at %q: primitive string field '%s' mismatch\n  expected: %s\n  actual: %s", p, field.Name, exp, act)
			}
		}
	}

	// Collect child nodes from both sides.
	type child struct {
		node    core.AstNode
		feature unique.Handle[string]
		index   int
	}
	var (
		expectedChildren = make([]child, 0, 10)
		actualChildren   = make([]child, 0, 10)
	)
	expected.ForEachNode(func(node core.AstNode, feature unique.Handle[string], index int) {
		expectedChildren = append(expectedChildren, child{node, feature, index})
	})
	actual.ForEachNode(func(node core.AstNode, feature unique.Handle[string], index int) {
		actualChildren = append(actualChildren, child{node, feature, index})
	})

	// 4. Recurse into each corresponding child pair.
	//  * check presence (amount) of children on each side
	//  * check containment (feature, index)
	//  * check deep equality of child nodes
	expectedIter, stopE := iter.Pull(slices.Values(expectedChildren))
	defer stopE()
	actualIter, stopA := iter.Pull(slices.Values(actualChildren))
	defer stopA()
	for {
		itemE, validE := expectedIter()
		itemA, validA := actualIter()

		if !validE && !validA {
			break
		}
		if validE != validA {
			if validE {
				pChild, _ := core.PathOf(itemE.node)
				pActual, _ := core.PathOf(actual)
				t.Errorf("at %q: child element mismatch\n  no counter part for %s of expected (field: %s, index: %d) in actual ", pActual, pChild, itemE.feature.Value(), itemE.index)
			} else {
				pChild, _ := core.PathOf(itemA.node)
				pExpected, _ := core.PathOf(expected)
				t.Errorf("at %q: child element mismatch\n  no counter part for %s of actual (field: %s, index: %d) in expected ", pExpected, pChild, itemA.feature.Value(), itemA.index)
			}
		} else if itemE.feature != itemA.feature {
			p, _ := core.PathOf(expected)
			t.Errorf("at %q: child feature mismatch\n  child contained in %s (expected, index: %d) vs %s (actual, index: %d)", p, itemE.feature.Value(), itemE.index, itemA.feature.Value(), itemA.index)

			// no need to check for index diff/equality, as position mismatches will cause other diffs being checked above or below, like container feature diffs or child type/primitive field value diffs
			// therefore, continue with comparing the individual children
		} else if !assertEqualAst(t, itemE.node, itemA.node) {
			return false
		}
	}

	return true
}
