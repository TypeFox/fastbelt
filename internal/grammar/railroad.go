// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"

	"typefox.dev/fastbelt/internal/railroad"
)

// BuildRuleDiagram builds a railroad diagram for one rule declaration.
// ok is false for rule kinds with no meaningful concrete-syntax diagram
// (terminal/token rules), so callers can skip them without erroring.
func BuildRuleDiagram(ctx context.Context, rule AbstractRule) (diagram *railroad.Diagram, ok bool) {
	switch r := rule.(type) {
	case ParserRule:
		return &railroad.Diagram{Item: MapElement(ctx, r.Body())}, true
	case CompositeRule:
		return &railroad.Diagram{Item: MapElement(ctx, r.Body())}, true
	case InfixRule:
		return &railroad.Diagram{Item: mapInfixRule(ctx, r)}, true
	default:
		// Token / TokenGroup (terminal/lexer rules): diagramming a regex is
		// a different feature, out of scope here.
		return nil, false
	}
}

// MapElement recursively maps one Element - and its cardinality - to a
// railroad.Node. Exported so tests can assert on the resulting Node tree
// directly, without round-tripping through SVG text.
func MapElement(ctx context.Context, el Element) railroad.Node {
	if el == nil {
		return railroad.Skip{}
	}
	base := mapBase(ctx, el)
	switch el.Cardinality() {
	case "?":
		return railroad.Optional(base)
	case "*":
		return railroad.ZeroOrMore(base)
	case "+":
		return &railroad.OneOrMore{Item: base}
	default:
		return base
	}
}

func mapBase(ctx context.Context, el Element) railroad.Node {
	switch e := el.(type) {
	case Alternatives:
		items := make([]railroad.Node, 0, len(e.Alts()))
		for _, alt := range e.Alts() {
			items = append(items, MapElement(ctx, alt))
		}
		if len(items) == 0 {
			// A malformed/partial Alternatives (e.g. mid-edit) - render as
			// a no-op rather than an empty Choice.
			return railroad.Skip{}
		}
		return &railroad.Choice{Normal: 0, Items: items}
	case Group:
		items := make([]railroad.Node, 0, len(e.Elements()))
		for _, item := range e.Elements() {
			items = append(items, MapElement(ctx, item))
		}
		return &railroad.Sequence{Items: items}
	case Assignment:
		// The property name and =/+=/?= operator are a tree-building
		// concern, not part of the concrete syntax the diagram depicts.
		return MapElement(ctx, e.Value())
	case CrossRef:
		// A CrossRef's own Rule() is itself a RuleCall, so this reduces to
		// the RuleCall case below.
		return MapElement(ctx, e.Rule())
	case RuleCall:
		return mapRuleCall(ctx, e)
	case Keyword:
		return &railroad.Terminal{Text: e.Value()}
	case Action:
		// Tree-rewriting action: consumes no tokens.
		return railroad.Skip{}
	default:
		return railroad.Skip{}
	}
}

func mapRuleCall(ctx context.Context, call RuleCall) railroad.Node {
	ref := call.Rule()
	name := ref.Text()
	target := ref.Ref(ctx)
	if target == nil {
		// A broken/unresolved reference (e.g. mid-edit) still renders,
		// flagged with a distinguishing style, rather than erroring.
		return &railroad.NonTerminal{Text: name, Unresolved: true}
	}
	if _, isToken := target.(AbstractTokenRule); isToken {
		return &railroad.Terminal{Text: name}
	}
	return &railroad.NonTerminal{Text: name}
}

// mapInfixRule synthesizes the same desugared shape ExpandInfixRules would
// write into an InfixRule's Body - "operand (operator operand)*".
// Precedence and associativity are intentionally flattened.
func mapInfixRule(ctx context.Context, rule InfixRule) railroad.Node {
	operand := MapElement(ctx, rule.Call())

	var operators []railroad.Node
	for _, group := range rule.Groups() {
		for _, op := range group.Operators() {
			operators = append(operators, MapElement(ctx, op))
		}
	}
	if len(operators) == 0 {
		return operand
	}

	var operatorNode railroad.Node
	if len(operators) == 1 {
		operatorNode = operators[0]
	} else {
		operatorNode = &railroad.Choice{Normal: 0, Items: operators}
	}

	secondOperand := MapElement(ctx, rule.Call())
	repeat := railroad.ZeroOrMore(&railroad.Sequence{Items: []railroad.Node{operatorNode, secondOperand}})
	return &railroad.Sequence{Items: []railroad.Node{operand, repeat}}
}
