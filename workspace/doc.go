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
//   - [Lock] — read/write coordination with atomic write-to-read downgrade
//   - [Builder] — runs the build pipeline that takes documents to a linked, validated state
//   - [DocumentManager] — concurrent in-memory store of all documents, keyed by URI
//   - [DocumentUpdater] — entry point for edits; serializes mutations and triggers builds
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
// Phases 1 and 2 are the write phase: they mutate document data and require
// exclusive access. Phase 3 only reads document data, so between phase 2 and
// phase 3 the builder calls the downgrade function passed to [Builder.Build],
// releasing exclusive access (see Concurrency below) so that read requests can
// proceed while validation runs.
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
// Read-only requests such as completion, hover, and go-to-definition run under
// [Lock.Read] and consult [DocumentManager] for the current document state.
//
// # Concurrency
//
// [Lock] serializes writes against reads and provides the atomic write-to-read
// downgrade the build pipeline relies on. [Lock.Write] grants exclusive access
// for the write phase, then the downgrade callback atomically converts that
// exclusive hold into a shared read hold, with no window for another writer to
// intervene. This guarantees validation (phase 3) observes a consistent,
// fully linked snapshot and completes before the next write phase begins.
// Writes have priority: a pending write blocks new readers, and starting a
// write cancels any write still in progress so the freshest edit wins.
package workspace
