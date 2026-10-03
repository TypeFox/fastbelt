// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package completion_test

import (
	"context"
	"fmt"
	"iter"
	"strings"
	"testing"
	"unique"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/languages/completion"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// The Completion grammar is a hand-crafted test bed for the completion
// engine. Each rule targets one simulator/parser edge case and its keyword is
// the lowercase rule name; see the comments in completion.fb.

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
	for _, want := range []string{"declare", "seq", "alt", "prefix", "call", "fqn", "list", "ref"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q at entry; got %v", want, itemLabels(items))
		}
	}
}

// Straight-line continuation: only "first" follows "seq".
func TestCompletion_AfterSeq(t *testing.T) {
	items := completionAt(t, "seq <|cursor>")
	if !hasLabel(items, "first") {
		t.Errorf("expected 'first' after 'seq'; got %v", itemLabels(items))
	}
	if hasLabel(items, "second") {
		t.Errorf("did not expect 'second' after 'seq'; got %v", itemLabels(items))
	}
	if hasLabel(items, "a") {
		t.Errorf("did not expect 'seq' to repeat; got %v", itemLabels(items))
	}
}

// Flat alternative: after "alt" both branches must stay live.
func TestCompletion_AfterAlt(t *testing.T) {
	items := completionAt(t, "alt <|cursor>")
	if !hasLabel(items, "first") {
		t.Errorf("expected 'first' after 'alt'; got %v", itemLabels(items))
	}
	if !hasLabel(items, "second") {
		t.Errorf("expected 'second' after 'alt'; got %v", itemLabels(items))
	}
}

// Shared-prefix entry: only "common" is valid before disambiguation.
func TestCompletion_AfterPrefix(t *testing.T) {
	items := completionAt(t, "prefix <|cursor>")
	if !hasLabel(items, "common") {
		t.Errorf("expected 'common' after 'prefix'; got %v", itemLabels(items))
	}
	if hasLabel(items, "first") || hasLabel(items, "second") {
		t.Errorf("did not expect branch tokens before 'common'; got %v", itemLabels(items))
	}
}

// Headline regression: both Prefix branches stay live past the shared prefix.
func TestCompletion_AfterCCommon(t *testing.T) {
	items := completionAt(t, "prefix common <|cursor>")
	if !hasLabel(items, "first") {
		t.Errorf("expected 'first' after 'prefix common'; got %v", itemLabels(items))
	}
	if !hasLabel(items, "second") {
		t.Errorf("expected 'second' after 'prefix common'; got %v", itemLabels(items))
	}
}

// End-of-rule pop-back: Root loop re-enters after "a first".
func TestCompletion_AfterAFirst(t *testing.T) {
	items := completionAt(t, "seq first <|cursor>")
	for _, want := range []string{"declare", "seq", "alt", "prefix", "call", "fqn", "list", "ref"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q after 'seq first'; got %v", want, itemLabels(items))
		}
	}
}

// Rule-call alternative entry: both CallLong/CallShort share "common".
func TestCompletion_AfterCall(t *testing.T) {
	items := completionAt(t, "call <|cursor>")
	if !hasLabel(items, "common") {
		t.Errorf("expected 'common' after 'call'; got %v", itemLabels(items))
	}
	if hasLabel(items, "then") || hasLabel(items, "long") {
		t.Errorf("did not expect later CallLong tokens; got %v", itemLabels(items))
	}
}

// Rule-call shared-prefix regression: CallLong's "then" must stay live
// even though CallShort is already a complete parse.
func TestCompletion_AfterDCommon(t *testing.T) {
	items := completionAt(t, "call common <|cursor>")
	if !hasLabel(items, "then") {
		t.Errorf("expected 'then' after 'call common'; got %v", itemLabels(items))
	}
}

// Single-path stretch: only "long" follows "d common then".
func TestCompletion_AfterDCommonThen(t *testing.T) {
	items := completionAt(t, "call common then <|cursor>")
	if !hasLabel(items, "long") {
		t.Errorf("expected 'long'; got %v", itemLabels(items))
	}
	if hasLabel(items, "common") {
		t.Errorf("did not expect 'common' to repeat; got %v", itemLabels(items))
	}
}

// Cross-reference entry with empty scope - Root keywords must not leak
// (the ID atom inherits the RefFQN.Ref hint and HintedOnlyIDs suppresses it).
func TestCompletion_AfterRefFQN(t *testing.T) {
	items := completionAt(t, "fqn <|cursor>")
	for _, leaked := range []string{"declare", "seq", "alt", "prefix", "call", "fqn", "list", "ref"} {
		if hasLabel(items, leaked) {
			t.Errorf("did not expect Root keyword %q mid-E; got %v", leaked, itemLabels(items))
		}
	}
}

// Cross-reference dispatch surfaces every Declare in scope.
func TestCompletion_AfterRefFQN_WithDeclares(t *testing.T) {
	items := completionAt(t, "declare foo declare bar fqn <|cursor>")
	for _, want := range []string{"foo", "bar"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as RefFQN.Ref candidate; got %v", want, itemLabels(items))
		}
	}
}

// Multi-segment FQN names surface as a single composite label.
func TestCompletion_AfterRefFQN_WithFQNDeclare(t *testing.T) {
	items := completionAt(t, "declare foo.bar fqn <|cursor>")
	if !hasLabel(items, "foo.bar") {
		t.Errorf("expected 'foo.bar'; got %v", itemLabels(items))
	}
	if hasLabel(items, ".") {
		t.Errorf("did not expect '.' as a standalone keyword; got %v", itemLabels(items))
	}
}

// Mid-FQN cursor: composite candidate still surfaces as one label.
func TestCompletion_AfterRefFQN_WithFQNDeclare_InFQN(t *testing.T) {
	items := completionAt(t, "declare foo.bar fqn foo<|cursor>")
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
		return
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
		return
	}
	if declare.TextEdit == nil {
		t.Errorf("expected REPLACE TextEdit on full-keyword item; got nil")
	}
}

// Dispatch enumerates every Declare, not just the first match.
func TestCompletion_AfterRefFQN_MultipleDeclares(t *testing.T) {
	items := completionAt(t, "declare foo.bar declare foo.baz fqn <|cursor>")
	for _, want := range []string{"foo.bar", "foo.baz"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as RefFQN.Ref candidate; got %v", want, itemLabels(items))
		}
	}
}

// Cursor right after the FQN separator dot.
func TestCompletion_AfterRefFQN_FQNTrailingDot(t *testing.T) {
	items := completionAt(t, "declare foo.bar fqn foo.<|cursor>")
	if !hasLabel(items, "foo.bar") {
		t.Errorf("expected 'foo.bar' after trailing dot; got %v", itemLabels(items))
	}
}

// Cursor inside the second FQN segment.
func TestCompletion_AfterRefFQN_FQNMidSegment(t *testing.T) {
	items := completionAt(t, `
declare foo.bar
fqn foo.<|cursor>`)
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
	textEdit := *edit.TextEdit
	// The edit should replace the full "foo." segment
	if textEdit.NewText != "foo.bar" {
		t.Errorf("expected replacement text 'foo.bar'; got %q", textEdit.NewText)
	}
	if textEdit.Range.Start.Character != 4 || textEdit.Range.End.Character != 8 {
		t.Errorf("expected replacement range {4,8}; got {%d,%d}", textEdit.Range.Start.Character, textEdit.Range.End.Character)
	}
}

// Error recovery: a stray token before the cursor must not erase the
// completions the user would expect at the cursor. Without recovery in the
// completion parser, the parse stops at the first mismatch and the simulator
// has no snapshot near the cursor to drive completions from.
func TestCompletion_AfterRefFQN_WithSyntaxErrorMidPrefix(t *testing.T) {
	items := completionAt(t, "declare foo bar fqn <|cursor>")
	for _, want := range []string{"foo"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as RefFQN.Ref candidate despite stray 'bar'; got %v", want, itemLabels(items))
		}
	}
}

// Rule re-entry through the Root loop: cursor inside Alt after a complete Seq.
func TestCompletion_MidRootSequence(t *testing.T) {
	items := completionAt(t, "seq first alt <|cursor>")
	for _, want := range []string{"first", "second"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q after 'seq first alt'; got %v", want, itemLabels(items))
		}
	}
}

// Plain-ID cross-reference path (no FQN composite involved).
func TestCompletion_AfterRefID_SimpleRef(t *testing.T) {
	items := completionAt(t, `
	declare alpha
	declare beta
	ref <|cursor>
	`)
	for _, want := range []string{"alpha", "beta"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as RefID.Ref candidate; got %v", want, itemLabels(items))
		}
	}
	if hasLabel(items, ".") {
		t.Errorf("did not expect '.' in plain-ID cross-ref result; got %v", itemLabels(items))
	}
}

// Cursor-in-token heuristic for plain-ID cross-refs.
func TestCompletion_AfterRefID_SimpleRef_Partial(t *testing.T) {
	items := completionAt(t, "declare alpha ref <|cursor>")
	if itemWithLabel(items, "alpha") == nil {
		t.Fatalf("expected 'alpha' for partial simple-ID cross-ref; got %v", itemLabels(items))
	}
}

// Use local scope for first member call completion
func TestCompletion_AfterMember_Simple(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		member <|cursor>
	`)
	assert.Len(t, items, 1)
	if itemWithLabel(items, "alpha") == nil {
		t.Fatalf("expected 'alpha' for first member call completion; got %v", itemLabels(items))
	}
}

// Use previous member's scope for subsequent member call completion.
// Also tests that actions are properly evaluated to populate the scope.
func TestCompletion_AfterMember_MemberCall(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		member alpha.beta.<|cursor>
	`)
	// Two contexts: the just-typed "." can be replaced (REPLACE-shaped
	// "." suggestion) AND the next MemberCall.Ref can be inserted
	// (INSERT-shaped "gamma"). Nothing else should leak in - the
	// ReplaceRange filter prunes the Root-loop keywords (seq, alt, prefix, ...)
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
func TestCompletion_AfterMemberNoDot_MemberCall(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		nodot alpha beta <|cursor>
	`)
	// Do not assert length, as the rule could end at this point
	// and surface all the other rule start keywords
	if itemWithLabel(items, "gamma") == nil {
		t.Fatalf("expected 'gamma' for member call completion; got %v", itemLabels(items))
	}
}

// Same as previous test, but with the cursor before a member call segment
func TestCompletion_AfterMemberNoDot_MemberCallWithExisting(t *testing.T) {
	items := completionAt(t, `
		declare alpha {
			declare beta {
				declare gamma
			}
		}
		nodot alpha beta <|cursor>gamma
	`)
	if itemWithLabel(items, "gamma") == nil {
		t.Fatalf("expected 'gamma' for member call completion; got %v", itemLabels(items))
	}
}

// Mixed alternative: at the cursor after "choice", both the cross-reference
// candidates and the literal keyword in the sibling branch must surface.
func TestCompletion_AfterRefOrKeyword_RefAndKeyword(t *testing.T) {
	items := completionAt(t, "declare foo choice <|cursor>")
	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo' as RefOrKeyword.Ref candidate; got %v", itemLabels(items))
	}
	if !hasLabel(items, "self") {
		t.Errorf("expected 'self' keyword alternative; got %v", itemLabels(items))
	}
}

// Mixed alternative without any declared symbols: the keyword branch must
// still surface even when the ref branch contributes nothing.
func TestCompletion_AfterRefOrKeyword_KeywordWithoutDeclares(t *testing.T) {
	items := completionAt(t, "declare some choice <|cursor>")
	assert.Len(t, items, 2)
	if !hasLabel(items, "self") {
		t.Errorf("expected 'self' keyword alternative; got %v", itemLabels(items))
	}
	if !hasLabel(items, "some") {
		t.Errorf("expected 'some' as RefOrKeyword.Ref candidate; got %v", itemLabels(items))
	}
}

// Two alternatives both assign a ref to Declare. The same candidate must
// appear once, not once per alternative.
func TestCompletion_AfterDedup_NoDuplicates(t *testing.T) {
	items := completionAt(t, "declare foo dedup <|cursor>")
	if got := countLabel(items, "foo"); got != 1 {
		t.Errorf("expected 'foo' exactly once across Dedup alternatives; got %d in %v", got, itemLabels(items))
	}
}

// Fully optional prefix: at the cursor after "opt", both the optional's
// opener ("anno") and the required follow-up ("doc") past the skipped
// group must surface.
func TestCompletion_AfterOpt_OptionalSkipped(t *testing.T) {
	items := completionAt(t, "opt <|cursor>")
	assert.Len(t, items, 2)
	for _, want := range []string{"optional", "then"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q after 'opt'; got %v", want, itemLabels(items))
		}
	}
}

// Once the optional group is entered, only its continuation is valid -
// the required follow-up ("then") must not leak past the unfinished group.
func TestCompletion_AfterOpt_OptionalEntered(t *testing.T) {
	items := completionAt(t, "opt optional <|cursor>")
	assert.Len(t, items, 1)
	if !hasLabel(items, "and") {
		t.Errorf("expected 'and' after 'opt optional'; got %v", itemLabels(items))
	}
	if hasLabel(items, "then") {
		t.Errorf("did not expect 'then' mid-optional; got %v", itemLabels(items))
	}
}

// After completing the optional group, the required follow-up must
// resurface.
func TestCompletion_AfterOpt_OptionalCompleted(t *testing.T) {
	items := completionAt(t, "opt optional and <|cursor>")
	assert.Len(t, items, 1)
	if !hasLabel(items, "then") {
		t.Errorf("expected 'then' after completed optional; got %v", itemLabels(items))
	}
	if hasLabel(items, "optional") {
		t.Errorf("did not expect 'optional' to repeat; got %v", itemLabels(items))
	}
}

// Even though the token group following the "group" keyword contains
// some keywords, none of them should be proposed by default.
func TestCompletion_AfterGroup_NoTokenGroupProposal(t *testing.T) {
	items := completionAt(t, "group <|cursor>")
	if len(items) != 0 {
		t.Errorf("expected no completion items for empty token group; got %v", itemLabels(items))
	}
}

// We use a token group as a cross reference terminal
// It should propose the names as usual
func TestCompletion_AfterRefGroup_TokenGroupCrossRef(t *testing.T) {
	items := completionAt(t, "declare some refgroup <|cursor>")
	assert.Len(t, items, 1)
	if !hasLabel(items, "some") {
		t.Errorf("expected 'some' as RefGroup.Ref candidate; got %v", itemLabels(items))
	}
}

// The token group cross-reference must enumerate the full scope, not just the
// first match, and each candidate must surface exactly once.
func TestCompletion_AfterRefGroup_TokenGroupCrossRef_MultipleDeclares(t *testing.T) {
	items := completionAt(t, "declare foo declare bar refgroup <|cursor>")
	for _, want := range []string{"foo", "bar"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as RefGroup.Ref candidate; got %v", want, itemLabels(items))
		}
		if got := countLabel(items, want); got != 1 {
			t.Errorf("expected %q exactly once; got %d in %v", want, got, itemLabels(items))
		}
	}
}

// We use a token group as a cross reference terminal
// It should propose the names as usual
func TestCompletion_AfterRefAction_ActionCrossRef(t *testing.T) {
	items := completionAt(t, "declare some action <|cursor>")
	assert.Len(t, items, 1)
	if !hasLabel(items, "some") {
		t.Errorf("expected 'some' as RefAction.Ref candidate; got %v", itemLabels(items))
	}
}

// The reference is completed on an existing RefItem node, so the scope is
// computed from the actual AST and contains the local symbols of Scope.
func TestCompletion_AfterScope_LocalSymbols_ExistingOwner(t *testing.T) {
	items := completionAt(t, "scope { declare local use local<|cursor> }")
	if !hasLabel(items, "local") {
		t.Errorf("expected 'local' as RefItem.Ref candidate; got %v", itemLabels(items))
	}
}

// The RefItem node does not exist yet, so the reference is completed on a
// synthetic owner. Its scope must still contain the local symbols of the
// enclosing Scope node.
func TestCompletion_AfterScope_LocalSymbols_SyntheticOwner(t *testing.T) {
	sources := []string{
		"declare global scope { declare local use <|cursor>",
		"declare global scope { declare local use <|cursor> }",
		// Syntactically valid: the cursor is in front of the existing reference.
		"declare global scope { declare local use <|cursor>local }",
	}
	for _, src := range sources {
		t.Run(src, func(t *testing.T) {
			items := completionAt(t, src)
			for _, want := range []string{"global", "local"} {
				if !hasLabel(items, want) {
					t.Errorf("expected %q as RefItem.Ref candidate; got %v", want, itemLabels(items))
				}
			}
		})
	}
}

// ownerRecordingFilter records the owner of every reference of the rules Scope
// to Retype that is completed.
type ownerRecordingFilter struct {
	completion.DefaultCompletionCompletionFilter
	owners []core.AstNode
}

func (f *ownerRecordingFilter) FilterWrapRefRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

func (f *ownerRecordingFilter) FilterNestRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

func (f *ownerRecordingFilter) FilterAmbigNameRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

func (f *ownerRecordingFilter) FilterAmbigRefsRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

func (f *ownerRecordingFilter) FilterRefItemRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

func (f *ownerRecordingFilter) FilterChainItemRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

func (f *ownerRecordingFilter) FilterOperandRef(_ context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	f.owners = append(f.owners, ref.Owner())
	return in
}

// completeOwners requests the completion at the cursor and returns the
// owners of the references that were completed.
func completeOwners(t *testing.T, src string) (*test.Doc, []core.AstNode) {
	t.Helper()
	filter := &ownerRecordingFilter{}
	sc := service.NewContainer()
	completion.SetupServices(sc)
	service.Override[completion.CompletionCompletionFilter](sc, filter)
	sc.Seal()

	doc := test.New(t, sc).Parse(src)
	doc.CompletionItems("cursor")
	if len(filter.owners) == 0 {
		t.Fatalf("expected the completion of a reference")
	}
	return doc, filter.owners
}

// describeOwner returns the types and the containment data of the owner and
// of its containers below the root, starting with the owner. The nodes that
// are not part of the document are marked as new.
func describeOwner(doc *test.Doc, owner core.AstNode) string {
	existing := collections.NewSet[core.AstNode]()
	for node := range core.AllNodes(doc.Root()) {
		existing.Add(node)
	}
	segments := []string{}
	for node := owner; node != nil && node != doc.Root(); node = node.Container() {
		segment := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%T", node), "*completion."), "Impl")
		if !existing.Has(node) {
			segment = "new " + segment
		}
		field, index := node.ContainmentData()
		if field != (unique.Handle[string]{}) && field.Value() != "" {
			segment += " at " + field.Value()
			if index >= 0 {
				segment += fmt.Sprintf("@%d", index)
			}
		}
		segments = append(segments, segment)
	}
	return strings.Join(segments, ", ")
}

// assertOwners requests the completion at the cursor of each source and
// compares the owners of the references with the expected descriptions, see
// describeOwner. The root contains a single node.
func assertOwners(t *testing.T, cases map[string][]string) {
	t.Helper()
	for src, expected := range cases {
		t.Run(src, func(t *testing.T) {
			doc, owners := completeOwners(t, src)
			actual := []string{}
			for _, owner := range owners {
				actual = append(actual, describeOwner(doc, owner))
			}
			assert.ElementsMatch(t, expected, actual)
		})
	}
}

// The synthetic owner must know the property and the index that it would
// have in the enclosing Scope node, as scope providers may depend on them.
func TestCompletion_AfterScope_SyntheticOwnerContainment(t *testing.T) {
	assertOwners(t, map[string][]string{
		"scope { declare local use <|cursor>":                            {"new RefItem at item, Scope at objects@0"},
		"scope { declare local use <|cursor>local }":                     {"new RefItem at item, Scope at objects@0"},
		"scope { declare local use lo<|cursor>":                          {"new RefItem at item, Scope at objects@0"},
		"scope { declare local use local and <|cursor>":                  {"new RefItem at others@0, Scope at objects@0"},
		"scope { declare local use local and <|cursor>local and local }": {"new RefItem at others@0, Scope at objects@0"},
		"scope { declare local use local and local and <|cursor>":        {"new RefItem at others@1, Scope at objects@0"},
		"scope { declare local use local and local and <|cursor>local }": {"new RefItem at others@1, Scope at objects@0"},
		"scope { declare local use local and local and lo<|cursor>":      {"new RefItem at others@1, Scope at objects@0"},
	})
}

// The items of Loop have no separator, so the parser does not enter the loop
// for an item that is missing. The owner is only known from the rule calls
// that lead to the cross-reference.
func TestCompletion_AfterLoop_SyntheticOwnerContainment(t *testing.T) {
	assertOwners(t, map[string][]string{
		"loop { <|cursor>":               {"new RefItem at items@0, Loop at objects@0"},
		"loop { <|cursor> }":             {"new RefItem at items@0, Loop at objects@0"},
		"loop { declare local <|cursor>": {"new RefItem at items@0, Loop at objects@0"},
		"loop { foo <|cursor>":           {"new RefItem at items@1, Loop at objects@0"},
		"loop { foo <|cursor>foo }":      {"new RefItem at items@1, Loop at objects@0"},
		"loop { foo foo <|cursor>":       {"new RefItem at items@2, Loop at objects@0"},
		// The cursor is behind the nested Loop, so the owner belongs to the outer one
		"loop { foo loop { foo foo } <|cursor>":   {"new RefItem at items@1, Loop at objects@0"},
		"loop { foo loop { foo foo } <|cursor> }": {"new RefItem at items@1, Loop at objects@0"},
		"loop { loop { loop { } <|cursor> } }":    {"new RefItem at items@0, Loop at nested@0, Loop at objects@0"},
		// The cursor is inside of the nested Loop
		"loop { foo loop { foo foo <|cursor>":     {"new RefItem at items@2, Loop at nested@0, Loop at objects@0"},
		"loop { foo loop { foo foo <|cursor> } }": {"new RefItem at items@2, Loop at nested@0, Loop at objects@0"},
	})
}

// A tree-rewriting action wraps the items of Chain, also for the input that
// follows the cursor. The owner of a new item still belongs to the Chain.
func TestCompletion_AfterChain_SyntheticOwnerContainment(t *testing.T) {
	assertOwners(t, map[string][]string{
		"chain { <|cursor>":               {"new ChainItem at items@0, Chain at objects@0"},
		"chain { foo and bar <|cursor>":   {"new ChainItem at items@1, Chain at objects@0"},
		"chain { foo and bar <|cursor> }": {"new ChainItem at items@1, Chain at objects@0"},
		// The items in front of the cursor are wrapped because of the input that follows it
		"chain { foo <|cursor> and bar }":             {"new ChainItem at items@1, Chain at objects@0"},
		"chain { foo and bar <|cursor> and baz }":     {"new ChainItem at items@1, Chain at objects@0"},
		"chain { foo foo and bar <|cursor> and baz }": {"new ChainItem at items@2, Chain at objects@0"},
	})
}

// The operands of the infix rule Binary are contained in the node of their
// operator, which depends on the precedence of the operators.
func TestCompletion_AfterInfix_SyntheticOwnerContainment(t *testing.T) {
	assertOwners(t, map[string][]string{
		"infix { <|cursor>": {"new Operand at items@0, Infix at objects@0"},
		// The owner is the right operand of the operator in front of the cursor
		"infix { foo plus <|cursor>":                {"new Operand at right, Binary at items@0, Infix at objects@0"},
		"infix { foo plus bar times <|cursor>":      {"new Operand at right, Binary at right, Binary at items@0, Infix at objects@0"},
		"infix { foo times bar plus <|cursor>":      {"new Operand at right, Binary at items@0, Infix at objects@0"},
		"infix { foo plus <|cursor>bar times baz }": {"new Operand at right, Binary at items@0, Infix at objects@0"},
		"infix { foo times <|cursor>bar plus baz }": {"new Operand at right, Binary at left, Binary at items@0, Infix at objects@0"},
		// The owner is a new item behind the operands
		"infix { foo plus bar times baz <|cursor>":   {"new Operand at items@1, Infix at objects@0"},
		"infix { foo times bar plus baz <|cursor> }": {"new Operand at items@1, Infix at objects@0"},
		// The operands in front of the cursor belong to operators that follow it
		"infix { foo <|cursor> plus bar }":               {"new Operand at items@1, Infix at objects@0"},
		"infix { foo plus bar <|cursor> times baz }":     {"new Operand at items@1, Infix at objects@0"},
		"infix { foo foo times bar <|cursor> plus baz }": {"new Operand at items@2, Infix at objects@0"},
	})
}

// The tree-rewriting action of WrapGroup wraps the node in front of the cursor
// into a list, and the owner is the next item of that list. The main parser
// only created the list if another item follows the cursor.
func TestCompletion_AfterWrap_OwnerInActionList(t *testing.T) {
	assertOwners(t, map[string][]string{
		"wrap { <|cursor>":          {"new WrapRef at item, Wrap at objects@0"},
		"wrap { foo <|cursor>":      {"new WrapRef at elements@1, new WrapGroup at item, Wrap at objects@0"},
		"wrap { foo <|cursor> }":    {"new WrapRef at elements@1, new WrapGroup at item, Wrap at objects@0"},
		"wrap { foo <|cursor>bar }": {"new WrapRef at elements@1, WrapGroup at item, Wrap at objects@0"},
		"wrap { foo bar <|cursor>":  {"new WrapRef at elements@2, WrapGroup at item, Wrap at objects@0"},
		// The cursor is behind the nested list
		"wrap { { foo bar } <|cursor>":       {"new WrapRef at elements@1, new WrapGroup at item, Wrap at objects@0"},
		"wrap { baz { foo bar } <|cursor> }": {"new WrapRef at elements@2, WrapGroup at item, Wrap at objects@0"},
		// The cursor is inside of the nested list
		"wrap { baz { foo <|cursor>":       {"new WrapRef at elements@1, new WrapGroup at elements@1, WrapGroup at item, Wrap at objects@0"},
		"wrap { baz { foo bar <|cursor> }": {"new WrapRef at elements@2, WrapGroup at elements@1, WrapGroup at item, Wrap at objects@0"},
	})
}

// The property Right of Shadow has the name of an operand of the infix rule that
// it calls, which must not be confused when leaving the operands.
func TestCompletion_AfterShadow_OwnerBehindInfixRule(t *testing.T) {
	assertOwners(t, map[string][]string{
		"shadow foo <|cursor>":                    {"new RefItem at items@0, Shadow at objects@0"},
		"shadow foo plus bar times baz <|cursor>": {"new RefItem at items@0, Shadow at objects@0"},
		"shadow foo times bar plus baz <|cursor>": {"new RefItem at items@0, Shadow at objects@0"},
		"shadow foo plus bar foo <|cursor>":       {"new RefItem at items@1, Shadow at objects@0"},
		"shadow foo plus <|cursor>":               {"new Operand at right, Binary at right, Shadow at objects@0"},
	})
}

// The owner of the reference in Nest is the Nest that the cursor is in, and not a
// nested Nest that ends in front of the cursor.
func TestCompletion_AfterNest_ExistingOwner(t *testing.T) {
	assertOwners(t, map[string][]string{
		"nest { <|cursor>":                     {"Nest at objects@0"},
		"nest { nest { } <|cursor>":            {"Nest at objects@0"},
		"nest { nest { nest { } } <|cursor> }": {"Nest at objects@0"},
		"nest { nest { nest { } <|cursor> } }": {"Nest at children@0, Nest at objects@0"},
		"nest { nest { } nest { <|cursor> } }": {"Nest at children@1, Nest at objects@0"},
	})
	items := completionAt(t, "declare global nest { declare outer nest { declare inner } <|cursor> }")
	for _, want := range []string{"global", "outer"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as Nest.Ref candidate; got %v", want, itemLabels(items))
		}
	}
	if hasLabel(items, "inner") {
		t.Errorf("did not expect 'inner' as Nest.Ref candidate; got %v", itemLabels(items))
	}
}

// The main parser reads the token in front of the cursor as the name of a
// AmbigName. The reference of AmbigRefs reads it as a reference, so its owner is
// unrelated to the node of the main parser.
func TestCompletion_AfterAmbig_OwnerOfOtherAlternative(t *testing.T) {
	doc, owners := completeOwners(t, "ambig foo <|cursor>bar first")
	if _, ok := test.FindNode[completion.Ambig](doc); !ok {
		t.Fatalf("expected the main parser to create an AmbigName")
	}
	actual := []string{}
	for _, owner := range owners {
		actual = append(actual, describeOwner(doc, owner))
	}
	assert.ElementsMatch(t, []string{
		"AmbigName at name, Ambig at objects@0",
		"new AmbigRefs, new Root",
	}, actual)
}

// The action of RetypeItem changes the type of the node that contains the owner.
func TestCompletion_AfterRetype_OwnerInNodeOfAction(t *testing.T) {
	assertOwners(t, map[string][]string{
		"retype { <|cursor>":     {"new RefItem at inner, new RetypeWrapper at items@0, Retype at objects@0"},
		"retype { foo <|cursor>": {"new RefItem at inner, new RetypeWrapper at items@1, Retype at objects@0"},
	})
}

// The cursor is further away from the start of the nested Loop than the
// simulator looks back, so it has to leave a rule that it didn't enter.
func TestCompletion_AfterLoop_LongInput(t *testing.T) {
	src := "declare global loop { declare outer loop { declare inner " + strings.Repeat("foo ", 40) + "} <|cursor>"
	items := completionAt(t, src)
	for _, want := range []string{"global", "outer"} {
		if !hasLabel(items, want) {
			t.Errorf("expected %q as RefItem.Ref candidate; got %v", want, itemLabels(items))
		}
	}
	if hasLabel(items, "inner") {
		t.Errorf("did not expect 'inner' as RefItem.Ref candidate; got %v", itemLabels(items))
	}
	assertOwners(t, map[string][]string{src: {"new RefItem at items@0, Loop at objects@1"}})
}

// The local symbols of a nested Loop are not visible behind it.
func TestCompletion_AfterLoop_LocalSymbols(t *testing.T) {
	cases := []struct {
		src      string
		expected []string
		excluded []string
	}{
		{"declare global loop { declare outer loop { declare inner } <|cursor>", []string{"global", "outer"}, []string{"inner"}},
		{"declare global loop { declare outer loop { declare inner } <|cursor> }", []string{"global", "outer"}, []string{"inner"}},
		{"declare global loop { declare outer loop { declare inner <|cursor>", []string{"global", "outer", "inner"}, nil},
		{"declare global loop { declare outer loop { declare inner } } <|cursor>", nil, []string{"global", "outer", "inner"}},
	}
	for _, c := range cases {
		t.Run(c.src, func(t *testing.T) {
			items := completionAt(t, c.src)
			for _, want := range c.expected {
				if !hasLabel(items, want) {
					t.Errorf("expected %q as RefItem.Ref candidate; got %v", want, itemLabels(items))
				}
			}
			for _, unwanted := range c.excluded {
				if hasLabel(items, unwanted) {
					t.Errorf("did not expect %q as RefItem.Ref candidate; got %v", unwanted, itemLabels(items))
				}
			}
		})
	}
}

// A token group used as a cross-reference terminal must resolve like any other
// cross-reference, not just drive completion.
func TestTokenGroupReference_Resolves(t *testing.T) {
	sc := completion.CreateServices(&SimpleCompletionContributor{})
	doc := test.New(t, sc).Parse("declare foo refgroup foo")
	doc.AssertNoErrors()

	ref := test.MustFindReferenceWithText[completion.Declare](doc, "foo")
	target := ref.Ref(doc.Ctx())
	assert.Nil(t, ref.Error())
	if assert.NotNil(t, target) {
		assert.Equal(t, "foo", target.Name())
	}
}

// The reference parses through the token group but points at no symbol, so it
// must report a resolution error.
func TestTokenGroupReference_Unresolved(t *testing.T) {
	sc := completion.CreateServices(&SimpleCompletionContributor{})
	doc := test.New(t, sc).Parse("refgroup missing")

	ref := test.MustFindReferenceWithText[completion.Declare](doc, "missing")
	assert.Nil(t, ref.Ref(doc.Ctx()))
	assert.NotNil(t, ref.Error())
}

type hidingCompletionFilter struct {
	completion.DefaultCompletionCompletionFilter
	hide string
}

func (h *hidingCompletionFilter) FilterRefFQNRef(ctx context.Context, ref *core.Reference[completion.Declare], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {
	return func(yield func(*core.SymbolDescription) bool) {
		for d := range in {
			if d.Unit.String() == h.hide {
				continue
			}
			if !yield(d) {
				return
			}
		}
	}
}

// FilterRefFQNRef override hides one declare while siblings remain.
func TestCompletion_FilterOverride(t *testing.T) {
	sc := service.NewContainer()
	completion.SetupServices(sc)
	service.Override[completion.CompletionCompletionFilter](sc, &hidingCompletionFilter{hide: "bar"})
	sc.Seal()

	doc := test.New(t, sc).Parse("declare foo declare bar fqn <|cursor>")
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
				docs := "Declares a named symbol."
				item.Documentation = &lsp.CompletionItemDocumentation{String: &docs}
			}
			accept(item)
		},
	}
	items := completionAtWith(t, "<|cursor>", contrib)

	declare := itemWithLabel(items, "declare")
	if declare == nil {
		t.Fatalf("expected 'declare'; got %v", itemLabels(items))
		return
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

// When a contributor opts in and accepts an item for the token group, that
// item must surface - the "no proposal by default" behaviour is purely a
// default, not a hard suppression.
func TestCompletion_ContributorAcceptsTokenGroup(t *testing.T) {
	contrib := &recordingContributor{
		onToken: func(tt *core.TokenType, _ int, _ server.ContributorContext, accept server.CompletionAcceptor) {
			if tt.Name == "SomeTokenGroup" {
				accept(lsp.CompletionItem{Label: "SomeTokenGroup"})
			}
		},
	}
	items := completionAtWith(t, "group <|cursor>", contrib)
	assert.Len(t, items, 1)
	if !hasLabel(items, "SomeTokenGroup") {
		t.Errorf("expected accepted token group item to surface; got %v", itemLabels(items))
	}
}

// Reference hook receives hint.Field, atnState, and a synthetic owner.
// cc.Node must be non-nil despite [Root, RefFQN, FQN] containing a composite
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
				name:     d.Unit.String(),
				field:    hint.Field,
				atnState: atnState,
				node:     cc.Node,
			})
			accept(lsp.CompletionItem{})
		},
	}
	items := completionAtWith(t, "declare foo fqn <|cursor>", contrib)

	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo'; got %v", itemLabels(items))
	}
	if len(observations) == 0 {
		t.Fatalf("expected at least one reference observation")
	}
	for _, o := range observations {
		if o.field != "RefFQN.Ref" {
			t.Errorf("expected hint.Field=\"RefFQN.Ref\"; got %+v", o)
		}
		if o.atnState <= 0 {
			t.Errorf("expected positive atnState; got %+v", o)
		}
		if o.node == nil {
			t.Errorf("expected cc.Node non-nil (composite skip); got %+v", o)
		} else if _, ok := o.node.(completion.RefFQN); !ok {
			t.Errorf("expected cc.Node to be a synthetic RefFQN; got %T", o.node)
		}
	}
}

// Multi-level synthetic chain: cc.Node lands on RefListItem despite [Root, RefList]
// not yet containing the RefListItem frame (cursor sits where one could begin).
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
	items := completionAtWith(t, "declare foo list <|cursor>", contrib)

	if !hasLabel(items, "foo") {
		t.Errorf("expected 'foo'; got %v", itemLabels(items))
	}
	if len(observations) == 0 {
		t.Fatalf("expected at least one reference observation")
	}
	for _, o := range observations {
		if o.field != "RefListItem.Ref" {
			t.Errorf("expected hint.Field=\"RefListItem.Ref\"; got %+v", o)
		}
		if _, ok := o.node.(completion.RefListItem); !ok {
			t.Errorf("expected cc.Node to be a synthetic RefListItem; got %T", o.node)
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
	c := itemWithLabel(items, "call")
	if c == nil {
		t.Fatalf("expected 'prefix' keyword; got %v", itemLabels(items))
		return
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
