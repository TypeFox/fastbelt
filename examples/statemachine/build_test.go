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

// syncFixture drives the [server.DocumentSynchronizer] of a fresh service
// container the way an LSP client would.
type syncFixture struct {
	t            *testing.T
	synchronizer server.DocumentSynchronizer
	lock         workspace.Lock
	docs         workspace.DocumentManager
}

func newSyncFixture(t *testing.T) *syncFixture {
	sc := CreateLspServices(nil)
	return &syncFixture{
		t:            t,
		synchronizer: sc.MustGet[server.DocumentSynchronizer](),
		lock:         sc.MustGet[workspace.Lock](),
		docs:         sc.MustGet[workspace.DocumentManager](),
	}
}

func (f *syncFixture) open(uri lsp.DocumentURI, text string) {
	f.synchronizer.DidOpen(context.Background(), &lsp.DidOpenTextDocumentParams{TextDocument: lsp.TextDocumentItem{
		URI: uri, LanguageID: "statemachine", Version: 1, Text: text,
	}})
}

// insert inserts text at the start of the document.
func (f *syncFixture) insert(uri lsp.DocumentURI, version int32, text string) {
	f.synchronizer.DidChange(context.Background(), &lsp.DidChangeTextDocumentParams{
		TextDocument:   lsp.VersionedTextDocumentIdentifier{URI: uri, Version: version},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{{Range: &lsp.Range{}, Text: text}},
	})
}

func (f *syncFixture) close(uri lsp.DocumentURI) {
	f.synchronizer.DidClose(context.Background(), &lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// requireBuilt fails the test unless both a reader scoped to uri and a
// workspace-wide reader are admitted at the validated state.
func (f *syncFixture) requireBuilt(uri lsp.DocumentURI) {
	f.t.Helper()
	coreURI := fastbelt.ParseURI(string(uri))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(f.t, f.lock.ReadAt(ctx, fastbelt.DocStateValidated, []fastbelt.URI{coreURI}, func(context.Context) {
		assert.True(f.t, f.docs.Get(coreURI).State().IsComplete())
	}), "the document was not built after the last write")
	require.NoError(f.t, f.lock.ReadAt(ctx, fastbelt.DocStateValidated, nil, func(context.Context) {}),
		"the workspace was not built after the last write")
}

func TestBurstOfEditsEndsWithBuiltDocument(t *testing.T) {
	f := newSyncFixture(t)
	uri := lsp.DocumentURI("file:///test/example.statemachine")
	f.open(uri, "")
	const edits = 200
	for version := int32(2); version <= edits+1; version++ {
		f.insert(uri, version, " ")
	}
	f.requireBuilt(uri)
}

// TestWriteWithoutChangeKeepsWorkspaceBuilt covers writes that change no
// document. Acquiring the write lock cancels the running build and resets the
// workspace floor, so such a write must still finish the build.
func TestWriteWithoutChangeKeepsWorkspaceBuilt(t *testing.T) {
	uri := lsp.DocumentURI("file:///test/example.statemachine")
	other := lsp.DocumentURI("file:///test/other.statemachine")
	tests := []struct {
		name  string
		write func(f *syncFixture)
	}{
		{"didClose", func(f *syncFixture) { f.close(other) }},
		{"didChange without effect", func(f *syncFixture) { f.insert(uri, 100, "") }},
		{"didChange of unopened document", func(f *syncFixture) {
			f.insert("file:///test/unopened.statemachine", 2, " ")
		}},
		{"didOpen with unchanged content", func(f *syncFixture) { f.open(uri, " ") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newSyncFixture(t)
			f.open(other, "")
			f.open(uri, "")
			// Start a build for the write under test to cancel.
			f.insert(uri, 2, " ")
			tt.write(f)
			f.requireBuilt(uri)
		})
	}
}

// TestBuildStepListenersRunAfterCancellation verifies that a document which
// completed a build step notifies every listener of that step, even if the
// build was cancelled in the meantime. The step is not repeated by a later
// build, so a skipped listener would never see the document.
func TestBuildStepListenersRunAfterCancellation(t *testing.T) {
	sc := CreateLspServices(nil)
	f := test.New(t, sc)
	doc := f.Parse(`
		statemachine Toggle

		initialState a

		state a
		end
	`).AssertNoErrors()

	builder := sc.MustGet[workspace.Builder]()
	ctx, cancel := context.WithCancel(f.Ctx())
	defer cancel()
	// Listeners are called in registration order: the first one stands in for
	// a newer write that cancels the build right after validation finished.
	builder.AddBuildStepListener(fastbelt.DocStateValidated, func(context.Context, *fastbelt.Document) error {
		cancel()
		return nil
	})
	called := 0
	builder.AddBuildStepListener(fastbelt.DocStateValidated, func(listenerCtx context.Context, _ *fastbelt.Document) error {
		assert.ErrorIs(t, listenerCtx.Err(), context.Canceled, "the listener should observe the cancellation")
		called++
		return nil
	})

	builder.Reset(doc.Document, fastbelt.DocStateComplete&^fastbelt.DocStateValidated)
	err := builder.Build(ctx, []*fastbelt.Document{doc.Document})

	assert.ErrorIs(t, err, context.Canceled)
	assert.True(t, doc.Document.State().Has(fastbelt.DocStateValidated), "validation finished before the cancellation")
	assert.Equal(t, 1, called, "the listener was skipped for a validated document")
}
