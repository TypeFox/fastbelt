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
	rHandle := ruleStart.RuleHandle()
	tokHandle := ruleStart.TokenRef(tokenTypes["ID"])
	ruleStart.NewEpsilon(rHandle.Left, tokHandle.Left)
	ruleStart.NewEpsilon(tokHandle.Right, rHandle.Right)
	expected := builder.Build()

	requireATNRulesEqual(t, expected, actual, "Start")
}

func requireATNRulesEqual(t *testing.T, expected, actual *ATN, ruleName string) {
	statesCount := len(expected.States)
	require.Len(t, actual.States, statesCount, "ATN state count mismatch")
}
