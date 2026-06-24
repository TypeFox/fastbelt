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
