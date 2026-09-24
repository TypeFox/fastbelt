// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import "testing"

func TestSequenceGeometry(t *testing.T) {
	a := &Terminal{Text: "if"}   // width = 2*8+20 = 36, needsSpace
	b := &Terminal{Text: "then"} // width = 4*8+20 = 52, needsSpace
	seq := &Sequence{Items: []Node{a, b}}

	wantWidth := a.Width() + 20 + b.Width() + 20
	if got := seq.Width(); got != wantWidth {
		t.Errorf("Width() = %v, want %v", got, wantWidth)
	}
	if got := seq.Up(); got != 11 {
		t.Errorf("Up() = %v, want 11", got)
	}
	if got := seq.Down(); got != 11 {
		t.Errorf("Down() = %v, want 11", got)
	}
}

func TestOptionalGeometry(t *testing.T) {
	item := &Terminal{Text: "x"} // width = 1*8+20 = 28, up=down=11
	opt := Optional(item)

	wantWidth := item.Width() + arcRadius*4
	if got := opt.Width(); got != wantWidth {
		t.Errorf("Width() = %v, want %v", got, wantWidth)
	}
	// Skip is index 0 (< normal=1): contributes max(AR, 0+0+VS) = max(10, 8) = 10.
	// item is the normal branch: it also contributes its own up-clearance,
	// max(AR, item.Up()) = max(10, 11) = 11. Total: 10 + 11 = 21.
	if got := opt.Up(); got != 21 {
		t.Errorf("Up() = %v, want 21", got)
	}
	// item is the normal branch: down contribution is max(AR, item.Down()) = max(10, 11) = 11.
	if got := opt.Down(); got != 11 {
		t.Errorf("Down() = %v, want 11", got)
	}
}
