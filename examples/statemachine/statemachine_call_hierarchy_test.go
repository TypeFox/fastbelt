// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// statemachine used across the call hierarchy tests:
//
//	off --flick--> on --flick--> off
//	on  --alarm--> panic
//
// So: off has 1 outgoing call (-> on) and 1 incoming call (from on).
//     on  has 2 outgoing calls (-> off, -> panic) and 1 incoming call (from off).
//     panic has 0 outgoing calls and 1 incoming call (from on).
const callHierarchySource = `
statemachine Test
events flick alarm
initialState off

state off
  flick => on
end

state on
  flick => off
  alarm => panic
end

state panic
end
`

func prepareStateItem(t *testing.T, provider server.CallHierarchyProvider, uri lsp.DocumentURI, line, character uint32) lsp.CallHierarchyItem {
	t.Helper()
	items, err := provider.HandlePrepareCallHierarchyRequest(
		context.Background(),
		&lsp.CallHierarchyPrepareParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{URI: uri},
				Position:     lsp.Position{Line: line, Character: character},
			},
		},
	)
	require.NoError(t, err)
	require.Len(t, items, 1)
	return items[0]
}

func TestCallHierarchy_PrepareOnStateName(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(callHierarchySource).AssertNoErrors()

	provider := service.MustGet[server.CallHierarchyProvider](f.Services())

	// Line 4 (0-indexed): "state off" - cursor on "off"
	item := prepareStateItem(t, provider, doc.Document.URI.DocumentURI(), 5, 7)
	assert.Equal(t, "off", item.Name)
	assert.Equal(t, lsp.Class, item.Kind)
}

func TestCallHierarchy_PrepareOnNonState(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(callHierarchySource).AssertNoErrors()

	provider := service.MustGet[server.CallHierarchyProvider](f.Services())

	// Line 1: "statemachine Test" - not a state
	items, err := provider.HandlePrepareCallHierarchyRequest(
		context.Background(),
		&lsp.CallHierarchyPrepareParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
				Position:     lsp.Position{Line: 1, Character: 15},
			},
		},
	)
	assert.NoError(t, err)
	assert.Empty(t, items)
}

func TestCallHierarchy_OutgoingCalls(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(callHierarchySource).AssertNoErrors()

	provider := service.MustGet[server.CallHierarchyProvider](f.Services())
	uri := doc.Document.URI.DocumentURI()

	onItem := prepareStateItem(t, provider, uri, 9, 7) // "state on"
	outgoing, err := provider.HandleOutgoingCallsRequest(
		context.Background(),
		&lsp.CallHierarchyOutgoingCallsParams{Item: onItem},
	)
	assert.NoError(t, err)
	assert.Len(t, outgoing, 2, "state 'on' transitions to both 'off' and 'panic'")

	names := map[string]int{}
	for _, call := range outgoing {
		names[call.To.Name]++
		assert.Len(t, call.FromRanges, 1)
	}
	assert.Equal(t, 1, names["off"])
	assert.Equal(t, 1, names["panic"])
}

func TestCallHierarchy_OutgoingCallsGroupsMultipleTransitionsToSameTarget(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(`
statemachine Test
events a b
initialState off

state off
  a => on
  b => on
end

state on
end
`).AssertNoErrors()

	provider := service.MustGet[server.CallHierarchyProvider](f.Services())
	uri := doc.Document.URI.DocumentURI()

	offItem := prepareStateItem(t, provider, uri, 5, 7) // "state off"
	outgoing, err := provider.HandleOutgoingCallsRequest(
		context.Background(),
		&lsp.CallHierarchyOutgoingCallsParams{Item: offItem},
	)
	assert.NoError(t, err)
	require.Len(t, outgoing, 1, "both transitions target 'on', so they should be grouped into one call entry")
	assert.Equal(t, "on", outgoing[0].To.Name)
	assert.Len(t, outgoing[0].FromRanges, 2, "grouped entry should list both transition ranges")
}

func TestCallHierarchy_IncomingCalls(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	doc := f.Parse(callHierarchySource).AssertNoErrors()

	provider := service.MustGet[server.CallHierarchyProvider](f.Services())
	uri := doc.Document.URI.DocumentURI()

	// "panic" is only reachable from "on" (via alarm).
	panicItem := prepareStateItem(t, provider, uri, 14, 7) // "state panic"
	incoming, err := provider.HandleIncomingCallsRequest(
		context.Background(),
		&lsp.CallHierarchyIncomingCallsParams{Item: panicItem},
	)
	assert.NoError(t, err)
	require.Len(t, incoming, 1)
	assert.Equal(t, "on", incoming[0].From.Name)
	assert.Len(t, incoming[0].FromRanges, 1)

	// "off" is only reachable from "on" (via flick).
	offItem := prepareStateItem(t, provider, uri, 5, 7) // "state off"
	incoming, err = provider.HandleIncomingCallsRequest(
		context.Background(),
		&lsp.CallHierarchyIncomingCallsParams{Item: offItem},
	)
	assert.NoError(t, err)
	require.Len(t, incoming, 1)
	assert.Equal(t, "on", incoming[0].From.Name)

	// "on" is only reachable from "off" (via flick).
	onItem := prepareStateItem(t, provider, uri, 9, 7) // "state on"
	incoming, err = provider.HandleIncomingCallsRequest(
		context.Background(),
		&lsp.CallHierarchyIncomingCallsParams{Item: onItem},
	)
	assert.NoError(t, err)
	require.Len(t, incoming, 1)
	assert.Equal(t, "off", incoming[0].From.Name)
}
