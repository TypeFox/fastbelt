// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package json

import (
	"encoding/json"
	"fmt"
	"reflect"

	core "typefox.dev/fastbelt"
)

// Unmarshal decodes data into T by reading the "$type" field to select a factory from factories.
func Unmarshal[T any](data []byte, factories map[string]func() core.AstNode) (T, error) {
	node := &struct {
		Type string `json:"$type"`
	}{}
	if err := json.Unmarshal(data, node); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal: %w", err)
	}

	factory, ok := factories[node.Type]
	if !ok {
		var zero T
		return zero, fmt.Errorf("unmarshal: unknown type %q", node.Type)
	}

	instance := factory()
	casted, ok := instance.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("unmarshal: %T is not convertible to type %s", instance, reflect.TypeFor[T]())
	}

	if err := json.Unmarshal(data, casted); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal %s: %w", node.Type, err)
	}

	return casted, nil
}
