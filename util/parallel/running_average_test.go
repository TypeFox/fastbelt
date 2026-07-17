// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parallel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunningAverage(t *testing.T) {
	t.Run("returns default value before first update", func(t *testing.T) {
		def := 42.0
		ra := NewRunningAverage(def)
		assert.Equal(t, def, ra.Value(), "expected default value before first update")
	})

	t.Run("updates average after multiple runs", func(t *testing.T) {
		ra := NewRunningAverage(10)
		ra.Update(20)
		ra.Update(30)
		ra.Update(40)
		for range 20 {
			ra.Update(50)
		}
		assert.Greater(t, ra.Value(), 45.0, "average should be close to target value after multiple updates")
	})
}
