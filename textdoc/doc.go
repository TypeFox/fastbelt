// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package textdoc provides text documents and LSP position mapping for Fastbelt.
//
// textdoc is the text layer of the framework. It holds source text together
// with LSP metadata (URI, language identifier, version) and converts between
// byte offsets and [typefox.dev/lsp] positions. It does not parse text or
// attach semantic data such as ASTs, symbols, or diagnostics. That role
// belongs to [typefox.dev/fastbelt.Document] and the
// [typefox.dev/fastbelt/workspace] package.
//
// # Document handles
//
// [Handle] is the common read interface for document text. Lexers, parsers,
// and other consumers work against Handle so they do not depend on whether
// the text came from disk or from an open editor buffer.
//
//   - [File] is an immutable snapshot of on-disk content.
//   - [Overlay] is a mutable editor buffer that may contain unsaved changes.
//
// When a workspace is initialized, files read from disk become [File] values
// and are wrapped in [typefox.dev/fastbelt.Document] instances by the workspace
// initializer. When a language client opens a document, the LSP server creates
// an [Overlay] for the same URI. Incremental didChange notifications update
// that overlay in place via [Overlay.Update].
//
// [Store] caches open overlays and, optionally, file snapshots. [Store.Get]
// returns the overlay when one exists for a URI, otherwise the cached file,
// so callers always see the effective editor content for open documents.
//
// On didClose, the server removes the overlay and reverts the workspace to
// on-disk [File] content when the URI refers to a local file.
//
// # Service registration
//
// Call [SetupDefaultServices] during service container setup to register a
// [DefaultStore]. Language projects typically invoke it from their
// SetupServices function alongside other framework SetupDefaultServices
// functions.
//
// For tests and standalone tooling that do not run an LSP server, create
// [File] or [Overlay] values directly and pass them to
// [typefox.dev/fastbelt.NewDocument].
package textdoc
