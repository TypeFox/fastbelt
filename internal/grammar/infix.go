// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"fmt"

	core "typefox.dev/fastbelt"
)

// InfixOperatorGroupName returns the name of the token group synthesized by
// [ExpandInfixRules] to unify all operator tokens of the given infix rule.
func InfixOperatorGroupName(rule InfixRule) string {
	return rule.Name() + "Operator"
}

// FindInfixReturnType returns the effective return type of an infix rule:
// the explicit "returns" type if present, otherwise the return type of the
// operand rule (a lone operand is returned unchanged, so the static type must
// cover both the operand and the binary node type).
func FindInfixReturnType(rule InfixRule, ctx context.Context) Interface {
	if rule == nil {
		return nil
	}
	if typeRef := rule.ReturnType(); typeRef != nil {
		return typeRef.Ref(ctx)
	}
	if rule.Call() == nil {
		return nil
	}
	operand, _ := rule.Call().Rule().Ref(ctx).(ParserRule)
	return FindReturnType(operand, ctx)
}

// FindInfixNodeType returns the interface the binary nodes of an infix rule
// are instances of. It is always the interface named like the rule itself.
func FindInfixNodeType(rule InfixRule) Interface {
	grammar, ok := rule.Container().(Grammar)
	if !ok || grammar == nil {
		return nil
	}
	return FindInterfaceByName(grammar, rule.Name())
}

// ExpandInfixRules desugars every infix rule of the grammar in place:
//
//  1. A token group named "<RuleName>Operator" is appended to the grammar,
//     containing all operator keywords and token references of the rule. The
//     generated parser consumes any operator with a single bitset check
//     against this group, without an alternatives decision.
//
//  2. The rule's body is set to the flat form the ATN, lookahead, and
//     completion machinery operate on:
//
//     [ RuleCall(operand), [ RuleCall(operator), RuleCall(operand) ]* ]
//
// The expansion is idempotent (rules with a non-nil body are skipped) and is
// only invoked by the code generator; documents managed by the language
// server are never expanded.
func ExpandInfixRules(g Grammar) error {
	ctx := context.Background()
	for _, rule := range g.InfixRules() {
		if rule.Body() != nil {
			continue
		}
		groupName := InfixOperatorGroupName(rule)
		if findRuleByName(g, groupName) != nil {
			return fmt.Errorf("cannot generate operator token group for infix rule '%s': the name '%s' is already taken", rule.Name(), groupName)
		}
		operand := resolveOperandRule(rule, ctx)
		if operand == nil {
			return fmt.Errorf("could not resolve the operand rule of infix rule '%s'", rule.Name())
		}

		tokenGroup, err := synthesizeOperatorGroup(g, rule, groupName, ctx)
		if err != nil {
			return err
		}
		synthesizeInfixBody(rule, operand, tokenGroup)
	}
	return nil
}

func findRuleByName(g Grammar, name string) AbstractRule {
	for _, rule := range g.Rules() {
		if rule.Name() == name {
			return rule
		}
	}
	for _, rule := range g.Composites() {
		if rule.Name() == name {
			return rule
		}
	}
	for _, rule := range g.InfixRules() {
		if rule.Name() == name {
			return rule
		}
	}
	for _, token := range g.Terminals() {
		if token.Name() == name {
			return token
		}
	}
	for _, tokenGroup := range g.TokenGroups() {
		if tokenGroup.Name() == name {
			return tokenGroup
		}
	}
	return nil
}

func resolveOperandRule(rule InfixRule, ctx context.Context) ParserRule {
	call := rule.Call()
	if call == nil {
		return nil
	}
	operand, _ := call.Rule().Ref(ctx).(ParserRule)
	return operand
}

// synthesizeOperatorGroup builds the unified operator token group of an infix
// rule and appends it to the grammar.
func synthesizeOperatorGroup(g Grammar, rule InfixRule, groupName string, ctx context.Context) (TokenGroup, error) {
	tokenGroup := NewTokenGroup()
	tokenGroup.SetName(&core.Token{Image: groupName, Range: rule.NameToken().Range})
	tokenGroup.SetTextRange(rule.TextRange())
	tokenGroup.SetDocument(rule.Document())
	keywordIndex := 0
	for _, precedence := range rule.Groups() {
		for _, operator := range precedence.Operators() {
			switch op := operator.(type) {
			case Keyword:
				// A fresh node sharing the original value token: the original
				// keyword stays in place below the infix rule, and keyword
				// deduplication works on the (quoted) token value.
				keyword := NewKeyword()
				keyword.SetValue(op.ValueToken())
				keyword.SetTextRange(op.TextRange())
				keyword.SetDocument(rule.Document())
				keyword.SetContainer(tokenGroup, fieldNameKeywords, keywordIndex)
				tokenGroup.SetKeywordsItem(keyword)
				keywordIndex++
			case RuleCall:
				target, ok := op.Rule().Ref(ctx).(AbstractTokenRule)
				if !ok {
					return nil, fmt.Errorf("operator '%s' of infix rule '%s' does not reference a token rule", op.Rule().Text(), rule.Name())
				}
				tokenGroup.SetTokenRefsItem(resolvedReference(tokenGroup, op.Rule().Unit(), target))
			default:
				return nil, fmt.Errorf("unsupported operator element in infix rule '%s'", rule.Name())
			}
		}
	}
	tokenGroup.SetContainer(g, fieldNameTokenGroups, len(g.TokenGroups()))
	g.SetTokenGroupsItem(tokenGroup)
	return tokenGroup, nil
}

// synthesizeInfixBody attaches the flat "operand (operator operand)*" body to
// the infix rule.
func synthesizeInfixBody(rule InfixRule, operand ParserRule, operatorGroup TokenGroup) {
	nameUnit := rule.Call().Rule().Unit()

	operandCall := func() RuleCall {
		call := NewRuleCall()
		call.SetRule(resolvedReference[AbstractRule](call, nameUnit, operand))
		call.SetTextRange(rule.Call().TextRange())
		call.SetDocument(rule.Document())
		return call
	}

	operatorCall := NewRuleCall()
	operatorCall.SetRule(resolvedReference[AbstractRule](operatorCall, operatorGroup.NameToken(), operatorGroup))
	operatorCall.SetTextRange(rule.TextRange())
	operatorCall.SetDocument(rule.Document())

	inner := NewGroup()
	inner.SetCardinality(&core.Token{Image: "*", Range: rule.TextRange()})
	inner.SetTextRange(rule.TextRange())
	inner.SetDocument(rule.Document())
	first := operandCall()
	second := operandCall()

	outer := NewGroup()
	outer.SetTextRange(rule.TextRange())
	outer.SetDocument(rule.Document())

	first.SetContainer(outer, fieldNameElements, 0)
	outer.SetElementsItem(first)
	inner.SetContainer(outer, fieldNameElements, 1)
	outer.SetElementsItem(inner)
	operatorCall.SetContainer(inner, fieldNameElements, 0)
	inner.SetElementsItem(operatorCall)
	second.SetContainer(inner, fieldNameElements, 1)
	inner.SetElementsItem(second)

	outer.SetContainer(rule, fieldNameBody, -1)
	rule.SetBody(outer)
}

// resolvedReference creates a reference that resolves to an already known
// target node, bypassing scoping and linking.
func resolvedReference[T core.AstNode](owner core.AstNode, unit core.StringUnit, target T) *core.Reference[T] {
	description := &core.SymbolDescription{
		Node: target,
		Name: any(target).(core.NamedTokenNode).NameToken(),
	}
	if doc := target.Document(); doc != nil {
		description.URI = doc.URI
	}
	return core.NewReference(owner, unit, func(context.Context, *core.Reference[T]) (*core.SymbolDescription, *core.ReferenceError) {
		return description, nil
	})
}
