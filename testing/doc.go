package fbtest

import (
	"context"
	"strings"
	"testing"

	core "typefox.dev/fastbelt"
)

// Doc wraps a built [core.Document] with assertion methods and marker positions.
// Access the underlying document via the exported Document field.
type Doc struct {
	Document *core.Document
	Ranges   []RangeMarker
	Indices  []IndexMarker
	ctx      context.Context
	t        testing.TB
}

// Ctx returns the context for this document. Pass it to [core.Reference.Ref].
func (d *Doc) Ctx() context.Context { return d.ctx }

// Root returns the root AST node of the document.
func (d *Doc) Root() core.AstNode { return d.Document.Root }

// Diagnostics returns all diagnostics collected during validation.
func (d *Doc) Diagnostics() []*core.Diagnostic { return d.Document.Diagnostics }

// AssertNoErrors fails the test if any LexerErrors, ParserErrors, or error-severity
// Diagnostics exist.
func (d *Doc) AssertNoErrors() *Doc {
	d.t.Helper()
	for _, e := range d.Document.LexerErrors {
		d.t.Errorf("fbtest: unexpected lexer error: %v", e)
	}
	for _, e := range d.Document.ParserErrors {
		d.t.Errorf("fbtest: unexpected parser error: %v", e)
	}
	for _, diag := range d.Document.Diagnostics {
		if diag.Severity == core.SeverityError {
			d.t.Errorf("fbtest: unexpected error diagnostic: %s", diag.Message)
		}
	}
	return d
}

// AssertNoParseErrors fails the test if any LexerErrors or ParserErrors exist.
// Diagnostics are not checked.
func (d *Doc) AssertNoParseErrors() *Doc {
	d.t.Helper()
	for _, e := range d.Document.LexerErrors {
		d.t.Errorf("fbtest: unexpected lexer error: %v", e)
	}
	for _, e := range d.Document.ParserErrors {
		d.t.Errorf("fbtest: unexpected parser error: %v", e)
	}
	return d
}

// AssertDiagnostic fails the test unless at least one diagnostic has the given
// severity and a message containing msgSubstring.
func (d *Doc) AssertDiagnostic(severity core.DiagnosticSeverity, msgSubstring string) *Doc {
	d.t.Helper()
	for _, diag := range d.Document.Diagnostics {
		if diag.Severity == severity && strings.Contains(diag.Message, msgSubstring) {
			return d
		}
	}
	d.t.Errorf("fbtest: no %s diagnostic containing %q (got %d total)",
		severity, msgSubstring, len(d.Document.Diagnostics))
	return d
}

// AssertDiagnostic fails the test unless at least one diagnostic has the given
// severity and code.
func (d *Doc) AssertDiagnosticCode(severity core.DiagnosticSeverity, code string) *Doc {
	d.t.Helper()
	for _, diag := range d.Document.Diagnostics {
		if diag.Severity == severity && diag.Code == code {
			return d
		}
	}
	d.t.Errorf("fbtest: no %s diagnostic with code %q (got %d total)",
		severity, code, len(d.Document.Diagnostics))
	return d
}

// AssertDiagnosticAtLabel fails the test unless at least one diagnostic has the given
// severity, a message containing msgSubstring, and a range that contains the position
// of the named marker.
func (d *Doc) AssertDiagnosticAtLabel(label string, severity core.DiagnosticSeverity, msgSubstring string) *Doc {
	d.t.Helper()
	loc, ok := d.markerRange(label)
	if !ok {
		d.t.Errorf("fbtest: no marker with label %q", label)
		return d
	}
	for _, diag := range d.Document.Diagnostics {
		if diag.Severity == severity && strings.Contains(diag.Message, msgSubstring) && loc == diag.Range {
			return d
		}
	}
	d.t.Errorf("fbtest: no %s diagnostic containing %q at label %q", severity, msgSubstring, label)
	return d
}

// AssertNoDiagnostics fails the test if any diagnostics exist at any severity.
func (d *Doc) AssertNoDiagnostics() *Doc {
	d.t.Helper()
	for _, diag := range d.Document.Diagnostics {
		d.t.Errorf("fbtest: unexpected diagnostic [%s]: %s", diag.Severity, diag.Message)
	}
	return d
}

// AssertState fails the test unless the document state includes the given flag.
func (d *Doc) AssertState(flag core.DocumentState) *Doc {
	d.t.Helper()
	if !d.Document.State.Has(flag) {
		d.t.Errorf("fbtest: document state does not include %v (actual: %v)", flag, d.Document.State)
	}
	return d
}

func (d *Doc) markerRange(label string) (core.TextRange, bool) {
	for _, r := range d.Ranges {
		if r.Label == label {
			start := d.Document.TextDoc.PositionAt(r.Start)
			end := d.Document.TextDoc.PositionAt(r.End)
			return core.TextRange{Start: core.TextLocation{
				Line:   core.TextLine(start.Line),
				Column: core.TextColumn(start.Character),
			}, End: core.TextLocation{
				Line:   core.TextLine(end.Line),
				Column: core.TextColumn(end.Character),
			}}, true
		}
	}
	return core.TextRange{}, false
}

// markerLocation returns the TextLocation of the named marker.
// For range markers it returns the start; for index markers the offset position.
func (d *Doc) markerLocation(label string) (core.TextLocation, bool) {
	text := d.Document.TextDoc.Text(nil)
	for _, r := range d.Ranges {
		if r.Label == label {
			return offsetToLocation(text, r.Start), true
		}
	}
	for _, idx := range d.Indices {
		if idx.Label == label {
			return offsetToLocation(text, idx.Offset), true
		}
	}
	return core.TextLocation{}, false
}

// FindNode returns the first node of type T in the document tree, or (zero, false).
func FindNode[T core.AstNode](d *Doc) (T, bool) {
	d.t.Helper()
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
	d.t.Helper()
	n, ok := FindNode[T](d)
	if !ok {
		d.t.Fatal("fbtest: MustFindNode: no node of the requested type found")
	}
	return n
}

// FindAll returns all nodes of type T in the document tree.
func FindAll[T core.AstNode](d *Doc) []T {
	d.t.Helper()
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
	d.t.Helper()
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
	d.t.Helper()
	n, ok := FindNodeAtOffset[T](d, offset)
	if !ok {
		d.t.Fatalf("fbtest: MustFindNodeAtOffset: no node of the requested type at offset %d", offset)
	}
	return n
}

// FindNodeAtLocation returns the most specific (smallest span) node of type T whose
// text range contains the given line/column position.
func FindNodeAtLocation[T core.AstNode](d *Doc, location core.TextLocation) (T, bool) {
	d.t.Helper()
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
	d.t.Helper()
	n, ok := FindNodeAtLocation[T](d, location)
	if !ok {
		d.t.Fatalf("fbtest: MustFindNodeAtLocation: no node of the requested type at %v", location)
	}
	return n
}

// FindNodeAtLabel returns the most specific node of type T at the position of the
// named range marker (start offset) or index marker.
func FindNodeAtLabel[T core.AstNode](d *Doc, label string) (T, bool) {
	d.t.Helper()
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
	d.t.Helper()
	n, ok := FindNodeAtLabel[T](d, label)
	if !ok {
		d.t.Fatalf("fbtest: MustFindNodeAtLabel: no node of the requested type at label %q", label)
	}
	return n
}

// FindNamedNode returns the first node of type T whose Name() method returns name.
// T must implement [core.NamedNode] (or have a Name() string method) to find results.
func FindNamedNode[T core.AstNode](d *Doc, name string) (T, bool) {
	d.t.Helper()
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
	d.t.Helper()
	n, ok := FindNamedNode[T](d, name)
	if !ok {
		d.t.Fatalf("fbtest: MustFindNamedNode: no node of the requested type with name %q", name)
	}
	return n
}

// FindReferenceWithText returns the first reference of type T whose cross-reference
// text matches the given string exactly.
func FindReferenceWithText[T core.AstNode](d *Doc, text string) (*core.Reference[T], bool) {
	d.t.Helper()
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
	d.t.Helper()
	r, ok := FindReferenceWithText[T](d, text)
	if !ok {
		d.t.Fatalf("fbtest: MustFindReferenceWithText: no reference of the requested type with text %q", text)
	}
	return r
}

// FindReference returns the first reference of type T whose cross-reference text is located at the given label.
func FindReference[T core.AstNode](d *Doc, label string) (*core.Reference[T], bool) {
	d.t.Helper()
	targetRange, found := d.markerRange(label)
	if !found || d.Document.Root == nil {
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
	d.t.Helper()
	r, ok := FindReference[T](d, label)
	if !ok {
		d.t.Fatalf("fbtest: MustFindReference: no reference of the requested type at label %q", label)
	}
	return r
}
