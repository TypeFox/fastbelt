// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"fmt"

	"typefox.dev/fastbelt/internal/grammar"
	allstar "typefox.dev/fastbelt/parser/allstar"
)

// TokenInfo carries the numeric type ID and optional category-match IDs for
// one token type.  The ID must match the constant emitted by GenerateLexer for
// the same terminal or keyword.
type TokenInfo struct {
	ID              int
	CategoryMatches []int
}

// buildTokenTypes constructs the token-type map required by FromParserRules.
// Keyword IDs are assigned first (sorted by raw value, matching GenerateLexer),
// followed by terminal IDs.
func buildTokenTypes(grammr grammar.Grammar) map[string]TokenInfo {
	m := make(map[string]TokenInfo)
	id := 1
	for _, kw := range GetAllKeywords(grammr) {
		m[kw.Value()] = TokenInfo{ID: id}
		id++
	}
	for _, tok := range grammr.Terminals() {
		m[tok.Name()] = TokenInfo{ID: id}
		id++
	}
	return m
}

// BuildRuntimeATNFromGrammar converts a Grammar into a RuntimeATN ready for
// the ALL(*) prediction engine.  It mirrors the token-ID assignment used by
// GenerateLexer so that the emitted ATN and the generated lexer agree on every
// numeric token type ID.
func BuildRuntimeATNFromGrammar(grammr grammar.Grammar) (*allstar.RuntimeATN, error) {
	tokenTypes := buildTokenTypes(grammr)
	rules, err := FromParserRules(grammr.Rules(), tokenTypes)
	if err != nil {
		return nil, fmt.Errorf("ATN build: %w", err)
	}
	return allstar.BuildRuntimeATN(allstar.CreateATN(rules)), nil
}

// FromParserRules converts a slice of grammar.ParserRule into the allstar Rule
// slice that CreateATN expects.  tokenTypes maps terminal/keyword name →
// TokenInfo.
func FromParserRules(
	rules []grammar.ParserRule,
	tokenTypes map[string]TokenInfo,
) ([]*allstar.Rule, error) {
	// First pass: create all Rule objects so forward references can be resolved.
	rulesByName := map[string]*allstar.Rule{}
	allstarRules := make([]*allstar.Rule, 0, len(rules))
	for _, gr := range rules {
		r := &allstar.Rule{Name: gr.Name()}
		rulesByName[gr.Name()] = r
		allstarRules = append(allstarRules, r)
	}

	// Second pass: convert bodies.
	for i, gr := range rules {
		r := allstarRules[i]
		if gr.Body() == nil {
			continue
		}
		counters := map[allstar.ProductionKind]int{}
		prods, err := convertElement(gr.Body(), counters, tokenTypes, rulesByName)
		if err != nil {
			return nil, fmt.Errorf("rule %q: %w", gr.Name(), err)
		}
		r.Definition = prods
	}

	return allstarRules, nil
}

func convertElement(
	el grammar.Element,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*allstar.Rule,
) ([]allstar.Production, error) {
	switch e := el.(type) {
	case grammar.Alternatives:
		return convertAlternatives(e, counters, tokenTypes, rulesByName)
	case grammar.Group:
		return convertGroup(e, counters, tokenTypes, rulesByName)
	case grammar.Assignment:
		if e.Value() == nil {
			return nil, nil
		}
		return convertAssignable(e.Value(), e.Cardinality(), counters, tokenTypes, rulesByName)
	case grammar.CrossRef:
		return convertCrossRef(e, counters, tokenTypes)
	case grammar.RuleCall:
		return convertRuleCall(e, e.Cardinality(), counters, tokenTypes, rulesByName)
	case grammar.Keyword:
		return convertKeyword(e, e.Cardinality(), counters, tokenTypes)
	case grammar.Action:
		return nil, nil
	default:
		return nil, nil
	}
}

func convertAssignable(
	a grammar.Assignable,
	cardinality string,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*allstar.Rule,
) ([]allstar.Production, error) {
	switch v := a.(type) {
	case grammar.Keyword:
		return convertKeyword(v, cardinality, counters, tokenTypes)
	case grammar.RuleCall:
		return convertRuleCall(v, cardinality, counters, tokenTypes, rulesByName)
	case grammar.CrossRef:
		return convertCrossRef(v, counters, tokenTypes)
	case grammar.Alternatives:
		return convertAlternatives(v, counters, tokenTypes, rulesByName)
	default:
		return nil, nil
	}
}

func convertKeyword(
	kw grammar.Keyword,
	cardinality string,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
) ([]allstar.Production, error) {
	name := kw.Value()
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown token %q", name)
	}
	term := &allstar.Terminal{
		TokenName:       name,
		TokenTypeID:     info.ID,
		CategoryMatches: info.CategoryMatches,
		Idx:             nextCounter(counters, allstar.ProdTerminal),
	}
	return wrapWithCardinality(term, cardinality, counters)
}

func convertRuleCall(
	rc grammar.RuleCall,
	cardinality string,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*allstar.Rule,
) ([]allstar.Production, error) {
	if rc.Rule() == nil {
		return nil, fmt.Errorf("RuleCall has no rule reference")
	}
	name := rc.Rule().Text
	if name == "" {
		return nil, fmt.Errorf("RuleCall rule reference has empty text")
	}

	if rule, ok := rulesByName[name]; ok {
		nt := &allstar.NonTerminal{
			ReferencedRule: rule,
			Idx:            nextCounter(counters, allstar.ProdNonTerminal),
		}
		return wrapWithCardinality(nt, cardinality, counters)
	}

	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown rule or token %q", name)
	}
	term := &allstar.Terminal{
		TokenName:       name,
		TokenTypeID:     info.ID,
		CategoryMatches: info.CategoryMatches,
		Idx:             nextCounter(counters, allstar.ProdTerminal),
	}
	return wrapWithCardinality(term, cardinality, counters)
}

func convertCrossRef(
	cr grammar.CrossRef,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
) ([]allstar.Production, error) {
	name := "ID"
	if cr.Rule() != nil && cr.Rule().Rule() != nil && cr.Rule().Rule().Text != "" {
		name = cr.Rule().Rule().Text
	}
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("cross-reference token %q not found in tokenTypes", name)
	}
	term := &allstar.Terminal{
		TokenName:       name,
		TokenTypeID:     info.ID,
		CategoryMatches: info.CategoryMatches,
		Idx:             nextCounter(counters, allstar.ProdTerminal),
	}
	return []allstar.Production{term}, nil
}

func convertAlternatives(
	alts grammar.Alternatives,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*allstar.Rule,
) ([]allstar.Production, error) {
	// Mirror the pre-order counter assignment used by populateContextWithNode:
	//   1. Alternation counter (only when there are 2+ alternatives)
	//   2. Cardinality wrapper counter (if cardinality is set)
	//   3. Recurse into alternatives
	// This ensures the ATN decision-map keys agree with the keys embedded in the
	// generated predict / predictOpt calls.
	var altIdx int
	if len(alts.Alts()) > 1 {
		altIdx = nextCounter(counters, allstar.ProdAlternation)
	}

	var wrapKind allstar.ProductionKind
	var wrapIdx int
	switch alts.Cardinality() {
	case "?":
		wrapKind = allstar.ProdOption
		wrapIdx = nextCounter(counters, wrapKind)
	case "*":
		wrapKind = allstar.ProdRepetition
		wrapIdx = nextCounter(counters, wrapKind)
	case "+":
		wrapKind = allstar.ProdRepetitionMandatory
		wrapIdx = nextCounter(counters, wrapKind)
	}

	alternation := &allstar.Alternation{Idx: altIdx}
	for _, alt := range alts.Alts() {
		prods, err := convertElement(alt, counters, tokenTypes, rulesByName)
		if err != nil {
			return nil, err
		}
		alternation.Alternatives = append(alternation.Alternatives, &allstar.Alternative{Definition: prods})
	}

	switch wrapKind {
	case allstar.ProdOption:
		return []allstar.Production{&allstar.Option{
			Definition: []allstar.Production{alternation},
			Idx:        wrapIdx,
		}}, nil
	case allstar.ProdRepetition:
		return []allstar.Production{&allstar.Repetition{
			Definition: []allstar.Production{alternation},
			Idx:        wrapIdx,
		}}, nil
	case allstar.ProdRepetitionMandatory:
		return []allstar.Production{&allstar.RepetitionMandatory{
			Definition: []allstar.Production{alternation},
			Idx:        wrapIdx,
		}}, nil
	default:
		return []allstar.Production{alternation}, nil
	}
}

func convertGroup(
	g grammar.Group,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*allstar.Rule,
) ([]allstar.Production, error) {
	switch g.Cardinality() {
	case "?":
		opt := &allstar.Option{Idx: nextCounter(counters, allstar.ProdOption)}
		for _, child := range g.Elements() {
			childProds, err := convertElement(child, counters, tokenTypes, rulesByName)
			if err != nil {
				return nil, err
			}
			opt.Definition = append(opt.Definition, childProds...)
		}
		return []allstar.Production{opt}, nil
	case "*":
		rep := &allstar.Repetition{Idx: nextCounter(counters, allstar.ProdRepetition)}
		for _, child := range g.Elements() {
			childProds, err := convertElement(child, counters, tokenTypes, rulesByName)
			if err != nil {
				return nil, err
			}
			rep.Definition = append(rep.Definition, childProds...)
		}
		return []allstar.Production{rep}, nil
	case "+":
		rep := &allstar.RepetitionMandatory{Idx: nextCounter(counters, allstar.ProdRepetitionMandatory)}
		for _, child := range g.Elements() {
			childProds, err := convertElement(child, counters, tokenTypes, rulesByName)
			if err != nil {
				return nil, err
			}
			rep.Definition = append(rep.Definition, childProds...)
		}
		return []allstar.Production{rep}, nil
	default:
		var prods []allstar.Production
		for _, child := range g.Elements() {
			childProds, err := convertElement(child, counters, tokenTypes, rulesByName)
			if err != nil {
				return nil, err
			}
			prods = append(prods, childProds...)
		}
		return prods, nil
	}
}

func wrapWithCardinality(prod allstar.Production, cardinality string, counters map[allstar.ProductionKind]int) ([]allstar.Production, error) {
	switch cardinality {
	case "?":
		return []allstar.Production{&allstar.Option{
			Definition: []allstar.Production{prod},
			Idx:        nextCounter(counters, allstar.ProdOption),
		}}, nil
	case "*":
		return []allstar.Production{&allstar.Repetition{
			Definition: []allstar.Production{prod},
			Idx:        nextCounter(counters, allstar.ProdRepetition),
		}}, nil
	case "+":
		return []allstar.Production{&allstar.RepetitionMandatory{
			Definition: []allstar.Production{prod},
			Idx:        nextCounter(counters, allstar.ProdRepetitionMandatory),
		}}, nil
	default:
		return []allstar.Production{prod}, nil
	}
}

func nextCounter(counters map[allstar.ProductionKind]int, kind allstar.ProductionKind) int {
	counters[kind]++
	return counters[kind]
}
