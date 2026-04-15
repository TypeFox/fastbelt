// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"golang.org/x/exp/jsonrpc2"
	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

type LanguageServer interface {
	lsp.Server
}

// DefaultLanguageServer implements the LanguageServer interface
type DefaultLanguageServer struct {
	srv ServerSrvCont
}

// NewDefaultLanguageServer creates a new default language server.
func NewDefaultLanguageServer(srv ServerSrvCont) *DefaultLanguageServer {
	return &DefaultLanguageServer{srv: srv}
}

func (s *DefaultLanguageServer) Initialize(ctx context.Context, params *lsp.ParamInitialize) (*lsp.InitializeResult, error) {
	slogHandler := s.srv.Server().SlogHandler
	if slogHandler != nil {
		// Set the default logger to use the configured slog handler
		// It will send logs to the client via the LSP connection
		slog.SetDefault(slog.New(slogHandler))
	}
	s.srv.Server().WorkspaceFolders = params.WorkspaceFolders
	definitionProvider := s.srv.Server().DefinitionProvider
	referencesProvider := s.srv.Server().ReferencesProvider
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: lsp.Incremental,
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider: false,
			},
			DefinitionProvider: &lsp.Or_ServerCapabilities_definitionProvider{
				Value: definitionProvider != nil,
			},
			ReferencesProvider: &lsp.Or_ServerCapabilities_referencesProvider{
				Value: referencesProvider != nil,
			},
		},
	}, nil
}

func (s *DefaultLanguageServer) Initialized(ctx context.Context, params *lsp.InitializedParams) error {
	initializer := s.srv.Workspace().Initializer
	if initializer != nil {
		return initializer.Initialize(ctx, s.srv.Server().WorkspaceFolders)
	}
	return nil
}

func (s *DefaultLanguageServer) Shutdown(ctx context.Context) error {
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

func (s *DefaultLanguageServer) DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidOpen(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidChange(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidClose(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) Completion(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
	return &lsp.CompletionList{
		IsIncomplete: false,
		Items:        []lsp.CompletionItem{},
	}, nil
}

// Implement other required Server interface methods with no-op implementations
func (s *DefaultLanguageServer) Progress(ctx context.Context, params *lsp.ProgressParams) error {
	return nil
}
func (s *DefaultLanguageServer) SetTrace(ctx context.Context, params *lsp.SetTraceParams) error {
	return nil
}
func (s *DefaultLanguageServer) IncomingCalls(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) OutgoingCalls(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveCodeAction(ctx context.Context, params *lsp.CodeAction) (*lsp.CodeAction, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveCodeLens(ctx context.Context, params *lsp.CodeLens) (*lsp.CodeLens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveCompletionItem(ctx context.Context, params *lsp.CompletionItem) (*lsp.CompletionItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveDocumentLink(ctx context.Context, params *lsp.DocumentLink) (*lsp.DocumentLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Resolve(ctx context.Context, params *lsp.InlayHint) (*lsp.InlayHint, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DidChangeNotebookDocument(ctx context.Context, params *lsp.DidChangeNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidCloseNotebookDocument(ctx context.Context, params *lsp.DidCloseNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidOpenNotebookDocument(ctx context.Context, params *lsp.DidOpenNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidSaveNotebookDocument(ctx context.Context, params *lsp.DidSaveNotebookDocumentParams) error {
	return nil
}
func (s *DefaultLanguageServer) CodeAction(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) CodeLens(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ColorPresentation(ctx context.Context, params *lsp.ColorPresentationParams) ([]lsp.ColorPresentation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Declaration(ctx context.Context, params *lsp.DeclarationParams) ([]lsp.DefinitionLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Definition(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.DefinitionLink, error) {
	definitionProvider := s.srv.Server().DefinitionProvider
	if definitionProvider == nil {
		return nil, nil
	}
	var result []lsp.DefinitionLink
	var providerErr error
	if err := s.srv.Workspace().Lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = definitionProvider.HandleDefinitionRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Diagnostic(ctx context.Context, params *lsp.DocumentDiagnosticParams) (*lsp.DocumentDiagnosticReport, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.DidSave(ctx, params)
	}
	return nil
}
func (s *DefaultLanguageServer) DocumentColor(ctx context.Context, params *lsp.DocumentColorParams) ([]lsp.ColorInformation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentHighlight(ctx context.Context, params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentLink(ctx context.Context, params *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentSymbol(ctx context.Context, params *lsp.DocumentSymbolParams) ([]any, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) FoldingRange(ctx context.Context, params *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Formatting(ctx context.Context, params *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Hover(ctx context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Implementation(ctx context.Context, params *lsp.ImplementationParams) ([]lsp.DefinitionLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) InlayHint(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) InlineCompletion(ctx context.Context, params *lsp.InlineCompletionParams) (*lsp.Or_Result_textDocument_inlineCompletion, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) InlineValue(ctx context.Context, params *lsp.InlineValueParams) ([]lsp.InlineValue, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) LinkedEditingRange(ctx context.Context, params *lsp.LinkedEditingRangeParams) (*lsp.LinkedEditingRanges, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DiagnosticWorkspace(ctx context.Context, params *lsp.WorkspaceDiagnosticParams) (*lsp.WorkspaceDiagnosticReport, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DidChangeConfiguration(ctx context.Context, params *lsp.DidChangeConfigurationParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidChangeWatchedFiles(ctx context.Context, params *lsp.DidChangeWatchedFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidChangeWorkspaceFolders(ctx context.Context, params *lsp.DidChangeWorkspaceFoldersParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidCreateFiles(ctx context.Context, params *lsp.CreateFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidDeleteFiles(ctx context.Context, params *lsp.DeleteFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) DidRenameFiles(ctx context.Context, params *lsp.RenameFilesParams) error {
	return nil
}
func (s *DefaultLanguageServer) ExecuteCommand(ctx context.Context, params *lsp.ExecuteCommandParams) (any, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Symbol(ctx context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) TextDocumentContent(ctx context.Context, params *lsp.TextDocumentContentParams) (*lsp.TextDocumentContentResult, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillCreateFiles(ctx context.Context, params *lsp.CreateFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillDeleteFiles(ctx context.Context, params *lsp.DeleteFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillRenameFiles(ctx context.Context, params *lsp.RenameFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) ResolveWorkspaceSymbol(ctx context.Context, params *lsp.WorkspaceSymbol) (*lsp.WorkspaceSymbol, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Moniker(ctx context.Context, params *lsp.MonikerParams) ([]lsp.Moniker, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) OnTypeFormatting(ctx context.Context, params *lsp.DocumentOnTypeFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) PrepareCallHierarchy(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) PrepareRename(ctx context.Context, params *lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) PrepareTypeHierarchy(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) RangeFormatting(ctx context.Context, params *lsp.DocumentRangeFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) RangesFormatting(ctx context.Context, params *lsp.DocumentRangesFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) References(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error) {
	referencesProvider := s.srv.Server().ReferencesProvider
	if referencesProvider == nil {
		return nil, nil
	}
	var result []lsp.Location
	var providerErr error
	if err := s.srv.Workspace().Lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = referencesProvider.HandleReferencesRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Rename(ctx context.Context, params *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SelectionRange(ctx context.Context, params *lsp.SelectionRangeParams) ([]lsp.SelectionRange, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensFull(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensFullDelta(ctx context.Context, params *lsp.SemanticTokensDeltaParams) (any, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensRange(ctx context.Context, params *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SignatureHelp(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) TypeDefinition(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams) error {
	if s.srv.Server().DocumentSyncher != nil {
		s.srv.Server().DocumentSyncher.WillSave(ctx, params)
	}
	return nil
}
func (s *DefaultLanguageServer) WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error) {
	if s.srv.Server().DocumentSyncher != nil {
		return s.srv.Server().DocumentSyncher.WillSaveWaitUntil(ctx, params)
	}
	return []lsp.TextEdit{}, nil
}
func (s *DefaultLanguageServer) Subtypes(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Supertypes(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) WorkDoneProgressCancel(ctx context.Context, params *lsp.WorkDoneProgressCancelParams) error {
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

	// Register build step listener to publish diagnostics after validation
	client := lsp.ClientDispatcher(conn)
	srv.Workspace().Builder.AddBuildStepListener(core.DocStateValidated, func(ctx context.Context, doc *core.Document) error {
		lspDiags := make([]lsp.Diagnostic, 0, len(doc.Diagnostics))
		for _, d := range doc.Diagnostics {
			lspDiags = append(lspDiags, toLspDiagnostic(*d))
		}
		params := &lsp.PublishDiagnosticsParams{
			URI:         doc.URI.DocumentURI(),
			Version:     doc.TextDoc.Version(),
			Diagnostics: lspDiags,
		}
		return client.PublishDiagnostics(ctx, params)
	})

	// Wait for the connection to close
	return conn.Wait()
}

func toLspDiagnostic(d core.Diagnostic) lsp.Diagnostic {
	result := lsp.Diagnostic{
		Range:    d.Range.LspRange(),
		Severity: lsp.DiagnosticSeverity(d.Severity),
		Message:  d.Message,
	}
	if d.Source != "" {
		result.Source = d.Source
	}
	if d.Code != "" {
		result.Code = d.Code
	}
	if d.CodeDescription != nil {
		result.CodeDescription = &lsp.CodeDescription{
			Href: d.CodeDescription.Href,
		}
	}
	if len(d.Tags) > 0 {
		tags := make([]lsp.DiagnosticTag, len(d.Tags))
		for i, t := range d.Tags {
			tags[i] = lsp.DiagnosticTag(t)
		}
		result.Tags = tags
	}
	if d.Data != nil {
		raw, err := json.Marshal(d.Data)
		if err == nil {
			rawMsg := json.RawMessage(raw)
			result.Data = &rawMsg
		}
	}
	return result
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
		Handler: lsp.ServerHandler(b.srv.Server().LanguageServer),
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
