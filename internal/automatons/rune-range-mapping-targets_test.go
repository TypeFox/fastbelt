package automatons

import (
	"testing"
)

func TestNFATargets_OneRange(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	Expect(len(targets.Ranges)).ToEqual(3)
	Expect(targets.GetRuneValues('a')).ToContain(1)
}

func TestNFATargets_TwoSeparatedRunes(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneValues('a', Targets{1})
	targets.AddRuneValues('c', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('c')).ToContain(2)
}

func TestNFATargets_TwoIdenticalRunes(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneValues('a', Targets{1})
	targets.AddRuneValues('a', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(3)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('a')).ToContain(2)
}

func TestNFATargets_TwoConnectedRunes(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneValues('a', Targets{1})
	targets.AddRuneValues('b', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(4)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('b')).ToContain(2)
}

func TestNFATargets_TwoSeparatedRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	targets.AddRuneRangeValues('e', 'g', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('g')).ToContain(2)
}

func TestNFATargets_TwoConnectedRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	targets.AddRuneRangeValues('d', 'f', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(4)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('d')).ToContain(2)
}

func TestNFATargets_TwoConnectedRangesReversed(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('d', 'f', Targets{2})
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	Expect(len(targets.Ranges)).ToEqual(4)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('d')).ToContain(2)
}

func TestNFATargets_TwoOverlappedRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	targets.AddRuneRangeValues('b', 'd', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('b')).ToContain(1)
	Expect(targets.GetRuneValues('b')).ToContain(2)
	Expect(targets.GetRuneValues('d')).ToContain(2)
}

func TestNFATargets_TwoOverlappedRangesReversed(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('b', 'd', Targets{2})
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneValues('a')).ToContain(1)
	Expect(targets.GetRuneValues('b')).ToContain(1)
	Expect(targets.GetRuneValues('b')).ToContain(2)
	Expect(targets.GetRuneValues('d')).ToContain(2)
}

func TestNFATargets_TwoContainingRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('e', 'f', Targets{1})
	targets.AddRuneRangeValues('a', 'h', Targets{2})
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneValues('a')).ToContain(2)
	Expect(targets.GetRuneValues('e')).ToContain(1)
	Expect(targets.GetRuneValues('e')).ToContain(2)
	Expect(targets.GetRuneValues('h')).ToContain(2)
}
