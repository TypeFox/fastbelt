// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"fmt"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// stateMachineCodeLensProvider shows above each Event/Command declaration
// how many places reference it, reusing the shared ReferencesFinder service
// rather than re-implementing reference counting.
//
// States intentionally don't get a lens here: their transition count is
// already shown inline via stateMachineInlayHintProvider, and showing the
// same number in two places on the same line would be redundant.
//
// CodeLensProvider has no default implementation to embed or delegate to
// (see server.CodeLensProvider's doc comment), so this implements the full
// interface directly.
type stateMachineCodeLensProvider struct {
	sc *service.Container
}

var _ server.CodeLensProvider = (*stateMachineCodeLensProvider)(nil)

func (p *stateMachineCodeLensProvider) HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	referencesFinder := service.MustGet[server.ReferencesFinder](p.sc)

	var lenses []lsp.CodeLens
	for node := range core.AllNodes(doc.Root) {
		switch n := node.(type) {
		case *EventImpl:
			count := countReferences(ctx, referencesFinder, n)
			lenses = append(lenses, lsp.CodeLens{
				Range: n.TextRange().LspRange(doc.TextDoc),
				Command: &lsp.Command{
					Title: fmt.Sprintf("%d reference(s)", count),
				},
			})
		case *CommandImpl:
			count := countReferences(ctx, referencesFinder, n)
			lenses = append(lenses, lsp.CodeLens{
				Range: n.TextRange().LspRange(doc.TextDoc),
				Command: &lsp.Command{
					Title: fmt.Sprintf("%d reference(s)", count),
				},
			})
		}
	}

	return lenses, nil
}

func countReferences(ctx context.Context, finder server.ReferencesFinder, target core.AstNode) int {
	count := 0
	for range finder.Find(ctx, target, server.FindReferencesOptions{IncludeDeclaration: false}) {
		count++
	}
	return count
}
