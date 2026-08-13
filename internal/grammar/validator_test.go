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
	for _, label := range []string{"1", "2"} {
		diag := doc.ExpectDiagnostic(label)
		diag.WithSeverity(core.SeverityError)
		diag.WithCode(ValidateUniqueRuleName)
	}
}

func TestDuplicateRuleNamesDifferentTypes(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		<|1:Foo|>: Name=ID;
		token <|2:Foo|>: ID;
	` + commonTokens)
	for _, label := range []string{"1", "2"} {
		diag := doc.ExpectDiagnostic(label)
		diag.WithSeverity(core.SeverityError)
		diag.WithCode(ValidateUniqueRuleName)
	}
}

func TestDuplicateRuleNamesTokenGroup(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		token group <|1:Foo|> { ID };
		token <|2:Foo|>: ID;
	` + commonTokens)
	for _, label := range []string{"1", "2"} {
		diag := doc.ExpectDiagnostic(label)
		diag.WithSeverity(core.SeverityError)
		diag.WithCode(ValidateUniqueRuleName)
	}
}

func TestDuplicateInterfaceNames(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface <|1:Foo|> { Name string }
		interface <|2:Foo|> { Other string }
	` + commonTokens)
	for _, label := range []string{"1", "2"} {
		diag := doc.ExpectDiagnostic(label)
		diag.WithSeverity(core.SeverityError)
		diag.WithCode(ValidateUniqueInterfaceName)
	}
}

// --- Interface field names uniqueness, capitalization ---

func TestFieldNameUppercaseStart(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Name string
		}
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestFieldNameLowercaseStart(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:name|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateFieldNameCapitalLetter)
}

func TestDuplicateFieldNames(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Name string
			<|1:Name|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateUniqueFieldName)
}

func TestDuplicateFieldNamesCaseInsensitive(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Name string
			<|1:NAME|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateUniqueFieldName)
}

func TestDuplicateFieldNamesCaseInsensitiveAndCapitalLetter(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Name string
			<|1:name|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithCode(ValidateUniqueFieldName)
	diag.WithSeverity(core.SeverityError)

	diag = doc.ExpectDiagnostic("1")
	diag.WithCode(ValidateFieldNameCapitalLetter)
	diag.WithSeverity(core.SeverityError)
}

func TestReservedFieldNameDocument(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:Document|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameTokens(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:Tokens|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameSetPrefix(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:SetName|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameIsPrefix(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:IsActive|> bool
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameTokenArrayAllowed(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Token []string
		}
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestReservedFieldNameText(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:Text|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameForEachNode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:ForEachNode|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameSameAsInterface(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			<|1:Foo|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameSuffixTokenConflict(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Name string
			<|1:NameToken|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestReservedFieldNameSuffixNodeConflictInherited(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base {
			Target composite
		}
		interface Derived extends Base {
			<|1:TargetNode|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateReservedFieldName)
}

func TestNestedArrayType(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Items <|1:[][]string|>
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateNestedArrayType)
}

func TestNestedArrayTypeWithInterface(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Bar { Name string }
		interface Foo {
			Items <|1:[][]Bar|>
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateNestedArrayType)
}

func TestSingleArrayTypeAllowed(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo {
			Items []string
			Refs []*Foo
		}
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestInheritedFieldNameValid(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base {
			Name string
		}
		interface Derived extends Base {
			Other string
		}
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestInheritedFieldNameDuplicate(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base {
			Name string
		}
		interface Derived extends Base {
			<|1:Name|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateUniqueFieldName)
}

func TestInheritedFieldNameDuplicateCaseInsensitive(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Base {
			Name string
		}
		interface Derived extends Base {
			<|1:NAME|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateUniqueFieldName)
}

func TestInheritedFieldNameDuplicateDeepHierarchy(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface C {
			Foo string
		}
		interface B extends C {}
		interface A extends B {
			<|1:Foo|> string
		}
	` + commonTokens)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateUniqueFieldName)
}

func TestInheritedFieldNameDuplicateInDeepHierarchyNoErrorAtTopElement(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface C {
			Foo string
		}
		interface B extends C {
			Foo string
		}
		interface A extends B {
			<|1:Bar|> string
		}
	` + commonTokens)
	doc.AssertNoDiagnostic("1")
}

// --- Terminal ---

func TestTerminalMatchesEmptyString(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
		token <|EMPTY|>: /a*/;
		token <|EMPTY2|>: "";
		hidden token WS: /[ \n\r\t]+/;
	`)
	diag := doc.ExpectDiagnostic("EMPTY")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateEmptyToken)
	diag = doc.ExpectDiagnostic("EMPTY2")
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
	for _, label := range []string{"A", "B"} {
		diag := doc.ExpectDiagnostic(label)
		diag.WithSeverity(core.SeverityError)
		diag.WithCode(ValidateInterfaceExtends)
	}
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

func TestUnassignedRuleCallToParserRuleAfterAction(t *testing.T) {
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

func TestUnassignedRuleCallToNonParserRuleAfterAction(t *testing.T) {
	f := test.New(t, CreateServices())
	// Bar extends Foo, so {Bar.Items+=current} is type-valid; the only error is the position check.
	doc := f.Parse(`
		grammar Test;
		interface Foo { Items []Foo }
		interface Bar extends Foo {}
		Bar: ({Bar.Items+=current} <|Car|>);
		token Car: "car";
	` + commonTokens)
	doc.AssertNoDiagnostic("Car")
}

func TestUnassignedRuleCallToParserRuleAfterAssignment(t *testing.T) {
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

func TestUnassignedRuleCallToNonParserRuleAfterAssignment(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		Foo: Name=ID <|SubRule|>;
		token SubRule: "subrule";
	` + commonTokens)
	doc.AssertNoDiagnostic("SubRule")
}

// --- Action type assignability ---

func TestActionTypeNotAssignableToRuleReturn(t *testing.T) {
	f := test.New(t, CreateServices())
	// Action type is Bar; rule Foo returns Foo. Bar does not extend Foo -> type error.
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

func TestCompositeRuleAssignedToStringField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		composite QualifiedName: ID ("." ID)*;
		Foo: Name=<|QualifiedName|>;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("QualifiedName")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

func TestCompositeRuleAssignedToCompositeField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name composite }
		composite QualifiedName: ID ("." ID)*;
		Foo: Name=QualifiedName;
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestCompositeRuleAssignedToNonCompositeField(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Flag bool }
		composite QualifiedName: ID ("." ID)*;
		Foo: Flag=<|QualifiedName|>;
	` + commonTokens)
	diag := doc.ExpectDiagnostic("QualifiedName")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateAssignmentType)
}

// --- Token groups ---

func TestTokenGroupRecursiveDirect(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		token group <|X|> { X }
	` + commonTokens)
	diag := doc.ExpectDiagnostic("X")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateRecursiveTokenGroup)
}

func TestTokenGroupRecursiveIndirect(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		token group <|X|> { Y }
		token group <|Y|> { X }
	` + commonTokens)
	for _, label := range []string{"X", "Y"} {
		diag := doc.ExpectDiagnostic(label)
		diag.WithSeverity(core.SeverityError)
		diag.WithCode(ValidateRecursiveTokenGroup)
	}
}

// Negative test - validation does not trigger on standalone token group
func TestTokenGroupRecursiveNegative(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
	    grammar Test;
		token group <|X|> { ID }
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestTokenGroupWithInvalidToken(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		token group X { <|WS|> }
	` + commonTokens)
	diag := doc.ExpectDiagnostic("WS")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateInvalidTokenInGroup)
}

// --- Cross-references ---

func TestCrossRefWithHiddenToken(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Ref *Target }
		Foo: Ref=[Target:<|WS|>];
	` + commonTokens)
	diag := doc.ExpectDiagnostic("WS")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateInvalidTokenInCrossRef)
}

func TestCrossRefWithCommentToken(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Ref *Target }
		comment token SL_COMMENT: /\/\/[^\r\n]*/;
		Foo: Ref=[Target:<|SL_COMMENT|>];
	` + commonTokens)
	diag := doc.ExpectDiagnostic("SL_COMMENT")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateInvalidTokenInCrossRef)
}

func TestCrossRefWithValidToken(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Ref *Target }
		Foo: Ref=[Target:ID];
	` + commonTokens)
	doc.AssertNoErrors()
}

func TestCrossRefMissingTerminal(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Target { Name string }
		interface Foo { Ref *Target }
		Foo: Ref=[<|Target|>];
	` + commonTokens)
	diag := doc.ExpectDiagnostic("Target")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateMissingCrossRefTerminal)
}

func TestDuplicateTokenModeNames(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting="hello";

		token mode <|1:default|> { "hello" }
		token mode <|2:default|> { }
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityError)
	diag1.WithCode(ValidateUniqueTokenModeName)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityError)
	diag2.WithCode(ValidateUniqueTokenModeName)
}

func TestDuplicateTokenNamesAcrossDifferentTokenScopes(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token <|1:ID|>: /[a-z_][a-z0-9_]*/
		hidden token WS: /[ \n\r\t]+/

		token mode default {
			token <|2:ID|>: /[A-Z_][A-Z0-9_]*/
			token group <|3:ID|> {}
		}
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityError)
	diag1.WithCode(ValidateUniqueRuleName)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityError)
	diag2.WithCode(ValidateUniqueRuleName)
	diag3 := doc.ExpectDiagnostic("3")
	diag3.WithSeverity(core.SeverityError)
	diag3.WithCode(ValidateUniqueRuleName)
}

func TestNonTokenModeRequiresDefaultTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[A-Z_][A-Z0-9_]*/

		token mode <|1:FirstMode|> {
			ID
		}
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityError)
	diag1.WithCode(ValidateDefaultTokenModeRequired)
}

// --- Token mode commands ---

func TestPushCommandWithoutTargetMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> <|1:push|>
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenCommandMode)
	diag.WithMessageContaining("requires a target token mode")
}

func TestModeCommandWithoutTargetMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> <|1:mode|>
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenCommandMode)
}

func TestPushCommandWithTargetModeIsValid(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Other)
			hidden WS
		}

		token mode Other {
			ID -> pop
		}
	`).AssertNoDiagnostics()
}

func TestPushCommandTargetingDefaultModeIsValid(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Other)
			hidden WS
		}

		token mode Other {
			ID -> mode(default)
		}
	`).AssertNoDiagnostics()
}

func TestPushCommandWithUnknownTargetMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(<|1:Missing|>)
			hidden WS
		}
	`)
	// The linker reports the dangling mode reference; the command itself has a
	// target and must not additionally be flagged as incomplete.
	doc.ExpectDiagnostic("1").WithSeverity(core.SeverityError)
	for _, diag := range doc.Diagnostics() {
		if diag.Code == ValidateTokenCommandMode {
			t.Errorf("unexpected %s diagnostic: %s", ValidateTokenCommandMode, diag.Message)
		}
	}
}

func TestPopCommandWithTargetMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Other)
			hidden WS
		}

		token mode Other {
			ID -> pop(<|1:default|>)
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenCommandMode)
	diag.WithMessageContaining("cannot take a target mode")
}

func TestPopCommandWithTargetModeReference(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Other)
			hidden WS
		}

		token mode Other {
			ID -> pop(<|1:Other|>)
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenCommandMode)
}

// --- Token mode reachability and emptiness ---

func TestUnreachableTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}

		token mode <|1:Orphan|> {
			ID
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityWarning)
	diag.WithCode(ValidateUnreachableTokenMode)
}

func TestUnreachableTokenModeNotReportedForPopTarget(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}

		token mode <|1:Orphan|> {
			ID -> pop
		}
	`)
	// A 'pop' inside a mode does not make that mode reachable.
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityWarning)
	diag.WithCode(ValidateUnreachableTokenMode)
}

func TestDefaultTokenModeIsAlwaysReachable(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`).AssertNoDiagnostics()
}

func TestUnreachableTokenModeReachedFromNestedMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Middle)
			hidden WS
		}

		token mode <|1:Middle|> {
			ID -> push(Inner)
		}

		token mode Inner {
			ID -> pop
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityWarning)
	diag.WithCode(ValidateNonDefaultTokenModeNoPop)
}

func TestEmptyTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode <|1:default|> { }
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityWarning)
	diag.WithCode(ValidateEmptyTokenMode)
}

func TestEmptyNamedTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Other)
			hidden WS
		}

		token mode <|1:Other|> { }
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityWarning)
	diag.WithCode(ValidateEmptyTokenMode)
	diag.WithMessageContaining("'Other'")
}

// --- Token mode coverage of parser tokens ---

func TestKeywordNotInAnyTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID <|1:"world"|>;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateKeywordNotInTokenMode)
}

func TestKeywordListedInTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			"world"
			hidden WS
		}
	`).AssertNoDiagnostics()
}

func TestKeywordCoveredByKeywordSelector(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			keywords /^[a-z]+$/
			ID
			hidden WS
		}
	`).AssertNoDiagnostics()
}

func TestKeywordNotMatchedByKeywordSelector(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID <|1:"!"|>;

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			keywords /^[a-z]+$/
			ID
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateKeywordNotInTokenMode)
}

func TestKeywordCoveredByKeywordBackedToken(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token WORLD: "world"
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			WORLD
			hidden WS
		}
	`).AssertNoErrors()
}

func TestKeywordCoveredByTokenGroupInMode(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token group Greetings {
			"world"
		}
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			Greetings
			hidden WS
		}
	`).AssertNoErrors()
}

func TestKeywordCoveredByNestedTokenGroupInMode(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token group Inner {
			"world"
		}
		token group Outer {
			Inner
		}
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			Outer
			hidden WS
		}
	`).AssertNoErrors()
}

func TestKeywordCoveredByModeLocalTokenDeclaration(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			token WORLD: "world"
			hidden WS
		}
	`).AssertNoErrors()
}

func TestKeywordCoveredByNonDefaultTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID -> push(Other)
			hidden WS
		}

		token mode Other {
			"world" -> pop
		}
	`).AssertNoDiagnostics()
}

func TestKeywordInTokenDeclarationNotReportedAsUncovered(t *testing.T) {
	f := test.New(t, CreateServices())
	// GREETING is never used by a parser rule, so its keyword does not have to
	// be reachable through a token mode.
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token GREETING: "hello"
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`).AssertNoErrors()
}

func TestKeywordUncoveredReportedOnlyOnce(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world" "world" "world";

		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`)
	count := 0
	for _, diag := range doc.Diagnostics() {
		if diag.Code == ValidateKeywordNotInTokenMode {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 %s diagnostic, got %d", ValidateKeywordNotInTokenMode, count)
	}
}

func TestKeywordCoverageNotCheckedWithoutTokenModes(t *testing.T) {
	f := test.New(t, CreateServices())
	// Without token modes the generator registers every keyword automatically.
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID "world";
	` + commonTokens).AssertNoDiagnostics()
}

func TestTokenNotInAnyTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string Name string }
		Foo: Greeting=NAME Name=ID;

		token <|1:NAME|>: /[A-Z]+/
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenNotInTokenMode)
	diag.WithMessageContaining("'NAME'")
}

func TestTokenGroupNotInAnyTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=Greetings;

		token group <|1:Greetings|> {
			ID
		}
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`)
	// The group's members are registered, but the group has its own token id
	// and is therefore still unreachable.
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenNotInTokenMode)
}

func TestKeywordBackedTokenCoveredByKeywordSelector(t *testing.T) {
	f := test.New(t, CreateServices())
	// WORLD shares the token id of the keyword "world", so a selector that
	// covers the keyword covers the token too.
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=WORLD;

		token WORLD: "world"
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			keywords /^[a-z]+$/
			ID
			hidden WS
		}
	`).AssertNoErrors()
}

func TestTokenInCrossRefNotInAnyTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name string }
		interface Bar { Ref *Foo }
		Foo: Name=ID;
		Bar: Ref=[Foo:NAME];

		token <|1:NAME|>: /[A-Z]+/
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenNotInTokenMode)
}

func TestTokenInCompositeRuleNotInAnyTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Name composite }
		Foo: Name=FQN;
		composite FQN: ID (DOT ID)*;

		token <|1:DOT|>: /\./
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			ID
			hidden WS
		}
	`)
	diag := doc.ExpectDiagnostic("1")
	diag.WithSeverity(core.SeverityError)
	diag.WithCode(ValidateTokenNotInTokenMode)
}

func TestTokenOnlyUsedInsideTokenGroupNotReported(t *testing.T) {
	f := test.New(t, CreateServices())
	// NAME is referenced from a token group, not from a parser rule, so it does
	// not need its own entry in a token mode.
	f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=Greetings;

		token group Greetings {
			ID
			NAME
		}
		token NAME: /[A-Z]+/
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			Greetings
			hidden WS
		}
	`).AssertNoDiagnostics()
}

func TestRecursiveTokenGroupCoverageTerminates(t *testing.T) {
	f := test.New(t, CreateServices())
	// The recursion is reported elsewhere; coverage collection must not hang.
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=Outer;

		token group Outer {
			Inner
		}
		token group Inner {
			Outer
			ID
		}
		token ID: /[a-z]+/
		hidden token WS: /\s+/

		token mode default {
			Outer
			hidden WS
		}
	`)
	for _, diag := range doc.Diagnostics() {
		if diag.Code == ValidateTokenNotInTokenMode {
			t.Errorf("unexpected %s diagnostic: %s", ValidateTokenNotInTokenMode, diag.Message)
		}
	}
}

func TestUnreachableTokenModeCycle(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/

		token mode <|0:default|> {
			ID
		}

		token mode <|1:LEFT|> {
			ID -> mode(RIGHT)
        }
		token mode <|2:RIGHT|> {
			ID -> mode(LEFT)
        }
	`)
	doc.AssertNoDiagnostic("0")
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityWarning)
	diag1.WithCode(ValidateUnreachableTokenMode)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityWarning)
	diag2.WithCode(ValidateUnreachableTokenMode)
}

func TestNonDefaultTokenModeWithNoExitCommand(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/

		token mode <|0:default|> {
			ID -> push(LEFT)
		}

		token mode <|1:LEFT|> {
			ID
			"123"
        }
	`)
	doc.AssertNoDiagnostic("0")
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityWarning)
	diag1.WithCode(ValidateNonDefaultTokenModeNoPop)
}

func TestTokenNotUsedInTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token <|1:ID|>: /[a-z]+/

		token mode default { "hello" }
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityError)
	diag1.WithCode(ValidateTokenNotInTokenMode)
}

func TestMembersUniqueInTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting=ID;

		token ID: /[a-z]+/

		token mode default {
			ID
			<|1:ID|>
			
			"hello"
			<|2:"hello"|>

			token X: /lol/
			token <|3:X|>: /lol2/

			token group G {
				ID
			}
			token group <|4:G|> {
				ID
			}
		}
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityError)
	diag1.WithCode(ValidateUniqueRuleNameInTokenMode)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityError)
	diag2.WithCode(ValidateUniqueRuleNameInTokenMode)
	diag3 := doc.ExpectDiagnostic("3")
	diag3.WithSeverity(core.SeverityError)
	diag3.WithCode(ValidateUniqueRuleName)
	diag4 := doc.ExpectDiagnostic("4")
	diag4.WithSeverity(core.SeverityError)
	diag4.WithCode(ValidateUniqueRuleName)
}

func TestNotCoveredByParserRuleUsingTokenModes(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting="x";

		token NUM: /[0-9]+/

		token mode default {
			token <|1:ID|>: /[a-z]+/
			<|2:NUM|>
			token group <|3:Cardinality|> {
				"?"
				"*"
				"+"
			}
			<|4:"unused"|>
		}
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityWarning)
	diag1.WithCode(ValidateTerminalNotCoveredByParserRule)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityWarning)
	diag2.WithCode(ValidateTerminalNotCoveredByParserRule)
	diag3 := doc.ExpectDiagnostic("3")
	diag3.WithSeverity(core.SeverityWarning)
	diag3.WithCode(ValidateTokenGroupNotCoveredByParserRule)
	diag4 := doc.ExpectDiagnostic("4")
	diag4.WithSeverity(core.SeverityWarning)
	diag4.WithCode(ValidateTerminalNotCoveredByParserRule)
}

func TestNotCoveredByParserRuleUsingImplicitTokenMode(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting="x";

		token <|1:ID|>: /[a-z]+/
		token group <|2:Cardinality|> {
			"?"
			"*"
			"+"
		}
	`)
	diag1 := doc.ExpectDiagnostic("1")
	diag1.WithSeverity(core.SeverityWarning)
	diag1.WithCode(ValidateTerminalNotCoveredByParserRule)
	diag2 := doc.ExpectDiagnostic("2")
	diag2.WithSeverity(core.SeverityWarning)
	diag2.WithCode(ValidateTokenGroupNotCoveredByParserRule)
}

func TestInvisibleTokenDontNeedToBeCoveredByParserRule(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Greeting string }
		Foo: Greeting="x";

		hidden token WS: /\s+/
		comment token group Comment {
			"//comment"
			"/*comment*/"
		}
	`)
	doc.AssertNoDiagnostics()
}
