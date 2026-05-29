// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"unicode"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/lsp"
)

// CompletionAcceptor is the callback adopters invoke (zero, one, or many
// times) to emit a completion item from a contributor method. Items
// passed to it are auto-enriched by the provider - any zero-valued
// fields are filled with sensible per-stage defaults derived from the
// contributor's parameters (tt / d / sn). Mirrors Langium's
// CompletionAcceptor.
type CompletionAcceptor func(item lsp.CompletionItem)

// CompletionContributor is the language-supplied policy for generating
// completion items. The completion provider drives three stages -
// token, cross-reference candidate, snippet - and asks the contributor
// once per emission what items to surface. The default contributor
// reproduces today's behavior byte-for-byte; adopters override
// individual methods to inject documentation, drop entries, surface
// terminal tokens, emit multiple variants, etc.
type CompletionContributor interface {
	// CompletionForToken is called once per TokenType the simulator
	// reports as valid at the cursor. atnState is the source ATN state
	// index - the fine-grained grammar position; two cursors inside
	// the same rule at different positions yield different values.
	//
	// The acceptor may be called zero or more times. The default
	// contributor calls accept(lsp.CompletionItem{}) when
	// tt.IsKeyword() and does nothing otherwise. Override to surface
	// terminal tokens (e.g. an in-progress ID), token groups,
	// or to specialize per atnState.
	CompletionForToken(ctx context.Context, tt *core.TokenType, atnState int, cc ContributorContext, accept CompletionAcceptor)

	// CompletionForReference is called once per cross-reference
	// candidate the scope/dispatch pass produced. hint.Field
	// identifies the reference field (e.g. "Transition.Event");
	// atnState is the source ATN state of the rule transition that
	// carried the hint. The default emits one item per d with all
	// fields auto-filled.
	CompletionForReference(ctx context.Context, d *core.SymbolDescription, hint *parser.CompletionHint, atnState int, cc ContributorContext, accept CompletionAcceptor)

	// CompletionForSnippet is called once per snippet whose
	// Applicable predicate already passed. Snippets aren't
	// ATN-driven, so no atnState parameter. The default emits one
	// item.
	CompletionForSnippet(ctx context.Context, sn SnippetTemplate, cc ContributorContext, accept CompletionAcceptor)

	// PostProcess runs once per emitted item after enrichment and
	// deduplication. Returning false drops the item; mutating *item
	// rewrites it in place.
	PostProcess(ctx context.Context, item *lsp.CompletionItem, cc ContributorContext) bool
}

// ContributorContext is the read-only view a contributor sees on each
// stage call. Mirrors Langium's CompletionContext - exposes the
// document, cursor offset, the synthetic AST node representing the
// rule being completed at the cursor, and the per-cursor REPLACE
// range / sort rank the framework would use by default.
type ContributorContext struct {
	Doc    *core.Document
	Cursor int

	// Node is the synthetic AST node representing the rule being
	// completed at the cursor, built by buildSyntheticOwnerChain from
	// the parser's RuleStack. May be nil at the document root or when
	// the rule stack is empty.
	Node core.AstNode

	// ReplaceRange is the LSP range items should REPLACE
	// (nil = INSERT at cursor with no replacement).
	ReplaceRange *lsp.Range

	// SortRank seeds the items' SortText prefix; 0 = complete-current,
	// 1 = complete-next.
	SortRank int
}

// DefaultCompletionContributor reproduces today's completion behavior:
// only keyword tokens surface, references and snippets always surface
// with default item shapes, and no post-process filtering happens.
// Adopters typically embed this and override the methods they need.
type DefaultCompletionContributor struct{}

// NewDefaultCompletionContributor returns a CompletionContributor that
// matches today's default item generation behavior.
func NewDefaultCompletionContributor() CompletionContributor {
	return &DefaultCompletionContributor{}
}

// CompletionForToken emits a default item only for keyword tokens that contains a letter or digit.
// This prevents trivial punctuation tokens from appearing in the default completion list.
func (*DefaultCompletionContributor) CompletionForToken(_ context.Context, tt *core.TokenType, _ int, _ ContributorContext, accept CompletionAcceptor) {
	if tt != nil && tt.IsKeyword() {
		for _, r := range tt.Name {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				accept(lsp.CompletionItem{})
				break
			}
		}
	}
}

// CompletionForReference emits one default item per cross-reference candidate.
func (*DefaultCompletionContributor) CompletionForReference(_ context.Context, _ *core.SymbolDescription, _ *parser.CompletionHint, _ int, _ ContributorContext, accept CompletionAcceptor) {
	accept(lsp.CompletionItem{})
}

// CompletionForSnippet emits one default item per applicable snippet.
func (*DefaultCompletionContributor) CompletionForSnippet(_ context.Context, _ SnippetTemplate, _ ContributorContext, accept CompletionAcceptor) {
	accept(lsp.CompletionItem{})
}

// PostProcess keeps every item unchanged.
func (*DefaultCompletionContributor) PostProcess(_ context.Context, _ *lsp.CompletionItem, _ ContributorContext) bool {
	return true
}
