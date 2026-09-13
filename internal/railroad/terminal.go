// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "math"

// Terminal renders a literal keyword/token as a rounded, filled box.
type Terminal struct {
	Text string
}

func (t *Terminal) Width() float64   { return math.Max(0, float64(len([]rune(t.Text))))*charWidth + 20 }
func (t *Terminal) Up() float64      { return 11 }
func (t *Terminal) Down() float64    { return 11 }
func (t *Terminal) needsSpace() bool { return true }

func (t *Terminal) render(x, y, width float64) *svgElement {
	return renderBox(x, y, width, t.Width(), t.Text, classTerminal, true)
}

// NonTerminal renders a reference to another rule as a square-cornered,
// filled box. Unresolved marks a reference that could not be resolved (e.g.
// a broken cross-reference in a document being edited), rendered with a
// distinguishing dashed style rather than failing.
type NonTerminal struct {
	Text       string
	Unresolved bool
}

func (n *NonTerminal) Width() float64 {
	return math.Max(0, float64(len([]rune(n.Text))))*charWidth + 20
}
func (n *NonTerminal) Up() float64      { return 11 }
func (n *NonTerminal) Down() float64    { return 11 }
func (n *NonTerminal) needsSpace() bool { return true }

func (n *NonTerminal) render(x, y, width float64) *svgElement {
	class := classNonTerminal
	if n.Unresolved {
		class = classNonTerminalUnresolved
	}
	return renderBox(x, y, width, n.Width(), n.Text, class, false)
}

func renderBox(x, y, width, boxWidth float64, text, class string, rounded bool) *svgElement {
	gapL, gapR := determineGaps(width, boxWidth)
	g := el("g")
	g.add(newPath(x, y).h(gapL).element())
	g.add(newPath(x+gapL+boxWidth, y).h(gapR).element())
	x += gapL

	rect := el("rect").
		attrf("x", x).
		attrf("y", y-11).
		attrf("width", boxWidth).
		attrf("height", 22).
		attr("class", class)
	if rounded {
		rect.attr("rx", "10").attr("ry", "10")
	}
	g.add(rect)
	g.add(el("text").
		attrf("x", x+boxWidth/2).
		attrf("y", y+4).
		attr("class", classBoxText).
		setText(text))
	return g
}

// Skip is a zero-width passthrough, used for the empty branch of Optional
// and for grammar elements (such as actions) that consume no tokens.
type Skip struct{}

func (Skip) Width() float64   { return 0 }
func (Skip) Up() float64      { return 0 }
func (Skip) Down() float64    { return 0 }
func (Skip) needsSpace() bool { return false }

func (Skip) render(x, y, width float64) *svgElement {
	return el("g").add(newPath(x, y).right(width).element())
}

// Start is the fixed marker at the beginning of a Diagram.
type Start struct{}

func (Start) Width() float64   { return 20 }
func (Start) Up() float64      { return 10 }
func (Start) Down() float64    { return 10 }
func (Start) needsSpace() bool { return false }

func (Start) render(x, y, _ float64) *svgElement {
	d := "M " + formatNum(x) + " " + formatNum(y-10) + " v20 m10 -20 v20 m-10 -10 h20.5"
	return el("path").attr("d", d)
}

// End is the fixed marker at the end of a Diagram.
type End struct{}

func (End) Width() float64   { return 20 }
func (End) Up() float64      { return 10 }
func (End) Down() float64    { return 10 }
func (End) needsSpace() bool { return false }

func (End) render(x, y, _ float64) *svgElement {
	d := "M " + formatNum(x) + " " + formatNum(y) + " h20 m-10 -10 v20 m10 -20 v20"
	return el("path").attr("d", d)
}
