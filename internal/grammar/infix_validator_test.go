// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
)

func TestInfixRuleValid(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		infix BinaryExpression on PrimaryExpression:
			"*" | "/"
			> "+" | "-"
			> right assoc "="
	` + commonTokens)
	doc.AssertNoDiagnostics()
}

func TestInfixRuleMissingNodeType(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		infix <|1:Binary|> on PrimaryExpression returns Expression: "+"
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixNodeType)
}

func TestInfixRuleMissingProperty(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Expression { Value string }
		interface Binary extends Expression {
			Left Expression
			Operator string
		}
		PrimaryExpression returns Expression: Value=ID
		infix <|1:Binary|> on PrimaryExpression returns Expression: "+"
	` + commonTokens)
	// Missing 'Right' field.
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixProperty)
}

func TestInfixRuleOperatorFieldNotString(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Expression { Value string }
		interface Binary extends Expression {
			Left Expression
			Operator Expression
			Right Expression
		}
		PrimaryExpression returns Expression: Value=ID
		infix <|1:Binary|> on PrimaryExpression returns Expression: "+"
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixProperty)
}

func TestInfixRuleOperandNotParserRule(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		infix BinaryExpression on <|1:ID|>: "+"
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixOperandRule)
}

func TestInfixRuleIncompatibleReturnType(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Expression { Value string }
		interface Other {}
		interface BinaryExpression extends Expression {
			Left Expression
			Operator string
			Right Expression
		}
		PrimaryExpression returns Expression: Value=ID
		infix <|1:BinaryExpression|> on PrimaryExpression returns Other: "+"
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateInfixReturnType)
}

func TestInfixRuleDuplicateOperator(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		infix BinaryExpression on PrimaryExpression:
			"+" | "*"
			> <|1:"+"|> | "-"
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixDuplicateOperator)
}

func TestInfixRuleDuplicateOperatorViaTokenGroup(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		token group MulOp { "*" "/" }
		infix BinaryExpression on PrimaryExpression:
			"*"
			> <|1:MulOp|>
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixDuplicateOperator)
}

func TestInfixRuleHiddenTokenOperator(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		hidden token COMMENT: /#[^\n]*/;
		infix BinaryExpression on PrimaryExpression:
			"+" | <|1:COMMENT|>
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixOperator)
}

func TestInfixRuleOperatorGroupNameTaken(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
	` + infixInterfaces + `
		token BinaryExpressionOperator: /@+/;
		infix <|1:BinaryExpression|> on PrimaryExpression: "+"
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateInfixOperatorGroupName)
}
