// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package completion

import (
	"context"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

type CompletionScopeProviderImpl struct {
	DefaultCompletionScopeProvider
}

func NewCompletionScopeProviderImpl(sc *service.Container) *CompletionScopeProviderImpl {
	return &CompletionScopeProviderImpl{DefaultCompletionScopeProvider{sc: sc}}
}

func (s *CompletionScopeProviderImpl) ScopeMemberCallRef(ctx context.Context, reference *fastbelt.Reference[Declare]) fastbelt.Scope {
	owner := reference.Owner()
	if memberCall, ok := owner.(MemberCall); ok {
		previous := memberCall.Previous()
		if previous == nil {
			// Use local scope if no previous member is set
			return s.DefaultCompletionScopeProvider.ScopeMemberCallRef(ctx, reference)
		}
		// Otherwise, use the scope of the previous member
		return childrenScope(previous.Ref().Ref(ctx))
	} else {
		return fastbelt.EmptyScope
	}
}

// ScopeLeadRef returns the Leads of the root that contains the owner. The
// scope depends on the container of the owner, which the completion provider
// has to supply at the start of the document, where the cross-reference of
// Lead has no token in front of it.
func (s *CompletionScopeProviderImpl) ScopeLeadRef(ctx context.Context, reference *fastbelt.Reference[Lead]) fastbelt.Scope {
	root, ok := reference.Owner().Container().(Root)
	if !ok {
		return fastbelt.EmptyScope
	}
	symbols := []*fastbelt.SymbolDescription{}
	for _, obj := range root.Objects() {
		if lead, ok := obj.(Lead); ok && lead.NameToken() != nil {
			symbols = append(symbols, fastbelt.NewSymbolDescription(lead, lead.NameToken()))
		}
	}
	return fastbelt.NewMapScopeFromSlice(symbols, nil)
}

// ScopeDepB returns the children of the Declare that A refers to.
func (s *CompletionScopeProviderImpl) ScopeDepB(ctx context.Context, reference *fastbelt.Reference[Declare]) fastbelt.Scope {
	dep, ok := reference.Owner().(Dep)
	if !ok || dep.A() == nil {
		return fastbelt.EmptyScope
	}
	return childrenScope(dep.A().Ref(ctx))
}

// childrenScope returns the children of the Declare, or the empty scope if
// the Declare is nil.
func childrenScope(decl Declare) fastbelt.Scope {
	if decl == nil {
		return fastbelt.EmptyScope
	}
	symbols := []*fastbelt.SymbolDescription{}
	for _, member := range decl.Children() {
		symbols = append(symbols, fastbelt.NewSymbolDescription(member, member.NameNode()))
	}
	return fastbelt.NewMapScopeFromSlice(symbols, nil)
}
