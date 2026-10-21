// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package test

import (
	"context"
	"iter"
	"reflect"
	"slices"
	"strings"
	"testing"
	"unique"

	core "typefox.dev/fastbelt"
)

var (
	tokenRefType         = reflect.TypeFor[*core.Token]()
	compositeNodeType    = reflect.TypeFor[core.CompositeNode]()
	untypedReferenceType = reflect.TypeFor[core.UntypedReference]()
)

// AssertEqualAst recursively compares two AST subtrees node by node. Diffs are
// reported via t.Errorf with the node's document path from core.PathOf.
// Will early exit and return false if concrete types of expected and actual mismatch.
// Will return true otherwise.
func AssertEqualAst(t *testing.T, expected, actual core.AstNode) bool {
	t.Helper()
	pExpected, _ := core.PathOf(expected)

	// 1. Concrete type must match.
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Errorf("at %q: type mismatch: expected %T, got %T", pExpected, expected, actual)
		return false
	}

	// 2. If named names must match.
	type named interface{ Name() string }
	expectedNamed, expectedIsNamed := expected.(named)
	actualNamed, actualIsNamed := actual.(named)
	if expectedIsNamed && actualIsNamed && expectedNamed.Name() != actualNamed.Name() {
		t.Errorf("at %q: Name mismatch: expected %s, got %s", pExpected, expectedNamed.Name(), actualNamed.Name())
	}

	// 3. Primitive field values must match
	//  Recall that the actual AstNode is composed of at least two other structs, maybe further ones:
	//    * AstNodeBase
	//    * ...Data
	//   [* ...Data ]
	// Hence, iterate the components (fields) of 'expected', start at field '1' since we need to ignore 'AstNodeBase' here,
	//  for each component then iterate its declared fields;
	// Note: We do _not_ iterate the methods as client code might contribute arbitrary additional methods

	// First, prepare some data required under way:
	expectedValueRef := reflect.ValueOf(expected)
	expectedValueRefType := expectedValueRef.Type()
	expectedValue := expectedValueRef.Elem() // description of the dereferenced value struct
	expectedValueType := expectedValue.Type()

	// Wrap expectedValueRef into an array for being used as argument while calling the corresponding field getters,
	// they're defined with pointer receivers (see `types_gen.go`), so we use the ...Ref here.
	// We can't rely on the fields directly as we wanna compare string and bool values instead of pure tokens.
	// Note that we call the getter on the composed type of 'expectedValue', not the particular component types.
	// That makes a difference if a getter is overridden in client code, and the generated getter is not promoted to the combined type.
	expectedGetterArg := []reflect.Value{expectedValueRef}

	actualValueRef := reflect.ValueOf(actual)
	actualGetterArg := []reflect.Value{actualValueRef}

	for subStructIndex := 1; subStructIndex != expectedValue.NumField(); subStructIndex++ {
		fields := expectedValue.FieldByIndex([]int{subStructIndex}).Type().Fields()

		for field := range fields {
			fieldType := field.Type
			isSlice := fieldType.Kind() == reflect.Slice
			elemType := fieldType
			if isSlice {
				elemType = fieldType.Elem()
			}
			isRelevant := elemType == tokenRefType || elemType == compositeNodeType || elemType.Implements(untypedReferenceType)
			if !isRelevant {
				continue
			}

			// for each field create a sub test, use the composed type's name instead of the component type's name
			// to allow us seeing immediately that the all fields of the composed types are hit
			t.Run(expectedValueType.Name()+"."+field.Name, func(t *testing.T) {
				stringGetterName := strings.ToUpper(field.Name[0:1]) + field.Name[1:]
				if stringGetterName[0] == '_' {
					stringGetterName = stringGetterName[1:]
				}

				getter, exists := expectedValueRefType.MethodByName(stringGetterName)
				if !exists {
					// if no 'string' getter exists, look for a 'bool' getter,
					//  we can't distinguish this on the struct field decls
					getter, exists = expectedValueRefType.MethodByName("Is" + stringGetterName)
					if !exists {
						t.Errorf("at %q, type %T: value getter for field %s missing", pExpected, expected, field.Name)
						return
					}
				}
				expRet := getter.Func.Call(expectedGetterArg)[0]
				actRet := getter.Func.Call(actualGetterArg)[0]

				switch {
				// slice values 1st:
				// no slice of boolean literals expected here, we cannot assign them in parser rules (as of now),
				// proceed with '*core.Token' items:
				case isSlice && elemType == tokenRefType:
					assertEqualTokenSlices(t, pExpected, field.Name, expRet, actRet)
				case isSlice && elemType == compositeNodeType:
					assertEqualCompositeSlices(t, pExpected, field.Name, expRet, actRet)
				case isSlice && elemType.Implements(untypedReferenceType): // reference slice
					if expRet.Len() != actRet.Len() {
						t.Errorf("at %q: reference list field '%s' length mismatch\n  expected: %d\n  actual: %d", pExpected, field.Name, expRet.Len(), actRet.Len())
						return
					}
					for i := 0; i < expRet.Len(); i++ {
						assertEqualReference(t, pExpected, field.Name, expRet.Index(i), actRet.Index(i))
					}
				// non-slice values 2nd:
				case expRet.Kind() == reflect.Bool:
					if exp, act := expRet.Bool(), actRet.Bool(); exp != act {
						t.Errorf("at %q: primitive bool field '%s' mismatch\n  expected: %v\n  actual: %v", pExpected, field.Name, exp, act)
					}
				case elemType == tokenRefType || elemType == compositeNodeType:
					// the getters of singular token and composite fields both return string
					if exp, act := expRet.String(), actRet.String(); exp != act {
						t.Errorf("at %q: primitive string field '%s' mismatch\n  expected: %v\n  actual: %v", pExpected, field.Name, exp, act)
					}
				case elemType.Implements(untypedReferenceType):
					assertEqualReference(t, pExpected, field.Name, expRet, actRet)
				default:
					t.Errorf("at %q: unexpected type %s of field %s", pExpected, field.Type, field.Name)
				}
			})
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
				t.Errorf("at %q: child element mismatch\n  no counter part for %s of actual (field: %s, index: %d) in expected ", pExpected, pChild, itemA.feature.Value(), itemA.index)
			}
		} else if itemE.feature != itemA.feature {
			t.Errorf("at %q: child feature mismatch\n  child contained in %s (expected, index: %d) vs %s (actual, index: %d)", pExpected, itemE.feature.Value(), itemE.index, itemA.feature.Value(), itemA.index)

			// no need to check for index diff/equality, as position mismatches will cause other diffs being checked above or below, like container feature diffs or child type/primitive field value diffs
			// therefore, continue with comparing the individual children
		} else if !AssertEqualAst(t, itemE.node, itemA.node) {
			return false
		}
	}

	return true
}

// assertEqualTokenSlices compares the .Image of each *core.Token in two []*core.Token field values.
func assertEqualTokenSlices(t *testing.T, path any, fieldName string, expected, actual reflect.Value) {
	t.Helper()
	expImages := tokenImages(expected)
	actImages := tokenImages(actual)
	if !slices.Equal(expImages, actImages) {
		t.Errorf("at %q: token list field '%s' mismatch\n  expected: %v\n  actual: %v", path, fieldName, expImages, actImages)
	}
}

func tokenImages(value reflect.Value) []string {
	images := make([]string, value.Len())
	for i := range images {
		if tok, ok := value.Index(i).Interface().(*core.Token); ok && tok != nil {
			images[i] = tok.Image
		}
	}
	return images
}

// assertEqualCompositeSlices compares the .String() of each core.CompositeNode in two
// []core.CompositeNode field values.
func assertEqualCompositeSlices(t *testing.T, path any, fieldName string, expected, actual reflect.Value) {
	t.Helper()
	expImages := compositeImages(expected)
	actImages := compositeImages(actual)
	if !slices.Equal(expImages, actImages) {
		t.Errorf("at %q: composite list field '%s' mismatch\n  expected: %v\n  actual: %v", path, fieldName, expImages, actImages)
	}
}

func compositeImages(value reflect.Value) []string {
	images := make([]string, value.Len())
	for i := range images {
		if cn, ok := value.Index(i).Interface().(core.CompositeNode); ok && cn != nil {
			images[i] = cn.String()
		}
	}
	return images
}

// assertEqualReference compares a single core.UntypedReference field value: the cross-reference
// text, whether resolution failed on both sides, and (if resolved on both) the resolved target's
// document-relative path. A nil reference on either side is treated as unresolved with no text.
func assertEqualReference(t *testing.T, path any, fieldName string, expected, actual reflect.Value) {
	t.Helper()
	expRef, _ := expected.Interface().(core.UntypedReference)
	actRef, _ := actual.Interface().(core.UntypedReference)

	expText, actText := refText(expRef), refText(actRef)
	if expText != actText {
		t.Errorf("at %q: reference field '%s' text mismatch\n  expected: %q\n  actual: %q", path, fieldName, expText, actText)
	}

	expErr, actErr := refError(expRef), refError(actRef)
	if (expErr == nil) != (actErr == nil) {
		t.Errorf("at %q: reference field '%s' resolution mismatch\n  expected error: %v\n  actual error: %v", path, fieldName, expErr, actErr)
		return
	}
	if expErr != nil {
		// Both sides failed to resolve; the exact message may differ if resolution failed for a
		// different reason (e.g. an unavailable JsonLinkingHelper), so there's nothing further to compare.
		return
	}

	expTargetPath, expOk := refTargetPath(expRef)
	actTargetPath, actOk := refTargetPath(actRef)
	if expOk != actOk || expTargetPath != actTargetPath {
		t.Errorf("at %q: reference field '%s' target mismatch\n  expected: %v\n  actual: %v", path, fieldName, expTargetPath, actTargetPath)
	}
}

func refText(ref core.UntypedReference) string {
	if ref == nil {
		return ""
	}
	return ref.Text()
}

func refError(ref core.UntypedReference) *core.ReferenceError {
	if ref == nil {
		return nil
	}
	return ref.Error()
}

// refTargetPath returns the document-relative path of ref's resolved target, or ok=false if ref
// is nil or has no resolved target.
func refTargetPath(ref core.UntypedReference) (path string, ok bool) {
	if ref == nil {
		return "", false
	}
	target := ref.RefNode(context.Background())
	if target == nil {
		return "", false
	}
	fragmentPath, err := core.PathOf(target)
	if err != nil {
		return "", false
	}
	return fragmentPath.String(), true
}
