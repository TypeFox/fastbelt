package test

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// Performs a rename operation using the server's RenameProvider at the position of the given marker label.
// Updates the affected documents with the edits returned and rebuilds the workspace.
func (d *Doc) RunRename(label string, newName string) *Doc {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, true)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	renameProvider, err := service.Get[server.RenameProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no rename provider available: %v", err)
	}
	edits, err := renameProvider.HandleRenameRequest(d.fixture.ctx, &lsp.RenameParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: d.Document.URI.DocumentURI()},
			Position:     location.LspPosition(),
		},
		NewName: newName,
	})
	if err != nil {
		d.fixture.t.Fatalf("fbtest: rename failed: %v", err)
	}
	documents := service.MustGet[workspace.DocumentManager](d.fixture.sc)
	toUpdate := []*core.Document{}
	for uri, fileEdits := range edits.Changes {
		doc := documents.Get(core.ParseURI(string(uri)))
		if doc == nil {
			d.fixture.t.Fatalf("fbtest: rename edits include unknown document URI %q", uri)
			return nil
		}
		changeEvents := make([]lsp.TextDocumentContentChangeEvent, len(fileEdits))
		for i, edit := range fileEdits {
			changeEvents[i] = lsp.TextDocumentContentChangeEvent{
				Range: &edit.Range,
				Text:  edit.NewText,
			}
		}
		overlay := doc.TextDoc.(*textdoc.Overlay)
		updateErr := overlay.Update(changeEvents, overlay.Version()+1)
		if updateErr != nil {
			d.fixture.t.Fatalf("fbtest: failed to apply rename edits: %v", updateErr)
		}
		toUpdate = append(toUpdate, doc)
	}
	// Rebuild the affected documents to ensure that the changes are reflected in the AST
	builder := service.MustGet[workspace.Builder](d.fixture.sc)
	for _, doc := range toUpdate {
		// Fully reset each document
		builder.Reset(doc, 0)
	}
	if err := builder.Build(d.fixture.ctx, toUpdate, nil); err != nil {
		d.fixture.t.Fatalf("fbtest: build failed after rename: %v", err)
	}
	return d
}

func (d *Doc) AssertDefinition(label string) {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, false)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	defProvider, err := service.Get[server.DefinitionProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no definition provider available: %v", err)
	}
	links, err := defProvider.HandleDefinitionRequest(d.fixture.ctx, &lsp.DefinitionParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: d.Document.URI.DocumentURI()},
			Position:     location.LspPosition(),
		},
	})
	if err != nil {
		d.fixture.t.Fatalf("fbtest: definition request failed: %v", err)
	}
	expectedLinks := []lsp.DefinitionLink{}
	for _, doc := range d.fixture.docs {
		ranges := doc.markerRanges(label)
		for _, r := range ranges {
			expectedLinks = append(expectedLinks, lsp.DefinitionLink{
				TargetURI:   doc.Document.URI.DocumentURI(),
				TargetRange: r.LspRange(),
			})
		}
	}
	if len(links) != len(expectedLinks) {
		d.fixture.t.Fatalf("fbtest: expected %d definition links, got %d", len(expectedLinks), len(links))
	}
	for _, expected := range expectedLinks {
		found := false
		for _, link := range links {
			if link.TargetURI == expected.TargetURI && link.TargetRange == expected.TargetRange {
				found = true
				break
			}
		}
		if !found {
			d.fixture.t.Errorf("fbtest: expected definition link not found: target URI %q, target range %v", expected.TargetURI, expected.TargetRange)
		}
	}
}

func (d *Doc) AssertReferences(label string) {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, false)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	refsProvider, err := service.Get[server.ReferencesProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no references provider available: %v", err)
	}
	locations, err := refsProvider.HandleReferencesRequest(d.fixture.ctx, &lsp.ReferenceParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: d.Document.URI.DocumentURI()},
			Position:     location.LspPosition(),
		},
	})
	if err != nil {
		d.fixture.t.Fatalf("fbtest: references request failed: %v", err)
	}
	expectedLocations := []lsp.Location{}
	for _, doc := range d.fixture.docs {
		ranges := doc.markerRanges(label)
		for _, r := range ranges {
			expectedLocations = append(expectedLocations, lsp.Location{
				URI:   doc.Document.URI.DocumentURI(),
				Range: r.LspRange(),
			})
		}
	}
	if len(locations) != len(expectedLocations) {
		d.fixture.t.Fatalf("fbtest: expected %d reference locations, got %d", len(expectedLocations), len(locations))
	}
	for _, expected := range expectedLocations {
		found := false
		for _, loc := range locations {
			if loc.URI == expected.URI && loc.Range == expected.Range {
				found = true
				break
			}
		}
		if !found {
			d.fixture.t.Errorf("fbtest: expected reference location not found: URI %q, range %v", expected.URI, expected.Range)
		}
	}
}

// AssertFoldingRanges verifies that folding ranges exist for all specified marker labels.
// Returns the Doc for chaining.
func (d *Doc) AssertFoldingRanges(labels ...string) *Doc {
	d.fixture.t.Helper()
	frProvider, err := service.Get[server.FoldingRangeProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no folding range provider available: %v", err)
	}
	params := &lsp.FoldingRangeParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: lsp.DocumentURI(d.Document.URI.DocumentURI()),
		},
	}
	result, err := frProvider.HandleFoldingRangeRequest(d.fixture.ctx, params)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: folding range request failed: %v", err)
	}
	for _, label := range labels {
		expectedRange, err := d.markerRange(label)
		if err != nil {
			d.fixture.t.Fatalf("fbtest: %v", err)
		}
		found := false
		for _, fr := range result {
			if fr.StartLine != nil && fr.EndLine != nil {
				if *fr.StartLine == uint32(expectedRange.Start.Line) &&
					*fr.EndLine == uint32(expectedRange.End.Line) {
					// For comment ranges, also verify the kind
					if label == "comment" && fr.Kind != "comment" {
						continue
					}
					found = true
					break
				}
			}
		}
		if !found {
			d.fixture.t.Errorf(
				"fbtest: expected folding range at label %q (lines %d-%d) not found",
				label, expectedRange.Start.Line, expectedRange.End.Line,
			)
		}
	}
	return d
}
