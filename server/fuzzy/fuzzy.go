// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fuzzy

import "unicode"

// Match performs fuzzy matching on the query and text.
//
// Fuzzy matching improves search/completion user experience by allowing to omit characters.
// For example, a query such as "FuMa" matches the text "FuzzyMatcher".
//
// Algorithm:
//   - Empty query matches everything
//   - Query characters must appear in order in text
//   - First query character must match at word boundary (start, camelCase, snake_case)
//   - Case insensitive (Unicode-aware)
//   - Allows omitted characters
func Match(query, text string) bool {
	if query == "" {
		return true
	}

	queryRunes := []rune(query)
	matchedFirst := false
	previous := rune(-1) // Sentinel value for start of text
	charIdx := 0

	for _, r := range text {
		q := queryRunes[charIdx]
		// Case-insensitive comparison using Unicode support
		if r == q || unicode.ToUpper(r) == unicode.ToUpper(q) {
			if !matchedFirst {
				// First character must match at word boundary
				if previous == -1 || isWordTransition(previous, r) {
					matchedFirst = true
				}
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

// isWordTransition checks if there's a word boundary between previous and current character.
// Returns true for:
//   - camelCase transition (lowercase -> uppercase): e.g., "myVariable"
//   - snake_case transition (underscore -> non-underscore): e.g., "my_variable"
func isWordTransition(previous, current rune) bool {
	// camelCase transition: any Unicode lowercase -> any Unicode uppercase
	if unicode.IsLower(previous) && unicode.IsUpper(current) {
		return true
	}
	// snake_case transition: underscore -> non-underscore
	if previous == '_' && current != '_' {
		return true
	}
	return false
}
