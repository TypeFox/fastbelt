// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package glob provides a minimal glob matcher over slash-separated paths.
package glob

import (
	"path"
	"strings"
)

// Match reports whether name matches a glob pattern. Paths are matched
// segment-by-segment on '/'. Within a segment the [path.Match] syntax applies
// ('*', '?', character classes). The '**' segment matches zero or more whole
// segments (crossing separators).
func Match(pattern, name string) bool {
	return matchSegments(strings.Split(pattern, "/"), strings.Split(name, "/"))
}

func matchSegments(pat, name []string) bool {
	if len(pat) == 0 {
		return len(name) == 0
	}
	if pat[0] == "**" {
		// '**' consumes zero or more segments; try every split.
		for i := 0; i <= len(name); i++ {
			if matchSegments(pat[1:], name[i:]) {
				return true
			}
		}
		return false
	}
	if len(name) == 0 {
		return false
	}
	if ok, _ := path.Match(pat[0], name[0]); !ok {
		return false
	}
	return matchSegments(pat[1:], name[1:])
}
