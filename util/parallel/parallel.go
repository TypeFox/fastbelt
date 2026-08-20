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

// maxWorkers caches runtime.GOMAXPROCS(0), which takes the scheduler lock on
// every call and shows up as contention when ForEach is called frequently.
// GOMAXPROCS changes at runtime are rare enough to ignore here.
var maxWorkers = runtime.GOMAXPROCS(0)

// ForEach calls action once for each element in the given slice,
// distributing the calls across the available CPU cores, while trying to
// limit the number of goroutines to a minimum to reduce scheduling overhead.
//
// The first argument to action is the element, and the second argument is the
// zero-based index of the element. This function blocks until every action
// call has returned.
func ForEach[T any](elements []T, action func(T, int)) {
	ForEachWithSetup(elements, func() struct{} { return struct{}{} }, func(_ struct{}, element T, i int) {
		action(element, i)
	})
}

// ForEachWithSetup behaves like [ForEach], but additionally calls setup once
// per worker goroutine and passes its result to every action call performed by
// that worker. Use it to give each worker goroutine-local state.
func ForEachWithSetup[T, S any](elements []T, setup func() S, action func(S, T, int)) {
	total := len(elements)
	if total == 0 {
		return
	}
	workers := min(maxWorkers, total)
	var next atomic.Int64
	// Store -1 so that the first Add(1) returns 0, the index of the first element.
	next.Store(-1)
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			state := setup()
			for {
				i := int(next.Add(1))
				if i >= total {
					return
				}
				action(state, elements[i], i)
			}
		})
	}
	wg.Wait()
}
