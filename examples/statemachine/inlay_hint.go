// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"fmt"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// stateMachineInlayHintProvider shows the number of outgoing transitions
// inline right after each State's own name, e.g. "state off: 1 transition(s)".
//
// This surfaces the same underlying data as the Code Lens example
// (stateMachineCodeLensProvider's "N transition(s)" lens), just through a
// different LSP mechanism: an inline decoration attached to the name itself
// instead of a clickable lens above the declaration.
type stateMachineInlayHintProvider struct {
	sc *service.Container
}

var _ server.InlayHintProvider = (*stateMachineInlayHintProvider)(nil)

func (p *stateMachineInlayHintProvider) HandleInlayHintRequest(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	var hints []lsp.InlayHint
	for node := range server.NodesInRange(doc, params.Range) {
		state, ok := node.(*StateImpl)
		if !ok {
			continue
		}

		// linking.Name resolves the node's "Name" attribute generically (works
		// for any grammar-generated type with a Name field), giving us the text
		// range of just the identifier itself, not the whole "state ... end" block.
		name := linking.Name(state)
		if name == nil {
			continue
		}

		count := len(state.Transitions())
		hints = append(hints, lsp.InlayHint{
			Position: name.TextRange().LspRange(doc.TextDoc).End,
			Label:    []lsp.InlayHintLabelPart{{Value: fmt.Sprintf(": %d transition(s)", count)}},
			Kind:     lsp.Type,
		})
	}

	return hints, nil
}
