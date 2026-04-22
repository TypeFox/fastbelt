// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/parser/allstar"
)

// tokenTypes used in all convert tests.
var testTokenTypes = map[string]TokenInfo{
	"a":   {ID: 1},
	"b":   {ID: 2},
	"ID":  {ID: 3, CategoryMatches: []int{10}},
	"int": {ID: 4},
}

// makeKeyword creates a grammar.Keyword with the given value and cardinality.
func makeKeyword(value, cardinality string) grammar.Keyword {
	kw := grammar.NewKeyword()
	kw.SetValue(makeToken(value))
	if cardinality != "" {
		kw.SetCardinality(makeToken(cardinality))
	}
	return kw
}

func makeToken(image string) *core.Token {
	return &core.Token{Image: image, TypeId: 0}
}

// makeRuleWithBody wraps elements into a ParserRule body.
// For a single-element body, it wraps directly; for multiple elements use a
// Group with no cardinality.
func makeParserRule(name string, body grammar.Element) grammar.ParserRule {
	r := grammar.NewParserRule()
	r.SetName(makeToken(name))
	r.SetBody(body)
	return r
}

// ────────────────────────────────────────────────────────────────────────────

func TestConvert_SingleKeyword(t *testing.T) {
	rules := []grammar.ParserRule{
		makeParserRule("R", makeKeyword("a", "")),
	}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Len(t, result[0].Definition, 1)

	term, ok := result[0].Definition[0].(*allstar.Terminal)
	require.True(t, ok)
	assert.Equal(t, 1, term.TokenTypeID)
	assert.Equal(t, 1, term.Idx)
}

func TestConvert_Group_Option(t *testing.T) {
	g := grammar.NewGroup()
	g.SetCardinality(makeToken("?"))
	g.SetElementsItem(makeKeyword("a", ""))

	rules := []grammar.ParserRule{makeParserRule("R", g)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	require.Len(t, result[0].Definition, 1)

	opt, ok := result[0].Definition[0].(*allstar.Option)
	require.True(t, ok)
	assert.Equal(t, 1, opt.Idx)
	require.Len(t, opt.Definition, 1)
}

func TestConvert_Group_Repetition(t *testing.T) {
	g := grammar.NewGroup()
	g.SetCardinality(makeToken("*"))
	g.SetElementsItem(makeKeyword("a", ""))

	rules := []grammar.ParserRule{makeParserRule("R", g)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	rep, ok := result[0].Definition[0].(*allstar.Repetition)
	require.True(t, ok)
	assert.Equal(t, 1, rep.Idx)
}

func TestConvert_Group_RepetitionMandatory(t *testing.T) {
	g := grammar.NewGroup()
	g.SetCardinality(makeToken("+"))
	g.SetElementsItem(makeKeyword("a", ""))

	rules := []grammar.ParserRule{makeParserRule("R", g)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	rep, ok := result[0].Definition[0].(*allstar.RepetitionMandatory)
	require.True(t, ok)
	assert.Equal(t, 1, rep.Idx)
}

func TestConvert_Group_Sequence(t *testing.T) {
	// Group with no cardinality → inline sequence.
	g := grammar.NewGroup()
	g.SetElementsItem(makeKeyword("a", ""))
	g.SetElementsItem(makeKeyword("b", ""))

	rules := []grammar.ParserRule{makeParserRule("R", g)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	// Two terminals, no wrapper.
	require.Len(t, result[0].Definition, 2)
	_, ok1 := result[0].Definition[0].(*allstar.Terminal)
	_, ok2 := result[0].Definition[1].(*allstar.Terminal)
	assert.True(t, ok1)
	assert.True(t, ok2)
}

func TestConvert_Alternatives(t *testing.T) {
	alts := grammar.NewAlternatives()
	alts.SetAltsItem(makeKeyword("a", ""))
	alts.SetAltsItem(makeKeyword("b", ""))

	rules := []grammar.ParserRule{makeParserRule("R", alts)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	require.Len(t, result[0].Definition, 1)

	alt, ok := result[0].Definition[0].(*allstar.Alternation)
	require.True(t, ok)
	assert.Equal(t, 1, alt.Idx)
	assert.Len(t, alt.Alternatives, 2)
}

func TestConvert_MultipleAlternations(t *testing.T) {
	// A Group with two Alternatives elements → two Alternations with Idx 1 and 2.
	g := grammar.NewGroup()

	alts1 := grammar.NewAlternatives()
	alts1.SetAltsItem(makeKeyword("a", ""))
	alts1.SetAltsItem(makeKeyword("b", ""))
	g.SetElementsItem(alts1)

	alts2 := grammar.NewAlternatives()
	alts2.SetAltsItem(makeKeyword("a", ""))
	alts2.SetAltsItem(makeKeyword("b", ""))
	g.SetElementsItem(alts2)

	rules := []grammar.ParserRule{makeParserRule("R", g)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	require.Len(t, result[0].Definition, 2)

	alt1 := result[0].Definition[0].(*allstar.Alternation)
	alt2 := result[0].Definition[1].(*allstar.Alternation)
	assert.Equal(t, 1, alt1.Idx)
	assert.Equal(t, 2, alt2.Idx)
}

func TestConvert_Assignment_Transparent(t *testing.T) {
	// Assignment wrapping a Keyword → Terminal (assignment stripped).
	assign := grammar.NewAssignment()
	assign.SetValue(makeKeyword("a", ""))

	rules := []grammar.ParserRule{makeParserRule("R", assign)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)
	require.Len(t, result[0].Definition, 1)
	_, ok := result[0].Definition[0].(*allstar.Terminal)
	assert.True(t, ok, "assignment should be transparent, yielding a Terminal")
}

func TestConvert_MixedCounters(t *testing.T) {
	// Rule with two Alternatives and two Option groups → indices 1 and 2 per kind.
	g := grammar.NewGroup()

	alts1 := grammar.NewAlternatives()
	alts1.SetAltsItem(makeKeyword("a", ""))
	alts1.SetAltsItem(makeKeyword("b", ""))
	g.SetElementsItem(alts1)

	opt1 := grammar.NewGroup()
	opt1.SetCardinality(makeToken("?"))
	opt1.SetElementsItem(makeKeyword("a", ""))
	g.SetElementsItem(opt1)

	alts2 := grammar.NewAlternatives()
	alts2.SetAltsItem(makeKeyword("a", ""))
	alts2.SetAltsItem(makeKeyword("b", ""))
	g.SetElementsItem(alts2)

	opt2 := grammar.NewGroup()
	opt2.SetCardinality(makeToken("?"))
	opt2.SetElementsItem(makeKeyword("b", ""))
	g.SetElementsItem(opt2)

	rules := []grammar.ParserRule{makeParserRule("R", g)}
	result, err := FromParserRules(rules, testTokenTypes)
	require.NoError(t, err)

	prods := result[0].Definition
	require.Len(t, prods, 4)

	assert.Equal(t, 1, prods[0].(*allstar.Alternation).Idx)
	assert.Equal(t, 1, prods[1].(*allstar.Option).Idx)
	assert.Equal(t, 2, prods[2].(*allstar.Alternation).Idx)
	assert.Equal(t, 2, prods[3].(*allstar.Option).Idx)
}

func TestConvert_UnknownToken(t *testing.T) {
	rules := []grammar.ParserRule{
		makeParserRule("R", makeKeyword("unknown_token", "")),
	}
	_, err := FromParserRules(rules, testTokenTypes)
	assert.Error(t, err, "unknown token should produce an error")
}

func TestConvert_RuleCall(t *testing.T) {
	// ruleB is referenced from ruleA via RuleCall.
	ruleB := grammar.NewParserRule()
	ruleB.SetName(makeToken("B"))
	ruleB.SetBody(makeKeyword("a", ""))

	rc := grammar.NewRuleCall()
	ref := core.NewReference[grammar.AbstractRule](nil, makeToken("B"), func(_ context.Context, _ *core.Reference[grammar.AbstractRule]) (*core.AstNodeDescription, *core.ReferenceError) {
		return nil, nil
	})
	rc.SetRule(ref)

	ruleA := makeParserRule("A", rc)

	result, err := FromParserRules([]grammar.ParserRule{ruleA, ruleB}, testTokenTypes)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Len(t, result[0].Definition, 1)

	nt, ok := result[0].Definition[0].(*allstar.NonTerminal)
	require.True(t, ok)
	assert.Equal(t, "B", nt.ReferencedRule.Name)
	assert.Equal(t, 1, nt.Idx)
}

func TestConvert_ForwardRef(t *testing.T) {
	// ruleA references ruleB which is defined after ruleA.
	ruleB := grammar.NewParserRule()
	ruleB.SetName(makeToken("B"))
	ruleB.SetBody(makeKeyword("b", ""))

	rc := grammar.NewRuleCall()
	ref := core.NewReference[grammar.AbstractRule](nil, makeToken("B"), func(_ context.Context, _ *core.Reference[grammar.AbstractRule]) (*core.AstNodeDescription, *core.ReferenceError) {
		return nil, nil
	})
	rc.SetRule(ref)

	ruleA := makeParserRule("A", rc)

	// Pass ruleA first, ruleB second — forward reference.
	result, err := FromParserRules([]grammar.ParserRule{ruleA, ruleB}, testTokenTypes)
	require.NoError(t, err)

	nt, ok := result[0].Definition[0].(*allstar.NonTerminal)
	require.True(t, ok)
	assert.Equal(t, "B", nt.ReferencedRule.Name)
}
