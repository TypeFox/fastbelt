// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "math"

// OneOrMore renders Item on the through-path with a Repeat branch looping
// back underneath it. Repeat defaults to Skip{} if nil, matching
// railroad-diagrams@1.0.0's OneOrMore(item, rep = Skip()).
type OneOrMore struct {
	Item   Node
	Repeat Node
}

func (o *OneOrMore) repeat() Node {
	if o.Repeat != nil {
		return o.Repeat
	}
	return Skip{}
}

func (o *OneOrMore) Width() float64 {
	return math.Max(o.Item.Width(), o.repeat().Width()) + arcRadius*2
}

func (o *OneOrMore) Up() float64 { return o.Item.Up() }

func (o *OneOrMore) Down() float64 {
	return math.Max(arcRadius*2, o.Item.Down()+verticalSeparation+o.repeat().Up()+o.repeat().Down())
}

func (o *OneOrMore) needsSpace() bool { return true }

func (o *OneOrMore) render(x, y, width float64) *svgElement {
	g := el("g")
	gapL, gapR := determineGaps(width, o.Width())
	g.add(newPath(x, y).h(gapL).element())
	g.add(newPath(x+gapL+o.Width(), y).h(gapR).element())
	x += gapL

	g.add(newPath(x, y).right(arcRadius).element())
	g.add(o.Item.render(x+arcRadius, y, o.Width()-arcRadius*2))
	g.add(newPath(x+o.Width()-arcRadius, y).right(arcRadius).element())

	distanceFromY := math.Max(arcRadius*2, o.Item.Down()+verticalSeparation+o.repeat().Up())
	g.add(newPath(x+arcRadius, y).arc("nw").down(distanceFromY - arcRadius*2).arc("ws").element())
	g.add(o.repeat().render(x+arcRadius, y+distanceFromY, o.Width()-arcRadius*2))
	g.add(newPath(x+o.Width()-arcRadius, y+distanceFromY).arc("se").up(distanceFromY - arcRadius*2).arc("en").element())

	return g
}

// ZeroOrMore is Optional(OneOrMore(item, Skip{})).
func ZeroOrMore(item Node) Node {
	return Optional(&OneOrMore{Item: item})
}
