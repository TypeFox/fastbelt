// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// CompletionProvider handles LSP textDocument/completion requests. The
// default implementation drives the language-specific completion parser, the
// ATN simulator, and the per-field dispatch table to produce a flat list of
// keyword + cross-reference + snippet items.
type CompletionProvider interface {
	HandleCompletionRequest(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error)
}

// DefaultCompletionProvider is the framework's CompletionProvider
// implementation. It is generic across languages and delegates everything
// language-specific to a registered LanguageCompletionAdapter and to a
// CompletionContributor that decides how items are generated.
type DefaultCompletionProvider struct {
	sc *service.Container
}

// NewDefaultCompletionProvider returns a CompletionProvider backed by sc
// and using the DefaultCompletionContributor as the fallback when the
// service container has no CompletionContributor registered.
func NewDefaultCompletionProvider(sc *service.Container) CompletionProvider {
	return &DefaultCompletionProvider{sc: sc}
}

// HandleCompletionRequest fulfils textDocument/completion. The flow is:
//
//  1. resolve doc + cursor offset;
//  2. build 1-2 completion contexts via buildCompletionContexts -
//     "complete current token" and/or "complete next token" depending on
//     where the cursor sits relative to surrounding tokens;
//  3. for each context, run the language CompletionParser over its prefix,
//     simulate the ATN forward, gather TokenTypes + cross-reference hints,
//     and build per-context CompletionItems;
//  4. deduplicate items by (Kind, Label) - both contexts often surface
//     the same keyword at boundaries between rules;
//  5. apply the CompletionContributor's PostProcess hook.
//
// Each context's TokenTypes/hints share the same downstream pipeline
// (keyword pass, dispatch pass, snippet pass); only the prefix token
// slice and the resulting TextEdit shape differ.
func (s *DefaultCompletionProvider) HandleCompletionRequest(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
	empty := &lsp.CompletionList{IsIncomplete: false, Items: []lsp.CompletionItem{}}

	adapter, err := service.Get[parser.LanguageCompletionAdapter](s.sc)
	if err != nil || adapter == nil {
		return empty, nil
	}
	documentManager, err := service.Get[workspace.DocumentManager](s.sc)
	if err != nil {
		return empty, nil
	}

	uri := core.ParseURI(string(params.TextDocument.URI))
	doc := documentManager.Get(uri)
	if doc == nil {
		return empty, nil
	}

	offset := int(doc.TextDoc.OffsetAt(params.Position))
	atn := adapter.ATN()
	if atn == nil {
		return empty, nil
	}

	contributor := s.resolveContributor()
	matcher := s.resolveFuzzyMatcher()

	contexts := buildCompletionContexts(doc, offset)
	items := make([]lsp.CompletionItem, 0, len(contexts)*8)
	for _, cc := range contexts {
		items = append(items, s.completionsForContext(ctx, contributor, matcher, adapter, atn, doc, offset, cc)...)
	}
	items = deduplicateItems(items)

	// PostProcess pass: contributor may drop or rewrite individual items.
	postCtx := ContributorContext{Doc: doc, Cursor: offset}
	filtered := make([]lsp.CompletionItem, 0, len(items))
	for _, item := range items {
		if contributor.PostProcess(ctx, &item, postCtx) {
			filtered = append(filtered, item)
		}
	}
	items = filtered

	return &lsp.CompletionList{IsIncomplete: false, Items: items}, nil
}

// resolveContributor returns the contributor registered in the service
// container if any, otherwise falls back to the one stored on the
// provider (set by the constructor).
func (s *DefaultCompletionProvider) resolveContributor() CompletionContributor {
	if c, err := service.Get[CompletionContributor](s.sc); err == nil && c != nil {
		return c
	}
	return &DefaultCompletionContributor{}
}

// resolveFuzzyMatcher returns the FuzzyMatcher registered in the service
// container, falling back to a fresh DefaultFuzzyMatcher when the
// container has none. The fallback keeps the completion pipeline
// usable in tests/embedded setups that haven't called
// SetupDefaultServices.
func (s *DefaultCompletionProvider) resolveFuzzyMatcher() FuzzyMatcher {
	if m, err := service.Get[FuzzyMatcher](s.sc); err == nil && m != nil {
		return m
	}
	return &DefaultFuzzyMatcher{}
}

// completionsForContext runs the CompletionParser+simulator+dispatch
// pipeline for one cursor context and returns the items it produced.
// Snippets are gathered here (rather than in HandleCompletionRequest) so
// that snippet applicability predicates can inspect the simulator's
// per-context TokenTypes.
//
// Each pass builds an acceptor closure that captures the appropriate
// enrichment helper and the per-context state; the contributor decides
// whether and how often to emit, the provider enriches the emission.
func (s *DefaultCompletionProvider) completionsForContext(
	ctx context.Context,
	contributor CompletionContributor,
	matcher FuzzyMatcher,
	adapter parser.LanguageCompletionAdapter,
	atn *parser.RuntimeATN,
	doc *core.Document,
	cursorOffset int,
	cc CompletionContext,
) []lsp.CompletionItem {
	prefixTokens := doc.Tokens[:cc.PrefixLen]
	result := adapter.Parse(prefixTokens)
	live, _, ok := result.SimulateAt(atn, cc.PrefixLen)
	if !ok {
		return nil
	}
	info := atn.NextCompletionsFromSet(live)

	contribCtx := ContributorContext{
		Doc:          doc,
		Cursor:       cursorOffset,
		Node:         buildSyntheticOwnerChain(adapter, doc, result.RuleStack),
		ReplaceRange: cc.ReplaceRange,
		SortRank:     cc.SortRank,
	}

	items := make([]lsp.CompletionItem, 0, len(info.Tokens)+len(info.Hints)+4)

	// Token pass: contributor decides per (TokenType, atnState) what to emit.
	for _, tc := range info.Tokens {
		if _, hidden := info.HintedOnlyIDs[tc.TokenType.Id]; hidden {
			continue
		}
		if !matcher.Match(cc.ReplaceText, tokenLabel(tc.TokenType)) {
			continue
		}
		tcCopy := tc
		accept := func(item lsp.CompletionItem) {
			items = append(items, EnrichTokenCompletionItem(item, tcCopy.TokenType, contribCtx))
		}
		contributor.CompletionForToken(ctx, tcCopy.TokenType, tcCopy.ATNStateIdx, contribCtx, accept)
	}

	// Cross-reference pass: dispatch per CompletionHint.Field; contributor
	// decides per (SymbolDescription, hint, atnState) what to emit.
	for _, hc := range info.Hints {
		owner := buildSyntheticOwnerChainFor(adapter, doc, result.RuleStack, hc.Hint.Field, cursorOffset)
		if owner == nil {
			continue
		}
		owner = applyPrecedingAction(adapter, owner, hc.Hint)
		seq, ok := adapter.DispatchCompletion(ctx, hc.Hint.Field, owner)
		if !ok {
			continue
		}
		hcCopy := hc
		for d := range seq {
			dCopy := d
			if !matcher.Match(cc.ReplaceText, dCopy.Name.String()) {
				continue
			}
			accept := func(item lsp.CompletionItem) {
				items = append(items, EnrichReferenceCompletionItem(item, dCopy, contribCtx))
			}
			contributor.CompletionForReference(ctx, dCopy, hcCopy.Hint, hcCopy.ATNStateIdx, contribCtx, accept)
		}
	}

	// Snippet pass.
	if reg, err := service.Get[SnippetRegistry](s.sc); err == nil && reg != nil {
		tokenTypes := make([]*core.TokenType, 0, len(info.Tokens))
		seen := make(collections.Set[int], len(info.Tokens))
		for _, tc := range info.Tokens {
			if !seen.Add(tc.TokenType.Id) {
				continue
			}
			tokenTypes = append(tokenTypes, tc.TokenType)
		}
		sctx := SnippetContext{
			Doc:        doc,
			Cursor:     cursorOffset,
			TokenTypes: tokenTypes,
			RuleStack:  result.RuleStack,
		}
		for _, sn := range reg.All() {
			if sn.Applicable != nil && !sn.Applicable(sctx) {
				continue
			}
			if !matcher.Match(cc.ReplaceText, sn.Label) {
				continue
			}
			snCopy := sn
			accept := func(item lsp.CompletionItem) {
				items = append(items, EnrichSnippetCompletionItem(item, snCopy, contribCtx))
			}
			contributor.CompletionForSnippet(ctx, snCopy, contribCtx, accept)
		}
	}
	return items
}

// CompletionContext describes one cursor interpretation. A single LSP
// completion request may produce several: typically one for "complete the
// in-progress token" (REPLACE semantics) and one for "complete what could
// follow" (INSERT semantics). See buildCompletionContexts for the rules.
type CompletionContext struct {
	// PrefixLen is the number of tokens the simulator treats as committed
	// (prefixTokens = doc.Tokens[:PrefixLen]). The simulator then walks
	// SimulateAt(atn, PrefixLen) to land at the cursor.
	PrefixLen int
	// ReplaceRange, when non-nil, is the LSP range the produced items
	// should REPLACE (the in-progress token). nil means INSERT at the
	// cursor with no replacement.
	ReplaceRange *lsp.Range
	// ReplaceText is the typed text inside ReplaceRange (from the token
	// start up to the cursor). Used to filter completion items whose
	// label cannot reasonably substitute for what the user typed.
	// Empty string disables filtering for this context.
	ReplaceText string
	// SortRank seeds the items' SortText prefix; "complete current"
	// contexts use 0 so they rank above "complete next" (rank 1) in the
	// client's display order.
	SortRank int
}

// cursorTokenInfo describes the token layout around the cursor for one
// completion request. CurrentIdx is the token containing or ending at
// the cursor (-1 if the cursor is in whitespace); NextIdx is the first
// token whose Start >= offset.
type cursorTokenInfo struct {
	CurrentIdx   int
	CurrentAtEnd bool // offset == doc.Tokens[CurrentIdx].End
	NextIdx      int
}

// backtrackToToken classifies the cursor's position relative to the
// document's tokens, in terms the rest of the pipeline understands
// (token indices, not byte offsets).
//
// Token boundaries: TextSegment.Indices.Start is inclusive, End is
// exclusive. So "cursor inside a token" means Start < offset < End, and
// "cursor at the end of a token" means offset == End. A cursor at
// Start (offset == Start) sits BEFORE the token (it could still be part
// of preceding whitespace), so we treat it as a between-tokens position.
func backtrackToToken(tokens core.TokenSlice, offset int) cursorTokenInfo {
	info := cursorTokenInfo{CurrentIdx: -1, NextIdx: len(tokens)}
	// Binary search for the first token whose Start >= offset.
	lo, hi := 0, len(tokens)
	for lo < hi {
		mid := (lo + hi) / 2
		if int(tokens[mid].TextSegment.Indices.Start) < offset {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	info.NextIdx = lo
	// The token possibly covering the cursor is the one immediately before
	// NextIdx - check whether the cursor lies inside or at its end.
	if lo > 0 {
		prev := lo - 1
		end := int(tokens[prev].TextSegment.Indices.End)
		if offset < end {
			info.CurrentIdx = prev
		} else if offset == end {
			info.CurrentIdx = prev
			info.CurrentAtEnd = true
		}
	}
	return info
}

// buildCompletionContexts produces the per-request cursor interpretations.
// At most two contexts are returned; many cursor positions yield exactly
// one.
func buildCompletionContexts(doc *core.Document, offset int) []CompletionContext {
	info := backtrackToToken(doc.Tokens, offset)

	if info.CurrentIdx < 0 {
		// Cursor sits between tokens / past EOF - only "complete next".
		return []CompletionContext{{
			PrefixLen:    info.NextIdx,
			ReplaceRange: nil,
			SortRank:     0,
		}}
	}

	curToken := &doc.Tokens[info.CurrentIdx]
	replace := curToken.TextSegment.Range.LspRange()
	prefixLen := info.CurrentIdx
	typed := curToken.Image
	if !info.CurrentAtEnd {
		typed = curToken.Image[:offset-int(curToken.TextSegment.Indices.Start)]
	}

	// When the current token is part of a multi-token CompositeNode,
	// the REPLACE range must cover the whole composite span - otherwise accepting
	// a candidate would overwrite only the trailing token and leave earlier segments
	// behind, producing duplicated text. Walk back through tokens that
	// belong to the same composite and widen the replace range / prefix.
	if startIdx, ok := compositeStart(doc.Tokens, info.CurrentIdx); ok {
		startToken := &doc.Tokens[startIdx]
		// Replace start position with the start of the first token in the composite
		replace.Start = startToken.TextSegment.Range.LspRange().Start
		// Regenerate the content of the composite up to the cursor position
		var sb strings.Builder
		for i := startIdx; i < info.CurrentIdx; i++ {
			sb.WriteString(doc.Tokens[i].Image)
		}
		sb.WriteString(typed)
		typed = sb.String()
		prefixLen = startIdx
	}

	if info.CurrentAtEnd {
		// Fully-typed token: two interpretations are both legitimate.
		// Either the user wants to replace what they just typed (REPLACE
		// shape, rank 0), or they want to continue with the grammar's
		// follow-set (INSERT shape, rank 1). The downstream filter prunes
		// the replace context to items whose label can substitute for the
		// content the user actually typed.
		return []CompletionContext{
			{
				PrefixLen:    prefixLen,
				ReplaceRange: &replace,
				ReplaceText:  typed,
				SortRank:     0,
			},
			{
				PrefixLen:    info.CurrentIdx + 1,
				ReplaceRange: nil,
				SortRank:     1,
			},
		}
	}

	return []CompletionContext{{
		PrefixLen:    prefixLen,
		ReplaceRange: &replace,
		ReplaceText:  typed,
		SortRank:     0,
	}}
}

// compositeStart returns the index of the first token sharing a CompositeNode
// element with tokens[idx], or (idx, false) if tokens[idx] is not part of a
// multi-token composite.
func compositeStart(tokens core.TokenSlice, idx int) (int, bool) {
	if idx < 0 || idx >= len(tokens) {
		return idx, false
	}
	composite, ok := tokens[idx].Element.(core.CompositeNode)
	if !ok || composite == nil {
		return idx, false
	}
	start := idx
	for start > 0 {
		prev, ok := tokens[start-1].Element.(core.CompositeNode)
		if !ok || prev != composite {
			break
		}
		start--
	}
	if start == idx {
		return idx, false
	}
	return start, true
}

// tokenLabel returns the user-visible label for a TokenType, mirroring
// the default chosen by EnrichTokenCompletionItem so the pre-emission
// filter sees the same string the client would.
func tokenLabel(tt *core.TokenType) string {
	if tt == nil {
		return ""
	}
	if tt.Label != "" {
		return tt.Label
	}
	return tt.Name
}

// deduplicateItems removes duplicate CompletionItems by (Kind, Label).
// Multiple contexts can surface the same item (e.g. cursor right after a
// keyword that's also valid at the next-token position); the client
// should see each suggestion exactly once. First occurrence wins so
// "complete current" items (which are emitted first) keep their
// TextEdit replacement range.
func deduplicateItems(items []lsp.CompletionItem) []lsp.CompletionItem {
	seen := make(collections.Set[string], len(items))
	out := items[:0]
	for _, it := range items {
		key := fmt.Sprintf("%d|%s", it.Kind, it.Label)
		if !seen.Add(key) {
			continue
		}
		out = append(out, it)
	}
	return out
}

// buildSyntheticOwnerChain walks the rule stack outermost-first, instantiates
// a synthetic node per frame via the adapter, and wires Container() +
// SetDocument() so scope providers reading node.Document() or walking
// ContainerOfType[T] resolve correctly. Returns the innermost synthetic -
// the owner of the in-progress reference.
func buildSyntheticOwnerChain(adapter parser.LanguageCompletionAdapter, doc *core.Document, ruleStack []parser.RuleContext) core.AstNode {
	if len(ruleStack) == 0 {
		return nil
	}
	var parent core.AstNode
	for _, frame := range ruleStack {
		node, ok := adapter.SyntheticOwnerFor(frame.RuleKey)
		if !ok || node == nil {
			// Some rules have no associated AST node (composite nodes)
			// Those are simply skipped
			continue
		}
		if parent != nil {
			node.SetContainer(parent)
		}
		node.SetDocument(doc)
		parent = node
	}
	return parent
}

// buildSyntheticOwnerChainFor extends buildSyntheticOwnerChain with the
// hint's owner rule when the rule stack doesn't already end at that rule.
//
// At cursor positions where a new rule could begin but hasn't yet, the
// parser's rule stack stops one level above the rule the hint's field
// belongs to. The hint's Field has the form "<OwnerRule>.<Property>" -
// we split on '.' to get the owner rule name and append a synthetic for
// it if the stack doesn't already end there.
//
// Returns nil if the adapter doesn't know one of the rule keys; the
// completion request then yields no candidates for this hint rather than
// silently returning the wrong scope.
func buildSyntheticOwnerChainFor(adapter parser.LanguageCompletionAdapter, doc *core.Document, ruleStack []parser.RuleContext, hintField string, cursorOffset int) core.AstNode {
	ownerRule, _, ok := splitHintField(hintField)
	if !ok {
		return buildSyntheticOwnerChain(adapter, doc, ruleStack)
	}
	// Prefer the AST that the main parser already built. Tree-rewrite
	// actions are executed during normal parsing, so the partial AST at
	// the cursor already carries the chain a scope provider needs -
	// reusing it avoids resynthesising state the parser tracked
	// correctly.
	if real := findExistingOwnerAtCursor(adapter, doc, ownerRule, cursorOffset); real != nil {
		return real
	}
	// If the owner rule appears anywhere on the stack, the parser is
	// already inside it - slice off any deeper frames. This covers
	// cross-references whose text form is a separate rule: at the cursor
	// the rule stack might end inside a string/composite rule that
	// carries no AST node, while the owner rule (which holds the
	// reference) sits one frame higher. Searching from the top down
	// picks the innermost matching frame.
	for i := len(ruleStack) - 1; i >= 0; i-- {
		if ruleStack[i].RuleKey == ownerRule {
			return buildSyntheticOwnerChain(adapter, doc, ruleStack[:i+1])
		}
	}
	// The owner rule isn't on the stack - the parser hasn't entered it
	// yet (cursor sits at a position where it could begin). Append a
	// synthetic frame for the owner so the chain has somewhere to
	// dispatch on.
	extended := make([]parser.RuleContext, 0, len(ruleStack)+1)
	extended = append(extended, ruleStack...)
	extended = append(extended, parser.RuleContext{RuleKey: ownerRule})
	return buildSyntheticOwnerChain(adapter, doc, extended)
}

// applyPrecedingAction handles tree-rewrite actions whose effect the main
// parser couldn't materialise (because the action's trigger token wasn't
// typed yet). If the hint carries action metadata AND the existing owner
// already has the hint's property filled, we mirror what the main parser
// would have done on the next iteration: allocate a new node of the
// action's TargetType and assign the existing owner to its action
// property slot. The new node becomes the owner the scope provider sees.
//
// When the hint has no action, or the owner's assignment slot is still
// empty (the main parser already created the post-action node), the
// owner is returned unchanged.
func applyPrecedingAction(adapter parser.LanguageCompletionAdapter, owner core.AstNode, hint *parser.CompletionHint) core.AstNode {
	if hint == nil || hint.PrecedingAction == nil {
		return owner
	}
	_, property, ok := splitHintField(hint.Field)
	if !ok {
		return owner
	}
	if !adapter.HasAssignment(owner, property) {
		return owner
	}
	action := hint.PrecedingAction
	wrapper := adapter.ApplyAction(action.TargetType, action.Property, owner)
	if wrapper == nil {
		return owner
	}
	wrapper.SetDocument(owner.Document())
	if container := owner.Container(); container != nil {
		wrapper.SetContainer(container)
	}
	return wrapper
}

// findExistingOwnerAtCursor returns the AST node of the hint's owner-rule
// type that is closest to (and contains) the cursor, by walking up from
// the owner of the token immediately preceding the cursor. Returns nil
// if no such node exists - typically because the cursor is at a position
// where the rule hasn't been entered yet, in which case the synthetic
// chain path applies.
func findExistingOwnerAtCursor(adapter parser.LanguageCompletionAdapter, doc *core.Document, ownerRule string, cursorOffset int) core.AstNode {
	template, ok := adapter.SyntheticOwnerFor(ownerRule)
	if !ok || template == nil {
		return nil
	}
	wantType := reflect.TypeOf(template)
	info := backtrackToToken(doc.Tokens, cursorOffset)
	idx := info.CurrentIdx
	if idx < 0 {
		idx = info.NextIdx - 1
	}
	if idx < 0 || idx >= len(doc.Tokens) {
		return nil
	}
	node := doc.Tokens[idx].Owner()
	for node != nil {
		if reflect.TypeOf(node) == wantType {
			return node
		}
		node = node.Container()
	}
	return nil
}

// splitHintField separates a CompletionHint.Field key of the form
// "<OwnerRule>.<Property>" into its two halves. Returns ok=false for
// malformed keys (no dot, or empty halves) so the caller can fall back
// to the unextended rule-stack chain.
func splitHintField(field string) (owner, property string, ok bool) {
	for i := 0; i < len(field); i++ {
		if field[i] == '.' {
			if i == 0 || i == len(field)-1 {
				return "", "", false
			}
			return field[:i], field[i+1:], true
		}
	}
	return "", "", false
}

// EnrichTokenCompletionItem fills zero-valued fields on item with
// per-stage defaults derived from tt and cc. The contributor returned
// item is the source of truth - any field already set on it is
// preserved verbatim.
func EnrichTokenCompletionItem(item lsp.CompletionItem, tt *core.TokenType, cc ContributorContext) lsp.CompletionItem {
	if tt == nil {
		return item
	}
	defaultLabel := tt.Label
	if defaultLabel == "" {
		defaultLabel = tt.Name
	}
	if item.Label == "" {
		item.Label = defaultLabel
	}
	if item.Kind == 0 {
		if tt.IsKeyword() {
			item.Kind = lsp.KeywordCompletion
		} else {
			item.Kind = lsp.TextCompletion
		}
	}
	if item.SortText == "" {
		item.SortText = fmt.Sprintf("%d-1-keyword", cc.SortRank)
	}
	fillInsertion(&item, cc.ReplaceRange)
	return item
}

// EnrichReferenceCompletionItem fills zero-valued fields on item with
// per-stage defaults derived from d and cc.
func EnrichReferenceCompletionItem(item lsp.CompletionItem, d *core.SymbolDescription, cc ContributorContext) lsp.CompletionItem {
	if d == nil {
		return item
	}
	defaultLabel := d.Name.String()
	if item.Label == "" {
		item.Label = defaultLabel
	}
	if item.Kind == 0 {
		item.Kind = lsp.ReferenceCompletion
	}
	if item.SortText == "" {
		item.SortText = fmt.Sprintf("%d-0-ref", cc.SortRank)
	}
	fillInsertion(&item, cc.ReplaceRange)
	return item
}

// EnrichSnippetCompletionItem fills zero-valued fields on item with
// per-stage defaults derived from sn and cc.
func EnrichSnippetCompletionItem(item lsp.CompletionItem, sn SnippetTemplate, cc ContributorContext) lsp.CompletionItem {
	if item.Label == "" {
		item.Label = sn.Label
	}
	if item.Kind == 0 {
		item.Kind = lsp.SnippetCompletion
	}
	if item.Detail == "" {
		item.Detail = sn.Detail
	}
	if item.InsertTextFormat == nil {
		fmt := lsp.SnippetTextFormat
		item.InsertTextFormat = &fmt
	}
	if item.SortText == "" {
		item.SortText = fmt.Sprintf("%d-2-snippet", cc.SortRank)
	}
	if item.TextEdit == nil && item.InsertText == "" {
		// Snippets default to inserting the body, not the label.
		applyInsertion(&item, sn.Body, cc.ReplaceRange, item.InsertTextFormat)
	} else {
		fillInsertion(&item, cc.ReplaceRange)
	}
	return item
}

// fillInsertion sets a default TextEdit / InsertText derived from Label
// when neither was set by the contributor.
func fillInsertion(item *lsp.CompletionItem, replace *lsp.Range) {
	if item.TextEdit != nil || item.InsertText != "" {
		return
	}
	applyInsertion(item, item.Label, replace, item.InsertTextFormat)
}

// applyInsertion populates either TextEdit (when replace != nil) or
// InsertText so the client knows how to apply the chosen completion. The
// LSP CompletionItem.TextEdit field is an Or<TextEdit|InsertReplaceEdit>
// wrapper; we always emit the plain TextEdit shape because it works for
// every client that supports completion.
func applyInsertion(item *lsp.CompletionItem, text string, replace *lsp.Range, insertFormat *lsp.InsertTextFormat) {
	if replace != nil {
		item.TextEdit = &lsp.Or_CompletionItem_textEdit{
			Value: lsp.TextEdit{Range: *replace, NewText: text},
		}
		// InsertText is ignored by the client when TextEdit is present,
		// but we set it anyway for clients that fall back to it.
		item.InsertText = text
		return
	}
	item.InsertText = text
	_ = insertFormat // kept for future per-callsite extension; SnippetTextFormat is already set on the item itself.
}
