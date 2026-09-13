// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package grammarhover provides a server.HoverProvider for fastbelt's own
// self-hosted grammar language (internal/grammar): the same
// documentation-comment content as server.DefaultHoverProvider, plus - for
// grammar rule nodes - a railroad syntax diagram of the rule's concrete
// syntax.
//
// It exists as its own package - rather than living in package server - for
// two reasons. First, railroad diagrams only ever make sense for a
// language whose AST has a "grammar rule" concept, i.e. only
// internal/grammar; no other language server built with fastbelt has any
// use for it, so it doesn't belong in the generic framework package.
// Second, it couldn't live in server even if that were desirable:
// internal/railroad already imports internal/grammar (to map grammar AST
// nodes to diagram nodes), and internal/grammar's own test files import
// package server (to exercise the real language server) - so server
// importing internal/grammar or internal/railroad, directly or transitively,
// would close an import cycle for those tests. This package is the one
// place allowed to depend on all three, since it's imported only by the
// grammar language server's entrypoint (internal/grammar/server) and by
// this package's own tests.
package grammarhover

import (
	"context"
	"encoding/base64"
	"fmt"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/internal/railroad"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// HoverProvider is the grammar language's implementation of
// [server.HoverProvider].
type HoverProvider struct {
	sc *service.Container
}

func NewHoverProvider(sc *service.Container) server.HoverProvider {
	return &HoverProvider{sc: sc}
}

func (p *HoverProvider) HandleHoverRequest(ctx context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	target, sourceRange, ok := server.ResolveHoverTarget(ctx, p.sc, params)
	if !ok {
		return nil, nil
	}

	content := ""
	if docProvider, err := service.Get[server.DocumentationProvider](p.sc); err == nil {
		content = docProvider.Documentation(target)
	}
	if diagram := diagramMarkdown(ctx, target); diagram != "" {
		if content != "" {
			content += "\n\n---\n\n"
		}
		content += diagram
	}
	if content == "" {
		return nil, nil
	}

	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: content,
		},
		Range: sourceRange,
	}, nil
}

// diagramMarkdown returns a Markdown image embedding a base64-encoded SVG
// railroad diagram for node, or "" if node isn't a grammar rule with a
// meaningful concrete-syntax diagram (e.g. a terminal/token rule).
func diagramMarkdown(ctx context.Context, node core.AstNode) string {
	rule, ok := node.(grammar.AbstractRule)
	if !ok {
		return ""
	}
	diagram, ok := railroad.BuildRuleDiagram(ctx, rule)
	if !ok {
		return ""
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(diagram.SVG()))
	return fmt.Sprintf("![Railroad diagram](data:image/svg+xml;base64,%s)", encoded)
}
