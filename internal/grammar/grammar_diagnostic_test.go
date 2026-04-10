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
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name=ID;
		Foo: Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError, "A rule's name has to be unique. 'Foo' is used multiple times.")
}

func TestDuplicateInterfaceNames(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		interface Foo { Other string }
	`+commonTokens).AssertDiagnostic(core.SeverityError, "An interface name has to be unique. 'Foo' is used multiple times.")
}

// --- Terminal ---

func TestTerminalMatchesEmptyString(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
		token EMPTY: /a*/;
		hidden token WS: /[ \n\r\t]+/;
	`).AssertDiagnostic(core.SeverityError, "This terminal could match an empty string.")
}

// --- Keywords ---

func TestKeywordEmpty(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: "" Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError, "Keywords cannot be empty.")
}

func TestKeywordWhitespaceOnly(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: " " Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError, "Keywords cannot only consist of whitespace characters.")
}

func TestKeywordContainsWhitespace(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: "hello world" Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityWarning, "Keywords should not contain whitespace characters.")
}

// --- Parser rule return type ---

func TestRuleWithoutReturnType(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		OrphanRule: Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError, "Unable to find return type for rule 'OrphanRule'.")
}

// --- Interface circular inheritance ---

func TestCircularInterfaceExtensionDirect(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo extends Foo {}
	`+commonTokens).AssertDiagnostic(core.SeverityError, "An interface cannot extend itself, neither directly nor indirectly.")
}

func TestCircularInterfaceExtensionIndirect(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface A extends B {}
		interface B extends A {}
	`+commonTokens).AssertDiagnostic(core.SeverityError, "An interface cannot extend itself, neither directly nor indirectly.")
}

// --- Unassigned rule call ---

func TestUnassignedRuleCallReturnTypeMismatch(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo {}
		interface Bar {}
		Foo: SubRule;
		SubRule returns Bar: ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError,
		"The return type 'Bar' of the called rule is not assignable to the return type 'Foo' of the current rule.")
}

func TestUnassignedRuleCallAfterAction(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	// Bar extends Foo, so {Bar.Items+=current} is type-valid; the only error is the position check.
	f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar extends Foo {}
		Bar: ({Bar.Items+=current} Bar);
	`+commonTokens).AssertDiagnostic(core.SeverityError,
		"An unassigned rule call cannot be preceded by an assigned action.")
}

func TestUnassignedRuleCallAfterAssignment(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name=ID SubRule;
		SubRule returns Foo: Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError,
		"An unassigned rule call cannot be preceded by an assignment.")
}

// --- Action type assignability ---

func TestActionTypeNotAssignableToRuleReturn(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	// Action type is Bar; rule Foo returns Foo. Bar does not extend Foo → type error.
	f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar { Items []Foo }
		Foo: ({Bar.Items+=current} ID);
	`+commonTokens).AssertDiagnostic(
		core.SeverityError,
		"The type 'Bar' of the action is not assignable to the rule's return type 'Foo'.",
	)
}

// --- Assignment operator mismatches ---

func TestBooleanOperatorOnNonBoolField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name?=ID;
	`+commonTokens).AssertDiagnostic(
		core.SeverityError,
		"The '?=' operator can only be used on boolean fields.",
	)
}

func TestArrayOperatorOnNonArrayField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name+=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError, "The '+=' operator can only be used on array fields.")
}

// --- Assignment value type compatibility ---

func TestCrossRefToNonReferenceField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Name string }
		Foo: Name=[Target:ID];
	`+commonTokens).AssertDiagnostic(core.SeverityError, "Cannot assign a cross-reference value to a non-reference field.")
}

func TestCrossRefTypeMismatch(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Baz { Name string }
		interface Foo { Ref *Bar }
		Foo: Ref=[Baz:ID];
	`+commonTokens).AssertDiagnostic(core.SeverityError,
		"The type 'Baz' of the cross-reference value is not assignable to the target field type 'Bar'.")
}

func TestTokenAssignedToNonStringField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Child { Name string }
		interface Foo { Child Child }
		Foo: Child=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError, "Cannot assign a token to a non-string field.")
}

func TestParserRuleReturnTypeMismatch(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Baz { Name string }
		interface Foo { Child Bar }
		Foo: Child=BazRule;
		Bar returns Bar: Name=ID;
		BazRule returns Baz: Name=ID;
	`+commonTokens).AssertDiagnostic(core.SeverityError,
		"The return type 'Baz' of the called rule is not assignable to the target field type 'Bar'.")
}

func TestKeywordAssignedToNonStringField(t *testing.T) {
	f := fbtest.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Child { Name string }
		interface Foo { Child Child }
		Foo: Child="keyword";
	`+commonTokens).AssertDiagnostic(core.SeverityError, "Cannot assign a keyword value to a non-string field.")
}
