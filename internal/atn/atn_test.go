// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package atn

import (
	"testing"
)

func TestTokenRef(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Name string }
		Start: Name=ID;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID"}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
}

func TestRuleRef(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Property Rule }
		interface Rule { Name string }
		Start: Property=Rule;
		Rule: Name=ID;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
}

func TestPlusRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Property []Rule }
		interface Rule { Name string }
		Start: Property+=Rule+;
		Rule: Name=ID;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "ID", "ID", "ID", "ID"}, 1)
}

func TestStarRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Property []Rule }
		interface Rule { Name string }
		Start: Property+=Rule*;
		Rule: Name=ID;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "ID", "ID", "ID", "ID"}, 1)
}

func TestAlternativesRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Property string }
		Start: Property=(ID|NUMBER);
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"NUMBER"}, 1)
}

func TestNestedAlternativesRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Property AbstractRule }
		interface AbstractRule {}
		interface RuleId extends AbstractRule {
    		Id string
		}
		interface RuleNumber extends AbstractRule {
			Num string
		}

		Start: Property=(RuleId|RuleNumber);
		RuleId: Id=ID;
		RuleNumber: Num=NUMBER;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"NUMBER"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"NUMBER", "ID"}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "NUMBER"}, 0)
}

func TestNestedStarAlternativesRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Properties []AbstractRule }
		interface AbstractRule {}
		interface RuleId extends AbstractRule {
    		Id string
		}
		interface RuleNumber extends AbstractRule {
			Num string
		}

		Start: Properties+=(RuleId|RuleNumber)*;
		RuleId: Id=ID;
		RuleNumber: Num=NUMBER;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"NUMBER"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"NUMBER", "ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "NUMBER"}, 1)
}

func TestCrossRefsRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { 
			Decl RuleDeclaration
			Defs []RuleDefinition
		}
		interface AbstractRule {}
		interface RuleDeclaration extends AbstractRule {
    		Name string
		}
		interface RuleDefinition extends AbstractRule {
			IdRef *RuleDeclaration
			Value string
		}

		Start: Decl=RuleDeclaration Defs+=RuleDefinition*;
		RuleDeclaration: Name=ID;
		RuleDefinition: IdRef=[RuleDeclaration:ID] Value=NUMBER;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "NUMBER"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "NUMBER", "ID"}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "NUMBER", "NUMBER"}, 0)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "NUMBER", "ID", "NUMBER"}, 1)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "NUMBER", "ID", "NUMBER", "ID", "NUMBER"}, 1)
}

func TestThreePaths(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Property Rule }
		interface Rule {}
		interface Rule1 extends Rule { Id string }
		interface Rule2 extends Rule { Name string }
		interface Rule3 extends Rule { Ref *Rule2 }
		Start: Property=Rule;
		Rule: Rule1 | Rule2 | Rule3;
		Rule1: Id=ID;
		Rule2: Name=ID;
		Rule3: Ref=[Rule2:ID];
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 3)
}

func TestCompositeRuleRef(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Name string }
		composite SimpleComposite: ID;
		Start: Name=SimpleComposite;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
}

func TestCompositeRuleWithAlternativesRef(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, `
		interface Start { Name string }
		composite IDOrNumber: ID | NUMBER;
		Start: Name=IDOrNumber;
	`, false)
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, 1)
}
