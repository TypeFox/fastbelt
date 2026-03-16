package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/lsp"
)

type ReferencesProvider interface {
	HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error)
}

type DefaultReferencesProvider struct {
	srv ServerSrvCont
}

func NewDefaultReferencesProvider(srv ServerSrvCont) ReferencesProvider {
	return &DefaultReferencesProvider{srv: srv}
}

func (rp *DefaultReferencesProvider) HandleReferencesRequest(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error) {
	uri := core.ParseURI(string(params.TextDocument.URI))
	targetDoc := rp.srv.Workspace().DocumentManager.Get(uri)
	if targetDoc == nil {
		return nil, nil // Document not found
	}
	targetDoc.RLock()
	defer targetDoc.RUnlock()
	offset := targetDoc.TextDoc.OffsetAt(params.Position)
	tokens := targetDoc.Tokens
	// This token represents the name of the symbol at the given position
	sourceToken := tokens.SearchOffset(offset)
	if sourceToken == nil {
		return nil, nil // No token at the given position
	}
	// We now need to figure out whether this token is really the name token
	// It could also be just a random keyword or something else
	owner := sourceToken.Element
	if owner == nil {
		return nil, nil // No AST node associated with the token
	}
	namer := rp.srv.Linking().Namer
	_, nameToken := namer.Name(owner)
	if nameToken == nil || nameToken != sourceToken {
		return nil, nil // The token at the position is not the name token
	}
	locations := []lsp.Location{
		// Include the definition location itself
		{
			URI:   owner.Document().URI.DocumentURI(),
			Range: nameToken.Segment.Range.LspRange(),
		},
	}
	documentManager := rp.srv.Workspace().DocumentManager
	// Iterate through all documents and collect references to the symbol
	for doc := range documentManager.All() {
		refDescriptions := doc.ReferenceDescriptions.ForTarget(owner)
		for refDesc := range refDescriptions {
			location := lsp.Location{
				URI:   refDesc.SourceURI().DocumentURI(),
				Range: refDesc.Segment.Range.LspRange(),
			}
			locations = append(locations, location)
		}
	}
	return locations, nil
}
