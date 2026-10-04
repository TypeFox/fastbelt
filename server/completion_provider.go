// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unique"

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
// container if any, otherwise the default contributor.
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
	// The rule calls of a hint are relative to the token in front of the cursor
	var lastToken *core.Token
	if cc.PrefixLen > 0 {
		lastToken = &doc.Tokens[cc.PrefixLen-1]
	}

	contribCtx := ContributorContext{
		Doc:          doc,
		Cursor:       cursorOffset,
		ReplaceRange: cc.ReplaceRange,
		SortRank:     cc.SortRank,
	}

	items := make([]lsp.CompletionItem, 0, len(info.Tokens)+len(info.Hints)+4)

	// Token pass: contributor decides per (TokenType, atnState) what to emit.
	for _, tc := range info.Tokens {
		if !matcher.Match(cc.ReplaceText, tokenLabel(tc.TokenType)) {
			continue
		}
		tcCopy := tc
		accept := func(item lsp.CompletionItem) {
			items = append(items, EnrichTokenCompletionItem(item, tcCopy.TokenType, contribCtx))
		}
		contributor.CompletionForToken(ctx, tcCopy.TokenType, tcCopy.ATNStateIdx, contribCtx, accept)
	}

	// Cross-reference pass: dispatch per hint on its owner; contributor
	// decides per (SymbolDescription, hint, atnState) what to emit.
	for _, hc := range info.Hints {
		owner := buildOwner(adapter, doc, result.RuleStack, hc, lastToken)
		if owner == nil {
			continue
		}
		seq, ok := adapter.DispatchCompletion(ctx, hc.Hint.Key(), owner)
		if !ok {
			continue
		}
		hcCopy := hc
		refCtx := contribCtx
		refCtx.Node = owner
		for d := range seq {
			dCopy := d
			if !matcher.Match(cc.ReplaceText, dCopy.Name) {
				continue
			}
			accept := func(item lsp.CompletionItem) {
				items = append(items, EnrichReferenceCompletionItem(item, dCopy, refCtx))
			}
			contributor.CompletionForReference(ctx, dCopy, hcCopy.Hint, hcCopy.ATNStateIdx, refCtx, accept)
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
// Token boundaries: Range.Start is inclusive, End is
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
		if int(tokens[mid].Range.Start) < offset {
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
		end := int(tokens[prev].Range.End)
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
	replace := curToken.Range.LspRange(doc.TextDoc)
	prefixLen := info.CurrentIdx
	typed := curToken.Image
	if !info.CurrentAtEnd {
		typed = curToken.Image[:offset-int(curToken.Range.Start)]
	}

	// When the current token is part of a multi-token CompositeNode,
	// the REPLACE range must cover the whole composite span - otherwise accepting
	// a candidate would overwrite only the trailing token and leave earlier segments
	// behind, producing duplicated text. Walk back through tokens that
	// belong to the same composite and widen the replace range / prefix.
	if startIdx, ok := compositeStart(doc.Tokens, info.CurrentIdx); ok {
		startToken := &doc.Tokens[startIdx]
		// Replace start position with the start of the first token in the composite
		replace.Start = startToken.Range.LspRange(doc.TextDoc).Start
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
			node.SetContainer(parent, unique.Make(""), 0)
		}
		node.SetDocument(doc)
		parent = node
	}
	return parent
}

// buildOwner returns the node that owns the cross-reference of the hint.
//
// The owner is derived from the node that owns the lastToken, which is the
// token in front of the cursor, and from the rule calls that lead from there
// to the cross-reference. The result is a node of the AST that the main
// parser built, or a new node with the container and the containment data
// that the main parser would give to it. Scope providers find the same
// symbols as for a node of the main parser then.
//
// The lastToken is nil if there is no such token. The owner is derived from
// the rule stack of the completion parser then, and is not connected to the
// AST.
//
// If the hint carries an action that precedes the cross-reference, the
// action is applied to the owner as the main parser would have on the next
// token, see applyPrecedingAction.
//
// Returns nil if the adapter doesn't know one of the rule keys; the
// completion request then yields no candidates for this hint rather than
// silently returning the wrong scope.
func buildOwner(adapter parser.LanguageCompletionAdapter, doc *core.Document, ruleStack []parser.RuleContext, hc parser.HintCompletion, lastToken *core.Token) core.AstNode {
	var owner core.AstNode
	if lastToken != nil {
		if base := baseOf(adapter, hc, lastToken); base != nil {
			owner = buildOwnerAt(adapter, base, hc, lastToken)
		}
	}
	if owner == nil {
		owner = buildSyntheticOwnerChainFor(adapter, doc, ruleStack, hc.Hint)
	}
	if owner == nil {
		return nil
	}
	return applyPrecedingAction(adapter, owner, hc.Hint)
}

// baseOf returns the node of the rule that was in progress at the last
// token, after leaving the rules that ended since then. Returns nil if that
// node is unknown.
func baseOf(adapter parser.LanguageCompletionAdapter, hc parser.HintCompletion, lastToken *core.Token) core.AstNode {
	if hc.ConsumedAt >= 0 && hc.ConsumedAt != lastToken.Kind {
		// The hint reads the token in another way than the main parser did,
		// so the node that owns the token is unrelated to the hint.
		return nil
	}
	base := lastToken.Owner()
	for _, call := range hc.Left {
		if base == nil {
			return nil
		}
		base = containerOfCall(adapter, base, call)
	}
	return base
}

// containerOfCall returns the container of the node that the assigned rule
// call created, which is the given node or one of its containers.
//
// A rule call creates several nodes if the rule contains tree-rewriting
// actions or is an infix rule, and the outermost of them is the node of the
// rule call. The main parser also creates them for the input that follows
// the cursor. Returns nil if that node is not assigned like the rule call.
func containerOfCall(adapter parser.LanguageCompletionAdapter, node core.AstNode, call *parser.RuleCallInfo) core.AstNode {
	for {
		container := node.Container()
		if container == nil {
			return nil
		}
		field, index := node.ContainmentData()
		if !adapter.AssignsCurrent(container, field, index) {
			if field == unique.Make(call.Property) {
				return container
			}
			return nil
		}
		node = container
	}
}

// buildOwnerAt returns the owner of a cross-reference that is reached from
// the base through the rule calls of the hint. The base is the node of the
// rule that the first rule call belongs to.
//
// An unassigned rule call continues the node of the calling rule, and an
// assigned rule call creates a node that is contained in it. The owner is
// the node of the last rule call. The lastToken is the token in front of it.
//
// Returns nil if the type of a node is unknown.
func buildOwnerAt(adapter parser.LanguageCompletionAdapter, base core.AstNode, hc parser.HintCompletion, lastToken *core.Token) core.AstNode {
	ownerType := hc.Hint.Owner
	doc := base.Document()
	// The rule call that ended last
	var previous *parser.RuleCallInfo
	if len(hc.Left) > 0 {
		previous = hc.Left[len(hc.Left)-1]
	}
	// The node of the rule in progress. If assigned is set, it is the
	// container of the node of that rule call, which is not created yet.
	node := base
	var assigned *parser.RuleCallInfo
	nodeType := ""
	create := func(nodeType string) core.AstNode {
		created := newNode(adapter, nodeType, doc)
		if created == nil {
			return nil
		}
		field := unique.Make(assigned.Property)
		index := -1
		if assigned.List {
			index = listIndexAfter(node, field, lastToken)
		}
		created.SetContainer(node, field, index)
		return created
	}
	for _, call := range hc.Calls {
		if call == nil {
			return nil
		}
		if action := call.PrecedingAction; action != nil {
			if assigned != nil {
				// The action determines the type of the node to create
				nodeType = action.TargetType
			} else if call.Repeated && call == previous {
				// The action was executed in front of the first of the repeated rule calls
			} else if node = applyAction(adapter, node, action); node == nil {
				return nil
			}
		}
		if call.Property != "" {
			if assigned != nil {
				if node = create(nodeType); node == nil {
					return nil
				}
			}
			assigned, nodeType = call, call.Type
		} else if assigned != nil && call.Type != "" {
			// The called rule determines the type of the node to create,
			// e.g. as one of several alternatives.
			nodeType = call.Type
		}
	}
	if assigned != nil {
		return create(ownerType)
	}
	// No rule call is assigned, so the owner is the node of the rule in progress
	if isOfType(adapter, node, ownerType) {
		return node
	}
	return replaceNode(adapter, node, ownerType)
}

// applyAction returns the node that results from executing the action with
// the given node as the current node.
func applyAction(adapter parser.LanguageCompletionAdapter, node core.AstNode, action *parser.ActionInfo) core.AstNode {
	if action.Property == "" {
		if isOfType(adapter, node, action.TargetType) {
			return node
		}
		return replaceNode(adapter, node, action.TargetType)
	}
	// The main parser executed the action if the input continues behind the cursor
	if container := node.Container(); container != nil && isOfType(adapter, container, action.TargetType) {
		if field, index := node.ContainmentData(); field == unique.Make(action.Field) && index <= 0 {
			return container
		}
	}
	wrapper := adapter.ApplyAction(action.TargetType, action.Property, node)
	if wrapper == nil {
		return nil
	}
	takePlace(wrapper, node)
	return wrapper
}

// newNode creates a node of the given type. Returns nil if the type is unknown.
func newNode(adapter parser.LanguageCompletionAdapter, nodeType string, doc *core.Document) core.AstNode {
	node, ok := adapter.SyntheticOwnerFor(nodeType)
	if !ok || node == nil {
		return nil
	}
	node.SetDocument(doc)
	return node
}

// replaceNode creates a node of the given type that takes the place of the
// given node in the AST. Returns nil if the type is unknown.
func replaceNode(adapter parser.LanguageCompletionAdapter, node core.AstNode, nodeType string) core.AstNode {
	replacement := newNode(adapter, nodeType, node.Document())
	if replacement != nil {
		takePlace(replacement, node)
	}
	return replacement
}

// takePlace gives the node the container and the containment data of another
// node. The container does not refer to the node.
func takePlace(node, other core.AstNode) {
	field, index := other.ContainmentData()
	node.SetContainer(other.Container(), field, index)
	node.SetDocument(other.Document())
}

// isOfType reports whether the node has the given type, and not a subtype of it.
func isOfType(adapter parser.LanguageCompletionAdapter, node core.AstNode, nodeType string) bool {
	template, ok := adapter.SyntheticOwnerFor(nodeType)
	return ok && template != nil && reflect.TypeOf(node) == reflect.TypeOf(template)
}

// listIndexAfter returns the index for a node that is added to the list of
// the container after the given token, in front of the items that follow it.
// An item that starts in front of the token and ends after it is split by
// the new node, which follows the first part.
func listIndexAfter(container core.AstNode, field unique.Handle[string], token *core.Token) int {
	index := 0
	container.ForEachNode(func(child core.AstNode, childField unique.Handle[string], childIndex int) {
		if childField == field && !isEmptyNode(child) && child.TextRange().Start <= token.Range.Start {
			index = childIndex + 1
		}
	})
	return index
}

// isEmptyNode reports whether the node has neither tokens, child nodes nor
// cross-reference text, like the nodes that the parser creates for input
// that is missing. The tokens of a composite cross-reference belong to its
// unit, not to the node.
func isEmptyNode(node core.AstNode) bool {
	if len(node.Tokens()) > 0 {
		return false
	}
	for range core.ChildNodes(node) {
		return false
	}
	empty := true
	node.ForEachReference(func(ref core.UntypedReference, _ unique.Handle[string], _ int) {
		if unit := ref.Unit(); unit != nil && unit.String() != "" {
			empty = false
		}
	})
	return empty
}

// buildSyntheticOwnerChainFor extends buildSyntheticOwnerChain with the
// hint's owner rule when the rule stack doesn't already end at that rule.
//
// At cursor positions where a new rule could begin but hasn't yet, the
// parser's rule stack stops one level above the rule the hint's owner
// belongs to, so a synthetic frame for the owner is appended if the stack
// doesn't already end there.
//
// Returns nil if the adapter doesn't know one of the rule keys; the
// completion request then yields no candidates for this hint rather than
// silently returning the wrong scope.
func buildSyntheticOwnerChainFor(adapter parser.LanguageCompletionAdapter, doc *core.Document, ruleStack []parser.RuleContext, hint *parser.CompletionHint) core.AstNode {
	ownerType := hint.Owner
	// If the owner rule appears anywhere on the stack, the parser is
	// already inside it - slice off any deeper frames. This covers
	// cross-references whose text form is a separate rule: at the cursor
	// the rule stack might end inside a string/composite rule that
	// carries no AST node, while the owner rule (which holds the
	// reference) sits one frame higher. Searching from the top down
	// picks the innermost matching frame.
	for i := len(ruleStack) - 1; i >= 0; i-- {
		if ruleStack[i].RuleKey == ownerType {
			return buildSyntheticOwnerChain(adapter, doc, ruleStack[:i+1])
		}
	}
	// The owner rule isn't on the stack - the parser hasn't entered it
	// yet (cursor sits at a position where it could begin). Append a
	// synthetic frame for the owner so the chain has somewhere to
	// dispatch on.
	extended := make([]parser.RuleContext, 0, len(ruleStack)+1)
	extended = append(extended, ruleStack...)
	extended = append(extended, parser.RuleContext{RuleKey: ownerType})
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
	if hint == nil || hint.PrecedingAction == nil || !adapter.HasAssignment(owner, hint.Property) {
		return owner
	}
	action := hint.PrecedingAction
	wrapper := adapter.ApplyAction(action.TargetType, action.Property, owner)
	if wrapper == nil {
		return owner
	}
	takePlace(wrapper, owner)
	return wrapper
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
	defaultLabel := d.Name
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
		item.TextEdit = &lsp.CompletionItemTextEdit{
			TextEdit: &lsp.TextEdit{Range: *replace, NewText: text},
		}
		// InsertText is ignored by the client when TextEdit is present,
		// but we set it anyway for clients that fall back to it.
		item.InsertText = text
		return
	}
	item.InsertText = text
	_ = insertFormat // kept for future per-callsite extension; SnippetTextFormat is already set on the item itself.
}
