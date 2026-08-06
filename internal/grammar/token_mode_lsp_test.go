// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// tokenModeMarking replaces the fixture's `>` index-marker terminator with `@>`,
// so that the `>` of a `->` command inside a marked range is not mistaken for
// the end of a marker. Markers still must not be nested.
var tokenModeMarking = &test.TestMarking{
	StartRange: "<|",
	EndRange:   "|>",
	StartIndex: "",
	EndIndex:   "@>",
	Delimiter:  ":",
}

// lspFixture returns a fixture with the grammar language services, the generated
// server services (needed for completion) and the default LSP providers.
func lspFixture(t *testing.T) *test.Fixture {
	sc := service.NewContainer()
	SetupServices(sc)
	SetupGeneratedServerServices(sc)
	server.SetupDefaultServices(sc)
	sc.Seal()
	return test.New(t, sc).WithMarking(tokenModeMarking)
}

func completionLabels(items []lsp.CompletionItem) []string {
	labels := make([]string, len(items))
	for i, item := range items {
		labels[i] = item.Label
	}
	return labels
}

// --- Go to definition ---

func TestDefinitionOnTokenModeInPushCommand(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(<|Inner@>Inner)
	hidden WS
}
<|Inner:token mode Inner {
	ID -> pop
}|>`, "file:///definition.fb")
	doc.AssertNoErrors()
	// The definition of a mode reference is the mode declaration.
	doc.AssertDefinition("Inner")
}

func TestDefinitionOnTokenModeInSetModeCommand(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> mode(<|Inner@>Inner)
	hidden WS
}
<|Inner:token mode Inner {
	ID -> pop
}|>`, "file:///definition-set.fb")
	doc.AssertNoErrors()
	doc.AssertDefinition("Inner")
}

// --- Find references ---

func TestReferencesOnTokenModeDeclaration(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(<|Inner|>)
	hidden WS
}
token mode <|Inner@><|Inner|> {
	ID -> pop
}`, "file:///references.fb")
	doc.AssertNoErrors()
	doc.AssertReferences("Inner")
}

func TestReferencesOnTokenModeFromUsage(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
token NAME: /[A-Z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(<|Inner@><|Inner|>)
	NAME -> push(<|Inner|>)
	hidden WS
}
token mode <|Inner|> {
	ID -> pop
}`, "file:///references-usage.fb")
	doc.AssertNoErrors()
	doc.AssertReferences("Inner")
}

// --- Rename ---

func TestRenameTokenMode(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(Inner)
	hidden WS
}
token mode <|target@>Inner {
	ID -> pop
}`, "file:///rename.fb")
	doc.AssertNoErrors()
	doc.RunRename("target", "Renamed")

	text := doc.Document.TextDoc.Text(nil)
	assert.Contains(t, text, "token mode Renamed {")
	assert.Contains(t, text, "push(Renamed)")
	assert.NotContains(t, text, "Inner")
	doc.AssertNoErrors()
}

func TestRenameTokenModeFromUsage(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(<|target@>Inner)
	hidden WS
}
token mode Inner {
	ID -> pop
}`, "file:///rename-usage.fb")
	doc.AssertNoErrors()
	doc.RunRename("target", "Renamed")

	text := doc.Document.TextDoc.Text(nil)
	assert.Contains(t, text, "token mode Renamed {")
	assert.Contains(t, text, "push(Renamed)")
	doc.AssertNoErrors()
}

// --- Document symbols ---

func TestDocumentSymbolForTokenMode(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID
	hidden WS
}
<|mode:token mode Inner {
	ID
}|>`, "file:///symbols-mode.fb")
	doc.AssertNoParseErrors()
	doc.AssertDocumentSymbol("mode", "Inner", lsp.Namespace)
}

func TestDocumentSymbolForTokenGroup(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=Punct;
<|group:token group Punct { "!" }|>
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	Punct
	hidden WS
}`, "file:///symbols-group.fb")
	doc.AssertNoParseErrors()
	doc.AssertDocumentSymbol("group", "Punct", lsp.Enum)
}

func TestDocumentSymbolForModeLocalTokenDeclaration(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID
	hidden WS
}
token mode Inner {
	<|inner:token INNER: /[A-Z]+/|>
}`, "file:///symbols-inner.fb")
	doc.AssertNoParseErrors()
	// A token declared inside a mode is a symbol in its own right, nested under
	// the mode it belongs to.
	doc.AssertDocumentSymbol("inner", "INNER", lsp.Constant)
	symbol := doc.MustFindSymbolAtLabel("inner")
	assert.Equal(t, "INNER", symbol.Name)
}

func TestDocumentSymbolSkipsDefaultTokenMode(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID
	hidden WS
}`, "file:///symbols-default.fb")
	doc.AssertNoParseErrors()
	// The default mode has no name of its own, so it contributes no symbol.
	provider := service.MustGet[server.DocumentSymbolProvider](f.Services())
	result, err := provider.HandleDocumentSymbolRequest(f.Ctx(), &lsp.DocumentSymbolParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: doc.Document.URI.DocumentURI()},
	})
	assert.NoError(t, err)
	for _, symbol := range result {
		assert.NotEqual(t, "default", symbol.Name)
		for _, child := range symbol.Children {
			assert.NotEqual(t, "default", child.Name)
		}
	}
}

// --- Folding ranges ---

func TestFoldingRangeForTokenMode(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
<|default:token mode default {
	ID -> push(Inner)
	hidden WS|>
}
<|inner:token mode Inner {
	ID
	ID -> pop|>
}`, "file:///folding.fb")
	doc.AssertNoParseErrors()
	doc.AssertFoldingRanges("default", "inner")
}

// --- Completion ---

func TestCompletionOfTokenModeNamesInPushCommand(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(<|cursor@>)
	hidden WS
}
token mode Inner {
	ID -> pop
}
token mode Other {
	ID -> pop
}`, "file:///completion-mode.fb")
	labels := completionLabels(doc.CompletionItems("cursor"))
	assert.Contains(t, labels, "Inner")
	assert.Contains(t, labels, "Other")
	// The default mode can be targeted by name as well.
	assert.Contains(t, labels, "default")
}

func TestCompletionOfTokenModeNamesInSetModeCommand(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(Inner)
	hidden WS
}
token mode Inner {
	ID -> mode(<|cursor@>)
}`, "file:///completion-set.fb")
	labels := completionLabels(doc.CompletionItems("cursor"))
	assert.Contains(t, labels, "default")
	assert.Contains(t, labels, "Inner")
}

func TestCompletionOfCommandKeywords(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> <|cursor@>
	hidden WS
}`, "file:///completion-command.fb")
	labels := completionLabels(doc.CompletionItems("cursor"))
	assert.Contains(t, labels, "push")
	assert.Contains(t, labels, "pop")
	assert.Contains(t, labels, "mode")
}

// --- Hover ---

func TestHoverOnTokenModeDeclaration(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(Inner)
	hidden WS
}
// The inner mode
token mode <|inner:Inner|> {
	ID -> pop
}`, "file:///hover.fb")
	doc.AssertNoErrors()
	doc.ExpectHoverAt("inner", "The inner mode")
}

func TestHoverOnTokenModeReference(t *testing.T) {
	f := lspFixture(t)
	doc := f.ParseURI(`grammar Test;
interface Foo { Greeting string }
Foo: Greeting=ID;
token ID: /[a-z]+/
hidden token WS: /\s+/
token mode default {
	ID -> push(<|usage:Inner|>)
	hidden WS
}
// The inner mode
token mode Inner {
	ID -> pop
}`, "file:///hover-usage.fb")
	doc.AssertNoErrors()
	// Hovering a reference shows the documentation of the mode it points to.
	doc.ExpectHoverAt("usage", "The inner mode")
}
