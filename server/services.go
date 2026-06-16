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
	service.Put(sc, &WorkspaceFolders{})
	service.Put(sc, &Connection{})
	if !service.Has[slog.Handler](sc) {
		service.Put(sc, NewSlogHandler(sc))
	}
	if !service.Has[lsp.Server](sc) {
		service.Put(sc, NewDefaultLanguageServer(sc))
	}
	if !service.Has[DocumentSyncher](sc) {
		service.Put(sc, NewDefaultDocumentSyncher(sc))
	}
	if !service.Has[DefinitionProvider](sc) {
		service.Put(sc, NewDefaultDefinitionProvider(sc))
	}
	if !service.Has[ReferencesProvider](sc) {
		service.Put(sc, NewDefaultReferencesProvider(sc))
	}
	if !service.Has[FoldingRangeProvider](sc) {
		service.Put(sc, NewDefaultFoldingRangeProvider(sc))
	}
	if !service.Has[RenameProvider](sc) {
		service.Put(sc, NewDefaultRenameProvider(sc))
	}
	if !service.Has[NameFinder](sc) {
		service.Put(sc, NewDefaultNameFinder(sc))
	}
	if !service.Has[ReferencesFinder](sc) {
		service.Put(sc, NewDefaultReferencesFinder(sc))
	}
	if !service.Has[DocumentSymbolProvider](sc) {
		service.Put(sc, NewDefaultDocumentSymbolProvider(sc))
	}
	if !service.Has[FuzzyMatcher](sc) {
		service.Put(sc, NewDefaultFuzzyMatcher())
	}
	if !service.Has[CompletionProvider](sc) {
		service.Put(sc, NewDefaultCompletionProvider(sc))
	}
	if !service.Has[SnippetRegistry](sc) {
		service.Put(sc, NewDefaultSnippetRegistry())
	}
	if !service.Has[CompletionTriggers](sc) {
		service.Put(sc, NewDefaultCompletionTriggers())
	}
	if !service.Has[CompletionContributor](sc) {
		service.Put(sc, NewDefaultCompletionContributor())
	}
	if !service.Has[DocumentHighlightProvider](sc) {
		service.Put(sc, NewDefaultDocumentHighlightProvider(sc))
	}
	if !service.Has[WorkspaceSymbolProvider](sc) {
		service.Put(sc, NewDefaultWorkspaceSymbolProvider(sc))
	}
}
