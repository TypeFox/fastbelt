// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parallel

import (
	"iter"
	"runtime"
	"slices"
	"sync"
	"sync/atomic"
)

// ForEachIter calls action once for each element in seq, distributing the
// calls across at most [runtime.GOMAXPROCS](0) goroutines instead of one
// goroutine per element.
//
// The second argument to action is the zero-based index of the element. seq
// is fully consumed before any action call begins. This function blocks
// until every action call has returned.
func ForEachIter[T any](seq iter.Seq[T], action func(T, int)) {
	ForEach(slices.Collect(seq), action)
}

// ForEach calls action once for each element in the given slice,
// distributing the calls across at most [runtime.GOMAXPROCS](0) goroutines
// instead of one goroutine per element.
//
// The second argument to action is the zero-based index of the element.
// This function blocks until every action call has returned.
func ForEach[T any](elements []T, action func(T, int)) {
	total := len(elements)
	if total == 0 {
		return
	}
	workers := min(runtime.GOMAXPROCS(0), total)
	var next atomic.Int64
	// Store -1 so that the first Add(1) returns 0, the index of the first element.
	next.Store(-1)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for {
				i := int(next.Add(1))
				if i >= total {
					return
				}
				action(elements[i], i)
			}
		})
	}
	wg.Wait()
}
