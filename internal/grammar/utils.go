// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"errors"
)

var errMissingKeywordValue = errors.New("keyword has no token value")

// convertString converts a keyword token value to its semantic string value.
//
// It strips surrounding double quotes when present.
// It returns errMissingKeywordValue only when the keyword has no token value
// (i.e. Keyword.Value() is an empty string).
// Quoted empty content (e.g. "\"\"") is valid and converts to "" without error.
func convertString(keyword Keyword) (string, error) {
	value := keyword.Value()
	if value == "" {
		return "", errMissingKeywordValue
	}
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}
	return value, nil
}

func FindReturnType(rule ParserRule, ctx context.Context) Interface {
	if rule == nil {
		return nil
	}
	typeRef := rule.ReturnType()
	if typeRef != nil {
		return typeRef.Ref(ctx)
	}
	grammar := rule.Container().(Grammar)
	if grammar == nil {
		return nil
	}
	return FindInterfaceByName(grammar, rule.Name())
}

func FindInterfaceByName(grammar Grammar, name string) Interface {
	for _, iface := range grammar.Interfaces() {
		if iface.Name() == name {
			return iface
		}
	}
	return nil
}
