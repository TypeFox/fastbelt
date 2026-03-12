// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"fmt"

	"typefox.dev/fastbelt/internal/grammar"
)

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
) ([]*Rule, error) {
	// First pass: create all Rule objects so forward references can be resolved.
	rulesByName := map[string]*Rule{}
	allstarRules := make([]*Rule, 0, len(rules))
	for _, gr := range rules {
		r := &Rule{Name: gr.Name()}
		rulesByName[gr.Name()] = r
		allstarRules = append(allstarRules, r)
	}

	// Second pass: convert bodies. NonTerminal.ReferencedRule is resolved via rulesByName.
	for i, gr := range rules {
		r := allstarRules[i]
		if gr.Body() == nil {
			continue
		}
		counters := map[ProductionKind]int{}
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
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*Rule,
) ([]Production, error) {
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
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*Rule,
) ([]Production, error) {
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
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
) ([]Production, error) {
	name := kw.Value()
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown token %q", name)
	}
	term := &Terminal{
		TokenName:       name,
		TokenTypeID:     info.ID,
		CategoryMatches: info.CategoryMatches,
		Idx:             nextCounter(counters, ProdTerminal),
	}
	return wrapWithCardinality(term, cardinality, counters)
}

func convertRuleCall(
	rc grammar.RuleCall,
	cardinality string,
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*Rule,
) ([]Production, error) {
	if rc.Rule() == nil {
		return nil, fmt.Errorf("RuleCall has no rule reference")
	}
	name := rc.Rule().Text
	if name == "" {
		return nil, fmt.Errorf("RuleCall rule reference has empty text")
	}

	// Check if the referenced rule is a parser rule.
	if rule, ok := rulesByName[name]; ok {
		nt := &NonTerminal{
			ReferencedRule: rule,
			Idx:            nextCounter(counters, ProdNonTerminal),
		}
		return wrapWithCardinality(nt, cardinality, counters)
	}

	// Otherwise treat it as a terminal (lexer rule reference).
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("unknown rule or token %q", name)
	}
	term := &Terminal{
		TokenName:       name,
		TokenTypeID:     info.ID,
		CategoryMatches: info.CategoryMatches,
		Idx:             nextCounter(counters, ProdTerminal),
	}
	return wrapWithCardinality(term, cardinality, counters)
}

func convertCrossRef(
	cr grammar.CrossRef,
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
) ([]Production, error) {
	// Use the explicitly named rule if present, otherwise fall back to "ID".
	name := "ID"
	if cr.Rule() != nil && cr.Rule().Rule() != nil && cr.Rule().Rule().Text != "" {
		name = cr.Rule().Rule().Text
	}
	info, ok := tokenTypes[name]
	if !ok {
		return nil, fmt.Errorf("cross-reference token %q not found in tokenTypes", name)
	}
	term := &Terminal{
		TokenName:       name,
		TokenTypeID:     info.ID,
		CategoryMatches: info.CategoryMatches,
		Idx:             nextCounter(counters, ProdTerminal),
	}
	return []Production{term}, nil
}

func convertAlternatives(
	alts grammar.Alternatives,
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*Rule,
) ([]Production, error) {
	alternatives := make([]*Alternative, 0, len(alts.Alts()))
	for _, alt := range alts.Alts() {
		prods, err := convertElement(alt, counters, tokenTypes, rulesByName)
		if err != nil {
			return nil, err
		}
		alternatives = append(alternatives, &Alternative{Definition: prods})
	}
	alternation := &Alternation{
		Alternatives: alternatives,
		Idx:          nextCounter(counters, ProdAlternation),
	}
	return []Production{alternation}, nil
}

func convertGroup(
	g grammar.Group,
	counters map[ProductionKind]int,
	tokenTypes map[string]TokenInfo,
	rulesByName map[string]*Rule,
) ([]Production, error) {
	// Convert children.
	var prods []Production
	for _, child := range g.Elements() {
		childProds, err := convertElement(child, counters, tokenTypes, rulesByName)
		if err != nil {
			return nil, err
		}
		prods = append(prods, childProds...)
	}

	switch g.Cardinality() {
	case "?":
		opt := &Option{
			Definition: prods,
			Idx:        nextCounter(counters, ProdOption),
		}
		return []Production{opt}, nil
	case "*":
		rep := &Repetition{
			Definition: prods,
			Idx:        nextCounter(counters, ProdRepetition),
		}
		return []Production{rep}, nil
	case "+":
		rep := &RepetitionMandatory{
			Definition: prods,
			Idx:        nextCounter(counters, ProdRepetitionMandatory),
		}
		return []Production{rep}, nil
	default:
		// No cardinality → inline sequence.
		return prods, nil
	}
}

// wrapWithCardinality wraps a leaf production with an optional/repetition based
// on the cardinality string. An empty cardinality returns [prod] as-is.
func wrapWithCardinality(prod Production, cardinality string, counters map[ProductionKind]int) ([]Production, error) {
	switch cardinality {
	case "?":
		return []Production{&Option{
			Definition: []Production{prod},
			Idx:        nextCounter(counters, ProdOption),
		}}, nil
	case "*":
		return []Production{&Repetition{
			Definition: []Production{prod},
			Idx:        nextCounter(counters, ProdRepetition),
		}}, nil
	case "+":
		return []Production{&RepetitionMandatory{
			Definition: []Production{prod},
			Idx:        nextCounter(counters, ProdRepetitionMandatory),
		}}, nil
	default:
		return []Production{prod}, nil
	}
}

// nextCounter returns the next 1-based counter value for kind and increments it.
func nextCounter(counters map[ProductionKind]int, kind ProductionKind) int {
	counters[kind]++
	return counters[kind]
}
