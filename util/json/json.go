// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"context"
	"encoding/json"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// UnmarshalAndBuildDocument uses the "encoding/json" entry point to unmarshal rootNode based on the given data string and builds the document.
// Similar to parsing-based document loading the ast build up first including the reference, while the resolution of the references is done during the linking phase of the building process.
// For properly linking references to other documents, a helper object is attached to the context given to builder providing access to other documents.
func UnmarshalAndBuildDocument[T core.AstNode](sc *service.Container, document *core.Document, rootNode T, data []byte, ctx context.Context) error {
	documents, err := service.Get[workspace.DocumentManager](sc)
	if err != nil {
		return err
	}
	builder, err := service.Get[workspace.Builder](sc)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, rootNode); err != nil {
		return err
	}
	core.AssignContainers(document, rootNode)
	document.Root = rootNode
	document.State = core.DocStateParsed
	documents.Set(document)

	if err := builder.Build(
		context.WithValue(
			ctx,
			core.JsonLinkingHelperKey(),
			NewJsonLinkingHelper(documents),
		),
		[]*core.Document{document}, nil,
	); err != nil {
		return err
	}

	return nil
}

type defaultJsonLinkingHelper struct {
	documentManager workspace.DocumentManager
}

func NewJsonLinkingHelper(docs workspace.DocumentManager) core.JsonLinkingHelper {
	return defaultJsonLinkingHelper{docs}
}

func (h defaultJsonLinkingHelper) GetDocument(uri core.URI) *core.Document {
	return h.documentManager.Get(uri)
}
