// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parallel

import (
	"math"
	"sync/atomic"
)

// RunningAverage is a concurrency-safe exponential moving average of a
// float64 ratio. It starts out returning a caller-supplied default until
// [RunningAverage.Update] has been called at least once.
type RunningAverage struct {
	def  float64
	bits atomic.Uint64
}

// NewRunningAverage creates a RunningAverage that returns def until the first
// sample is recorded via [RunningAverage.Update].
func NewRunningAverage(def float64) *RunningAverage {
	return &RunningAverage{def: def}
}

// Value returns the current average, or the default passed to
// [NewRunningAverage] if no sample has been recorded yet.
func (r *RunningAverage) Value() float64 {
	bits := r.bits.Load()
	if bits == 0 {
		return r.def
	}
	return math.Float64frombits(bits)
}

// Update folds sample into the running average, weighting the existing
// average at 90% and the new sample at 10%.
func (r *RunningAverage) Update(sample float64) {
	next := r.Value()*0.9 + sample*0.1
	r.bits.Store(math.Float64bits(next))
}

// Capacity returns the ideal slice capacity for a new slice relative
// to the given size, based on the current average ratio.
func (r *RunningAverage) Capacity(size int) int {
	return int(float64(size) * r.Value() * 1.1)
}
