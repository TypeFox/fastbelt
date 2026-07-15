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
		// Note: 7919 is very likely to be more than the available amount of CPU cores
		// So this test will exercise the chunking feature of ForEach.
		// Also, it is a prime number, so it will not be divisible by any chunk size.
		elements := make([]int, 7919)
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

	t.Run("single element", func(t *testing.T) {
		called := false
		ForEach([]int{42}, func(value int, index int) {
			called = true
			if value != 42 || index != 0 {
				t.Errorf("expected value 42 at index 0, got value %d at index %d", value, index)
			}
		})
		if !called {
			t.Error("expected action to be called for single element")
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
