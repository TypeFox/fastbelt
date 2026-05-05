// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"fmt"
	"reflect"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
)

func CreateATN(grammr grammar.Grammar) (*ATN, map[string]*grammar.ParserRule, map[string]TokenType) {
	tokenTypes := GetTokenTypes(grammr)
	lookaheadNames := GetLookaheadNames(grammr)

	builder := NewATNBuilder()

	entries := map[grammar.ParserRule]ATNRuleBuilder{}
	byName := map[string]*grammar.ParserRule{}
	for _, gr := range grammr.Rules() {
		entries[gr] = builder.DeclareRule(&gr)
		byName[gr.Name()] = &gr
	}
	for _, gr := range grammr.Rules() {
		ruleBuilder := entries[gr]
		handle, err := convertElement(ruleBuilder, gr.Body(), lookaheadNames, tokenTypes, byName)
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
	return atn, byName, tokenTypes
}

func GetLookaheadNames(grammr grammar.Grammar) map[grammar.Element]string {
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

func GetLookaheadName(el grammar.Element, names map[grammar.Element]string) string {
	if name, ok := names[el]; ok {
		return name
	}
	panic("no lookahead name found for element")
}

func GetTokenTypes(grammr grammar.Grammar) map[string]TokenType {
	tokens := grammr.Terminals()
	keywords := GetAllKeywords(grammr)
	nodes := make(map[string]TokenType, len(tokens)+len(keywords))
	id := 1
	for _, keyword := range keywords {
		nodes[keyword.Text()] = TokenType{ID: id}
		id++
	}
	for _, token := range tokens {
		nodes[token.Name()] = TokenType{ID: id}
		id++
	}
	return nodes
}

func convertElement(
	rb ATNRuleBuilder,
	el grammar.Element,
	names map[grammar.Element]string,
	tokenTypes map[string]TokenType,
	rulesByName map[string]*grammar.ParserRule,
) (*ATNHandle, error) {
	switch e := el.(type) {
	case grammar.Alternatives:
		return convertAlternatives(rb, e, e.Cardinality(), names, tokenTypes, rulesByName)

	case grammar.Group:
		return convertGroup(rb, e, names, tokenTypes, rulesByName)

	case grammar.Assignment:
		// Assignment is transparent: recurse into Value().
		if e.Value() == nil {
			return nil, nil
		}
		return convertAssignable(rb, e.Value(), e.Cardinality(), names, tokenTypes, rulesByName)

	case grammar.CrossRef:
		return convertCrossRef(rb, e, names, tokenTypes)

	case grammar.RuleCall:
		return convertRuleCall(rb, e, e.Cardinality(), names, tokenTypes, rulesByName)

	case grammar.Keyword:
		return convertKeyword(rb, e, e.Cardinality(), names, tokenTypes)

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
	names map[grammar.Element]string,
	tokenTypes map[string]TokenType,
	rulesByName map[string]*grammar.ParserRule,
) (*ATNHandle, error) {
	switch v := a.(type) {
	case grammar.Keyword:
		return convertKeyword(rb, v, cardinality, names, tokenTypes)
	case grammar.RuleCall:
		return convertRuleCall(rb, v, cardinality, names, tokenTypes, rulesByName)
	case grammar.CrossRef:
		//TODO true?
		// CrossRef cardinality comes from the outer element, not from the Assignable itself.
		// The CrossRef has its own Cardinality() but the assignment wrapper's cardinality
		// takes precedence when present.
		return convertCrossRef(rb, v, names, tokenTypes)
	case grammar.Alternatives:
		return convertAlternatives(rb, v, cardinality, names, tokenTypes, rulesByName)
	default:
		return nil, nil
	}
}

func convertKeyword(
	rb ATNRuleBuilder,
	kw grammar.Keyword,
	cardinality string,
	lookaheadNames map[grammar.Element]string,
	tokenTypes map[string]TokenType,
) (*ATNHandle, error) {
	name := kw.Value()
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown token %q", name)
	}
	handle := rb.TokenRef(info)
	lookaheadName := GetLookaheadName(kw, lookaheadNames)
	return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
}

func convertRuleCall(
	rb ATNRuleBuilder,
	rc grammar.RuleCall,
	cardinality string,
	names map[grammar.Element]string,
	tokenTypes map[string]TokenType,
	rulesByName map[string]*grammar.ParserRule,
) (*ATNHandle, error) {
	if rc.Rule() == nil {
		return nil, fmt.Errorf("RuleCall has no rule reference")
	}
	name := rc.Rule().Text()
	if name == "" {
		return nil, fmt.Errorf("RuleCall rule reference has empty text")
	}

	lookaheadName := GetLookaheadName(rc, names)

	if rule, ok := rulesByName[name]; ok {
		handle := rb.RuleRef(rule)
		return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
	}

	// Otherwise treat it as a terminal (lexer rule reference).
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown rule or token %q", name)
	}
	termHandle := rb.TokenRef(info)
	return wrapWithCardinality(rb, termHandle, cardinality, lookaheadName), nil
}

func convertCrossRef(
	rb ATNRuleBuilder,
	cr grammar.CrossRef,
	names map[grammar.Element]string,
	tokenTypes map[string]TokenType,
) (*ATNHandle, error) {
	// Use the explicitly named rule if present, otherwise fall back to "ID".
	name := "ID"
	if cr.Rule() != nil && cr.Rule().Rule() != nil && cr.Rule().Rule().Text() != "" {
		name = cr.Rule().Rule().Text()
	}
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("cross-reference token %q not found in tokenTypes", name)
	}
	termHandle := rb.TokenRef(info)
	lookaheadName := GetLookaheadName(cr, names)
	//TODO true? CrossRef cardinality comes from the outer element, not from the CrossRef itself.
	return wrapWithCardinality(rb, termHandle, cr.Cardinality(), lookaheadName), nil
}

func convertAlternatives(
	rb ATNRuleBuilder,
	alts grammar.Alternatives,
	cardinality string,
	names map[grammar.Element]string,
	tokenTypes map[string]TokenType,
	rulesByName map[string]*grammar.ParserRule,
) (*ATNHandle, error) {
	handles := make([]*ATNHandle, 0, len(alts.Alts()))
	for _, alt := range alts.Alts() {
		handle, err := convertElement(rb, alt, names, tokenTypes, rulesByName)
		if err != nil {
			return nil, err
		}
		if handle == nil {
			continue
		}
		handles = append(handles, handle)
	}
	start := rb.NewState(parser.ATNBasic)
	lookaheadName := GetLookaheadName(alts, names)
	handle := rb.MakeAlts(lookaheadName, start, handles)
	return wrapWithCardinality(rb, handle, cardinality, lookaheadName), nil
}

func convertGroup(
	rb ATNRuleBuilder,
	g grammar.Group,
	names map[grammar.Element]string,
	tokenTypes map[string]TokenType,
	rulesByName map[string]*grammar.ParserRule,
) (*ATNHandle, error) {
	elementHandles := make([]*ATNHandle, 0, len(g.Elements()))
	for _, element := range g.Elements() {
		elementHandle, err := convertElement(rb, element, names, tokenTypes, rulesByName)
		if err != nil {
			return nil, err
		}
		if elementHandle == nil {
			continue
		}
		elementHandles = append(elementHandles, elementHandle)
	}
	handle := rb.MakeBlock(elementHandles)
	lookaheadName := GetLookaheadName(g, names)
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
