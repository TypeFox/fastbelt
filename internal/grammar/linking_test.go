// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"testing"

	core "typefox.dev/fastbelt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalReferenceLinking(t *testing.T) {
	srv := CreateServices()
	builder := srv.Workspace().Builder

	doc, err := core.NewDocumentFromString("inmemory://test.fb", "fastbelt", `
grammar Test;

interface Foo {
    Bar *Bar
}

interface Bar {
    Name string
}

Entry returns Foo: bar=SubRule;
SubRule returns Bar: Name=ID;

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`)
	require.NoError(t, err)

	ctx := context.Background()
	err = builder.Build(ctx, []*core.Document{doc})
	require.NoError(t, err)

	assert.True(t, doc.State.Has(core.DocStateLinked), "document should be in Linked state")

	grammar, ok := doc.Root.(Grammar)
	require.True(t, ok, "root should be a Grammar node")

	rules := grammar.Rules()
	require.Len(t, rules, 2, "grammar should have 2 parser rules")
	interfaces := grammar.Interfaces()
	require.Len(t, interfaces, 2, "grammar should have 2 interfaces")

	entryRule := rules[0]
	assert.Equal(t, "Entry", entryRule.Name())

	// ParserRule.ReturnType -> Interface: Entry returns *Foo*
	returnTypeRef := entryRule.ReturnType()
	require.NotNil(t, returnTypeRef, "Entry rule should have a ReturnType reference")
	resolvedInterface := returnTypeRef.Ref(ctx)
	require.NotNil(t, resolvedInterface, "ReturnType reference should resolve")
	assert.Nil(t, returnTypeRef.Error(), "ReturnType reference should have no error")
	assert.Equal(t, "Foo", resolvedInterface.Name())

	// RuleCall.Rule -> AbstractRule: bar=*SubRule*
	body := entryRule.Body()
	require.NotNil(t, body, "Entry rule should have a body")
	assignment, ok := body.(Assignment)
	require.True(t, ok, "body should be an Assignment")
	assignable := assignment.Value()
	require.NotNil(t, assignable, "assignment should have a value")
	ruleCall, ok := assignable.(RuleCall)
	require.True(t, ok, "assignment value should be a RuleCall")
	ruleRef := ruleCall.Rule()
	require.NotNil(t, ruleRef, "RuleCall should have a Rule reference")
	resolvedRule := ruleRef.Ref(ctx)
	require.NotNil(t, resolvedRule, "Rule reference should resolve")
	assert.Nil(t, ruleRef.Error(), "Rule reference should have no error")
	assert.Equal(t, "SubRule", resolvedRule.Name())
}

func TestAssignmentPropertyScope(t *testing.T) {
	srv := CreateServices()
	builder := srv.Workspace().Builder

	doc, err := core.NewDocumentFromString("inmemory://test.fb", "fastbelt", `
grammar Test;

interface Person {
    Name string
    Age string
}

Person returns Person: "person" Name=ID Age=ID;

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`)
	require.NoError(t, err)

	ctx := context.Background()
	err = builder.Build(ctx, []*core.Document{doc})
	require.NoError(t, err)

	grammar, ok := doc.Root.(Grammar)
	require.True(t, ok, "root should be a Grammar node")

	rules := grammar.Rules()
	require.Len(t, rules, 1)
	personRule := rules[0]
	assert.Equal(t, "Person", personRule.Name())

	iface := grammar.Interfaces()[0]
	require.Equal(t, "Person", iface.Name())
	fields := iface.Fields()
	require.Len(t, fields, 2)

	// The rule body is: "person" Name=ID Age=ID
	// With multiple elements this becomes a Group node.
	body := personRule.Body()
	require.NotNil(t, body)
	group, ok := body.(Group)
	require.True(t, ok, "body with multiple elements should be a Group")
	elements := group.Elements()
	require.Len(t, elements, 3, "group should have keyword + 2 assignments")

	// Assignment.Property -> Field: *Name*=ID
	nameAssignment, ok := elements[1].(Assignment)
	require.True(t, ok, "second element should be an Assignment")
	namePropertyRef := nameAssignment.Property()
	require.NotNil(t, namePropertyRef, "assignment should have a Property reference")
	resolvedNameField := namePropertyRef.Ref(ctx)
	require.NotNil(t, resolvedNameField, "Name property reference should resolve")
	assert.Nil(t, namePropertyRef.Error(), "Name property reference should have no error")
	assert.Equal(t, "Name", resolvedNameField.Name())
	assert.Same(t, fields[0], resolvedNameField, "should resolve to the Name field of Person interface")

	// Assignment.Property -> Field: *Age*=ID
	ageAssignment, ok := elements[2].(Assignment)
	require.True(t, ok, "third element should be an Assignment")
	agePropertyRef := ageAssignment.Property()
	require.NotNil(t, agePropertyRef, "assignment should have a Property reference")
	resolvedAgeField := agePropertyRef.Ref(ctx)
	require.NotNil(t, resolvedAgeField, "Age property reference should resolve")
	assert.Nil(t, agePropertyRef.Error(), "Age property reference should have no error")
	assert.Equal(t, "Age", resolvedAgeField.Name())
	assert.Same(t, fields[1], resolvedAgeField, "should resolve to the Age field of Person interface")
}

func TestCrossDocumentReferenceLinking(t *testing.T) {
	srv := CreateServices()
	builder := srv.Workspace().Builder
	docManager := srv.Workspace().DocumentManager

	typesDoc, err := core.NewDocumentFromString("inmemory://types.fb", "fastbelt", `
grammar Types;

interface Animal {
    Name string
}
`)
	require.NoError(t, err)

	rulesDoc, err := core.NewDocumentFromString("inmemory://rules.fb", "fastbelt", `
grammar Rules;

Animal returns Animal: Name=ID;

token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`)
	require.NoError(t, err)

	docManager.Set(typesDoc)
	docManager.Set(rulesDoc)

	ctx := context.Background()
	docs := []*core.Document{typesDoc, rulesDoc}
	err = builder.Build(ctx, docs)
	require.NoError(t, err)

	assert.True(t, typesDoc.State.Has(core.DocStateLinked))
	assert.True(t, rulesDoc.State.Has(core.DocStateLinked))

	// Grab the Animal interface from the types document.
	typesGrammar, ok := typesDoc.Root.(Grammar)
	require.True(t, ok)
	require.Len(t, typesGrammar.Interfaces(), 1)
	animalIface := typesGrammar.Interfaces()[0]
	assert.Equal(t, "Animal", animalIface.Name())

	// Grab the Animal rule from the rules document.
	rulesGrammar, ok := rulesDoc.Root.(Grammar)
	require.True(t, ok)
	require.Len(t, rulesGrammar.Rules(), 1)
	animalRule := rulesGrammar.Rules()[0]
	assert.Equal(t, "Animal", animalRule.Name())

	// ParserRule.ReturnType should resolve across documents to the Animal interface.
	returnTypeRef := animalRule.ReturnType()
	require.NotNil(t, returnTypeRef)
	resolvedIface := returnTypeRef.Ref(ctx)
	require.NotNil(t, resolvedIface, "ReturnType should resolve to interface from other document")
	assert.Nil(t, returnTypeRef.Error(), "ReturnType reference should have no error")
	assert.Equal(t, "Animal", resolvedIface.Name())
	assert.Same(t, animalIface, resolvedIface,
		"should resolve to the exact Interface node from the types document")
}
