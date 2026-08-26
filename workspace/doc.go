// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package workspace manages the lifecycle of Fastbelt documents: loading files,
// applying edits, and rebuilding documents through parsing, cross-reference
// linking, and validation as their contents change.
//
// A workspace is the in-memory collection of all documents that belong to one
// language, held by [DocumentManager]. Each document is a
// [typefox.dev/fastbelt.Document] that accumulates text, tokens, an AST, symbol
// tables, resolved references, and diagnostics. This package drives every
// document from raw text to a fully linked and validated state and keeps that
// state consistent as files are opened, changed, and deleted. It is the layer
// between the Language Server Protocol handlers in
// [typefox.dev/fastbelt/server] and the cross-reference resolution in
// [typefox.dev/fastbelt/linking].
//
// # Services
//
// Most behavior is provided through services in a
// [typefox.dev/fastbelt/util/service.Container]. [SetupDefaultServices]
// registers the framework defaults:
//
//   - [Initializer] — loads files matching [FileExtensions] from workspace folders on startup
//   - [Lock] — read/write coordination with per-document state admission for readers
//   - [Builder] — runs the build pipeline that takes documents to a linked, validated state
//   - [DocumentManager] — concurrent in-memory store of all documents, keyed by URI
//   - [DocumentUpdater] — entry point for edits; serializes mutations and triggers builds
//   - [DocumentChangeImpact] — reports cross-document reference dependencies on changed files
//   - [DocumentParser] — lexes and parses a single document into tokens and an AST
//   - [DocumentValidator] — collects diagnostics for a single document
//
// # The build lifecycle
//
// [Builder] is the centerpiece of the package. It governs how a document moves
// through a sequence of build steps, each recorded as a bit in the document's
// [typefox.dev/fastbelt.DocumentState]. [Builder.Build] processes a batch of
// documents in three phases, running documents in parallel within each phase
// and checking the context for cancellation between steps:
//
//   - Phase 1 (per document): parse into an AST ([DocumentParser]) and compute
//     the symbols the document exports to others.
//   - Phase 2 (per document): import symbols from other documents, compute
//     local symbols, link all cross-references, and index reference
//     descriptions. This phase needs every document's exports from phase 1, so
//     all documents must finish phase 1 before any enters phase 2.
//   - Phase 3 (per document): run validations and store the resulting
//     diagnostics on the document.
//
// The linking-related steps in phases 1 and 2 are performed by the services in
// [typefox.dev/fastbelt/linking], which the builder resolves from the container;
// see that package for how each step works.
//
// The whole build runs under the exclusive [Lock.Write], but requests do not
// have to wait for it to finish: each completed step is published through the
// document's [typefox.dev/fastbelt.DocumentState], and [Lock.ReadAt] admits
// readers as soon as the documents they target reach the states they need
// (see Concurrency below). A document-symbol request, for example, only needs
// a parsed AST and is served while linking and validation are still running.
//
// Because each step is guarded by its [typefox.dev/fastbelt.DocumentState] bit,
// builds are incremental. [Builder.Reset] clears selected steps of a document
// by keeping a bitmask of states and resetting the rest; a later [Builder.Build]
// re-runs only the cleared steps. For example, after a text edit the updater
// keeps the parse and symbol steps of unaffected documents and resets only
// linking and validation, so unchanged work is not repeated.
//
// [Builder.AddBuildStepListener] registers callbacks that fire as documents
// complete selected steps, which the server uses to publish diagnostics as soon
// as validation finishes.
//
// # How the services work together
//
// On startup the server calls [Initializer.Initialize], which walks the open
// workspace folders, reads files whose extension matches [FileExtensions], and
// registers them with [DocumentManager].
//
// When a file changes, the LSP document sync layer calls [DocumentUpdater.Update]
// with the changed and deleted handles. The updater runs under [Lock.Write]: it
// updates [DocumentManager], collects the documents to rebuild, calls
// [Builder.Reset] on them, and then calls [Builder.Build]. A newer edit cancels
// the context of an in-progress build, so superseded builds stop quickly while
// the latest one runs to completion.
//
// Read-only requests run under [Lock.Read] when they may touch arbitrary
// documents (find references, workspace symbols), or under [Lock.ReadAt] when
// they target a known set of documents and a known build state.
//
// # Concurrency
//
// [Lock.Write] grants exclusive access for a build cycle. A write starts in
// the mutation phase, in which documents may be created, deleted, and reset;
// once [Builder.Build] begins, document states only advance, and every
// advance is signalled through [Lock.StateChanged]. [Lock.ReadAt] relies on
// this monotonicity: it admits a reader as soon as every requested document
// carries the requested state bits, even while the build is still writing
// later phases or other documents. Setting a state bit (an atomic store)
// publishes all data written by that step, so an admitted reader always
// observes complete data for the states it requested — provided it touches
// only the documents it listed and only the data those states produce.
// For whole-workspace access at a specific state, ReadAt accepts an empty
// document list: it is then admitted based on the workspace floor — the state
// every document is guaranteed to have reached — which the builder reports
// through [Lock.StateChanged] at each phase barrier, without any document
// scan. [Lock.Read] is simply the state-agnostic variant: it waits until no
// write is active or pending.
//
// Writes have priority: a pending write blocks new readers (including
// ReadAt), and starting a write cancels any write still in progress so the
// freshest edit wins. A write acquires the lock only after all readers have
// drained, so a reader admitted via ReadAt never observes a document reset.
package workspace
