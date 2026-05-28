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
		decl := previous.Ref().Ref(ctx)
		if decl == nil {
			return fastbelt.EmptyScope
		}
		symbols := []*fastbelt.SymbolDescription{}
		for _, member := range decl.Children() {
			symbols = append(symbols, fastbelt.NewSymbolDescription(member, member.NameNode()))
		}
		return fastbelt.NewMapScopeFromSlice(symbols, nil)
	} else {
		return fastbelt.EmptyScope
	}
}
