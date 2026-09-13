// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "math"

// Sequence lays out its items left to right on a shared centerline.
type Sequence struct {
	Items []Node
}

func (s *Sequence) Width() float64 {
	total := 0.0
	for _, it := range s.Items {
		total += it.Width()
		if it.needsSpace() {
			total += 20
		}
	}
	return total
}

func (s *Sequence) Up() float64 {
	up := 0.0
	for _, it := range s.Items {
		up = math.Max(up, it.Up())
	}
	return up
}

func (s *Sequence) Down() float64 {
	down := 0.0
	for _, it := range s.Items {
		down = math.Max(down, it.Down())
	}
	return down
}

func (s *Sequence) needsSpace() bool { return false }

func (s *Sequence) render(x, y, width float64) *svgElement {
	g := el("g")
	gapL, gapR := determineGaps(width, s.Width())
	g.add(newPath(x, y).h(gapL).element())
	g.add(newPath(x+gapL+s.Width(), y).h(gapR).element())
	x += gapL

	for _, it := range s.Items {
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
	return g
}
