// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "strings"

// pathBuilder builds the "d" attribute of an SVG <path> element via
// chainable relative move/line/arc commands, mirroring
// railroad-diagrams@1.0.0's Path helper.
type pathBuilder struct {
	d strings.Builder
}

func newPath(x, y float64) *pathBuilder {
	p := &pathBuilder{}
	p.d.WriteString("M" + formatNum(x) + " " + formatNum(y))
	return p
}

func (p *pathBuilder) h(val float64) *pathBuilder {
	p.d.WriteString("h" + formatNum(val))
	return p
}

func (p *pathBuilder) right(val float64) *pathBuilder { return p.h(val) }

func (p *pathBuilder) v(val float64) *pathBuilder {
	p.d.WriteString("v" + formatNum(val))
	return p
}

func (p *pathBuilder) down(val float64) *pathBuilder { return p.v(val) }
func (p *pathBuilder) up(val float64) *pathBuilder   { return p.v(-val) }

// arc draws a quarter circle of radius arcRadius. sweep is a two-letter
// compass code (e.g. "ne", "es", "sw", "wn") describing which quadrant the
// curve bends through; this exactly mirrors Path.prototype.arc from
// railroad-diagrams@1.0.0, including its clockwise-flag table.
func (p *pathBuilder) arc(sweep string) *pathBuilder {
	x, y := arcRadius, arcRadius
	if sweep[0] == 'e' || sweep[1] == 'w' {
		x = -x
	}
	if sweep[0] == 's' || sweep[1] == 'n' {
		y = -y
	}
	cw := 0
	switch sweep {
	case "ne", "es", "sw", "wn":
		cw = 1
	}
	p.d.WriteString("a" + formatNum(arcRadius) + " " + formatNum(arcRadius) + " 0 0 " + formatNum(float64(cw)) + " " + formatNum(x) + " " + formatNum(y))
	return p
}

func (p *pathBuilder) element() *svgElement {
	return el("path").attr("d", p.d.String())
}
