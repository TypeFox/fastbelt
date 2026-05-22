// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/parser"
)

// SnippetTemplate describes a multi-token completion (e.g. typing "state"
// inserts an entire `state ... end` block). The body is rendered in LSP
// snippet syntax so the client expands tab stops on accept.
type SnippetTemplate struct {
	// Label is the entry displayed in the completion list.
	Label string
	// Detail is the short right-aligned hint shown next to Label, if set.
	Detail string
	// Documentation is the long description rendered when the entry is
	// focused. It is sent as plain text; clients that prefer Markdown can
	// pre-format the content.
	Documentation string
	// Body is the snippet payload inserted on accept. LSP snippet syntax
	// is supported: `${1:placeholder}`, `${0}`, etc. See the LSP spec for
	// the full grammar.
	Body string
	// Applicable, when non-nil, gates whether the snippet appears for a
	// given cursor position. Default-nil snippets are always offered.
	Applicable func(ctx SnippetContext) bool
}

// SnippetContext is the predicate-side view of the cursor: the document and
// byte offset the request lands at, the live ATN TokenTypes the simulator
// reported (so snippets can require a specific keyword to be valid), and
// the CompletionParser's rule stack (so snippets can require, e.g.,
// "we're inside a State rule body").
type SnippetContext struct {
	Doc        *core.Document
	Cursor     int
	TokenTypes []*core.TokenType
	RuleStack  []parser.RuleContext
}

// SnippetRegistry is where adopters register their language's snippets. The
// default implementation is an in-memory slice; production deployments can
// override it to load snippets from a file/database/etc.
type SnippetRegistry interface {
	Register(s SnippetTemplate)
	All() []SnippetTemplate
}

// DefaultSnippetRegistry is the default in-memory implementation.
type DefaultSnippetRegistry struct {
	snippets []SnippetTemplate
}

// NewDefaultSnippetRegistry returns an empty registry.
func NewDefaultSnippetRegistry() SnippetRegistry {
	return &DefaultSnippetRegistry{}
}

// Register adds a snippet to the registry. Snippets are returned by All in
// registration order; the completion provider further sorts the resulting
// CompletionItems via SortText.
func (r *DefaultSnippetRegistry) Register(s SnippetTemplate) {
	r.snippets = append(r.snippets, s)
}

// All returns the registered snippets.
func (r *DefaultSnippetRegistry) All() []SnippetTemplate {
	return r.snippets
}

// RequiresValidToken returns a SnippetTemplate.Applicable predicate that
// fires only when the simulator says the given TokenType is among the valid
// next tokens at the cursor. This is the common pattern for "snippet that
// starts with keyword X":
//
//	snippets.Register(server.SnippetTemplate{
//	    Label:      "state block",
//	    Body:       "state ${1:name}\n  ${0}\nend",
//	    Applicable: server.RequiresValidToken(Keyword_state),
//	})
func RequiresValidToken(t *core.TokenType) func(SnippetContext) bool {
	if t == nil {
		return func(SnippetContext) bool { return false }
	}
	return func(ctx SnippetContext) bool {
		for _, tt := range ctx.TokenTypes {
			if tt != nil && tt.Id == t.Id {
				return true
			}
		}
		return false
	}
}
