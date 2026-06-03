// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package extiter provides utilities for working with [iter.Seq] sequences.
//
// Most functions fall into two groups. Sequence builders and lazy transforms
// such as [Of], [Empty], [Map], [Filter], [FlatMap], and [Concat] return a
// new sequence and defer work until a consumer iterates. Eager operations
// such as [Count], [IsEmpty], [Every], [Find], and [Reduce] walk the input
// sequence immediately and return a value.
//
// Fastbelt uses extiter throughout symbol lookup and linking: [Empty] and
// [IsEmpty] represent and test empty symbol iterators, [Map] and [Filter]
// reshape document and scope collections, [FlatMap] flattens nested symbol
// containers, and [Concat] chains local and outer scope elements.
package extiter
