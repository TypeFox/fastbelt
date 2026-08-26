// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// TestResetKeepsLocalSymbols is a regression test for the Reset fallthrough
// cascade: resetting a document with the updater's keep-state
// (Parsed|ExportedSymbols|LocalSymbols) used to nil doc.LocalSymbols while
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

	builder := service.MustGet[workspace.Builder](sc)
	keep := core.DocStateParsed | core.DocStateExportedSymbols | core.DocStateLocalSymbols
	builder.Reset(doc.Document, keep)

	if doc.Document.LocalSymbols == nil {
		t.Fatal("Reset dropped LocalSymbols although DocStateLocalSymbols was kept")
	}
	if doc.Document.ImportedSymbols != nil {
		t.Error("Reset kept ImportedSymbols although DocStateImportedSymbols was dropped")
	}

	if err := builder.Build(f.Ctx(), []*core.Document{doc.Document}); err != nil {
		t.Fatalf("rebuild after reset failed: %v", err)
	}
	for _, ref := range doc.Document.References {
		if err := ref.Error(); err != nil {
			t.Errorf("reference failed to resolve after relink: %v", err)
		}
	}
}
