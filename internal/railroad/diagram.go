// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

const diagramPadding = 20.0

// Diagram is a complete railroad diagram for a single grammar rule: a start
// marker, the rule's body, and an end marker, laid out left to right.
type Diagram struct {
	Item Node
}

// SVG renders the diagram to a standalone SVG document: no XML declaration
// (unnecessary for an inline/embedded image), a light <style> block, and no
// background fill so it reads correctly on both light and dark hosts (see
// style.go).
func (d *Diagram) SVG() string {
	seq := &Sequence{Items: []Node{start{}, d.Item, end{}}}
	width := seq.Width()
	up := seq.Up()

	g := el("g").attr("transform", "translate(.5 .5)").
		add(seq.render(diagramPadding, diagramPadding+up, width))

	// The +1 matches railroad-diagrams@1.0.0's Diagram width fudge.
	totalWidth := width + 1 + diagramPadding*2
	totalHeight := up + seq.Down() + diagramPadding*2

	svg := el("svg").
		attr("xmlns", "http://www.w3.org/2000/svg").
		attrf("width", totalWidth).
		attrf("height", totalHeight).
		attr("viewBox", "0 0 "+formatNum(totalWidth)+" "+formatNum(totalHeight)).
		attr("class", "railroad-diagram").
		add(styleElement(), g)

	return svg.String()
}
