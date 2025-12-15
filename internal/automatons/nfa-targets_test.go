package automatons

import (
	"testing"
)

func TestNFATargets_OneRange(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('a', 'c', 1)
	Expect(len(targets.Ranges)).ToEqual(3)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
}

func TestNFATargets_TwoSeparatedRunes(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneTarget('a', 1)
	targets.AddRuneTarget('c', 2)
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('c')).ToContain(2)
}

func TestNFATargets_TwoIdenticalRunes(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneTarget('a', 1)
	targets.AddRuneTarget('a', 2)
	Expect(len(targets.Ranges)).ToEqual(3)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('a')).ToContain(2)
}

func TestNFATargets_TwoConnectedRunes(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneTarget('a', 1)
	targets.AddRuneTarget('b', 2)
	Expect(len(targets.Ranges)).ToEqual(4)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('b')).ToContain(2)
}

func TestNFATargets_TwoSeparatedRanges(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('a', 'c', 1)
	targets.AddRuneRangeTarget('e', 'g', 2)
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('g')).ToContain(2)
}

func TestNFATargets_TwoConnectedRanges(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('a', 'c', 1)
	targets.AddRuneRangeTarget('d', 'f', 2)
	Expect(len(targets.Ranges)).ToEqual(4)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('d')).ToContain(2)
}

func TestNFATargets_TwoConnectedRangesReversed(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('d', 'f', 2)
	targets.AddRuneRangeTarget('a', 'c', 1)
	Expect(len(targets.Ranges)).ToEqual(4)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('d')).ToContain(2)
}

func TestNFATargets_TwoOverlappedRanges(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('a', 'c', 1)
	targets.AddRuneRangeTarget('b', 'd', 2)
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('b')).ToContain(1)
	Expect(targets.GetRuneTargets('b')).ToContain(2)
	Expect(targets.GetRuneTargets('d')).ToContain(2)
}

func TestNFATargets_TwoOverlappedRangesReversed(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('b', 'd', 2)
	targets.AddRuneRangeTarget('a', 'c', 1)
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneTargets('a')).ToContain(1)
	Expect(targets.GetRuneTargets('b')).ToContain(1)
	Expect(targets.GetRuneTargets('b')).ToContain(2)
	Expect(targets.GetRuneTargets('d')).ToContain(2)
}

func TestNFATargets(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('e', 'f', 1)
	targets.AddRuneRangeTarget('a', 'h', 2)
	Expect(len(targets.Ranges)).ToEqual(5)
	Expect(targets.GetRuneTargets('a')).ToContain(2)
	Expect(targets.GetRuneTargets('e')).ToContain(1)
	Expect(targets.GetRuneTargets('e')).ToContain(2)
	Expect(targets.GetRuneTargets('h')).ToContain(2)
}
