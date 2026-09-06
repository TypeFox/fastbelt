// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package test

import (
	"context"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/languages/completion"
	"typefox.dev/fastbelt/server"
	"typefox.dev/lsp"
)

// simpleCompletionContributor surfaces all keywords for completion tests.
type simpleCompletionContributor struct {
	server.DefaultCompletionContributor
}

func (c *simpleCompletionContributor) CompletionForToken(ctx context.Context, tt *core.TokenType, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor) {
	if tt != nil && tt.IsKeyword() {
		accept(lsp.CompletionItem{})
	}
}

// TestExpectCompletion_basics demonstrates the ExpectCompletion helper.
// Mirrors patterns from completion_engine_test.go but uses the new
// chainable assertion API similar to ExpectDiagnostic.
func TestExpectCompletion_basics(t *testing.T) {
	sc := completion.CreateServices(&simpleCompletionContributor{})
	f := New(t, sc)

	// At entry, Root keywords surface as completion items.
	doc := f.Parse("<|cursor>")
	doc.ExpectCompletion("cursor").
		Has("declare").
		Has("a").
		Has("b").
		Has("c").
		Has("d").
		Has("e").
		Has("f").
		Has("g").
		HasCount(16)
}

// TestExpectCompletion_kindFiltering demonstrates HasKind for LSP kind assertions.
func TestExpectCompletion_kindFiltering(t *testing.T) {
	sc := completion.CreateServices(&simpleCompletionContributor{})
	f := New(t, sc)

	// Keywords should surface as Keyword completion items.
	doc := f.Parse("a <|cursor>")
	doc.ExpectCompletion("cursor").
		Has("first").
		HasKind("first", lsp.KeywordCompletion)
}

// TestExpectCompletion_noItems verifies no items when token group has no default proposals.
func TestExpectCompletion_noItems(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	// Token group "m" produces no default proposals.
	doc := f.Parse("m <|cursor>")
	items := doc.CompletionItems("cursor")
	if len(items) != 0 {
		t.Fatalf("expected no items for token group, got %d", len(items))
	}
}

// TestExpectCompletion_afterDeclares verifies cross-reference completion.
func TestExpectCompletion_afterDeclares(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	doc := f.Parse("declare foo declare bar e <|cursor>")
	doc.ExpectCompletion("cursor").
		Has("foo").
		Has("bar").
		HasCount(2)
}

// TestExpectCompletion_fqnComposite verifies FQN composite labels.
func TestExpectCompletion_fqnComposite(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	doc := f.Parse("declare foo.bar e <|cursor>")
	doc.ExpectCompletion("cursor").
		Has("foo.bar").
		HasCount(1)
}

// TestExpectCompletion_memberCall verifies member call completion chains.
func TestExpectCompletion_memberCall(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	doc := f.Parse(`
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		h alpha.beta.<|cursor>
	`)
	// gamma is the expected completion item.
	doc.ExpectCompletion("cursor").
		Has("gamma").
		HasCount(1)
}

// TestExpectCompletion_optionalPrefix verifies optional group handling.
func TestExpectCompletion_optionalPrefix(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	doc := f.Parse("l <|cursor>")
	doc.ExpectCompletion("cursor").
		Has("optional").
		Has("then").
		HasCount(2)
}

// TestExpectCompletion_labelPrefix verifies WithLabelPrefix filtering.
func TestExpectCompletion_labelPrefix(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	doc := f.Parse("declare alpha declare beta g <|cursor>")
	doc.ExpectCompletion("cursor").
		Has("alpha").
		Has("beta").
		WithLabelPrefix("")
}

// TestExpectCompletion_hasKind_correct verifies that HasKind passes when kind matches.
func TestExpectCompletion_hasKind_correct(t *testing.T) {
	sc := completion.CreateServices(nil)
	f := New(t, sc)

	// "foo" is a symbol reference (snake_case), expected kind is Reference.
	doc := f.Parse("declare foo e <|cursor>")
	doc.ExpectCompletion("cursor").
		Has("foo").
		HasKind("foo", lsp.ReferenceCompletion)
}
