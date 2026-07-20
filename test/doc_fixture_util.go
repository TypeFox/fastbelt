// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package test

import (
	"iter"
	"reflect"
	"slices"
	"strings"
	"testing"
	"unique"

	core "typefox.dev/fastbelt"
)

var REF_TOKEN_TYPE = reflect.TypeFor[*core.Token]()

// AssertEqualAst recursively compares two AST subtrees node by node. Diffs are
// reported via t.Errorf with the node's document path from core.PathOf.
// Will early exit and return false if concrete types of expected and actual mismatch.
// Will return true otherwise.
func AssertEqualAst(t testing.TB, expected, actual core.AstNode) bool {
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
	//   * then derive a reference value via '.Addr()', and
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
		} else if !AssertEqualAst(t, itemE.node, itemA.node) {
			return false
		}
	}

	return true
}
