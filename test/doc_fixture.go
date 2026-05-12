package test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// Doc wraps a built [core.Document] with assertion methods and marker positions.
// Access the underlying document via the exported Document field.
type Doc struct {
	Document *core.Document
	Ranges   []RangeMarker
	Indices  []IndexMarker
	fixture  *Fixture
}

// Ctx returns the context for this document. Pass it to [core.Reference.Ref].
func (d *Doc) Ctx() context.Context { return d.fixture.ctx }

// Root returns the root AST node of the document.
func (d *Doc) Root() core.AstNode { return d.Document.Root }

// Diagnostics returns all diagnostics collected during validation.
func (d *Doc) Diagnostics() []*core.Diagnostic { return d.Document.Diagnostics }

// AssertNoErrors fails the test if any error-severity Diagnostics exist.
func (d *Doc) AssertNoErrors() *Doc {
	d.fixture.t.Helper()
	// The diagnostics contain all kinds of errors
	for _, diag := range d.Document.Diagnostics {
		if diag.Severity == core.SeverityError {
			d.fixture.t.Errorf("fbtest: unexpected error diagnostic: %s", diag.Message)
		}
	}
	return d
}

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

// AssertNoParseErrors fails the test if any LexerErrors or ParserErrors exist.
// Diagnostics are not checked.
func (d *Doc) AssertNoParseErrors() *Doc {
	d.fixture.t.Helper()
	for _, e := range d.Document.LexerErrors {
		d.fixture.t.Errorf("fbtest: unexpected lexer error: %v", e)
	}
	for _, e := range d.Document.ParserErrors {
		d.fixture.t.Errorf("fbtest: unexpected parser error: %v", e)
	}
	return d
}

type DiagnosticExpectation struct {
	t          testing.TB
	Diagnostic *core.Diagnostic
}

func (d *DiagnosticExpectation) WithMessage(msg string) *DiagnosticExpectation {
	d.t.Helper()
	if d.Diagnostic.Message != msg {
		d.t.Errorf("fbtest: expected diagnostic message %q, got %q", msg, d.Diagnostic.Message)
	}
	return d
}

func (d *DiagnosticExpectation) WithMessageContaining(substring string) *DiagnosticExpectation {
	d.t.Helper()
	if !strings.Contains(d.Diagnostic.Message, substring) {
		d.t.Errorf("fbtest: expected diagnostic message containing %q, got %q", substring, d.Diagnostic.Message)
	}
	return d
}

func (d *DiagnosticExpectation) WithCode(code string) *DiagnosticExpectation {
	d.t.Helper()
	if d.Diagnostic.Code != code {
		d.t.Errorf("fbtest: expected diagnostic code %q, got %q", code, d.Diagnostic.Code)
	}
	return d
}

func (d *DiagnosticExpectation) WithSource(source string) *DiagnosticExpectation {
	d.t.Helper()
	if d.Diagnostic.Source != source {
		d.t.Errorf("fbtest: expected diagnostic source %q, got %q", source, d.Diagnostic.Source)
	}
	return d
}

func (d *DiagnosticExpectation) WithSeverity(severity core.DiagnosticSeverity) *DiagnosticExpectation {
	d.t.Helper()
	if d.Diagnostic.Severity != severity {
		d.t.Errorf("fbtest: expected diagnostic severity %q, got %q", severity, d.Diagnostic.Severity)
	}
	return d
}

func (d *DiagnosticExpectation) WithTags(tags ...core.DiagnosticTag) *DiagnosticExpectation {
	d.t.Helper()
	if len(d.Diagnostic.Tags) != len(tags) {
		d.t.Errorf("fbtest: expected diagnostic tags %v, got %v", tags, d.Diagnostic.Tags)
		return d
	}
	for _, tag := range tags {
		found := slices.Contains(d.Diagnostic.Tags, tag)
		if !found {
			d.t.Errorf("fbtest: expected diagnostic tag %v not found in actual tags %v", tag, d.Diagnostic.Tags)
		}
	}
	return d
}

// AssertDiagnostic fails the test unless at least one diagnostic has the given
// severity and a message containing msgSubstring.
func (d *Doc) ExpectDiagnostic(label string) *DiagnosticExpectation {
	d.fixture.t.Helper()
	loc, err := d.markerRange(label)
	if err != nil {
		d.fixture.t.Fatalf("fbtest: %v", err)
		return nil
	}
	for _, diag := range d.Document.Diagnostics {
		if diag.Range == loc {
			return &DiagnosticExpectation{t: d.fixture.t, Diagnostic: diag}
		}
	}
	d.fixture.t.Fatalf("fbtest: no diagnostic at label %q", label)
	return nil
}

// AssertNoDiagnostics fails the test if any diagnostics exist at any severity.
func (d *Doc) AssertNoDiagnostics() *Doc {
	d.fixture.t.Helper()
	for _, diag := range d.Document.Diagnostics {
		d.fixture.t.Errorf("fbtest: unexpected diagnostic [%s]: %s", diag.Severity, diag.Message)
	}
	return d
}

// AssertState fails the test unless the document state includes the given flag.
func (d *Doc) AssertState(flag core.DocumentState) *Doc {
	d.fixture.t.Helper()
	if !d.Document.State.Has(flag) {
		d.fixture.t.Errorf("fbtest: document state does not include %v (actual: %v)", flag, d.Document.State)
	}
	return d
}

func (d *Doc) markerRange(label string) (core.TextRange, error) {
	ranges := d.markerRanges(label)
	if len(ranges) == 0 {
		return core.TextRange{}, fmt.Errorf("no marker with label %q", label)
	} else if len(ranges) > 1 {
		return core.TextRange{}, fmt.Errorf("multiple markers with label %q found; expected exactly one", label)
	}
	return ranges[0], nil
}

func (d *Doc) markerRanges(label string) []core.TextRange {
	var result []core.TextRange
	for _, r := range d.Ranges {
		if r.Label == label {
			start := d.Document.TextDoc.PositionAt(r.Start)
			end := d.Document.TextDoc.PositionAt(r.End)
			result = append(result, core.TextRange{Start: core.TextLocation{
				Line:   core.TextLine(start.Line),
				Column: core.TextColumn(start.Character),
			}, End: core.TextLocation{
				Line:   core.TextLine(end.Line),
				Column: core.TextColumn(end.Character),
			}})
		}
	}
	return result
}

// markerLocation returns the TextLocation of the named marker.
// For range markers it returns the start; for index markers the offset position.
func (d *Doc) markerLocations(label string, includeRanges bool) []core.TextLocation {
	var result []core.TextLocation
	for _, idx := range d.Indices {
		if idx.Label == label {
			position := d.Document.TextDoc.PositionAt(idx.Offset)
			result = append(result, core.TextLocation{
				Line:   core.TextLine(position.Line),
				Column: core.TextColumn(position.Character),
			})
		}
	}
	if includeRanges {
		for _, r := range d.Ranges {
			if r.Label == label {
				start := d.Document.TextDoc.PositionAt(r.Start)
				result = append(result, core.TextLocation{
					Line:   core.TextLine(start.Line),
					Column: core.TextColumn(start.Character),
				})
			}
		}
	}
	return result
}

func (d *Doc) markerLocation(label string, includeRanges bool) (core.TextLocation, error) {
	locations := d.markerLocations(label, includeRanges)
	if len(locations) == 0 {
		return core.TextLocation{}, fmt.Errorf("no marker with label %q", label)
	} else if len(locations) > 1 {
		return core.TextLocation{}, fmt.Errorf("multiple markers with label %q found; expected exactly one", label)
	}
	return locations[0], nil
}

// FindNode returns the first node of type T in the document tree, or (zero, false).
func FindNode[T core.AstNode](d *Doc) (T, bool) {
	d.fixture.t.Helper()
	var zero T
	if d.Document.Root == nil {
		return zero, false
	}
	for node := range core.AllNodes(d.Document.Root) {
		if n, ok := node.(T); ok {
			return n, true
		}
	}
	return zero, false
}

// MustFindNode fails the test if no node of type T is found in the document tree.
func MustFindNode[T core.AstNode](d *Doc) T {
	d.fixture.t.Helper()
	n, ok := FindNode[T](d)
	if !ok {
		d.fixture.t.Fatal("fbtest: MustFindNode: no node of the requested type found")
	}
	return n
}

// FindAll returns all nodes of type T in the document tree.
func FindAll[T core.AstNode](d *Doc) []T {
	d.fixture.t.Helper()
	if d.Document.Root == nil {
		return nil
	}
	var result []T
	for node := range core.AllNodes(d.Document.Root) {
		if n, ok := node.(T); ok {
			result = append(result, n)
		}
	}
	return result
}

// FindNodeAtOffset returns the most specific (smallest span) node of type T whose
// text span contains the given byte offset.
func FindNodeAtOffset[T core.AstNode](d *Doc, offset int) (T, bool) {
	d.fixture.t.Helper()
	var zero T
	if d.Document.Root == nil {
		return zero, false
	}
	var result T
	found := false
	bestSpan := 0
	for node := range core.AllNodes(d.Document.Root) {
		n, ok := node.(T)
		if !ok {
			continue
		}
		seg := node.Segment()
		if seg == nil {
			continue
		}
		start, end := int(seg.Indices.Start), int(seg.Indices.End)
		if start <= offset && offset < end {
			span := end - start
			if !found || span < bestSpan {
				result, bestSpan, found = n, span, true
			}
		}
	}
	return result, found
}

// MustFindNodeAtOffset fails the test if no node of type T contains the given offset.
func MustFindNodeAtOffset[T core.AstNode](d *Doc, offset int) T {
	d.fixture.t.Helper()
	n, ok := FindNodeAtOffset[T](d, offset)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindNodeAtOffset: no node of the requested type at offset %d", offset)
	}
	return n
}

// FindNodeAtLocation returns the most specific (smallest span) node of type T whose
// text range contains the given line/column position.
func FindNodeAtLocation[T core.AstNode](d *Doc, location core.TextLocation) (T, bool) {
	d.fixture.t.Helper()
	var zero T
	if d.Document.Root == nil {
		return zero, false
	}
	var result T
	found := false
	bestSpan := 0
	for node := range core.AllNodes(d.Document.Root) {
		n, ok := node.(T)
		if !ok {
			continue
		}
		seg := node.Segment()
		if seg == nil {
			continue
		}
		if locationInRange(location, seg.Range) {
			span := int(seg.Indices.End - seg.Indices.Start)
			if !found || span < bestSpan {
				result, bestSpan, found = n, span, true
			}
		}
	}
	return result, found
}

// MustFindNodeAtLocation fails the test if no node of type T contains the given location.
func MustFindNodeAtLocation[T core.AstNode](d *Doc, location core.TextLocation) T {
	d.fixture.t.Helper()
	n, ok := FindNodeAtLocation[T](d, location)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindNodeAtLocation: no node of the requested type at %v", location)
	}
	return n
}

// FindNodeAtLabel returns the most specific node of type T at the position of the
// named range marker (start offset) or index marker.
func FindNodeAtLabel[T core.AstNode](d *Doc, label string) (T, bool) {
	d.fixture.t.Helper()
	var zero T
	for _, r := range d.Ranges {
		if r.Label == label {
			return FindNodeAtOffset[T](d, r.Start)
		}
	}
	for _, idx := range d.Indices {
		if idx.Label == label {
			return FindNodeAtOffset[T](d, idx.Offset)
		}
	}
	return zero, false
}

// MustFindNodeAtLabel fails the test if no node of type T is found at the given label.
func MustFindNodeAtLabel[T core.AstNode](d *Doc, label string) T {
	d.fixture.t.Helper()
	n, ok := FindNodeAtLabel[T](d, label)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindNodeAtLabel: no node of the requested type at label %q", label)
	}
	return n
}

// FindNamedNode returns the first node of type T whose Name() method returns the given name.
// T must implement [core.NamedNode] (or have a Name() string method) to find results.
func FindNamedNode[T core.AstNode](d *Doc, name string) (T, bool) {
	d.fixture.t.Helper()
	var zero T
	if d.Document.Root == nil {
		return zero, false
	}
	for node := range core.AllNodes(d.Document.Root) {
		n, ok := node.(T)
		if !ok {
			continue
		}
		if named, ok := node.(core.NamedNode); ok && named.Name() == name {
			return n, true
		}
	}
	return zero, false
}

// MustFindNamedNode fails the test if no node of type T with the given name is found.
func MustFindNamedNode[T core.AstNode](d *Doc, name string) T {
	d.fixture.t.Helper()
	n, ok := FindNamedNode[T](d, name)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindNamedNode: no node of the requested type with name %q", name)
	}
	return n
}

// FindReferenceWithText returns the first reference of type T whose cross-reference
// text matches the given string exactly.
func FindReferenceWithText[T core.AstNode](d *Doc, text string) (*core.Reference[T], bool) {
	d.fixture.t.Helper()
	if d.Document.Root == nil {
		return nil, false
	}
	// Documents are always built, so References should be filled
	for _, ref := range d.Document.References {
		if ref.Text() == text {
			if typedRef, ok := ref.(*core.Reference[T]); ok {
				return typedRef, true
			}
		}
	}
	return nil, false
}

// MustFindReferenceWithText fails the test if no reference of type T with the given
// text is found.
func MustFindReferenceWithText[T core.AstNode](d *Doc, text string) *core.Reference[T] {
	d.fixture.t.Helper()
	r, ok := FindReferenceWithText[T](d, text)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindReferenceWithText: no reference of the requested type with text %q", text)
	}
	return r
}

// FindReference returns the first reference of type T whose cross-reference text is located at the given label.
// If the label can be found multiple times, the function returns (nil, false)
func FindReference[T core.AstNode](d *Doc, label string) (*core.Reference[T], bool) {
	d.fixture.t.Helper()
	targetRange, err := d.markerRange(label)
	if err != nil || d.Document.Root == nil {
		return nil, false
	}
	for _, ref := range d.Document.References {
		segment := ref.Segment()
		if segment == nil {
			continue
		}
		if targetRange == segment.Range {
			if typedRef, ok := ref.(*core.Reference[T]); ok {
				return typedRef, true
			}
		}
	}
	return nil, false
}

// MustFindReference fails the test if no reference of type T is found at the given label.
func MustFindReference[T core.AstNode](d *Doc, label string) *core.Reference[T] {
	d.fixture.t.Helper()
	r, ok := FindReference[T](d, label)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindReference: no reference of the requested type at label %q", label)
	}
	return r
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
		expectedRange, ok := d.MarkerRange(label)
		if !ok {
			d.fixture.t.Fatalf("fbtest: no marker with label %q", label)
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
			d.fixture.t.Errorf("fbtest: expected folding range at label %q (lines %d-%d) not found",
				label, expectedRange.Start.Line, expectedRange.End.Line)
		}
	}

	return d
}
