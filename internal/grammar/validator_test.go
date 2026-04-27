// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
)

// --- Rule and interface uniqueness ---

func TestDuplicateRuleNames(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		<|1:Foo|>: Name=ID;
		<|2:Foo|>: Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateUniqueRuleName)
}

func TestDuplicateInterfaceNames(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface <|1:Foo|> { Name string }
		interface <|2:Foo|> { Other string }
	` + commonTokens)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityError)
	diag1.WithCode(ValidateUniqueInterfaceName)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityError)
	diag2.WithCode(ValidateUniqueInterfaceName)
}

// --- Terminal ---

func TestTerminalMatchesEmptyString(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
		token <|EMPTY|>: /a*/;
		hidden token WS: /[ \n\r\t]+/;
	`)
	diag := doc.ExpectDiagnostic("EMPTY")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateEmptyToken)
}

// --- Keywords ---

func TestKeywordEmpty(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|empty:""|> Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("empty")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateEmptyKeyword)
}

func TestKeywordWhitespaceOnly(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|ws:" "|> Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("ws")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateWhitespaceOnlyKeyword)
}

func TestKeywordContainsWhitespace(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|keyword:"hello world"|> Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("keyword")
	diag.WithSeverity(core.SeverityWarning)
	diag.WithCode(ValidateKeywordWithWhitespace)
}

// --- Parser rule return type ---

func TestRuleWithoutReturnType(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		<|OrphanRule|>: Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("OrphanRule")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateRuleReturnType)
}

// --- Interface circular inheritance ---

func TestCircularInterfaceExtensionDirect(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo extends <|Foo|> {}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("Foo")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateInterfaceExtends)
}

func TestCircularInterfaceExtensionIndirect(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface A extends <|B|> {}
		interface B extends <|A|> {}
	` + commonTokens)
	diagB := doc.ExpectDiagnostic("B")
	diagB.WithSeverity(core.SeverityError)
	diagB.WithCode(ValidateInterfaceExtends)
	diagA := doc.ExpectDiagnostic("A")
	diagA.WithSeverity(core.SeverityError)
	diagA.WithCode(ValidateInterfaceExtends)
}

// --- Unassigned rule call ---

func TestUnassignedRuleCallReturnTypeMismatch(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {}
		interface Bar {}
		Foo: <|SubRule|>;
		SubRule returns Bar: ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("SubRule")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateRuleCallReturnType)
}

func TestUnassignedRuleCallAfterAction(t *testing.T) {
	f := test.New(t, CreateServices())
	// Bar extends Foo, so {Bar.Items+=current} is type-valid; the only error is the position check.
	doc := f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar extends Foo {}
		Bar: ({Bar.Items+=current} <|Bar|>);
	` + commonTokens)
	diag := doc.ExpectDiagnostic("Bar")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateRuleCallPosition)
}

func TestUnassignedRuleCallAfterAssignment(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name=ID <|SubRule|>;
		SubRule returns Foo: Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("SubRule")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateRuleCallPosition)
}

// --- Action type assignability ---

func TestActionTypeNotAssignableToRuleReturn(t *testing.T) {
	f := test.New(t, CreateServices())
	// Action type is Bar; rule Foo returns Foo. Bar does not extend Foo → type error.
	doc := f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar { Items []Foo }
		Foo: ({<|Bar|>.Items+=current} ID);
	` + commonTokens)
	diag := doc.ExpectDiagnostic("Bar")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateActionAssignmentType)
}

// --- Assignment operator mismatches ---

func TestBooleanOperatorOnNonBoolField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name<|?=|>ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("?=")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

func TestArrayOperatorOnNonArrayField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name<|+=|>ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("+=")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

// --- Assignment value type compatibility ---

func TestCrossRefToNonReferenceField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Name string }
		Foo: Name=<|tar:[Target:ID]|>;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("tar")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

func TestCrossRefTypeMismatch(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Baz { Name string }
		interface Foo { Ref *Bar }
		Foo: Ref=[<|Baz|>:ID];
	` + commonTokens)
	diag := doc.ExpectDiagnostic("Baz")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

func TestTokenAssignedToNonStringField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Child { Name string }
		interface Foo { Child Child }
		Foo: Child=<|ID|>;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("ID")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

func TestParserRuleReturnTypeMismatch(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Baz { Name string }
		interface Foo { Child Bar }
		Foo: Child=<|BazRule|>;
		Bar returns Bar: Name=ID;
		BazRule returns Baz: Name=ID;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("BazRule")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

func TestKeywordAssignedToNonStringField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Child { Name string }
		interface Foo { Child Child }
		Foo: Child=<|1:"keyword"|>;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}
