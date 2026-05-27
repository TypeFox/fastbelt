// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package completion_test

import (
	"context"
	"iter"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/languages/completion"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// The Completion grammar is a hand-crafted test bed for the completion
// engine. Each rule targets one simulator/parser edge case:
//
//   Root:   repeating entry rule
//   Declare: exports symbols into scope (FQN composite name)
//   A:      straight-line continuation
//   B:      flat alternative - both branches must stay live post-commit
//   C:      shared-prefix alternative - keep both alts past the prefix
//   D:      rule-call alternative with shared prefix (DLong/DShort)
//   E:      cross-reference via FQN composite
//   F:      nested rule with FItem children (synthetic chain coverage)
//   G:      cross-reference via plain ID
//   J:      cross-reference and keyword alternative at the same position
//   K:      two assignments to the same target type - dedup check
//   L:      fully optional group in front of mandatory content

func completionAt(t *testing.T, src string) []lsp.CompletionItem {
	t.Helper()
	return completionAtWith(t, src, &SimpleCompletionContributor{})
}

type SimpleCompletionContributor struct {
	server.DefaultCompletionContributor
}

func (c *SimpleCompletionContributor) CompletionForToken(ctx context.Context, tt *core.TokenType, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor) {
	if tt != nil && tt.IsKeyword() {
		accept(lsp.CompletionItem{})
	}
}

func completionAtWith(t *testing.T, src string, contributor server.CompletionContributor) []lsp.CompletionItem {
	t.Helper()
	sc := completion.CreateServices(contributor)
	doc := test.New(t, sc).Parse(src)
	return doc.CompletionItems("cursor")
}

func hasLabel(items []lsp.CompletionItem, label string) bool {
	for _, it := range items {
		if it.Label == label {
			return true
		}
	}
	return false
}

func itemLabels(items []lsp.CompletionItem) []string {
	out := make([]string, 0, len(items))
	for _, it := range items {
		out = append(out, it.Label)
	}
	return out
}

func itemWithLabel(items []lsp.CompletionItem, label string) *lsp.CompletionItem {
	for i := range items {
		if items[i].Label == label {
			return &items[i]
		}
	}
	return nil
}

func countLabel(items []lsp.CompletionItem, label string) int {
	n := 0
	for _, it := range items {
		if it.Label == label {
			n++
		}
	}
	return n
}

// Entry: every Root alternative starter must surface.
func TestCompletion_AtEntry(t *testing.T) {
	items := completionAt(t, "<|cursor>")
	for _, want := range []string{"declare", "a", "b", "c", "d", "e", "f", "g"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q at entry; got %v", want, itemLabels(items))
		}
	}
}

// Straight-line continuation: only "first" follows "a".
func TestCompletion_AfterA(t *testing.T) {
	items := completionAt(t, "a <|cursor>")
	if !hasLabel(items, "first") {
		t.Errorf("expected 'first' after 'a'; got %v", itemLabels(items))
	}
	if hasLabel(items, "second") {
		t.Errorf("did not expect 'second' after 'a'; got %v", itemLabels(items))
	}
	if hasLabel(items, "a") {
		t.Errorf("did not expect 'a' to repeat; got %v", itemLabels(items))
	}
}

// Flat alternative: after "b" both branches must stay live.
func TestCompletion_AfterB(t *testing.T) {
	items := completionAt(t, "b <|cursor>")
	if !hasLabel(items, "first") {
		t.Errorf("expected 'first' after 'b'; got %v", itemLabels(items))
	}
	if !hasLabel(items, "second") {
		t.Errorf("expected 'second' after 'b'; got %v", itemLabels(items))
	}
}

// Shared-prefix entry: only "common" is valid before disambiguation.
func TestCompletion_AfterC(t *testing.T) {
	items := completionAt(t, "c <|cursor>")
	if !hasLabel(items, "common") {
		t.Errorf("expected 'common' after 'c'; got %v", itemLabels(items))
	}
	if hasLabel(items, "first") || hasLabel(items, "second") {
		t.Errorf("did not expect branch tokens before 'common'; got %v", itemLabels(items))
	}
}

// Headline regression: both C branches stay live past the shared prefix.
func TestCompletion_AfterCCommon(t *testing.T) {
	items := completionAt(t, "c common <|cursor>")
	if !hasLabel(items, "first") {
		t.Errorf("expected 'first' after 'c common'; got %v", itemLabels(items))
	}
	if !hasLabel(items, "second") {
		t.Errorf("expected 'second' after 'c common'; got %v", itemLabels(items))
	}
}

// End-of-rule pop-back: Root loop re-enters after "a first".
func TestCompletion_AfterAFirst(t *testing.T) {
	items := completionAt(t, "a first <|cursor>")
	for _, want := range []string{"declare", "a", "b", "c", "d", "e", "f", "g"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q after 'a first'; got %v", want, itemLabels(items))
		}
	}
}

// Rule-call alternative entry: both DLong/DShort share "common".
func TestCompletion_AfterD(t *testing.T) {
	items := completionAt(t, "d <|cursor>")
	if !hasLabel(items, "common") {
		t.Errorf("expected 'common' after 'd'; got %v", itemLabels(items))
	}
	if hasLabel(items, "then") || hasLabel(items, "long") {
		t.Errorf("did not expect later DLong tokens; got %v", itemLabels(items))
	}
}

// Rule-call shared-prefix regression: DLong's "then" must stay live
// even though DShort is already a complete parse.
func TestCompletion_AfterDCommon(t *testing.T) {
	items := completionAt(t, "d common <|cursor>")
	if !hasLabel(items, "then") {
		t.Errorf("expected 'then' after 'd common'; got %v", itemLabels(items))
	}
}

// Single-path stretch: only "long" follows "d common then".
func TestCompletion_AfterDCommonThen(t *testing.T) {
	items := completionAt(t, "d common then <|cursor>")
	if !hasLabel(items, "long") {
		t.Errorf("expected 'long'; got %v", itemLabels(items))
	}
	if hasLabel(items, "common") {
		t.Errorf("did not expect 'common' to repeat; got %v", itemLabels(items))
	}
}

// Cross-reference entry with empty scope - Root keywords must not leak
// (the ID atom inherits the E.Ref hint and HintedOnlyIDs suppresses it).
func TestCompletion_AfterE(t *testing.T) {
	items := completionAt(t, "e <|cursor>")
	for _, leaked := range []string{"declare", "a", "b", "c", "d", "e", "f", "g"} {
		if hasLabel(items, leaked) {
			t.Errorf("did not expect Root keyword %q mid-E; got %v", leaked, itemLabels(items))
		}
	}
}

// Cross-reference dispatch surfaces every Declare in scope.
func TestCompletion_AfterE_WithDeclares(t *testing.T) {
	items := completionAt(t, "declare foo declare bar e <|cursor>")
	for _, want := range []string{"foo", "bar"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as E.Ref candidate; got %v", want, itemLabels(items))
		}
	}
}

// Multi-segment FQN names surface as a single composite label.
func TestCompletion_AfterE_WithFQNDeclare(t *testing.T) {
	items := completionAt(t, "declare foo.bar e <|cursor>")
	if !hasLabel(items, "foo.bar") {
		t.Errorf("expected 'foo.bar'; got %v", itemLabels(items))
	}
	if hasLabel(items, ".") {
		t.Errorf("did not expect '.' as a standalone keyword; got %v", itemLabels(items))
	}
}

// Mid-FQN cursor: composite candidate still surfaces as one label.
func TestCompletion_AfterE_WithFQNDeclare_InFQN(t *testing.T) {
	items := completionAt(t, "declare foo.bar e foo<|cursor>")
	if !hasLabel(items, "foo.bar") {
		t.Errorf("expected 'foo.bar'; got %v", itemLabels(items))
	}
	if hasLabel(items, ".") {
		t.Errorf("did not expect '.' as a standalone keyword; got %v", itemLabels(items))
	}
}

// Cursor strictly inside a partial keyword - item must be REPLACE-shaped.
func TestCompletion_PartialKeyword(t *testing.T) {
	items := completionAt(t, "decla<|cursor>")
	declare := itemWithLabel(items, "declare")
	if declare == nil {
		t.Fatalf("expected 'declare'; got %v", itemLabels(items))
	}
	if declare.TextEdit == nil {
		t.Errorf("expected REPLACE TextEdit on partial-keyword item; got nil")
	}
}

// Cursor strictly inside a full keyword - item must be REPLACE-shaped.
func TestCompletion_InsideKeyword(t *testing.T) {
	items := completionAt(t, "decla<|cursor>re")
	declare := itemWithLabel(items, "declare")
	if declare == nil {
		t.Fatalf("expected 'declare'; got %v", itemLabels(items))
	}
	if declare.TextEdit == nil {
		t.Errorf("expected REPLACE TextEdit on full-keyword item; got nil")
	}
}

// Dispatch enumerates every Declare, not just the first match.
func TestCompletion_AfterE_MultipleDeclares(t *testing.T) {
	items := completionAt(t, "declare foo.bar declare foo.baz e <|cursor>")
	for _, want := range []string{"foo.bar", "foo.baz"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as E.Ref candidate; got %v", want, itemLabels(items))
		}
	}
}

// Cursor right after the FQN separator dot.
func TestCompletion_AfterE_FQNTrailingDot(t *testing.T) {
	items := completionAt(t, "declare foo.bar e foo.<|cursor>")
	if !hasLabel(items, "foo.bar") {
		t.Errorf("expected 'foo.bar' after trailing dot; got %v", itemLabels(items))
	}
}

// Cursor inside the second FQN segment.
func TestCompletion_AfterE_FQNMidSegment(t *testing.T) {
	items := completionAt(t, `
declare foo.bar
e foo.<|cursor>`)
	item := itemWithLabel(items, "foo.bar")
	if item == nil {
		t.Errorf("expected 'foo.bar' for partial FQN mid-segment; got %v", itemLabels(items))
		return
	}
	if item.TextEdit == nil {
		t.Errorf("expected REPLACE TextEdit for partial FQN mid-segment; got nil")
		return
	}
	edit := item.TextEdit
	textEdit := edit.Value.(lsp.TextEdit)
	// The edit should replace the full "foo." segment
	if textEdit.NewText != "foo.bar" {
		t.Errorf("expected replacement text 'foo.bar'; got %q", textEdit.NewText)
	}
	if textEdit.Range.Start.Character != 2 || textEdit.Range.End.Character != 6 {
		t.Errorf("expected replacement range {2,6}; got {%d,%d}", textEdit.Range.Start.Character, textEdit.Range.End.Character)
	}
}

// Error recovery: a stray token before the cursor must not erase the
// completions the user would expect at the cursor. Without recovery in the
// completion parser, the parse stops at the first mismatch and the simulator
// has no snapshot near the cursor to drive completions from.
func TestCompletion_AfterE_WithSyntaxErrorMidPrefix(t *testing.T) {
	items := completionAt(t, "declare foo bar e <|cursor>")
	for _, want := range []string{"foo"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as E.Ref candidate despite stray 'bar'; got %v", want, itemLabels(items))
		}
	}
}

// Rule re-entry through the Root loop: cursor inside B after a complete A.
func TestCompletion_MidRootSequence(t *testing.T) {
	items := completionAt(t, "a first b <|cursor>")
	for _, want := range []string{"first", "second"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q after 'a first b'; got %v", want, itemLabels(items))
		}
	}
}

// Plain-ID cross-reference path (no FQN composite involved).
func TestCompletion_AfterG_SimpleRef(t *testing.T) {
	items := completionAt(t, `
	declare alpha
	declare beta
	g <|cursor>
	`)
	for _, want := range []string{"alpha", "beta"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as G.Ref candidate; got %v", want, itemLabels(items))
		}
	}
	if hasLabel(items, ".") {
		t.Errorf("did not expect '.' in plain-ID cross-ref result; got %v", itemLabels(items))
	}
}

// Cursor-in-token heuristic for plain-ID cross-refs.
func TestCompletion_AfterG_SimpleRef_Partial(t *testing.T) {
	items := completionAt(t, "declare alpha g <|cursor>")
	if itemWithLabel(items, "alpha") == nil {
		t.Fatalf("expected 'alpha' for partial simple-ID cross-ref; got %v", itemLabels(items))
	}
}

// Use local scope for first member call completion
func TestCompletion_AfterH_Simple(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		h <|cursor>
	`)
	assert.Len(t, items, 1)
	if itemWithLabel(items, "alpha") == nil {
		t.Fatalf("expected 'alpha' for first member call completion; got %v", itemLabels(items))
	}
}

// Use previous member's scope for subsequent member call completion.
// Also tests that actions are properly evaluated to populate the scope.
func TestCompletion_AfterH_MemberCall(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		h alpha.beta.<|cursor>
	`)
	// Two contexts: the just-typed "." can be replaced (REPLACE-shaped
	// "." suggestion) AND the next MemberCall.Ref can be inserted
	// (INSERT-shaped "gamma"). Nothing else should leak in - the
	// ReplaceRange filter prunes the Root-loop keywords (a, b, c, ...)
	// that are theoretically valid substitutes for "." but don't match
	// what the user actually typed.
	assert.Len(t, items, 2)
	dot := itemWithLabel(items, ".")
	gamma := itemWithLabel(items, "gamma")
	if dot == nil || gamma == nil {
		t.Fatalf("expected '.' and 'gamma'; got %v", itemLabels(items))
	}
	if dot.TextEdit == nil {
		t.Errorf("expected '.' to carry a REPLACE TextEdit")
	}
	if gamma.TextEdit != nil {
		t.Errorf("expected 'gamma' to be INSERT-shaped (no TextEdit); got %+v", gamma.TextEdit)
	}
}

// Use previous member's scope for subsequent member call completion.
// Also tests that actions are properly evaluated to populate the scope.
// This test covers the non-dot member call syntax, ensuring that the same
func TestCompletion_AfterI_MemberCall(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		i alpha beta <|cursor>
	`)
	// Do not assert length, as the rule could end at this point
	// and surface all the other rule start keywords
	if itemWithLabel(items, "gamma") == nil {
		t.Fatalf("expected 'gamma' for member call completion; got %v", itemLabels(items))
	}
}

// Same as previous test, but with the cursor before a member call segment
func TestCompletion_AfterI_MemberCallWithExisting(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		i alpha beta <|cursor>gamma
	`)
	if itemWithLabel(items, "gamma") == nil {
		t.Fatalf("expected 'gamma' for member call completion; got %v", itemLabels(items))
	}
}

// Mixed alternative: at the cursor after "j", both the cross-reference
// candidates and the literal keyword in the sibling branch must surface.
func TestCompletion_AfterJ_RefAndKeyword(t *testing.T) {
	items := completionAt(t, "declare foo j <|cursor>")
	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo' as J.Ref candidate; got %v", itemLabels(items))
	}
	if !hasLabel(items, "self") {
		t.Errorf("expected 'self' keyword alternative; got %v", itemLabels(items))
	}
}

// Mixed alternative without any declared symbols: the keyword branch must
// still surface even when the ref branch contributes nothing.
func TestCompletion_AfterJ_KeywordWithoutDeclares(t *testing.T) {
	items := completionAt(t, "declare some j <|cursor>")
	assert.Len(t, items, 2)
	if !hasLabel(items, "self") {
		t.Errorf("expected 'self' keyword alternative; got %v", itemLabels(items))
	}
	if !hasLabel(items, "some") {
		t.Errorf("expected 'some' as J.Ref candidate; got %v", itemLabels(items))
	}
}

// Two alternatives both assign a ref to Declare. The same candidate must
// appear once, not once per alternative.
func TestCompletion_AfterK_NoDuplicates(t *testing.T) {
	items := completionAt(t, "declare foo k <|cursor>")
	if got := countLabel(items, "foo"); got != 1 {
		t.Errorf("expected 'foo' exactly once across K alternatives; got %d in %v", got, itemLabels(items))
	}
}

// Fully optional prefix: at the cursor after "l", both the optional's
// opener ("anno") and the required follow-up ("doc") past the skipped
// group must surface.
func TestCompletion_AfterL_OptionalSkipped(t *testing.T) {
	items := completionAt(t, "l <|cursor>")
	assert.Len(t, items, 2)
	for _, want := range []string{"optional", "then"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q after 'l'; got %v", want, itemLabels(items))
		}
	}
}

// Once the optional group is entered, only its continuation is valid -
// the required follow-up ("then") must not leak past the unfinished group.
func TestCompletion_AfterL_OptionalEntered(t *testing.T) {
	items := completionAt(t, "l optional <|cursor>")
	assert.Len(t, items, 1)
	if !hasLabel(items, "and") {
		t.Errorf("expected 'and' after 'l optional'; got %v", itemLabels(items))
	}
	if hasLabel(items, "then") {
		t.Errorf("did not expect 'then' mid-optional; got %v", itemLabels(items))
	}
}

// After completing the optional group, the required follow-up must
// resurface.
func TestCompletion_AfterL_OptionalCompleted(t *testing.T) {
	items := completionAt(t, "l optional and <|cursor>")
	assert.Len(t, items, 1)
	if !hasLabel(items, "then") {
		t.Errorf("expected 'then' after completed optional; got %v", itemLabels(items))
	}
	if hasLabel(items, "optional") {
		t.Errorf("did not expect 'optional' to repeat; got %v", itemLabels(items))
	}
}

type hidingCompletionFilter struct {
	completion.DefaultCompletionCompletionFilter
	hide string
}

func (h *hidingCompletionFilter) FilterERef(ctx context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	return func(yield func(*core.SymbolDescription) bool) {
		for d := range in {
			if d.Name.String() == h.hide {
				continue
			}
			if !yield(d) {
				return
			}
		}
	}
}

// FilterERef override hides one declare while siblings remain.
func TestCompletion_FilterOverride(t *testing.T) {
	sc := service.NewContainer()
	completion.SetupServices(sc)
	service.Override[completion.CompletionCompletionFilter](sc, &hidingCompletionFilter{hide: "bar"})
	sc.Seal()

	doc := test.New(t, sc).Parse("declare foo declare bar e <|cursor>")
	items := doc.CompletionItems("cursor")

	if hasLabel(items, "bar") {
		t.Errorf("expected 'bar' to be filtered out; got %v", itemLabels(items))
	}
	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo' to remain; got %v", itemLabels(items))
	}
}

type recordingContributor struct {
	server.DefaultCompletionContributor
	onToken     func(tt *core.TokenType, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor)
	onReference func(d *core.SymbolDescription, hint *parser.CompletionHint, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor)
	postProcess func(item *lsp.CompletionItem, cc server.ContributorContext) bool
}

func (r *recordingContributor) CompletionForToken(ctx context.Context, tt *core.TokenType, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor) {
	if r.onToken != nil {
		r.onToken(tt, atnState, cc, accept)
		return
	}
	r.DefaultCompletionContributor.CompletionForToken(ctx, tt, atnState, cc, accept)
}

func (r *recordingContributor) CompletionForReference(ctx context.Context, d *core.SymbolDescription, hint *parser.CompletionHint, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor) {
	if r.onReference != nil {
		r.onReference(d, hint, atnState, cc, accept)
		return
	}
	r.DefaultCompletionContributor.CompletionForReference(ctx, d, hint, atnState, cc, accept)
}

func (r *recordingContributor) PostProcess(ctx context.Context, item *lsp.CompletionItem, cc server.ContributorContext) bool {
	if r.postProcess != nil {
		return r.postProcess(item, cc)
	}
	return r.DefaultCompletionContributor.PostProcess(ctx, item, cc)
}

// Token hook attaches Documentation; enrichment fills the rest.
func TestCompletion_ContributorTokenDocs(t *testing.T) {
	contrib := &recordingContributor{
		onToken: func(tt *core.TokenType, _ int, _ server.ContributorContext, accept server.CompletionAcceptor) {
			if !tt.IsKeyword() {
				return
			}
			item := lsp.CompletionItem{}
			if tt.Name == "declare" {
				item.Documentation = &lsp.Or_CompletionItem_documentation{Value: "Declares a named symbol."}
			}
			accept(item)
		},
	}
	items := completionAtWith(t, "<|cursor>", contrib)

	declare := itemWithLabel(items, "declare")
	if declare == nil {
		t.Fatalf("expected 'declare'; got %v", itemLabels(items))
	}
	if declare.Documentation == nil {
		t.Errorf("expected Documentation set; got nil")
	}
	if declare.Kind != lsp.KeywordCompletion {
		t.Errorf("expected Kind=KeywordCompletion; got %v", declare.Kind)
	}
	if declare.SortText == "" {
		t.Errorf("expected SortText filled; got empty")
	}
}

// Reference hook receives hint.Field, atnState, and a synthetic owner.
// cc.Node must be non-nil despite [Root, E, FQN] containing a composite
// frame - the chain builder skips frames with no synthetic factory.
func TestCompletion_ContributorReferenceBranching(t *testing.T) {
	type seen struct {
		name     string
		field    string
		atnState int
		node     core.AstNode
	}
	var observations []seen

	contrib := &recordingContributor{
		onReference: func(d *core.SymbolDescription, hint *parser.CompletionHint, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor) {
			observations = append(observations, seen{
				name:     d.Name.String(),
				field:    hint.Field,
				atnState: atnState,
				node:     cc.Node,
			})
			accept(lsp.CompletionItem{})
		},
	}
	items := completionAtWith(t, "declare foo e <|cursor>", contrib)

	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo'; got %v", itemLabels(items))
	}
	if len(observations) == 0 {
		t.Fatalf("expected at least one reference observation")
	}
	for _, o := range observations {
		if o.field != "E.Ref" {
			t.Errorf("expected hint.Field=\"E.Ref\"; got %+v", o)
		}
		if o.atnState <= 0 {
			t.Errorf("expected positive atnState; got %+v", o)
		}
		if o.node == nil {
			t.Errorf("expected cc.Node non-nil (composite skip); got %+v", o)
		} else if _, ok := o.node.(completion.E); !ok {
			t.Errorf("expected cc.Node to be a synthetic E; got %T", o.node)
		}
	}
}

// Multi-level synthetic chain: cc.Node lands on FItem despite [Root, F]
// not yet containing the FItem frame (cursor sits where one could begin).
func TestCompletion_ContributorSyntheticChain(t *testing.T) {
	type seen struct {
		field string
		node  core.AstNode
	}
	var observations []seen

	contrib := &recordingContributor{
		onReference: func(d *core.SymbolDescription, hint *parser.CompletionHint, atnState int, cc server.ContributorContext, accept server.CompletionAcceptor) {
			observations = append(observations, seen{field: hint.Field, node: cc.Node})
			accept(lsp.CompletionItem{})
		},
	}
	items := completionAtWith(t, "declare foo f <|cursor>", contrib)

	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo'; got %v", itemLabels(items))
	}
	if len(observations) == 0 {
		t.Fatalf("expected at least one reference observation")
	}
	for _, o := range observations {
		if o.field != "FItem.Ref" {
			t.Errorf("expected hint.Field=\"FItem.Ref\"; got %+v", o)
		}
		if _, ok := o.node.(completion.FItem); !ok {
			t.Errorf("expected cc.Node to be a synthetic FItem; got %T", o.node)
		}
	}
}

// PostProcess rewrites SortText and drops items by returning false.
func TestCompletion_ContributorPostProcess(t *testing.T) {
	contrib := &recordingContributor{
		postProcess: func(item *lsp.CompletionItem, _ server.ContributorContext) bool {
			if item.Label == "declare" {
				return false
			}
			if strings.HasPrefix(item.Label, "c") {
				item.SortText = "zzz-" + item.Label
			}
			return true
		},
	}
	items := completionAtWith(t, "<|cursor>", contrib)

	if hasLabel(items, "declare") {
		t.Errorf("expected 'declare' to be dropped; got %v", itemLabels(items))
	}
	c := itemWithLabel(items, "c")
	if c == nil {
		t.Fatalf("expected 'c' keyword; got %v", itemLabels(items))
	}
	if !strings.HasPrefix(c.SortText, "zzz-") {
		t.Errorf("expected rewritten SortText; got %q", c.SortText)
	}
}

// Override emits non-keyword tokens (default contributor drops them).
func TestCompletion_ContributorSurfacesTerminalTokens(t *testing.T) {
	contrib := &recordingContributor{
		onToken: func(tt *core.TokenType, _ int, _ server.ContributorContext, accept server.CompletionAcceptor) {
			accept(lsp.CompletionItem{})
		},
	}
	items := completionAtWith(t, "declare <|cursor>", contrib)
	if !hasLabel(items, "ID") {
		t.Errorf("expected 'ID' terminal token; got %v", itemLabels(items))
	}
}
