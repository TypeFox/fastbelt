// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "math"

// OneOrMore renders Item on the through-path with an empty branch looping
// back underneath it, matching railroad-diagrams@1.0.0's OneOrMore(item).
type OneOrMore struct {
	Item Node
}

func (o *OneOrMore) Width() float64 { return o.Item.Width() + arcRadius*2 }

func (o *OneOrMore) Up() float64 { return o.Item.Up() }

func (o *OneOrMore) Down() float64 {
	return math.Max(arcRadius*2, o.Item.Down()+verticalSeparation)
}

func (o *OneOrMore) needsSpace() bool { return true }

func (o *OneOrMore) render(x, y, width float64) *svgElement {
	g := el("g")
	w := o.Width()
	gap := (width - w) / 2
	g.add(newPath(x, y).h(gap).element())
	g.add(newPath(x+gap+w, y).h(gap).element())
	x += gap

	g.add(newPath(x, y).h(arcRadius).element())
	g.add(o.Item.render(x+arcRadius, y, w-arcRadius*2))
	g.add(newPath(x+w-arcRadius, y).h(arcRadius).element())

	distanceFromY := o.Down()
	g.add(newPath(x+arcRadius, y).arc("nw").down(distanceFromY - arcRadius*2).arc("ws").element())
	g.add(Skip{}.render(x+arcRadius, y+distanceFromY, w-arcRadius*2))
	g.add(newPath(x+w-arcRadius, y+distanceFromY).arc("se").up(distanceFromY - arcRadius*2).arc("en").element())

	return g
}

// ZeroOrMore is Optional(OneOrMore(item)).
func ZeroOrMore(item Node) Node {
	return Optional(&OneOrMore{Item: item})
}
