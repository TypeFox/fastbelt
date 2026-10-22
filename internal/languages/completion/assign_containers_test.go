// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package completion_test

import (
	"testing"
	"unique"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/languages/completion"
	"typefox.dev/fastbelt/test"
)

// The completion language references declares via the FQN composite rule, so
// this covers the composite reference unit branch of AssignContainers, which
// the statemachine example cannot exercise.
func TestAssignContainersConsistency(t *testing.T) {
	f := test.New(t, completion.CreateServices(&SimpleCompletionContributor{}))
	doc := f.Parse(`
declare foo { declare bar { declare baz } }
declare qux.quux
fqn qux.quux
list qux.quux qux.quux
ref foo
`)
	doc.AssertNoErrors()

	root := doc.Root()
	if root.Container() != nil {
		t.Errorf("root node: expected nil container, got %v", root.Container())
	}
	if root.Document() != doc.Document {
		t.Errorf("root node: expected document %p, got %p", doc.Document, root.Document())
	}

	check := func(kind string, node core.AstNode, parent core.AstNode, field unique.Handle[string], index int) {
		t.Helper()
		gotField, gotIndex := node.ContainmentData()
		if node.Container() != parent {
			t.Errorf("%s %T: expected container %T, got %T", kind, node, parent, node.Container())
		}
		if node.Document() != doc.Document {
			t.Errorf("%s %T: expected document %p, got %p", kind, node, doc.Document, node.Document())
		}
		if gotField != field {
			t.Errorf("%s %T: expected containment field %q, got %q", kind, node, field.Value(), gotField.Value())
		}
		if gotIndex != index {
			t.Errorf("%s %T: expected containment index %d, got %d", kind, node, index, gotIndex)
		}
	}

	var visited, visitedUnits int
	var walk func(parent core.AstNode)
	walk = func(parent core.AstNode) {
		parent.ForEachNode(func(child core.AstNode, field unique.Handle[string], index int) {
			visited++
			check("node", child, parent, field, index)
			walk(child)
		})
		parent.ForEachReference(func(ref core.UntypedReference, field unique.Handle[string], index int) {
			if unit, ok := ref.Unit().(core.CompositeNode); ok {
				visitedUnits++
				check("reference unit", unit, parent, field, index)
			}
		})
	}
	walk(root)

	if visited == 0 {
		t.Fatal("fixture produced no child nodes, test would trivially pass")
	}
	if visitedUnits == 0 {
		t.Fatal("fixture produced no composite reference units, reference branch not exercised")
	}
}
