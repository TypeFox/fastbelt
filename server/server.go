// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"encoding/json"
	"log"

	"golang.org/x/exp/jsonrpc2"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
	"typefox.dev/lsp"
)

// InitializeParticipant is an interface for services that want to
// participate in the LSP server initialization process.
//
// All services implementing this interface and registered in the [service.Container]
// instance will have their [InitializeParticipant.OnServerInitialize]
// method called during the [lsp.Server.Initialize] request.
type InitializeParticipant interface {
	OnServerInitialize(params *lsp.ParamInitialize)
}

// DefaultLanguageServer implements the [lsp.Server] interface
type DefaultLanguageServer struct {
	sc *service.Container
}

// NewDefaultLanguageServer creates a new default language server.
func NewDefaultLanguageServer(sc *service.Container) lsp.Server {
	return &DefaultLanguageServer{sc: sc}
}

func (s *DefaultLanguageServer) Initialize(ctx context.Context, params *lsp.ParamInitialize) (*lsp.InitializeResult, error) {
	// Initialize all participants first
	for service := range service.GetAll[InitializeParticipant](s.sc) {
		service.OnServerInitialize(params)
	}
	workspaceFolders, err := service.Get[*WorkspaceFolders](s.sc)
	if err != nil {
		return nil, err
	}
	workspaceFolders.Value = params.WorkspaceFolders
	var triggerChars []string
	if triggers, err := service.Get[CompletionTriggers](s.sc); err == nil && triggers != nil {
		triggerChars = triggers.TriggerCharacters()
	}
	positionEncoding := lsp.UTF16
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			PositionEncoding: &positionEncoding,
			TextDocumentSync: lsp.Incremental,
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: triggerChars,
			},
			DefinitionProvider: &lsp.Or_ServerCapabilities_definitionProvider{
				Value: service.Has[DefinitionProvider](s.sc),
			},
			DocumentSymbolProvider: &lsp.Or_ServerCapabilities_documentSymbolProvider{
				Value: service.Has[DocumentSymbolProvider](s.sc),
			},
			FoldingRangeProvider: &lsp.Or_ServerCapabilities_foldingRangeProvider{
				Value: service.Has[FoldingRangeProvider](s.sc),
			},
			DocumentHighlightProvider: &lsp.Or_ServerCapabilities_documentHighlightProvider{
				Value: service.Has[DocumentHighlightProvider](s.sc),
			},
			WorkspaceSymbolProvider: &lsp.Or_ServerCapabilities_workspaceSymbolProvider{
				Value: service.Has[WorkspaceSymbolProvider](s.sc),
			},
			HoverProvider: &lsp.Or_ServerCapabilities_hoverProvider{
				Value: service.Has[HoverProvider](s.sc),
			},
			ReferencesProvider: &lsp.Or_ServerCapabilities_referencesProvider{
				Value: service.Has[ReferencesProvider](s.sc),
			},
			RenameProvider: service.Has[RenameProvider](s.sc),
			DeclarationProvider: &lsp.Or_ServerCapabilities_declarationProvider{
				Value: service.Has[DeclarationProvider](s.sc),
			},
			ImplementationProvider: &lsp.Or_ServerCapabilities_implementationProvider{
				Value: service.Has[ImplementationProvider](s.sc),
			},
			TypeDefinitionProvider: &lsp.Or_ServerCapabilities_typeDefinitionProvider{
				Value: service.Has[TypeDefinitionProvider](s.sc),
			},
			SemanticTokensProvider: buildSemanticTokensOptions(s.sc),
			CallHierarchyProvider: &lsp.Or_ServerCapabilities_callHierarchyProvider{
				Value: service.Has[CallHierarchyProvider](s.sc),
			},
			TypeHierarchyProvider: &lsp.Or_ServerCapabilities_typeHierarchyProvider{
				Value: service.Has[TypeHierarchyProvider](s.sc),
			},
			InlayHintProvider: func() *lsp.Or_ServerCapabilities_inlayHintProvider {
				if service.Has[InlayHintProvider](s.sc) {
					return &lsp.Or_ServerCapabilities_inlayHintProvider{Value: &lsp.InlayHintOptions{}}
				}
				return nil
			}(),
			SignatureHelpProvider: buildSignatureHelpOptions(s.sc),
			CodeActionProvider: func() *lsp.CodeActionOptions {
				if service.Has[CodeActionProvider](s.sc) {
					return &lsp.CodeActionOptions{}
				}
				return nil
			}(),
			CodeLensProvider: func() *lsp.CodeLensOptions {
				if service.Has[CodeLensProvider](s.sc) {
					return &lsp.CodeLensOptions{ResolveProvider: false}
				}
				return nil
			}(),
			DocumentLinkProvider: func() *lsp.DocumentLinkOptions {
				if service.Has[DocumentLinkProvider](s.sc) {
					return &lsp.DocumentLinkOptions{ResolveProvider: false}
				}
				return nil
			}(),
			ExecuteCommandProvider: func() *lsp.ExecuteCommandOptions {
				if service.Has[CommandProvider](s.sc) {
					return &lsp.ExecuteCommandOptions{Commands: []string{}}
				}
				return nil
			}(),
		},
	}, nil
}

func buildSemanticTokensOptions(sc *service.Container) *lsp.SemanticTokensOptions {
	if !service.Has[SemanticTokensProvider](sc) {
		return nil
	}
	contributor, err := service.Get[SemanticTokensContributor](sc)
	if err != nil || len(contributor.TokenTypes()) == 0 {
		return nil
	}
	return &lsp.SemanticTokensOptions{
		Legend: lsp.SemanticTokensLegend{
			TokenTypes:     contributor.TokenTypes(),
			TokenModifiers: contributor.TokenModifiers(),
		},
		Full:  &lsp.Or_SemanticTokensOptions_full{Value: true},
		Range: &lsp.Or_SemanticTokensOptions_range{Value: true},
	}
}

func buildSignatureHelpOptions(sc *service.Container) *lsp.SignatureHelpOptions {
	provider, err := service.Get[SignatureHelpProvider](sc)
	if err != nil {
		return nil
	}
	triggerChars := provider.TriggerCharacters()
	if len(triggerChars) == 0 {
		return nil
	}
	return &lsp.SignatureHelpOptions{
		TriggerCharacters:   triggerChars,
		RetriggerCharacters: provider.RetriggerCharacters(),
	}
}

func (s *DefaultLanguageServer) Initialized(ctx context.Context, params *lsp.InitializedParams) error {
	if initializer, err := service.Get[workspace.Initializer](s.sc); err == nil {
		log.Print("LS Initializer running...")
		defer log.Println("done.")
		workspaceFolders := service.MustGet[*WorkspaceFolders](s.sc).Value
		return initializer.Initialize(ctx, workspaceFolders)
	}
	return nil
}

func (s *DefaultLanguageServer) Shutdown(ctx context.Context) error {
	return nil
}

func (s *DefaultLanguageServer) Exit(ctx context.Context) error {
	// Close the connection to allow the server to exit
	if connection := service.MustGet[*Connection](s.sc); connection.Value != nil {
		return connection.Value.Close()
	}
	return nil
}

func (s *DefaultLanguageServer) DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) error {
	if documentSyncher, err := service.Get[DocumentSyncher](s.sc); err == nil {
		documentSyncher.DidOpen(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams) error {
	if documentSyncher, err := service.Get[DocumentSyncher](s.sc); err == nil {
		documentSyncher.DidChange(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) error {
	if documentSyncher, err := service.Get[DocumentSyncher](s.sc); err == nil {
		documentSyncher.DidClose(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams) error {
	if documentSyncher, err := service.Get[DocumentSyncher](s.sc); err == nil {
		documentSyncher.WillSave(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error) {
	if documentSyncher, err := service.Get[DocumentSyncher](s.sc); err == nil {
		return documentSyncher.WillSaveWaitUntil(ctx, params)
	}
	return nil, nil
}

func (s *DefaultLanguageServer) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
	if documentSyncher, err := service.Get[DocumentSyncher](s.sc); err == nil {
		documentSyncher.DidSave(ctx, params)
	}
	return nil
}

func (s *DefaultLanguageServer) Completion(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error) {
	var result *lsp.CompletionList
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	completion, err := service.Get[CompletionProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = completion.HandleCompletionRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}

func (s *DefaultLanguageServer) Definition(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.DefinitionLink, error) {
	var result []lsp.DefinitionLink
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	definition, err := service.Get[DefinitionProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = definition.HandleDefinitionRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}

func (s *DefaultLanguageServer) References(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error) {
	var result []lsp.Location
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	references, err := service.Get[ReferencesProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = references.HandleReferencesRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}

// Implement other required Server interface methods with no-op implementations

func (s *DefaultLanguageServer) Progress(ctx context.Context, params *lsp.ProgressParams) error {
	return nil
}
func (s *DefaultLanguageServer) SetTrace(ctx context.Context, params *lsp.SetTraceParams) error {
	return nil
}
func (s *DefaultLanguageServer) IncomingCalls(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
	var result []lsp.CallHierarchyIncomingCall
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[CallHierarchyProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleIncomingCallsRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) OutgoingCalls(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
	var result []lsp.CallHierarchyOutgoingCall
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[CallHierarchyProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleOutgoingCallsRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
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
	var result []lsp.CodeAction
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[CodeActionProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleCodeActionRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) CodeLens(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	var result []lsp.CodeLens
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[CodeLensProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleCodeLensRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) ColorPresentation(ctx context.Context, params *lsp.ColorPresentationParams) ([]lsp.ColorPresentation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Declaration(ctx context.Context, params *lsp.DeclarationParams) ([]lsp.DefinitionLink, error) {
	var result []lsp.DefinitionLink
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[DeclarationProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleDeclarationRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Diagnostic(ctx context.Context, params *lsp.DocumentDiagnosticParams) (*lsp.DocumentDiagnosticReport, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentColor(ctx context.Context, params *lsp.DocumentColorParams) ([]lsp.ColorInformation, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) DocumentHighlight(ctx context.Context, params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
	var result []lsp.DocumentHighlight
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[DocumentHighlightProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleDocumentHighlightRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) DocumentLink(ctx context.Context, params *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error) {
	var result []lsp.DocumentLink
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[DocumentLinkProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleDocumentLinkRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) DocumentSymbol(ctx context.Context, params *lsp.DocumentSymbolParams) ([]any, error) {
	var result []lsp.DocumentSymbol
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[DocumentSymbolProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleDocumentSymbolRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return toAnySlice(result), providerErr
}
func (s *DefaultLanguageServer) FoldingRange(ctx context.Context, params *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error) {
	var result []lsp.FoldingRange
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[FoldingRangeProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleFoldingRangeRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Formatting(ctx context.Context, params *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Hover(ctx context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	var result *lsp.Hover
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[HoverProvider](s.sc)
	if err != nil {
		return nil, err
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleHoverRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Implementation(ctx context.Context, params *lsp.ImplementationParams) ([]lsp.DefinitionLink, error) {
	var result []lsp.DefinitionLink
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[ImplementationProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleImplementationRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) InlayHint(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	var result []lsp.InlayHint
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[InlayHintProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleInlayHintRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
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
	provider, err := service.Get[CommandProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	return provider.HandleExecuteCommandRequest(ctx, params)
}
func (s *DefaultLanguageServer) Symbol(ctx context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	var result []lsp.SymbolInformation
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[WorkspaceSymbolProvider](s.sc)
	if err != nil {
		return nil, nil // No provider registered, return empty
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleWorkspaceSymbolRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
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
	var result []lsp.CallHierarchyItem
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[CallHierarchyProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandlePrepareCallHierarchyRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) PrepareRename(ctx context.Context, params *lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error) {
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	renameProvider, err := service.Get[RenameProvider](s.sc)
	if err != nil {
		return nil, err
	}
	var result *lsp.PrepareRenameResult
	var providerErr error
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = renameProvider.PrepareRenameRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) PrepareTypeHierarchy(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error) {
	var result []lsp.TypeHierarchyItem
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[TypeHierarchyProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandlePrepareTypeHierarchyRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) RangeFormatting(ctx context.Context, params *lsp.DocumentRangeFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) RangesFormatting(ctx context.Context, params *lsp.DocumentRangesFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) Rename(ctx context.Context, params *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	renameProvider, err := service.Get[RenameProvider](s.sc)
	if err != nil {
		return nil, err
	}
	var result *lsp.WorkspaceEdit
	var providerErr error
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = renameProvider.HandleRenameRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) SelectionRange(ctx context.Context, params *lsp.SelectionRangeParams) ([]lsp.SelectionRange, error) {
	return nil, nil
}
func (s *DefaultLanguageServer) SemanticTokensFull(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	var result *lsp.SemanticTokens
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[SemanticTokensProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleSemanticTokensFullRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) SemanticTokensFullDelta(ctx context.Context, params *lsp.SemanticTokensDeltaParams) (any, error) {
	var result any
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[SemanticTokensProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleSemanticTokensFullDeltaRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) SemanticTokensRange(ctx context.Context, params *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error) {
	var result *lsp.SemanticTokens
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[SemanticTokensProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleSemanticTokensRangeRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) SignatureHelp(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	var result *lsp.SignatureHelp
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[SignatureHelpProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleSignatureHelpRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) TypeDefinition(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.DefinitionLink, error) {
	var result []lsp.DefinitionLink
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[TypeDefinitionProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleTypeDefinitionRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Subtypes(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error) {
	var result []lsp.TypeHierarchyItem
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[TypeHierarchyProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleSubtypesRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) Supertypes(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error) {
	var result []lsp.TypeHierarchyItem
	var providerErr error
	lock, err := service.Get[workspace.Lock](s.sc)
	if err != nil {
		return nil, err
	}
	provider, err := service.Get[TypeHierarchyProvider](s.sc)
	if err != nil {
		return nil, nil
	}
	if err := lock.Read(ctx, func(ctx context.Context) {
		result, providerErr = provider.HandleSupertypesRequest(ctx, params)
	}); err != nil {
		return nil, err
	}
	return result, providerErr
}
func (s *DefaultLanguageServer) WorkDoneProgressCancel(ctx context.Context, params *lsp.WorkDoneProgressCancelParams) error {
	return nil
}

// StartLanguageServer starts a language server using the service container.
// It sets up JSON-RPC communication over stdio and handles the essential LSP messages.
func StartLanguageServer(ctx context.Context, sc *service.Container) error {
	dialer, err := service.Get[jsonrpc2.Dialer](sc)
	if err != nil {
		return err
	}
	binder, err := service.Get[jsonrpc2.Binder](sc)
	if err != nil {
		return err
	}

	// Create a connection using the configured dialer and binder
	conn, err := jsonrpc2.Dial(ctx, dialer, binder)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close() // Ignore error in defer
	}()

	// Wait for the connection to close
	return conn.Wait()
}

// toAnySlice converts []lsp.DocumentSymbol to []any for the LSP response.
func toAnySlice(symbols []lsp.DocumentSymbol) []any {
	result := make([]any, len(symbols))
	for i, sym := range symbols {
		result[i] = sym
	}
	return result
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
	sc *service.Container
}

// NewDefaultBinder creates a new default binder.
func NewDefaultBinder(sc *service.Container) *DefaultBinder {
	return &DefaultBinder{sc: sc}
}

func (b *DefaultBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	// Store the JSON-RPC connection in the service container
	connection, err := service.Get[*Connection](b.sc)
	if err != nil {
		return jsonrpc2.ConnectionOptions{}, err
	}
	connection.Value = conn
	// Bind the LSP server implementation from the service container to the connection
	server, err := service.Get[lsp.Server](b.sc)
	if err != nil {
		return jsonrpc2.ConnectionOptions{}, err
	}
	return jsonrpc2.ConnectionOptions{
		Handler: lsp.ServerHandler(server),
	}, nil
}
