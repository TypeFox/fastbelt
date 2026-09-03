// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// stateMachineCallHierarchyProvider treats States and Transitions as a call
// graph: a State is like a function, and each Transition to another State is
// like a call site.
//
//   - HandleOutgoingCallsRequest returns the state's own Transitions() -
//     "calls made by this function".
//   - HandleIncomingCallsRequest scans every other state's transitions for
//     ones that target this state - "callers of this function".
//
// HandlePrepareCallHierarchyRequest uses the shared server.NameFinder service
// to resolve the cursor to a state.
type stateMachineCallHierarchyProvider struct {
	sc *service.Container
}

var _ server.CallHierarchyProvider = (*stateMachineCallHierarchyProvider)(nil)

func (p *stateMachineCallHierarchyProvider) document(uri lsp.DocumentURI) *core.Document {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	return documentManager.Get(core.ParseURI(string(uri)))
}

func (p *stateMachineCallHierarchyProvider) HandlePrepareCallHierarchyRequest(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error) {
	doc := p.document(params.TextDocument.URI)
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	offset := doc.TextDoc.OffsetAt(params.Position)
	first, second := doc.Tokens.SearchOffset2(offset)
	if first == nil {
		return nil, nil
	}

	nameFinder := service.MustGet[server.NameFinder](p.sc)
	foundName := nameFinder.Find(ctx, first, second)
	if foundName.Target == nil {
		return nil, nil
	}

	state, ok := foundName.Target.Owner().(State)
	if !ok {
		return nil, nil
	}

	return []lsp.CallHierarchyItem{stateHierarchyItem(doc, state)}, nil
}

func (p *stateMachineCallHierarchyProvider) HandleIncomingCallsRequest(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
	doc := p.document(params.Item.URI)
	if doc == nil || doc.Root == nil {
		return nil, nil
	}
	target, ok := findStateByName(doc, params.Item.Name)
	if !ok {
		return nil, nil
	}

	var calls []lsp.CallHierarchyIncomingCall
	for node := range core.AllNodes(doc.Root) {
		source, ok := node.(State)
		if !ok {
			continue
		}

		var fromRanges []lsp.Range
		for _, transition := range source.Transitions() {
			stateRef := transition.State()
			if stateRef == nil {
				continue
			}
			if resolved := stateRef.Ref(ctx); resolved == target {
				fromRanges = append(fromRanges, stateRef.TextRange().LspRange(doc.TextDoc))
			}
		}
		if len(fromRanges) > 0 {
			calls = append(calls, lsp.CallHierarchyIncomingCall{
				From:       stateHierarchyItem(doc, source),
				FromRanges: fromRanges,
			})
		}
	}
	return calls, nil
}

func (p *stateMachineCallHierarchyProvider) HandleOutgoingCallsRequest(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
	doc := p.document(params.Item.URI)
	if doc == nil || doc.Root == nil {
		return nil, nil
	}
	source, ok := findStateByName(doc, params.Item.Name)
	if !ok {
		return nil, nil
	}

	// Multiple transitions can target the same state; group them into a
	// single outgoing call entry with multiple FromRanges rather than
	// duplicating the "to" item once per transition.
	var order []State
	byTarget := map[State][]lsp.Range{}
	for _, transition := range source.Transitions() {
		stateRef := transition.State()
		if stateRef == nil {
			continue
		}
		target := stateRef.Ref(ctx)
		if target == nil {
			continue
		}
		if _, seen := byTarget[target]; !seen {
			order = append(order, target)
		}
		byTarget[target] = append(byTarget[target], stateRef.TextRange().LspRange(doc.TextDoc))
	}

	var calls []lsp.CallHierarchyOutgoingCall
	for _, target := range order {
		calls = append(calls, lsp.CallHierarchyOutgoingCall{
			To:         stateHierarchyItem(doc, target),
			FromRanges: byTarget[target],
		})
	}
	return calls, nil
}

// stateHierarchyItem builds the LSP representation of a State. The
// selection range narrows down to just the state's own name (via
// NameToken), while Range covers the full "state ... end" block.
func stateHierarchyItem(doc *core.Document, state State) lsp.CallHierarchyItem {
	blockRange := state.TextRange().LspRange(doc.TextDoc)
	selectionRange := blockRange
	if nameToken := state.NameToken(); nameToken != nil {
		selectionRange = nameToken.Range.LspRange(doc.TextDoc)
	}
	return lsp.CallHierarchyItem{
		Name:           state.Name(),
		Kind:           lsp.Class,
		URI:            doc.TextDoc.URI(),
		Range:          blockRange,
		SelectionRange: selectionRange,
	}
}

// findStateByName re-locates a State by name within doc. HandleIncomingCallsRequest
// and HandleOutgoingCallsRequest only receive the lsp.CallHierarchyItem produced
// earlier by HandlePrepareCallHierarchyRequest (name, URI, ranges), not the
// original AST node, so the node must be found again. State names are unique
// within a single statemachine document, so matching by name is sufficient.
func findStateByName(doc *core.Document, name string) (State, bool) {
	for node := range core.AllNodes(doc.Root) {
		if state, ok := node.(State); ok && state.Name() == name {
			return state, true
		}
	}
	return nil, false
}
