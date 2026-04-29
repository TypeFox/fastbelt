// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.
package generator

import (
	"fmt"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser/allstar"
)

func GenerateATN(grammr grammar.Grammar, packageName string) string {
	rules, err := FromParserRules(grammr.Rules(), GetTokenTypes(grammr))
	if err != nil {
		panic(err)
	}
	atn := allstar.CreateATN(rules)
	rtn := allstar.BuildRuntimeATN(atn)
	source := allstar.EmitGoSource(packageName, "BuildATN", "typefox.dev/fastbelt/parser/allstar", rtn)
	return FormatIfPossible(source.String())
}

func GetTokenTypes(grammr grammar.Grammar) map[string]TokenInfo {
	tokens := grammr.Terminals()
	keywords := GetAllKeywords(grammr)
	nodes := make(map[string]TokenInfo, len(tokens)+len(keywords))
	id := 1
	for _, keyword := range keywords {
		nodes[keyword.Text()] = TokenInfo{ID: id}
		id++
	}
	for _, token := range tokens {
		nodes[token.Name()] = TokenInfo{ID: id}
		id++
	}
	return nodes
}

// TokenInfo carries the type ID and category-match IDs for a token type.
type TokenInfo struct {
	ID              int
	CategoryMatches []int
}

// FromParserRules converts a slice of grammar.ParserRule into the allstar Rule
// slice that CreateATN expects. tokenTypes maps terminal name → TokenInfo.
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

	// Second pass: convert bodies. NonTerminal.ReferencedRule is resolved via rulesByName.
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

// convertElement converts a grammar.Element into zero or more allstar Productions.
// Multiple productions are returned for plain Groups (no cardinality) that inline
// their children into the parent sequence.
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
		// Assignment is transparent: recurse into Value().
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
		// Semantic actions have no ATN impact.
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
		// CrossRef cardinality comes from the outer element, not from the Assignable itself.
		// The CrossRef has its own Cardinality() but the assignment wrapper's cardinality
		// takes precedence when present.
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
	name := rc.Rule().Text()
	if name == "" {
		return nil, fmt.Errorf("RuleCall rule reference has empty text")
	}

	// Check if the referenced rule is a parser rule.
	if rule, ok := rulesByName[name]; ok {
		nt := &allstar.NonTerminal{
			ReferencedRule: rule,
			Idx:            nextCounter(counters, allstar.ProdNonTerminal),
		}
		return wrapWithCardinality(nt, cardinality, counters)
	}

	// Otherwise treat it as a terminal (lexer rule reference).
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
	// Use the explicitly named rule if present, otherwise fall back to "ID".
	name := "ID"
	if cr.Rule() != nil && cr.Rule().Rule() != nil && cr.Rule().Rule().Text() != "" {
		name = cr.Rule().Rule().Text()
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
	// Assign Idx before recursing so the occurrence index matches the
	// code generator's pre-order traversal.
	alternation := &allstar.Alternation{
		Idx: nextCounter(counters, allstar.ProdAlternation),
	}
	for _, alt := range alts.Alts() {
		prods, err := convertElement(alt, counters, tokenTypes, rulesByName)
		if err != nil {
			return nil, err
		}
		alternation.Alternatives = append(alternation.Alternatives, &allstar.Alternative{Definition: prods})
	}
	return []allstar.Production{alternation}, nil
}

func convertGroup(
	g grammar.Group,
	counters map[allstar.ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*allstar.Rule,
) ([]allstar.Production, error) {
	switch g.Cardinality() {
	case "?":
		// Assign Idx before recursing (pre-order) to match the code generator.
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
		// No cardinality → inline sequence.
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

// wrapWithCardinality wraps a leaf production with an optional/repetition based
// on the cardinality string. An empty cardinality returns [prod] as-is.
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

// nextCounter returns the next 1-based counter value for kind and increments it.
func nextCounter(counters map[allstar.ProductionKind]int, kind allstar.ProductionKind) int {
	counters[kind]++
	return counters[kind]
}
