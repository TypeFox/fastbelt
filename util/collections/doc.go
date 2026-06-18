// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package collections provides generic collection data structures shared across
// Fastbelt. Use [MultiMap] when a key may map to many values and callers need
// grouped lookup, ordered values per key, or a total element count. Use [Set]
// for membership tracking when only the presence of distinct values matters.
package collections
