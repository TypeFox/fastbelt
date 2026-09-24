// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"errors"

	core "typefox.dev/fastbelt"
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
		value = KeywordValue(keyword)
	}
	return value, nil
}

func FindReturnType(rule AbstractRuleWithReturnType, ctx context.Context) Interface {
	if rule == nil {
		return nil
	}
	typeRef := rule.ReturnType()
	if typeRef != nil {
		return typeRef.Ref(ctx)
	}
	if infixRule, ok := rule.(InfixRule); ok {
		// The return type of an infix rule is the return type of its operand rule.
		if operand, ok := infixRule.Call().Rule().Ref(ctx).(ParserRule); ok {
			return FindReturnType(operand, ctx)
		}
		// Indicates malformed input, no return type
		return nil
	}
	grammar, ok := rule.Container().(Grammar)
	if !ok || grammar == nil {
		return nil
	}
	return FindInterfaceByName(grammar, rule.Name())
}

// FindInterfaceByName looks up an interface by name in the folder-wide grammar
// that grammar belongs to (see [siblingGrammars]).
func FindInterfaceByName(grammar Grammar, name string) Interface {
	if view := viewOf(grammar); view != nil {
		return view.interfaces[name]
	}
	for _, iface := range grammar.Interfaces() {
		if iface.Name() == name {
			return iface
		}
	}
	return nil
}

// siblingGrammars returns the grammar roots of every .fb file in the same
// folder as g, g itself included, sorted by URI so that folder-wide
// diagnostics are deterministic. All files of a folder form one grammar; see
// [folderGrammar].
func siblingGrammars(g Grammar) []Grammar {
	if view := viewOf(g); view != nil {
		return view.grammars
	}
	return []Grammar{g}
}

// folderGrammar returns a grammar node that lists the declarations of every
// sibling of g (see [siblingGrammars]): g itself for a single-file folder, and
// otherwise the read-only aggregate that the importer builds once per folder.
func folderGrammar(g Grammar) Grammar {
	if view := viewOf(g); view != nil {
		return view.folder
	}
	return g
}

// ownedBy reports whether node belongs to doc. Folder-wide checks compute over
// the whole folder but only report on the nodes of the document being validated.
func ownedBy(doc *core.Document, node core.AstNode) bool {
	return node.Document() == doc
}
