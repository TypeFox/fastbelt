// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package glob

import "testing"

func TestMatch(t *testing.T) {
	cases := []struct {
		pattern string
		path    string
		want    bool
	}{
		{"**/*.some", "/C:/dir/file.some", true},
		{"**/*.some", "/C:/dir/sub/file.some", true},
		{"**/*.some", "file.some", true},
		{"**/*.some", "/C:/dir/file.other", false},
		{"*.some", "file.some", true},
		{"*.some", "dir/file.some", false},
		{"a?c", "abc", true},
		{"a?c", "ac", false},
		{"**", "any/path/here", true},
		{"src/**/*.go", "src/a/b/c.go", true},
		{"src/**/*.go", "test/a.go", false},
	}
	for _, c := range cases {
		if got := Match(c.pattern, c.path); got != c.want {
			t.Errorf("Match(%q, %q) = %v, want %v", c.pattern, c.path, got, c.want)
		}
	}
}
