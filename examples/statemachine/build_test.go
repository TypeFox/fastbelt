// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// TestResetKeepsLocalSymbols is a regression test for the Reset fallthrough
// cascade: resetting a document with the updater's keep-state
// (Lexed|Parsed|ExportedSymbols|LocalSymbols) used to nil doc.LocalSymbols while
// keeping the DocStateLocalSymbols bit, so the following build skipped
// recomputing local symbols and relinked against nil.
func TestResetKeepsLocalSymbols(t *testing.T) {
	sc := CreateLspServices(nil)
	f := test.New(t, sc)
	doc := f.Parse(`
		statemachine Toggle

		events flick

		initialState a

		state a
		  flick => b
		end

		state b
		  flick => a
		end
	`).AssertNoErrors()

	builder := sc.MustGet[workspace.Builder]()
	keep := fastbelt.DocStateLexed | fastbelt.DocStateParsed | fastbelt.DocStateExportedSymbols | fastbelt.DocStateLocalSymbols
	builder.Reset(doc.Document, keep)

	if doc.Document.LocalSymbols == nil {
		t.Fatal("Reset dropped LocalSymbols although DocStateLocalSymbols was kept")
	}
	if doc.Document.ImportedSymbols != nil {
		t.Error("Reset kept ImportedSymbols although DocStateImportedSymbols was dropped")
	}

	if err := builder.Build(f.Ctx(), []*fastbelt.Document{doc.Document}); err != nil {
		t.Fatalf("rebuild after reset failed: %v", err)
	}
	for _, ref := range doc.Document.References {
		if err := ref.Error(); err != nil {
			t.Errorf("reference failed to resolve after relink: %v", err)
		}
	}
}

func TestBurstOfEditsEndsWithBuiltDocument(t *testing.T) {
	sc := CreateLspServices(nil)
	synchronizer := sc.MustGet[server.DocumentSynchronizer]()
	lock := sc.MustGet[workspace.Lock]()
	docs := sc.MustGet[workspace.DocumentManager]()
	ctx := context.Background()
	uri := lsp.DocumentURI("file:///test/example.statemachine")
	coreURI := fastbelt.ParseURI(string(uri))
	synchronizer.DidOpen(ctx, &lsp.DidOpenTextDocumentParams{TextDocument: lsp.TextDocumentItem{
		URI: uri, LanguageID: "statemachine", Version: 1, Text: "", // example.statemachine
	}})
	const edits = 200
	for version := int32(2); version <= edits+1; version++ {
		synchronizer.DidChange(ctx, &lsp.DidChangeTextDocumentParams{
			TextDocument:   lsp.VersionedTextDocumentIdentifier{URI: uri, Version: version},
			ContentChanges: []lsp.TextDocumentContentChangeEvent{{Range: &lsp.Range{}, Text: " "}},
		})
	}
	readCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	require.NoError(t, lock.ReadAt(readCtx, fastbelt.DocStateValidated, []fastbelt.URI{coreURI}, func(context.Context) {
		assert.True(t, docs.Get(coreURI).State().IsComplete())
	}), "the document was not built after the last edit") // context deadline exceeded
}
