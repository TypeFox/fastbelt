package automatons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNFATargets_Simple(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('a', 'c', 1)
	assert.Equal(t, 1, len(targets.GetRuneTargets('a')))
}

func TestNFATargets(t *testing.T) {
	targets := NewNFAMutableTargets()
	targets.AddRuneRangeTarget('e', 'f', 1)
	targets.AddRuneRangeTarget('a', 'h', 2)
	assert.Equal(t, 1, len(targets.GetRuneTargets('a')))
	assert.Equal(t, 2, len(targets.GetRuneTargets('e')))
	assert.Equal(t, 1, len(targets.GetRuneTargets('h')))
}
