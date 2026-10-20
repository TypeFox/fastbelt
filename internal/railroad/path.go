// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "strconv"

// pathBuilder builds the "d" attribute of an SVG <path> element via
// chainable relative move/line/arc commands, mirroring
// railroad-diagrams@1.0.0's Path helper.
type pathBuilder struct {
	d []byte
}

func newPath(x, y float64) *pathBuilder {
	p := &pathBuilder{d: make([]byte, 0, 64)}
	p.d = append(p.d, 'M')
	p.d = appendNum(p.d, x)
	p.d = append(p.d, ' ')
	p.d = appendNum(p.d, y)
	return p
}

func (p *pathBuilder) h(val float64) *pathBuilder {
	p.d = append(p.d, 'h')
	p.d = appendNum(p.d, val)
	return p
}

func (p *pathBuilder) v(val float64) *pathBuilder {
	p.d = append(p.d, 'v')
	p.d = appendNum(p.d, val)
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
	cw := byte('0')
	switch sweep {
	case "ne", "es", "sw", "wn":
		cw = '1'
	}
	p.d = append(p.d, 'a')
	p.d = appendNum(p.d, arcRadius)
	p.d = append(p.d, ' ')
	p.d = appendNum(p.d, arcRadius)
	p.d = append(p.d, " 0 0 "...)
	p.d = append(p.d, cw, ' ')
	p.d = appendNum(p.d, x)
	p.d = append(p.d, ' ')
	p.d = appendNum(p.d, y)
	return p
}

func (p *pathBuilder) element() *svgElement {
	return el("path").attr("d", string(p.d))
}

func appendNum(b []byte, v float64) []byte {
	return strconv.AppendFloat(b, v, 'f', -1, 64)
}
