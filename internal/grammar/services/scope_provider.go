// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

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
			descriptions := generateInterfaceFieldsDescriptions(ctx, targetType, map[generated.Interface]bool{})
			return core.NewMapScopeFromSlice(descriptions, nil)
		}
	}
	return core.EmptyScope
}

func (s *FastbeltScopeProvider) ScopeAssignmentProperty(ctx context.Context, reference *core.Reference[generated.Field]) core.Scope {
	if assignment, ok := reference.Owner.(generated.Assignment); ok {
		iface := getCurrentType(ctx, assignment)
		if iface == nil {
			return nil
		}
		descriptions := generateInterfaceFieldsDescriptions(ctx, iface, map[generated.Interface]bool{})
		return core.NewMapScopeFromSlice(descriptions, nil)
	}
	return nil
}

func getCurrentType(ctx context.Context, node core.AstNode) generated.Interface {
	for node != nil {
		if rule, ok := node.(generated.ParserRule); ok {
			// Arrived at the parser rule, return its return type
			return rule.ReturnType().Ref(ctx)
		}
		container := node.Container()
		if group, ok := container.(generated.Group); ok {
			// Attempt to find the last action that was executed in the parser rule
			elem := group.Elements()
			var lastAction generated.Action = nil
			for i := range group.Elements() {
				if action, ok := elem[i].(generated.Action); ok {
					lastAction = action
				}
				if elem[i] == node {
					break
				}
			}
			if lastAction != nil {
				return lastAction.Type().Ref(ctx)
			}
		}
		node = container
	}
	return nil
}

func generateInterfaceFieldsDescriptions(ctx context.Context, iface generated.Interface, visited map[generated.Interface]bool) []*core.AstNodeDescription {
	fieldDesc := []*core.AstNodeDescription{}
	if visited[iface] {
		return fieldDesc
	}
	visited[iface] = true
	for _, field := range iface.Fields() {
		if field.Name() != "" {
			fieldDesc = append(fieldDesc, core.NewAstNodeDescription(field, field.Name(), &field.NameToken().Segment, field.Segment()))
		}
	}
	for _, super := range iface.Extends() {
		superType := super.Ref(ctx)
		if superType != nil {
			fieldDesc = append(fieldDesc, generateInterfaceFieldsDescriptions(ctx, superType, visited)...)
		}
	}
	return fieldDesc
}
