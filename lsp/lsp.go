// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package lsp

import (
	"context"
	"io"
	"os"

	"github.com/TypeFox/go-lsp/protocol"
	"golang.org/x/exp/jsonrpc2"
)

// LanguageServerHandlers contains the handlers for various LSP requests
type LanguageServerHandlers struct {
	// Initialize handles the initialize request
	Initialize func(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error)
	// Initialized handles the initialized notification
	Initialized func(ctx context.Context, params *protocol.InitializedParams) error
	// DidOpen handles textDocument/didOpen notifications
	DidOpen func(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error
	// DidChange handles textDocument/didChange notifications
	DidChange func(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error
	// DidClose handles textDocument/didClose notifications
	DidClose func(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error
	// Completion handles textDocument/completion requests
	Completion func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error)
	// Shutdown handles the shutdown request - server should shut down but not exit
	Shutdown func(ctx context.Context) error
	// Exit handles the exit notification - server should exit its process
	Exit func(ctx context.Context) error
}

// stdioDialer implements jsonrpc2.Dialer for stdio communication
type stdioDialer struct{}

func (d stdioDialer) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	return &stdioReadWriteCloser{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, nil
}

// stdioReadWriteCloser combines stdin/stdout into a ReadWriteCloser
type stdioReadWriteCloser struct {
	io.Reader
	io.Writer
}

func (rw *stdioReadWriteCloser) Close() error {
	// stdin/stdout don't need explicit closing
	return nil
}

// languageServer implements the protocol.Server interface
type languageServer struct {
	handlers *LanguageServerHandlers
}

func (s *languageServer) Initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	if s.handlers != nil && s.handlers.Initialize != nil {
		return s.handlers.Initialize(ctx, params)
	}
	// Default implementation with basic capabilities
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.Incremental,
			CompletionProvider: &protocol.CompletionOptions{
				ResolveProvider: false,
			},
		},
	}, nil
}

func (s *languageServer) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	if s.handlers != nil && s.handlers.Initialized != nil {
		return s.handlers.Initialized(ctx, params)
	}
	return nil
}

func (s *languageServer) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	if s.handlers != nil && s.handlers.DidOpen != nil {
		return s.handlers.DidOpen(ctx, params)
	}
	return nil
}

func (s *languageServer) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	if s.handlers != nil && s.handlers.DidChange != nil {
		return s.handlers.DidChange(ctx, params)
	}
	return nil
}

func (s *languageServer) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	if s.handlers != nil && s.handlers.DidClose != nil {
		return s.handlers.DidClose(ctx, params)
	}
	return nil
}

func (s *languageServer) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	if s.handlers != nil && s.handlers.Completion != nil {
		return s.handlers.Completion(ctx, params)
	}
	// Default empty completion list
	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        []protocol.CompletionItem{},
	}, nil
}

// Implement other required Server interface methods with no-op implementations
func (s *languageServer) Progress(ctx context.Context, params *protocol.ProgressParams) error { return nil }
func (s *languageServer) SetTrace(ctx context.Context, params *protocol.SetTraceParams) error { return nil }
func (s *languageServer) IncomingCalls(ctx context.Context, params *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) { return nil, nil }
func (s *languageServer) OutgoingCalls(ctx context.Context, params *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) { return nil, nil }
func (s *languageServer) ResolveCodeAction(ctx context.Context, params *protocol.CodeAction) (*protocol.CodeAction, error) { return nil, nil }
func (s *languageServer) ResolveCodeLens(ctx context.Context, params *protocol.CodeLens) (*protocol.CodeLens, error) { return nil, nil }
func (s *languageServer) ResolveCompletionItem(ctx context.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) { return nil, nil }
func (s *languageServer) ResolveDocumentLink(ctx context.Context, params *protocol.DocumentLink) (*protocol.DocumentLink, error) { return nil, nil }
func (s *languageServer) Exit(ctx context.Context) error {
	if s.handlers != nil && s.handlers.Exit != nil {
		return s.handlers.Exit(ctx)
	}
	return nil
}
func (s *languageServer) Resolve(ctx context.Context, params *protocol.InlayHint) (*protocol.InlayHint, error) { return nil, nil }
func (s *languageServer) DidChangeNotebookDocument(ctx context.Context, params *protocol.DidChangeNotebookDocumentParams) error { return nil }
func (s *languageServer) DidCloseNotebookDocument(ctx context.Context, params *protocol.DidCloseNotebookDocumentParams) error { return nil }
func (s *languageServer) DidOpenNotebookDocument(ctx context.Context, params *protocol.DidOpenNotebookDocumentParams) error { return nil }
func (s *languageServer) DidSaveNotebookDocument(ctx context.Context, params *protocol.DidSaveNotebookDocumentParams) error { return nil }
func (s *languageServer) Shutdown(ctx context.Context) error {
	if s.handlers != nil && s.handlers.Shutdown != nil {
		return s.handlers.Shutdown(ctx)
	}
	return nil
}
func (s *languageServer) CodeAction(ctx context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) { return nil, nil }
func (s *languageServer) CodeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) { return nil, nil }
func (s *languageServer) ColorPresentation(ctx context.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) { return nil, nil }
func (s *languageServer) Declaration(ctx context.Context, params *protocol.DeclarationParams) (*protocol.Or_textDocument_declaration, error) { return nil, nil }
func (s *languageServer) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) { return nil, nil }
func (s *languageServer) Diagnostic(ctx context.Context, params *protocol.DocumentDiagnosticParams) (*protocol.DocumentDiagnosticReport, error) { return nil, nil }
func (s *languageServer) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error { return nil }
func (s *languageServer) DocumentColor(ctx context.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) { return nil, nil }
func (s *languageServer) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) { return nil, nil }
func (s *languageServer) DocumentLink(ctx context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) { return nil, nil }
func (s *languageServer) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]any, error) { return nil, nil }
func (s *languageServer) FoldingRange(ctx context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) { return nil, nil }
func (s *languageServer) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) { return nil, nil }
func (s *languageServer) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) { return nil, nil }
func (s *languageServer) Implementation(ctx context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) { return nil, nil }
func (s *languageServer) InlayHint(ctx context.Context, params *protocol.InlayHintParams) ([]protocol.InlayHint, error) { return nil, nil }
func (s *languageServer) InlineCompletion(ctx context.Context, params *protocol.InlineCompletionParams) (*protocol.Or_Result_textDocument_inlineCompletion, error) { return nil, nil }
func (s *languageServer) InlineValue(ctx context.Context, params *protocol.InlineValueParams) ([]protocol.InlineValue, error) { return nil, nil }
func (s *languageServer) LinkedEditingRange(ctx context.Context, params *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) { return nil, nil }
func (s *languageServer) DiagnosticWorkspace(ctx context.Context, params *protocol.WorkspaceDiagnosticParams) (*protocol.WorkspaceDiagnosticReport, error) { return nil, nil }
func (s *languageServer) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error { return nil }
func (s *languageServer) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error { return nil }
func (s *languageServer) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error { return nil }
func (s *languageServer) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) error { return nil }
func (s *languageServer) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) error { return nil }
func (s *languageServer) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) error { return nil }
func (s *languageServer) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (any, error) { return nil, nil }
func (s *languageServer) Symbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) { return nil, nil }
func (s *languageServer) TextDocumentContent(ctx context.Context, params *protocol.TextDocumentContentParams) (*string, error) { return nil, nil }
func (s *languageServer) WillCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) { return nil, nil }
func (s *languageServer) WillDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) { return nil, nil }
func (s *languageServer) WillRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) { return nil, nil }
func (s *languageServer) ResolveWorkspaceSymbol(ctx context.Context, params *protocol.WorkspaceSymbol) (*protocol.WorkspaceSymbol, error) { return nil, nil }
func (s *languageServer) Moniker(ctx context.Context, params *protocol.MonikerParams) ([]protocol.Moniker, error) { return nil, nil }
func (s *languageServer) OnTypeFormatting(ctx context.Context, params *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) { return nil, nil }
func (s *languageServer) PrepareCallHierarchy(ctx context.Context, params *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) { return nil, nil }
func (s *languageServer) PrepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (*protocol.PrepareRenameResult, error) { return nil, nil }
func (s *languageServer) PrepareTypeHierarchy(ctx context.Context, params *protocol.TypeHierarchyPrepareParams) ([]protocol.TypeHierarchyItem, error) { return nil, nil }
func (s *languageServer) RangeFormatting(ctx context.Context, params *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) { return nil, nil }
func (s *languageServer) RangesFormatting(ctx context.Context, params *protocol.DocumentRangesFormattingParams) ([]protocol.TextEdit, error) { return nil, nil }
func (s *languageServer) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) { return nil, nil }
func (s *languageServer) Rename(ctx context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) { return nil, nil }
func (s *languageServer) SelectionRange(ctx context.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) { return nil, nil }
func (s *languageServer) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) { return nil, nil }
func (s *languageServer) SemanticTokensFullDelta(ctx context.Context, params *protocol.SemanticTokensDeltaParams) (any, error) { return nil, nil }
func (s *languageServer) SemanticTokensRange(ctx context.Context, params *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) { return nil, nil }
func (s *languageServer) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) { return nil, nil }
func (s *languageServer) TypeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) ([]protocol.Location, error) { return nil, nil }
func (s *languageServer) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) error { return nil }
func (s *languageServer) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) { return nil, nil }
func (s *languageServer) Subtypes(ctx context.Context, params *protocol.TypeHierarchySubtypesParams) ([]protocol.TypeHierarchyItem, error) { return nil, nil }
func (s *languageServer) Supertypes(ctx context.Context, params *protocol.TypeHierarchySupertypesParams) ([]protocol.TypeHierarchyItem, error) { return nil, nil }
func (s *languageServer) WorkDoneProgressCancel(ctx context.Context, params *protocol.WorkDoneProgressCancelParams) error { return nil }

// simpleBinder implements jsonrpc2.Binder
type simpleBinder struct {
	server protocol.Server
}

func (b *simpleBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	return jsonrpc2.ConnectionOptions{
		Handler: protocol.ServerHandler(b.server),
	}, nil
}

// StartLanguageServer starts a language server using the provided handlers.
// It sets up JSON-RPC communication over stdio and handles the essential LSP messages.
func StartLanguageServer(ctx context.Context, handlers *LanguageServerHandlers) error {
	server := &languageServer{handlers: handlers}
	binder := &simpleBinder{server: server}
	
	// Create a connection using stdio
	conn, err := jsonrpc2.Dial(ctx, stdioDialer{}, binder)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close() // Ignore error in defer
	}()
	
	// Wait for the connection to close
	return conn.Wait()
}
