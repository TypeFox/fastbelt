// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package parallel provides utilities for parallel iteration over slices and
// sequences, plus [RunningAverage], a concurrency-safe exponential moving
// average. [ForEach], [ForEachIter], and [ForEachWithSetup] run the work on at
// most [runtime.GOMAXPROCS] worker goroutines instead of one goroutine per
// element. This keeps scheduling overhead low for large collections, such as
// the documents processed by the build pipeline in
// [typefox.dev/fastbelt/workspace], while still balancing the load when
// elements take uneven time.
package parallel
