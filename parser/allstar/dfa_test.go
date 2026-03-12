// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeDummyState(num int) *ATNState {
	return &ATNState{StateNumber: num}
}

func TestATNConfigSet_Add_Dedup(t *testing.T) {
	s := NewATNConfigSet()
	state := makeDummyState(5)
	c := &ATNConfig{State: state, Alt: 0, Stack: []*ATNState{}}
	s.Add(c)
	s.Add(c) // duplicate
	assert.Equal(t, 1, s.Len())
}

func TestATNConfigSet_Add_DifferentAlt(t *testing.T) {
	s := NewATNConfigSet()
	state := makeDummyState(5)
	s.Add(&ATNConfig{State: state, Alt: 0, Stack: []*ATNState{}})
	s.Add(&ATNConfig{State: state, Alt: 1, Stack: []*ATNState{}})
	assert.Equal(t, 2, s.Len())
}

func TestATNConfigSet_Key_Consistency(t *testing.T) {
	s := NewATNConfigSet()
	state := makeDummyState(3)
	s.Add(&ATNConfig{State: state, Alt: 0, Stack: []*ATNState{}})
	key1 := s.Key()
	key2 := s.Key()
	assert.Equal(t, key1, key2)
}

func TestATNConfigSet_Finalize(t *testing.T) {
	s := NewATNConfigSet()
	s.Add(&ATNConfig{State: makeDummyState(1), Alt: 0, Stack: []*ATNState{}})
	s.Finalize()
	assert.Equal(t, 1, s.Len(), "Finalize should not change length")
}

func TestATNConfigSet_Alts(t *testing.T) {
	s := NewATNConfigSet()
	s.Add(&ATNConfig{State: makeDummyState(1), Alt: 2, Stack: []*ATNState{}})
	s.Add(&ATNConfig{State: makeDummyState(2), Alt: 3, Stack: []*ATNState{}})
	alts := s.Alts()
	assert.Equal(t, []int{2, 3}, alts)
}

func TestATNConfigKey_WithAlt(t *testing.T) {
	state := makeDummyState(7)
	c := &ATNConfig{State: state, Alt: 2, Stack: []*ATNState{}}
	key := atnConfigKey(c, true)
	assert.Contains(t, key, "a2")
	assert.Contains(t, key, "s7")
}

func TestATNConfigKey_WithoutAlt(t *testing.T) {
	state := makeDummyState(7)
	c := &ATNConfig{State: state, Alt: 2, Stack: []*ATNState{}}
	key := atnConfigKey(c, false)
	assert.NotContains(t, key, "a2")
	assert.Contains(t, key, "s7")
}

func TestATNConfigKey_WithStack(t *testing.T) {
	state := makeDummyState(3)
	stackState := makeDummyState(9)
	c := &ATNConfig{State: state, Alt: 0, Stack: []*ATNState{stackState}}
	key := atnConfigKey(c, true)
	assert.Contains(t, key, "s3")
	assert.Contains(t, key, "9")
}
