// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import "unicode"

// FuzzyMatcher performs approximate string matching between a user-typed
// query and a candidate text. It is a service-container service so any
// LSP feature that needs to filter labels by a typed prefix (completion,
// symbol search) can share one implementation across
// the language server.
type FuzzyMatcher interface {
	// Match reports whether query approximately matches text. An empty
	// query always matches.
	Match(query, text string) bool
}

// DefaultFuzzyMatcher is the framework's reference FuzzyMatcher
// implementation. Adopters can override it by registering a different
// FuzzyMatcher in the service container; the case-fold and
// word-transition helpers are exported as methods on the receiver so
// embedded types can selectively replace just those pieces.
type DefaultFuzzyMatcher struct{}

// NewDefaultFuzzyMatcher returns a ready-to-use FuzzyMatcher.
func NewDefaultFuzzyMatcher() FuzzyMatcher {
	return &DefaultFuzzyMatcher{}
}

// Match implements FuzzyMatcher. See the package-level interface doc for
// the matching rules.
func (m *DefaultFuzzyMatcher) Match(query, text string) bool {
	if query == "" {
		return true
	}
	// Casting to []rune automatically decodes UTF-8
	queryRunes := []rune(query)
	matchedFirst := false
	var previous rune = -1
	charIdx := 0
	for _, r := range text {
		q := queryRunes[charIdx]
		if r == q || unicode.ToUpper(r) == unicode.ToUpper(q) {
			if !matchedFirst && (previous == -1 || m.IsWordTransition(previous, r)) {
				matchedFirst = true
			}
			if matchedFirst {
				charIdx++
				if charIdx == len(queryRunes) {
					return true
				}
			}
		}
		previous = r
	}
	return false
}

// IsWordTransition reports whether the boundary between previous and
// current marks the start of a new "word" for anchoring: either a
// camelCase boundary (any Unicode lower-case letter followed by any
// Unicode upper-case letter) or a snake_case boundary (`_` followed
// by anything non-`_`).
func (m *DefaultFuzzyMatcher) IsWordTransition(previous, current rune) bool {
	if unicode.IsLower(previous) && unicode.IsUpper(current) {
		return true
	}
	return previous == '_' && current != '_'
}
