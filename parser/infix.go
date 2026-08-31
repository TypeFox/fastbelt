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
// The fold is the standard operator-precedence stack reduction: tighter-binding
// operators (smaller level) reduce before a looser incoming operator is pushed,
// so the loosest operator ends up at the root. Operators of the same level
// reduce eagerly when left-associative (yielding a left-leaning tree) and lazily
// when right-associative. Runs in O(len(operators)).
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
	operand := func(i int) T {
		if i < len(parts) {
			return parts[i]
		}
		return zero
	}
	type stackedOp struct {
		token *core.Token
		level int
	}
	operands := make([]T, 1, len(operators)+1)
	operands[0] = operand(0)
	ops := make([]stackedOp, 0, len(operators))
	reduce := func() {
		op := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		right := operands[len(operands)-1]
		operands = operands[:len(operands)-1]
		operands[len(operands)-1] = build(operands[len(operands)-1], op.token, right)
	}
	for i, token := range operators {
		level, rightAssoc := math.MaxInt, false
		if p, ok := precedence[token.Type.Id]; ok {
			level, rightAssoc = p.Level, p.Right
		}
		for len(ops) > 0 && (ops[len(ops)-1].level < level ||
			(ops[len(ops)-1].level == level && !rightAssoc)) {
			reduce()
		}
		ops = append(ops, stackedOp{token, level})
		operands = append(operands, operand(i+1))
	}
	for len(ops) > 0 {
		reduce()
	}
	return operands[0]
}
