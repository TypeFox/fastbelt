// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package multilang

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
)

func TestLanguageDispatch(t *testing.T) {
	sc := CreateServices()

	fixture := test.New(t, sc)
	docs := fixture.ParseAll(
		"file:///ws/a.hello",
		"hello world",
		"file:///ws/b.bye",
		// Should be able to reference the "world" value here
		"goodbye world",
	)

	require.Len(t, docs, 2)
	hello := docs[0].Document
	greeting, ok := hello.Root.(Greeting)
	require.True(t, ok, "hello.Root = %T, want Greeting", hello.Root)
	assert.Equal(t, "world", greeting.Name(), "Greeting name = %q, want %q", greeting.Name(), "world")

	bye := docs[1].Document
	farewell, ok := bye.Root.(Farewell)
	require.True(t, ok, "bye.Root = %T, want Farewell", bye.Root)
	refItem := farewell.To().Ref(t.Context())
	assert.Equal(t, greeting, refItem, "Farewell To = %v, want %v", refItem, greeting)
}

func TestLanguageAwareLexing(t *testing.T) {
	sc := CreateServices()

	fixture := test.New(t, sc)
	docs := fixture.ParseAll(
		"file:///ws/a.hello",
		"hello goodbye",
		"file:///ws/b.bye",
		"goodbye hello",
	)

	require.Len(t, docs, 2)

	// "goodbye" is only a keyword in the Farewell language; in a .hello
	// document it must lex as Token_ID and parse as a Greeting name.
	hello := docs[0].Document
	greeting, ok := hello.Root.(Greeting)
	require.True(t, ok, "hello.Root = %T, want Greeting", hello.Root)
	assert.Equal(t, "goodbye", greeting.Name(), "Greeting name = %q, want %q", greeting.Name(), "goodbye")
	for _, token := range hello.Tokens {
		if token.Type.Id == Keyword_goodbye_Idx {
			assert.Fail(t, "Keyword_goodbye lexed in a .hello document")
		}
	}

	// Symmetric: "hello" is not a keyword in the Farewell language.
	bye := docs[1].Document
	farewell, ok := bye.Root.(Farewell)
	require.True(t, ok, "bye.Root = %T, want Farewell", bye.Root)
	assert.Equal(t, "hello", farewell.To().Text(), "Farewell name = %q, want %q", farewell.To().Text(), "hello")
	for _, token := range bye.Tokens {
		if token.Type.Id == Keyword_hello_Idx {
			assert.Fail(t, "Keyword_hello lexed in a .bye document")
		}
	}
}

// completionAt returns the completion labels at the "cursor" marker of a
// document with the given URI, routed to its language by the selector.
func completionAt(t *testing.T, sc *service.Container, uri, src string) []string {
	t.Helper()
	items := test.New(t, sc).ParseURI(src, uri).CompletionItems("cursor")
	labels := make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.Label)
	}
	return labels
}

func TestLanguageAwareCompletion(t *testing.T) {
	sc := CreateServices()

	// At the entry of a .hello document only Greeting's keyword must surface.
	hello := completionAt(t, sc, "file:///ws/a.hello", "<|cursor>")
	if !slices.Contains(hello, "hello") {
		t.Errorf("expected 'hello' at .hello entry; got %v", hello)
	}
	if slices.Contains(hello, "goodbye") {
		t.Errorf("did not expect 'goodbye' at .hello entry; got %v", hello)
	}

	// Symmetric: only Farewell's keyword in a .bye document.
	bye := completionAt(t, sc, "file:///ws/b.bye", "<|cursor>")
	if !slices.Contains(bye, "goodbye") {
		t.Errorf("expected 'goodbye' at .bye entry; got %v", bye)
	}
	if slices.Contains(bye, "hello") {
		t.Errorf("did not expect 'hello' at .bye entry; got %v", bye)
	}
}

func TestCrossLanguageCompletion(t *testing.T) {
	sc := CreateServices()
	fixture := test.New(t, sc)
	docs := fixture.ParseAll(
		"file:///ws/a.hello",
		"hello world",
		"file:///ws/b.bye",
		// Should be able to complete the "world" value here
		"goodbye <|cursor>",
	)
	require.Len(t, docs, 2)
	farewell := docs[1]

	items := farewell.CompletionItems("cursor")
	require.Len(t, items, 1)
	assert.Equal(t, "world", items[0].Label)
}
