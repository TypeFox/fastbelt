// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parallel

import (
	"iter"
	"runtime"
	"slices"
	"sync"
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
	if len(elements) == 0 {
		return
	}
	workers := min(runtime.GOMAXPROCS(0), len(elements))
	// Note: we use a channel to distribute indices to workers
	// This ensures that there is no idle worker while there are still elements to process
	indices := make(chan int, workers)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for index := range indices {
				action(elements[index], index)
			}
		})
	}
	for i := range elements {
		indices <- i
	}
	close(indices)
	wg.Wait()
}
