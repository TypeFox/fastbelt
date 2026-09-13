// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "math"

const diagramPadding = 20.0

// Diagram is a complete railroad diagram for a single grammar rule: a Start
// marker, the rule's body, and an End marker, laid out left to right.
type Diagram struct {
	Item Node
}

func (d *Diagram) items() []Node {
	return []Node{Start{}, d.Item, End{}}
}

// SVG renders the diagram to a standalone SVG document: no XML declaration
// (unnecessary for an inline/embedded image), a light <style> block, and no
// background fill so it reads correctly on both light and dark hosts (see
// style.go).
func (d *Diagram) SVG() string {
	items := d.items()

	width := 1.0 // matches railroad-diagrams@1.0.0's Diagram width fudge
	up, down := 0.0, 0.0
	for _, it := range items {
		width += it.Width()
		if it.needsSpace() {
			width += 20
		}
		up = math.Max(up, it.Up())
		down = math.Max(down, it.Down())
	}

	x := diagramPadding
	y := diagramPadding + up

	g := el("g").attr("transform", "translate(.5 .5)")
	for _, it := range items {
		if it.needsSpace() {
			g.add(newPath(x, y).h(10).element())
			x += 10
		}
		g.add(it.render(x, y, it.Width()))
		x += it.Width()
		if it.needsSpace() {
			g.add(newPath(x, y).h(10).element())
			x += 10
		}
	}

	totalWidth := width + diagramPadding*2
	totalHeight := up + down + diagramPadding*2

	svg := el("svg").
		attr("xmlns", "http://www.w3.org/2000/svg").
		attrf("width", totalWidth).
		attrf("height", totalHeight).
		attr("viewBox", "0 0 "+formatNum(totalWidth)+" "+formatNum(totalHeight)).
		attr("class", "railroad-diagram").
		add(styleElement(), g)

	return svg.String()
}
