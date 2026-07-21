// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package server implements the Language Server Protocol (LSP) side of
// Fastbelt. It exposes the language services built by
// [typefox.dev/fastbelt/workspace] — parsed, linked, and validated documents —
// to editors and IDEs: text documents are synchronized through LSP
// notifications, and language features such as completion, go to definition,
// find references, hover, and rename are answered from the workspace.
//
// # The language server and its method services
//
// [DefaultLanguageServer] implements the [typefox.dev/lsp.Server] interface,
// but contains no feature logic of its own. Each supported LSP method is
// delegated to a dedicated service interface (for example
// [CompletionProvider] or [DocumentSyncher]) that is looked up in the
// [typefox.dev/fastbelt/util/service.Container] when the request arrives.
// Request-handling methods first acquire the
// [typefox.dev/fastbelt/workspace.Lock] for reading, so a request never
// observes a half-built workspace. The capabilities announced to the client
// in the initialize response are derived from which services are registered,
// so the set of advertised features always matches the set of available
// services.
//
// [SetupDefaultServices] registers the default implementation of every LSP
// method service that has not been registered yet. To customize a feature,
// register your own implementation of the corresponding interface. Services
// that need to inspect the client's initialize request implement
// [InitializeParticipant].
//
// # Shared name and reference lookup
//
// Most position-based LSP features start by answering the same question:
// which named symbol does the cursor point at? Two services encapsulate this
// so that all features resolve names consistently:
//
//   - [NameFinder] turns the token(s) at a cursor position into a
//     [FoundName], pairing the unit under the cursor with the unit that
//     holds the name of the referenced symbol — following cross-references
//     to their target, or staying in place when the cursor is already on a
//     declaration's name.
//   - [ReferencesFinder] yields every reference to a given AST node, drawing
//     on the reference descriptions indexed per document. [FindReferencesOptions]
//     controls whether the search is restricted to a single document and
//     whether the declaration itself is included.
//
// The default definition, references, document highlight, hover, and rename
// providers are thin compositions of these two services plus conversion to
// LSP result types. Registering a custom [NameFinder] or [ReferencesFinder]
// therefore adjusts all of these features at once.
//
// # Starting a server
//
// [StartLanguageServer] is the main entry point. It expects a fully
// configured service container: the language-specific and workspace
// services, the server services from [SetupDefaultServices], and a transport
// registered with [SetupStdioServices] (communication over stdin/stdout) or
// SetupWasmServices (communication with a JavaScript host, available in
// js/wasm builds). It dials a JSON-RPC connection using the registered
// [golang.org/x/exp/jsonrpc2.Dialer] and [golang.org/x/exp/jsonrpc2.Binder]
// — by default [DefaultBinder], which attaches the [typefox.dev/lsp.Server]
// from the container as the message handler — and blocks until the
// connection is closed, which the default server triggers on the LSP exit
// notification.
//
// A typical main function for a generated language looks like this:
//
//	func main() {
//		sc := service.NewContainer()
//		mylang.SetupServices(sc) // language, textdoc, linking, workspace services
//		mylang.SetupGeneratedServerServices(sc)
//		server.SetupDefaultServices(sc)
//		server.SetupStdioServices(sc)
//		sc.Seal()
//		if err := server.StartLanguageServer(context.Background(), sc); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// While the server runs, [DiagnosticsPublisher] pushes parser, linker, and
// validation diagnostics to the client whenever an open document is rebuilt,
// and [SlogHandler] forwards log output to the client as LSP log messages.
package server
