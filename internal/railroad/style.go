// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

// CSS classes shared between the node renderers above and the <style> block
// emitted by Diagram.SVG.
const (
	classTerminal              = "rr-terminal"
	classNonTerminal           = "rr-nonterminal"
	classNonTerminalUnresolved = "rr-nonterminal-unresolved"
	classBoxText               = "rr-text"
)

// diagramCSS colors every shape itself rather than relying on the SVG's
// background, since the diagram is embedded as a flat <img> (a data: URI)
// with no way to react to the viewer's light/dark theme: connector tracks
// use a mid-tone gray that reads on either background, and terminal/
// non-terminal boxes are always filled with a fixed color and white text so
// each box supplies its own contrast.
const diagramCSS = `
.railroad-diagram path {
	fill: none;
	stroke: #8f8f8f;
	stroke-width: 2;
}
.railroad-diagram .` + classTerminal + ` {
	fill: #3b6ea5;
	stroke: #3b6ea5;
	stroke-width: 2;
}
.railroad-diagram .` + classNonTerminal + ` {
	fill: #5a5a5a;
	stroke: #5a5a5a;
	stroke-width: 2;
}
.railroad-diagram .` + classNonTerminalUnresolved + ` {
	fill: #7a7a7a;
	stroke: #b0b0b0;
	stroke-width: 2;
	stroke-dasharray: 4 2;
}
.railroad-diagram .` + classBoxText + ` {
	fill: #ffffff;
	font: 12px monospace;
	text-anchor: middle;
	dominant-baseline: middle;
}
`

func styleElement() *svgElement {
	return el("style").setText(diagramCSS)
}
