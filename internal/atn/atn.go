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
		return convertCrossRef(rb, e)

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
		//TODO true?
		// CrossRef cardinality comes from the outer element, not from the Assignable itself.
		// The CrossRef has its own Cardinality() but the assignment wrapper's cardinality
		// takes precedence when present.
		return convertCrossRef(rb, v)
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
	lookaheadName := rb.GetLookaheadNameByElement(kw)
	return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
}

func convertRuleCall(
	rb ATNRuleBuilder,
	rc grammar.RuleCall,
	cardinality string,
) (*ATNHandle, error) {
	rule := rc.Rule().Ref(context.Background())
	lookaheadName := rb.GetLookaheadNameByElement(rc)

	switch typed := rule.(type) {
	case grammar.CompositeRule:
		handle := rb.RuleRef(typed)
		return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
	case grammar.ParserRule:
		handle := rb.RuleRef(typed)
		return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
	case grammar.Token:
		id := rb.GetTokenTypeByName(typed.Name())
		termHandle := rb.TokenRef(id)
		return wrapWithCardinality(rb, termHandle, cardinality, lookaheadName), nil
	}
	panic(fmt.Sprintf("unexpected rule type %T", rule))
}

func convertCrossRef(
	rb ATNRuleBuilder,
	cr grammar.CrossRef,
) (*ATNHandle, error) {
	name := cr.Rule().Rule().Text()
	id := rb.GetTokenTypeByName(name)
	termHandle := rb.TokenRef(id)
	lookaheadName := rb.GetLookaheadNameByElement(cr)
	//TODO true? CrossRef cardinality comes from the outer element, not from the CrossRef itself.
	return wrapWithCardinality(rb, termHandle, cr.Cardinality(), lookaheadName), nil
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
	return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
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
	return wrapWithCardinality(rb, handle, g.Cardinality(), lookaheadName), nil
}

func wrapWithCardinality(rb ATNRuleBuilder, prod *ATNHandle, cardinality string, lookaheadName string) *ATNHandle {
	switch cardinality {
	case "?":
		return rb.Optional(lookaheadName, prod)
	case "*":
		return rb.Star(lookaheadName, prod)
	case "+":
		return rb.Plus(lookaheadName, prod)
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
