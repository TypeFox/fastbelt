// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

// compare_test.go — side-by-side LL(k) vs ALL(*) comparison.
//
// Each test builds the same grammar twice: once with a static LL(1) strategy
// (which fails due to ambiguity) and once with ALL(*) (which succeeds).
// The Strategy interface is chosen at parser-construction time.

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	parser "typefox.dev/fastbelt/parser"
)

// ────────────────────────────────────────────────────────────────────────────
// Token alphabet
// ────────────────────────────────────────────────────────────────────────────

const (
	cmpA = 1 // token 'a'
	cmpB = 2 // token 'b'
	cmpC = 3 // token 'c'
)

var tokenNames = map[string]int{"a": cmpA, "b": cmpB, "c": cmpC}

// tokenize converts a space-separated string such as "a a b" into a
// *parser.ParserState backed by the corresponding token sequence.
func tokenize(s string) *parser.ParserState {
	var tokens []*core.Token
	for _, word := range strings.Fields(s) {
		if id, ok := tokenNames[word]; ok {
			tokens = append(tokens, &core.Token{TypeId: id})
		}
	}
	return parser.NewParserState(tokens)
}

// ────────────────────────────────────────────────────────────────────────────
// Static LL(1) strategy
// ────────────────────────────────────────────────────────────────────────────

// llkStrategy is a static LL(1) prediction strategy backed by pre-built tables.
// altTables maps decision key → (first token type ID → alternative index).
// When a token has no entry — or is shared by multiple alts — -1 is returned.
type llkStrategy struct {
	altTables map[string]map[int]int
	optTables map[string][]int
}

func (s *llkStrategy) PredictAlternation(src TokenSource, key string) int {
	table, ok := s.altTables[key]
	if !ok {
		return -1
	}
	tID := tokenTypeID(src.LA(1))
	alt, ok := table[tID]
	if !ok {
		return -1
	}
	return alt
}

func (s *llkStrategy) PredictOptional(src TokenSource, key string) bool {
	triggers, ok := s.optTables[key]
	if !ok {
		return false
	}
	tID := tokenTypeID(src.LA(1))
	for _, id := range triggers {
		if id == tID {
			return true
		}
	}
	return false
}

// ────────────────────────────────────────────────────────────────────────────
// Mini parser driven by Strategy
// ────────────────────────────────────────────────────────────────────────────

// streamParser drives token consumption using parser.ParserState and resolves
// decisions through a pluggable Strategy.
type streamParser struct {
	state    *parser.ParserState
	strategy Strategy
}

func newStreamParser(state *parser.ParserState, strategy Strategy) *streamParser {
	return &streamParser{state: state, strategy: strategy}
}

// predict returns the chosen alternative for key, or an error if no prediction
// is possible.
func (p *streamParser) predict(key string) (int, error) {
	alt := p.strategy.PredictAlternation(p.state, key)
	if alt < 0 {
		return -1, fmt.Errorf("no LL prediction for %q", key)
	}
	return alt, nil
}

// consumeMany consumes zero or more tokens of the given type.
func (p *streamParser) consumeMany(typ int) {
	for p.state.LAId(1) == typ {
		p.state.Consume(typ)
	}
}

// consumeAtLeastOne consumes one or more tokens of the given type.
// Records a parser error if none are present.
func (p *streamParser) consumeAtLeastOne(typ int) {
	if p.state.LAId(1) != typ {
		// trigger error via Consume
		p.state.Consume(typ)
		return
	}
	p.consumeMany(typ)
}

// parseErrors returns the first recorded parser error, or nil.
func (p *streamParser) parseErrors() error {
	errs := p.state.Errors()
	if len(errs) > 0 {
		return errors.New(errs[0].Msg)
	}
	return nil
}

// ────────────────────────────────────────────────────────────────────────────
// Grammar builders
// ────────────────────────────────────────────────────────────────────────────

// buildLongGrammar returns rules for: LongRule := OR(A+, A+·B)
func buildLongGrammar() []*Rule {
	rule := &Rule{
		Name: "LongRule",
		Definition: []Production{
			&Alternation{
				Alternatives: []*Alternative{
					{Definition: []Production{
						&RepetitionMandatory{Definition: []Production{
							&Terminal{TokenTypeID: cmpA, Idx: 1},
						}, Idx: 1},
					}},
					{Definition: []Production{
						&RepetitionMandatory{Definition: []Production{
							&Terminal{TokenTypeID: cmpA, Idx: 2},
						}, Idx: 2},
						&Terminal{TokenTypeID: cmpB, Idx: 1},
					}},
				},
				Idx: 1,
			},
		},
	}
	return []*Rule{rule}
}

// buildSubruleGrammar returns rules for: AltRule := OR(Ra, Rb); Ra := A+; Rb := A+·B
func buildSubruleGrammar() []*Rule {
	ra := &Rule{Name: "Ra", Definition: []Production{
		&RepetitionMandatory{Definition: []Production{
			&Terminal{TokenTypeID: cmpA, Idx: 1},
		}, Idx: 1},
	}}
	rb := &Rule{Name: "Rb", Definition: []Production{
		&RepetitionMandatory{Definition: []Production{
			&Terminal{TokenTypeID: cmpA, Idx: 1},
		}, Idx: 1},
		&Terminal{TokenTypeID: cmpB, Idx: 1},
	}}
	altRule := &Rule{Name: "AltRule", Definition: []Production{
		&Alternation{
			Alternatives: []*Alternative{
				{Definition: []Production{&NonTerminal{ReferencedRule: ra, Idx: 1}}},
				{Definition: []Production{&NonTerminal{ReferencedRule: rb, Idx: 2}}},
			},
			Idx: 1,
		},
	}}
	return []*Rule{ra, rb, altRule}
}

// buildManyGrammar returns rules for: ManyRule := OR(A*, A*·C)
func buildManyGrammar() []*Rule {
	rule := &Rule{
		Name: "ManyRule",
		Definition: []Production{
			&Alternation{
				Alternatives: []*Alternative{
					{Definition: []Production{
						&Repetition{Definition: []Production{
							&Terminal{TokenTypeID: cmpA, Idx: 1},
						}, Idx: 1},
					}},
					{Definition: []Production{
						&Repetition{Definition: []Production{
							&Terminal{TokenTypeID: cmpA, Idx: 2},
						}, Idx: 2},
						&Terminal{TokenTypeID: cmpC, Idx: 1},
					}},
				},
				Idx: 1,
			},
		},
	}
	return []*Rule{rule}
}

// buildOptionGrammar returns rules for: OptionRule := OPTION(A+) · A+
func buildOptionGrammar() []*Rule {
	rule := &Rule{
		Name: "OptionRule",
		Definition: []Production{
			&Option{
				Definition: []Production{
					&RepetitionMandatory{Definition: []Production{
						&Terminal{TokenTypeID: cmpA, Idx: 1},
					}, Idx: 1},
				},
				Idx: 1,
			},
			&RepetitionMandatory{Definition: []Production{
				&Terminal{TokenTypeID: cmpA, Idx: 2},
			}, Idx: 2},
		},
	}
	return []*Rule{rule}
}

// ────────────────────────────────────────────────────────────────────────────
// parse helpers (drive actual token consumption)
// ────────────────────────────────────────────────────────────────────────────

func (p *streamParser) parseLong() error {
	alt, err := p.predict("LongRule_Alternation_1")
	if err != nil {
		return err
	}
	p.consumeAtLeastOne(cmpA)
	if alt == 1 {
		p.state.Consume(cmpB)
	}
	return p.parseErrors()
}

func (p *streamParser) parseSubrule() error {
	alt, err := p.predict("AltRule_Alternation_1")
	if err != nil {
		return err
	}
	p.consumeAtLeastOne(cmpA)
	if alt == 1 {
		p.state.Consume(cmpB)
	}
	return p.parseErrors()
}

func (p *streamParser) parseMany() error {
	alt, err := p.predict("ManyRule_Alternation_1")
	if err != nil {
		return err
	}
	p.consumeMany(cmpA)
	if alt == 1 {
		p.state.Consume(cmpC)
	}
	return p.parseErrors()
}

func (p *streamParser) parseOption() error {
	if p.strategy.PredictOptional(p.state, "OptionRule_Option_1") {
		p.consumeAtLeastOne(cmpA)
	}
	p.consumeAtLeastOne(cmpA)
	return p.parseErrors()
}

// ────────────────────────────────────────────────────────────────────────────
// Scenario 1 — Greedy alternation: OR(A+, A+·B)
// ────────────────────────────────────────────────────────────────────────────

func TestGreedyOR_LLk_fails_ALLStar_succeeds(t *testing.T) {
	cases := []struct {
		input   string
		wantAlt int
	}{
		{"a a b", 1},
		{"a a a", 0},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			grammar := buildLongGrammar()
			// LL(1): first token 'a' maps to both alts → conflict, returns -1
			llk := &llkStrategy{altTables: map[string]map[int]int{
				"LongRule_Alternation_1": {}, // empty: no unambiguous entry for 'a'
			}}
			// act + assert LL(k)
			errLLk := newStreamParser(tokenize(tc.input), llk).parseLong()
			assert.Error(t, errLLk, "LL(k) must fail: %s", tc.input)

			// ALL(*): simulates both paths
			star := NewLLStarLookahead(grammar, nil)
			errStar := newStreamParser(tokenize(tc.input), star).parseLong()
			assert.NoError(t, errStar, "ALL(*) must succeed: %s", tc.input)
		})
	}
}

// ────────────────────────────────────────────────────────────────────────────
// Scenario 2 — Sub-rule indirection: OR(SUBRULE(Ra), SUBRULE(Rb))
// ────────────────────────────────────────────────────────────────────────────

func TestSubruleOR_LLk_fails_ALLStar_succeeds(t *testing.T) {
	cases := []struct {
		input   string
		wantAlt int
	}{
		{"a a b", 1},
		{"a a a", 0},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			grammar := buildSubruleGrammar()
			llk := &llkStrategy{altTables: map[string]map[int]int{
				"AltRule_Alternation_1": {}, // 'a' maps to both Ra and Rb → conflict
			}}

			errLLk := newStreamParser(tokenize(tc.input), llk).parseSubrule()
			assert.Error(t, errLLk, "LL(k) must fail: %s", tc.input)

			star := NewLLStarLookahead(grammar, nil)
			errStar := newStreamParser(tokenize(tc.input), star).parseSubrule()
			assert.NoError(t, errStar, "ALL(*) must succeed: %s", tc.input)
		})
	}
}

// ────────────────────────────────────────────────────────────────────────────
// Scenario 3 — Zero-or-more common prefix: OR(A*, A*·C)
// ────────────────────────────────────────────────────────────────────────────

func TestManyOR_LLk_fails_ALLStar_succeeds(t *testing.T) {
	cases := []struct {
		input string
	}{
		{"a a c"},
		{"a a"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			grammar := buildManyGrammar()
			llk := &llkStrategy{altTables: map[string]map[int]int{
				"ManyRule_Alternation_1": {}, // 'a' shared by both alts → conflict
			}}

			errLLk := newStreamParser(tokenize(tc.input), llk).parseMany()
			assert.Error(t, errLLk, "LL(k) must fail: %s", tc.input)

			star := NewLLStarLookahead(grammar, nil)
			errStar := newStreamParser(tokenize(tc.input), star).parseMany()
			assert.NoError(t, errStar, "ALL(*) must succeed: %s", tc.input)
		})
	}
}

// ────────────────────────────────────────────────────────────────────────────
// Scenario 4 — Option greediness: OPTION(A+) · A+
//
// Input "a": only one 'a' available. Taking the option (LL-1 greedy) consumes it,
// leaving the mandatory A+ with nothing → error.  ALL(*) correctly skips the
// option and lets the mandatory A+ consume the single token.
// ────────────────────────────────────────────────────────────────────────────

func TestOptionConflict_LLk_greedy_wrong_ALLStar_correct(t *testing.T) {
	grammar := buildOptionGrammar()
	// LL(1) greedy: sees 'a', takes option (standard LL table entry)
	llk := &llkStrategy{optTables: map[string][]int{
		"OptionRule_Option_1": {cmpA}, // sees 'a' → take option
	}}

	// act
	errLLk := newStreamParser(tokenize("a"), llk).parseOption()
	errStar := newStreamParser(tokenize("a"), NewLLStarLookahead(grammar, nil)).parseOption()

	// assert
	assert.Error(t, errLLk, "LL(1) greedy takes option, leaving outer A+ with no input")
	assert.NoError(t, errStar, "ALL(*) skips option, outer A+ consumes 'a'")
}

