// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"testing"
)

func TestTokenRef(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, GrammarTemplate(`
		interface Start { Name string }
		Start: Name=ID;
	`))
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, [][]int{
		[]int{-1},
	})
}

func TestRuleRef(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, GrammarTemplate(`
		interface Start { Property Rule }
		interface Rule { Name string }
		Start: Property=Rule;
		Rule: Name=ID;
	`))
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID"}, [][]int{
		[]int{-1},
	})
}

func TestPlusRule(t *testing.T) {
	atn, rules, tokenTypes := FixtureATN(t, GrammarTemplate(`
		interface Start { Property []Rule }
		interface Rule { Name string }
		Start: Property+=Rule+;
		Rule: Name=ID;
	`))
	RequireATNRecognizes(t, atn, rules, tokenTypes, "Start", []string{"ID", "ID", "ID"}, [][]int{
		[]int{0, 0, 0},
		[]int{0, 0},
		[]int{0},
	})
}
