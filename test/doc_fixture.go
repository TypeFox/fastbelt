package test

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/util/service"
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
			d.fixture.t.Fatalf("fbtest: unexpected error diagnostic: %s", diag.Message)
		}
	}
	return d
}

// AssertNoParseErrors fails the test if any LexerErrors or ParserErrors exist.
// Diagnostics are not checked.
func (d *Doc) AssertNoParseErrors() *Doc {
	d.fixture.t.Helper()
	for _, e := range d.Document.LexerErrors {
		d.fixture.t.Fatalf("fbtest: unexpected lexer error: %v", e)
	}
	for _, e := range d.Document.ParserErrors {
		d.fixture.t.Fatalf("fbtest: unexpected parser error: %v", e)
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
	loc, err := d.MarkerRange(label)
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

func (d *Doc) MarkerRange(label string) (core.TextRange, error) {
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
	targetRange, err := d.MarkerRange(label)
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

// FindSymbolAtLabel finds a document symbol at the given marker label.
// Returns the symbol and true if found, or nil and false if not found.
func (d *Doc) FindSymbolAtLabel(label string) (*lsp.DocumentSymbol, bool) {
	d.fixture.t.Helper()

	provider, err := service.Get[server.DocumentSymbolProvider](d.fixture.sc)
	if err != nil {
		return nil, false
	}

	params := &lsp.DocumentSymbolParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: lsp.DocumentURI(d.Document.URI.DocumentURI()),
		},
	}

	result, err := provider.HandleDocumentSymbolRequest(d.fixture.ctx, params)
	if err != nil {
		return nil, false
	}

	expectedRange, err := d.MarkerRange(label)
	if err != nil {
		return nil, false
	}

	sym := findSymbolAtRange(result, expectedRange)
	if sym == nil {
		return nil, false
	}
	return sym, true
}

// MustFindSymbolAtLabel finds a document symbol at the given marker label.
// Fails the test if the symbol is not found.
func (d *Doc) MustFindSymbolAtLabel(label string) *lsp.DocumentSymbol {
	d.fixture.t.Helper()
	sym, ok := d.FindSymbolAtLabel(label)
	if !ok {
		d.fixture.t.Fatalf("fbtest: MustFindSymbolAtLabel: no symbol found at label %q", label)
	}
	return sym
}

// AssertDocumentSymbol verifies that a document symbol exists with the given name at the marker label.
// Returns the Doc for chaining.
func (d *Doc) AssertDocumentSymbol(label string, expectedName string, expectedKind lsp.SymbolKind) *Doc {
	d.fixture.t.Helper()

	found := d.MustFindSymbolAtLabel(label)

	if found.Name != expectedName {
		d.fixture.t.Errorf("fbtest: symbol at %q has name %q, expected %q",
			label, found.Name, expectedName)
	}
	if found.Kind != expectedKind {
		d.fixture.t.Errorf("fbtest: symbol at %q has kind %v, expected %v",
			label, found.Kind, expectedKind)
	}

	return d
}

// findSymbolAtRange is a helper function for recursive symbol search.
func findSymbolAtRange(symbols []lsp.DocumentSymbol, targetRange core.TextRange) *lsp.DocumentSymbol {
	for i := range symbols {
		sym := &symbols[i]
		symRange := core.TextRange{
			Start: core.TextLocation{
				Line:   core.TextLine(sym.Range.Start.Line),
				Column: core.TextColumn(sym.Range.Start.Character),
			},
			End: core.TextLocation{
				Line:   core.TextLine(sym.Range.End.Line),
				Column: core.TextColumn(sym.Range.End.Character),
			},
		}

		if symRange == targetRange {
			return sym
		}

		if found := findSymbolAtRange(sym.Children, targetRange); found != nil {
			return found
		}
	}
	return nil
}
