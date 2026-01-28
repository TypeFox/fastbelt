// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"io"
	"os"

	"github.com/TypeFox/go-lsp/protocol"
	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/workspace"
)

// LanguageServerHandlers contains the handlers for various LSP requests.
// TODO extract these handlers into separate services instead of having them all here.
type LanguageServerHandlers struct {
	// Initialized handles the initialized notification
	Initialized func(ctx context.Context, params *protocol.InitializedParams) error
	// Completion handles textDocument/completion requests
	Completion func(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error)
	// Shutdown handles the shutdown request - server should shut down but not exit
	Shutdown func(ctx context.Context) error
}

type LanguageServer interface {
	protocol.Server
}

// DefaultLanguageServer implements the LanguageServer interface
type DefaultLanguageServer struct {
	srv ServerSrvCont
}

// NewDefaultLanguageServer creates a new default language server.
func NewDefaultLanguageServer(srv ServerSrvCont) LanguageServer {
	return &DefaultLanguageServer{srv: srv}
}

func (s *DefaultLanguageServer) Initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
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

func (s *DefaultLanguageServer) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	handlers := s.srv.Server().LanguageServerHandlers
	if handlers != nil && handlers.Initialized != nil {
		return handlers.Initialized(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) Shutdown(ctx context.Context) error {
	handlers := s.srv.Server().LanguageServerHandlers
	if handlers != nil && handlers.Shutdown != nil {
		return handlers.Shutdown(ctx)
	}
	return nil
}

func (s *DefaultLanguageServer) Exit(ctx context.Context) error {
	// Close the connection to allow the server to exit
	connection := s.srv.Server().Connection
	if connection != nil {
		return connection.Close()
	}
	return nil
}

func (s *DefaultLanguageServer) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidOpen(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidChange(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidClose(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	handlers := s.srv.Server().LanguageServerHandlers
	if handlers != nil && handlers.Completion != nil {
		return handlers.Completion(ctx, params)
	}
	// Default empty completion list
	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        []protocol.CompletionItem{},
	}, nil
}

// Implement other required Server interface methods with no-op implementations
func (s *DefaultLanguageServer) Progress(ctx context.Context, params *protocol.ProgressParams) error {
	return nil
}
func (s *DefaultLanguageServer) SetTrace(ctx context.Context, params *protocol.SetTraceParams) error {
	return nil
}
func (s *DefaultLanguageServer) IncomingCalls(ctx context.Context, params *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) OutgoingCalls(ctx context.Context, params *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveCodeAction(ctx context.Context, params *protocol.CodeAction) (*protocol.CodeAction, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveCodeLens(ctx context.Context, params *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveCompletionItem(ctx context.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveDocumentLink(ctx context.Context, params *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Resolve(ctx context.Context, params *protocol.InlayHint) (*protocol.InlayHint, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DidChangeNotebookDocument(ctx context.Context, params *protocol.DidChangeNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidCloseNotebookDocument(ctx context.Context, params *protocol.DidCloseNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidOpenNotebookDocument(ctx context.Context, params *protocol.DidOpenNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidSaveNotebookDocument(ctx context.Context, params *protocol.DidSaveNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) CodeAction(ctx context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) CodeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ColorPresentation(ctx context.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Declaration(ctx context.Context, params *protocol.DeclarationParams) (*protocol.Or_textDocument_declaration, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Diagnostic(ctx context.Context, params *protocol.DocumentDiagnosticParams) (*protocol.DocumentDiagnosticReport, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidSave(ctx, params)
	}
	return nil
}
func (s *DefaultLanguageServer) DocumentColor(ctx context.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentLink(ctx context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]any, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) FoldingRange(ctx context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Implementation(ctx context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) InlayHint(ctx context.Context, params *protocol.InlayHintParams) ([]protocol.InlayHint, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) InlineCompletion(ctx context.Context, params *protocol.InlineCompletionParams) (*protocol.Or_Result_textDocument_inlineCompletion, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) InlineValue(ctx context.Context, params *protocol.InlineValueParams) ([]protocol.InlineValue, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) LinkedEditingRange(ctx context.Context, params *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DiagnosticWorkspace(ctx context.Context, params *protocol.WorkspaceDiagnosticParams) (*protocol.WorkspaceDiagnosticReport, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (any, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Symbol(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) TextDocumentContent(ctx context.Context, params *protocol.TextDocumentContentParams) (*string, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveWorkspaceSymbol(ctx context.Context, params *protocol.WorkspaceSymbol) (*protocol.WorkspaceSymbol, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Moniker(ctx context.Context, params *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) OnTypeFormatting(ctx context.Context, params *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) PrepareCallHierarchy(ctx context.Context, params *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) PrepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (*protocol.PrepareRenameResult, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) PrepareTypeHierarchy(ctx context.Context, params *protocol.TypeHierarchyPrepareParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) RangeFormatting(ctx context.Context, params *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) RangesFormatting(ctx context.Context, params *protocol.DocumentRangesFormattingParams) ([]protocol.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Rename(ctx context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SelectionRange(ctx context.Context, params *protocol.SelectionRangeParams) ([]protocol.SelectionRange, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensFullDelta(ctx context.Context, params *protocol.SemanticTokensDeltaParams) (any, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensRange(ctx context.Context, params *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) TypeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.WillSave(ctx, params)
	}
	return nil
}
func (s *DefaultLanguageServer) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	if s.srv.Server().DocumentSyncher != nil {
		return s.srv.Server().DocumentSyncher.WillSaveWaitUntil(ctx, params)
	}
	return []protocol.TextEdit{}, nil
}
func (s *DefaultLanguageServer) Subtypes(ctx context.Context, params *protocol.TypeHierarchySubtypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Supertypes(ctx context.Context, params *protocol.TypeHierarchySupertypesParams) ([]protocol.TypeHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WorkDoneProgressCancel(ctx context.Context, params *protocol.WorkDoneProgressCancelParams) error {
	return nil
}

// StartLanguageServer starts a language server using the service container.
// It sets up JSON-RPC communication over stdio and handles the essential LSP messages.
func StartLanguageServer(ctx context.Context, srv ServerSrvCont) error {
	dialer := srv.Server().ConnectionDialer
	binder := srv.Server().ConnectionBinder

	// Create a connection using the configured dialer and binder
	conn, err := jsonrpc2.Dial(ctx, dialer, binder)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close() // Ignore error in defer
	}()

	// Register validation listener to publish diagnostics
	client := protocol.ClientDispatcher(conn)
	srv.Workspace().Builder.AddValidationListener(func(ctx context.Context, results []workspace.ValidationResult) error {
		for _, result := range results {
			// Collect diagnostics from lexer and parser errors
			diagnostics := []protocol.Diagnostic{}
			diagnostics = append(diagnostics, workspace.CreateLexerDiagnostics(result.Result.LexerErrors)...)
			diagnostics = append(diagnostics, workspace.CreateParserDiagnostics(result.Document, result.Result.ParserErrors)...)
			diagnostics = append(diagnostics, workspace.CreateLinkerDiagnostics(result.Document, result.Result.Root)...)
			// Publish diagnostics (empty array if no errors to clear previous diagnostics)
			params := &protocol.PublishDiagnosticsParams{
				URI:         result.Document.URI(),
				Version:     result.Document.Version(),
				Diagnostics: diagnostics,
			}
			if err := client.PublishDiagnostics(ctx, params); err != nil {
				return err
			}
		}
		return nil
	})

	// Wait for the connection to close
	return conn.Wait()
}

// DefaultBinder implements the jsonrpc2.Binder interface
type DefaultBinder struct {
	srv ServerSrvCont
}

// NewDefaultBinder creates a new default binder.
func NewDefaultBinder(srv ServerSrvCont) jsonrpc2.Binder {
	return &DefaultBinder{srv: srv}
}

func (b *DefaultBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	b.srv.Server().Connection = conn
	return jsonrpc2.ConnectionOptions{
		Handler: protocol.ServerHandler(b.srv.Server().LanguageServer),
	}, nil
}

// StdioDialer implements jsonrpc2.Dialer for stdio communication
type StdioDialer struct{}

func (d StdioDialer) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
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
