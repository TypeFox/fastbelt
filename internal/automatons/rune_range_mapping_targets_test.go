package automatons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNFATargets_OneRange(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	assert.Equal(t, 3, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
}

func TestNFATargets_TwoSeparatedRunes(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneValues('a', Targets{1})
	targets.AddRuneValues('c', Targets{2})
	assert.Equal(t, 5, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('c'), &Targets{2})
}

func TestNFATargets_TwoIdenticalRunes(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneValues('a', Targets{1})
	targets.AddRuneValues('a', Targets{2})
	assert.Equal(t, 3, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1, 2})
}

func TestNFATargets_TwoConnectedRunes(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneValues('a', Targets{1})
	targets.AddRuneValues('b', Targets{2})
	assert.Equal(t, 4, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('b'), &Targets{2})
}

func TestNFATargets_TwoSeparatedRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	targets.AddRuneRangeValues('e', 'g', Targets{2})
	assert.Equal(t, 5, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('g'), &Targets{2})
}

func TestNFATargets_TwoConnectedRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	targets.AddRuneRangeValues('d', 'f', Targets{2})
	assert.Equal(t, 4, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('d'), &Targets{2})
}

func TestNFATargets_TwoConnectedRangesReversed(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('d', 'f', Targets{2})
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	assert.Equal(t, 4, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('d'), &Targets{2})
}

func TestNFATargets_TwoOverlappedRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	targets.AddRuneRangeValues('b', 'd', Targets{2})
	assert.Equal(t, 5, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('b'), &Targets{1, 2})
	assert.EqualValues(t, targets.GetRuneValues('d'), &Targets{2})
}

func TestNFATargets_TwoOverlappedRangesReversed(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('b', 'd', Targets{2})
	targets.AddRuneRangeValues('a', 'c', Targets{1})
	assert.Equal(t, 5, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{1})
	assert.EqualValues(t, targets.GetRuneValues('b'), &Targets{2, 1})
	assert.EqualValues(t, targets.GetRuneValues('d'), &Targets{2})
}

func TestNFATargets_TwoContainingRanges(t *testing.T) {
	targets := NewRuneRangeTargets()
	targets.AddRuneRangeValues('e', 'f', Targets{1})
	targets.AddRuneRangeValues('a', 'h', Targets{2})
	assert.Equal(t, 5, len(targets.Ranges))
	assert.EqualValues(t, targets.GetRuneValues('a'), &Targets{2})
	assert.EqualValues(t, targets.GetRuneValues('e'), &Targets{1, 2})
	assert.EqualValues(t, targets.GetRuneValues('h'), &Targets{2})
}
