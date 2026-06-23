// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package atn

import (
	"context"
	"fmt"
	"reflect"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
)

// completionHintFor returns the per-field CompletionHint for a CrossRef whose
// container is an Assignment in the given rule. Returns nil if the CrossRef is
// not nested inside an Assignment (bare cross-references contribute no hint).
// completionHintFor returns the per-field CompletionHint for a CrossRef whose
// container is an Assignment in the given rule. Returns nil if the CrossRef is
// not nested inside an Assignment (bare cross-references contribute no hint).
func completionHintFor(rule grammar.AbstractRuleWithBody, cr grammar.CrossRef) *parser.CompletionHint {
	if rule == nil {
		return nil
	}
	assignment, ok := cr.Container().(grammar.Assignment)
	if !ok {
		return nil
	}
	prop := assignment.Property()
	if prop == nil {
		return nil
	}
	field := prop.Ref(context.Background())
	if field == nil {
		return nil
	}
	fieldName := field.Name()
	if fieldName == "" {
		return nil
	}
	iface := field.Container().(grammar.Interface)
	hint := &parser.CompletionHint{Field: iface.Name() + "." + fieldName}
	if action := findPrecedingAction(assignment); action != nil {
		typeName := ""
		if t := action.Type(); t != nil {
			typeName = t.Text()
		}
		actionProp := ""
		if p := action.Property(); p != nil {
			actionProp = p.Text()
		}
		hint.PrecedingAction = &parser.ActionInfo{
			TargetType: typeName,
			Property:   actionProp,
		}
	}
	return hint
}

// findPrecedingAction returns the grammar.Action that fires immediately
// before el's first token is consumed, or nil if no such action exists.
//
// The walk handles indirection through Alternatives and outer Groups: a CR
// inside `(a | b)` whose enclosing alternative has no token-consuming
// predecessor is still "first", so an action in the surrounding group
// applies to it. The walk stops at the rule body or as soon as it crosses
// any token-consuming element.
func findPrecedingAction(el grammar.Element) grammar.Action {
	node := fastbelt.AstNode(el)
	for node != nil {
		parent := node.Container()
		if parent == nil {
			return nil
		}
		switch p := parent.(type) {
		case grammar.Group:
			elems := p.Elements()
			idx := -1
			for i, e := range elems {
				if fastbelt.AstNode(e) == node {
					idx = i
					break
				}
			}
			if idx < 0 {
				return nil
			}
			for i := idx - 1; i >= 0; i-- {
				element := elems[i]
				if action, ok := element.(grammar.Action); ok {
					return action
				} else if element.Cardinality() == "" || element.Cardinality() == "+" {
					return nil // unskippable prior sibling consumes tokens, so no preceding action applies
				}
			}
			node = parent
		case grammar.Alternatives:
			node = parent
		case grammar.AbstractRuleWithBody:
			return nil
		default:
			node = parent
		}
	}
	return nil
}

func CreateATN(grammr grammar.Grammar, tokenTypeIds map[string]int) (*ATN, map[string]grammar.AbstractRuleWithBody) {
	lookaheadNames := ComputeLookaheadNames(grammr)
	byName := map[string]grammar.AbstractRuleWithBody{}
	for _, gr := range grammr.Rules() {
		byName[gr.Name()] = gr
	}
	builder := NewATNBuilder(lookaheadNames, tokenTypeIds, byName)

	entries := map[grammar.AbstractRuleWithBody]ATNRuleBuilder{}
	allRules := []grammar.AbstractRuleWithBody{}
	for _, gr := range grammr.Rules() {
		allRules = append(allRules, gr)
	}
	for _, gr := range grammr.Composites() {
		allRules = append(allRules, gr)
	}
	for _, gr := range allRules {
		entries[gr] = builder.DeclareRule(gr)
	}
	for _, gr := range allRules {
		ruleBuilder := entries[gr]
		handle, err := convertElement(ruleBuilder, gr.Body())
		if err != nil {
			// Handle the error appropriately, e.g., log it or return it
			fmt.Printf("Error converting element for rule %s: %v\n", gr.Name(), err)
		}
		if handle == nil {
			continue
		}
		ruleBuilder.Assign(handle)
	}
	atn := builder.Build()

	// Tag decision states with their grammar element so the generator can
	// map grammar.Element to ATN state index.
	reverseNames := make(map[string]grammar.Element, len(lookaheadNames))
	for el, name := range lookaheadNames {
		reverseNames[name] = el
	}
	for name, state := range atn.DecisionMap {
		if el, ok := reverseNames[name]; ok {
			state.Production = el
		}
	}

	return atn, byName
}

func ComputeLookaheadNames(grammr grammar.Grammar) map[grammar.Element]string {
	counters := make(map[reflect.Type]int)
	names := map[grammar.Element]string{}
	ruleName := ""
	for node := range fastbelt.AllNodes(grammr) {
		switch e := node.(type) {
		case grammar.ParserRule:
			counters = make(map[reflect.Type]int)
			ruleName = e.Name()
		case grammar.Keyword, grammar.RuleCall, grammar.CrossRef, grammar.Alternatives, grammar.Group:
			el := e.(grammar.Element)
			index := nextCounter(counters, el)
			var typeName string
			switch el.(type) {
			case grammar.Keyword:
				typeName = "Keyword"
			case grammar.RuleCall:
				typeName = "RuleCall"
			case grammar.CrossRef:
				typeName = "CrossRef"
			case grammar.Alternatives:
				typeName = "Alternatives"
			case grammar.Group:
				typeName = "Group"
			}
			name := fmt.Sprintf("%s_%s_%d", ruleName, typeName, index)
			names[el] = name
		}
	}
	return names
}

func convertElement(
	rb ATNRuleBuilder,
	el grammar.Element,
) (*ATNHandle, error) {
	switch e := el.(type) {
	case grammar.Alternatives:
		return convertAlternatives(rb, e, e.Cardinality())

	case grammar.Group:
		return convertGroup(rb, e)

	case grammar.Assignment:
		// Assignment is transparent: recurse into Value().
		if e.Value() == nil {
			return nil, nil
		}
		return convertAssignable(rb, e.Value(), e.Cardinality())

	case grammar.CrossRef:
		return convertCrossRef(rb, e, e.Cardinality())

	case grammar.RuleCall:
		return convertRuleCall(rb, e, e.Cardinality())

	case grammar.Keyword:
		return convertKeyword(rb, e, e.Cardinality())

	case grammar.Action:
		// Semantic actions have no ATN impact.
		return nil, nil

	default:
		return nil, nil
	}
}

func convertAssignable(
	rb ATNRuleBuilder,
	a grammar.Assignable,
	cardinality string,
) (*ATNHandle, error) {
	switch v := a.(type) {
	case grammar.Keyword:
		return convertKeyword(rb, v, cardinality)
	case grammar.RuleCall:
		return convertRuleCall(rb, v, cardinality)
	case grammar.CrossRef:
		// The assignment wrapper's cardinality takes precedence over the CrossRef's
		// own cardinality when present (e.g. `Actions+=[Command:ID]+`).
		return convertCrossRef(rb, v, cardinality)
	case grammar.Alternatives:
		return convertAlternatives(rb, v, cardinality)
	default:
		return nil, nil
	}
}

func convertKeyword(
	rb ATNRuleBuilder,
	kw grammar.Keyword,
	cardinality string,
) (*ATNHandle, error) {
	name := kw.Value()
	id := rb.GetTokenTypeByName(name)
	if id == -1 {
		return nil, fmt.Errorf("unknown token %q", name)
	}
	handle := rb.TokenRef(id)
	handle.Left.ConsumedElement = kw
	lookaheadName := rb.GetLookaheadNameByElement(kw)
	return wrapWithCardinality(rb, handle, cardinality, lookaheadName, kw), nil
}

func convertRuleCall(
	rb ATNRuleBuilder,
	rc grammar.RuleCall,
	cardinality string,
) (*ATNHandle, error) {
	rule := rc.Rule().Ref(context.Background())
	lookaheadName := rb.GetLookaheadNameByElement(rc)

	switch typed := rule.(type) {
	case grammar.AbstractRuleWithBody:
		handle := rb.RuleRef(typed)
		handle.Left.RuleCallEntry = rc // tag so generator can find follow state via RuleCallEntry
		return wrapWithCardinality(rb, handle, cardinality, lookaheadName, rc), nil
	case grammar.AbstractTokenRule:
		id := rb.GetTokenTypeByName(typed.Name())
		termHandle := rb.TokenRef(id)
		termHandle.Left.ConsumedElement = rc
		return wrapWithCardinality(rb, termHandle, cardinality, lookaheadName, rc), nil
	}
	panic(fmt.Sprintf("unexpected rule type %T", rule))
}

func convertCrossRef(
	rb ATNRuleBuilder,
	cr grammar.CrossRef,
	cardinality string,
) (*ATNHandle, error) {
	// An enclosing assignment's cardinality (e.g. `Actions+=[Command:ID]+`) takes
	// precedence over the CrossRef's own; fall back to the CrossRef's when the
	// caller passes none.
	if cardinality == "" {
		cardinality = cr.Cardinality()
	}
	rule := cr.Rule().Rule().Ref(context.Background())
	hint := completionHintFor(rb.Rule(), cr)
	if abstractRule, ok := rule.(grammar.AbstractRuleWithBody); ok {
		handle := rb.RuleRef(abstractRule)
		if hint != nil {
			for _, t := range handle.Left.Transitions {
				if rt, ok := t.(*RuleTransition); ok && rt.Rule == abstractRule {
					rt.CompletionHint = hint
				}
			}
		}
		return handle, nil
	}
	id := rb.GetTokenTypeByName(rule.Name())
	termHandle := rb.TokenRef(id)
	termHandle.Left.ConsumedElement = cr.Rule() // tag with inner RuleCall to match generator naming
	if hint != nil {
		for _, t := range termHandle.Left.Transitions {
			if at, ok := t.(*AtomTransition); ok && at.TokenTypeId == id {
				at.CompletionHint = hint
			}
		}
	}
	lookaheadName := rb.GetLookaheadNameByElement(cr)
	return wrapWithCardinality(rb, termHandle, cardinality, lookaheadName, cr), nil
}

func convertAlternatives(
	rb ATNRuleBuilder,
	alts grammar.Alternatives,
	cardinality string,
) (*ATNHandle, error) {
	handles := make([]*ATNHandle, 0, len(alts.Alts()))
	for _, alt := range alts.Alts() {
		handle, err := convertElement(rb, alt)
		if err != nil {
			return nil, err
		}
		if handle == nil {
			continue
		}
		handles = append(handles, handle)
	}
	start := rb.NewState(parser.ATNBasic)
	lookaheadName := rb.GetLookaheadNameByElement(alts)
	handle := rb.MakeAlternatives(lookaheadName, start, handles)
	// `start` is the alternative-choice decision state for this Alternatives.
	rb.RecordOrDecision(alts, start)
	return wrapWithCardinality(rb, handle, cardinality, lookaheadName, alts), nil
}

func convertGroup(
	rb ATNRuleBuilder,
	g grammar.Group,
) (*ATNHandle, error) {
	elementHandles := make([]*ATNHandle, 0, len(g.Elements()))
	for _, element := range g.Elements() {
		elementHandle, err := convertElement(rb, element)
		if err != nil {
			return nil, err
		}
		if elementHandle == nil {
			continue
		}
		elementHandles = append(elementHandles, elementHandle)
	}
	handle := rb.MakeConcatenation(elementHandles)
	lookaheadName := rb.GetLookaheadNameByElement(g)
	return wrapWithCardinality(rb, handle, g.Cardinality(), lookaheadName, g), nil
}

func wrapWithCardinality(rb ATNRuleBuilder, prod *ATNHandle, cardinality string, lookaheadName string, el grammar.Element) *ATNHandle {
	switch cardinality {
	case "?":
		h := rb.Optional(lookaheadName, prod)
		// Optional returns the same handle; its Left is the decision state.
		rb.RecordLoopDecision(el, h.Left)
		return h
	case "*":
		h := rb.Star(lookaheadName, prod)
		// Star returns {Left: entry, ...}; entry is the decision state.
		rb.RecordLoopDecision(el, h.Left)
		return h
	case "+":
		h := rb.Plus(lookaheadName, prod)
		// Plus returns {Left: blkStart, ...}; the PlusLoopBack (blkStart.Loopback)
		// is the decision state.
		rb.RecordLoopDecision(el, h.Left.Loopback)
		return h
	default:
		return prod
	}
}

func nextCounter(counters map[reflect.Type]int, obj any) int {
	// returns the next 1-based counter value for kind and increments it
	kind := reflect.TypeOf(obj)
	counters[kind]++
	return counters[kind]
}
