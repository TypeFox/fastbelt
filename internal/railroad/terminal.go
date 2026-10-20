// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "unicode/utf8"

const boxHalfHeight = 11.0

// Terminal renders a literal keyword/token as a rounded, filled box.
type Terminal struct {
	Text string
}

func (t *Terminal) Width() float64   { return boxWidth(t.Text) }
func (t *Terminal) Up() float64      { return boxHalfHeight }
func (t *Terminal) Down() float64    { return boxHalfHeight }
func (t *Terminal) needsSpace() bool { return true }

func (t *Terminal) render(x, y, width float64) *svgElement {
	return renderBox(x, y, width, t.Text, classTerminal, true)
}

// NonTerminal renders a reference to another rule as a square-cornered,
// filled box. Unresolved marks a reference that could not be resolved (e.g.
// a broken cross-reference in a document being edited), rendered with a
// distinguishing dashed style rather than failing.
type NonTerminal struct {
	Text       string
	Unresolved bool
}

func (n *NonTerminal) Width() float64   { return boxWidth(n.Text) }
func (n *NonTerminal) Up() float64      { return boxHalfHeight }
func (n *NonTerminal) Down() float64    { return boxHalfHeight }
func (n *NonTerminal) needsSpace() bool { return true }

func (n *NonTerminal) render(x, y, width float64) *svgElement {
	class := classNonTerminal
	if n.Unresolved {
		class = classNonTerminalUnresolved
	}
	return renderBox(x, y, width, n.Text, class, false)
}

func boxWidth(text string) float64 {
	return float64(utf8.RuneCountInString(text))*charWidth + 20
}

func renderBox(x, y, width float64, text, class string, rounded bool) *svgElement {
	w := boxWidth(text)
	gap := (width - w) / 2
	g := el("g")
	g.add(newPath(x, y).h(gap).element())
	g.add(newPath(x+gap+w, y).h(gap).element())
	x += gap

	rect := el("rect").
		attrf("x", x).
		attrf("y", y-boxHalfHeight).
		attrf("width", w).
		attrf("height", boxHalfHeight*2).
		attr("class", class)
	if rounded {
		rect.attr("rx", "10").attr("ry", "10")
	}
	g.add(rect)
	g.add(el("text").
		attrf("x", x+w/2).
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
	return el("g").add(newPath(x, y).h(width).element())
}

// start is the fixed marker at the beginning of a Diagram.
type start struct{}

func (start) Width() float64   { return 20 }
func (start) Up() float64      { return 10 }
func (start) Down() float64    { return 10 }
func (start) needsSpace() bool { return false }

func (start) render(x, y, _ float64) *svgElement {
	d := "M " + formatNum(x) + " " + formatNum(y-10) + " v20 m10 -20 v20 m-10 -10 h20.5"
	return el("path").attr("d", d)
}

// end is the fixed marker at the end of a Diagram.
type end struct{}

func (end) Width() float64   { return 20 }
func (end) Up() float64      { return 10 }
func (end) Down() float64    { return 10 }
func (end) needsSpace() bool { return false }

func (end) render(x, y, _ float64) *svgElement {
	d := "M " + formatNum(x) + " " + formatNum(y) + " h20 m-10 -10 v20 m10 -20 v20"
	return el("path").attr("d", d)
}
