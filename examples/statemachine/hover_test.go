// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// createHoverServices builds a service container with both the statemachine language
// and the default server services (DocumentationProvider, HoverProvider, etc.).
func createHoverServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	server.SetupDefaultServices(sc)
	sc.Seal()
	return sc
}

// hoverAt calls the hover provider at the position of the named marker.
// Accepts both index markers (<|label>) and range markers (<|label:text|>);
// for range markers the start of the range is used as the cursor position.
func hoverAt(t *testing.T, f *test.Fixture, doc *test.Doc, label string) *lsp.Hover {
	t.Helper()
	offset := -1
	for _, idx := range doc.Indices {
		if idx.Label == label {
			offset = idx.Offset
			break
		}
	}
	if offset == -1 {
		for _, r := range doc.Ranges {
			if r.Label == label {
				offset = r.Start
				break
			}
		}
	}
	if offset == -1 {
		t.Fatalf("no marker with label %q", label)
	}
	pos := doc.Document.TextDoc.PositionAt(offset)
	provider := service.MustGet[server.HoverProvider](f.Services())
	result, err := provider.HandleHoverRequest(f.Ctx(), &lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
			Position:     pos,
		},
	})
	require.NoError(t, err)
	return result
}

func TestDocumentationSingleLineComment(t *testing.T) {
	f := test.New(t, createHoverServices())
	docProvider := service.MustGet[server.DocumentationProvider](f.Services())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

// The off state
state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	off := test.MustFindNamedNode[State](doc, "off")
	assert.Equal(t, "The off state", docProvider.Documentation(off))
}

func TestDocumentationMultipleLineComments(t *testing.T) {
	f := test.New(t, createHoverServices())
	docProvider := service.MustGet[server.DocumentationProvider](f.Services())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

// line one
// line two
state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	off := test.MustFindNamedNode[State](doc, "off")
	assert.Equal(t, "line one  \nline two", docProvider.Documentation(off))
}

func TestDocumentationBlockComment(t *testing.T) {
	f := test.New(t, createHoverServices())
	docProvider := service.MustGet[server.DocumentationProvider](f.Services())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

/*
 * line one
 * line two
 */
state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	off := test.MustFindNamedNode[State](doc, "off")
	assert.Equal(t, "line one\nline two", docProvider.Documentation(off))
}

func TestDocumentationNoComment(t *testing.T) {
	f := test.New(t, createHoverServices())
	docProvider := service.MustGet[server.DocumentationProvider](f.Services())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	off := test.MustFindNamedNode[State](doc, "off")
	assert.Equal(t, "", docProvider.Documentation(off))
}

func TestHoverReturnsDocumentationForReference(t *testing.T) {
	f := test.New(t, createHoverServices())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

state off
  flick => <|toOn:on|>
end

// The on state
state on
  flick => off
end
`).AssertNoErrors()

	result := hoverAt(t, f, doc, "toOn")
	require.NotNil(t, result)
	assert.Equal(t, lsp.Markdown, result.Contents.Kind)
	assert.Equal(t, "The on state", result.Contents.Value)

	expectedRange, err := doc.MarkerRange("toOn")
	require.NoError(t, err)
	assert.Equal(t, expectedRange.LspRange(), result.Range)
}

func TestHoverNilForUndocumentedReference(t *testing.T) {
	f := test.New(t, createHoverServices())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

state off
  flick => <|toOn:on|>
end

state on
  flick => off
end
`).AssertNoErrors()

	assert.Nil(t, hoverAt(t, f, doc, "toOn"))
}

func TestHoverNilForWhitespace(t *testing.T) {
	f := test.New(t, createHoverServices())

	doc := f.Parse(`
statemachine<|ws> Test
events flick
initialState off

state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	assert.Nil(t, hoverAt(t, f, doc, "ws"))
}
