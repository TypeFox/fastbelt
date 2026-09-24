// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"errors"
	"sort"

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
	for _, g := range siblingGrammars(grammar) {
		for _, iface := range g.Interfaces() {
			if iface.Name() == name {
				return iface
			}
		}
	}
	return nil
}

// siblingGrammars returns the grammar roots of every .fb file in the same
// folder as g, g itself included and always first, the rest sorted by URI so
// that folder-wide diagnostics are deterministic. All files of a folder form
// one grammar; see [folderGrammar].
func siblingGrammars(g Grammar) []Grammar {
	result := []Grammar{g}
	docs := siblingDocuments(g.Document())
	if docs == nil {
		return result
	}
	ownDoc := g.Document()
	var others []Grammar
	for doc := range docs {
		if doc == ownDoc {
			continue
		}
		if other, ok := doc.Root.(Grammar); ok {
			others = append(others, other)
		}
	}
	sort.Slice(others, func(i, j int) bool {
		return others[i].Document().URI.String() < others[j].Document().URI.String()
	})
	return append(result, others...)
}

// folderGrammar returns a detached grammar node that lists the declarations of
// every sibling of g (see [siblingGrammars]). The generated setters only append
// to slices, so the listed nodes keep their real container and document; the
// aggregate is a read-only view for checks that must consider the whole folder.
//
// ponytail: every document of a folder builds this view again during its own
// validation, which is quadratic in the number of files. A grammar has a handful
// of files at most; cache per build if that ever changes.
func folderGrammar(g Grammar) Grammar {
	siblings := siblingGrammars(g)
	if len(siblings) == 1 {
		return g
	}
	folder := NewGrammar()
	folder.SetName(g.NameToken())
	for _, s := range siblings {
		for _, item := range s.Rules() {
			folder.SetRulesItem(item)
		}
		for _, item := range s.Composites() {
			folder.SetCompositesItem(item)
		}
		for _, item := range s.InfixRules() {
			folder.SetInfixRulesItem(item)
		}
		for _, item := range s.Terminals() {
			folder.SetTerminalsItem(item)
		}
		for _, item := range s.TokenGroups() {
			folder.SetTokenGroupsItem(item)
		}
		for _, item := range s.TokenModes() {
			folder.SetTokenModesItem(item)
		}
		for _, item := range s.Interfaces() {
			folder.SetInterfacesItem(item)
		}
	}
	return folder
}

// ownedBy reports whether node belongs to doc. Folder-wide checks compute over
// the whole folder but only report on the nodes of the document being validated.
func ownedBy(doc *core.Document, node core.AstNode) bool {
	return node.Document() == doc
}
