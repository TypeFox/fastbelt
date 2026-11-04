// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package server implements the Language Server Protocol (LSP) side of
// Fastbelt. It exposes the parsed, linked, and validated documents managed by
// [typefox.dev/fastbelt/workspace] to editors and IDEs: text documents are
// synchronized through LSP notifications, and language features such as
// completion, go to definition, find references, hover, and rename are
// computed from the documents in the workspace.
//
// # The language server and its method services
//
// Each LSP feature is provided by a service interface, such as
// [CompletionProvider] or [DefinitionProvider], that [DefaultLanguageServer]
// looks up in the [typefox.dev/fastbelt/util/service.Container].
//
// [SetupDefaultServices] registers the [DefaultLanguageServer], the
// [DiagnosticsPublisher], a [SlogHandler], and default implementations for
// document sync, completion, definition, references, document highlight,
// document and workspace symbols, folding ranges, hover, and rename, plus the
// helper services they use. A default feature implementation is only
// registered if no implementation of its interface was registered
// beforehand. Other features, such as semantic tokens or code actions, are
// not registered by default; register an implementation to enable them
// ([TokenBasedSemanticTokensProvider] is a ready-made one for semantic
// tokens). Services that need to inspect the client's initialize request
// implement [InitializeParticipant].
//
// # Shared name and reference lookup
//
// Most position-based features first determine which symbol the cursor
// points at. [NameFinder] resolves the name under the cursor to the symbol it
// refers to, and [ReferencesFinder] collects all occurrences of a symbol. The
// default providers are built on these two services, so all features resolve
// names consistently; the type documentation lists which features each one
// affects. Custom providers can reuse [TargetAtCursor], which returns the
// target node and the source range at a cursor position.
//
// # Starting a server
//
// [StartLanguageServer] is the main entry point. It expects a fully
// configured service container: the language-specific and workspace
// services, the server services from [SetupDefaultServices], and a transport
// registered with [SetupStdioServices] (stdin/stdout) or SetupWasmServices
// (JavaScript host, js/wasm builds only).
//
// A typical main function of a language server looks like this:
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
// While the server runs, [DiagnosticsPublisher] pushes the diagnostics of
// open documents to the client, and [SlogHandler] forwards log output to the
// client as LSP log messages.
package server
