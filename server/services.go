// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"log/slog"

	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// WorkspaceFolders is populated during the LSP initialize request.
type WorkspaceFolders struct {
	Value []lsp.WorkspaceFolder
}

// Connection is assigned by ConnectionBinder when the language server is started.
type Connection struct {
	Value *jsonrpc2.Connection
}

// SetupDefaultServices sets up the default services for the language server.
// It should be called together with [SetupStdioServices] or [SetupWasmServices].
// If any service is already set, it's not overwritten.
func SetupDefaultServices(sc *service.Container) {
	sc.Put(&WorkspaceFolders{})
	sc.Put(&Connection{})
	if !sc.Has[slog.Handler]() {
		sc.Put(NewSlogHandler(sc))
	}
	if !sc.Has[DiagnosticsPublisher]() {
		sc.Put(NewDiagnosticsPublisher(sc))
	}
	if !sc.Has[lsp.Server]() {
		sc.Put(NewDefaultLanguageServer(sc))
	}
	if !sc.Has[DocumentSyncher]() {
		sc.Put(NewDefaultDocumentSyncher(sc))
	}
	if !sc.Has[DefinitionProvider]() {
		sc.Put(NewDefaultDefinitionProvider(sc))
	}
	if !sc.Has[ReferencesProvider]() {
		sc.Put(NewDefaultReferencesProvider(sc))
	}
	if !sc.Has[FoldingRangeProvider]() {
		sc.Put(NewDefaultFoldingRangeProvider(sc))
	}
	if !sc.Has[RenameProvider]() {
		sc.Put(NewDefaultRenameProvider(sc))
	}
	if !sc.Has[NameFinder]() {
		sc.Put(NewDefaultNameFinder(sc))
	}
	if !sc.Has[ReferencesFinder]() {
		sc.Put(NewDefaultReferencesFinder(sc))
	}
	if !sc.Has[DocumentSymbolProvider]() {
		sc.Put(NewDefaultDocumentSymbolProvider(sc))
	}
	if !sc.Has[FuzzyMatcher]() {
		sc.Put(NewDefaultFuzzyMatcher())
	}
	if !sc.Has[CompletionProvider]() {
		sc.Put(NewDefaultCompletionProvider(sc))
	}
	if !sc.Has[SnippetRegistry]() {
		sc.Put(NewDefaultSnippetRegistry())
	}
	if !sc.Has[CompletionTriggers]() {
		sc.Put(NewDefaultCompletionTriggers())
	}
	if !sc.Has[CompletionContributor]() {
		sc.Put(NewDefaultCompletionContributor())
	}
	if !sc.Has[SignatureHelpTriggers]() {
		sc.Put(NewDefaultSignatureHelpTriggers())
	}
	if !sc.Has[DocumentHighlightProvider]() {
		sc.Put(NewDefaultDocumentHighlightProvider(sc))
	}
	if !sc.Has[WorkspaceSymbolProvider]() {
		sc.Put(NewDefaultWorkspaceSymbolProvider(sc))
	}
	if !sc.Has[DocumentationProvider]() {
		sc.Put(NewDefaultDocumentationProvider())
	}
	if !sc.Has[HoverProvider]() {
		sc.Put(NewDefaultHoverProvider(sc))
	}
}
