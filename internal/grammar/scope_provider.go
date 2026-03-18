// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"

	core "typefox.dev/fastbelt"
)

type scopeProviderImpl struct {
	*DefaultFastbeltScopeProvider
}

func newScopeProviderImpl(srv FastbeltLinkingSrvCont) *scopeProviderImpl {
	return &scopeProviderImpl{
		DefaultFastbeltScopeProvider: NewDefaultFastbeltScopeProvider(srv),
	}
}

func (s *scopeProviderImpl) ScopeActionProperty(ctx context.Context, reference *core.Reference[Field]) core.Scope {
	if action, ok := reference.Owner().(Action); ok && action.Type() != nil {
		targetType := action.Type().Ref(ctx)
		if targetType != nil {
			descriptions := generateInterfaceFieldsDescriptions(ctx, targetType, map[Interface]bool{})
			return core.NewMapScopeFromSlice(descriptions, nil)
		}
	}
	return core.EmptyScope
}

func (s *scopeProviderImpl) ScopeAssignmentProperty(ctx context.Context, reference *core.Reference[Field]) core.Scope {
	if assignment, ok := reference.Owner().(Assignment); ok {
		iface := getCurrentType(ctx, assignment)
		if iface == nil {
			return nil
		}
		descriptions := generateInterfaceFieldsDescriptions(ctx, iface, map[Interface]bool{})
		return core.NewMapScopeFromSlice(descriptions, nil)
	}
	return nil
}

func getCurrentType(ctx context.Context, node core.AstNode) Interface {
	for node != nil {
		if rule, ok := node.(ParserRule); ok {
			// Arrived at the parser rule, return its return type
			return rule.ReturnType().Ref(ctx)
		}
		container := node.Container()
		if group, ok := container.(Group); ok {
			// Attempt to find the last action that was executed in the parser rule
			elem := group.Elements()
			var lastAction Action = nil
			for i := range group.Elements() {
				if action, ok := elem[i].(Action); ok {
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

func generateInterfaceFieldsDescriptions(ctx context.Context, iface Interface, visited map[Interface]bool) []*core.AstNodeDescription {
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
