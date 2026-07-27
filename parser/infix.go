// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parser

import (
	"math"

	core "typefox.dev/fastbelt"
)

// InfixPrecedence describes one operator token type of an infix rule.
type InfixPrecedence struct {
	// Level is the precedence group index of the operator.
	// Groups are declared tightest-binding first, so a larger level binds looser.
	Level int
	// Right reports whether the operator's precedence group is right-associative.
	Right bool
}

// BuildInfixTree folds the flat parts/operators lists collected by a generated
// infix rule parser into a binary expression tree.
//
// The tree is built by recursively splitting at the loosest-binding operator
// (the one with the largest level). Between operators of the same level, the
// rightmost one becomes the split point for left-associative groups (yielding
// a left-leaning tree) and the leftmost one for right-associative groups.
//
// precedence is keyed by the leaf token type id; an operator
// token missing from the map is treated as binding loosest. build constructs a
// single node; right is the zero value of T when a trailing operator has no
// operand (parse error tolerance). A single part is returned unchanged.
func BuildInfixTree[T core.AstNode](parts []T, operators []*core.Token,
	precedence map[int]InfixPrecedence,
	build func(left T, operator *core.Token, right T) T) T {
	var zero T
	if len(parts) == 0 {
		return zero
	}
	// rec builds the tree for operators[lo:hi] and their surrounding parts.
	var rec func(lo, hi int) T
	rec = func(lo, hi int) T {
		if lo == hi {
			if lo < len(parts) {
				return parts[lo]
			}
			return zero
		}
		pivot := lo
		bestLevel := -1
		for i := lo; i < hi; i++ {
			level, right := math.MaxInt, false
			if p, ok := precedence[operators[i].Type.Id]; ok {
				level, right = p.Level, p.Right
			}
			if level > bestLevel {
				bestLevel = level
				pivot = i
			} else if level == bestLevel && !right {
				// Left-associative: the rightmost operator of the loosest
				// level becomes the root. Right-associative: keep the leftmost.
				pivot = i
			}
		}
		return build(rec(lo, pivot), operators[pivot], rec(pivot+1, hi))
	}
	return rec(0, len(operators))
}
