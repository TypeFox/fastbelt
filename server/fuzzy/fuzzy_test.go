// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fuzzy

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		text     string
		expected bool
	}{
		{
			name:     "Empty query matches",
			query:    "",
			text:     "Hello",
			expected: true,
		},

		{
			name:     "Matches full string",
			query:    "Hello",
			text:     "Hello",
			expected: true,
		},

		{
			name:     "Matches first few characters",
			query:    "He",
			text:     "Hello",
			expected: true,
		},

		{
			name:     "Does not match non-word-boundary",
			query:    "ell",
			text:     "Hello",
			expected: false,
		},

		{
			name:     "Matches omitted characters",
			query:    "Ho",
			text:     "Hello",
			expected: true,
		},

		{
			name:     "Does not match wrong omitted characters",
			query:    "Hi",
			text:     "Hello",
			expected: false,
		},

		{
			name:     "Matches first few characters (word transition, camelCase)",
			query:    "fb",
			text:     "helloFastbelt",
			expected: true,
		},

		{
			name:     "Matches first few characters (word transition, snake_case)",
			query:    "fb",
			text:     "hello_fastbelt",
			expected: true,
		},

		{
			name:     "Matches omitted characters (word transition, camelCase)",
			query:    "fb",
			text:     "helloFastbelt",
			expected: true,
		},

		{
			name:     "Matches omitted characters (word transition, snake_case)",
			query:    "fb",
			text:     "hello_fastbelt",
			expected: true,
		},

		{
			name:     "Case-insensitive match",
			query:    "hello",
			text:     "Hello",
			expected: true,
		},
		{
			name:     "Case-insensitive match uppercase query",
			query:    "HELLO",
			text:     "Hello",
			expected: true,
		},

		{
			name:     "FuMa matches FuzzyMatcher",
			query:    "FuMa",
			text:     "FuzzyMatcher",
			expected: true,
		},
		{
			name:     "PA matches PersonAddress",
			query:    "PA",
			text:     "PersonAddress",
			expected: true,
		},
		{
			name:     "PeAd matches PersonAddress",
			query:    "PeAd",
			text:     "PersonAddress",
			expected: true,
		},

		{
			name:     "Unicode: café matches café",
			query:    "café",
			text:     "café",
			expected: true,
		},
		{
			name:     "Unicode: case-insensitive Café",
			query:    "café",
			text:     "Café",
			expected: true,
		},
		{
			name:     "Unicode: German umlaut Über",
			query:    "über",
			text:     "Überschrift",
			expected: true,
		},
		{
			name:     "Unicode: Spanish ñ",
			query:    "señ",
			text:     "Señor",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Match(tt.query, tt.text)
			if result != tt.expected {
				t.Errorf("Match(%q, %q) = %v, expected %v", tt.query, tt.text, result, tt.expected)
			}
		})
	}
}
