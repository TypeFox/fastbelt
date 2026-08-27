// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

func TestTypeDefinitionProvider_CustomImplementation(t *testing.T) {
	sc := service.NewContainer()
	grammar.SetupServices(sc)
	server.SetupDefaultServices(sc)

	// Register custom provider:
	// Use Put, not Override, since there is no default to override.
	service.Put[server.TypeDefinitionProvider](sc, &grammarTypeDefinitionProvider{sc: sc})
	sc.Seal()

	f := test.New(t, sc)
	doc := f.Parse(`
		grammar Test;
		interface Person { Name string }
		Person: Name=ID;
	` + commonTokens)
	doc.AssertNoErrors()

	provider := service.MustGet[server.TypeDefinitionProvider](f.Services())

	// "Person" appears twice: once as the interface name, once as the rule
	// name ("Person: Name=ID;"). Find the rule occurrence specifically by
	// matching "Person:" (only the rule declaration is followed by a colon).
	text := doc.Document.TextDoc.Text(nil)
	position := doc.Document.TextDoc.PositionAt(indexOf(text, "Person:"))

	result, err := provider.HandleTypeDefinitionRequest(
		context.Background(),
		&lsp.TypeDefinitionParams{
			TextDocumentPositionParams: lsp.TextDocumentPositionParams{
				TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
				Position:     position,
			},
		},
	)
	assert.NoError(t, err)
	if assert.Len(t, result, 1) {
		// The rule's returns type is the "Person" interface, declared earlier.
		wantPosition := doc.Document.TextDoc.PositionAt(indexOf(text, "interface Person") + len("interface "))
		assert.Equal(t, wantPosition, result[0].TargetSelectionRange.Start)
	}
}

func indexOf(text, needle string) int {
	for i := 0; i+len(needle) <= len(text); i++ {
		if text[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

// grammarTypeDefinitionProvider resolves a grammar parser rule to its return
// type interface (explicit "returns" clause, or implicit same-name
// interface).
type grammarTypeDefinitionProvider struct {
	sc *service.Container
}

func (p *grammarTypeDefinitionProvider) HandleTypeDefinitionRequest(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error) {
	documentManager := service.MustGet[workspace.DocumentManager](p.sc)
	doc := documentManager.Get(core.ParseURI(string(params.TextDocument.URI)))
	if doc == nil || doc.Root == nil {
		return nil, nil
	}

	node := server.NodeAtCursor(doc, params.Position)
	rule, ok := node.(*grammar.ParserRuleImpl)
	if !ok {
		return nil, nil
	}

	iface := grammar.FindReturnType(rule, ctx)
	if iface == nil {
		return nil, nil
	}

	targetDoc := iface.Document().TextDoc
	sourceRange := rule.NameToken().TextRange().LspRange(doc.TextDoc)
	fullRange := iface.TextRange().LspRange(targetDoc)
	nameRange := iface.NameToken().TextRange().LspRange(targetDoc)

	return []lsp.DefinitionLink{
		{
			OriginSelectionRange: &sourceRange,
			TargetURI:            targetDoc.URI(),
			TargetRange:          fullRange,
			TargetSelectionRange: nameRange,
		},
	}, nil
}
