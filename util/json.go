// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package util

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"reflect"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// UnmarshalDecode reads the next JSON value from decoder into an instance of type T by
// reading the "$type" field, selecting a corresponding factory, creating an instance, and
// unmarshaling its content via the decoder API of package "encoding/json/v2".
func UnmarshalDecode[T core.AstNode](decoder *jsontext.Decoder, factories map[string]func() core.AstNode) (T, error) {
	var zero T
	raw, err := decoder.ReadValue()
	if err != nil {
		return zero, fmt.Errorf("unmarshalDecode: %w", err)
	}
	node := &struct {
		Type string `json:"$type"`
	}{}
	if err := json.Unmarshal(raw, node); err != nil {
		return zero, fmt.Errorf("unmarshalDecode: %w", err)
	}
	factory, ok := factories[node.Type]
	if !ok {
		return zero, fmt.Errorf("unmarshalDecode: unknown type %q", node.Type)
	}
	instance := factory()
	asT, ok := instance.(T)
	if !ok {
		return zero, fmt.Errorf("unmarshalDecode: %T is not convertible to type %s", instance, reflect.TypeFor[T]())
	}
	if unmarshaler, ok := instance.(json.UnmarshalerFrom); ok {
		if err := unmarshaler.UnmarshalJSONFrom(jsontext.NewDecoder(bytes.NewReader(raw))); err != nil {
			return zero, fmt.Errorf("unmarshalDecode %s: %w", node.Type, err)
		}
	} else {
		return zero, fmt.Errorf("unmarshalDecode: %T is not convertible to type json.UnmarshalerFrom", instance)
	}
	return asT, nil
}

// UnmarshalAndBuildDocument uses the "encoding/json/v2" entry point to unmarshal rootNode based on the given data string and builds the document.
// Similar to parsing-based document loading the ast build up first including the reference, while the resolution of the references is done during the linking phase of the building process.
// For properly linking references to other documents, a helper object is attached to the context given to builder providing access to other documents.
func UnmarshalAndBuildDocument(ctx context.Context, sc *service.Container, document *core.Document, input io.Reader, factories map[string]func() core.AstNode) error {
	documents, err := service.Get[workspace.DocumentManager](sc)
	if err != nil {
		return err
	}
	builder, err := service.Get[workspace.Builder](sc)
	if err != nil {
		return err
	}

	if document.Root, err = UnmarshalDecode[core.AstNode](jsontext.NewDecoder(input), factories); err != nil {
		return err
	}

	document.State = core.DocStateParsed
	core.AssignContainers(document)
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

func (h defaultJsonLinkingHelper) GetDocument(uri core.URI) *core.Document {
	return h.documentManager.Get(uri)
}

func NewJsonLinkingHelper(docs workspace.DocumentManager) core.JsonLinkingHelper {
	return defaultJsonLinkingHelper{docs}
}
