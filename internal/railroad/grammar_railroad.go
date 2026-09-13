// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import (
	"context"

	"typefox.dev/fastbelt/internal/grammar"
)

// BuildRuleDiagram builds a railroad diagram for one rule declaration.
// ok is false for rule kinds with no meaningful concrete-syntax diagram
// (terminal/token rules), so callers can skip them without erroring.
func BuildRuleDiagram(ctx context.Context, rule grammar.AbstractRule) (diagram *Diagram, ok bool) {
	switch r := rule.(type) {
	case grammar.ParserRule:
		return &Diagram{Item: MapElement(ctx, r.Body())}, true
	case grammar.CompositeRule:
		return &Diagram{Item: MapElement(ctx, r.Body())}, true
	case grammar.InfixRule:
		return &Diagram{Item: mapInfixRule(ctx, r)}, true
	default:
		// Token / TokenGroup (terminal/lexer rules): diagramming a regex is
		// a different feature, out of scope here.
		return nil, false
	}
}

// MapElement recursively maps one grammar.Element - and its cardinality -
// to a railroad.Node. Exported so tests can assert on the resulting Node
// tree directly, without round-tripping through SVG text.
func MapElement(ctx context.Context, el grammar.Element) Node {
	if el == nil {
		return Skip{}
	}
	base := mapBase(ctx, el)
	switch el.Cardinality() {
	case "?":
		return Optional(base)
	case "*":
		return ZeroOrMore(base)
	case "+":
		return &OneOrMore{Item: base}
	default:
		return base
	}
}

func mapBase(ctx context.Context, el grammar.Element) Node {
	switch e := el.(type) {
	case grammar.Alternatives:
		items := make([]Node, 0, len(e.Alts()))
		for _, alt := range e.Alts() {
			items = append(items, MapElement(ctx, alt))
		}
		if len(items) == 0 {
			// A malformed/partial Alternatives (e.g. mid-edit) - render as
			// a no-op rather than an empty Choice.
			return Skip{}
		}
		return &Choice{Normal: 0, Items: items}
	case grammar.Group:
		items := make([]Node, 0, len(e.Elements()))
		for _, item := range e.Elements() {
			items = append(items, MapElement(ctx, item))
		}
		return &Sequence{Items: items}
	case grammar.Assignment:
		// The property name and =/+=/?= operator are a tree-building
		// concern, not part of the concrete syntax the diagram depicts.
		return MapElement(ctx, e.Value())
	case grammar.CrossRef:
		// A CrossRef's own Rule() is itself a RuleCall, so this reduces to
		// the RuleCall case below.
		return MapElement(ctx, e.Rule())
	case grammar.RuleCall:
		return mapRuleCall(ctx, e)
	case grammar.Keyword:
		return &Terminal{Text: e.Value()}
	case grammar.Action:
		// Tree-rewriting action: consumes no tokens.
		return Skip{}
	default:
		return Skip{}
	}
}

func mapRuleCall(ctx context.Context, call grammar.RuleCall) Node {
	ref := call.Rule()
	name := ref.Text()
	target := ref.Ref(ctx)
	if target == nil {
		// A broken/unresolved reference (e.g. mid-edit) still renders,
		// flagged with a distinguishing style, rather than erroring.
		return &NonTerminal{Text: name, Unresolved: true}
	}
	if _, isToken := target.(grammar.AbstractTokenRule); isToken {
		return &Terminal{Text: name}
	}
	return &NonTerminal{Text: name}
}

// mapInfixRule synthesizes the same desugared shape ExpandInfixRules would
// write into an InfixRule's Body - "operand (operator operand)*".
// Precedence and associativity are intentionally flattened.
func mapInfixRule(ctx context.Context, rule grammar.InfixRule) Node {
	operand := MapElement(ctx, rule.Call())

	var operators []Node
	for _, group := range rule.Groups() {
		for _, op := range group.Operators() {
			operators = append(operators, MapElement(ctx, op))
		}
	}
	if len(operators) == 0 {
		return operand
	}

	var operatorNode Node
	if len(operators) == 1 {
		operatorNode = operators[0]
	} else {
		operatorNode = &Choice{Normal: 0, Items: operators}
	}

	secondOperand := MapElement(ctx, rule.Call())
	repeat := ZeroOrMore(&Sequence{Items: []Node{operatorNode, secondOperand}})
	return &Sequence{Items: []Node{operand, repeat}}
}
