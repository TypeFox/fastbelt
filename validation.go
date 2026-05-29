// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import "context"

// DiagnosticSeverity mirrors LSP DiagnosticSeverity values.
type DiagnosticSeverity int

// Diagnostic severity values used by [Diagnostic].
const (
	SeverityError   DiagnosticSeverity = 1
	SeverityWarning DiagnosticSeverity = 2
	SeverityInfo    DiagnosticSeverity = 3
	SeverityHint    DiagnosticSeverity = 4
)

// String returns the human-readable label of s.
func (s DiagnosticSeverity) String() string {
	switch s {
	case SeverityError:
		return "Error"
	case SeverityWarning:
		return "Warning"
	case SeverityInfo:
		return "Info"
	case SeverityHint:
		return "Hint"
	default:
		return "Unknown"
	}
}

// DiagnosticTag mirrors LSP DiagnosticTag values.
type DiagnosticTag int

// Diagnostic tag values used by [Diagnostic].
const (
	TagUnnecessary DiagnosticTag = 1
	TagDeprecated  DiagnosticTag = 2
)

// Diagnostic represents a diagnostic message such as an error or warning.
// The struct mirrors lsp.Diagnostic so the core package stays free of that dependency.
type Diagnostic struct {
	// Range is the source range where the diagnostic applies.
	Range TextRange
	// Severity is the diagnostic severity level; when unset, clients may use a default severity.
	Severity DiagnosticSeverity
	// Message is the primary human-readable diagnostic message.
	Message string
	// Source identifies the diagnostic source, such as a tool or check name.
	Source string
	// Code is an optional diagnostic identifier shown to users and used for filtering.
	Code string
	// CodeDescription provides additional metadata for Code, typically a documentation link.
	CodeDescription *DiagnosticCodeDescription
	// Tags carry additional semantic classification, such as unnecessary or deprecated code.
	Tags []DiagnosticTag
	// Data is preserved between publishDiagnostics and codeAction requests.
	Data any
}

// A DiagnosticCodeDescription provides additional metadata for a diagnostic code.
//
// It mirrors lsp.CodeDescription and currently only carries a documentation URL.
type DiagnosticCodeDescription struct {
	Href string
}

// ValidationAcceptor is a callback that collects diagnostics reported during validation.
type ValidationAcceptor func(diagnostic *Diagnostic)

// Validator can be implemented by AST node Impl structs to provide custom validation checks.
type Validator interface {
	// Validate performs validation on the receiver node.
	// The level parameter identifies when validation runs (e.g. "on-type", "on-save").
	// The accept callback is used to collect diagnostics.
	Validate(ctx context.Context, level string, accept ValidationAcceptor)
}

// DiagnosticOption configures optional fields of a [Diagnostic] created by [NewDiagnostic].
type DiagnosticOption func(d *Diagnostic)

// NewDiagnostic creates a [Diagnostic] anchored to the given node's text range.
func NewDiagnostic(severity DiagnosticSeverity, message string, node AstNode, opts ...DiagnosticOption) *Diagnostic {
	d := &Diagnostic{
		Severity: severity,
		Message:  message,
	}
	if seg := node.Segment(); seg != nil {
		d.Range = seg.Range
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// WithToken narrows the diagnostic range to the given token's text segment.
// NOTE: These options might clash with other options in this package. If that happens,
// we can either rename them to DiagnosticToken etc. or move them to a separate package.
func WithToken(token *Token) DiagnosticOption {
	return func(d *Diagnostic) {
		if token != nil {
			d.Range = token.TextSegment.Range
		}
	}
}

// WithStringUnit narrows the diagnostic range to the given [StringUnit], if available.
func WithStringUnit(unit StringUnit) DiagnosticOption {
	return func(d *Diagnostic) {
		if unit != nil {
			if seg := unit.Segment(); seg != nil {
				d.Range = seg.Range
			}
		}
	}
}

// WithReference narrows the diagnostic range to the given reference text, if available.
func WithReference(ref UntypedReference) DiagnosticOption {
	return func(d *Diagnostic) {
		if ref != nil {
			if seg := ref.Segment(); seg != nil {
				d.Range = seg.Range
			}
		}
	}
}

// WithRange sets an explicit range on the diagnostic, overriding any node or token range.
func WithRange(r TextRange) DiagnosticOption {
	return func(d *Diagnostic) {
		d.Range = r
	}
}

// WithCode sets the diagnostic code.
func WithCode(code string) DiagnosticOption {
	return func(d *Diagnostic) {
		d.Code = code
	}
}

// WithTags sets diagnostic tags (e.g. [TagUnnecessary], [TagDeprecated]).
func WithTags(tags ...DiagnosticTag) DiagnosticOption {
	return func(d *Diagnostic) {
		d.Tags = tags
	}
}

// WithData attaches arbitrary data to the diagnostic.
func WithData(data any) DiagnosticOption {
	return func(d *Diagnostic) {
		d.Data = data
	}
}

// WithCodeDescription sets additional metadata for the diagnostic code.
//
// The current metadata format mirrors LSP and supports a documentation URL.
func WithCodeDescription(codeDescription *DiagnosticCodeDescription) DiagnosticOption {
	return func(d *Diagnostic) {
		d.CodeDescription = codeDescription
	}
}

// WithCodeDescriptionHref sets a hyperlink for the error code description.
func WithCodeDescriptionHref(href string) DiagnosticOption {
	return func(d *Diagnostic) {
		d.CodeDescription = &DiagnosticCodeDescription{Href: href}
	}
}
