// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package railroad renders railroad syntax diagrams as SVG. It is a small,
// scoped Go port of the box-model layout algorithm from
// https://github.com/tabatkins/railroad-diagrams (railroad-diagrams@1.0.0,
// CC0), covering the primitives needed for context-free grammar rules:
// Sequence, Choice, Optional, ZeroOrMore, OneOrMore, Terminal, NonTerminal,
// and Skip.
package railroad

// Layout constants, ported verbatim from railroad-diagrams@1.0.0.
const (
	arcRadius          = 10.0
	verticalSeparation = 8.0
	charWidth          = 8.0
)

// Node is a railroad-diagram element: something that can report its box
// geometry (width, and vertical extent above/below its entry/exit
// centerline) and render itself into an SVG fragment at a given position.
type Node interface {
	// Width is the total horizontal extent of the node.
	Width() float64
	// Up is the vertical extent above the entry/exit centerline.
	Up() float64
	// Down is the vertical extent below the entry/exit centerline.
	Down() float64

	// needsSpace reports whether 20px of extra horizontal breathing room
	// should be reserved around this node when it's placed in a Sequence
	// or at the top level of a Diagram.
	needsSpace() bool
	// render draws the node at the given position, assuming the given
	// width (which may be wider than Width(), in which case the node is
	// centered within it).
	render(x, y, width float64) *svgElement
}
