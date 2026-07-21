// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"log"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// DiagnosticsPublisher is responsible for publishing diagnostics to the LSP client.
// It publishes diagnostics in the following scenarios:
//
//   - A document is opened (reuses diagnostics from build).
//   - A document is closed (clears diagnostics).
//   - A document is validated by the builder.
//
// When registered to a service container, it automatically hooks into the
// document lifecycle events during server initialization.
type DiagnosticsPublisher struct {
	sc *service.Container
}

// NewDiagnosticsPublisher creates a new instance of [DiagnosticsPublisher].
func NewDiagnosticsPublisher(sc *service.Container) *DiagnosticsPublisher {
	return &DiagnosticsPublisher{sc: sc}
}

// OnServerInitialize implements [InitializeParticipant] by subscribing to the
// document lifecycle events that trigger publishing diagnostics.
func (d *DiagnosticsPublisher) OnServerInitialize(_ *lsp.ParamInitialize) {
	store, err := service.Get[textdoc.Store](d.sc)
	if err != nil {
		return
	}
	syncher, err := service.Get[DocumentSyncher](d.sc)
	if err != nil {
		return
	}
	syncher.OnDidOpen(func(ctx context.Context, event *TextDocumentChangeEvent) {
		docManager, err := service.Get[workspace.DocumentManager](d.sc)
		if err != nil {
			return
		}
		uri := core.ParseURI(string(event.Document.URI()))
		document := docManager.Get(uri)
		if document != nil {
			// If the document already exists, publish diagnostics immediately.
			d.publishDocumentDiagnostics(ctx, document)
		}
	})
	syncher.OnDidClose(func(ctx context.Context, event *TextDocumentChangeEvent) {
		uri := core.NormalizeURI(event.Document.URI())
		// Clear diagnostics for the closed document.
		d.publishDiagnostics(ctx, uri, []lsp.Diagnostic{})
	})
	builder, err := service.Get[workspace.Builder](d.sc)
	if err != nil {
		return
	}
	builder.AddBuildStepListener(core.DocStateValidated, func(ctx context.Context, doc *core.Document) error {
		if store.GetOverlay(doc.URI.DocumentURI()) == nil {
			return nil // Document is not open, skip publishing diagnostics
		}
		d.publishDocumentDiagnostics(ctx, doc)
		return nil
	})
}

func (d *DiagnosticsPublisher) publishDocumentDiagnostics(ctx context.Context, doc *core.Document) {
	lspDiags := make([]lsp.Diagnostic, 0, len(doc.Diagnostics))
	for _, d := range doc.Diagnostics {
		lspDiags = append(lspDiags, toLspDiagnostic(*d))
	}
	d.publishDiagnostics(ctx, doc.TextDoc.URI(), lspDiags)
}

func (d *DiagnosticsPublisher) publishDiagnostics(ctx context.Context, uri lsp.DocumentURI, diagnostics []lsp.Diagnostic) {
	connection, err := service.Get[*Connection](d.sc)
	if err != nil {
		return
	}
	client := lsp.ClientDispatcher(connection.Value)
	err = client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
	if err != nil {
		log.Printf("failed to publish diagnostics: %v", err)
	}
}
