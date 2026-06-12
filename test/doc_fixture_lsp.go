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

func (d *Doc) AssertHighlights(label string) {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, false)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	highlightProvider, err := service.Get[server.DocumentHighlightProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no document highlight provider available: %v", err)
	}
	highlights, err := highlightProvider.HandleDocumentHighlightRequest(d.fixture.ctx, &lsp.DocumentHighlightParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: d.Document.URI.DocumentURI()},
			Position:     location.LspPosition(),
		},
	})
	if err != nil {
		d.fixture.t.Fatalf("fbtest: document highlight request failed: %v", err)
	}
	expectedLocations := []lsp.DocumentHighlight{}
	for _, doc := range d.fixture.docs {
		ranges := doc.markerRanges(label)
		for _, r := range ranges {
			expectedLocations = append(expectedLocations, lsp.DocumentHighlight{
				Range: r.LspRange(),
			})
		}
	}
	if len(highlights) != len(expectedLocations) {
		d.fixture.t.Fatalf("fbtest: expected %d document highlights, got %d", len(expectedLocations), len(highlights))
	}
	for _, expected := range expectedLocations {
		found := false
		for _, highlight := range highlights {
			if highlight.Range == expected.Range {
				found = true
				break
			}
		}
		if !found {
			d.fixture.t.Errorf("fbtest: expected document highlight not found: range %v", expected.Range)
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
		expectedRange, err := d.MarkerRange(label)
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

// ExpectHoverAt asserts that a hover at the named marker returns Markdown content equal to markup.
// For range markers the start of the range is used as the cursor position.
// Returns the Doc for chaining.
func (d *Doc) ExpectHoverAt(label, markup string) *Doc {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, true)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	hoverProvider, err := service.Get[server.HoverProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no hover provider available: %v", err)
	}
	result, err := hoverProvider.HandleHoverRequest(d.fixture.ctx, &lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: d.Document.URI.DocumentURI()},
			Position:     location.LspPosition(),
		},
	})
	if err != nil {
		d.fixture.t.Fatalf("fbtest: hover request failed: %v", err)
	}
	if result == nil {
		d.fixture.t.Fatalf("fbtest: expected hover result at %q, got nil", label)
		return d
	}
	if result.Contents.Kind != lsp.Markdown {
		d.fixture.t.Errorf("fbtest: expected Markdown content kind at %q, got %v", label, result.Contents.Kind)
	}
	if result.Contents.Value != markup {
		d.fixture.t.Errorf("fbtest: expected hover value %q at %q, got %q", markup, label, result.Contents.Value)
	}
	expectedRange, err := d.MarkerRange(label)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	if result.Range != expectedRange.LspRange() {
		d.fixture.t.Errorf("fbtest: expected hover range %v at %q, got %v", expectedRange.LspRange(), label, result.Range)
	}
	return d
}

// ExpectNoHoverAt asserts that a hover at the named marker returns nil.
// For range markers the start of the range is used as the cursor position.
// Returns the Doc for chaining.
func (d *Doc) ExpectNoHoverAt(label string) *Doc {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, true)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
	}
	hoverProvider, err := service.Get[server.HoverProvider](d.fixture.sc)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: no hover provider available: %v", err)
	}
	result, err := hoverProvider.HandleHoverRequest(d.fixture.ctx, &lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: d.Document.URI.DocumentURI()},
			Position:     location.LspPosition(),
		},
	})
	if err != nil {
		d.fixture.t.Fatalf("fbtest: hover request failed: %v", err)
	}
	if result != nil {
		d.fixture.t.Errorf("fbtest: expected nil hover at %q, got non-nil result", label)
	}
	return d
}

func (d *Doc) CompletionItems(label string) []lsp.CompletionItem {
	d.fixture.t.Helper()
	location, err := d.markerLocation(label, true)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: CompletionItems: no marker with label %q", label)
	}
	provider := service.MustGet[server.CompletionProvider](d.fixture.sc)
	resp, err := provider.HandleCompletionRequest(d.Ctx(), &lsp.CompletionParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{
				URI: d.Document.TextDoc.URI(),
			},
			Position: lsp.Position{Line: uint32(location.Line), Character: uint32(location.Column)},
		},
	})
	if err != nil {
		d.fixture.t.Fatalf("HandleCompletionRequest returned error: %v", err)
	}
	return resp.Items
}
