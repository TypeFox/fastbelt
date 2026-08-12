// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

const commonTokens = `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
comment token SL_COMMENT: /\/\/[^\r\n]*/;
comment token ML_COMMENT: /\/\*[\s\S]*?\*\//;
`

func TestTextmateGenerator(t *testing.T) {
	f := test.New(t, CreateServices())
	doc := f.Parse(`
		grammar Test;
		interface Foo { Bar *Bar }
		interface Bar { Name string }
		Entry returns Foo: "foo" bar=SubRule;
		SubRule returns Bar: Name=ID;
	` + commonTokens)
	grammar := test.MustFindNode[grammar.Grammar](doc)
	content := GenerateTextMate(grammar, TextMateGeneratorConfig{
		Id:              "fastbelt",
		FileExtensions:  []string{".fb"},
		CaseInsensitive: false,
	})
	require.Equal(t, `{
  "repository": {
    "comments": {
      "patterns": [
        {
          "name": "comment.line.fastbelt",
          "begin": "//",
          "beginCaptures": {
            "1": {
              "name": "punctuation.whitespace.comment.leading.fastbelt"
            }
          },
          "end": "(?=$)"
        },
        {
          "name": "comment.block.fastbelt",
          "begin": "/\\*",
          "beginCaptures": {
            "0": {
              "name": "punctuation.definition.comment.fastbelt"
            }
          },
          "end": "\\*/",
          "endCaptures": {
            "0": {
              "name": "punctuation.definition.comment.fastbelt"
            }
          }
        }
      ]
    }
  },
  "scopeName": "source.fastbelt",
  "patterns": [
    {
      "include": "#comments"
    },
    {
      "name": "keyword.control.fastbelt",
      "match": "\\B(\"foo\")\\B"
    }
  ],
  "fileTypes": [
    ".fb"
  ],
  "name": "fastbelt"
}`, content)
}

func CreateServices() *service.Container {
	sc := service.NewContainer()
	SetupServices(sc)
	sc.Seal()
	return sc
}

func SetupServices(sc *service.Container) {
	service.Put[workspace.LanguageID](sc, "fastbelt")
	service.Put[workspace.FileExtensions](sc, []string{".fb"})

	textdoc.SetupDefaultServices(sc)
	linking.SetupDefaultServices(sc)
	workspace.SetupDefaultServices(sc)
	grammar.SetupGeneratedServices(sc)
}
