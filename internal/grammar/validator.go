// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	core "typefox.dev/fastbelt"
)

// GrammarImpl.Validate checks grammar-level constraints:
//   - Rule names must be unique within the grammar.
//   - Interface names must be unique within the grammar.
func (g *GrammarImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkUniqueRuleNames(g, accept)
	checkUniqueInterfaceNames(g, accept)
}

func checkUniqueRuleNames(g *GrammarImpl, accept core.ValidationAcceptor) {
	seen := map[string][]core.NamedNode{}
	for _, rule := range g.Rules() {
		if rule.Name() != "" {
			seen[rule.Name()] = append(seen[rule.Name()], rule)
		}
	}
	for _, terminal := range g.Terminals() {
		if terminal.Name() != "" {
			seen[terminal.Name()] = append(seen[terminal.Name()], terminal)
		}
	}
	for name, nodes := range seen {
		if len(nodes) > 1 {
			for _, node := range nodes {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("A rule's name has to be unique. '%s' is used multiple times.", name),
					node,
					core.WithToken(node.NameToken()),
				))
			}
		}
	}
}

func checkUniqueInterfaceNames(g *GrammarImpl, accept core.ValidationAcceptor) {
	seen := map[string][]Interface{}
	for _, iface := range g.Interfaces() {
		if iface.Name() != "" {
			seen[iface.Name()] = append(seen[iface.Name()], iface)
		}
	}
	for name, ifaces := range seen {
		if len(ifaces) > 1 {
			for _, iface := range ifaces {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("An interface name has to be unique. '%s' is used multiple times.", name),
					iface,
					core.WithToken(iface.NameToken()),
				))
			}
		}
	}
}

// TokenImpl.Validate checks terminal rule constraints:
//   - The regular expression should not match the empty string.
func (t *TokenImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkEmptyTerminalRule(t, accept)
}

func checkEmptyTerminalRule(t *TokenImpl, accept core.ValidationAcceptor) {
	raw := t.Regexp()
	if raw == "" {
		return
	}
	// Strip surrounding slashes from the regex literal
	pattern := raw
	if len(pattern) >= 2 && pattern[0] == '/' && pattern[len(pattern)-1] == '/' {
		pattern = pattern[1 : len(pattern)-1]
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}
	if re.MatchString("") {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"This terminal could match an empty string.",
			t,
			core.WithToken(t.NameToken()),
		))
	}
}

// KeywordImpl.Validate checks keyword constraints:
//   - Keywords cannot be empty.
//   - Keywords cannot consist only of whitespace.
//   - Keywords should not contain whitespace characters (warning).
func (k *KeywordImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkKeyword(k, accept)
}

func checkKeyword(k *KeywordImpl, accept core.ValidationAcceptor) {
	raw := k.Value()
	if raw == "" {
		return
	}
	// Strip surrounding double quotes from the string literal
	value := raw
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}
	if len(value) == 0 {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Keywords cannot be empty.",
			k,
			core.WithToken(k.ValueToken()),
		))
	} else if strings.TrimSpace(value) == "" {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Keywords cannot only consist of whitespace characters.",
			k,
			core.WithToken(k.ValueToken()),
		))
	} else if strings.ContainsAny(value, " \t\n\r") {
		accept(core.NewDiagnostic(
			core.SeverityWarning,
			"Keywords should not contain whitespace characters.",
			k,
			core.WithToken(k.ValueToken()),
		))
	}
}
