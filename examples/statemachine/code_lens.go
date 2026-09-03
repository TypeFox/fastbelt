// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"encoding/json"
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
// Reference counting happens lazily, via server.ResolvingCodeLensProvider:
// HandleCodeLensRequest only returns the range and a small Data payload
// identifying the declaration, deferring the actual (ReferencesFinder-based)
// count to HandleCodeLensResolveRequest, which the client only calls for
// lenses currently visible in the viewport. This avoids counting references
// for declarations that are scrolled out of view.
//
// CodeLensProvider has no default implementation to embed or delegate to
// (see server.CodeLensProvider's doc comment), so this implements the full
// interface directly.
type stateMachineCodeLensProvider struct {
	sc *service.Container
}

var _ server.ResolvingCodeLensProvider = (*stateMachineCodeLensProvider)(nil)

// codeLensData identifies the declaration a lens belongs to.
type codeLensData struct {
	URI  string `json:"uri"`
	Kind string `json:"kind"` // "event" or "command"
	Name string `json:"name"`
}

func (p *stateMachineCodeLensProvider) HandleCodeLensRequest(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	var lenses []lsp.CodeLens
	for node := range core.AllNodes(doc.Root) {
		switch n := node.(type) {
		case *EventImpl:
			lenses = append(lenses, lsp.CodeLens{
				Range: n.TextRange().LspRange(doc.TextDoc),
				Data:  codeLensData{URI: string(params.TextDocument.URI), Kind: "event", Name: n.Name()},
			})
		case *CommandImpl:
			lenses = append(lenses, lsp.CodeLens{
				Range: n.TextRange().LspRange(doc.TextDoc),
				Data:  codeLensData{URI: string(params.TextDocument.URI), Kind: "command", Name: n.Name()},
			})
		}
	}

	return lenses, nil
}

func (p *stateMachineCodeLensProvider) HandleCodeLensResolveRequest(ctx context.Context, lens *lsp.CodeLens) (*lsp.CodeLens, error) {
	raw, err := json.Marshal(lens.Data)
	if err != nil {
		return lens, nil
	}
	var data codeLensData
	if err := json.Unmarshal(raw, &data); err != nil {
		return lens, nil
	}

	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	doc := documentManager.Get(core.ParseURI(data.URI))
	if doc == nil || doc.Root == nil {
		return lens, nil
	}

	target := findEventOrCommandByName(doc, data.Kind, data.Name)
	if target == nil {
		return lens, nil
	}

	referencesFinder := service.MustGet[server.ReferencesFinder](p.sc)
	count := countReferences(ctx, referencesFinder, target)
	lens.Command = &lsp.Command{Title: fmt.Sprintf("%d reference(s)", count)}
	return lens, nil
}

// findEventOrCommandByName re-locates an Event or Command by name within
// doc. HandleCodeLensResolveRequest only receives the Data payload set
// earlier by HandleCodeLensRequest, not the original AST node, so the node
// must be found again. Event/Command names are unique within a single
// statemachine document, so matching by name is sufficient.
func findEventOrCommandByName(doc *core.Document, kind, name string) core.AstNode {
	for node := range core.AllNodes(doc.Root) {
		switch n := node.(type) {
		case *EventImpl:
			if kind == "event" && n.Name() == name {
				return n
			}
		case *CommandImpl:
			if kind == "command" && n.Name() == name {
				return n
			}
		}
	}
	return nil
}

func countReferences(ctx context.Context, finder server.ReferencesFinder, target core.AstNode) int {
	count := 0
	for range finder.Find(ctx, target, server.FindReferencesOptions{IncludeDeclaration: false}) {
		count++
	}
	return count
}
