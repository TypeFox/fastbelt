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

type infixOp struct {
	token *core.Token
	level int
}

// InfixFold folds the operands and operator tokens of a generated infix rule
// into a binary expression tree while they are being parsed.
//
// The fold is the standard operator-precedence stack reduction: tighter-binding
// operators (smaller level) reduce before a looser incoming operator is pushed,
// so the loosest operator ends up at the root. Operators of the same level
// reduce eagerly when left-associative (yielding a left-leaning tree) and lazily
// when right-associative.
//
// A single InfixFold lives on the generated parser per infix rule and is reused
// by every invocation of that rule, including nested ones (an operand that
// itself contains the infix rule, e.g. inside parentheses): each invocation
// works above the [InfixFrame] returned by Begin and leaves the stacks as it
// found them. The zero value is ready to use.
//
// Operands and operators must alternate, starting with an operand; a missing
// operand (parse error tolerance) is substituted with the zero value of T.
type InfixFold[T core.AstNode] struct {
	operands []T
	ops      []infixOp
}

// InfixFrame marks the stack depths at the start of one infix rule invocation.
type InfixFrame struct {
	operands, ops int
}

// Begin starts a new invocation and returns its frame.
func (f *InfixFold[T]) Begin() InfixFrame {
	return InfixFrame{operands: len(f.operands), ops: len(f.ops)}
}

// Operand pushes the next operand.
func (f *InfixFold[T]) Operand(node T) {
	f.operands = append(f.operands, node)
}

// Operator pushes the next operator token, first reducing every operator of the
// current frame that binds at least as tight (or exactly as tight and
// left-associative). table is indexed by the operator's token type id; a nil
// token or an id outside the table binds loosest. build constructs one binary
// node from a reduced operator and its two operands.
func (f *InfixFold[T]) Operator(frame InfixFrame, token *core.Token, table []InfixPrecedence, build func(left T, operator *core.Token, right T) T) {
	level, right := math.MaxInt, false
	if token != nil && token.Type.Id < len(table) {
		p := table[token.Type.Id]
		level, right = p.Level, p.Right
	}
	if len(f.operands)-frame.operands <= len(f.ops)-frame.ops {
		// No operand preceded this operator: substitute the zero value.
		var zero T
		f.operands = append(f.operands, zero)
	}
	for len(f.ops) > frame.ops {
		top := f.ops[len(f.ops)-1]
		if top.level < level || (top.level == level && !right) {
			f.reduce(build)
		} else {
			break
		}
	}
	f.ops = append(f.ops, infixOp{token: token, level: level})
}

// End reduces every remaining operator of the frame, pops the frame, and
// returns the folded tree. A single operand is returned unchanged; no operands
// yield the zero value of T.
func (f *InfixFold[T]) End(frame InfixFrame, build func(left T, operator *core.Token, right T) T) T {
	var zero T
	if len(f.operands)-frame.operands <= len(f.ops)-frame.ops {
		// A trailing operator has no right operand: substitute the zero value.
		f.operands = append(f.operands, zero)
	}
	for len(f.ops) > frame.ops {
		f.reduce(build)
	}
	if len(f.operands) == frame.operands {
		return zero
	}
	result := f.operands[frame.operands]
	// Clear the popped slots so the stack does not retain nodes of finished invocations.
	clear(f.operands[frame.operands:])
	f.operands = f.operands[:frame.operands]
	return result
}

func (f *InfixFold[T]) reduce(build func(left T, operator *core.Token, right T) T) {
	op := f.ops[len(f.ops)-1]
	f.ops = f.ops[:len(f.ops)-1]
	n := len(f.operands)
	right := f.operands[n-1]
	var zero T
	f.operands[n-1] = zero
	f.operands = f.operands[:n-1]
	f.operands[n-2] = build(f.operands[n-2], op.token, right)
}
