// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package glob provides a minimal glob matcher over slash-separated paths.
package glob

import "strings"

// Match reports whether path matches a glob pattern. Paths are matched
// segment-by-segment on '/'. Within a segment, '*' matches any run of
// non-separator characters and '?' matches a single one. The '**' segment
// matches zero or more whole segments (crossing separators).
//
// ponytail: minimal glob covering the "**/*.ext" defaults + simple hand-written
// patterns; swap for a glob lib if callers ever need {a,b} / [..] classes.
func Match(pattern, path string) bool {
	return matchSegments(strings.Split(pattern, "/"), strings.Split(path, "/"))
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
	if !matchSegment(pat[0], name[0]) {
		return false
	}
	return matchSegments(pat[1:], name[1:])
}

// matchSegment matches a single path segment with '*' and '?' wildcards.
// Iterative backtracking on bytes ('?' matches any byte; the segment never
// contains a '/').
func matchSegment(p, s string) bool {
	pi, si := 0, 0
	star, ss := -1, 0
	for si < len(s) {
		switch {
		case pi < len(p) && (p[pi] == s[si] || p[pi] == '?'):
			pi++
			si++
		case pi < len(p) && p[pi] == '*':
			star = pi
			ss = si
			pi++
		case star != -1:
			pi = star + 1
			ss++
			si = ss
		default:
			return false
		}
	}
	for pi < len(p) && p[pi] == '*' {
		pi++
	}
	return pi == len(p)
}
