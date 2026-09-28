// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "math"

// Choice stacks its items vertically, with the item at index Normal drawn
// on the main centerline and the rest branching above (indices < Normal) or
// below (indices > Normal) via quarter-circle arcs. Optional and ZeroOrMore
// are built on top of this. Items must not be empty.
type Choice struct {
	Normal int
	Items  []Node
}

func (c *Choice) Width() float64 {
	w := 0.0
	for _, it := range c.Items {
		w = math.Max(w, it.Width())
	}
	return w + arcRadius*4
}

func (c *Choice) Up() float64 {
	up := 0.0
	for i, it := range c.Items {
		switch {
		case i < c.Normal:
			up += math.Max(arcRadius, it.Up()+it.Down()+verticalSeparation)
		case i == c.Normal:
			up += math.Max(arcRadius, it.Up())
		}
	}
	return up
}

func (c *Choice) Down() float64 {
	down := 0.0
	for i, it := range c.Items {
		switch {
		case i == c.Normal:
			down += math.Max(arcRadius, it.Down())
		case i > c.Normal:
			down += math.Max(arcRadius, verticalSeparation+it.Up()+it.Down())
		}
	}
	return down
}

func (c *Choice) needsSpace() bool { return false }

func (c *Choice) render(x, y, width float64) *svgElement {
	g := el("g")
	w := c.Width()
	gap := (width - w) / 2
	g.add(newPath(x, y).h(gap).element())
	g.add(newPath(x+gap+w, y).h(gap).element())
	x += gap

	last := len(c.Items) - 1
	innerWidth := w - arcRadius*4

	// Items above the main centerline, curving up and back down.
	distanceFromY := 0.0
	for i := c.Normal - 1; i >= 0; i-- {
		item := c.Items[i]
		if i == c.Normal-1 {
			distanceFromY = math.Max(arcRadius*2, c.Items[i+1].Up()+verticalSeparation+item.Down())
		}
		g.add(newPath(x, y).arc("se").up(distanceFromY - arcRadius*2).arc("wn").element())
		g.add(item.render(x+arcRadius*2, y-distanceFromY, innerWidth))
		g.add(newPath(x+arcRadius*2+innerWidth, y-distanceFromY).arc("ne").down(distanceFromY - arcRadius*2).arc("ws").element())
		prevDown := 0.0
		if i != 0 {
			prevDown = c.Items[i-1].Down()
		}
		distanceFromY += math.Max(arcRadius, item.Up()+verticalSeparation+prevDown)
	}

	// The straight-line path through the normal item.
	g.add(newPath(x, y).h(arcRadius * 2).element())
	g.add(c.Items[c.Normal].render(x+arcRadius*2, y, innerWidth))
	g.add(newPath(x+arcRadius*2+innerWidth, y).h(arcRadius * 2).element())

	// Items below the main centerline, curving down and back up.
	distanceFromY = 0.0
	for i := c.Normal + 1; i <= last; i++ {
		item := c.Items[i]
		if i == c.Normal+1 {
			distanceFromY = math.Max(arcRadius*2, c.Items[i-1].Down()+verticalSeparation+item.Up())
		}
		g.add(newPath(x, y).arc("ne").down(distanceFromY - arcRadius*2).arc("ws").element())
		g.add(item.render(x+arcRadius*2, y+distanceFromY, innerWidth))
		g.add(newPath(x+arcRadius*2+innerWidth, y+distanceFromY).arc("se").up(distanceFromY - arcRadius*2).arc("wn").element())
		nextUp := 0.0
		if i != last {
			nextUp = c.Items[i+1].Up()
		}
		distanceFromY += math.Max(arcRadius, item.Down()+verticalSeparation+nextUp)
	}

	return g
}

// Optional renders item as the through-path with a Skip branch curving
// above it, matching railroad-diagrams@1.0.0's default Optional(item).
func Optional(item Node) Node {
	return &Choice{Normal: 1, Items: []Node{Skip{}, item}}
}
