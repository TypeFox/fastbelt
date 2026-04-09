// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	fbtest "typefox.dev/fastbelt/testing"
)

// commonTokens is appended to every test grammar to avoid repeating token definitions.
const commonTokens = `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`

// --- Assignment.Property linking ---

// TestAssignmentPropertyImplicitReturnType verifies that Assignment.Property resolves
// when the rule name implicitly matches the interface name (no explicit "returns").
func TestAssignmentPropertyImplicitReturnType(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|nameRef:Name|>=ID;
	` + commonTokens).AssertNoErrors()

	nameRef := fbtest.MustFindReference[Field](doc, "nameRef")
	require.Nil(t, nameRef.Error(), "Assignment.Property should resolve without error")

	field := nameRef.Ref(doc.Ctx())
	require.NotNil(t, field)
	assert.Equal(t, "Name", field.Name())

	fooIface := fbtest.MustFindNamedNode[Interface](doc, "Foo")
	require.Len(t, fooIface.Fields(), 1)
	assert.Same(t, fooIface.Fields()[0], field, "Property should resolve to the exact Field node on Foo")
}

// TestAssignmentPropertyExplicitReturnType verifies that an explicit "returns" clause
// correctly scopes Assignment.Property resolution to the declared interface.
func TestAssignmentPropertyExplicitReturnType(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string Owner string }
		Bar returns Foo: <|nameRef:Name|>=ID <|ownerRef:Owner|>=ID;
	` + commonTokens).AssertNoErrors()

	nameRef := fbtest.MustFindReference[Field](doc, "nameRef")
	require.Nil(t, nameRef.Error())
	assert.Equal(t, "Name", nameRef.Ref(doc.Ctx()).Name())

	ownerRef := fbtest.MustFindReference[Field](doc, "ownerRef")
	require.Nil(t, ownerRef.Error())
	assert.Equal(t, "Owner", ownerRef.Ref(doc.Ctx()).Name())

	// Both should resolve to fields on Foo (the return type), not on Bar.
	fooIface := fbtest.MustFindNamedNode[Interface](doc, "Foo")
	assert.Same(t, fooIface.Fields()[0], nameRef.Ref(doc.Ctx()))
	assert.Same(t, fooIface.Fields()[1], ownerRef.Ref(doc.Ctx()))
}

// TestAssignmentPropertySingleLevelInheritance verifies that Assignment.Property can
// resolve a field declared on a direct parent interface.
func TestAssignmentPropertySingleLevelInheritance(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base { BaseName string }
		interface Child extends Base { OwnName string }
		Child: <|baseRef:BaseName|>=ID <|ownRef:OwnName|>=ID;
	` + commonTokens).AssertNoErrors()

	// BaseName is declared on Base — requires single-level inheritance traversal.
	baseRef := fbtest.MustFindReference[Field](doc, "baseRef")
	require.Nil(t, baseRef.Error(), "inherited field BaseName should resolve without error")
	baseField := baseRef.Ref(doc.Ctx())
	require.NotNil(t, baseField)
	assert.Equal(t, "BaseName", baseField.Name())
	baseIface := fbtest.MustFindNamedNode[Interface](doc, "Base")
	assert.Same(t, baseIface.Fields()[0], baseField, "BaseName should resolve to the Base interface field")

	// OwnName is declared directly on Child.
	ownRef := fbtest.MustFindReference[Field](doc, "ownRef")
	require.Nil(t, ownRef.Error())
	ownField := ownRef.Ref(doc.Ctx())
	require.NotNil(t, ownField)
	assert.Equal(t, "OwnName", ownField.Name())
	childIface := fbtest.MustFindNamedNode[Interface](doc, "Child")
	assert.Same(t, childIface.Fields()[0], ownField)
}

// TestAssignmentPropertyMultiLevelInheritance verifies that Assignment.Property can
// resolve a field from a grandparent interface through a chain of extensions.
func TestAssignmentPropertyMultiLevelInheritance(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface A { FieldA string }
		interface B extends A { FieldB string }
		interface C extends B { FieldC string }
		Rule returns C: <|refA:FieldA|>=ID <|refB:FieldB|>=ID <|refC:FieldC|>=ID;
	` + commonTokens).AssertNoErrors()

	refA := fbtest.MustFindReference[Field](doc, "refA")
	require.Nil(t, refA.Error(), "grandparent field FieldA should resolve")
	assert.Equal(t, "FieldA", refA.Ref(doc.Ctx()).Name())
	ifaceA := fbtest.MustFindNamedNode[Interface](doc, "A")
	assert.Same(t, ifaceA.Fields()[0], refA.Ref(doc.Ctx()), "FieldA should resolve to the A interface field")

	refB := fbtest.MustFindReference[Field](doc, "refB")
	require.Nil(t, refB.Error(), "parent field FieldB should resolve")
	assert.Equal(t, "FieldB", refB.Ref(doc.Ctx()).Name())
	ifaceB := fbtest.MustFindNamedNode[Interface](doc, "B")
	assert.Same(t, ifaceB.Fields()[0], refB.Ref(doc.Ctx()), "FieldB should resolve to the B interface field")

	refC := fbtest.MustFindReference[Field](doc, "refC")
	require.Nil(t, refC.Error(), "own field FieldC should resolve")
	assert.Equal(t, "FieldC", refC.Ref(doc.Ctx()).Name())
	ifaceC := fbtest.MustFindNamedNode[Interface](doc, "C")
	assert.Same(t, ifaceC.Fields()[0], refC.Ref(doc.Ctx()), "FieldC should resolve to the C interface field")
}

// --- Interface.Extends linking ---

// TestInterfaceExtendsReference verifies that Interface.Extends cross-references
// resolve to the correct parent interface nodes, including multiple extends.
func TestInterfaceExtendsReference(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base {}
		interface Mixin {}
		interface Child extends <|baseRef:Base|>, <|mixinRef:Mixin|> {}
	` + commonTokens).AssertNoErrors()

	baseRef := fbtest.MustFindReference[Interface](doc, "baseRef")
	require.Nil(t, baseRef.Error())
	baseIface := baseRef.Ref(doc.Ctx())
	require.NotNil(t, baseIface)
	assert.Equal(t, "Base", baseIface.Name())

	mixinRef := fbtest.MustFindReference[Interface](doc, "mixinRef")
	require.Nil(t, mixinRef.Error())
	mixinIface := mixinRef.Ref(doc.Ctx())
	require.NotNil(t, mixinIface)
	assert.Equal(t, "Mixin", mixinIface.Name())

	// Verify the Child interface reports both resolved parents.
	childIface := fbtest.MustFindNamedNode[Interface](doc, "Child")
	require.Len(t, childIface.Extends(), 2)
	assert.Same(t, baseIface, childIface.Extends()[0].Ref(doc.Ctx()))
	assert.Same(t, mixinIface, childIface.Extends()[1].Ref(doc.Ctx()))
}

// --- Action linking ---

// TestActionTypeAndPropertyReferences verifies that Action.Type resolves to the
// declared interface and Action.Property resolves to the field on that interface.
func TestActionTypeAndPropertyReferences(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Node { Name string }
		interface List { Head Node Rest []Node }
		List: Head=Node ({List.Rest+=current} Rest+=Node)*;
		Node returns Node: Name=ID;
	` + commonTokens).AssertNoParseErrors()

	action := fbtest.MustFindNode[Action](doc)

	// Action.Type → List interface.
	typeRef := action.Type()
	require.NotNil(t, typeRef)
	require.Nil(t, typeRef.Error(), "Action.Type should resolve without error")
	resolvedType := typeRef.Ref(doc.Ctx())
	require.NotNil(t, resolvedType)
	assert.Equal(t, "List", resolvedType.Name())
	listIface := fbtest.MustFindNamedNode[Interface](doc, "List")
	assert.Same(t, listIface, resolvedType)

	// Action.Property → List.Rest field.
	propRef := action.Property()
	require.NotNil(t, propRef)
	require.Nil(t, propRef.Error(), "Action.Property should resolve without error")
	resolvedField := propRef.Ref(doc.Ctx())
	require.NotNil(t, resolvedField)
	assert.Equal(t, "Rest", resolvedField.Name())
	assert.Equal(t, "+=", action.Operator())
	assert.Same(t, listIface.Fields()[1], resolvedField)
}

// TestActionPropertyInheritedField verifies that Action.Property can resolve a field
// declared on a parent interface of the action's declared Type.
func TestActionPropertyInheritedField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base { Items []Node }
		interface Child extends Base {}
		interface Node { Name string }
		Child: ({Child.Items+=current} Items+=Node)+;
		Node returns Node: Name=ID;
	` + commonTokens).AssertNoParseErrors()

	action := fbtest.MustFindNode[Action](doc)

	// Action.Type = Child, but Action.Property "Items" is declared on Base.
	typeRef := action.Type()
	require.Nil(t, typeRef.Error())
	assert.Equal(t, "Child", typeRef.Ref(doc.Ctx()).Name())

	propRef := action.Property()
	require.NotNil(t, propRef)
	require.Nil(t, propRef.Error(), "Action.Property should resolve inherited field from Base")
	field := propRef.Ref(doc.Ctx())
	require.NotNil(t, field)
	assert.Equal(t, "Items", field.Name())

	baseIface := fbtest.MustFindNamedNode[Interface](doc, "Base")
	assert.Same(t, baseIface.Fields()[0], field, "Items should resolve to the Base interface field")
}

// --- Post-action scope ---

// TestAssignmentPropertyAfterAction verifies that an assignment appearing after an
// action in the same parser rule group uses the action's Type for scope resolution,
// not the rule's original return type.
func TestAssignmentPropertyAfterAction(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Node { Name string }
		interface Container { Head Node Rest []Node }
		Container: Head=Node ({Container.Rest+=current} <|restRef:Rest|>+=Node)*;
		Node returns Node: Name=ID;
	` + commonTokens).AssertNoParseErrors()

	// "Rest" appears after the action {Container.Rest+=current}.
	// Its scope should come from the action's type (Container), not the outer rule return type.
	restRef := fbtest.MustFindReference[Field](doc, "restRef")
	require.Nil(t, restRef.Error(), "post-action assignment should resolve using action's type as scope")
	field := restRef.Ref(doc.Ctx())
	require.NotNil(t, field)
	assert.Equal(t, "Rest", field.Name())

	containerIface := fbtest.MustFindNamedNode[Interface](doc, "Container")
	assert.Same(t, containerIface.Fields()[1], field)
}

// --- CrossRef.Type linking ---

// TestCrossRefTypeReference verifies that the interface type reference inside a
// cross-reference expression ([Interface:Rule]) resolves to the correct interface.
func TestCrossRefTypeReference(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface State {
			Name string
		}
		interface Machine {
			States []*State
			Active *State
		}
		Machine: "machine" States+=[<|stateTypeRef:State|>:ID] "active" Active=[State:ID];
		State returns State: Name=ID;
	` + commonTokens).AssertNoErrors()

	// The marker pins to the Type reference inside the first [State:ID] cross-ref.
	stateTypeRef := fbtest.MustFindReference[Interface](doc, "stateTypeRef")
	require.Nil(t, stateTypeRef.Error(), "CrossRef.Type should resolve to the State interface")
	resolvedIface := stateTypeRef.Ref(doc.Ctx())
	require.NotNil(t, resolvedIface)
	assert.Equal(t, "State", resolvedIface.Name())

	stateIface := fbtest.MustFindNamedNode[Interface](doc, "State")
	assert.Same(t, stateIface, resolvedIface)
}

func TestLocalReferenceLinking(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Bar *Bar }
		interface Bar { Name string }
		Entry returns <|returnRef:Foo|>: bar=<|ruleRef:SubRule|>;
		SubRule returns Bar: Name=ID;
	` + commonTokens).AssertState(core.DocStateLinked)

	grammar := fbtest.MustFindNode[Grammar](doc)
	assert.Len(t, grammar.Rules(), 2, "grammar should have 2 parser rules")
	assert.Len(t, grammar.Interfaces(), 2, "grammar should have 2 interfaces")

	// ParserRule.ReturnType → Foo interface
	returnRef := fbtest.MustFindReference[Interface](doc, "returnRef")
	require.Nil(t, returnRef.Error())
	assert.Equal(t, "Foo", returnRef.Ref(doc.Ctx()).Name())

	// RuleCall.Rule → SubRule
	ruleRef := fbtest.MustFindReference[AbstractRule](doc, "ruleRef")
	require.Nil(t, ruleRef.Error())
	assert.Equal(t, "SubRule", ruleRef.Ref(doc.Ctx()).Name())
}

func TestAssignmentPropertyScope(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string Age string }
		Person returns Person: "person" <|nameRef:Name|>=ID <|ageRef:Age|>=ID;
	` + commonTokens).AssertNoErrors()

	iface := fbtest.MustFindNamedNode[Interface](doc, "Person")
	require.Len(t, iface.Fields(), 2)

	nameRef := fbtest.MustFindReference[Field](doc, "nameRef")
	require.Nil(t, nameRef.Error())
	assert.Equal(t, "Name", nameRef.Ref(doc.Ctx()).Name())
	assert.Same(t, iface.Fields()[0], nameRef.Ref(doc.Ctx()))

	ageRef := fbtest.MustFindReference[Field](doc, "ageRef")
	require.Nil(t, ageRef.Error())
	assert.Equal(t, "Age", ageRef.Ref(doc.Ctx()).Name())
	assert.Same(t, iface.Fields()[1], ageRef.Ref(doc.Ctx()))
}

func TestCrossDocumentReferenceLinking(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	docs := f.ParseAll(
		"inmemory://types.fb", `
			grammar Types;
			interface Animal { Name string }
		`,
		"inmemory://rules.fb", `
			grammar Rules;
			Animal returns <|Animal|>: Name=ID;
		`+commonTokens,
	)
	typesDoc, rulesDoc := docs[0], docs[1]
	rulesDoc.AssertState(core.DocStateLinked)
	typesDoc.AssertState(core.DocStateLinked)

	// ParserRule.ReturnType should resolve across documents to the Animal interface.
	animalRef := fbtest.MustFindReference[Interface](rulesDoc, "Animal")
	require.Nil(t, animalRef.Error())
	resolvedIface := animalRef.Ref(rulesDoc.Ctx())
	assert.Equal(t, "Animal", resolvedIface.Name())
	assert.Same(t, fbtest.MustFindNamedNode[Interface](typesDoc, "Animal"), resolvedIface)
}
