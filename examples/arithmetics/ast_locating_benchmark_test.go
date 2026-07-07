// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package arithmetics

import (
	"testing"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/test"
)

// BenchmarkPathOf benchmarks [core.PathOf] at four containment depths:
// self (0 segments), shallow (1 segment), medium (3 segments), and deep (5 segments).
// Node pointers are resolved once before the timer starts so only PathOf itself is measured.
func BenchmarkPathOf(b *testing.B) {
	doc := loadPriceCalcDoc(b)
	module := test.MustFindNode[Module](doc)
	stmts := module.Statements()

	// /statements@3/expression/left
	def3 := stmts[3].(Definition)
	bin3 := def3.Expression().(BinaryExpression)

	// /statements@7/expression/left/left/left
	def7 := stmts[7].(Definition)
	outerAdd := def7.Expression().(BinaryExpression)
	div := outerAdd.Left().(BinaryExpression)
	innerAdd := div.Left().(BinaryExpression)

	cases := []struct {
		name string
		node core.AstNode
	}{
		{"self/0-segments", module},
		{"shallow/1-segment", stmts[3]},
		{"medium/3-segments", bin3.Left()},
		{"deep/5-segments", innerAdd.Left()},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = core.PathOf(tc.node)
			}
		})
	}
}

// BenchmarkResolve benchmarks [core.Resolve] at four path depths:
// self (0 segments), shallow (1 segment), medium (3 segments), and deep (5 segments).
// The root node is obtained once before the timer starts.
func BenchmarkResolve(b *testing.B) {
	doc := loadPriceCalcDoc(b)
	root := doc.Root()

	cases := []struct {
		name string
		path string
	}{
		{"self/0-segments", ""},
		{"self/single-slash", "/"},
		{"shallow/1-segment", "/statements@3"},
		{"medium/3-segments", "/statements@3/expression/left"},
		{"deep/5-segments", "/statements@7/expression/left/left/left"},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_, _ = core.Resolve(tc.path, root)
			}
		})
	}
}
