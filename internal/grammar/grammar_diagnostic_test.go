// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	core "typefox.dev/fastbelt"
	fbtest "typefox.dev/fastbelt/testing"
)

// --- Rule and interface uniqueness ---

func TestDuplicateRuleNames(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		<|1:Foo|>: Name=ID;
		<|2:Foo|>: Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateUniqueRuleName)
}

func TestDuplicateInterfaceNames(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface <|1:Foo|> { Name string }
		interface <|2:Foo|> { Other string }
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateUniqueInterfaceName)
	doc.ExpectDiagnostic("2").WithSeverity(core.SeverityError).WithCode(ValidateUniqueInterfaceName)
}

// --- Terminal ---

func TestTerminalMatchesEmptyString(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
		token <|EMPTY|>: /a*/;
		hidden token WS: /[ \n\r\t]+/;
	`)
	doc.ExpectDiagnostic("EMPTY").WithSeverity(core.SeverityError).WithCode(ValidateEmptyToken)
}

// --- Keywords ---

func TestKeywordEmpty(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|empty:""|> Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("empty").WithSeverity(core.SeverityError).WithCode(ValidateEmptyKeyword)
}

func TestKeywordWhitespaceOnly(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|ws:" "|> Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("ws").WithSeverity(core.SeverityError).WithCode(ValidateWhitespaceOnlyKeyword)
}

func TestKeywordContainsWhitespace(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: <|keyword:"hello world"|> Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("keyword").WithSeverity(core.SeverityWarning).WithCode(ValidateKeywordWithWhitespace)
}

// --- Parser rule return type ---

func TestRuleWithoutReturnType(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		<|OrphanRule|>: Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("OrphanRule").WithSeverity(core.SeverityError).WithCode(ValidateRuleReturnType)
}

// --- Interface circular inheritance ---

func TestCircularInterfaceExtensionDirect(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo extends <|Foo|> {}
	` + commonTokens)
	doc.ExpectDiagnostic("Foo").WithSeverity(core.SeverityError).WithCode(ValidateInterfaceExtends)
}

func TestCircularInterfaceExtensionIndirect(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface A extends <|B|> {}
		interface B extends <|A|> {}
	` + commonTokens)
	doc.ExpectDiagnostic("B").WithSeverity(core.SeverityError).WithCode(ValidateInterfaceExtends)
	doc.ExpectDiagnostic("A").WithSeverity(core.SeverityError).WithCode(ValidateInterfaceExtends)
}

// --- Unassigned rule call ---

func TestUnassignedRuleCallReturnTypeMismatch(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {}
		interface Bar {}
		Foo: <|SubRule|>;
		SubRule returns Bar: ID;
	` + commonTokens)
	doc.ExpectDiagnostic("SubRule").WithSeverity(core.SeverityError).WithCode(ValidateRuleCallReturnType)
}

func TestUnassignedRuleCallAfterAction(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	// Bar extends Foo, so {Bar.Items+=current} is type-valid; the only error is the position check.
	doc := f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar extends Foo {}
		Bar: ({Bar.Items+=current} <|Bar|>);
	` + commonTokens)
	doc.ExpectDiagnostic("Bar").WithSeverity(core.SeverityError).WithCode(ValidateRuleCallPosition)
}

func TestUnassignedRuleCallAfterAssignment(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name=ID <|SubRule|>;
		SubRule returns Foo: Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("SubRule").WithSeverity(core.SeverityError).WithCode(ValidateRuleCallPosition)
}

// --- Action type assignability ---

func TestActionTypeNotAssignableToRuleReturn(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	// Action type is Bar; rule Foo returns Foo. Bar does not extend Foo → type error.
	doc := f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar { Items []Foo }
		Foo: ({<|Bar|>.Items+=current} ID);
	` + commonTokens)
	doc.ExpectDiagnostic("Bar").WithSeverity(core.SeverityError).WithCode(ValidateActionAssignmentType)
}

// --- Assignment operator mismatches ---

func TestBooleanOperatorOnNonBoolField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name<|?=|>ID;
	` + commonTokens)
	doc.ExpectDiagnostic("?=").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}

func TestArrayOperatorOnNonArrayField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name<|+=|>ID;
	` + commonTokens)
	doc.ExpectDiagnostic("+=").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}

// --- Assignment value type compatibility ---

func TestCrossRefToNonReferenceField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Name string }
		Foo: Name=<|tar:[Target:ID]|>;
	` + commonTokens)
	doc.ExpectDiagnostic("tar").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}

func TestCrossRefTypeMismatch(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Baz { Name string }
		interface Foo { Ref *Bar }
		Foo: Ref=[<|Baz|>:ID];
	` + commonTokens)
	doc.ExpectDiagnostic("Baz").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}

func TestTokenAssignedToNonStringField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Child { Name string }
		interface Foo { Child Child }
		Foo: Child=<|ID|>;
	` + commonTokens)
	doc.ExpectDiagnostic("ID").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}

func TestParserRuleReturnTypeMismatch(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Baz { Name string }
		interface Foo { Child Bar }
		Foo: Child=<|BazRule|>;
		Bar returns Bar: Name=ID;
		BazRule returns Baz: Name=ID;
	` + commonTokens)
	doc.ExpectDiagnostic("BazRule").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}

func TestKeywordAssignedToNonStringField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Child { Name string }
		interface Foo { Child Child }
		Foo: Child=<|1:"keyword"|>;
	` + commonTokens)
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError).WithCode(ValidateAssignmentType)
}
