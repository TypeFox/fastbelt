package services

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar/generated"
)

type FastbeltScopeProvider struct {
	*generated.DefaultFastbeltScopeProvider
}

func NewFastbeltScopeProvider(srv generated.FastbeltLinkingSrvCont) *FastbeltScopeProvider {
	return &FastbeltScopeProvider{
		DefaultFastbeltScopeProvider: generated.NewDefaultFastbeltScopeProvider(srv),
	}
}

func (s *FastbeltScopeProvider) ScopeActionProperty(ctx context.Context, reference *core.Reference[generated.Field]) core.Scope {
	if action, ok := reference.Owner.(generated.Action); ok && action.Type() != nil {
		targetType := action.Type().Ref(ctx)
		if targetType != nil {
			fieldDesc := []*core.AstNodeDescription{}
			for _, field := range targetType.Fields() {
				if field.Name() != "" {
					fieldDesc = append(fieldDesc, core.NewAstNodeDescription(field, field.Name(), &field.NameToken().Segment, field.Segment()))
				}
			}
			return core.NewMapScopeFromSlice(fieldDesc, nil)
		}
	}
	return core.EmptyScope
}
