// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package allstar implements the ALL(*) adaptive lookahead algorithm.
// See: https://www.antlr.org/papers/allstar-techreport.pdf
package allstar

// ProductionKind is a discriminator for the Production union.
type ProductionKind int

const (
	ProdTerminal            ProductionKind = iota
	ProdNonTerminal         ProductionKind = iota
	ProdAlternative         ProductionKind = iota // one branch inside an Alternation; no occurrence index
	ProdAlternation         ProductionKind = iota // OR(...)
	ProdOption              ProductionKind = iota // OPTION(...)
	ProdRepetition          ProductionKind = iota // MANY(...)
	ProdRepetitionMandatory ProductionKind = iota // AT_LEAST_ONE(...)
)

// Production is the interface for all grammar elements.
type Production interface {
	Kind() ProductionKind
	// Children returns the sub-productions, or nil for leaves.
	Children() []Production
	// Occurrence returns the 1-based index distinguishing multiple uses of the
	// same production kind within one rule (0 for Alternative which is not indexed).
	Occurrence() int
}

// Rule is a named parser rule.
type Rule struct {
	Name       string
	Definition []Production
}

// Terminal consumes a single token.
type Terminal struct {
	TokenTypeID int
	TokenName   string
	Idx         int
	// CategoryMatches holds token type IDs matched via category inheritance.
	CategoryMatches []int
}

func (t *Terminal) Kind() ProductionKind   { return ProdTerminal }
func (t *Terminal) Children() []Production { return nil }
func (t *Terminal) Occurrence() int        { return t.Idx }

// NonTerminal calls another rule.
type NonTerminal struct {
	ReferencedRule *Rule
	Idx            int
}

func (n *NonTerminal) Kind() ProductionKind   { return ProdNonTerminal }
func (n *NonTerminal) Children() []Production { return nil }
func (n *NonTerminal) Occurrence() int        { return n.Idx }

// Alternative is one branch inside an Alternation (no own occurrence index).
type Alternative struct {
	Definition []Production
}

func (a *Alternative) Kind() ProductionKind   { return ProdAlternative }
func (a *Alternative) Children() []Production { return a.Definition }
func (a *Alternative) Occurrence() int        { return 0 }

// Alternation is an OR decision point.
type Alternation struct {
	Alternatives []*Alternative
	Idx          int
}

func (a *Alternation) Kind() ProductionKind { return ProdAlternation }
func (a *Alternation) Children() []Production {
	out := make([]Production, len(a.Alternatives))
	for i, alt := range a.Alternatives {
		out[i] = alt
	}
	return out
}
func (a *Alternation) Occurrence() int { return a.Idx }

// Option wraps an optional sequence.
type Option struct {
	Definition []Production
	Idx        int
}

func (o *Option) Kind() ProductionKind   { return ProdOption }
func (o *Option) Children() []Production { return o.Definition }
func (o *Option) Occurrence() int        { return o.Idx }

// Repetition is a zero-or-more loop.
type Repetition struct {
	Definition []Production
	Idx        int
}

func (r *Repetition) Kind() ProductionKind   { return ProdRepetition }
func (r *Repetition) Children() []Production { return r.Definition }
func (r *Repetition) Occurrence() int        { return r.Idx }

// RepetitionMandatory is a one-or-more loop.
type RepetitionMandatory struct {
	Definition []Production
	Idx        int
}

func (r *RepetitionMandatory) Kind() ProductionKind   { return ProdRepetitionMandatory }
func (r *RepetitionMandatory) Children() []Production { return r.Definition }
func (r *RepetitionMandatory) Occurrence() int        { return r.Idx }

// ProductionTypeName returns the ATN decision-map key segment for p, and
// whether p is a decision-producing kind at all.
// Returns ("", false) for Terminal, NonTerminal, and Alternative.
func ProductionTypeName(p Production) (string, bool) {
	switch p.Kind() {
	case ProdAlternation:
		return "Alternation", true
	case ProdOption:
		return "Option", true
	case ProdRepetition:
		return "Repetition", true
	case ProdRepetitionMandatory:
		return "RepetitionMandatory", true
	default:
		return "", false
	}
}
