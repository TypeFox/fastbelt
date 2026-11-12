// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"log"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// DiagnosticsPublisher is a service for pushing diagnostics to the LSP client.
//
// The language server does not call it itself: an implementation decides when
// to publish, typically by implementing [InitializeParticipant] and
// subscribing to document lifecycle events, as [DefaultDiagnosticsPublisher]
// does. Other services can call it to publish or clear diagnostics on demand.
type DiagnosticsPublisher interface {
	// PublishDiagnostics sends the current [core.Document.Diagnostics] of doc
	// to the client, replacing the diagnostics previously published for it.
	PublishDiagnostics(ctx context.Context, doc *core.Document)
	// ClearDiagnostics removes all diagnostics the client shows for uri.
	ClearDiagnostics(ctx context.Context, uri core.URI)
}

// DiagnosticsPublisherOptions configures a [DefaultDiagnosticsPublisher].
type DiagnosticsPublisherOptions struct {
	// AllDocuments makes the publisher report diagnostics for every validated
	// document in the workspace. By default, only documents that are open in
	// the editor get diagnostics.
	AllDocuments bool
}

// DefaultDiagnosticsPublisher is the default implementation of
// [DiagnosticsPublisher]. During server initialization it hooks into the
// document lifecycle and publishes diagnostics in the following scenarios:
//
//   - A document is opened (reuses diagnostics from build).
//   - A document is closed (clears diagnostics).
//   - A document is validated by the builder.
//
// By default, only documents that are open in the editor are covered. With
// [DiagnosticsPublisherOptions.AllDocuments], every validated document is
// published, and closing a document that is still part of the workspace
// keeps its diagnostics. Diagnostics of a document that is removed from the
// workspace are not cleared automatically; call
// [DiagnosticsPublisher.ClearDiagnostics] where documents are removed.
type DefaultDiagnosticsPublisher struct {
	sc      *service.Container
	options DiagnosticsPublisherOptions
}

// Ensure DefaultDiagnosticsPublisher implements the expected interfaces.
var _ DiagnosticsPublisher = (*DefaultDiagnosticsPublisher)(nil)
var _ InitializeParticipant = (*DefaultDiagnosticsPublisher)(nil)

// NewDefaultDiagnosticsPublisher creates a [DefaultDiagnosticsPublisher] that
// publishes diagnostics for open documents only.
func NewDefaultDiagnosticsPublisher(sc *service.Container) DiagnosticsPublisher {
	return NewDiagnosticsPublisherWithOptions(sc, DiagnosticsPublisherOptions{})
}

// NewDiagnosticsPublisherWithOptions creates a [DefaultDiagnosticsPublisher]
// with custom [DiagnosticsPublisherOptions].
func NewDiagnosticsPublisherWithOptions(sc *service.Container, options DiagnosticsPublisherOptions) DiagnosticsPublisher {
	return &DefaultDiagnosticsPublisher{sc: sc, options: options}
}

// NoopDiagnosticsPublisher is a [DiagnosticsPublisher] that does nothing.
// Register it to turn off pushed diagnostics, for example when diagnostics
// are delivered in another way.
type NoopDiagnosticsPublisher struct{}

// NewNoopDiagnosticsPublisher creates a new instance of [NoopDiagnosticsPublisher].
func NewNoopDiagnosticsPublisher() DiagnosticsPublisher {
	return NoopDiagnosticsPublisher{}
}

// PublishDiagnostics does nothing.
func (NoopDiagnosticsPublisher) PublishDiagnostics(context.Context, *core.Document) {}

// ClearDiagnostics does nothing.
func (NoopDiagnosticsPublisher) ClearDiagnostics(context.Context, core.URI) {}

// OnServerInitialize implements [InitializeParticipant] by subscribing to the
// document lifecycle events that trigger publishing diagnostics.
func (d *DefaultDiagnosticsPublisher) OnServerInitialize(_ *lsp.ParamInitialize) {
	store, err := d.sc.Get[textdoc.Store]()
	if err != nil {
		return
	}
	synchronizer, err := d.sc.Get[DocumentSynchronizer]()
	if err != nil {
		return
	}
	synchronizer.OnDidOpen(func(ctx context.Context, event *TextDocumentChangeEvent) {
		docManager, err := d.sc.Get[workspace.DocumentManager]()
		if err != nil {
			return
		}
		uri := core.ParseURI(string(event.Document.URI()))
		document := docManager.Get(uri)
		if document != nil && document.State().Has(core.DocStateValidated) &&
			bytes.Equal(document.TextDoc.Content(), event.Document.Content()) {
			// If the user opened a document that has already gone through the workspace build
			// we can immediately send out the diagnostics.
			// In all other cases, a build will be triggered, which will publish its own diagnostics.
			d.PublishDiagnostics(ctx, document)
		}
	})
	synchronizer.OnDidClose(func(ctx context.Context, event *TextDocumentChangeEvent) {
		uri := core.ParseURI(string(event.Document.URI()))
		if d.options.AllDocuments {
			if docManager, err := d.sc.Get[workspace.DocumentManager](); err == nil && docManager.Has(uri) {
				// The document is still part of the workspace, keep its diagnostics.
				return
			}
		}
		// Clear diagnostics for the closed document.
		d.ClearDiagnostics(ctx, uri)
	})
	builder, err := d.sc.Get[workspace.Builder]()
	if err != nil {
		return
	}
	builder.AddBuildStepListener(core.DocStateValidated, func(ctx context.Context, doc *core.Document) error {
		if !d.options.AllDocuments && store.GetOverlay(doc.URI.DocumentURI()) == nil {
			return nil // Document is not open, skip publishing diagnostics
		}
		d.PublishDiagnostics(ctx, doc)
		return nil
	})
}

func (d *DefaultDiagnosticsPublisher) PublishDiagnostics(ctx context.Context, doc *core.Document) {
	lspDiags := make([]lsp.Diagnostic, 0, len(doc.Diagnostics))
	for _, d := range doc.Diagnostics {
		lspDiags = append(lspDiags, toLspDiagnostic(doc.TextDoc, d))
	}
	d.publishDiagnostics(ctx, doc.TextDoc.URI(), lspDiags)
}

func (d *DefaultDiagnosticsPublisher) ClearDiagnostics(ctx context.Context, uri core.URI) {
	d.publishDiagnostics(ctx, uri.DocumentURI(), []lsp.Diagnostic{})
}

func toLspDiagnostic(textDoc textdoc.Handle, d *core.Diagnostic) lsp.Diagnostic {
	result := lsp.Diagnostic{
		Range:    d.Range.LspRange(textDoc),
		Severity: lsp.DiagnosticSeverity(d.Severity),
		Message:  lsp.DiagnosticMessage{String: &d.Message},
	}
	if d.Source != "" {
		result.Source = d.Source
	}
	if d.Code != "" {
		result.Code = lsp.DiagnosticCode{String: &d.Code}
	}
	if d.CodeDescription != nil {
		result.CodeDescription = &lsp.CodeDescription{
			Href: d.CodeDescription.Href,
		}
	}
	if len(d.Tags) > 0 {
		tags := make([]lsp.DiagnosticTag, len(d.Tags))
		for i, t := range d.Tags {
			tags[i] = lsp.DiagnosticTag(t)
		}
		result.Tags = tags
	}
	if d.Data != nil {
		raw, err := json.Marshal(d.Data)
		if err == nil {
			rawMsg := jsontext.Value(raw)
			result.Data = &rawMsg
		}
	}
	return result
}

func (d *DefaultDiagnosticsPublisher) publishDiagnostics(ctx context.Context, uri lsp.DocumentURI, diagnostics []lsp.Diagnostic) {
	connection, err := d.sc.Get[*Connection]()
	if err != nil {
		return
	}
	client := lsp.ClientDispatcher(connection.Value)
	// Never cancel publishing, even if the build that computed the diagnostics
	// is already cancelled: nothing publishes them again unless the document changes.
	ctx = context.WithoutCancel(ctx)
	err = client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
	if err != nil {
		log.Printf("failed to publish diagnostics: %v", err)
	}
}
