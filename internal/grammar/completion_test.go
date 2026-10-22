// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The fields that are offered for an assignment depend on the actions in
// front of it. The element that is completed doesn't exist yet, so the
// completion has to find its place among the elements in front of the cursor.
func TestCompletionOfFieldsBehindAction(t *testing.T) {
	cases := []struct {
		rule     string
		expected []string
		excluded []string
	}{
		{"Foo: <|cursor@>;", []string{"A"}, []string{"B"}},
		{"Foo: A=ID <|cursor@>;", []string{"A"}, []string{"B"}},
		{"Foo: {Bar} <|cursor@>;", []string{"A", "B"}, nil},
		{"Foo: {Bar} B=ID <|cursor@>;", []string{"A", "B"}, nil},
		{"Foo: A=ID | {Bar} <|cursor@>;", []string{"A", "B"}, nil},
		// The cursor is inside of the group of the action
		{"Foo: A=ID ({Bar} <|cursor@>);", []string{"A", "B"}, nil},
		// The cursor is behind the group of the action
		{"Foo: ({Bar} B=ID)? <|cursor@>;", []string{"A"}, []string{"B"}},
	}
	for _, c := range cases {
		t.Run(c.rule, func(t *testing.T) {
			doc := lspFixture(t).ParseURI(`grammar Test;
interface Foo { A string }
interface Bar extends Foo { B string }
`+c.rule+`
token ID: /[a-z]+/
hidden token WS: /\s+/`, "file:///completion-fields.fb")
			labels := completionLabels(doc.CompletionItems("cursor"))
			for _, expected := range c.expected {
				assert.Contains(t, labels, expected)
			}
			for _, excluded := range c.excluded {
				assert.NotContains(t, labels, excluded)
			}
		})
	}
}
