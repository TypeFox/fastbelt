// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
)

// All .fb files of a folder form one grammar. These tests cover the
// grammar-level validations that must therefore look at the whole folder while
// reporting on each file's own nodes.

func TestFolderDuplicateRuleNames(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface Foo { Name string }
			<|dup:Foo|>: Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar Test;
			<|dup:Foo|>: Name=ID;
		`,
	)
	for _, doc := range docs {
		doc.ExpectDiagnostic("dup").
			WithSeverity(core.SeverityError).
			WithCode(ValidateUniqueRuleName)
	}
}

func TestFolderDuplicateInterfaceNames(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface <|dup:Foo|> { Name string }
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar Test;
			interface <|dup:Foo|> { Other string }
		`,
	)
	for _, doc := range docs {
		doc.ExpectDiagnostic("dup").
			WithSeverity(core.SeverityError).
			WithCode(ValidateUniqueInterfaceName)
	}
}

func TestFolderGrammarNameMismatch(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar <|name:Alpha|>;
			interface Foo { Name string }
			Foo: Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar <|name:Beta|>;
			interface Bar { Name string }
			Bar: Name=ID;
		`,
	)
	for _, doc := range docs {
		doc.ExpectDiagnostic("name").
			WithSeverity(core.SeverityError).
			WithCode(ValidateGrammarNameMismatch).
			WithMessage("All grammar files in a folder form one grammar and must declare the same name, but found: 'Alpha' and 'Beta'.")
	}
	// A file in another folder is a different grammar and does not interfere.
	other := f.ParseURI(`
		grammar Gamma;
		interface Baz { Name string }
		Baz: Name=ID;
	`+commonTokens, "file:///elsewhere/c.fb")
	other.AssertNoErrors()
}

func TestFolderImplicitReturnTypeFromSibling(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/types.fb", `
			grammar Test;
			interface Person { Name string Age string }
		`,
		"file:///ws/rules.fb", `
			grammar Test;
			Person: "person" <|nameRef:Name|>=ID <|ageRef:Age|>=ID;
		`+commonTokens,
	)
	typesDoc, rulesDoc := docs[0], docs[1]
	rulesDoc.AssertNoErrors()
	typesDoc.AssertNoErrors()

	rule := test.MustFindNamedNode[ParserRule](rulesDoc, "Person")
	iface := test.MustFindNamedNode[Interface](typesDoc, "Person")
	assert.Same(t, iface, FindReturnType(rule, rulesDoc.Ctx()))

	// Field assignments of the rule scope against the sibling's interface.
	nameRef := test.MustFindReference[Field](rulesDoc, "nameRef")
	require.Nil(t, nameRef.Error())
	assert.Same(t, iface.Fields()[0], nameRef.Ref(rulesDoc.Ctx()))
	ageRef := test.MustFindReference[Field](rulesDoc, "ageRef")
	require.Nil(t, ageRef.Error())
	assert.Same(t, iface.Fields()[1], ageRef.Ref(rulesDoc.Ctx()))
}

func TestFolderTwoDefaultTokenModes(t *testing.T) {
	f := test.New(t, CreateServices())
	// a.fb declares nothing that is exported (an unnamed mode with a nested
	// token) and must still take part in the folder-wide check.
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			token mode <|dup:default|> {
				token NUM: /[0-9]+/
			}
		`,
		"file:///ws/b.fb", `
			grammar Test;
			interface Foo { Value string }
			Foo: Value=NUM;
			token mode <|dup:default|> {
				NUM
			}
		`,
	)
	for _, doc := range docs {
		doc.ExpectDiagnostic("dup").
			WithSeverity(core.SeverityError).
			WithCode(ValidateUniqueTokenModeName)
	}
}

func TestFolderDefaultTokenModeInSibling(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface Foo { Name string }
			Foo: Name=ID INNER;
			token ID: /[a-z]+/
			hidden token WS: /\s+/
			token mode default {
				ID -> push(Inner)
				hidden WS
			}
		`,
		"file:///ws/b.fb", `
			grammar Test;
			token mode Inner {
				token INNER: /[A-Z]+/ -> pop
			}
		`,
	)
	for _, doc := range docs {
		doc.AssertNoDiagnostics()
	}
}

func TestFolderTokenUsedOnlyBySibling(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/tokens.fb", `
			grammar Test;
			token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
			token NUM: /[0-9]+/;
			hidden token WS: /[ \n\r\t]+/;
		`,
		"file:///ws/rules.fb", `
			grammar Test;
			interface Foo { Name string Value string }
			Foo: Name=ID Value=NUM;
		`,
	)
	for _, doc := range docs {
		doc.AssertNoDiagnostics()
	}
}

func TestFolderKeywordCoveredBySiblingMode(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/modes.fb", `
			grammar Test;
			token ID: /[a-z]+/
			hidden token WS: /\s+/
			token mode default {
				"person"
				ID
				hidden WS
			}
		`,
		"file:///ws/rules.fb", `
			grammar Test;
			interface Person { Name string }
			Person: "person" Name=ID;
		`,
	)
	for _, doc := range docs {
		doc.AssertNoDiagnostics()
	}
}

func TestFolderKeywordInlineAndAsTokenDecl(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface Foo { Name string }
			Foo: <|kw:"person"|> Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar Test;
			token PERSON: <|kw:"person"|>;
		`,
	)
	for _, doc := range docs {
		doc.ExpectDiagnostic("kw").
			WithSeverity(core.SeverityError).
			WithCode(ValidateKeywordPureStandaloneOrTokenDecl)
	}
}

func TestFolderViewIsSharedAcrossSiblings(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface Foo { Name string }
			Foo: Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar Test;
			interface Bar { Name string }
			Bar: Name=ID;
		`,
	)
	a := docs[0].Root().(Grammar)
	b := docs[1].Root().(Grammar)
	// One view per folder: both documents resolve to the same aggregate.
	assert.Same(t, viewOf(a), viewOf(b))
	assert.Same(t, folderGrammar(a), folderGrammar(b))
	assert.Equal(t, []Grammar{a, b}, siblingGrammars(b))
	assert.Same(t, test.MustFindNamedNode[Interface](docs[1], "Bar"), FindInterfaceByName(a, "Bar"))
}

func TestFolderGrammarNameMismatchListsAllNames(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar <|name:Gamma|>;
			interface Foo { Name string }
			Foo: Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar <|name:Alpha|>;
			interface Bar { Name string }
			Bar: Name=ID;
		`,
		"file:///ws/c.fb", `
			grammar <|name:Beta|>;
			interface Baz { Name string }
			Baz: Name=ID;
		`,
		// Same name as a.fb: no new entry in the list.
		"file:///ws/d.fb", `
			grammar <|name:Gamma|>;
			interface Qux { Name string }
			Qux: Name=ID;
		`,
	)
	// Every file reports the full, sorted, de-duplicated set of names.
	for _, doc := range docs {
		doc.ExpectDiagnostic("name").
			WithSeverity(core.SeverityError).
			WithCode(ValidateGrammarNameMismatch).
			WithMessageContaining("'Alpha', 'Beta' and 'Gamma'")
	}
}

func TestFolderGrammarNamesMatchNoDiagnostic(t *testing.T) {
	f := test.New(t, CreateServices())
	docs := f.ParseAll(
		"file:///ws/a.fb", `
			grammar Test;
			interface Foo { Name string }
			Foo: Name=ID;
		`+commonTokens,
		"file:///ws/b.fb", `
			grammar Test;
			interface Bar { Name string }
			Bar: Name=ID;
		`,
	)
	for _, doc := range docs {
		doc.AssertNoErrors()
	}
}
