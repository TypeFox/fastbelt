// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parallel

import (
	"sync/atomic"
	"testing"
)

func TestParallelForEach(t *testing.T) {
	t.Run("visits every element exactly once with correct indices", func(t *testing.T) {
		elements := make([]int, 500)
		for i := range elements {
			elements[i] = i
		}
		seen := make([]int32, len(elements))
		ForEach(elements, func(value int, index int) {
			atomic.AddInt32(&seen[index], 1)
			if value != index {
				t.Errorf("expected value %d at index %d, got %d", index, index, value)
			}
		})
		for i, count := range seen {
			if count != 1 {
				t.Errorf("expected index %d to be visited once, got %d", i, count)
			}
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		called := false
		ForEach([]int{}, func(int, int) { called = true })
		if called {
			t.Error("expected action not to be called for empty sequence")
		}
	})
}
