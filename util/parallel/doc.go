// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package parallel provides utilities for parallel iteration over slices and
// other collections.
//
// The naive way to parallelize work in Go is to spawn one goroutine per
// element. Goroutines are cheap, but not free: each one costs stack
// allocation and scheduler bookkeeping, and running far more goroutines
// than there are CPU cores adds overhead without adding throughput.
// For workloads that iterate over large collections — such as the
// document build pipeline in [typefox.dev/fastbelt/workspace], which
// parses, links, and validates every document in a workspace — this
// overhead is measurable.
//
// This package therefore bounds concurrency instead of spawning
// unboundedly: [ForEach] and [ForEachIter] distribute the elements across
// at most [runtime.GOMAXPROCS](0) worker goroutines, which pull elements
// from a shared index counter. This caps the goroutine count at the
// number of usable CPU cores while still balancing the load when
// individual elements take different amounts of time to process.
package parallel
