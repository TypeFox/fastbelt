// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package atn

import (
	"context"
	"fmt"
	"strconv"

	fastbelt "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser"
)

// BuildElementNames returns a map from token-consuming grammar elements to their base name.
// The naming scheme mirrors the generator's populateContextWithNode:
// prefix starts as the rule name and grows with each Assignment property on the path.
// Keyword -> prefix + KeywordName; token RuleCall -> prefix + token.Name().
// CrossRef delegates to its inner RuleCall so the key matches the generator's accessNames key.
func BuildElementNames(grammr grammar.Grammar) map[grammar.Element]string {
	result := make(map[grammar.Element]string)
	for _, rule := range grammr.Rules() {
		collectElementNames(result, rule.Name(), rule.Body())
	}
	for _, composite := range grammr.Composites() {
		collectElementNames(result, composite.Name(), composite.Body())
	}
	return result
}

func collectElementNames(names map[grammar.Element]string, prefix string, node fastbelt.AstNode) {
	switch n := node.(type) {
	case grammar.Alternatives:
		for _, alt := range n.Alts() {
			collectElementNames(names, prefix, alt)
		}
	case grammar.Group:
		for _, elem := range n.Elements() {
			collectElementNames(names, prefix, elem)
		}
	case grammar.Keyword:
		names[n] = prefix + "_" + grammar.KeywordName(n)
	case grammar.Assignment:
		if n.Value() != nil {
			childPrefix := prefix + "_" + n.Property().Text()
			collectElementNames(names, childPrefix, n.Value())
		}
	case grammar.CrossRef:
		// CrossRef delegates to its inner RuleCall so the key matches what
		// generateCrossReferenceParser uses (crossRef.Rule()).
		collectElementNames(names, prefix, n.Rule())
	case grammar.RuleCall:
		ruleRef := n.Rule().Ref(context.Background())
		if token, ok := ruleRef.(grammar.Token); ok {
			names[n] = prefix + "_" + token.Name()
		}
	}
}

// BuildStateNameMap returns a slice of Go identifier constant names for every
// ATN state, indexed by state array index. Names follow the unified scheme:
//
//	Rule + Element [+ _Index]
//
// where Element is the grammar element name for token-consuming states, or the
// state-type role for structural states. The _Index suffix is appended only
// when two or more states share the same base name (omitted when unique).
func BuildStateNameMap(a *ATN, elementNames map[grammar.Element]string) []string {
	// First pass: compute base names for every state.
	bases := make([]string, len(a.States))
	for i, s := range a.States {
		bases[i] = stateBaseName(s, elementNames)
	}

	// Count occurrences of each base name.
	counts := make(map[string]int, len(bases))
	for _, b := range bases {
		counts[b]++
	}

	// Second pass: append _Index only for genuinely ambiguous bases.
	counters := make(map[string]int, len(bases))
	names := make([]string, len(a.States))
	for i, b := range bases {
		if counts[b] == 1 {
			names[i] = b
		} else {
			idx := counters[b]
			names[i] = b + "_" + strconv.Itoa(idx)
			counters[b]++
		}
	}
	return names
}

func stateBaseName(s *ATNState, elementNames map[grammar.Element]string) string {
	ruleName := ""
	if s.Rule != nil {
		ruleName = s.Rule.Name()
	}

	// Token-consuming state: use the property-style element name (already includes rule name).
	if s.ConsumedElement != nil {
		if name, ok := elementNames[s.ConsumedElement]; ok {
			return name
		}
	}

	// Fallback: rule name + state-type role.
	// Use double underscore to separate from token-consuming names
	return ruleName + "__" + atnStateRoleName(s.Type)
}

func atnStateRoleName(t parser.ATNStateType) string {
	switch t {
	case parser.ATNRuleStart:
		return "Start"
	case parser.ATNRuleStop:
		return "Stop"
	case parser.ATNBlockEnd:
		return "BlockEnd"
	case parser.ATNBasic:
		return "Basic"
	case parser.ATNLoopEntry:
		return "LoopEntry"
	case parser.ATNLoopBack:
		return "LoopBack"
	case parser.ATNLoopEnd:
		return "LoopEnd"
	default:
		return fmt.Sprintf("State%d", int(t))
	}
}
