// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
)

// tokenTypeA and tokenTypeB are fake token type IDs used across all integration tests.
const (
	tokenA = 1
	tokenB = 2
)

// mockTokenSource wraps a []int of token type IDs into a TokenSource.
// When the index is out of range, it returns the EOF sentinel.
type mockTokenSource struct {
	ids []int
}

func newMockSource(ids ...int) *mockTokenSource {
	return &mockTokenSource{ids: ids}
}

func (m *mockTokenSource) LA(offset int) *core.Token {
	idx := offset - 1
	if idx < 0 || idx >= len(m.ids) {
		return core.EOFToken
	}
	// Return a token whose TypeId matches the requested id.
	return &core.Token{TypeId: m.ids[idx]}
}

// ────────────────────────────────────────────────────────────────────────────
// 6.5.1 LL(*) lookahead (unbounded)
//
// LongRule := OR(
//
//	alt0: ε
//	alt1: AT_LEAST_ONE(A)
//	alt2: AT_LEAST_ONE(A) CONSUME(B)
//
// )
// ────────────────────────────────────────────────────────────────────────────

func buildLongRule() *LLStarLookahead {
	// alt0: empty alternative
	alt0 := &Alternative{}
	// alt1: AT_LEAST_ONE(A)
	alt1 := &Alternative{
		Definition: []Production{
			&RepetitionMandatory{
				Definition: []Production{
					&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1},
				},
				Idx: 1,
			},
		},
	}
	// alt2: AT_LEAST_ONE(A) CONSUME(B)
	alt2 := &Alternative{
		Definition: []Production{
			&RepetitionMandatory{
				Definition: []Production{
					&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 2},
				},
				Idx: 2,
			},
			&Terminal{TokenTypeID: tokenB, TokenName: "B", Idx: 1},
		},
	}

	rule := &Rule{
		Name: "LongRule",
		Definition: []Production{
			&Alternation{
				Alternatives: []*Alternative{alt0, alt1, alt2},
				Idx:          1,
			},
		},
	}
	return NewLLStarLookahead([]*Rule{rule}, nil)
}

func TestLL_Star_LongestAlt1(t *testing.T) {
	s := buildLongRule()
	// decision 0 is the Alternation (OR) at decision index 0
	// We need the decision index for the alternation in LongRule.
	decision := s.atn.DecisionMap["LongRule_Alternation_1"].Decision
	src := newMockSource(tokenA, tokenA, tokenA)
	result := s.AdaptivePredict(src, decision, EmptyPredicates)
	// Should pick alt1 (index 1) — greedy, no terminating B.
	assert.Equal(t, 1, result)
}

func TestLL_Star_LongestAlt2(t *testing.T) {
	s := buildLongRule()
	decision := s.atn.DecisionMap["LongRule_Alternation_1"].Decision
	src := newMockSource(tokenA, tokenA, tokenB)
	result := s.AdaptivePredict(src, decision, EmptyPredicates)
	// Should pick alt2 (index 2) — has terminating B.
	assert.Equal(t, 2, result)
}

func TestLL_Star_ShortestAlt(t *testing.T) {
	s := buildLongRule()
	decision := s.atn.DecisionMap["LongRule_Alternation_1"].Decision
	src := newMockSource() // empty input
	result := s.AdaptivePredict(src, decision, EmptyPredicates)
	// Should pick alt0 (index 0) — empty alternative.
	assert.Equal(t, 0, result)
}

// ────────────────────────────────────────────────────────────────────────────
// 6.5.2 Ambiguity detection
//
// Mirrors the AmbigiousParser test suite from atn.test.ts.
// ────────────────────────────────────────────────────────────────────────────

// ambigRules builds the full set of rules used in the ambiguity tests:
//
//	OptionRule       := OPTION(AT_LEAST_ONE(A)) AT_LEAST_ONE(A)
//	AltRule          := OR(SUBRULE(RuleB), SUBRULE(RuleC))
//	RuleB            := MANY(A)
//	RuleC            := MANY(A) OPTION(B)
//	AltRuleWithEOF   := OR(SUBRULE(RuleEOF), SUBRULE(RuleEOF))
//	RuleEOF          := MANY(A) CONSUME(EOF)
//	AltRuleWithPred  := OR(GATE(pred,CONSUME(A)), GATE(!pred,CONSUME(A)), CONSUME(B))
//	AltWithOption    := OR(CONSUME(A), CONSUME(B)) OPTION(CONSUME(A))
type ambigRules struct {
	s              *LLStarLookahead
	ambiguityMsgs  []string
	optionRule     *Rule
	altRule        *Rule
	altRuleEOF     *Rule
	altRulePred    *Rule
	altWithOption  *Rule
}

func buildAmbigRules() *ambigRules {
	msgs := &[]string{}
	logging := func(m string) { *msgs = append(*msgs, m) }

	// RuleB: MANY(A)
	ruleB := &Rule{Name: "RuleB", Definition: []Production{
		&Repetition{
			Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1}},
			Idx:        1,
		},
	}}

	// RuleC: MANY(A) OPTION(B)
	ruleC := &Rule{Name: "RuleC", Definition: []Production{
		&Repetition{
			Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1}},
			Idx:        1,
		},
		&Option{
			Definition: []Production{&Terminal{TokenTypeID: tokenB, TokenName: "B", Idx: 1}},
			Idx:        1,
		},
	}}

	// OptionRule: OPTION(AT_LEAST_ONE(A)) AT_LEAST_ONE(A)
	optionRule := &Rule{Name: "OptionRule", Definition: []Production{
		&Option{
			Definition: []Production{
				&RepetitionMandatory{
					Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1}},
					Idx:        1,
				},
			},
			Idx: 1,
		},
		&RepetitionMandatory{
			Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 2}},
			Idx:        2,
		},
	}}

	// AltRule: OR(SUBRULE(RuleB), SUBRULE(RuleC))
	altRule := &Rule{Name: "AltRule", Definition: []Production{
		&Alternation{
			Alternatives: []*Alternative{
				{Definition: []Production{&NonTerminal{ReferencedRule: ruleB, Idx: 1}}},
				{Definition: []Production{&NonTerminal{ReferencedRule: ruleC, Idx: 2}}},
			},
			Idx: 1,
		},
	}}

	// RuleEOF: MANY(A) CONSUME(EOF)  (EOF token ID is 0)
	ruleEOF := &Rule{Name: "RuleEOF", Definition: []Production{
		&Repetition{
			Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1}},
			Idx:        1,
		},
		&Terminal{TokenTypeID: core.EOF.Id, TokenName: "EOF", Idx: 2},
	}}

	// AltRuleWithEOF: OR(SUBRULE(RuleEOF), SUBRULE(RuleEOF))
	altRuleEOF := &Rule{Name: "AltRuleWithEOF", Definition: []Production{
		&Alternation{
			Alternatives: []*Alternative{
				{Definition: []Production{&NonTerminal{ReferencedRule: ruleEOF, Idx: 1}}},
				{Definition: []Production{&NonTerminal{ReferencedRule: ruleEOF, Idx: 2}}},
			},
			Idx: 1,
		},
	}}

	// AltRuleWithPred: OR(CONSUME(A), CONSUME(A), CONSUME(B)) — predicates applied externally
	altRulePred := &Rule{Name: "AltRuleWithPred", Definition: []Production{
		&Alternation{
			Alternatives: []*Alternative{
				{Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1}}},
				{Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 2}}},
				{Definition: []Production{&Terminal{TokenTypeID: tokenB, TokenName: "B", Idx: 2}}},
			},
			Idx: 1,
		},
	}}

	// AltWithOption: OR(CONSUME(A), CONSUME(B)) OPTION(CONSUME(A))
	// The rule returns "intermediate + option":
	//   OR returns 2 for A, 4 for B
	//   OPTION returns 1 if A follows, 0 otherwise
	// So B,A → 4+1 = 5
	altWithOption := &Rule{Name: "AltWithOption", Definition: []Production{
		&Alternation{
			Alternatives: []*Alternative{
				{Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 1}}},
				{Definition: []Production{&Terminal{TokenTypeID: tokenB, TokenName: "B", Idx: 1}}},
			},
			Idx: 1,
		},
		&Option{
			Definition: []Production{&Terminal{TokenTypeID: tokenA, TokenName: "A", Idx: 2}},
			Idx:        1,
		},
	}}

	allRules := []*Rule{ruleB, ruleC, optionRule, altRule, ruleEOF, altRuleEOF, altRulePred, altWithOption}
	s := NewLLStarLookahead(allRules, &LLStarLookaheadOptions{
		Logging: logging,
	})

	return &ambigRules{
		s:             s,
		ambiguityMsgs: *msgs,
		optionRule:    optionRule,
		altRule:       altRule,
		altRuleEOF:    altRuleEOF,
		altRulePred:   altRulePred,
		altWithOption: altWithOption,
	}
}

func TestAmbig_OptionAmbiguity(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	optionDecision := ar.s.atn.DecisionMap["OptionRule_Option_1"].Decision
	src := newMockSource(tokenA, tokenA, tokenA)
	// The OPTION is taken when there's input (ambiguous with the AT_LEAST_ONE outside).
	result := ar.s.AdaptivePredict(src, optionDecision, EmptyPredicates)
	// OPTION taken → alt 0
	assert.Equal(t, 0, result, "option should be taken (alt 0)")
	assert.True(t, len(*reports) > 0, "expected ambiguity report for OPTION")
	assert.Contains(t, (*reports)[0], "<0, 1> in <OPTION>")
}

func TestAmbig_FirstAltOnAmbiguity(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRule_Alternation_1"].Decision
	src := newMockSource(tokenA, tokenA, tokenA)
	result := ar.s.AdaptivePredict(src, altDecision, EmptyPredicates)
	assert.Equal(t, 0, result, "first alt should win on ambiguity")
	assert.True(t, len(*reports) > 0)
	assert.Contains(t, (*reports)[0], "<0, 1> in <OR>")
}

func TestAmbig_EOFAmbiguity(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRuleWithEOF_Alternation_1"].Decision
	src := newMockSource() // empty input → EOF
	result := ar.s.AdaptivePredict(src, altDecision, EmptyPredicates)
	assert.Equal(t, 0, result, "first alt should win on EOF ambiguity")
	assert.True(t, len(*reports) > 0)
	assert.Contains(t, (*reports)[0], "<0, 1> in <OR>")
}

func TestAmbig_LongPrefixNoAmbiguity(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRule_Alternation_1"].Decision
	src := newMockSource(tokenA, tokenA, tokenB)
	result := ar.s.AdaptivePredict(src, altDecision, EmptyPredicates)
	assert.Equal(t, 1, result, "should pick alt1 when B terminates")
	assert.Equal(t, 0, len(*reports), "no ambiguity expected")
}

func TestAmbig_PredicateOverride_Auto(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRuleWithPred_Alternation_1"].Decision
	src := newMockSource(tokenA)
	// No predicates → auto-resolve to first ambiguous alt
	result := ar.s.AdaptivePredict(src, altDecision, EmptyPredicates)
	assert.Equal(t, 0, result)
	assert.True(t, len(*reports) > 0)
	assert.Contains(t, (*reports)[0], "<0, 1> in <OR>")
}

func TestAmbig_PredicateOverride_True(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRuleWithPred_Alternation_1"].Decision
	src := newMockSource(tokenA)
	// pred=true: alt0 enabled, alt1 disabled
	preds := &PredicateSet{}
	preds.Set(0, true)
	preds.Set(1, false)
	result := ar.s.AdaptivePredict(src, altDecision, preds)
	assert.Equal(t, 0, result)
	assert.Equal(t, 0, len(*reports), "no ambiguity with predicate")
}

func TestAmbig_PredicateOverride_False(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRuleWithPred_Alternation_1"].Decision
	src := newMockSource(tokenA)
	// pred=false: alt0 disabled, alt1 enabled
	preds := &PredicateSet{}
	preds.Set(0, false)
	preds.Set(1, true)
	result := ar.s.AdaptivePredict(src, altDecision, preds)
	assert.Equal(t, 1, result)
	assert.Equal(t, 0, len(*reports), "no ambiguity with predicate")
}

func TestAmbig_NonAmbigInPredicated(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	altDecision := ar.s.atn.DecisionMap["AltRuleWithPred_Alternation_1"].Decision
	src := newMockSource(tokenB)
	result := ar.s.AdaptivePredict(src, altDecision, EmptyPredicates)
	assert.Equal(t, 2, result)
	assert.Equal(t, 0, len(*reports))
}

func TestAmbig_AltFollowedByOption(t *testing.T) {
	ar := buildAmbigRules()
	reports := &[]string{}
	ar.s.logging = func(m string) { *reports = append(*reports, m) }

	// AltWithOption: OR(A→2, B→4) + OPTION(A→1)
	// Input: B, A → OR picks B (alt1, value=4), OPTION picks A (value=1) → total=5
	altDecision := ar.s.atn.DecisionMap["AltWithOption_Alternation_1"].Decision
	optDecision := ar.s.atn.DecisionMap["AltWithOption_Option_1"].Decision

	src1 := newMockSource(tokenB, tokenA)
	altResult := ar.s.AdaptivePredict(src1, altDecision, EmptyPredicates)
	assert.Equal(t, 1, altResult, "B → alt1 (value offset 4 in TS test)")

	// advance source past the B
	src2 := newMockSource(tokenA)
	optResult := ar.s.AdaptivePredict(src2, optDecision, EmptyPredicates)
	assert.Equal(t, 0, optResult, "A follows → OPTION taken (alt 0)")

	// Total value: intermediate=4 (alt1+1 in TS), option=1 → 5.
	// We just check both decisions are correct; the arithmetic is done by the caller.
	_ = reports
}
