package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

func TestOptionalSemicolonBetweenParserRules(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		interface Bar { Name string }
		A: Name=ID
		B: Name=ID
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.Rules(), 2)
	assert.Equal(t, "A", g.Rules()[0].Name())
	assert.Equal(t, "B", g.Rules()[1].Name())
}

func TestOptionalSemicolonBetweenCompositeRules(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		composite A: ID
		composite B: ID
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.Composites(), 2)
	assert.Equal(t, "A", g.Composites()[0].Name())
	assert.Equal(t, "B", g.Composites()[1].Name())
}

func TestOptionalSemicolonCompositeBeforeParserRule(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		composite A: ID
		B: Name=ID
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.Composites(), 1)
	require.Len(t, g.Rules(), 1)
	assert.Equal(t, "B", g.Rules()[0].Name())
}

func TestOptionalSemicolonAfterStarBeforeReturnsRule(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Expression {}
		interface BinaryExpression extends Expression {
			Left Expression
			Operator string
			Right Expression
		}
		Addition returns Expression:
			Multiplication ({BinaryExpression.Left=current} Operator=("+" | "-") Right=Multiplication)*
		Multiplication returns Expression:
			PrimaryExpression
		PrimaryExpression returns Expression:
			Number=ID
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.Rules(), 3)
	assert.Equal(t, "Addition", g.Rules()[0].Name())
	assert.Equal(t, "Multiplication", g.Rules()[1].Name())
}

const infixInterfaces = `
	interface Expression { Value string }
	interface BinaryExpression extends Expression {
		Left Expression
		Operator string
		Right Expression
	}
	PrimaryExpression returns Expression: Value=ID
`

func TestInfixRule(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		infix BinaryExpression on PrimaryExpression:
			"%"
			> "^"
			> "*" | "/"
			> "+" | "-"
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.InfixRules(), 1)
	rule := g.InfixRules()[0]
	assert.Equal(t, "BinaryExpression", rule.Name())
	assert.Equal(t, "PrimaryExpression", rule.Call().Rule().Ref(doc.Ctx()).Name())
	assert.Nil(t, rule.ReturnType())
	require.Len(t, rule.Groups(), 4)
	require.Len(t, rule.Groups()[2].Operators(), 2)
	assert.Equal(t, "", rule.Groups()[2].Associativity())
	assert.Equal(t, `"*"`, rule.Groups()[2].Operators()[0].(Keyword).Value())
	assert.Equal(t, `"/"`, rule.Groups()[2].Operators()[1].(Keyword).Value())
}

func TestInfixRuleReturnsAndAssociativity(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		infix BinaryExpression on PrimaryExpression returns Expression:
			"*" | "/"
			> "+" | "-"
			> right "="
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.InfixRules(), 1)
	rule := g.InfixRules()[0]
	assert.Equal(t, "Expression", rule.ReturnType().Ref(doc.Ctx()).Name())
	require.Len(t, rule.Groups(), 3)
	assert.Equal(t, "", rule.Groups()[0].Associativity())
	assert.Equal(t, "right", rule.Groups()[2].Associativity())
	require.Len(t, rule.Groups()[2].Operators(), 1)
	assert.Equal(t, `"="`, rule.Groups()[2].Operators()[0].(Keyword).Value())
}

func TestInfixRuleTokenOperators(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		token group MulOp { "*" "/" }
		token POW: /\^+/;
		infix BinaryExpression on PrimaryExpression:
			POW
			> MulOp
			> left "+" | "-"
	` + commonTokens)
	doc.AssertNoParseErrors()
	doc.AssertNoLinkingErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.InfixRules(), 1)
	rule := g.InfixRules()[0]
	require.Len(t, rule.Groups(), 3)
	pow := rule.Groups()[0].Operators()[0].(RuleCall)
	assert.Equal(t, "POW", pow.Rule().Ref(doc.Ctx()).Name())
	mul := rule.Groups()[1].Operators()[0].(RuleCall)
	assert.Equal(t, "MulOp", mul.Rule().Ref(doc.Ctx()).Name())
	assert.Equal(t, "left", rule.Groups()[2].Associativity())
	require.Len(t, rule.Groups()[2].Operators(), 2)
}

func TestOptionalSemicolonStarRuleBeforePlainRule(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		interface Bar { Names []string }
		A: Names+=ID*
		B: Name=ID
	` + commonTokens)
	doc.AssertNoParseErrors()
	g := doc.Document.Root.(Grammar)
	require.Len(t, g.Rules(), 2)
	assert.Equal(t, "B", g.Rules()[1].Name())
}
