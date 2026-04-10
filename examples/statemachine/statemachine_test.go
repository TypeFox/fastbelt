// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	fbtest "typefox.dev/fastbelt/testing"
)

// lightSwitch is a minimal two-state machine used across several tests.
const lightSwitch = `
statemachine LightSwitch

events
  flick

commands
  turnOn
  turnOff

initialState off

state off
  actions { turnOff }
  flick => on
end

state on
  actions { turnOn }
  flick => off
end
`

func TestParsing(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(lightSwitch).
		AssertNoErrors().
		AssertState(core.DocStateLinked)
}

func TestAST(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(lightSwitch).AssertNoErrors()

	sm := fbtest.MustFindNode[Statemachine](doc)
	assert.Equal(t, "LightSwitch", sm.Name())
	assert.Len(t, sm.Events(), 1)
	assert.Len(t, sm.Commands(), 2)
	assert.Len(t, sm.States(), 2)
}

func TestFindAll(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(lightSwitch).AssertNoErrors()

	events := fbtest.FindAll[Event](doc)
	require.Len(t, events, 1)
	assert.Equal(t, "flick", events[0].Name())

	states := fbtest.FindAll[State](doc)
	require.Len(t, states, 2)
}

func TestFindNamedNode(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(lightSwitch).AssertNoErrors()

	off := fbtest.MustFindNamedNode[State](doc, "off")
	assert.Len(t, off.Actions(), 1)
	assert.Len(t, off.Transitions(), 1)

	on := fbtest.MustFindNamedNode[State](doc, "on")
	assert.Len(t, on.Actions(), 1)
	assert.Len(t, on.Transitions(), 1)
}

func TestInitialStateReference(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(lightSwitch).AssertNoErrors()

	sm := fbtest.MustFindNode[Statemachine](doc)
	initRef := sm.Init()
	require.NotNil(t, initRef)
	require.Nil(t, initRef.Error(), "Init reference should resolve without error")

	resolved := initRef.Ref(doc.Ctx())
	require.NotNil(t, resolved)
	assert.Equal(t, "off", resolved.Name())
}

func TestTransitionReferences(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(lightSwitch).AssertNoErrors()

	offState := fbtest.MustFindNamedNode[State](doc, "off")
	transitions := offState.Transitions()
	require.Len(t, transitions, 1)

	t0 := transitions[0]

	eventRef := t0.Event()
	require.NotNil(t, eventRef)
	assert.Nil(t, eventRef.Error())
	assert.Equal(t, "flick", eventRef.Ref(doc.Ctx()).Name())

	stateRef := t0.State()
	require.NotNil(t, stateRef)
	assert.Nil(t, stateRef.Error())
	assert.Equal(t, "on", stateRef.Ref(doc.Ctx()).Name())
}

func TestFindReferenceWithText(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(lightSwitch).AssertNoErrors()

	// Locate the "off" cross-reference (the one in the "on" state's transition).
	ref := fbtest.MustFindReferenceWithText[State](doc, "off")
	assert.Nil(t, ref.Error())
	resolved := ref.Ref(doc.Ctx())
	require.NotNil(t, resolved)
	assert.Equal(t, "off", resolved.Name())
}

// TestFindReferenceAtMarker demonstrates the range marker syntax for pinning
// a reference lookup to a specific token when multiple references share the
// same text (both transitions reference the event "flick" here).
func TestFindReferenceAtMarker(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		statemachine Toggle

		events
		flick

		initialState a

		state a
		<|aFlick:flick|> => b
		end

		state b
		<|bFlick:flick|> => a
		end
		`).AssertNoErrors()

	// Both references use the text "flick", so FindReferenceWithText would
	// return whichever comes first. Markers let us pick the one we want.
	aRef := fbtest.MustFindReference[Event](doc, "aFlick")
	assert.Nil(t, aRef.Error())
	assert.Equal(t, "flick", aRef.Ref(doc.Ctx()).Name())

	bRef := fbtest.MustFindReference[Event](doc, "bFlick")
	assert.Nil(t, bRef.Error())
	assert.Equal(t, "flick", bRef.Ref(doc.Ctx()).Name())
}

func TestFindNodeAtLabel(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		statemachine S
		events e
		initialState <|cursor>idle
		state idle
		e => idle
		end
		`).AssertNoErrors()

	// The index marker sits at the start of the "idle" reference token.
	// FindNodeAtLabel resolves to the smallest AST node containing that offset.
	// Since "idle" is a reference token owned by the Statemachine node, we
	// expect the Statemachine to be the containing node.
	sm, ok := fbtest.FindNodeAtLabel[Statemachine](doc, "cursor")
	assert.True(t, ok)
	assert.Equal(t, "S", sm.Name())
}

func TestParseError(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`this is not a valid statemachine`)
	assert.NotEmpty(t, doc.Document.ParserErrors, "invalid input should produce parser errors")
}

func TestUnresolvedReference(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		statemachine Broken
		events click
		initialState missing
		state real
		click => doesNotExist
		end
		`).AssertNoParseErrors()

	sm := fbtest.MustFindNode[Statemachine](doc)
	initRef := sm.Init()
	require.NotNil(t, initRef)
	// "missing" does not name any state — resolution should fail.
	assert.NotNil(t, initRef.Error(), "reference to nonexistent state should have an error")
}
