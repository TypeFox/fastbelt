// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/service"
)

type scopeProviderImpl struct {
	DefaultFastbeltScopeProvider
}

func newScopeProviderImpl(_ *service.Container) *scopeProviderImpl {
	return &scopeProviderImpl{}
}

// Token modes need no override: every .fb file of a folder is part of one
// grammar with a single mode table, so the generated default scope (all named
// modes of the folder) is the right one. The unnamed default mode is targeted
// as `push(default)` and is never exported.

func (s *scopeProviderImpl) ScopeRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) core.Scope {
	root, ok := reference.Owner().Document().Root.(Grammar)
	if !ok {
		return core.EmptyScope
	}
	return extractTokenReference(root, reference)
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
			return core.EmptyScope
		}
		descriptions := generateInterfaceFieldsDescriptions(ctx, iface, map[Interface]bool{})
		return core.NewMapScopeFromSlice(descriptions, nil)
	}
	return core.EmptyScope
}

func (s *scopeProviderImpl) ScopeTokenGroupTokenRefs(ctx context.Context, reference *core.Reference[AbstractTokenRule]) core.Scope {
	root, ok := reference.Owner().Document().Root.(Grammar)
	if !ok {
		return core.EmptyScope
	}
	return extractTokenReference(root, reference)
}

// extractTokenReference scopes a token reference to the tokens and groups
// declared inside the token modes of every grammar file in the folder (those are
// nested too deep to be exported), followed by the default scope.
func extractTokenReference[T AbstractRule](root Grammar, reference *core.Reference[T]) core.Scope {
	symbols := []*core.SymbolDescription{}
	for _, g := range siblingGrammars(root) {
		for _, tokenMode := range g.TokenModes() {
			for _, member := range tokenMode.Members() {
				if tokenDeclUsage, ok := member.(TokenDeclUsage); ok {
					tokenDecl := tokenDeclUsage.Declaration()
					if tokenDecl.NameToken() == nil {
						continue
					}
					symbols = append(symbols, core.NewSymbolDescription(tokenDecl, tokenDecl.NameToken()))
				} else if tokenGroupUsage, ok := member.(TokenGroupUsage); ok {
					tokenGroup := tokenGroupUsage.Group()
					if tokenGroup.NameToken() == nil {
						continue
					}
					symbols = append(symbols, core.NewSymbolDescription(tokenGroup, tokenGroup.NameToken()))
				}
			}
		}
	}
	outer := linking.DefaultScopeOfType[T](reference.Owner())
	return core.NewMapScopeFromSlice(symbols, outer)
}

func (s *scopeProviderImpl) ScopeTokenUsageTokenRef(ctx context.Context, reference *core.Reference[AbstractTokenRule]) core.Scope {
	root, ok := reference.Owner().Document().Root.(Grammar)
	if !ok {
		return core.EmptyScope
	}
	return extractTokenReference(root, reference)
}

func getCurrentType(ctx context.Context, node core.AstNode) Interface {
	for node != nil {
		if rule, ok := node.(ParserRule); ok {
			// Arrived at the parser rule, return its return type
			return FindReturnType(rule, ctx)
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

func generateInterfaceFieldsDescriptions(ctx context.Context, iface Interface, visited map[Interface]bool) []*core.SymbolDescription {
	fieldDesc := []*core.SymbolDescription{}
	if visited[iface] {
		return fieldDesc
	}
	visited[iface] = true
	for _, field := range iface.Fields() {
		if field.Name() != "" {
			fieldDesc = append(fieldDesc, core.NewTokenSymbolDescription(field))
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
