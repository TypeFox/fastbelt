// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package token_groups

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

func assertIdentifier(t *testing.T, doc *test.Doc, expected string) {
	doc.AssertNoErrors()
	model := doc.Document.Root.(Model)
	id := model.Item().ValueToken()
	if id.Image != expected {
		t.Fatalf("Expected Item.ValueToken text to be '%s', but got '%s'", expected, id.Image)
	}
}

func TestIdentifierParsesId(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("a JohnDoe")
	assertIdentifier(t, doc, "JohnDoe")
}

func TestIdentifierParsesInt(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("a 123")
	assertIdentifier(t, doc, "123")
}

func TestIdentifierParsesKeyword(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("a a")
	assertIdentifier(t, doc, "a")
}

func TestIdentifierInAlternative(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("b JohnDoe")
	assertIdentifier(t, doc, "JohnDoe")
}

func TestIdentifierInAlternativeWithKeyword(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("b b")
	assertIdentifier(t, doc, "b")
}

func TestNestedIdentifierA(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("c a")
	assertIdentifier(t, doc, "a")
}

func TestNestedIdentifierKeyword(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("c nested")
	assertIdentifier(t, doc, "nested")
}

func TestNestedIdentifierID(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("c JohnDoe")
	assertIdentifier(t, doc, "JohnDoe")
}

func TestNestedIdentifierInt(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("c 123")
	assertIdentifier(t, doc, "123")
}

func TestOptionalIdentifierMissing(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("d")
	doc.AssertNoErrors()
	model := doc.Document.Root.(Model)
	id := model.Item().ValueToken()
	assert.Nil(t, id, "Expected Item.ValueToken to be nil when optional identifier is missing")
}

func TestOptionalIdentifierPresent(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("d JohnDoe")
	assertIdentifier(t, doc, "JohnDoe")
}

func TestRecoveryMatchesBothIdentifiers(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// Happy path for rule E: both group assignments must bind in sequence.
	doc := fixture.Parse("e foo nested")
	doc.AssertNoErrors()
	model := doc.Document.Root.(Model)
	recovery := model.Item().(Recovery)
	first := recovery.FirstToken()
	second := recovery.SecondToken()
	require.NotNil(t, first)
	require.NotNil(t, second)
	assert.Equal(t, "foo", first.Image)
	assert.Equal(t, "nested", second.Image)
}

func TestIdentifierRejectsNonMemberKeyword(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// "b" is a keyword, not a member of the Identifier group {ID, INT, "a"},
	// so it must not be accepted where an Identifier is expected.
	doc := fixture.Parse("a b")
	require.NotEmpty(t, doc.Document.ParserErrors)
}

func TestKeywordGroupMatchesFirst(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("f one")
	assertIdentifier(t, doc, "one")
}

func TestKeywordGroupMatchesSecond(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("f two")
	assertIdentifier(t, doc, "two")
}

func TestKeywordGroupRejectsNonMember(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// "three" lexes as an ID, which is not a member of KeywordGroup {"one", "two"}.
	doc := fixture.Parse("f three")
	require.NotEmpty(t, doc.Document.ParserErrors)
}

func TestRecoverySkipMissingIdentifier(t *testing.T) {
	fixture := test.New(t, CreateServices())
	// The grammar expects an identifier after 'e'
	// the recovery should correctly skip one token and assign the "nested" keyword to "Second"
	doc := fixture.Parse("e nested")
	require.Len(t, doc.Document.ParserErrors, 1)
	model := doc.Document.Root.(Model)
	recovery := model.Item().(Recovery)
	second := recovery.SecondToken()
	require.NotNil(t, second)
	assert.Equal(t, "nested", second.Image)
}

func TestRegexpGroupMatchesOne(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("g one")
	assertIdentifier(t, doc, "one")
}

func TestRegexpGroupMatchesTwo(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("g two")
	assertIdentifier(t, doc, "two")
}

func TestRegexpGroupRejectsNonMember(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("g three")
	require.NotEmpty(t, doc.Document.ParserErrors)
}
