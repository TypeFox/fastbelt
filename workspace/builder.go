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
)

// Builder is the interface for building workspace-related structures.
type Builder interface {
	// Build processes the provided documents through all build phases (parse, compute
	// symbol table, link). It should regularly check ctx for cancellation between phases.
	Build(ctx context.Context, docs []*core.Document) error
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

// BuildStepListener is called right after a document has completed a build step.
// If the listener returns an error, it will be logged but will not prevent other listeners from being called.
// The listener may be called while the document's write lock is held; it must not acquire any document locks.
type BuildStepListener func(ctx context.Context, doc *core.Document) error

type buildStepEntry struct {
	states   core.DocumentState
	listener BuildStepListener
}

// DefaultBuilder is the default implementation of Builder.
type DefaultBuilder struct {
	srv         WorkspaceSrvCont
	buildMu     sync.Mutex
	listeners   []buildStepEntry
	listenersMu sync.RWMutex
}

// NewDefaultBuilder creates a new default builder.
func NewDefaultBuilder(srv WorkspaceSrvCont) Builder {
	return &DefaultBuilder{
		srv: srv,
	}
}

// Build processes the provided documents through all build phases.
func (b *DefaultBuilder) Build(ctx context.Context, docs []*core.Document) error {
	b.buildMu.Lock()
	defer b.buildMu.Unlock()

	// PHASE 1: Lock each document, parse, and compute exports (parallel per document).
	parser := b.srv.Workspace().DocumentParser
	exportedSymbols := b.srv.Linking().ExportedSymbolsProvider
	var phase1 sync.WaitGroup
	for _, doc := range docs {
		phase1.Go(func() {
			doc.Lock()
			if ctx.Err() != nil {
				return
			}
			// STEP 1.1: Parse the document and create the AST.
			if !doc.State.Has(core.DocStateParsed) {
				parser.Parse(doc)
				doc.State = doc.State.With(core.DocStateParsed)
				b.notifyListeners(ctx, core.DocStateParsed, doc)
			}
			if ctx.Err() != nil {
				return
			}
			// STEP 1.2: Compute the exported symbols for cross-document references.
			if !doc.State.Has(core.DocStateExportedSymbols) {
				exportedSymbols.Provide(ctx, doc)
				doc.State = doc.State.With(core.DocStateExportedSymbols)
				b.notifyListeners(ctx, core.DocStateExportedSymbols, doc)
			}
		})
	}
	phase1.Wait()

	if err := ctx.Err(); err != nil {
		for _, doc := range docs {
			doc.Unlock()
		}
		return err
	}

	// PHASE 2: Compute imported/local symbols and link (parallel per document).
	// This requires the exported symbols of all documents to be available.
	importedSymbols := b.srv.Linking().ImportedSymbolsProvider
	localSymbols := b.srv.Linking().LocalSymbolsProvider
	linker := b.srv.Linking().Linker
	referenceDescriptions := b.srv.Linking().ReferenceDescriptionsProvider
	var phase2 sync.WaitGroup
	for _, doc := range docs {
		phase2.Go(func() {
			if ctx.Err() != nil {
				return
			}
			// STEP 2.1: Collect imported symbols from all other documents.
			if !doc.State.Has(core.DocStateImportedSymbols) {
				allDocs := b.srv.Workspace().DocumentManager.All()
				importedSymbols.Provide(ctx, doc, allDocs)
				doc.State = doc.State.With(core.DocStateImportedSymbols)
				b.notifyListeners(ctx, core.DocStateImportedSymbols, doc)
			}
			if ctx.Err() != nil {
				return
			}
			// STEP 2.2: Compute the local symbols for intra-document references.
			if !doc.State.Has(core.DocStateLocalSymbols) {
				localSymbols.Provide(ctx, doc)
				doc.State = doc.State.With(core.DocStateLocalSymbols)
				b.notifyListeners(ctx, core.DocStateLocalSymbols, doc)
			}
			if ctx.Err() != nil {
				return
			}
			// STEP 2.3: Link the document to resolve all references.
			if !doc.State.Has(core.DocStateLinked) {
				linker.Link(ctx, doc)
				doc.State = doc.State.With(core.DocStateLinked)
				b.notifyListeners(ctx, core.DocStateLinked, doc)
			}
			if ctx.Err() != nil {
				return
			}
			// STEP 2.4: Provide reference descriptions for the document.
			if !doc.State.Has(core.DocStateReferences) {
				referenceDescriptions.Provide(ctx, doc)
				doc.State = doc.State.With(core.DocStateReferences)
				b.notifyListeners(ctx, core.DocStateReferences, doc)
			}
		})
	}
	phase2.Wait()

	for _, doc := range docs {
		doc.Unlock()
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	// PHASE 3: Run custom validations (parallel per document).
	validator := b.srv.Workspace().DocumentValidator
	var phase3 sync.WaitGroup
	for _, doc := range docs {
		phase3.Go(func() {
			if ctx.Err() != nil {
				return
			}
			doc.RLock()
			if !doc.State.Has(core.DocStateValidated) {
				diagnostics := validator.Validate(ctx, doc, "on-save")
				doc.RUnlock()
				if ctx.Err() != nil {
					return
				}
				doc.Lock()
				doc.Diagnostics = diagnostics
				doc.State = doc.State.With(core.DocStateValidated)
				doc.Unlock()
				b.notifyListeners(ctx, core.DocStateValidated, doc)
			} else {
				doc.RUnlock()
			}
		})
	}
	phase3.Wait()

	return ctx.Err()
}

// Reset selectively clears build results of a document. See [Builder.Reset].
func (b *DefaultBuilder) Reset(doc *core.Document, state core.DocumentState) {
	doc.Lock()
	defer doc.Unlock()
	if !state.Has(core.DocStateParsed) {
		doc.Root = nil
		doc.Tokens = core.TokenSlice{}
		doc.ParserErrors = []*core.ParserError{}
		doc.LexerErrors = []*core.LexerError{}
	}
	if !state.Has(core.DocStateExportedSymbols) {
		doc.ExportedSymbols = nil
	}
	if !state.Has(core.DocStateImportedSymbols) {
		doc.ImportedSymbols = nil
	}
	if !state.Has(core.DocStateLocalSymbols) {
		doc.LocalSymbols = nil
	}
	if !state.Has(core.DocStateLinked) {
		for _, ref := range doc.References {
			ref.Reset()
		}
		doc.References = []core.UntypedReference{}
	}
	if !state.Has(core.DocStateReferences) {
		doc.ReferenceDescriptions = nil
	}
	if !state.Has(core.DocStateValidated) {
		doc.Diagnostics = []*core.Diagnostic{}
	}
	doc.State = doc.State & state
}

// AddBuildStepListener registers a listener to be called after the specified build steps.
func (b *DefaultBuilder) AddBuildStepListener(states core.DocumentState, listener BuildStepListener) {
	if listener == nil {
		return
	}
	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()
	b.listeners = append(b.listeners, buildStepEntry{states: states, listener: listener})
}

// RemoveBuildStepListener unregisters a previously registered listener from all steps.
func (b *DefaultBuilder) RemoveBuildStepListener(listener BuildStepListener) {
	if listener == nil {
		return
	}
	b.listenersMu.Lock()
	defer b.listenersMu.Unlock()
	listenerPtr := reflect.ValueOf(listener).Pointer()
	for i, entry := range b.listeners {
		if reflect.ValueOf(entry.listener).Pointer() == listenerPtr {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			return
		}
	}
}

func (b *DefaultBuilder) notifyListeners(ctx context.Context, state core.DocumentState, doc *core.Document) {
	b.listenersMu.RLock()
	var matched []BuildStepListener
	for _, entry := range b.listeners {
		if entry.states.Has(state) {
			matched = append(matched, entry.listener)
		}
	}
	b.listenersMu.RUnlock()

	for _, listener := range matched {
		if ctx.Err() != nil {
			return
		}
		if err := listener(ctx, doc); err != nil {
			log.Printf("build step listener error (%s): %v", state, err)
		}
	}
}
