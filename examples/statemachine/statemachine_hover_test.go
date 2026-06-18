// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
)

func TestDocumentationSingleLineComment(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
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
	f := test.New(t, CreateLspServices(nil))
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

func TestDocumentationBlankLinePreserved(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	docProvider := service.MustGet[server.DocumentationProvider](f.Services())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

// line one
//
// line two
state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors()

	off := test.MustFindNamedNode[State](doc, "off")
	assert.Equal(t, "line one  \n  \nline two", docProvider.Documentation(off))
}

func TestDocumentationBlockComment(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
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

func TestDocumentationIgnoresCommentsSeparatedByBlankLine(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
	docProvider := service.MustGet[server.DocumentationProvider](f.Services())

	doc := f.Parse(`
statemachine Test
events flick
initialState off

/* Section header */

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

func TestDocumentationNoComment(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))
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
	f := test.New(t, CreateLspServices(nil))

	f.Parse(`
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
`).AssertNoErrors().ExpectHoverAt("toOn", "The on state")
}

func TestHoverNilForUndocumentedReference(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))

	f.Parse(`
statemachine Test
events flick
initialState off

state off
  flick => <|toOn:on|>
end

state on
  flick => off
end
`).AssertNoErrors().ExpectNoHoverAt("toOn")
}

func TestHoverNilForWhitespace(t *testing.T) {
	f := test.New(t, CreateLspServices(nil))

	f.Parse(`
statemachine<|ws> Test
events flick
initialState off

state off
  flick => on
end

state on
  flick => off
end
`).AssertNoErrors().ExpectNoHoverAt("ws")
}
