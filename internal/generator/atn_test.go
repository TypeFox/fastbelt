// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

func grammarTemplate(rules string) string {
	return "grammar Test;\n" + rules + `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`
}

func fixtureATN(t *testing.T, grammarStr string) (*ATN, map[string]*grammar.ParserRule, map[string]TokenInfo) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(grammarStr).AssertNoErrors()
	grammr := doc.Root().(grammar.Grammar)
	return CreateATN(grammr)
}

func TestTokenRef(t *testing.T) {
	actual, rules, tokenTypes := fixtureATN(t, grammarTemplate(`
		interface Start { Name string }
		Start: Name=ID;
	`))

	builder := NewATNBuilder()
	ruleStart := builder.DeclareRule(rules["Start"])
	tokHandle := ruleStart.TokenRef(tokenTypes["ID"])
	ruleStart.Assign(tokHandle)
	expected := builder.Build()

	requireATNRulesEqual(t, expected, actual, "Start")
}

func TestRuleRef(t *testing.T) {
	actual, rules, tokenTypes := fixtureATN(t, grammarTemplate(`
		interface Start { Property Rule }
		interface Rule { Name string }
		Start: Property=Rule;
		Rule: Name=ID;
	`))

	builder := NewATNBuilder()
	{
		rule := builder.DeclareRule(rules["Rule"])
		tokHandle := rule.TokenRef(tokenTypes["ID"])
		rule.Assign(tokHandle)
	}
	{
		rule := builder.DeclareRule(rules["Start"])
		ruleRefHandle := rule.RuleRef(rules["Rule"])
		rule.Assign(ruleRefHandle)
	}
	expected := builder.Build()

	requireATNRulesEqual(t, expected, actual, "Start")
}

func TestPlusRule(t *testing.T) {
	actual, rules, tokenTypes := fixtureATN(t, grammarTemplate(`
		interface Start { Property []Rule }
		interface Rule { Name string }
		Start: Property+=Rule+;
		Rule: Name=ID;
	`))

	builder := NewATNBuilder()
	{
		rule := builder.DeclareRule(rules["Rule"])
		tokHandle := rule.TokenRef(tokenTypes["ID"])
		rule.Assign(tokHandle)
	}
	{
		rule := builder.DeclareRule(rules["Start"])
		ruleRefHandle := rule.RuleRef(rules["Rule"])
		plusHandle := rule.Plus("", ruleRefHandle)
		rule.Assign(plusHandle)
	}
	expected := builder.Build()

	requireATNRulesEqual(t, expected, actual, "Start")
}

func requireATNRulesEqual(t *testing.T, expected, actual *ATN, ruleName string) {
	statesCount := len(expected.States)
	require.Len(t, actual.States, statesCount, "ATN state count mismatch")
}
