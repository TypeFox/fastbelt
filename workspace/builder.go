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
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/parallel"
	"typefox.dev/fastbelt/util/service"
)

// Builder runs the workspace document build pipeline: parse, symbol export/import,
// linking, reference indexing, and validation.
type Builder interface {
	// Build processes the provided documents through all build phases (parse, compute
	// symbol table, link, validate). It must be called under [Lock.Write] and should
	// regularly check ctx for cancellation between phases. Calling Build marks the
	// transition from the write's mutation phase to its build phase (see [Lock.Write]):
	// implementations signal [Lock.StateChanged] as document states advance so that
	// [Lock.ReadAt] requests are admitted as soon as the states they need are reached,
	// and report the workspace-wide floor after each phase barrier. The floor relies
	// on every document outside docs already being complete, which the caller must
	// ensure (the updater collects every incomplete document into the build set).
	Build(ctx context.Context, docs []*core.Document) error
	// Reset selectively clears build results of a document. The state parameter is a
	// bitmask of states to keep; for every bit that is not set, the corresponding document
	// fields are reset to their initial values and the bit is cleared from doc.State.
	// Reset must be called under [Lock.Write], during the mutation phase - that is,
	// before Build starts advancing document states.
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

// NewDefaultBuilder returns a [Builder] that runs the standard four-phase
// build pipeline in parallel per document within each phase.
func NewDefaultBuilder(sc *service.Container) Builder {
	return &DefaultBuilder{sc: sc}
}

func (s *DefaultBuilder) Build(ctx context.Context, docs []*core.Document) error {
	lock := service.MustGet[Lock](s.sc)
	// Build has started, signal the lock to be ready to admit readers after
	// advancing document states.
	lock.StateChanged(0)
	// Wakes up any ReadAt calls to make sure they can see the new document state.
	// Documents might reach the requested state before the full workspace does.
	advance := func(doc *core.Document, state core.DocumentState) {
		doc.AddState(state)
		lock.StateChanged(0)
		s.notifyListeners(ctx, state, doc)
	}

	// PHASE 1: Parse, and compute exports (parallel per document).
	lexer := service.MustGet[lexer.Lexer](s.sc)
	parser := service.MustGet[parser.Parser](s.sc)
	exporter := service.MustGet[linking.SymbolExporter](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		// STEP 1.1: Lex the document and produce tokens.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateLexed) {
			lexer.Exec(doc)
			if ctx.Err() != nil {
				s.resetLexed(doc)
				return
			}
			advance(doc, core.DocStateLexed)
		}
		// STEP 1.2: Parse the document and create the AST.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateParsed) {
			parser.Parse(doc)
			if ctx.Err() != nil {
				s.resetParsed(doc)
				return
			}
			advance(doc, core.DocStateParsed)
		}
		// STEP 1.3: Compute the exported symbols for cross-document references.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateExportedSymbols) {
			exporter.ExportSymbols(ctx, doc)
			if ctx.Err() != nil {
				s.resetExportedSymbols(doc)
				return
			}
			advance(doc, core.DocStateExportedSymbols)
		}
	})

	if err := ctx.Err(); err != nil {
		return err
	}
	// Phase 1 barrier: every document in docs is now parsed with exports
	// computed, and documents outside docs were already complete when the
	// updater collected the build set. Report the workspace floor so
	// workspace-wide ReadAt calls are admitted.
	lock.StateChanged(core.DocStateParsed | core.DocStateExportedSymbols)

	// PHASE 2: Compute imported and local symbols (parallel per document and step).
	// This requires the exported symbols of all documents to be available.
	documentManager := service.MustGet[DocumentManager](s.sc)
	importer := service.MustGet[linking.SymbolImporter](s.sc)
	localSymbols := service.MustGet[linking.LocalSymbolsProvider](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		// STEP 2.1: Collect imported symbols from all other documents.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateImportedSymbols) {
			allDocs := documentManager.All()
			importer.ImportSymbols(ctx, doc, allDocs)
			if ctx.Err() != nil {
				s.resetImportedSymbols(doc)
				return
			}
			advance(doc, core.DocStateImportedSymbols)
		}
		// STEP 2.2: Compute the local symbols for intra-document references.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateLocalSymbols) {
			localSymbols.LocalSymbols(ctx, doc)
			if ctx.Err() != nil {
				s.resetLocalSymbols(doc)
				return
			}
			advance(doc, core.DocStateLocalSymbols)
		}
	})

	if err := ctx.Err(); err != nil {
		return err
	}
	// Phase 2 barrier: the whole workspace has its static scopes computed.
	lock.StateChanged(core.DocStateImportedSymbols | core.DocStateLocalSymbols)

	// PHASE 3: Link and compute reference descriptions (parallel per document).
	// Resolving a reference may resolve references of other documents, which
	// reads their imported and local symbols. Keep this separate from phase 2
	// so those scopes are complete before any document starts linking.
	linker := service.MustGet[linking.Linker](s.sc)
	referenceDescriptions := service.MustGet[linking.ReferenceDescriptionsProvider](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		// STEP 3.1: Link the document to resolve all references.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateLinked) {
			linker.Link(ctx, doc)
			if ctx.Err() != nil {
				s.resetLinking(doc)
				return
			}
			advance(doc, core.DocStateLinked)
		}
		// STEP 3.2: Provide reference descriptions for the document.
		if ctx.Err() == nil && !doc.State().Has(core.DocStateReferences) {
			referenceDescriptions.ReferenceDescriptions(ctx, doc)
			if ctx.Err() != nil {
				s.resetReferences(doc)
				return
			}
			advance(doc, core.DocStateReferences)
		}
	})

	if err := ctx.Err(); err != nil {
		return err
	}
	// Phase 3 barrier: the whole workspace is now linked and indexed.
	lock.StateChanged(core.DocStateLinked | core.DocStateReferences)

	// PHASE 4: Run custom validations (parallel per document).
	validator := service.MustGet[DocumentValidator](s.sc)
	parallel.ForEach(docs, func(doc *core.Document, _ int) {
		if ctx.Err() == nil && !doc.State().Has(core.DocStateValidated) {
			diagnostics := validator.Validate(ctx, doc)
			if ctx.Err() != nil {
				s.resetValidation(doc)
				return
			}
			doc.Diagnostics = diagnostics
			advance(doc, core.DocStateValidated)
		}
	})

	if err := ctx.Err(); err != nil {
		return err
	}
	// Phase 4 barrier: the whole workspace is validated.
	lock.StateChanged(core.DocStateValidated)
	return nil
}

func (s *DefaultBuilder) Reset(doc *core.Document, state core.DocumentState) {
	switch {
	case !state.Has(core.DocStateLexed):
		s.resetLexed(doc)
		fallthrough
	case !state.Has(core.DocStateParsed):
		s.resetParsed(doc)
		fallthrough
	case !state.Has(core.DocStateExportedSymbols):
		s.resetExportedSymbols(doc)
		fallthrough
	case !state.Has(core.DocStateLocalSymbols):
		s.resetLocalSymbols(doc)
		fallthrough
	case !state.Has(core.DocStateImportedSymbols):
		s.resetImportedSymbols(doc)
		fallthrough
	case !state.Has(core.DocStateLinked):
		s.resetLinking(doc)
		fallthrough
	case !state.Has(core.DocStateReferences):
		s.resetReferences(doc)
		fallthrough
	case !state.Has(core.DocStateValidated):
		s.resetValidation(doc)
	}
	doc.SetState(doc.State() & state)
}

func (s *DefaultBuilder) resetLexed(doc *core.Document) {
	doc.Tokens = core.TokenSlice{}
	doc.Comments = core.TokenSlice{}
	doc.LexerErrors = []*core.LexerError{}
}

func (s *DefaultBuilder) resetParsed(doc *core.Document) {
	doc.Root = nil
	doc.ParserErrors = []*core.ParserError{}
	doc.References = []core.UntypedReference{}
}

func (s *DefaultBuilder) resetExportedSymbols(doc *core.Document) {
	doc.ExportedSymbols = nil
}

func (s *DefaultBuilder) resetLocalSymbols(doc *core.Document) {
	doc.LocalSymbols = nil
}

func (s *DefaultBuilder) resetImportedSymbols(doc *core.Document) {
	doc.ImportedSymbols = nil
}

func (s *DefaultBuilder) resetLinking(doc *core.Document) {
	parallel.ForEach(doc.References, func(ref core.UntypedReference, _ int) {
		ref.Reset()
	})
}

func (s *DefaultBuilder) resetReferences(doc *core.Document) {
	doc.ReferenceDescriptions = nil
}

func (s *DefaultBuilder) resetValidation(doc *core.Document) {
	doc.Diagnostics = nil
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
