// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"
	"reflect"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/parallel"
	"typefox.dev/fastbelt/util/service"
)

// Builder runs the workspace document build pipeline: parse, symbol export/import,
// linking, reference indexing, and validation.
type Builder interface {
	// Build processes the provided documents through all build phases (parse, compute
	// symbol table, link, validate). It should regularly check ctx for cancellation
	// between phases. downgrade must be called by the implementation once phases 1 and 2
	// (the write phase) are complete; this transitions the workspace lock to readable so
	// that read requests can proceed while phase 3 (validation) runs.
	Build(ctx context.Context, docs []*core.Document, downgrade func()) error
	// Reset selectively clears build results of a document. The state parameter is a
	// bitmask of states to keep; for every bit that is not set, the corresponding document
	// fields are reset to their initial values and the bit is cleared from doc.State.
	Reset(doc *core.Document, state core.DocumentState)
	// AddBuildStepListener registers a listener to be called after documents complete the
	// specified build steps. The states parameter is a bitmask, so multiple steps can be
	// selected with e.g. DocStateParsed | DocStateLinked.
	AddBuildStepListener(states core.DocumentState, listener BuildStepListener)
	// RemoveBuildStepListener unregisters a previously registered listener from all steps.
	RemoveBuildStepListener(listener BuildStepListener)
}

// BuildStepListener is called right after a document completes a build step
// selected by [Builder.AddBuildStepListener]. If the listener returns an error,
// it is logged and does not stop other listeners or the build.
type BuildStepListener func(ctx context.Context, doc *core.Document) error

type buildStepEntry struct {
	states   core.DocumentState
	listener BuildStepListener
}

// DefaultBuilder is the default implementation of [Builder].
type DefaultBuilder struct {
	sc          *service.Container
	listeners   []buildStepEntry
	listenersMu sync.RWMutex
}

// NewDefaultBuilder returns a [Builder] that runs the standard three-phase
// build pipeline in parallel per document within each phase.
func NewDefaultBuilder(sc *service.Container) Builder {
	return &DefaultBuilder{sc: sc}
}

func (s *DefaultBuilder) Build(ctx context.Context, docs []*core.Document, downgrade func()) error {
	// PHASE 1: Parse, and compute exports (parallel per document).
	parser := service.MustGet[DocumentParser](s.sc)
	exporter := service.MustGet[linking.SymbolExporter](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		if ctx.Err() != nil {
			return
		}
		// STEP 1.1: Parse the document and create the AST.
		if !doc.State.Has(core.DocStateParsed) {
			parser.Parse(doc)
			doc.State = doc.State.With(core.DocStateParsed)
			s.notifyListeners(ctx, core.DocStateParsed, doc)
		}
		if ctx.Err() != nil {
			return
		}
		// STEP 1.2: Compute the exported symbols for cross-document references.
		if !doc.State.Has(core.DocStateExportedSymbols) {
			exporter.ExportSymbols(ctx, doc)
			doc.State = doc.State.With(core.DocStateExportedSymbols)
			s.notifyListeners(ctx, core.DocStateExportedSymbols, doc)
		}
	})

	if err := ctx.Err(); err != nil {
		return err
	}

	// PHASE 2: Compute imported/local symbols and link (parallel per document).
	// This requires the exported symbols of all documents to be available.
	documentManager := service.MustGet[DocumentManager](s.sc)
	importer := service.MustGet[linking.SymbolImporter](s.sc)
	localSymbols := service.MustGet[linking.LocalSymbolsProvider](s.sc)
	linker := service.MustGet[linking.Linker](s.sc)
	referenceDescriptions := service.MustGet[linking.ReferenceDescriptionsProvider](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		if ctx.Err() != nil {
			return
		}
		// STEP 2.1: Collect imported symbols from all other documents.
		if !doc.State.Has(core.DocStateImportedSymbols) {
			allDocs := documentManager.All()
			importer.ImportSymbols(ctx, doc, allDocs)
			doc.State = doc.State.With(core.DocStateImportedSymbols)
			s.notifyListeners(ctx, core.DocStateImportedSymbols, doc)
		}
		if ctx.Err() != nil {
			return
		}
		// STEP 2.2: Compute the local symbols for intra-document references.
		if !doc.State.Has(core.DocStateLocalSymbols) {
			localSymbols.LocalSymbols(ctx, doc)
			doc.State = doc.State.With(core.DocStateLocalSymbols)
			s.notifyListeners(ctx, core.DocStateLocalSymbols, doc)
		}
		if ctx.Err() != nil {
			return
		}
		// STEP 2.3: Link the document to resolve all references.
		if !doc.State.Has(core.DocStateLinked) {
			linker.Link(ctx, doc)
			doc.State = doc.State.With(core.DocStateLinked)
			s.notifyListeners(ctx, core.DocStateLinked, doc)
		}
		if ctx.Err() != nil {
			return
		}
		// STEP 2.4: Provide reference descriptions for the document.
		if !doc.State.Has(core.DocStateReferences) {
			referenceDescriptions.ReferenceDescriptions(ctx, doc)
			doc.State = doc.State.With(core.DocStateReferences)
			s.notifyListeners(ctx, core.DocStateReferences, doc)
		}
	})

	if err := ctx.Err(); err != nil {
		// Important note: Do not downgrade the lock here!
		// If we downgrade the lock here, we would allow read access to
		// the workspace while the documents are in an inconsistent state.
		// In most cases, the error has been triggered by a new change,
		// which will trigger a new build with a re-acquired read-lock.
		return err
	}

	// Transition from write phase to readable: releases the exclusive lock so
	// read requests can proceed while validation (phase 3) runs concurrently.
	if downgrade != nil {
		downgrade()
	}

	// PHASE 3: Run custom validations (parallel per document).
	validator := service.MustGet[DocumentValidator](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		if ctx.Err() != nil {
			return
		}
		if !doc.State.Has(core.DocStateValidated) {
			diagnostics := validator.Validate(ctx, doc, "on-save")
			if ctx.Err() != nil {
				return
			}
			doc.Diagnostics = diagnostics
			doc.State = doc.State.With(core.DocStateValidated)
			s.notifyListeners(ctx, core.DocStateValidated, doc)
		}
	})

	return ctx.Err()
}

func (s *DefaultBuilder) Reset(doc *core.Document, state core.DocumentState) {
	switch {
	case !state.Has(core.DocStateParsed):
		doc.Root = nil
		doc.Tokens = core.TokenSlice{}
		doc.ParserErrors = []*core.ParserError{}
		doc.LexerErrors = []*core.LexerError{}
		doc.References = []core.UntypedReference{}
		fallthrough
	case !state.Has(core.DocStateExportedSymbols):
		doc.ExportedSymbols = nil
		fallthrough
	case !state.Has(core.DocStateImportedSymbols):
		doc.ImportedSymbols = nil
		fallthrough
	case !state.Has(core.DocStateLocalSymbols):
		doc.LocalSymbols = nil
	case !state.Has(core.DocStateLinked):
		// Note: NOOP if the document should be completely reset,
		// because the references slice is already cleared above.
		// If not cleared, the references are reset to allow re-resolution
		// on the next build.
		for _, ref := range doc.References {
			ref.Reset()
		}
	case !state.Has(core.DocStateReferences):
		doc.ReferenceDescriptions = nil
	case !state.Has(core.DocStateValidated):
		doc.Diagnostics = []*core.Diagnostic{}
	}
	doc.State = doc.State & state
}

func (s *DefaultBuilder) AddBuildStepListener(states core.DocumentState, listener BuildStepListener) {
	if listener == nil {
		return
	}
	s.listenersMu.Lock()
	defer s.listenersMu.Unlock()
	s.listeners = append(s.listeners, buildStepEntry{states: states, listener: listener})
}

func (s *DefaultBuilder) RemoveBuildStepListener(listener BuildStepListener) {
	if listener == nil {
		return
	}
	s.listenersMu.Lock()
	defer s.listenersMu.Unlock()
	listenerPtr := reflect.ValueOf(listener).Pointer()
	for i, entry := range s.listeners {
		if reflect.ValueOf(entry.listener).Pointer() == listenerPtr {
			s.listeners = append(s.listeners[:i], s.listeners[i+1:]...)
			return
		}
	}
}

func (s *DefaultBuilder) notifyListeners(ctx context.Context, state core.DocumentState, doc *core.Document) {
	s.listenersMu.RLock()
	var matched []BuildStepListener
	for _, entry := range s.listeners {
		if entry.states.Has(state) {
			matched = append(matched, entry.listener)
		}
	}
	s.listenersMu.RUnlock()

	for _, listener := range matched {
		if ctx.Err() != nil {
			return
		}
		if err := listener(ctx, doc); err != nil {
			log.Printf("build step listener error (%s): %v", state, err)
		}
	}
}
