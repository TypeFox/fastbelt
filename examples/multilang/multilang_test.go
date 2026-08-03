// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package multilang

import (
	"context"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// parse builds a single document through the multilang container and returns
// its root node, routed to an entry rule by the registered LanguageSelector.
func parse(t *testing.T, sc *service.Container, path, text string) *core.Document {
	t.Helper()
	// A deliberately non-matching language id: the document is not in the
	// textdoc.Store, so the selector must route it by the URI-path glob.
	file := textdoc.NewFile(lsp.URIFromPath(path), "unrelated", 0, text)
	doc := core.NewDocument(file)
	service.MustGet[workspace.DocumentManager](sc).Set(doc)
	if err := service.MustGet[workspace.Builder](sc).Build(context.Background(), []*core.Document{doc}, nil); err != nil {
		t.Fatalf("Build: %v", err)
	}
	if len(doc.ParserErrors) != 0 {
		t.Fatalf("unexpected parser errors for %s: %v", path, doc.ParserErrors)
	}
	return doc
}

func TestLanguageDispatch(t *testing.T) {
	sc := CreateServices()

	hello := parse(t, sc, "/ws/a.hello", "hello world")
	if _, ok := hello.Root.(Greeting); !ok {
		t.Fatalf(".hello routed to %T, want Greeting", hello.Root)
	}

	bye := parse(t, sc, "/ws/b.bye", "goodbye world")
	if _, ok := bye.Root.(Farewell); !ok {
		t.Fatalf(".bye routed to %T, want Farewell", bye.Root)
	}
}

func TestLanguageAwareLexing(t *testing.T) {
	sc := CreateServices()

	// "goodbye" is only a keyword in the Farewell language; in a .hello
	// document it must lex as Token_ID and parse as a Greeting name.
	hello := parse(t, sc, "/ws/c.hello", "hello goodbye")
	greeting, ok := hello.Root.(Greeting)
	if !ok {
		t.Fatalf(".hello routed to %T, want Greeting", hello.Root)
	}
	if greeting.Name() != "goodbye" {
		t.Fatalf("Greeting name = %q, want %q", greeting.Name(), "goodbye")
	}
	for _, token := range hello.Tokens {
		if token.TypeId == Keyword_goodbye_Idx {
			t.Fatalf("Keyword_goodbye lexed in a .hello document")
		}
	}

	// Symmetric: "hello" is not a keyword in the Farewell language.
	bye := parse(t, sc, "/ws/d.bye", "goodbye hello")
	farewell, ok := bye.Root.(Farewell)
	if !ok {
		t.Fatalf(".bye routed to %T, want Farewell", bye.Root)
	}
	if farewell.Name() != "hello" {
		t.Fatalf("Farewell name = %q, want %q", farewell.Name(), "hello")
	}
	for _, token := range bye.Tokens {
		if token.TypeId == Keyword_hello_Idx {
			t.Fatalf("Keyword_hello lexed in a .bye document")
		}
	}
}
