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
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal: %w", err)
	}

	typeRaw, ok := raw["$type"]
	if !ok {
		var zero T
		return zero, fmt.Errorf("unmarshal: missing $type field")
	}

	var typeName string
	if err := json.Unmarshal(typeRaw, &typeName); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal $type: %w", err)
	}

	factory, ok := factories[typeName]
	if !ok {
		var zero T
		return zero, fmt.Errorf("unmarshal: unknown type %q", typeName)
	}

	instance := factory()
	casted, ok := instance.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("unmarshal: %T is not convertible to type %s", instance, reflect.TypeFor[T]())
	}

	if err := json.Unmarshal(data, casted); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal %s: %w", typeName, err)
	}

	return casted, nil
}
